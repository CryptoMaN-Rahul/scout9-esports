package grid

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ListFiles lists available files for a series
func (c *Client) ListFiles(ctx context.Context, seriesID string) ([]FileInfo, error) {
	url := fmt.Sprintf("%s/file-download/list/%s", FileDownloadURL, seriesID)

	resp, err := c.doFileDownloadRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("list files request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("list files failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Files []FileInfo `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Files, nil
}

// DownloadEvents downloads and parses the events JSONL file for a series
// Returns EventWrapper slice with the correct GRID event structure
func (c *Client) DownloadEvents(ctx context.Context, seriesID string) ([]EventWrapper, error) {
	cacheKey := fmt.Sprintf("events:%s", seriesID)

	// Check cache
	var events []EventWrapper
	if c.getCached(ctx, cacheKey, &events) {
		return events, nil
	}

	// First list files to get the events URL
	files, err := c.ListFiles(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}

	// Find the events file (could be events-grid or events-grid-compressed)
	var eventsURL string
	for _, f := range files {
		if (f.ID == "events-grid" || f.ID == "events-grid-compressed") && f.Status == "ready" {
			eventsURL = f.FullURL
			break
		}
	}

	if eventsURL == "" {
		return nil, fmt.Errorf("events file not available for series %s", seriesID)
	}

	// Download the file
	data, err := c.downloadFile(ctx, eventsURL)
	if err != nil {
		return nil, fmt.Errorf("download events: %w", err)
	}

	// The file is a ZIP, need to decompress
	events, err = parseEventsZip(data)
	if err != nil {
		return nil, fmt.Errorf("parse events: %w", err)
	}

	// Cache result
	c.setCache(ctx, cacheKey, events, 24*time.Hour)

	return events, nil
}

// parseEventsZip decompresses and parses the events JSONL file
func parseEventsZip(data []byte) ([]EventWrapper, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	var wrappers []EventWrapper

	for _, file := range reader.File {
		if !strings.HasSuffix(file.Name, ".jsonl") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open file in zip: %w", err)
		}

		scanner := bufio.NewScanner(rc)
		// Increase buffer size for large lines - GRID events can have full state snapshots
		// that exceed 10MB per line
		buf := make([]byte, 0, 256*1024)
		scanner.Buffer(buf, 20*1024*1024) // 20MB max line size

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var wrapper EventWrapper
			if err := json.Unmarshal(line, &wrapper); err != nil {
				// Skip malformed lines
				continue
			}
			wrappers = append(wrappers, wrapper)
		}

		rc.Close()

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scan jsonl: %w", err)
		}
	}

	return wrappers, nil
}

// DownloadEndState downloads the end state JSON file for a series
func (c *Client) DownloadEndState(ctx context.Context, seriesID string) (*SeriesState, error) {
	cacheKey := fmt.Sprintf("endstate:%s", seriesID)

	// Check cache
	var state SeriesState
	if c.getCached(ctx, cacheKey, &state) {
		return &state, nil
	}

	// First list files to get the end state URL
	files, err := c.ListFiles(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}

	// Find the end state file
	var stateURL string
	for _, f := range files {
		if f.ID == "state-grid" && f.Status == "ready" {
			stateURL = f.FullURL
			break
		}
	}

	if stateURL == "" {
		return nil, fmt.Errorf("end state file not available for series %s", seriesID)
	}

	// Download the file
	data, err := c.downloadFile(ctx, stateURL)
	if err != nil {
		return nil, fmt.Errorf("download end state: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parse end state: %w", err)
	}

	// Cache result
	c.setCache(ctx, cacheKey, state, 24*time.Hour)

	return &state, nil
}

// LoLEventData contains parsed LoL events with rich data
type LoLEventData struct {
	DragonKills    []DragonKillEvent
	BaronKills     []ObjectiveKillEvent
	HeraldKills    []ObjectiveKillEvent
	VoidGrubKills  []ObjectiveKillEvent
	TowerDestroys  []TowerDestroyEvent
	Kills          []KillEvent
	DraftActions   []DraftAction
	FirstBloodTime time.Time
	GameStartTime  time.Time
}

// VALEventData contains parsed VALORANT events with rich data
type VALEventData struct {
	RoundEnds []RoundEndEvent
	Plants    []PlantEvent
	Defuses   []DefuseEvent
	Kills     []VALKillEvent
	MapName   string
}

// ParseLoLEvents extracts LoL-specific events from EventWrappers
func ParseLoLEvents(wrappers []EventWrapper) (*LoLEventData, error) {
	data := &LoLEventData{
		DragonKills:   make([]DragonKillEvent, 0),
		BaronKills:    make([]ObjectiveKillEvent, 0),
		HeraldKills:   make([]ObjectiveKillEvent, 0),
		VoidGrubKills: make([]ObjectiveKillEvent, 0),
		TowerDestroys: make([]TowerDestroyEvent, 0),
		Kills:         make([]KillEvent, 0),
		DraftActions:  make([]DraftAction, 0),
	}

	// Regex for parsing tower IDs like "red-turret-mid-2"
	towerRegex := regexp.MustCompile(`(red|blue)-(?:turret|inhibitor)-(\w+)-?(\d*)`)

	isFirstKill := true
	var gameStartTime time.Time

	// First pass: find game start time
	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			if event.Action == "started" {
				targetType := ""
				if event.Target != nil {
					targetType = event.Target.Type
				}
				if targetType == "game" {
					gameStartTime = wrapper.OccurredAt
					data.GameStartTime = gameStartTime
					break
				}
			}
		}
		if !gameStartTime.IsZero() {
			break
		}
	}

	// Helper to extract in-game clock time from seriesState (more accurate than wall-clock diff)
	// Returns time in milliseconds
	getGameClockTime := func(event GridEvent) int {
		if event.SeriesState == nil {
			return 0
		}
		// Navigate: seriesState.games[0].clock.currentSeconds
		if games, ok := event.SeriesState["games"].([]interface{}); ok && len(games) > 0 {
			if game, ok := games[0].(map[string]interface{}); ok {
				if clock, ok := game["clock"].(map[string]interface{}); ok {
					if currentSeconds, ok := clock["currentSeconds"].(float64); ok {
						return int(currentSeconds * 1000) // Convert to milliseconds
					}
				}
			}
		}
		return 0
	}

	// Fallback: calculate game time from wall-clock difference
	calcGameTimeFromWallClock := func(eventTime time.Time) int {
		if gameStartTime.IsZero() {
			return 0
		}
		return int(eventTime.Sub(gameStartTime).Milliseconds())
	}

	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			actorType := ""
			if event.Actor != nil {
				actorType = event.Actor.Type
			}
			targetType := ""
			if event.Target != nil {
				targetType = event.Target.Type
			}

			// Get game time - prefer in-game clock, fallback to wall-clock calculation
			gameTime := getGameClockTime(event)
			if gameTime == 0 {
				gameTime = calcGameTimeFromWallClock(wrapper.OccurredAt)
			}

			switch {
			// Player kills
			case actorType == "player" && event.Action == "killed" && targetType == "player":
				kill := parseKillEvent(event, wrapper.OccurredAt, isFirstKill)
				kill.GameTime = gameTime
				data.Kills = append(data.Kills, kill)
				if isFirstKill {
					data.FirstBloodTime = wrapper.OccurredAt
					isFirstKill = false
				}

			// Dragon kills (A-tier NPC with dragon in ID)
			case event.Action == "killed" && targetType == "ATierNPC":
				targetID := ""
				if event.Target != nil {
					targetID = strings.ToLower(event.Target.ID)
				}
				
				if strings.Contains(targetID, "drake") || strings.Contains(targetID, "dragon") {
					dragon := parseDragonKillEvent(event, wrapper.OccurredAt)
					dragon.GameTime = gameTime
					data.DragonKills = append(data.DragonKills, dragon)
				} else if strings.Contains(targetID, "baron") || strings.Contains(targetID, "nashor") {
					baron := parseObjectiveKillEvent(event, wrapper.OccurredAt, "baron")
					baron.GameTime = gameTime
					data.BaronKills = append(data.BaronKills, baron)
				} else if strings.Contains(targetID, "herald") || strings.Contains(targetID, "rift") {
					herald := parseObjectiveKillEvent(event, wrapper.OccurredAt, "herald")
					herald.GameTime = gameTime
					data.HeraldKills = append(data.HeraldKills, herald)
				} else if strings.Contains(targetID, "voidgrub") || strings.Contains(targetID, "grub") {
					grub := parseObjectiveKillEvent(event, wrapper.OccurredAt, "voidGrub")
					grub.GameTime = gameTime
					data.VoidGrubKills = append(data.VoidGrubKills, grub)
				}

			// Tower destruction
			case event.Action == "destroyed" && (targetType == "tower" || targetType == "fortifier"):
				tower := parseTowerDestroyEvent(event, wrapper.OccurredAt, towerRegex)
				tower.GameTime = gameTime
				data.TowerDestroys = append(data.TowerDestroys, tower)

			// Draft picks and bans
			case event.Action == "picked" && targetType == "character":
				draft := parseDraftAction(event, wrapper.OccurredAt, "picked")
				data.DraftActions = append(data.DraftActions, draft)

			case event.Action == "banned" && targetType == "character":
				draft := parseDraftAction(event, wrapper.OccurredAt, "banned")
				data.DraftActions = append(data.DraftActions, draft)

			// Game start
			case event.Action == "started" && actorType == "series" && targetType == "game":
				data.GameStartTime = wrapper.OccurredAt
			}
		}
	}

	return data, nil
}

// ParseVALEvents extracts VALORANT-specific events from EventWrappers
func ParseVALEvents(wrappers []EventWrapper) (*VALEventData, error) {
	data := &VALEventData{
		RoundEnds: make([]RoundEndEvent, 0),
		Plants:    make([]PlantEvent, 0),
		Defuses:   make([]DefuseEvent, 0),
		Kills:     make([]VALKillEvent, 0),
	}

	// Track current round, map, and round start times
	currentRound := 0
	currentGameID := ""
	roundStartTimes := make(map[int]time.Time)

	// Helper to calculate game time from round start
	calcRoundGameTime := func(eventTime time.Time, roundNum int) int {
		if startTime, ok := roundStartTimes[roundNum]; ok {
			return int(eventTime.Sub(startTime).Milliseconds())
		}
		return 0
	}

	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			actorType := ""
			if event.Actor != nil {
				actorType = event.Actor.Type
			}
			targetType := ""
			if event.Target != nil {
				targetType = event.Target.Type
			}

			switch {
			// Round start - track round number and start time
			case event.Action == "started" && targetType == "round":
				if event.Target != nil {
					if seq, ok := event.Target.State["sequenceNumber"].(float64); ok {
						currentRound = int(seq)
						roundStartTimes[currentRound] = wrapper.OccurredAt
					}
				}
				// Extract map name from actor (game) state
				if event.Actor != nil && event.Actor.State != nil {
					if mapInfo, ok := event.Actor.State["map"].(map[string]interface{}); ok {
						if name, ok := mapInfo["name"].(string); ok {
							data.MapName = name
						}
					}
				}
				// Track game ID
				if event.Actor != nil {
					currentGameID = event.Actor.ID
				}

			// Round end with winner
			case event.Action == "won" && targetType == "round":
				roundEnd := parseRoundEndEvent(event, wrapper.OccurredAt, currentGameID)
				data.RoundEnds = append(data.RoundEnds, roundEnd)

			// Spike plant
			case event.Action == "completed" && targetType == "plantBomb":
				plant := parsePlantEvent(event, wrapper.OccurredAt, currentRound, data.MapName)
				plant.GameTime = calcRoundGameTime(wrapper.OccurredAt, currentRound)
				data.Plants = append(data.Plants, plant)

			// Spike defuse
			case event.Action == "completed" && targetType == "defuseBomb":
				defuse := parseDefuseEvent(event, wrapper.OccurredAt, currentRound)
				defuse.GameTime = calcRoundGameTime(wrapper.OccurredAt, currentRound)
				data.Defuses = append(data.Defuses, defuse)

			// Player kills
			case actorType == "player" && event.Action == "killed" && targetType == "player":
				kill := parseVALKillEvent(event, wrapper.OccurredAt, currentRound, data.MapName)
				kill.GameTime = calcRoundGameTime(wrapper.OccurredAt, currentRound)
				data.Kills = append(data.Kills, kill)
			}
		}
	}

	return data, nil
}

// Helper functions for parsing specific event types

func parseKillEvent(event GridEvent, occurredAt time.Time, isFirstBlood bool) KillEvent {
	kill := KillEvent{
		OccurredAt: occurredAt,
		FirstBlood: isFirstBlood,
	}

	// Extract killer info
	if event.Actor != nil {
		kill.KillerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				kill.KillerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				kill.KillerTeamID = teamID
			}
		}
		kill.KillerPosition = event.GetActorPosition()
	}

	// Extract victim info
	if event.Target != nil {
		kill.VictimID = event.Target.ID
		if event.Target.State != nil {
			if name, ok := event.Target.State["name"].(string); ok {
				kill.VictimName = name
			}
			if teamID, ok := event.Target.State["teamId"].(string); ok {
				kill.VictimTeamID = teamID
			}
		}
		kill.VictimPosition = event.GetTargetPosition()
	}

	// Extract assists from actor's game state
	if gameState := event.GetActorGameState(); gameState != nil {
		if assists, ok := gameState["killAssistsReceivedFromPlayer"].([]interface{}); ok {
			for _, a := range assists {
				if assist, ok := a.(map[string]interface{}); ok {
					if playerID, ok := assist["playerId"].(string); ok {
						kill.AssistIDs = append(kill.AssistIDs, playerID)
					}
				}
			}
		}
	}

	return kill
}

func parseDragonKillEvent(event GridEvent, occurredAt time.Time) DragonKillEvent {
	dragon := DragonKillEvent{
		OccurredAt: occurredAt,
	}

	// Extract dragon type from target ID
	if event.Target != nil {
		targetID := strings.ToLower(event.Target.ID)
		switch {
		case strings.Contains(targetID, "infernal"):
			dragon.DragonType = "infernal"
		case strings.Contains(targetID, "mountain"):
			dragon.DragonType = "mountain"
		case strings.Contains(targetID, "ocean"):
			dragon.DragonType = "ocean"
		case strings.Contains(targetID, "cloud"):
			dragon.DragonType = "cloud"
		case strings.Contains(targetID, "chemtech"):
			dragon.DragonType = "chemtech"
		case strings.Contains(targetID, "hextech"):
			dragon.DragonType = "hextech"
		case strings.Contains(targetID, "elder"):
			dragon.DragonType = "elder"
		default:
			dragon.DragonType = "unknown"
		}
	}

	// Extract player/team info
	if event.Actor != nil {
		dragon.PlayerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				dragon.PlayerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				dragon.TeamID = teamID
			}
		}
		dragon.Position = event.GetActorPosition()
	}

	return dragon
}

func parseObjectiveKillEvent(event GridEvent, occurredAt time.Time, objType string) ObjectiveKillEvent {
	obj := ObjectiveKillEvent{
		OccurredAt:    occurredAt,
		ObjectiveType: objType,
	}

	if event.Target != nil {
		obj.ObjectiveID = event.Target.ID
	}

	if event.Actor != nil {
		obj.PlayerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				obj.PlayerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				obj.TeamID = teamID
			}
		}
		obj.Position = event.GetActorPosition()
	}

	return obj
}

func parseTowerDestroyEvent(event GridEvent, occurredAt time.Time, towerRegex *regexp.Regexp) TowerDestroyEvent {
	tower := TowerDestroyEvent{
		OccurredAt: occurredAt,
	}

	// Extract tower info from target
	if event.Target != nil {
		tower.TowerID = event.Target.ID
		
		// Parse tower ID for lane and number
		matches := towerRegex.FindStringSubmatch(event.Target.ID)
		if len(matches) >= 3 {
			tower.Lane = matches[2] // top, mid, bot
			if len(matches) >= 4 && matches[3] != "" {
				fmt.Sscanf(matches[3], "%d", &tower.TowerNum)
			}
		}
	}

	// Extract team/player info from actor
	if event.Actor != nil {
		if event.Actor.Type == "team" {
			tower.TeamID = event.Actor.ID
			if event.Actor.State != nil {
				if name, ok := event.Actor.State["name"].(string); ok {
					tower.TeamName = name
				}
			}
		} else if event.Actor.Type == "player" {
			tower.PlayerID = event.Actor.ID
			if event.Actor.State != nil {
				if name, ok := event.Actor.State["name"].(string); ok {
					tower.PlayerName = name
				}
				if teamID, ok := event.Actor.State["teamId"].(string); ok {
					tower.TeamID = teamID
				}
			}
		}
	}

	return tower
}

func parseDraftAction(event GridEvent, occurredAt time.Time, action string) DraftAction {
	draft := DraftAction{
		Action:     action,
		OccurredAt: occurredAt,
	}

	// Extract team info from actor
	if event.Actor != nil {
		draft.TeamID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				draft.TeamName = name
			}
		}
	}

	// Extract character info from target
	if event.Target != nil {
		draft.CharacterID = event.Target.ID
		if event.Target.State != nil {
			if name, ok := event.Target.State["name"].(string); ok {
				draft.CharacterName = name
			}
		}
	}

	return draft
}

func parseRoundEndEvent(event GridEvent, occurredAt time.Time, gameID string) RoundEndEvent {
	roundEnd := RoundEndEvent{
		GameID:     gameID,
		OccurredAt: occurredAt,
	}

	// Extract winner info from actor (team)
	if event.Actor != nil {
		roundEnd.WinnerTeam = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				roundEnd.WinnerName = name
			}
		}
	}

	// Extract round number and win type from target
	if event.Target != nil {
		if event.Target.State != nil {
			if seq, ok := event.Target.State["sequenceNumber"].(float64); ok {
				roundEnd.RoundNum = int(seq)
			}
			// Win type is in teams array
			if teams, ok := event.Target.State["teams"].([]interface{}); ok {
				for _, t := range teams {
					if team, ok := t.(map[string]interface{}); ok {
						if won, ok := team["won"].(bool); ok && won {
							if winType, ok := team["winType"].(string); ok {
								roundEnd.WinType = winType
							}
						}
					}
				}
			}
		}
	}

	return roundEnd
}

func parsePlantEvent(event GridEvent, occurredAt time.Time, roundNum int, mapName string) PlantEvent {
	plant := PlantEvent{
		RoundNum:   roundNum,
		MapName:    mapName,
		OccurredAt: occurredAt,
	}

	// Extract player info from actor
	if event.Actor != nil {
		plant.PlayerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				plant.PlayerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				plant.TeamID = teamID
			}
			// Get agent from character
			if game, ok := event.Actor.State["game"].(map[string]interface{}); ok {
				if char, ok := game["character"].(map[string]interface{}); ok {
					if name, ok := char["name"].(string); ok {
						plant.Agent = name
					}
				}
			}
		}
		plant.Position = event.GetActorPosition()
		
		// Infer site from position if available
		if plant.Position != nil && mapName != "" {
			plant.Site = inferSiteFromPosition(plant.Position, mapName)
		}
	}

	return plant
}

func parseDefuseEvent(event GridEvent, occurredAt time.Time, roundNum int) DefuseEvent {
	defuse := DefuseEvent{
		RoundNum:   roundNum,
		OccurredAt: occurredAt,
	}

	if event.Actor != nil {
		defuse.PlayerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				defuse.PlayerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				defuse.TeamID = teamID
			}
			if game, ok := event.Actor.State["game"].(map[string]interface{}); ok {
				if char, ok := game["character"].(map[string]interface{}); ok {
					if name, ok := char["name"].(string); ok {
						defuse.Agent = name
					}
				}
			}
		}
		defuse.Position = event.GetActorPosition()
	}

	return defuse
}

func parseVALKillEvent(event GridEvent, occurredAt time.Time, roundNum int, mapName string) VALKillEvent {
	kill := VALKillEvent{
		RoundNum:   roundNum,
		MapName:    mapName,
		OccurredAt: occurredAt,
	}

	// Extract killer info
	if event.Actor != nil {
		kill.KillerID = event.Actor.ID
		if event.Actor.State != nil {
			if name, ok := event.Actor.State["name"].(string); ok {
				kill.KillerName = name
			}
			if teamID, ok := event.Actor.State["teamId"].(string); ok {
				kill.KillerTeamID = teamID
			}
			if game, ok := event.Actor.State["game"].(map[string]interface{}); ok {
				if char, ok := game["character"].(map[string]interface{}); ok {
					if name, ok := char["name"].(string); ok {
						kill.KillerAgent = name
					}
				}
			}
		}
		kill.KillerPosition = event.GetActorPosition()
	}

	// Extract victim info
	if event.Target != nil {
		kill.VictimID = event.Target.ID
		if event.Target.State != nil {
			if name, ok := event.Target.State["name"].(string); ok {
				kill.VictimName = name
			}
			if teamID, ok := event.Target.State["teamId"].(string); ok {
				kill.VictimTeamID = teamID
			}
			if game, ok := event.Target.State["game"].(map[string]interface{}); ok {
				if char, ok := game["character"].(map[string]interface{}); ok {
					if name, ok := char["name"].(string); ok {
						kill.VictimAgent = name
					}
				}
			}
		}
		kill.VictimPosition = event.GetTargetPosition()
	}

	return kill
}

// inferSiteFromPosition determines the spike site based on plant position
// This uses approximate site boundaries for each VALORANT map
func inferSiteFromPosition(pos *Position, mapName string) string {
	if pos == nil {
		return ""
	}

	// Site boundaries are approximate and based on map analysis
	// These would need to be refined with more data
	switch strings.ToLower(mapName) {
	case "ascent":
		// A site is roughly in the upper area, B site in the lower
		if pos.Y > 0 {
			return "A"
		}
		return "B"
	case "bind":
		// A site is on the left, B site on the right
		if pos.X < 0 {
			return "A"
		}
		return "B"
	case "haven":
		// Three sites: A (left), B (middle), C (right)
		if pos.X < -2000 {
			return "A"
		} else if pos.X > 2000 {
			return "C"
		}
		return "B"
	case "split":
		// A site is upper, B site is lower
		if pos.Y > 0 {
			return "A"
		}
		return "B"
	case "breeze":
		// A site is on one side, B on the other
		if pos.X > 5000 {
			return "A"
		}
		return "B"
	case "icebox":
		// A site is upper, B site is lower
		if pos.Y > 0 {
			return "A"
		}
		return "B"
	case "fracture":
		// A site is on one side, B on the other
		if pos.X > 0 {
			return "A"
		}
		return "B"
	case "pearl":
		// A site is upper, B site is lower
		if pos.Y > 0 {
			return "A"
		}
		return "B"
	case "lotus":
		// Three sites: A, B, C
		if pos.X < -2000 {
			return "A"
		} else if pos.X > 2000 {
			return "C"
		}
		return "B"
	case "sunset":
		// A site is on one side, B on the other
		if pos.X > 0 {
			return "A"
		}
		return "B"
	case "abyss":
		// A site is on one side, B on the other
		if pos.X > 0 {
			return "A"
		}
		return "B"
	default:
		return ""
	}
}

// GetEventsForSeries is a convenience method that downloads and parses events
func (c *Client) GetEventsForSeries(ctx context.Context, seriesID string, titleID string) (interface{}, error) {
	wrappers, err := c.DownloadEvents(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	// Parse based on game type
	switch titleID {
	case "3", "lol", "league-of-legends":
		return ParseLoLEvents(wrappers)
	case "25", "val", "valorant":
		return ParseVALEvents(wrappers)
	default:
		return wrappers, nil
	}
}

// ParseLoLEvents is a method wrapper for the package function (for backward compatibility)
func (c *Client) ParseLoLEvents(wrappers []EventWrapper) (*LoLEventData, error) {
	return ParseLoLEvents(wrappers)
}

// ParseVALEvents is a method wrapper for the package function (for backward compatibility)
func (c *Client) ParseVALEvents(wrappers []EventWrapper) (*VALEventData, error) {
	return ParseVALEvents(wrappers)
}


// DownloadAndParseVALEvents downloads and parses VALORANT events for a series
// This is a convenience method that combines DownloadEvents and ParseVALEvents
func (c *Client) DownloadAndParseVALEvents(ctx context.Context, seriesID string) (*VALEventData, error) {
	wrappers, err := c.DownloadEvents(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	return ParseVALEvents(wrappers)
}

// DownloadAndParseLoLEvents downloads and parses LoL events for a series
// This is a convenience method that combines DownloadEvents and ParseLoLEvents
func (c *Client) DownloadAndParseLoLEvents(ctx context.Context, seriesID string) (*LoLEventData, error) {
	wrappers, err := c.DownloadEvents(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	return ParseLoLEvents(wrappers)
}
