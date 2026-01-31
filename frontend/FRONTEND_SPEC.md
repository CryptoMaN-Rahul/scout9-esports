# Scout9 Frontend Specification

> **Purpose**: This spec file contains all details needed to continue frontend development without re-analyzing the codebase. Use this as the primary reference for any AI agent continuing this work.

**Last Updated**: January 31, 2026  
**Status**: 🚧 In Development

---

## Quick Context

**Project**: Scout9 - Automated Scouting Report Generator for Esports  
**Hackathon**: Cloud9 x JetBrains - Category 2  
**Prize**: $6,000 + GDC trip  
**Backend**: Go REST API at `http://localhost:8080`  
**Frontend**: Next.js 14 + TypeScript + Tailwind CSS + shadcn/ui

---

## Progress Tracker

| Task | Status | Notes |
|------|--------|-------|
| **1. Project Setup** | ✅ Complete | Next.js 14 + Tailwind + deps installed |
| **2. Type Definitions** | ✅ Complete | types/index.ts matches backend |
| **3. API Client** | ✅ Complete | lib/api.ts with error handling |
| **4. Layout & Navigation** | ✅ Complete | Header, Footer, GameToggle |
| **5. Landing Page** | ✅ Complete | app/page.tsx with team search |
| **6. Team Selection** | ✅ Complete | Search + demo fallback |
| **7. Report Generation UI** | ✅ Complete | LoL + VALORANT generate pages |
| **8. Report Dashboard** | ✅ Complete | Report pages for both games |
| **9. Executive Summary** | ✅ Complete | ExecutiveSummary.tsx |
| **10. Common Strategies** | ✅ Complete | CommonStrategies.tsx |
| **11. Player Tendencies** | ✅ Complete | PlayerTendencies.tsx |
| **12. Compositions** | ✅ Complete | Compositions.tsx |
| **13. How to Win Section** | ✅ Complete | HowToWin.tsx - KEY DIFFERENTIATOR |
| **14. Data Visualizations** | ✅ Complete | RadarChart, BarChart, PieChart, LineChart |
| **15. Head-to-Head Page** | ✅ Complete | app/compare/page.tsx |
| **16. PDF Export** | 🔄 In Progress | Print CSS added, button hooked up |
| **17. Game Theming** | ✅ Complete | CSS variables for LoL/VALORANT |
| **18. Animations** | ✅ Complete | Framer Motion throughout |
| **19. Responsive Design** | ✅ Complete | Mobile tabs, responsive grids |
| **20. Final Polish** | 🔄 In Progress | Testing needed |

**Legend**: ⬜ Not Started | 🔄 In Progress | ✅ Complete | ❌ Blocked

---

## Backend API Reference

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Health Check
```
GET /health
Response: { "status": "healthy", "service": "scout9" }
```

#### Titles (Games)
```
GET /api/titles
Response: [{ "id": "3", "name": "League of Legends" }, { "id": "6", "name": "VALORANT" }]
```

#### Tournaments
```
GET /api/tournaments?title={lol|valorant}
Response: [{ "id": "758024", "name": "LCK - Spring 2024", "logoUrl": "...", "titleId": "3" }]
```

#### Teams
```
GET /api/teams?tournament={tournamentId}
GET /api/teams/search?q={query}&title={lol|valorant}
GET /api/teams/{teamId}
GET /api/teams/{teamId}/series?limit={1-50}
Response: { "id": "123", "name": "T1", "logoUrl": "..." }
```

#### Reports
```
POST /api/reports
Body: { "teamId": "123", "teamName": "T1", "matchCount": 10, "titleId": "3" }
Response: Full ScoutingReport object

GET /api/reports/{reportId}
GET /api/reports
```

#### Matchups (Head-to-Head)
```
GET /api/matchup?team1={id}&team2={id}&title={lol|valorant}&matches={count}
Response: HeadToHeadAnalysis object
```

---

## TypeScript Types (Frontend)

```typescript
// === Core Types ===
type GameTitle = 'lol' | 'valorant';

interface Title {
  id: string;
  name: string;
}

interface Tournament {
  id: string;
  name: string;
  logoUrl?: string;
  titleId: string;
}

interface Team {
  id: string;
  name: string;
  logoUrl?: string;
  colorHex?: string;
}

// === Report Types ===
interface ScoutingReport {
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

interface DigestibleReport {
  teamName: string;
  matchesAnalyzed: number;
  generatedAt: string;
  executiveSummary: string;
  commonStrategies: CommonStrategiesSection;
  playerTendencies: PlayerTendencyInsight[];
  recentCompositions: CompositionInsight[];
  howToWin: HowToWinSection;
}

// === How to Win (THE DIFFERENTIATOR) ===
interface HowToWinSection {
  winCondition: string;
  actionableInsights: ActionableInsight[];
  draftStrategy: DraftStrategySection;
  inGameStrategy: InGameStrategyInsight[];
  confidenceScore: number;
}

interface ActionableInsight {
  recommendation: string;
  dataBacking: string;
  impact: 'HIGH' | 'MEDIUM' | 'LOW';
  actionType: 'BAN' | 'PICK' | 'TARGET_PLAYER' | 'FORCE_MAP' | 'STRATEGY';
  confidence: number;
}

interface DraftStrategySection {
  priorityBans: DraftInsight[];
  recommendedPicks: DraftInsight[];
  targetPicks: DraftInsight[];
}

interface DraftInsight {
  character: string;
  reason: string;
  playerName?: string;
  winRate?: number;
  pickRate?: number;
}

interface InGameStrategyInsight {
  phase: string;
  timing: string;
  strategy: string;
  priority: 'HIGH' | 'MEDIUM' | 'LOW';
}

// === Team Analysis ===
interface TeamAnalysis {
  teamId: string;
  teamName: string;
  title: GameTitle;
  matchesAnalyzed: number;
  gamesAnalyzed: number;
  winRate: number;
  strengths: Insight[];
  weaknesses: Insight[];
  lolMetrics?: LoLTeamMetrics;
  valMetrics?: VALTeamMetrics;
}

interface LoLTeamMetrics {
  firstBloodRate: number;
  firstDragonRate: number;
  firstTowerRate: number;
  firstTowerAvgTime: number;
  goldDiff15: number;
  dragonControlRate: number;
  heraldControlRate: number;
  baronControlRate: number;
  elderDragonRate: number;
  avgGameDuration: number;
  earlyGameRating: number;
  midGameRating: number;
  lateGameRating: number;
  aggressionScore: number;
  winConditions: string[];
}

interface VALTeamMetrics {
  attackWinRate: number;
  defenseWinRate: number;
  pistolWinRate: number;
  attackPistolWinRate: number;
  defensePistolWinRate: number;
  ecoRoundWinRate: number;
  forceBuyWinRate: number;
  fullBuyWinRate: number;
  firstBloodRate: number;
  mapStats: Record<string, MapStats>;
  mapPool: MapPoolEntry[];
  aggressionScore: number;
  clutchRate: number;
}

// === Player Profile ===
interface PlayerProfile {
  playerId: string;
  nickname: string;
  role: string;
  teamId: string;
  gamesPlayed: number;
  kda: number;
  avgKills: number;
  avgDeaths: number;
  avgAssists: number;
  characterPool: CharacterStats[];
  signaturePicks: string[];
  threatLevel: number;
  threatReason: string;
  weaknesses: Insight[];
  tendencies: string[];
}

interface CharacterStats {
  name: string;
  games: number;
  wins: number;
  losses: number;
  winRate: number;
  avgKDA: number;
}

// === Common Strategies ===
interface CommonStrategiesSection {
  attackPatterns: StrategyInsight[];
  defenseSetups: StrategyInsight[];
  objectivePriorities: StrategyInsight[];
  timingPatterns: StrategyInsight[];
}

interface StrategyInsight {
  text: string;
  metric: string;
  value: number;
  sampleSize: number;
  context: string;
}

// === Compositions ===
interface CompositionAnalysis {
  topCompositions: CompositionInsight[];
  firstPickPriorities: string[];
  commonBans: string[];
}

interface CompositionInsight {
  characters: string[];
  frequency: number;
  winRate: number;
  gamesPlayed: number;
  archetype?: string;
}

// === Head to Head ===
interface HeadToHeadAnalysis {
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
interface Insight {
  text: string;
  importance: 'HIGH' | 'MEDIUM' | 'LOW';
  dataPoints?: string[];
}

interface MapStats {
  mapName: string;
  gamesPlayed: number;
  wins: number;
  winRate: number;
  attackWinRate: number;
  defenseWinRate: number;
}

interface MapPoolEntry {
  map: string;
  games: number;
  winRate: number;
  comfort: number;
}

interface PlayerTendencyInsight {
  playerName: string;
  role: string;
  insights: string[];
  signaturePick?: { character: string; pickRate: number; winRate: number };
  kda: number;
  threatLevel: number;
}
```

---

## Component Structure

```
frontend/
├── app/
│   ├── layout.tsx              # Root layout with providers
│   ├── page.tsx                # Landing page (game selection)
│   ├── globals.css             # Global styles + themes
│   ├── lol/
│   │   ├── page.tsx            # LoL team selection
│   │   └── report/[id]/
│   │       └── page.tsx        # LoL report view
│   ├── valorant/
│   │   ├── page.tsx            # VALORANT team selection
│   │   └── report/[id]/
│   │       └── page.tsx        # VALORANT report view
│   └── compare/
│       └── page.tsx            # Head-to-head comparison
├── components/
│   ├── ui/                     # shadcn/ui components
│   ├── layout/
│   │   ├── Header.tsx
│   │   ├── GameToggle.tsx
│   │   └── Footer.tsx
│   ├── team/
│   │   ├── TeamSearch.tsx
│   │   ├── TeamCard.tsx
│   │   └── TeamSelector.tsx
│   ├── report/
│   │   ├── ReportHeader.tsx
│   │   ├── ExecutiveSummary.tsx
│   │   ├── CommonStrategies.tsx
│   │   ├── PlayerTendencies.tsx
│   │   ├── PlayerCard.tsx
│   │   ├── Compositions.tsx
│   │   ├── HowToWin.tsx        # THE KEY COMPONENT
│   │   ├── WinCondition.tsx
│   │   ├── ActionableInsights.tsx
│   │   ├── DraftStrategy.tsx
│   │   └── InGameStrategy.tsx
│   ├── compare/
│   │   ├── ComparisonView.tsx
│   │   └── StyleRadar.tsx
│   ├── charts/
│   │   ├── WinRateBar.tsx
│   │   ├── PlaystyleRadar.tsx
│   │   ├── MapHeatmap.tsx
│   │   └── ThreatMeter.tsx
│   └── shared/
│       ├── LoadingReport.tsx
│       ├── ConfidenceBadge.tsx
│       ├── ImpactBadge.tsx
│       ├── CharacterIcon.tsx
│       └── StatCard.tsx
├── lib/
│   ├── api.ts                  # API client
│   ├── utils.ts                # Utility functions
│   └── constants.ts            # Game data, colors, icons
├── hooks/
│   ├── useReport.ts
│   ├── useTeams.ts
│   └── useMatchup.ts
└── types/
    └── index.ts                # All TypeScript types
```

---

## Design System

### Color Themes

#### LoL Theme (Gold/Blue)
```css
--lol-primary: #C89B3C;      /* Gold */
--lol-secondary: #0A1428;    /* Dark Blue */
--lol-accent: #0AC8B9;       /* Teal */
--lol-bg: #010A13;           /* Near Black */
--lol-card: #1E2328;         /* Card Background */
--lol-text: #F0E6D2;         /* Light Gold */
--lol-muted: #785A28;        /* Muted Gold */
```

#### VALORANT Theme (Red/Cyan)
```css
--val-primary: #FF4655;      /* Valorant Red */
--val-secondary: #0F1923;    /* Dark */
--val-accent: #00D4AA;       /* Cyan */
--val-bg: #0F1923;           /* Background */
--val-card: #1F2326;         /* Card */
--val-text: #ECE8E1;         /* Light */
--val-muted: #768079;        /* Muted */
```

### Typography
- **Headings**: Inter (bold, tight tracking)
- **Body**: Inter (regular)
- **Code/Stats**: JetBrains Mono

### Key UI Patterns
- Dark theme default (esports aesthetic)
- Glassmorphism cards with subtle borders
- Neon accent glows on important elements
- Smooth animations (Framer Motion)
- Impact badges (HIGH=red, MEDIUM=yellow, LOW=green)

---

## Key Features to Implement

### 1. How to Win Section (Priority #1)
This is the hackathon differentiator. Must include:
- **Win Condition Banner**: Large, prominent, one-liner
- **Confidence Score**: Visual indicator (0-100%)
- **Actionable Insights Grid**: Cards with impact badges
- **Draft Strategy Visualizer**: Ban/pick recommendations with icons
- **In-Game Strategy Timeline**: Phase-by-phase recommendations

### 2. Team Selection Flow
- Game toggle (LoL/VALORANT) with smooth transition
- Searchable team dropdown with logos
- Tournament filter
- Match count slider (5-50)
- Generate button with loading state

### 3. Report Dashboard Sections
1. **Header**: Team logo, name, report metadata
2. **Executive Summary**: Key paragraph with stats highlighted
3. **Common Strategies**: Game-specific grid (4 columns)
4. **Player Tendencies**: Horizontal scroll of player cards
5. **Recent Compositions**: Table with character icons
6. **How to Win**: Full-width hero section

### 4. Data Visualizations
- **Playstyle Radar**: Early/Mid/Late game + Aggression
- **Win Rate Bars**: Horizontal progress bars
- **Map Pool Heatmap**: VALORANT map comfort levels
- **Threat Meter**: Player threat visualization
- **Trend Sparklines**: Recent form indicators

### 5. Head-to-Head Comparison
- Side-by-side team cards
- Style comparison radar overlay
- Historical record display
- Key insights list
- "Generate Report for Team X" buttons

### 6. PDF Export
- Print-friendly CSS
- One-page summary option
- Full report export
- Uses browser print or html2pdf

---

## Development Commands

```bash
# Navigate to frontend
cd /Users/rahulm/Downloads/skysthelimit\ hack/frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Run linting
npm run lint
```

---

## Backend Startup (Reference)

```bash
# From project root
cd /Users/rahulm/Downloads/skysthelimit\ hack

# Set environment variables
export GRID_API_KEY=your_key_here
export PORT=8080

# Run backend
go run cmd/server/main.go
```

---

## Files Created Log

| File | Purpose | Status |
|------|---------|--------|
| `frontend/FRONTEND_SPEC.md` | This specification file | ✅ |
| `frontend/package.json` | Project dependencies | ⬜ |
| `frontend/tailwind.config.ts` | Tailwind configuration | ⬜ |
| `frontend/app/layout.tsx` | Root layout | ⬜ |
| `frontend/app/page.tsx` | Landing page | ⬜ |
| `frontend/app/globals.css` | Global styles | ⬜ |
| `frontend/types/index.ts` | TypeScript types | ⬜ |
| `frontend/lib/api.ts` | API client | ⬜ |
| `frontend/lib/constants.ts` | Game constants | ⬜ |
| ... | ... | ... |

---

## Notes & Decisions

### Architecture Decisions
- **Next.js 14 App Router**: Modern React with server components
- **shadcn/ui**: High-quality, customizable components
- **Tailwind CSS**: Rapid styling with design system
- **Framer Motion**: Professional animations
- **Recharts**: Data visualizations

### Development Notes
- (Add notes here as development proceeds)

---

## Resuming Development

If context is lost, use this checklist:
1. Read this FRONTEND_SPEC.md file
2. Check the Progress Tracker above
3. Look at Files Created Log for current state
4. Continue from the next ⬜ Not Started item
5. Update progress and add notes as you go

---
