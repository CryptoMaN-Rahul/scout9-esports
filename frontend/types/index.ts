// Scout9 Frontend Types
// Matches backend Go types from pkg/intelligence/types.go and pkg/report/formatter.go

// === Core Types ===
export type GameTitle = 'lol' | 'valorant';

export interface Title {
  id: string;
  name: string;
}

export interface Tournament {
  id: string;
  name: string;
  logoUrl?: string;
  startDate?: string;
  endDate?: string;
  titleId: string;
}

export interface Team {
  id: string;
  name: string;
  logoUrl?: string;
  colorHex?: string;
}

export interface Player {
  id: string;
  nickname: string;
  role?: string;
  teamId?: string;
}

export interface Series {
  id: string;
  tournamentId: string;
  startTime: string;
  format?: string;
  teams: Team[];
  titleId: string;
}

// === Report Types ===
export interface ScoutingReport {
  id: string;
  generatedAt: string;
  opponentTeam: Team;
  title: GameTitle;
  matchesAnalyzed: number;
  executiveSummary: string;
  howToWin: HowToWinSection;
  teamStrategy: TeamAnalysis;
  playerProfiles: PlayerProfile[];
  compositions: CompositionAnalysis;
  trendAnalysis?: TrendAnalysis;
}

export interface DigestibleReport {
  teamName: string;
  matchesAnalyzed: number;
  generatedAt: string;
  executiveSummary: string;
  commonStrategies: CommonStrategiesSection;
  playerTendencies: PlayerTendencyInsight[];
  recentCompositions: CompositionInsight[];
  howToWin: HowToWinSection;
}

export interface ReportListItem {
  id: string;
  teamName: string;
  title: GameTitle;
  generatedAt: string;
  matchesAnalyzed: number;
}

// === How to Win (THE DIFFERENTIATOR) ===
// This matches the backend's CounterStrategy type
export interface HowToWinSection {
  teamId?: string;
  teamName?: string;
  winCondition: string;
  weaknesses?: WeaknessTarget[];
  draftRecommendations?: DraftRecommendation[];
  inGameStrategies?: InGameStrategy[];
  targetPlayers?: PlayerTarget[];
  confidenceScore: number;
  warnings?: string[];
  // Legacy fields for backward compatibility
  actionableInsights?: ActionableInsight[];
  draftStrategy?: DraftStrategySection;
  inGameStrategy?: InGameStrategyInsight[];
  sampleSizeWarning?: string;
}

export interface WeaknessTarget {
  title: string;
  description: string;
  evidence?: string;
  impact: number;
}

export interface DraftRecommendation {
  type: string; // "ban", "pick", "target"
  character: string;
  reason: string;
  priority: number;
}

export interface InGameStrategy {
  title: string;
  description: string;
  timing?: string;
  evidence?: string;
}

export interface PlayerTarget {
  playerName: string;
  role: string;
  reason: string;
  priority: number;
}

export type ImpactLevel = 'HIGH' | 'MEDIUM' | 'LOW';
export type ActionType = 'BAN' | 'PICK' | 'TARGET_PLAYER' | 'FORCE_MAP' | 'STRATEGY';

export interface ActionableInsight {
  recommendation: string;
  dataBacking: string;
  impact: ImpactLevel;
  actionType: ActionType;
  confidence: number;
}

export interface DraftStrategySection {
  priorityBans: DraftInsight[];
  recommendedPicks: DraftInsight[];
  targetPicks: DraftInsight[];
}

export interface DraftInsight {
  character: string;
  reason: string;
  playerName?: string;
  winRate?: number;
  pickRate?: number;
}

export interface InGameStrategyInsight {
  phase: string;
  timing: string;
  strategy: string;
  priority: ImpactLevel;
}

// === Team Analysis ===
export interface TeamAnalysis {
  teamId: string;
  teamName: string;
  title: GameTitle;
  matchesAnalyzed: number;
  gamesAnalyzed: number;
  winRate: number;
  recentForm: string;
  strengths: Insight[];
  weaknesses: Insight[];
  lolMetrics?: LoLTeamMetrics;
  valMetrics?: VALTeamMetrics;
}

export interface LoLTeamMetrics {
  // Early game
  firstBloodRate: number;
  firstDragonRate: number;
  firstTowerRate: number;
  firstTowerAvgTime: number;
  
  // Mid game
  goldDiff15: number;
  dragonControlRate: number;
  heraldControlRate: number;
  
  // Late game
  baronControlRate: number;
  elderDragonRate: number;
  avgGameDuration: number;
  
  // Playstyle ratings
  earlyGameRating: number;
  midGameRating: number;
  lateGameRating: number;
  aggressionScore: number;
  winConditions: string[];
}

export interface VALTeamMetrics {
  attackWinRate: number;
  defenseWinRate: number;
  pistolWinRate: number;
  attackPistolWinRate: number;
  defensePistolWinRate: number;
  ecoRoundWinRate: number;
  forceBuyWinRate: number;
  fullBuyWinRate: number;
  avgTeamLoadout?: number;
  economyStats?: EconomyRoundStats;
  firstBloodRate: number;
  firstDeathRate?: number;
  mapStats: Record<string, MapStats>;
  mapPool: MapPoolEntry[];
  aggressionScore: number;
  clutchRate: number;
}

export interface EconomyRoundStats {
  // Object-style (original definition)
  ecoRounds?: { won: number; total: number; winRate: number } | number;
  forceRounds?: { won: number; total: number; winRate: number } | number;
  fullBuyRounds?: { won: number; total: number; winRate: number } | number;
  // Flat-style (actual backend format)
  ecoWins?: number;
  ecoWinRate?: number;
  forceWins?: number;
  forceWinRate?: number;
  fullBuyWins?: number;
  fullBuyWinRate?: number;
  avgLoadoutValue?: number;
}

// === Player Profile ===
export interface PlayerProfile {
  playerId: string;
  nickname: string;
  role: string;
  teamId: string;
  gamesPlayed: number;
  
  // Performance
  kda: number;
  avgKills: number;
  avgDeaths: number;
  avgAssists: number;
  
  // Champion/Agent pool
  characterPool: CharacterStats[];
  signaturePicks: string[];
  
  // Threat assessment
  threatLevel: number;
  threatReason: string;
  
  weaknesses: Insight[];
  tendencies: string[];
  
  // Enhanced metrics
  multikillStats?: MultikillStats;
  weaponStats?: WeaponStat[];
  synergyPartners?: SynergyPartner[];
  assistRatio?: number;
  abilityUsage?: AbilityUsageStats[];
  itemBuilds?: ItemBuildStats[];
  objectiveFocus?: ObjectiveFocusStats;
}

export interface CharacterStats {
  name?: string;
  character?: string;
  games?: number;
  gamesPlayed?: number;
  wins?: number;
  losses?: number;
  winRate: number;
  avgKDA?: number;
  kda?: number;
  pickRate?: number;
}

export interface MultikillStats {
  doubles?: number;
  triples?: number;
  quadras?: number;
  pentas?: number;
  avgPerGame?: number;
  // Backend field names (alternative)
  doubleKills?: number;
  tripleKills?: number;
  quadraKills?: number;
  pentaKills?: number;
  totalMultikills?: number;
}

export interface WeaponStat {
  weapon?: string;
  weaponName?: string;
  kills: number;
  percentage?: number;
  killShare?: number;
}

export interface SynergyPartner {
  playerId: string;
  playerName: string;
  assistsGiven: number;
  assistsReceived: number;
  synergyScore: number;
}

export interface AbilityUsageStats {
  ability?: string;
  abilityId?: string;
  abilityName?: string;
  uses?: number;
  usageCount?: number;
  usagePerGame?: number;
  kills?: number;
}

export interface ItemBuildStats {
  item: string;
  frequency: number;
  winRate: number;
}

export interface ObjectiveFocusStats {
  splitPushScore: number;
  teamfightScore: number;
  objectiveParticipation: number;
}

// === Common Strategies ===
export interface CommonStrategiesSection {
  attackPatterns: StrategyInsight[];
  defenseSetups: StrategyInsight[];
  objectivePriorities: StrategyInsight[];
  timingPatterns: StrategyInsight[];
}

export interface StrategyInsight {
  text: string;
  metric: string;
  value: number;
  sampleSize: number;
  context: string;
}

// === Compositions ===
export interface CompositionAnalysis {
  teamId?: string;
  title?: string;
  topCompositions: CompositionInsight[];
  firstPickPriorities: FirstPickPriority[];
  commonBans: string[];
  flexPicks?: string[];
  archetypeBreakdown?: Record<string, number>;
  playerSynergies?: PlayerSynergyItem[];
}

export interface FirstPickPriority {
  character: string;
  rate: number;
  winRate: number;
  gamesPlayed: number;
}

export interface PlayerSynergyItem {
  playerName: string;
  character: string;
  gamesPlayed: number;
  winRate: number;
  kda: number;
}

export interface CompositionInsight {
  characters: string[];
  frequency: number;
  winRate: number;
  gamesPlayed: number;
  archetype?: string;
}

// === Trends ===
export interface TrendAnalysis {
  formTrend: 'improving' | 'declining' | 'stable';
  recentResults: GameResult[];
  winRateTrend: number[];
  kdaTrend: number[];
}

export interface GameResult {
  date: string;
  opponent: string;
  won: boolean;
  score: string;
}

// === Head to Head Data (for Compare page) ===
export interface HeadToHeadData {
  team1: Team;
  team2: Team;
  matchHistory?: {
    matchId: string;
    date: string;
    winner: string;
    score: string;
    tournamentName: string;
    games: {
      gameNumber: number;
      winnerId: string;
      duration: number;
    }[];
  }[];
  stats: {
    totalMatches: number;
    team1Wins: number;
    team2Wins: number;
    avgGameDuration?: number;
    commonPicks?: {
      team1: string[];
      team2: string[];
    };
    keyMatchups?: {
      player1: { playerId: string; nickname: string; role: string; winRate: number; avgKDA: number };
      player2: { playerId: string; nickname: string; role: string; winRate: number; avgKDA: number };
      gamesPlayed: number;
      significance: string;
    }[];
  };
  styleComparison?: StyleComparisonData;
  insights?: string[];
  warnings?: string[];
  confidenceScore?: number;
}

export interface StyleComparisonData {
  team1EarlyGameRating?: number;
  team2EarlyGameRating?: number;
  earlyGameAdvantage?: string;
  earlyGameInsight?: string;
  team1MidGameRating?: number;
  team2MidGameRating?: number;
  midGameAdvantage?: string;
  midGameInsight?: string;
  team1LateGameRating?: number;
  team2LateGameRating?: number;
  lateGameAdvantage?: string;
  lateGameInsight?: string;
  team1Aggression?: number;
  team2Aggression?: number;
  styleInsight?: string;
}

// === Head to Head Analysis (from API) ===
export interface HeadToHeadAnalysis {
  team1: TeamAnalysis;
  team2: TeamAnalysis;
  historicalRecord: {
    team1Wins: number;
    team2Wins: number;
    lastMet?: string;
  };
  styleComparison: {
    earlyGame: { team1: number; team2: number };
    midGame: { team1: number; team2: number };
    lateGame: { team1: number; team2: number };
    aggression: { team1: number; team2: number };
  };
  insights: string[];
  confidenceScore: number;
}

// === Utility Types ===
export interface Insight {
  text?: string;
  title?: string;
  description?: string;
  value?: number;
  sampleSize?: number;
  importance?: ImpactLevel;
  dataPoints?: string[];
}

export interface MapStats {
  mapName: string;
  gamesPlayed: number;
  wins: number;
  winRate: number;
  attackWinRate: number;
  defenseWinRate: number;
}

export interface MapPoolEntry {
  // Frontend format
  map?: string;
  games?: number;
  comfort?: number;
  // Backend format
  mapName?: string;
  gamesPlayed?: number;
  strength?: string; // "strong", "average", "weak"
  // Common
  winRate: number;
}

export interface PlayerTendencyInsight {
  playerName: string;
  role: string;
  insights: string[];
  signaturePick?: { character: string; pickRate: number; winRate: number };
  kda: number;
  threatLevel: number;
}

// === API Request/Response Types ===
export interface GenerateReportRequest {
  teamId: string;
  teamName?: string;
  matchCount?: number;
  titleId?: string;
}

export interface ApiError {
  error: string;
  message?: string;
  code?: number;
}
