package grid

import "time"

// Title represents a game title (LoL, VALORANT)
type Title struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Tournament represents a tournament/league
type Tournament struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LogoURL   string    `json:"logoUrl,omitempty"`
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
	TitleID   string    `json:"titleId"`
}

// Team represents an esports team
type Team struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LogoURL  string `json:"logoUrl,omitempty"`
	ColorHex string `json:"colorHex,omitempty"`
}

// Player represents an individual player
type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role,omitempty"`
	TeamID   string `json:"teamId,omitempty"`
}

// Series represents a match series (Bo1, Bo3, Bo5)
type Series struct {
	ID           string    `json:"id"`
	TournamentID string    `json:"tournamentId"`
	StartTime    time.Time `json:"startTime"`
	Format       string    `json:"format,omitempty"`
	Teams        []Team    `json:"teams"`
	TitleID      string    `json:"titleId"`
}

// SeriesState represents the detailed state of a series
type SeriesState struct {
	ID       string      `json:"id"`
	Started  bool        `json:"started"`
	Finished bool        `json:"finished"`
	Teams    []TeamState `json:"teams"`
	Games    []Game      `json:"games"`
}

// TeamState represents team-level state in a series
type TeamState struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Won    bool   `json:"won"`
	Score  int    `json:"score"`
	Kills  int    `json:"kills"`
	Deaths int    `json:"deaths"`
}

// Game represents a single game within a series
type Game struct {
	ID           string        `json:"id"`
	Sequence     int           `json:"sequence"`
	Map          string        `json:"map,omitempty"` // VALORANT only
	Duration     int           `json:"duration"`      // seconds
	Finished     bool          `json:"finished"`
	Started      bool          `json:"started"`
	Paused       bool          `json:"paused,omitempty"`
	Teams        []GameTeam    `json:"teams"`
	StartTime    time.Time     `json:"startTime,omitempty"`
	DraftActions []DraftAction `json:"draftActions,omitempty"` // LoL draft picks/bans
	Segments     []Segment     `json:"segments,omitempty"`     // Rounds/segments (VALORANT)
}

// GameTeam represents team-level stats for a game
type GameTeam struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	Side                string          `json:"side,omitempty"` // "blue"/"red" for LoL, "attack"/"defense" for VAL
	Score               int             `json:"score"`
	Won                 bool            `json:"won"`
	Kills               int             `json:"kills"`
	Deaths              int             `json:"deaths"`
	NetWorth            int             `json:"netWorth,omitempty"`     // LoL only
	Money               int             `json:"money,omitempty"`        // Team money
	LoadoutValue        int             `json:"loadoutValue,omitempty"` // Team loadout value
	StructuresDestroyed int             `json:"structuresDestroyed"`
	Players             []GamePlayer    `json:"players"`
	Objectives          []TeamObjective `json:"objectives,omitempty"` // Team objectives (dragons, barons, etc.)
}

// TeamObjective represents an objective completed by a team
type TeamObjective struct {
	ID              string `json:"id"`
	Type            string `json:"type"` // e.g., "slayInfernalDrake", "destroyTower", "slayBaron"
	CompletionCount int    `json:"completionCount"`
}

// GamePlayer represents player-level stats for a game
type GamePlayer struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Character     string         `json:"character"` // champion or agent
	Kills         int            `json:"kills"`
	Deaths        int            `json:"deaths"`
	Assists       int            `json:"assists"`                 // killAssistsReceived from API
	AssistsGiven  int            `json:"assistsGiven"`            // killAssistsGiven from API
	AssistDetails []AssistDetail `json:"assistDetails,omitempty"` // Who assisted on kills (killAssistsReceivedFromPlayer)
	Teamkills     int            `json:"teamkills,omitempty"`     // Friendly fire kills
	Selfkills     int            `json:"selfkills,omitempty"`     // Self-kills
	// NOTE: Damage field is NOT available via GRID API - removed to avoid confusion
	NetWorth            int               `json:"netWorth,omitempty"`     // LoL only
	Money               int               `json:"money,omitempty"`        // Current money
	LoadoutValue        int               `json:"loadoutValue,omitempty"` // Loadout value
	Items               []Item            `json:"items,omitempty"`        // Player inventory items
	StructuresDestroyed int               `json:"structuresDestroyed"`    // Towers/structures destroyed
	Objectives          []PlayerObjective `json:"objectives,omitempty"`   // Objectives completed
	Multikills          []Multikill       `json:"multikills,omitempty"`   // Multi-kill stats
	WeaponKills         []WeaponKill      `json:"weaponKills,omitempty"`  // Kills by weapon (VALORANT)
	Abilities           []AbilityUsage    `json:"abilities,omitempty"`    // Ability usage (VALORANT only - IDs)
}

// AssistDetail represents who assisted on a kill
type AssistDetail struct {
	PlayerID        string `json:"playerId"`
	PlayerName      string `json:"playerName,omitempty"`
	AssistsReceived int    `json:"assistsReceived"`
}

// AbilityUsage represents ability usage (VALORANT only - just IDs)
type AbilityUsage struct {
	ID          string `json:"id"`
	AbilityName string `json:"abilityName"` // Same as ID for VALORANT
}

// Item represents an inventory item
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PlayerObjective represents an objective completed by a player
type PlayerObjective struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	CompletionCount int    `json:"completionCount"`
}

// Multikill represents multi-kill statistics
type Multikill struct {
	ID            string `json:"id"`
	NumberOfKills int    `json:"numberOfKills"` // 2=double, 3=triple, 4=quadra, 5=penta
	Count         int    `json:"count"`
}

// WeaponKill represents kills by weapon type
type WeaponKill struct {
	ID         string `json:"id"`
	WeaponName string `json:"weaponName"`
	Count      int    `json:"count"`
}

// Segment represents a round/segment within a game (VALORANT rounds, LoL phases)
type Segment struct {
	ID             string        `json:"id"`
	SequenceNumber int           `json:"sequenceNumber"`
	Type           string        `json:"type"` // "round", "half", etc.
	Finished       bool          `json:"finished"`
	Teams          []SegmentTeam `json:"teams"`
}

// SegmentTeam represents team state for a specific segment/round
type SegmentTeam struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Side       string          `json:"side"` // "attack"/"defense"
	Won        bool            `json:"won"`
	Kills      int             `json:"kills"`
	Deaths     int             `json:"deaths"`
	Objectives []TeamObjective `json:"objectives,omitempty"`
	Players    []SegmentPlayer `json:"players,omitempty"`
}

// SegmentPlayer represents player state for a specific segment/round
type SegmentPlayer struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Kills      int               `json:"kills"`
	Deaths     int               `json:"deaths"`
	Objectives []PlayerObjective `json:"objectives,omitempty"`
}

// FileInfo represents a downloadable file from File Download API
type FileInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	FileName    string `json:"fileName"`
	FullURL     string `json:"fullURL"`
}

// EventWrapper represents a single line in the GRID JSONL file
// Each line contains one or more events in the Events array
type EventWrapper struct {
	ID             string      `json:"id"`
	CorrelationID  string      `json:"correlationId"`
	OccurredAt     time.Time   `json:"occurredAt"`
	SeriesID       string      `json:"seriesId"`
	SequenceNumber int         `json:"sequenceNumber"`
	Events         []GridEvent `json:"events"`
}

// GridEvent represents a single event within an EventWrapper
// Event type is constructed as: {actor.type}-{action}-{target.type}
type GridEvent struct {
	ID                string                 `json:"id"`
	Action            string                 `json:"action"`
	Actor             *EventEntity           `json:"actor,omitempty"`
	Target            *EventEntity           `json:"target,omitempty"`
	SeriesState       map[string]interface{} `json:"seriesState,omitempty"`
	SeriesStateDelta  map[string]interface{} `json:"seriesStateDelta,omitempty"`
	IncludesFullState bool                   `json:"includesFullState"`
}

// EventEntity represents an actor or target in a GridEvent
type EventEntity struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	State      map[string]interface{} `json:"state,omitempty"`
	StateDelta map[string]interface{} `json:"stateDelta,omitempty"`
}

// GetEventType returns the constructed event type string
// Format: {actor.type}-{action}-{target.type}
func (e *GridEvent) GetEventType() string {
	actorType := ""
	if e.Actor != nil {
		actorType = e.Actor.Type
	}
	targetType := ""
	if e.Target != nil {
		targetType = e.Target.Type
	}
	return actorType + "-" + e.Action + "-" + targetType
}

// GetActorState safely retrieves a value from actor.state
func (e *GridEvent) GetActorState(key string) interface{} {
	if e.Actor == nil || e.Actor.State == nil {
		return nil
	}
	return e.Actor.State[key]
}

// GetTargetState safely retrieves a value from target.state
func (e *GridEvent) GetTargetState(key string) interface{} {
	if e.Target == nil || e.Target.State == nil {
		return nil
	}
	return e.Target.State[key]
}

// GetActorGameState retrieves nested game state from actor
func (e *GridEvent) GetActorGameState() map[string]interface{} {
	if e.Actor == nil || e.Actor.State == nil {
		return nil
	}
	if game, ok := e.Actor.State["game"].(map[string]interface{}); ok {
		return game
	}
	return nil
}

// GetTargetGameState retrieves nested game state from target
func (e *GridEvent) GetTargetGameState() map[string]interface{} {
	if e.Target == nil || e.Target.State == nil {
		return nil
	}
	if game, ok := e.Target.State["game"].(map[string]interface{}); ok {
		return game
	}
	return nil
}

// GetPosition extracts position from game state
func (e *GridEvent) GetActorPosition() *Position {
	gameState := e.GetActorGameState()
	if gameState == nil {
		return nil
	}
	if pos, ok := gameState["position"].(map[string]interface{}); ok {
		x, _ := pos["x"].(float64)
		y, _ := pos["y"].(float64)
		return &Position{X: x, Y: y}
	}
	return nil
}

// GetTargetPosition extracts position from target game state
func (e *GridEvent) GetTargetPosition() *Position {
	gameState := e.GetTargetGameState()
	if gameState == nil {
		return nil
	}
	if pos, ok := gameState["position"].(map[string]interface{}); ok {
		x, _ := pos["x"].(float64)
		y, _ := pos["y"].(float64)
		return &Position{X: x, Y: y}
	}
	return nil
}

// GameEvent is kept for backward compatibility but deprecated
// Use EventWrapper and GridEvent instead
type GameEvent struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	GameTime  int                    `json:"gameTime,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// LoL-specific event types
type DragonKillEvent struct {
	TeamID     string    `json:"teamId"`
	TeamName   string    `json:"teamName"`
	PlayerID   string    `json:"playerId"`
	PlayerName string    `json:"playerName"`
	DragonType string    `json:"dragonType"` // infernal, mountain, ocean, cloud, chemtech, hextech, elder
	Position   *Position `json:"position,omitempty"`
	OccurredAt time.Time `json:"occurredAt"`
	GameTime   int       `json:"gameTime"` // milliseconds from game start (for backward compatibility)
}

type TowerDestroyEvent struct {
	TeamID     string    `json:"teamId"`
	TeamName   string    `json:"teamName"`
	PlayerID   string    `json:"playerId,omitempty"` // empty if team credit
	PlayerName string    `json:"playerName,omitempty"`
	TowerID    string    `json:"towerId"` // e.g., "red-turret-mid-2"
	Lane       string    `json:"lane"`    // top, mid, bot
	TowerNum   int       `json:"towerNum"`
	OccurredAt time.Time `json:"occurredAt"`
	GameTime   int       `json:"gameTime"` // milliseconds from game start
}

type KillEvent struct {
	KillerID       string    `json:"killerId"`
	KillerName     string    `json:"killerName"`
	KillerTeamID   string    `json:"killerTeamId"`
	KillerPosition *Position `json:"killerPosition,omitempty"`
	VictimID       string    `json:"victimId"`
	VictimName     string    `json:"victimName"`
	VictimTeamID   string    `json:"victimTeamId"`
	VictimPosition *Position `json:"victimPosition,omitempty"`
	AssistIDs      []string  `json:"assistIds"`
	FirstBlood     bool      `json:"firstBlood"`
	OccurredAt     time.Time `json:"occurredAt"`
	GameTime       int       `json:"gameTime"` // milliseconds from game start
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// VALORANT-specific event types
type RoundEndEvent struct {
	GameID      string    `json:"gameId"`
	RoundNum    int       `json:"roundNum"`
	WinnerTeam  string    `json:"winnerTeam"`
	WinnerName  string    `json:"winnerName"`
	WinType     string    `json:"winType"`     // "opponentEliminated", "defuse", "detonate", "time"
	AttackTeam  string    `json:"attackTeam"`  // For backward compatibility
	DefenseTeam string    `json:"defenseTeam"` // For backward compatibility
	OccurredAt  time.Time `json:"occurredAt"`
}

// IsDefendingSide returns true if the given teamID was defending this round
func (r *RoundEndEvent) IsDefendingSide(teamID string) bool {
	return r.DefenseTeam == teamID
}

// IsAttackingSide returns true if the given teamID was attacking this round
func (r *RoundEndEvent) IsAttackingSide(teamID string) bool {
	return r.AttackTeam == teamID
}

type PlantEvent struct {
	PlayerID   string    `json:"playerId"`
	PlayerName string    `json:"playerName"`
	TeamID     string    `json:"teamId"`
	Agent      string    `json:"agent"`
	Position   *Position `json:"position,omitempty"` // Used to infer site
	Site       string    `json:"site,omitempty"`     // Inferred from position: "A", "B", "C"
	RoundNum   int       `json:"roundNum"`
	MapName    string    `json:"mapName"`
	OccurredAt time.Time `json:"occurredAt"`
	GameTime   int       `json:"gameTime"` // milliseconds from round start
}

type DefuseEvent struct {
	PlayerID   string    `json:"playerId"`
	PlayerName string    `json:"playerName"`
	TeamID     string    `json:"teamId"`
	Agent      string    `json:"agent"`
	Position   *Position `json:"position,omitempty"`
	RoundNum   int       `json:"roundNum"`
	OccurredAt time.Time `json:"occurredAt"`
	GameTime   int       `json:"gameTime"` // milliseconds from round start
}

// VALKillEvent is VALORANT-specific kill event with additional context
type VALKillEvent struct {
	KillerID       string    `json:"killerId"`
	KillerName     string    `json:"killerName"`
	KillerTeamID   string    `json:"killerTeamId"`
	KillerAgent    string    `json:"killerAgent"`
	KillerPosition *Position `json:"killerPosition,omitempty"`
	VictimID       string    `json:"victimId"`
	VictimName     string    `json:"victimName"`
	VictimTeamID   string    `json:"victimTeamId"`
	VictimAgent    string    `json:"victimAgent"`
	VictimPosition *Position `json:"victimPosition,omitempty"`
	RoundNum       int       `json:"roundNum"`
	MapName        string    `json:"mapName"`
	OccurredAt     time.Time `json:"occurredAt"`
	GameTime       int       `json:"gameTime"` // milliseconds from round start
}

// DraftAction represents a pick or ban in draft phase
type DraftAction struct {
	TeamID        string    `json:"teamId"`
	TeamName      string    `json:"teamName"`
	Action        string    `json:"action"` // "picked" or "banned"
	CharacterName string    `json:"characterName"`
	CharacterID   string    `json:"characterId"`
	Sequence      int       `json:"sequence"`
	OccurredAt    time.Time `json:"occurredAt"`
}

// ObjectiveKillEvent represents Baron, Herald, or other major objective kills
type ObjectiveKillEvent struct {
	TeamID        string    `json:"teamId"`
	TeamName      string    `json:"teamName"`
	PlayerID      string    `json:"playerId"`
	PlayerName    string    `json:"playerName"`
	ObjectiveID   string    `json:"objectiveId"`   // e.g., "baron", "riftHerald", "voidGrub1"
	ObjectiveType string    `json:"objectiveType"` // "baron", "herald", "voidGrub"
	Position      *Position `json:"position,omitempty"`
	OccurredAt    time.Time `json:"occurredAt"`
	GameTime      int       `json:"gameTime"` // milliseconds from game start
}
