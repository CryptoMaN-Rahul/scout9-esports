package intelligence

import (
	"sort"

	"scout9/pkg/grid"
)

// Game phase time boundaries in seconds
const (
	EarlyGameEndSeconds = 900  // 15 minutes
	MidGameEndSeconds   = 1500 // 25 minutes
)

// EventAnalyzer analyzes JSONL events to extract first blood, dragon, tower,
// game phase, ability usage, and objective timing data.
type EventAnalyzer struct{}

// NewEventAnalyzer creates a new EventAnalyzer
func NewEventAnalyzer() *EventAnalyzer {
	return &EventAnalyzer{}
}

// FirstsAnalysis contains first blood, dragon, tower data for a team
type FirstsAnalysis struct {
	// First Blood
	FirstBloodCount   int
	FirstBloodGames   int
	FirstBloodRate    float64
	AvgFirstBloodTime float64        // seconds
	FirstBloodPlayers map[string]int // playerID -> count

	// First Dragon
	FirstDragonCount   int
	FirstDragonGames   int
	FirstDragonRate    float64
	AvgFirstDragonTime float64        // seconds
	FirstDragonTypes   map[string]int // dragonType -> count

	// First Tower
	FirstTowerCount   int
	FirstTowerGames   int
	FirstTowerRate    float64
	AvgFirstTowerTime float64        // seconds
	FirstTowerLanes   map[string]int // lane -> count
}

// PhaseAnalysis contains game phase performance data
type PhaseAnalysis struct {
	// Early game (0-15 min)
	EarlyKills      int
	EarlyDeaths     int
	EarlyObjectives int
	EarlyKillDiff   int

	// Mid game (15-25 min)
	MidKills      int
	MidDeaths     int
	MidObjectives int
	MidKillDiff   int

	// Late game (25+ min)
	LateKills      int
	LateDeaths     int
	LateObjectives int
	LateKillDiff   int

	// Derived
	StrongestPhase  string  // "early", "mid", "late"
	EarlyGameRating float64 // 0-100
	MidGameRating   float64 // 0-100
	LateGameRating  float64 // 0-100
}

// AbilityAnalysis contains ability usage data from JSONL events
type AbilityAnalysis struct {
	PlayerAbilities map[string]map[string]int // playerID -> abilityID -> count
	TotalGames      int
}

// EventObjectiveTimings contains objective timing data from JSONL events
type EventObjectiveTimings struct {
	DragonTimings []int // seconds
	BaronTimings  []int // seconds
	HeraldTimings []int // seconds
	TowerTimings  []int // seconds

	AvgFirstDragonTime float64
	AvgBaronTime       float64
	AvgHeraldTime      float64
	AvgFirstTowerTime  float64

	EarlyObjectivePriority bool // first dragon < 360 seconds (6 min)
}

// FirstObjectiveData tracks first objective data for a single game
type FirstObjectiveData struct {
	GameID          string
	TeamID          string
	ObjectiveType   string // "blood", "dragon", "tower"
	GameTimeSeconds int
	PlayerID        string
	PlayerName      string
	DragonType      string // for dragons only
	Lane            string // for towers only
}

// ClassifyGamePhase returns the game phase for a given time in seconds
func ClassifyGamePhase(gameTimeSeconds int) string {
	if gameTimeSeconds < EarlyGameEndSeconds {
		return "early"
	}
	if gameTimeSeconds < MidGameEndSeconds {
		return "mid"
	}
	return "late"
}

// AnalyzeFirstBlood finds the first blood in a game and returns data about it
func (ea *EventAnalyzer) AnalyzeFirstBlood(events *grid.LoLEventData, teamID string) *FirstObjectiveData {
	if events == nil || len(events.Kills) == 0 {
		return nil
	}

	// Sort kills by game time to find first
	kills := make([]grid.KillEvent, len(events.Kills))
	copy(kills, events.Kills)
	sort.Slice(kills, func(i, j int) bool {
		return kills[i].GameTime < kills[j].GameTime
	})

	// First kill is first blood
	firstKill := kills[0]

	// Check if this team got first blood
	if firstKill.KillerTeamID != teamID {
		return nil
	}

	return &FirstObjectiveData{
		ObjectiveType:   "blood",
		TeamID:          firstKill.KillerTeamID,
		GameTimeSeconds: firstKill.GameTime / 1000, // Convert ms to seconds
		PlayerID:        firstKill.KillerID,
		PlayerName:      firstKill.KillerName,
	}
}

// AnalyzeFirstDragon finds the first dragon in a game and returns data about it
func (ea *EventAnalyzer) AnalyzeFirstDragon(events *grid.LoLEventData, teamID string) *FirstObjectiveData {
	if events == nil || len(events.DragonKills) == 0 {
		return nil
	}

	// Sort dragons by game time to find first
	dragons := make([]grid.DragonKillEvent, len(events.DragonKills))
	copy(dragons, events.DragonKills)
	sort.Slice(dragons, func(i, j int) bool {
		return dragons[i].GameTime < dragons[j].GameTime
	})

	// First dragon
	firstDragon := dragons[0]

	// Check if this team got first dragon
	if firstDragon.TeamID != teamID {
		return nil
	}

	return &FirstObjectiveData{
		ObjectiveType:   "dragon",
		TeamID:          firstDragon.TeamID,
		GameTimeSeconds: firstDragon.GameTime / 1000, // Convert ms to seconds
		PlayerID:        firstDragon.PlayerID,
		PlayerName:      firstDragon.PlayerName,
		DragonType:      firstDragon.DragonType,
	}
}

// AnalyzeFirstTower finds the first tower in a game and returns data about it
func (ea *EventAnalyzer) AnalyzeFirstTower(events *grid.LoLEventData, teamID string) *FirstObjectiveData {
	if events == nil || len(events.TowerDestroys) == 0 {
		return nil
	}

	// Sort towers by game time to find first
	towers := make([]grid.TowerDestroyEvent, len(events.TowerDestroys))
	copy(towers, events.TowerDestroys)
	sort.Slice(towers, func(i, j int) bool {
		return towers[i].GameTime < towers[j].GameTime
	})

	// First tower
	firstTower := towers[0]

	// Check if this team got first tower
	if firstTower.TeamID != teamID {
		return nil
	}

	return &FirstObjectiveData{
		ObjectiveType:   "tower",
		TeamID:          firstTower.TeamID,
		GameTimeSeconds: firstTower.GameTime / 1000, // Convert ms to seconds
		PlayerID:        firstTower.PlayerID,
		PlayerName:      firstTower.PlayerName,
		Lane:            firstTower.Lane,
	}
}

// AnalyzeFirsts analyzes JSONL events to extract first blood/dragon/tower data
// for a specific team across multiple games
func (ea *EventAnalyzer) AnalyzeFirsts(eventsPerGame []*grid.LoLEventData, teamID string) *FirstsAnalysis {
	analysis := &FirstsAnalysis{
		FirstBloodPlayers: make(map[string]int),
		FirstDragonTypes:  make(map[string]int),
		FirstTowerLanes:   make(map[string]int),
	}

	var totalFirstBloodTime, totalFirstDragonTime, totalFirstTowerTime int
	gamesWithKills, gamesWithDragons, gamesWithTowers := 0, 0, 0

	for _, events := range eventsPerGame {
		if events == nil {
			continue
		}

		// First Blood
		if len(events.Kills) > 0 {
			gamesWithKills++
			fb := ea.AnalyzeFirstBlood(events, teamID)
			if fb != nil {
				analysis.FirstBloodCount++
				totalFirstBloodTime += fb.GameTimeSeconds
				if fb.PlayerID != "" {
					analysis.FirstBloodPlayers[fb.PlayerID]++
				}
			}
		}

		// First Dragon
		if len(events.DragonKills) > 0 {
			gamesWithDragons++
			fd := ea.AnalyzeFirstDragon(events, teamID)
			if fd != nil {
				analysis.FirstDragonCount++
				totalFirstDragonTime += fd.GameTimeSeconds
				if fd.DragonType != "" {
					analysis.FirstDragonTypes[fd.DragonType]++
				}
			}
		}

		// First Tower
		if len(events.TowerDestroys) > 0 {
			gamesWithTowers++
			ft := ea.AnalyzeFirstTower(events, teamID)
			if ft != nil {
				analysis.FirstTowerCount++
				totalFirstTowerTime += ft.GameTimeSeconds
				if ft.Lane != "" {
					analysis.FirstTowerLanes[ft.Lane]++
				}
			}
		}
	}

	// Calculate rates
	analysis.FirstBloodGames = gamesWithKills
	if gamesWithKills > 0 {
		analysis.FirstBloodRate = float64(analysis.FirstBloodCount) / float64(gamesWithKills)
	}
	if analysis.FirstBloodCount > 0 {
		analysis.AvgFirstBloodTime = float64(totalFirstBloodTime) / float64(analysis.FirstBloodCount)
	}

	analysis.FirstDragonGames = gamesWithDragons
	if gamesWithDragons > 0 {
		analysis.FirstDragonRate = float64(analysis.FirstDragonCount) / float64(gamesWithDragons)
	}
	if analysis.FirstDragonCount > 0 {
		analysis.AvgFirstDragonTime = float64(totalFirstDragonTime) / float64(analysis.FirstDragonCount)
	}

	analysis.FirstTowerGames = gamesWithTowers
	if gamesWithTowers > 0 {
		analysis.FirstTowerRate = float64(analysis.FirstTowerCount) / float64(gamesWithTowers)
	}
	if analysis.FirstTowerCount > 0 {
		analysis.AvgFirstTowerTime = float64(totalFirstTowerTime) / float64(analysis.FirstTowerCount)
	}

	return analysis
}


// AnalyzePhases analyzes JSONL events to extract game phase performance
func (ea *EventAnalyzer) AnalyzePhases(events *grid.LoLEventData, teamID string) *PhaseAnalysis {
	analysis := &PhaseAnalysis{}

	if events == nil {
		return analysis
	}

	// Count kills per phase
	for _, kill := range events.Kills {
		gameTimeSeconds := kill.GameTime / 1000 // Convert ms to seconds
		phase := ClassifyGamePhase(gameTimeSeconds)

		isTeamKill := kill.KillerTeamID == teamID
		isTeamDeath := kill.VictimTeamID == teamID

		switch phase {
		case "early":
			if isTeamKill {
				analysis.EarlyKills++
			}
			if isTeamDeath {
				analysis.EarlyDeaths++
			}
		case "mid":
			if isTeamKill {
				analysis.MidKills++
			}
			if isTeamDeath {
				analysis.MidDeaths++
			}
		case "late":
			if isTeamKill {
				analysis.LateKills++
			}
			if isTeamDeath {
				analysis.LateDeaths++
			}
		}
	}

	// Count objectives per phase (dragons, barons, heralds)
	for _, dragon := range events.DragonKills {
		if dragon.TeamID != teamID {
			continue
		}
		gameTimeSeconds := dragon.GameTime / 1000
		phase := ClassifyGamePhase(gameTimeSeconds)
		switch phase {
		case "early":
			analysis.EarlyObjectives++
		case "mid":
			analysis.MidObjectives++
		case "late":
			analysis.LateObjectives++
		}
	}

	for _, baron := range events.BaronKills {
		if baron.TeamID != teamID {
			continue
		}
		gameTimeSeconds := baron.GameTime / 1000
		phase := ClassifyGamePhase(gameTimeSeconds)
		switch phase {
		case "early":
			analysis.EarlyObjectives++
		case "mid":
			analysis.MidObjectives++
		case "late":
			analysis.LateObjectives++
		}
	}

	for _, herald := range events.HeraldKills {
		if herald.TeamID != teamID {
			continue
		}
		gameTimeSeconds := herald.GameTime / 1000
		phase := ClassifyGamePhase(gameTimeSeconds)
		switch phase {
		case "early":
			analysis.EarlyObjectives++
		case "mid":
			analysis.MidObjectives++
		case "late":
			analysis.LateObjectives++
		}
	}

	// Calculate kill differentials
	analysis.EarlyKillDiff = analysis.EarlyKills - analysis.EarlyDeaths
	analysis.MidKillDiff = analysis.MidKills - analysis.MidDeaths
	analysis.LateKillDiff = analysis.LateKills - analysis.LateDeaths

	// Identify strongest phase
	analysis.StrongestPhase = ea.identifyStrongestPhase(analysis)

	// Calculate phase ratings (0-100 scale based on kill diff)
	analysis.EarlyGameRating = ea.calculatePhaseRating(analysis.EarlyKillDiff)
	analysis.MidGameRating = ea.calculatePhaseRating(analysis.MidKillDiff)
	analysis.LateGameRating = ea.calculatePhaseRating(analysis.LateKillDiff)

	return analysis
}

// identifyStrongestPhase returns the phase with the highest kill differential
func (ea *EventAnalyzer) identifyStrongestPhase(analysis *PhaseAnalysis) string {
	maxDiff := analysis.EarlyKillDiff
	strongest := "early"

	if analysis.MidKillDiff > maxDiff {
		maxDiff = analysis.MidKillDiff
		strongest = "mid"
	}

	if analysis.LateKillDiff > maxDiff {
		strongest = "late"
	}

	return strongest
}

// calculatePhaseRating converts kill differential to 0-100 rating
// +10 diff = 100, 0 diff = 50, -10 diff = 0
func (ea *EventAnalyzer) calculatePhaseRating(killDiff int) float64 {
	// Clamp to -10 to +10 range
	if killDiff > 10 {
		killDiff = 10
	}
	if killDiff < -10 {
		killDiff = -10
	}

	// Convert to 0-100 scale: -10 -> 0, 0 -> 50, +10 -> 100
	return float64(killDiff+10) * 5.0
}

// AggregatePhases combines phase analysis from multiple games
func (ea *EventAnalyzer) AggregatePhases(analyses []*PhaseAnalysis) *PhaseAnalysis {
	if len(analyses) == 0 {
		return &PhaseAnalysis{}
	}

	total := &PhaseAnalysis{}

	for _, a := range analyses {
		if a == nil {
			continue
		}
		total.EarlyKills += a.EarlyKills
		total.EarlyDeaths += a.EarlyDeaths
		total.EarlyObjectives += a.EarlyObjectives
		total.MidKills += a.MidKills
		total.MidDeaths += a.MidDeaths
		total.MidObjectives += a.MidObjectives
		total.LateKills += a.LateKills
		total.LateDeaths += a.LateDeaths
		total.LateObjectives += a.LateObjectives
	}

	// Recalculate differentials and ratings
	total.EarlyKillDiff = total.EarlyKills - total.EarlyDeaths
	total.MidKillDiff = total.MidKills - total.MidDeaths
	total.LateKillDiff = total.LateKills - total.LateDeaths

	total.StrongestPhase = ea.identifyStrongestPhase(total)
	total.EarlyGameRating = ea.calculatePhaseRating(total.EarlyKillDiff)
	total.MidGameRating = ea.calculatePhaseRating(total.MidKillDiff)
	total.LateGameRating = ea.calculatePhaseRating(total.LateKillDiff)

	return total
}

// AnalyzeAbilities analyzes JSONL events to extract ability usage counts
func (ea *EventAnalyzer) AnalyzeAbilities(wrappers []grid.EventWrapper) *AbilityAnalysis {
	analysis := &AbilityAnalysis{
		PlayerAbilities: make(map[string]map[string]int),
		TotalGames:      1, // Assume single game for now
	}

	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			// Look for ability usage events
			if event.Action != "used" {
				continue
			}

			// Check if target is an ability
			if event.Target == nil || event.Target.Type != "ability" {
				continue
			}

			// Get player ID from actor
			if event.Actor == nil || event.Actor.Type != "player" {
				continue
			}

			playerID := event.Actor.ID
			abilityID := event.Target.ID

			// Initialize player map if needed
			if analysis.PlayerAbilities[playerID] == nil {
				analysis.PlayerAbilities[playerID] = make(map[string]int)
			}

			analysis.PlayerAbilities[playerID][abilityID]++
		}
	}

	return analysis
}

// GetMostUsedAbility returns the most-used ability for a player
func (ea *EventAnalyzer) GetMostUsedAbility(analysis *AbilityAnalysis, playerID string) (string, int) {
	if analysis == nil || analysis.PlayerAbilities == nil {
		return "", 0
	}

	abilities := analysis.PlayerAbilities[playerID]
	if abilities == nil {
		return "", 0
	}

	var maxAbility string
	var maxCount int

	for abilityID, count := range abilities {
		if count > maxCount {
			maxCount = count
			maxAbility = abilityID
		}
	}

	return maxAbility, maxCount
}

// AnalyzeObjectiveTimings analyzes JSONL events to extract objective timing patterns
func (ea *EventAnalyzer) AnalyzeObjectiveTimings(events *grid.LoLEventData, teamID string) *EventObjectiveTimings {
	analysis := &EventObjectiveTimings{
		DragonTimings: make([]int, 0),
		BaronTimings:  make([]int, 0),
		HeraldTimings: make([]int, 0),
		TowerTimings:  make([]int, 0),
	}

	if events == nil {
		return analysis
	}

	// Collect dragon timings
	for _, dragon := range events.DragonKills {
		if dragon.TeamID == teamID {
			analysis.DragonTimings = append(analysis.DragonTimings, dragon.GameTime/1000)
		}
	}

	// Collect baron timings
	for _, baron := range events.BaronKills {
		if baron.TeamID == teamID {
			analysis.BaronTimings = append(analysis.BaronTimings, baron.GameTime/1000)
		}
	}

	// Collect herald timings
	for _, herald := range events.HeraldKills {
		if herald.TeamID == teamID {
			analysis.HeraldTimings = append(analysis.HeraldTimings, herald.GameTime/1000)
		}
	}

	// Collect tower timings
	for _, tower := range events.TowerDestroys {
		if tower.TeamID == teamID {
			analysis.TowerTimings = append(analysis.TowerTimings, tower.GameTime/1000)
		}
	}

	// Calculate averages
	if len(analysis.DragonTimings) > 0 {
		// First dragon time
		sort.Ints(analysis.DragonTimings)
		analysis.AvgFirstDragonTime = float64(analysis.DragonTimings[0])

		// Check early objective priority (first dragon < 6 min = 360 seconds)
		analysis.EarlyObjectivePriority = analysis.DragonTimings[0] < 360
	}

	if len(analysis.BaronTimings) > 0 {
		sum := 0
		for _, t := range analysis.BaronTimings {
			sum += t
		}
		analysis.AvgBaronTime = float64(sum) / float64(len(analysis.BaronTimings))
	}

	if len(analysis.HeraldTimings) > 0 {
		sum := 0
		for _, t := range analysis.HeraldTimings {
			sum += t
		}
		analysis.AvgHeraldTime = float64(sum) / float64(len(analysis.HeraldTimings))
	}

	if len(analysis.TowerTimings) > 0 {
		sort.Ints(analysis.TowerTimings)
		analysis.AvgFirstTowerTime = float64(analysis.TowerTimings[0])
	}

	return analysis
}

// AggregateObjectiveTimings combines timing analysis from multiple games
func (ea *EventAnalyzer) AggregateObjectiveTimings(analyses []*EventObjectiveTimings) *EventObjectiveTimings {
	if len(analyses) == 0 {
		return &EventObjectiveTimings{}
	}

	total := &EventObjectiveTimings{
		DragonTimings: make([]int, 0),
		BaronTimings:  make([]int, 0),
		HeraldTimings: make([]int, 0),
		TowerTimings:  make([]int, 0),
	}

	var firstDragonTimes, baronTimes, heraldTimes, firstTowerTimes []float64
	earlyPriorityCount := 0

	for _, a := range analyses {
		if a == nil {
			continue
		}

		total.DragonTimings = append(total.DragonTimings, a.DragonTimings...)
		total.BaronTimings = append(total.BaronTimings, a.BaronTimings...)
		total.HeraldTimings = append(total.HeraldTimings, a.HeraldTimings...)
		total.TowerTimings = append(total.TowerTimings, a.TowerTimings...)

		if a.AvgFirstDragonTime > 0 {
			firstDragonTimes = append(firstDragonTimes, a.AvgFirstDragonTime)
		}
		if a.AvgBaronTime > 0 {
			baronTimes = append(baronTimes, a.AvgBaronTime)
		}
		if a.AvgHeraldTime > 0 {
			heraldTimes = append(heraldTimes, a.AvgHeraldTime)
		}
		if a.AvgFirstTowerTime > 0 {
			firstTowerTimes = append(firstTowerTimes, a.AvgFirstTowerTime)
		}
		if a.EarlyObjectivePriority {
			earlyPriorityCount++
		}
	}

	// Calculate averages
	if len(firstDragonTimes) > 0 {
		sum := 0.0
		for _, t := range firstDragonTimes {
			sum += t
		}
		total.AvgFirstDragonTime = sum / float64(len(firstDragonTimes))
	}

	if len(baronTimes) > 0 {
		sum := 0.0
		for _, t := range baronTimes {
			sum += t
		}
		total.AvgBaronTime = sum / float64(len(baronTimes))
	}

	if len(heraldTimes) > 0 {
		sum := 0.0
		for _, t := range heraldTimes {
			sum += t
		}
		total.AvgHeraldTime = sum / float64(len(heraldTimes))
	}

	if len(firstTowerTimes) > 0 {
		sum := 0.0
		for _, t := range firstTowerTimes {
			sum += t
		}
		total.AvgFirstTowerTime = sum / float64(len(firstTowerTimes))
	}

	// Early priority if majority of games had early first dragon
	total.EarlyObjectivePriority = earlyPriorityCount > len(analyses)/2

	return total
}
