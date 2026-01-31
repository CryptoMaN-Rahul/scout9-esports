package intelligence

import "time"

// TeamAnalysis contains macro-level team insights
type TeamAnalysis struct {
	TeamID          string  `json:"teamId"`
	TeamName        string  `json:"teamName"`
	Title           string  `json:"title"` // "lol" or "valorant"
	MatchesAnalyzed int     `json:"matchesAnalyzed"`
	GamesAnalyzed   int     `json:"gamesAnalyzed"`
	WinRate         float64 `json:"winRate"`
	
	// Common metrics
	Strengths   []Insight `json:"strengths"`
	Weaknesses  []Insight `json:"weaknesses"`
	
	// LoL-specific metrics
	LoLMetrics *LoLTeamMetrics `json:"lolMetrics,omitempty"`
	
	// VALORANT-specific metrics
	VALMetrics *VALTeamMetrics `json:"valMetrics,omitempty"`
}

// Insight represents a data-backed observation
type Insight struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	SampleSize  int     `json:"sampleSize"`
	Comparison  string  `json:"comparison,omitempty"` // e.g., "above average", "below average"
}

// LoLTeamMetrics contains League of Legends specific team metrics
type LoLTeamMetrics struct {
	// Early game
	FirstBloodRate      float64 `json:"firstBloodRate"`
	FirstDragonRate     float64 `json:"firstDragonRate"`
	FirstTowerRate      float64 `json:"firstTowerRate"`
	FirstTowerAvgTime   float64 `json:"firstTowerAvgTime"` // minutes
	
	// Mid game
	GoldDiff15          float64 `json:"goldDiff15"`          // average gold diff at 15 min
	DragonControlRate   float64 `json:"dragonControlRate"`   // % of dragons taken
	HeraldControlRate   float64 `json:"heraldControlRate"`   // % of heralds taken
	
	// Late game
	BaronControlRate    float64 `json:"baronControlRate"`
	ElderDragonRate     float64 `json:"elderDragonRate"`
	AvgGameDuration     float64 `json:"avgGameDuration"` // minutes
	
	// Playstyle indicators
	EarlyGameRating     float64 `json:"earlyGameRating"`  // 0-100
	MidGameRating       float64 `json:"midGameRating"`    // 0-100
	LateGameRating      float64 `json:"lateGameRating"`   // 0-100
	AggressionScore     float64 `json:"aggressionScore"`  // 0-100
	
	// Win conditions
	WinConditions       []string `json:"winConditions"`
}

// VALTeamMetrics contains VALORANT specific team metrics
type VALTeamMetrics struct {
	// Overall
	AttackWinRate       float64            `json:"attackWinRate"`
	DefenseWinRate      float64            `json:"defenseWinRate"`
	
	// Pistol rounds
	PistolWinRate       float64            `json:"pistolWinRate"`
	AttackPistolWinRate float64            `json:"attackPistolWinRate"`
	DefensePistolWinRate float64           `json:"defensePistolWinRate"`
	
	// Economy
	EcoRoundWinRate     float64            `json:"ecoRoundWinRate"`
	ForceBuyWinRate     float64            `json:"forceBuyWinRate"`
	FullBuyWinRate      float64            `json:"fullBuyWinRate"`      // NEW: Full buy round win rate
	AvgTeamLoadout      float64            `json:"avgTeamLoadout"`      // NEW: Average team loadout value
	EconomyStats        *EconomyRoundStats `json:"economyStats,omitempty"` // NEW: Detailed economy stats
	
	// First blood
	FirstBloodRate      float64            `json:"firstBloodRate"`
	FirstDeathRate      float64            `json:"firstDeathRate"`
	
	// Map-specific
	MapStats            map[string]*MapStats `json:"mapStats"`
	MapPool             []MapPoolEntry       `json:"mapPool"`
	
	// Playstyle
	AggressionScore     float64            `json:"aggressionScore"`
	ClutchRate          float64            `json:"clutchRate"`
}

// MapStats contains per-map statistics for VALORANT
type MapStats struct {
	MapName         string  `json:"mapName"`
	GamesPlayed     int     `json:"gamesPlayed"`
	WinRate         float64 `json:"winRate"`
	AttackWinRate   float64 `json:"attackWinRate"`
	DefenseWinRate  float64 `json:"defenseWinRate"`
	AvgRoundsWon    float64 `json:"avgRoundsWon"`
	AvgRoundsLost   float64 `json:"avgRoundsLost"`
}

// MapPoolEntry represents a team's performance on a specific map
type MapPoolEntry struct {
	MapName     string  `json:"mapName"`
	GamesPlayed int     `json:"gamesPlayed"`
	WinRate     float64 `json:"winRate"`
	Strength    string  `json:"strength"` // "strong", "average", "weak"
}

// PlayerProfile contains individual player analysis
type PlayerProfile struct {
	PlayerID        string            `json:"playerId"`
	Nickname        string            `json:"nickname"`
	Role            string            `json:"role"`
	TeamID          string            `json:"teamId"`
	
	// Performance metrics
	GamesPlayed     int               `json:"gamesPlayed"`
	KDA             float64           `json:"kda"`
	AvgKills        float64           `json:"avgKills"`
	AvgDeaths       float64           `json:"avgDeaths"`
	AvgAssists      float64           `json:"avgAssists"`
	
	// Champion/Agent pool
	CharacterPool   []CharacterStats  `json:"characterPool"`
	SignaturePicks  []string          `json:"signaturePicks"`
	
	// Threat assessment
	ThreatLevel     int               `json:"threatLevel"` // 1-10
	ThreatReason    string            `json:"threatReason"`
	
	// Weaknesses
	Weaknesses      []Insight         `json:"weaknesses"`
	Tendencies      []string          `json:"tendencies"`
	
	// LoL-specific
	CSPerMin        float64           `json:"csPerMin,omitempty"`      // NOT available via GRID API
	GoldPerMin      float64           `json:"goldPerMin,omitempty"`
	DamageShare     float64           `json:"damageShare,omitempty"`   // NOT available via GRID API
	
	// VALORANT-specific
	ACS             float64           `json:"acs,omitempty"`      // Average Combat Score
	FirstBloodRate  float64           `json:"firstBloodRate,omitempty"`
	ClutchRate      float64           `json:"clutchRate,omitempty"`
	
	// NEW: Enhanced metrics from GRID API
	MultikillStats  *MultikillStats   `json:"multikillStats,omitempty"`  // Multi-kill breakdown
	WeaponStats     []WeaponStat      `json:"weaponStats,omitempty"`     // VALORANT weapon usage
	SynergyPartners []SynergyPartner  `json:"synergyPartners,omitempty"` // Players who assist them most
	AssistRatio     float64           `json:"assistRatio,omitempty"`     // Assists given vs received
	
	// NEW: Priority 2 enhancements
	AbilityUsage    []AbilityUsageStats   `json:"abilityUsage,omitempty"`    // Ability usage patterns
	ItemBuilds      []ItemBuildStats      `json:"itemBuilds,omitempty"`      // LoL item build patterns
	ObjectiveFocus  *ObjectiveFocusStats  `json:"objectiveFocus,omitempty"`  // Objective focus stats
}

// CharacterStats represents stats for a specific champion/agent
type CharacterStats struct {
	Character   string  `json:"character"`
	GamesPlayed int     `json:"gamesPlayed"`
	WinRate     float64 `json:"winRate"`
	KDA         float64 `json:"kda"`
	PickRate    float64 `json:"pickRate"` // % of games this character was picked
}

// MultikillStats tracks multi-kill performance
type MultikillStats struct {
	DoubleKills int `json:"doubleKills"`
	TripleKills int `json:"tripleKills"`
	QuadraKills int `json:"quadraKills"`
	PentaKills  int `json:"pentaKills"`
	TotalMultikills int `json:"totalMultikills"`
}

// WeaponStat tracks weapon usage for VALORANT
type WeaponStat struct {
	WeaponName string  `json:"weaponName"`
	Kills      int     `json:"kills"`
	KillShare  float64 `json:"killShare"` // % of total kills with this weapon
}

// SynergyPartner tracks which teammates assist a player most
type SynergyPartner struct {
	PlayerID     string  `json:"playerId"`
	PlayerName   string  `json:"playerName,omitempty"`
	AssistCount  int     `json:"assistCount"`
	SynergyScore float64 `json:"synergyScore"` // Normalized score
}

// AbilityUsageStats tracks ability usage patterns
type AbilityUsageStats struct {
	AbilityID    string  `json:"abilityId"`
	AbilityName  string  `json:"abilityName"`
	UsageCount   int     `json:"usageCount"`
	UsagePerGame float64 `json:"usagePerGame"`
}

// ItemBuildStats tracks item build patterns for LoL
type ItemBuildStats struct {
	ItemName     string  `json:"itemName"`
	ItemID       string  `json:"itemId"`
	BuildCount   int     `json:"buildCount"`   // Times this item was built
	BuildRate    float64 `json:"buildRate"`    // % of games this item was built
	AvgBuildTime float64 `json:"avgBuildTime"` // Average game time when built (if available)
}

// ObjectiveFocusStats tracks player objective focus (towers, dragons, etc.)
type ObjectiveFocusStats struct {
	TowersDestroyed    int     `json:"towersDestroyed"`
	TowersPerGame      float64 `json:"towersPerGame"`
	DragonsSecured     int     `json:"dragonsSecured"`
	DragonsPerGame     float64 `json:"dragonsPerGame"`
	BaronsSecured      int     `json:"baronsSecured"`
	HeraldsSecured     int     `json:"heraldsSecured"`
	ObjectiveFocused   bool    `json:"objectiveFocused"`   // > 1.5 towers/game or high objective participation
	ObjectiveFocusType string  `json:"objectiveFocusType"` // "split-pusher", "objective-focused", "teamfighter"
}

// EconomyRoundStats tracks economy round performance for VALORANT
type EconomyRoundStats struct {
	EcoRounds       int     `json:"ecoRounds"`       // Rounds with loadout < 2000
	EcoWins         int     `json:"ecoWins"`
	EcoWinRate      float64 `json:"ecoWinRate"`
	ForceRounds     int     `json:"forceRounds"`     // Rounds with loadout 2000-4000
	ForceWins       int     `json:"forceWins"`
	ForceWinRate    float64 `json:"forceWinRate"`
	FullBuyRounds   int     `json:"fullBuyRounds"`   // Rounds with loadout > 4000
	FullBuyWins     int     `json:"fullBuyWins"`
	FullBuyWinRate  float64 `json:"fullBuyWinRate"`
	AvgLoadoutValue float64 `json:"avgLoadoutValue"` // Average loadout value across all rounds
}

// CompositionAnalysis contains team composition insights
type CompositionAnalysis struct {
	TeamID              string              `json:"teamId"`
	Title               string              `json:"title"`
	
	// Most played compositions
	TopCompositions     []Composition       `json:"topCompositions"`
	
	// Draft patterns
	FirstPickPriorities []DraftPriority     `json:"firstPickPriorities"`
	CommonBans          []DraftPriority     `json:"commonBans"`
	FlexPicks           []string            `json:"flexPicks"`
	
	// Composition archetypes (LoL)
	ArchetypeBreakdown  map[string]float64  `json:"archetypeBreakdown,omitempty"`
	
	// Synergies
	PlayerSynergies     []Synergy           `json:"playerSynergies"`
}

// Composition represents a team composition
type Composition struct {
	Characters  []string `json:"characters"`
	GamesPlayed int      `json:"gamesPlayed"`
	WinRate     float64  `json:"winRate"`
	Frequency   float64  `json:"frequency"` // % of games
	Archetype   string   `json:"archetype,omitempty"` // e.g., "teamfight", "pick", "split-push"
}

// DraftPriority represents a draft priority (pick or ban)
type DraftPriority struct {
	Character   string  `json:"character"`
	Rate        float64 `json:"rate"` // pick/ban rate
	WinRate     float64 `json:"winRate,omitempty"`
	GamesPlayed int     `json:"gamesPlayed"`
}

// Synergy represents a player-character synergy
type Synergy struct {
	PlayerName  string  `json:"playerName"`
	Character   string  `json:"character"`
	GamesPlayed int     `json:"gamesPlayed"`
	WinRate     float64 `json:"winRate"`
	KDA         float64 `json:"kda"`
}

// CounterStrategy contains "How to Win" recommendations
type CounterStrategy struct {
	TeamID              string              `json:"teamId"`
	TeamName            string              `json:"teamName"`
	
	// Win condition summary
	WinCondition        string              `json:"winCondition"`
	
	// Exploitable weaknesses
	Weaknesses          []WeaknessTarget    `json:"weaknesses"`
	
	// Draft recommendations
	DraftRecommendations []DraftRecommendation `json:"draftRecommendations"`
	
	// In-game strategies
	InGameStrategies    []Strategy          `json:"inGameStrategies"`
	
	// Target players
	TargetPlayers       []PlayerTarget      `json:"targetPlayers"`
	
	// Confidence score
	ConfidenceScore     float64             `json:"confidenceScore"` // 0-100
	
	// Warnings (e.g., low sample size)
	Warnings            []string            `json:"warnings,omitempty"`
}

// WeaknessTarget represents an exploitable weakness
type WeaknessTarget struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Evidence    string  `json:"evidence"`
	Impact      float64 `json:"impact"` // 0-100 expected impact
}

// DraftRecommendation represents a draft recommendation
type DraftRecommendation struct {
	Type        string `json:"type"` // "ban", "pick", "target"
	Character   string `json:"character"`
	Reason      string `json:"reason"`
	Priority    int    `json:"priority"` // 1 = highest
}

// Strategy represents an in-game strategy recommendation
type Strategy struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Timing      string `json:"timing,omitempty"` // e.g., "early game", "post-15 min"
	Evidence    string `json:"evidence"`
}

// PlayerTarget represents a player to target
type PlayerTarget struct {
	PlayerName  string `json:"playerName"`
	Role        string `json:"role"`
	Reason      string `json:"reason"`
	Priority    int    `json:"priority"`
}

// TrendAnalysis contains trend and form analysis
type TrendAnalysis struct {
	TeamID          string       `json:"teamId"`
	
	// Win rate trends
	Last5WinRate    float64      `json:"last5WinRate"`
	Last10WinRate   float64      `json:"last10WinRate"`
	OverallWinRate  float64      `json:"overallWinRate"`
	
	// Form indicator
	FormIndicator   string       `json:"formIndicator"` // "hot", "stable", "cold"
	FormScore       float64      `json:"formScore"`     // -100 to 100
	
	// Trend changes
	TrendChanges    []TrendChange `json:"trendChanges"`
	
	// Metric trends
	MetricTrends    []MetricTrend `json:"metricTrends"`
}

// TrendChange represents a significant change in team performance
type TrendChange struct {
	Metric      string    `json:"metric"`
	OldValue    float64   `json:"oldValue"`
	NewValue    float64   `json:"newValue"`
	ChangeDate  time.Time `json:"changeDate"`
	Explanation string    `json:"explanation"`
}

// MetricTrend represents the trend of a specific metric
type MetricTrend struct {
	Metric    string    `json:"metric"`
	Direction string    `json:"direction"` // "improving", "stable", "declining"
	Values    []float64 `json:"values"`    // chronological values
}

// MatchupAnalysis contains head-to-head analysis
type MatchupAnalysis struct {
	Team1ID         string          `json:"team1Id"`
	Team1Name       string          `json:"team1Name"`
	Team2ID         string          `json:"team2Id"`
	Team2Name       string          `json:"team2Name"`
	
	// Head-to-head record
	TotalMatches    int             `json:"totalMatches"`
	Team1Wins       int             `json:"team1Wins"`
	Team2Wins       int             `json:"team2Wins"`
	
	// Historical patterns
	Patterns        []MatchupPattern `json:"patterns"`
	
	// Draft history
	DraftPatterns   []DraftPattern   `json:"draftPatterns"`
	
	// Recommendations
	Recommendations []string         `json:"recommendations"`
}

// MatchupPattern represents a pattern from historical matchups
type MatchupPattern struct {
	Description string  `json:"description"`
	Frequency   float64 `json:"frequency"`
	Impact      string  `json:"impact"`
}

// DraftPattern represents a draft pattern from historical matchups
type DraftPattern struct {
	Description string   `json:"description"`
	Characters  []string `json:"characters"`
	Outcome     string   `json:"outcome"`
}

// ScoutingReport is the complete output structure
type ScoutingReport struct {
	ID              string              `json:"id"`
	GeneratedAt     time.Time           `json:"generatedAt"`
	
	// Target team info
	OpponentTeam    TeamInfo            `json:"opponentTeam"`
	Title           string              `json:"title"` // "lol" or "valorant"
	MatchesAnalyzed int                 `json:"matchesAnalyzed"`
	
	// Report sections
	ExecutiveSummary    string              `json:"executiveSummary"`
	HowToWin            *CounterStrategy    `json:"howToWin"`
	TeamStrategy        *TeamAnalysis       `json:"teamStrategy"`
	PlayerProfiles      []*PlayerProfile    `json:"playerProfiles"`
	Compositions        *CompositionAnalysis `json:"compositions"`
	TrendAnalysis       *TrendAnalysis      `json:"trendAnalysis"`
	MatchupHistory      *MatchupAnalysis    `json:"matchupHistory,omitempty"`
}

// TeamInfo contains basic team information
type TeamInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
}


// =============================================================================
// HACKATHON WINNING TRANSFORMATION - NEW TYPES
// =============================================================================
// These types are designed to match the EXACT output format required by the
// Cloud9 x JetBrains Hackathon Category 2: Automated Scouting Report Generator
// =============================================================================

// DigestibleReport is the hackathon-compliant output format
// This matches the exact structure shown in the hackathon examples
type DigestibleReport struct {
	TeamName        string `json:"teamName"`
	MatchesAnalyzed int    `json:"matchesAnalyzed"`
	GeneratedAt     string `json:"generatedAt"`

	// Section 1: Executive Summary (1 paragraph)
	ExecutiveSummary string `json:"executiveSummary"`

	// Section 2: Common Strategies (hackathon format)
	CommonStrategies CommonStrategiesSection `json:"commonStrategies"`

	// Section 3: Player Tendencies (hackathon format)
	PlayerTendencies []PlayerTendencyInsight `json:"playerTendencies"`

	// Section 4: Recent Compositions (hackathon format)
	RecentCompositions []CompositionInsight `json:"recentCompositions"`

	// Section 5: How to Win (THE DIFFERENTIATOR)
	HowToWin HowToWinSection `json:"howToWin"`
}

// CommonStrategiesSection matches hackathon example format
type CommonStrategiesSection struct {
	// Attack patterns: "On Attack, 70% of pistol rounds are a 5-man fast-hit on B-Site (Ascent)"
	AttackPatterns []StrategyInsight `json:"attackPatterns"`

	// Defense setups: "On Defense, they default to a 1-3-1 setup, rotating their Sentinel to mid"
	DefenseSetups []StrategyInsight `json:"defenseSetups"`

	// Objective priorities (LoL): "Prioritizes first Drake (82% contest rate)"
	ObjectivePriorities []StrategyInsight `json:"objectivePriorities"`

	// Timing patterns: "4-man group for first tower push, usually in bot lane at ~13 mins"
	TimingPatterns []StrategyInsight `json:"timingPatterns"`
}

// StrategyInsight is a single strategy observation in hackathon format
type StrategyInsight struct {
	// The insight text exactly as it should appear in the report
	// Example: "On Attack, 70% of pistol rounds are a 5-man fast-hit on B-Site (Ascent)"
	Text string `json:"text"`

	// Supporting data
	Metric     string  `json:"metric"`     // e.g., "pistol_attack_pattern"
	Value      float64 `json:"value"`      // e.g., 0.70
	SampleSize int     `json:"sampleSize"` // e.g., 10 games
	Context    string  `json:"context"`    // e.g., "Ascent", "bot lane"
}

// PlayerTendencyInsight matches hackathon example format
// Example: "Player 'Jett' has a 75% first-duel rate with an Operator on A-main defense"
type PlayerTendencyInsight struct {
	// The insight text exactly as it should appear
	Text string `json:"text"`

	// Player info
	PlayerName string `json:"playerName"`
	Role       string `json:"role"`

	// The specific tendency
	TendencyType string  `json:"tendencyType"` // "first_duel", "champion_pool", "matchup"
	Value        float64 `json:"value"`
	Context      string  `json:"context"` // "with Operator on A-main defense"
	SampleSize   int     `json:"sampleSize"`
}

// CompositionInsight matches hackathon example format
// Example: "Most-played comp (68% on Split): Jett, Raze, Brimstone, Skye, Cypher"
type CompositionInsight struct {
	// The insight text exactly as it should appear
	Text string `json:"text"`

	// Composition details
	Characters  []string `json:"characters"`
	Frequency   float64  `json:"frequency"` // e.g., 0.68
	WinRate     float64  `json:"winRate"`
	MapContext  string   `json:"mapContext,omitempty"` // e.g., "Split"
	Archetype   string   `json:"archetype,omitempty"`  // e.g., "late-game, team-fight-oriented"
	GamesPlayed int      `json:"gamesPlayed"`
}

// HowToWinSection is THE KEY DIFFERENTIATOR
// This provides specific, actionable, data-backed recommendations
type HowToWinSection struct {
	// One-line win condition
	// Example: "This team's win condition is their bot lane. Their jungler paths to bot 75% of the time pre-10 mins."
	WinCondition string `json:"winCondition"`

	// Specific actionable insights
	// Example: "Opponent X has a 15% win rate when their mid-laner is against an assassin. Recommend picking LeBlanc and Zed."
	ActionableInsights []ActionableInsight `json:"actionableInsights"`

	// Draft strategy
	DraftStrategy DraftStrategySection `json:"draftStrategy"`

	// In-game strategy
	InGameStrategy []InGameStrategyInsight `json:"inGameStrategy"`

	// Confidence score
	ConfidenceScore float64 `json:"confidenceScore"`
}

// ActionableInsight is a specific, data-backed recommendation
type ActionableInsight struct {
	// The recommendation
	// Example: "Target their mid-laner when he's on control mages"
	Recommendation string `json:"recommendation"`

	// The data backing
	// Example: "He has 2.1 KDA on control mages vs 8.5 on assassins"
	DataBacking string `json:"dataBacking"`

	// Impact level
	Impact string `json:"impact"` // "HIGH", "MEDIUM", "LOW"

	// Action type
	ActionType string `json:"actionType"` // "BAN", "PICK", "TARGET_PLAYER", "FORCE_MAP", "STRATEGY"

	// Confidence
	Confidence float64 `json:"confidence"`
}

// DraftStrategySection contains specific draft recommendations
type DraftStrategySection struct {
	// Priority bans with reasons
	// Example: "Ban Jett - their duelist 'Aspas' has 78% win rate on Jett"
	PriorityBans []DraftInsight `json:"priorityBans"`

	// Recommended picks
	// Example: "Pick assassins mid - opponent mid-laner has 15% win rate vs assassins"
	RecommendedPicks []DraftInsight `json:"recommendedPicks"`

	// Target picks (force opponent onto weak picks)
	// Example: "Force their top-laner onto tanks - 42% win rate vs 68% on carries"
	TargetPicks []DraftInsight `json:"targetPicks"`
}

// DraftInsight is a single draft recommendation
type DraftInsight struct {
	// The recommendation text
	Text string `json:"text"`

	// Character involved
	Character string `json:"character"`

	// Player involved (if targeting specific player)
	PlayerName string `json:"playerName,omitempty"`

	// Data backing
	WinRate    float64 `json:"winRate"`
	SampleSize int     `json:"sampleSize"`

	// Priority (1 = highest)
	Priority int `json:"priority"`
}

// InGameStrategyInsight is a specific in-game recommendation
type InGameStrategyInsight struct {
	// The strategy
	// Example: "Aggressive counter-jungling on their top side"
	Strategy string `json:"strategy"`

	// Timing
	// Example: "Pre-10 minutes"
	Timing string `json:"timing"`

	// Reason
	// Example: "Their jungler paths to bot 75% of the time, leaving top side vulnerable"
	Reason string `json:"reason"`

	// Expected impact
	Impact string `json:"impact"`
}

// =============================================================================
// MATCHUP ANALYSIS TYPES
// =============================================================================

// PlayerMatchupProfile tracks how a player performs against specific champions/agents
type PlayerMatchupProfile struct {
	PlayerID   string         `json:"playerId"`
	PlayerName string         `json:"playerName"`
	Role       string         `json:"role"`
	Matchups   []MatchupStats `json:"matchups"`

	// Aggregated insights
	StrongAgainst []MatchupStats `json:"strongAgainst"` // Characters they dominate
	WeakAgainst   []MatchupStats `json:"weakAgainst"`   // Characters they struggle against
}

// MatchupStats tracks performance in a specific matchup
type MatchupStats struct {
	// The character the player was playing
	PlayedCharacter string `json:"playedCharacter"`

	// The opponent character (what they were against)
	VsCharacter string `json:"vsCharacter"`

	// Performance
	GamesPlayed int     `json:"gamesPlayed"`
	WinRate     float64 `json:"winRate"`
	KDA         float64 `json:"kda"`
	AvgKills    float64 `json:"avgKills"`
	AvgDeaths   float64 `json:"avgDeaths"`

	// Classification
	MatchupType string `json:"matchupType"` // "favorable", "even", "unfavorable"
}

// =============================================================================
// VALORANT SITE-SPECIFIC ANALYSIS
// =============================================================================

// SiteAnalysis contains site-specific performance data for VALORANT
type SiteAnalysis struct {
	MapName string                  `json:"mapName"`
	Sites   map[string]*SiteStats   `json:"sites"` // "A", "B", "C", "Mid"
}

// SiteStats contains performance data for a specific site
type SiteStats struct {
	SiteName string `json:"siteName"`

	// Attack performance on this site
	AttackAttempts   int     `json:"attackAttempts"`
	AttackSuccesses  int     `json:"attackSuccesses"`
	AttackWinRate    float64 `json:"attackWinRate"`

	// Defense performance on this site
	DefenseAttempts  int     `json:"defenseAttempts"`
	DefenseSuccesses int     `json:"defenseSuccesses"`
	DefenseWinRate   float64 `json:"defenseWinRate"`

	// Retake performance
	RetakeAttempts   int     `json:"retakeAttempts"`
	RetakeSuccesses  int     `json:"retakeSuccesses"`
	RetakeWinRate    float64 `json:"retakeWinRate"`

	// Common patterns
	CommonAttackPatterns  []AttackPattern  `json:"commonAttackPatterns"`
	CommonDefenseSetups   []DefenseSetup   `json:"commonDefenseSetups"`
}

// AttackPattern describes a common attack pattern
type AttackPattern struct {
	// Description: "5-man fast-hit", "split A-B", "slow default"
	Description string `json:"description"`

	// Target site
	TargetSite string `json:"targetSite"`

	// Frequency and success
	Frequency   float64 `json:"frequency"`   // How often they use this
	SuccessRate float64 `json:"successRate"` // Win rate when using this

	// Timing
	AvgExecuteTime float64 `json:"avgExecuteTime"` // Seconds into round
}

// DefenseSetup describes a common defensive setup
type DefenseSetup struct {
	// Description: "1-3-1", "2-1-2", "stack A"
	Description string `json:"description"`

	// Player positions
	Positions map[string]string `json:"positions"` // Role -> Position

	// Frequency and success
	Frequency   float64 `json:"frequency"`
	SuccessRate float64 `json:"successRate"`
}

// =============================================================================
// LOL TIMING-BASED ANALYSIS
// =============================================================================

// TimingAnalysis contains timing-based patterns for LoL
type TimingAnalysis struct {
	TeamID string `json:"teamId"`

	// Jungle pathing
	JunglePathing JunglePathingAnalysis `json:"junglePathing"`

	// Objective timings
	ObjectiveTimings ObjectiveTimingAnalysis `json:"objectiveTimings"`

	// Lane patterns
	LanePatterns LanePatternAnalysis `json:"lanePatterns"`
}

// JunglePathingAnalysis tracks jungle pathing patterns
type JunglePathingAnalysis struct {
	// Pre-6 minute pathing preference
	PreSixPathPreference string  `json:"preSixPathPreference"` // "bot-focused", "top-focused", "vertical"
	PreSixBotRate        float64 `json:"preSixBotRate"`        // % of time pathing bot pre-6
	PreSixTopRate        float64 `json:"preSixTopRate"`        // % of time pathing top pre-6

	// Gank patterns
	GanksByLane map[string]float64 `json:"ganksByLane"` // Lane -> gank frequency

	// First clear patterns
	FirstClearPatterns []ClearPattern `json:"firstClearPatterns"`

	// Counter-jungle frequency
	CounterJungleRate float64 `json:"counterJungleRate"`
}

// ClearPattern describes a jungle clear pattern
type ClearPattern struct {
	Path        string  `json:"path"`        // e.g., "red-krugs-raptors-wolves-blue-gromp"
	Frequency   float64 `json:"frequency"`   // How often this path is used
	EndLocation string  `json:"endLocation"` // Where they end up
	EndTime     float64 `json:"endTime"`     // Minutes when clear completes
}

// ObjectiveTimingAnalysis tracks objective timing patterns
type ObjectiveTimingAnalysis struct {
	// Dragon
	FirstDragonAvgTime   float64 `json:"firstDragonAvgTime"`   // Minutes
	FirstDragonContestRate float64 `json:"firstDragonContestRate"`

	// Herald
	FirstHeraldAvgTime   float64 `json:"firstHeraldAvgTime"`
	HeraldUsagePattern   string  `json:"heraldUsagePattern"` // "top", "mid", "bot"

	// Tower
	FirstTowerAvgTime    float64 `json:"firstTowerAvgTime"`
	FirstTowerLane       string  `json:"firstTowerLane"` // Most common lane

	// Baron
	BaronAttemptTimings  []float64 `json:"baronAttemptTimings"` // List of attempt times
}

// LanePatternAnalysis tracks lane-specific patterns
type LanePatternAnalysis struct {
	// Lane priority patterns
	LanePriority map[string]float64 `json:"lanePriority"` // Lane -> priority score

	// Roaming patterns
	MidRoamRate    float64 `json:"midRoamRate"`    // How often mid roams
	SupportRoamRate float64 `json:"supportRoamRate"` // How often support roams

	// Grouping patterns
	FirstGroupTime float64 `json:"firstGroupTime"` // When they first group (minutes)
	GroupLocation  string  `json:"groupLocation"`  // Where they typically group
}

// =============================================================================
// ENHANCED COUNTER-STRATEGY TYPES
// =============================================================================
// These types support the enhanced "How to Win" insights with specific,
// actionable recommendations based on champion class matchup analysis.
// =============================================================================

// ClassMatchupInsight represents a class-based matchup insight
// Example: "Faker has 2.1 KDA on control mages vs 8.5 on assassins"
type ClassMatchupInsight struct {
	PlayerName     string  `json:"playerName"`
	Role           string  `json:"role"`
	WeakClass      string  `json:"weakClass"`
	WeakClassKDA   float64 `json:"weakClassKda"`
	WeakClassWR    float64 `json:"weakClassWinRate"`
	StrongClass    string  `json:"strongClass"`
	StrongClassKDA float64 `json:"strongClassKda"`
	StrongClassWR  float64 `json:"strongClassWinRate"`
	Text           string  `json:"text"` // Hackathon format text
	SampleSize     int     `json:"sampleSize"`
	Confidence     float64 `json:"confidence"`
}

// SpecificDraftRecommendation includes specific champion names, not just class names
// Example: "Pick Zed, LeBlanc, Akali vs Faker - they have 2.1 KDA on mages"
type SpecificDraftRecommendation struct {
	Type         string   `json:"type"`         // "pick", "ban", "target"
	TargetPlayer string   `json:"targetPlayer"` // Player being targeted
	TargetRole   string   `json:"targetRole"`   // Role of target player
	TargetClass  string   `json:"targetClass"`  // Class they're weak on
	Champions    []string `json:"champions"`    // Specific champion names to pick/ban
	Reason       string   `json:"reason"`       // Full explanation
	WinRate      float64  `json:"winRate"`      // Target's win rate on weak class
	KDA          float64  `json:"kda"`          // Target's KDA on weak class
	SampleSize   int      `json:"sampleSize"`   // Games analyzed
	Priority     int      `json:"priority"`     // 1 = highest priority
	Confidence   float64  `json:"confidence"`   // Confidence in this recommendation
}

// EconomyAnalysis contains VALORANT economy round analysis
type EconomyAnalysis struct {
	TeamID          string                     `json:"teamId"`
	EcoRoundWinRate float64                    `json:"ecoRoundWinRate"`
	EcoRounds       int                        `json:"ecoRounds"`
	ForceWinRate    float64                    `json:"forceWinRate"`
	ForceRounds     int                        `json:"forceRounds"`
	FullBuyWinRate  float64                    `json:"fullBuyWinRate"`
	FullBuyRounds   int                        `json:"fullBuyRounds"`
	ByMap           map[string]*MapEconomyStats `json:"byMap,omitempty"`
}

// MapEconomyStats contains per-map economy statistics
type MapEconomyStats struct {
	MapName        string  `json:"mapName"`
	EcoWinRate     float64 `json:"ecoWinRate"`
	EcoRounds      int     `json:"ecoRounds"`
	ForceWinRate   float64 `json:"forceWinRate"`
	ForceRounds    int     `json:"forceRounds"`
	FullBuyWinRate float64 `json:"fullBuyWinRate"`
	FullBuyRounds  int     `json:"fullBuyRounds"`
}

// EconomyInsight is a hackathon-format economy insight
// Example: "They only win 15% of eco rounds - play aggressive on their saves"
type EconomyInsight struct {
	Text       string  `json:"text"`
	Type       string  `json:"type"` // "eco_weak", "force_strong", "force_weak", "full_buy_weak"
	Value      float64 `json:"value"`
	SampleSize int     `json:"sampleSize"`
	Impact     string  `json:"impact"` // "HIGH", "MEDIUM", "LOW"
	MapContext string  `json:"mapContext,omitempty"`
}

// HeadToHeadReport contains team comparison for /api/matchups endpoint
type HeadToHeadReport struct {
	Team1ID   string `json:"team1Id"`
	Team1Name string `json:"team1Name"`
	Team2ID   string `json:"team2Id"`
	Team2Name string `json:"team2Name"`
	Title     string `json:"title"` // "lol" or "valorant"

	// Historical record
	TotalMatches int `json:"totalMatches"`
	Team1Wins    int `json:"team1Wins"`
	Team2Wins    int `json:"team2Wins"`

	// Style comparison
	StyleComparison *StyleComparison `json:"styleComparison"`

	// Key insights
	Insights []HeadToHeadInsight `json:"insights"`

	// Draft patterns from historical matches
	DraftPatterns []DraftPattern `json:"draftPatterns,omitempty"`

	// Confidence
	ConfidenceScore float64 `json:"confidenceScore"`
	
	// Warnings
	Warnings []string `json:"warnings,omitempty"`
}

// StyleComparison compares team styles across game phases
type StyleComparison struct {
	// Early game comparison
	Team1EarlyGameRating float64 `json:"team1EarlyGameRating"`
	Team2EarlyGameRating float64 `json:"team2EarlyGameRating"`
	EarlyGameAdvantage   string  `json:"earlyGameAdvantage"` // "team1", "team2", "even"
	EarlyGameInsight     string  `json:"earlyGameInsight"`   // "Team A wins early game 65% vs Team B's 40%"

	// Mid game comparison
	Team1MidGameRating float64 `json:"team1MidGameRating"`
	Team2MidGameRating float64 `json:"team2MidGameRating"`
	MidGameAdvantage   string  `json:"midGameAdvantage"`
	MidGameInsight     string  `json:"midGameInsight"`

	// Late game comparison
	Team1LateGameRating float64 `json:"team1LateGameRating"`
	Team2LateGameRating float64 `json:"team2LateGameRating"`
	LateGameAdvantage   string  `json:"lateGameAdvantage"`
	LateGameInsight     string  `json:"lateGameInsight"`

	// Aggression comparison
	Team1Aggression float64 `json:"team1Aggression"`
	Team2Aggression float64 `json:"team2Aggression"`

	// Overall style insight
	StyleInsight string `json:"styleInsight"`
}

// HeadToHeadInsight is a specific insight from head-to-head analysis
type HeadToHeadInsight struct {
	Text       string  `json:"text"`
	Type       string  `json:"type"` // "early_game", "draft", "player_matchup", "historical"
	Confidence float64 `json:"confidence"`
}

// InsightConfidence tracks confidence for individual insights
type InsightConfidence struct {
	InsightID  string  `json:"insightId"`
	Confidence float64 `json:"confidence"`
	SampleSize int     `json:"sampleSize"`
	Warning    string  `json:"warning,omitempty"` // "Low sample size" etc.
}

// CounterClassRecommendation maps weak classes to counter picks
type CounterClassRecommendation struct {
	WeakClass    string   `json:"weakClass"`
	CounterClass string   `json:"counterClass"`
	Champions    []string `json:"champions"`
}
