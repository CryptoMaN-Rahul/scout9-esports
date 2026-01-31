# Scout9 - Automated Scouting Report Generator

## Cloud9 x JetBrains Hackathon - Category 2 Submission

**Scout9** is an automated scouting report generator for esports teams, built using official GRID API data for League of Legends and VALORANT. It analyzes opponent match data and generates actionable, data-backed scouting reports to help coaches and players prepare for upcoming matches.

### ✅ Validation Status (January 1, 2026)

| Feature | Status | Details |
|---------|--------|---------|
| **LoL Counter-Strategy** | ✅ 5/5 | Class matchup insights, draft recommendations, win conditions |
| **VALORANT Counter-Strategy** | ✅ 4/4 | Economy analysis, site-specific analysis with JSONL events |
| **Head-to-Head Comparison** | ✅ 4/4 | Historical record, style comparison with team analyses |
| **Economy Analyzer** | ✅ 3/3 | Per-map economy stats, eco/force/full buy analysis |
| **GRID API Integration** | ✅ | Native `teamId` filter - no hardcoded tournament IDs |

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Data Flow](#data-flow)
4. [Features & Capabilities](#features--capabilities)
5. [API Endpoints](#api-endpoints)
6. [Report Output Format](#report-output-format)
7. [GRID API Integration](#grid-api-integration)
8. [Intelligence Modules](#intelligence-modules)
9. [How to Run](#how-to-run)
10. [Technology Stack](#technology-stack)

---

## Project Overview

Scout9 solves a real problem for esports teams: **manual scouting is time-consuming and often misses patterns that data can reveal**. By automating the analysis of opponent match data, Scout9 provides:

- **Team-wide strategy analysis** - Macro patterns, objective priorities, timing tendencies
- **Individual player profiling** - Champion/agent pools, KDA, playstyle identification
- **Composition analysis** - Draft patterns, synergies, win rates by comp
- **Counter-strategy recommendations** - The "How to Win" section with specific, actionable insights

### Target Users
- Esports coaches and analysts
- Professional players preparing for matches
- Team strategists and data analysts

---

## Architecture

```
scout9/
├── cmd/server/main.go          # HTTP server entry point
├── internal/api/router.go      # REST API routes
├── pkg/
│   ├── grid/                   # GRID API integration
│   │   ├── client.go           # Main API client with rate limiting
│   │   ├── central_data.go     # Teams, tournaments, series queries
│   │   ├── series_state.go     # Match details, player stats
│   │   ├── file_download.go    # JSONL event file parsing
│   │   └── types.go            # Data structures
│   ├── intelligence/           # Analysis engines
│   │   ├── lol_analyzer.go     # LoL team/player analysis
│   │   ├── val_analyzer.go     # VALORANT team/player analysis
│   │   ├── composition_analyzer.go
│   │   ├── counter_strategy.go # "How to Win" generation (ENHANCED)
│   │   ├── economy_analyzer.go # VALORANT economy analysis (NEW)
│   │   ├── head_to_head_analyzer.go # Team comparison (NEW)
│   │   ├── event_analyzer.go   # JSONL event parsing
│   │   ├── timing_analyzer.go  # Objective timing patterns
│   │   ├── matchup_analyzer.go # Player matchup analysis (ENHANCED)
│   │   ├── site_analyzer.go    # VALORANT site patterns
│   │   ├── trend_analyzer.go   # Performance trends
│   │   ├── role_detector.go    # Role inference from data
│   │   ├── lane_detector.go    # LoL lane classification
│   │   ├── stats_engine.go     # Statistical utilities
│   │   ├── champion_data.go    # LoL champion database
│   │   ├── validator.go        # Data validation
│   │   └── types.go            # Analysis type definitions
│   ├── report/                 # Report generation
│   │   ├── generator.go        # Orchestrates report creation
│   │   ├── formatter.go        # Converts to hackathon format
│   │   └── html_export.go      # HTML export (placeholder)
│   ├── llm/                    # AI-powered insights
│   │   ├── openai.go           # OpenAI integration
│   │   ├── service.go          # LLM service interface
│   │   └── template.go         # Template-based fallback
│   └── cache/
│       └── redis.go            # Redis caching layer
├── scripts/                    # Test and validation scripts
└── docker-compose.yml          # Infrastructure setup
```

---

## Data Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           GRID API                                       │
├─────────────────┬─────────────────────┬─────────────────────────────────┤
│  Central Data   │   Series State      │   File Download                 │
│  (GraphQL)      │   (GraphQL)         │   (REST)                        │
│                 │                     │                                 │
│  • Teams        │  • Match results    │  • JSONL event files            │
│  • Tournaments  │  • Player stats     │  • Kill positions               │
│  • Series IDs   │  • Draft actions    │  • Objective timings            │
│                 │  • Objectives       │  • Detailed events              │
└────────┬────────┴──────────┬──────────┴────────────────┬────────────────┘
         │                   │                           │
         ▼                   ▼                           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        pkg/grid/client.go                               │
│                    (Rate-limited API client)                            │
└─────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     pkg/intelligence/                                    │
├─────────────────┬─────────────────────┬─────────────────────────────────┤
│  LoL Analyzer   │   VAL Analyzer      │   Composition Analyzer          │
│  • Team metrics │   • Team metrics    │   • Draft patterns              │
│  • Player stats │   • Player stats    │   • Synergies                   │
│  • Objectives   │   • Economy         │   • Win rates                   │
├─────────────────┼─────────────────────┼─────────────────────────────────┤
│  Event Analyzer │   Timing Analyzer   │   Counter Strategy Engine       │
│  • Kill events  │   • Jungle pathing  │   • Weakness exploitation       │
│  • Objectives   │   • Objective times │   • Draft recommendations       │
│  • Positions    │   • Game pace       │   • In-game strategies          │
├─────────────────┼─────────────────────┼─────────────────────────────────┤
│  Economy        │   Head-to-Head      │   Matchup Analyzer              │
│  Analyzer (NEW) │   Analyzer (NEW)    │   (ENHANCED)                    │
│  • Eco/force/   │   • Historical      │   • Class performance           │
│    full buy WR  │     record          │   • Counter-class picks         │
│  • Per-map      │   • Style compare   │   • Specific champion           │
│    economy      │   • Insights        │     recommendations             │
└─────────────────┴─────────────────────┴─────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     pkg/report/generator.go                             │
│                   (Orchestrates all analyzers)                          │
└─────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     pkg/report/formatter.go                             │
│              (Converts to hackathon-compliant format)                   │
└─────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      DigestibleReport                                    │
│  • Executive Summary                                                     │
│  • Common Strategies (attack/defense/objectives/timing)                 │
│  • Player Tendencies                                                     │
│  • Recent Compositions                                                   │
│  • HOW TO WIN (actionable insights, draft strategy, in-game strategy)   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Features & Capabilities

### League of Legends Analysis

| Feature | Data Source | Description |
|---------|-------------|-------------|
| **Win Rate & Form** | Series State | Overall win rate, recent form (last 5 games) |
| **First Blood Rate** | Series State | Early game aggression indicator |
| **Dragon Control** | Series State + JSONL | First dragon rate, dragon soul patterns |
| **Baron Control** | Series State | Baron secure rate, timing patterns |
| **Herald Priority** | Series State | Herald control rate |
| **Tower Priority** | Series State | First tower rate, average timing |
| **Game Duration** | Series State | Average game length, pace indicator |
| **Draft Patterns** | Series State | First pick priorities, common bans |
| **Champion Pools** | Series State | Per-player champion stats, win rates |
| **Multikill Stats** | Series State | Double/triple/quadra/penta kills |
| **Synergy Analysis** | Series State | Assist network, duo identification |
| **Objective Focus** | Series State | Split-pusher vs teamfighter detection |
| **Item Builds** | Series State | Core item patterns |
| **Role Detection** | Champion Data | Inferred from champion pool |

### VALORANT Analysis

| Feature | Data Source | Description |
|---------|-------------|-------------|
| **Win Rate & Form** | Series State | Overall win rate, recent form |
| **Attack Win Rate** | Series State | Attack side performance |
| **Defense Win Rate** | Series State | Defense side performance |
| **Pistol Win Rate** | Series State | Pistol round performance |
| **First Blood Rate** | Series State | Opening duel success |
| **Economy Analysis** | Series State | Eco/force/full buy win rates |
| **Map Pool** | Series State | Per-map win rates, comfort picks |
| **Agent Pools** | Series State | Per-player agent stats |
| **Weapon Stats** | Series State | Kills by weapon, Operator dependency |
| **Synergy Analysis** | Series State | Assist network, duo identification |
| **Site Patterns** | JSONL Events | Plant positions, site preferences |
| **Round Win Types** | JSONL Events | Elimination vs defuse vs explosion |

### Counter-Strategy Generation ("How to Win")

The **key differentiator** - specific, actionable, data-backed recommendations:

1. **Win Condition** - One-sentence summary of how to beat this opponent
2. **Exploitable Weaknesses** - Data-backed vulnerabilities with evidence
3. **Draft Recommendations**:
   - Priority bans (high win-rate champions/agents)
   - Target picks (force opponent onto weak picks)
   - Recommended picks (counter their style)
4. **In-Game Strategies** - Timing-specific recommendations
5. **Target Players** - Who to focus, who to avoid

### Enhanced Counter-Strategy Features (NEW)

Scout9 now includes advanced analysis capabilities for more specific, hackathon-winning insights:

#### Champion Class Matchup Analysis (LoL)
- Analyzes player performance ON different champion classes (assassin, mage, tank, fighter, marksman, support)
- Generates insights like: "Faker has 2.1 KDA on control mages vs 8.5 on assassins"
- Provides specific champion recommendations: "Pick Zed, LeBlanc, Akali vs [player]"
- Thresholds: KDA < 2.5 with ≥2 games → pick recommendation; KDA > 4.0 with WR > 60% and ≥3 games → ban recommendation

#### Economy Exploitation (VALORANT)
- Tracks eco/force/full buy round performance
- Generates insights like: "They only win 15% of eco rounds - play aggressive on their saves"
- Thresholds:
  - Eco WR < 15% → "aggressive on saves" insight
  - Force WR > 45% → "respect force buys" insight
  - Force WR < 25% → "punish force buys" insight
- Per-map economy breakdowns

#### Site-Specific Attack Patterns (VALORANT)
- Identifies sites with <40% win rate and ≥3 attempts
- Generates attack recommendations: "Attack B-Site on Ascent - they only hold it 35% of the time"
- Tracks pistol round patterns by site
- Detects predictable attack patterns (>40% frequency)

#### Head-to-Head Comparison
- New `/api/matchups` endpoint comparing two teams
- Historical record (wins/losses between teams)
- Style comparison (early/mid/late game ratings)
- Aggression comparison
- Generates insights even without historical matches

#### Enhanced Win Conditions
- Includes specific data points (percentages, player names)
- Includes jungle pathing recommendations when available
- Includes site recommendations when available
- Limited to 2-3 sentences for clarity

#### Confidence Scoring
- Sample size-based confidence boosts:
  - ≥15 games: +25 points
  - ≥10 games: +15 points
- Per-insight confidence levels
- Low sample size warnings (<5 games analyzed)

---

## API Endpoints

### Teams & Tournaments

```
GET /api/teams?title={lol|valorant}
    Returns list of teams for the specified game

GET /api/teams/{teamId}
    Returns detailed team information

GET /api/tournaments?title={lol|valorant}
    Returns available tournaments
```

### Reports

```
POST /api/reports/generate
    Body: { "teamId": "123", "titleId": "3", "matchCount": 10 }
    Generates full scouting report

GET /api/reports/{reportId}
    Retrieves a generated report

POST /api/reports/text
    Body: { "teamId": "123", "titleId": "3", "matchCount": 10 }
    Returns plain text report in hackathon format
```

### Matchups

```
GET /api/matchup?team1={teamId}&team2={teamId}&title={lol|valorant}&matches={count}
    Returns head-to-head analysis using HeadToHeadAnalyzer
    Response: HeadToHeadReport with:
      - Historical record (wins/losses)
      - Style comparison (early/mid/late game ratings)
      - Key insights
      - Confidence score
      - Warnings (if low sample size)
    
    Error Handling:
      - 400: Missing team1 or team2 parameters
      - 404: Team not found
      - 502: GRID API errors
```

### Health

```
GET /health
    Returns server health status
```

---

## Report Output Format

### DigestibleReport Structure

```go
type DigestibleReport struct {
    TeamName           string
    MatchesAnalyzed    int
    GeneratedAt        string
    ExecutiveSummary   string                    // 1-paragraph summary
    CommonStrategies   CommonStrategiesSection   // Attack/defense/objectives/timing
    PlayerTendencies   []PlayerTendencyInsight   // Per-player insights
    RecentCompositions []CompositionInsight      // Top compositions
    HowToWin           HowToWinSection           // THE KEY SECTION
}
```

### Example Text Output

```
═══════════════════════════════════════════════════════════════
  SCOUTING REPORT: Team Liquid
  Generated: 2026-01-01 12:00:00 | Matches Analyzed: 10
═══════════════════════════════════════════════════════════════

📋 EXECUTIVE SUMMARY
───────────────────────────────────────────────────────────────
Team Liquid is an aggressive early-game team with 65% win rate (10 matches).
Key strength: strong dragon control (72% first dragon). Exploitable weakness:
passive late game (only 45% win rate in games over 35 minutes). Current form: improving.

🎯 COMMON STRATEGIES
───────────────────────────────────────────────────────────────
Objective Priorities:
  • Prioritizes first Drake (72% contest rate)
  • Strong Herald priority (68% control rate)
Timing Patterns:
  • Average first tower at ~12 mins (65% first tower rate)
  • Average game duration: 31 mins (standard pace)

👤 PLAYER TENDENCIES
───────────────────────────────────────────────────────────────
  • Player 'CoreJJ' (Support) signature pick: Nautilus (80% pick rate, 75% win rate)
  • Player 'CoreJJ' has exceptional synergy with Yeon (52% of assists)
  • Player 'Yeon' averages 1.8 multikills/game (strong teamfight presence)

🎮 RECENT COMPOSITIONS
───────────────────────────────────────────────────────────────
  • Comp #1 (40% frequency, 75% win rate): K'Sante, Xin Zhao, Ahri, Jinx, Nautilus

🏆 HOW TO WIN
═══════════════════════════════════════════════════════════════
WIN CONDITION: Survive early game and outscale - they struggle after 35 minutes
Confidence: 78%

Actionable Insights:
  [HIGH] Exploit: Weak late game execution
         Data: Only 45% win rate in games over 35 minutes
  [HIGH] Play aggressive early - invade and force fights
         Data: Opponent has only 38% first blood rate

Draft - Priority Bans:
  🚫 Ban Nautilus from CoreJJ - 75% win rate

In-Game Strategy:
  ⚔️  Draft scaling compositions and survive early (Draft phase and 0-15 minutes)
      Reason: Opponent wins fast (avg 31 min) - they may struggle in extended games
```

---

## GRID API Integration

### Endpoints Used

| API | URL | Rate Limit | Purpose |
|-----|-----|------------|---------|
| Central Data | `https://api-op.grid.gg/central-data/graphql` | 40/min | Teams, tournaments, series |
| Series State | `https://api-op.grid.gg/live-data-feed/series-state/graphql` | 1200/min | Match details, stats |
| File Download | `https://api.grid.gg/file-download/...` | 20/min | JSONL event files |

### Key API Features Used

**SeriesFilter (Central Data API):**
- `teamId: ID` - Direct team filtering (no hardcoded tournament IDs needed)
- `teamIds: IdFilter` - Multiple teams with `in` operator
- `tournament: SeriesTournamentFilter` - Tournament-based filtering
- `titleId: ID` - Filter by game title (LoL/VALORANT)
- `startTimeScheduled: DateTimeFilter` - Date range filtering

**Validated via API Introspection (January 2026):**
```graphql
# Query series for a specific team
query SeriesForTeam($teamId: ID!, $limit: Int!) {
  allSeries(
    filter: { teamId: $teamId }
    orderBy: StartTimeScheduled
    orderDirection: DESC
    first: $limit
  ) {
    edges {
      node {
        id
        startTimeScheduled
        teams { baseInfo { id name logoUrl } }
        tournament { id name }
        title { id name }
      }
    }
  }
}
```

### Data Fields Captured

**From Series State API:**
- Player: kills, deaths, assists, netWorth, character, items, objectives, multikills, weaponKills, abilities
- Team: score, kills, deaths, objectives, structuresDestroyed
- Game: map, draftActions, segments (VALORANT rounds)

**From JSONL Events:**
- Kill events with positions
- Dragon/Baron/Herald kills with timing
- Tower destructions
- Spike plants/defuses with positions
- Round win types

### Available Tournaments

**League of Legends:**
- LCK, LCS, LEC, LPL, LTA (2024-2025 seasons)

**VALORANT:**
- VCT Americas (2024-2025 seasons)
- Masters Madrid

---

## Intelligence Modules

### LoL Analyzer (`pkg/intelligence/lol_analyzer.go`)

Analyzes League of Legends team and player data:
- Calculates team metrics (win rate, first blood, dragon control, etc.)
- Builds player profiles with champion pools, KDA, multikills
- Tracks assist networks for synergy analysis
- Detects roles from champion pools

### VAL Analyzer (`pkg/intelligence/val_analyzer.go`)

Analyzes VALORANT team and player data:
- Calculates team metrics (attack/defense win rates, pistol rounds)
- Tracks economy round performance (eco/force/full buy)
- Builds player profiles with agent pools, weapon stats
- Analyzes map pool strength

### Composition Analyzer (`pkg/intelligence/composition_analyzer.go`)

Analyzes draft patterns and team compositions:
- Identifies top compositions by frequency and win rate
- Tracks first pick priorities
- Identifies common bans
- Detects composition archetypes

### Counter Strategy Engine (`pkg/intelligence/counter_strategy.go`)

Generates the "How to Win" section:
- Identifies exploitable weaknesses
- Generates draft recommendations (bans, target picks)
- Creates in-game strategy recommendations
- Identifies target players
- Analyzes synergy to find duos to break up
- **NEW: Enhanced win condition generation** with specific data points
- **NEW: Confidence scoring** with sample size boosts
- **NEW: Sample size warnings** for low-data reports
- **NEW: Integration with EconomyAnalyzer** for VALORANT insights
- **NEW: Integration with class matchup analysis** for LoL insights

### Economy Analyzer (`pkg/intelligence/economy_analyzer.go`) - NEW

Analyzes VALORANT economy round performance:
- Classifies rounds as eco/force/full buy based on round number patterns
- Calculates win rates for each economy state
- Tracks per-map economy statistics
- Generates hackathon-format insights:
  - "They only win 15% of eco rounds - play aggressive on their saves"
  - "They win 48% of force buys - respect their force rounds"
  - "They only win 22% of force buys - punish their force rounds"
- Identifies site weaknesses (<40% win rate with ≥3 attempts)

### Head-to-Head Analyzer (`pkg/intelligence/head_to_head_analyzer.go`) - NEW

Compares two teams for matchup analysis:
- Queries GRID API for historical matches between teams
- Calculates head-to-head record (wins/losses)
- Compares team styles:
  - Early/mid/late game ratings
  - Aggression scores
  - Attack vs defense focus (VALORANT)
- Generates style insights even without historical matches
- Calculates confidence score based on data availability

### Matchup Analyzer (`pkg/intelligence/matchup_analyzer.go`) - ENHANCED

Analyzes player matchups and class performance:
- **NEW: GenerateClassMatchupInsights()** - Compares best vs worst class KDA
- **NEW: GenerateSpecificDraftRecommendations()** - Specific champion recommendations
- **NEW: CounterClassMap** - Maps weak classes to counter picks
- Generates insights like: "Faker has 2.1 KDA on mages vs 8.5 on assassins"

### Event Analyzer (`pkg/intelligence/event_analyzer.go`)

Parses JSONL event files for detailed analysis:
- Extracts kill events with positions
- Parses objective events (dragons, barons, towers)
- Extracts VALORANT plant/defuse events
- Infers site from plant positions

### Timing Analyzer (`pkg/intelligence/timing_analyzer.go`)

Analyzes timing patterns:
- Jungle pathing from kill positions
- Objective timing patterns
- Game pace analysis

---

## How to Run

### Prerequisites

- Go 1.25+
- Redis (for caching)
- PostgreSQL (optional, for persistence)
- GRID API key

### Environment Variables

```bash
# Copy example env file
cp .env.example .env

# Required variables:
GRID_API_KEY=your_api_key_here
REDIS_URL=localhost:6379
PORT=8080
OPENAI_API_KEY=your_openai_key  # Optional, for AI summaries
```

### Running with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### Running Locally

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go

# Server starts on http://localhost:8080
```

### Testing

```bash
# Run all tests
go test ./...

# Test GRID API connection
GRID_API_KEY=xxx go run scripts/test_grid_data_capture.go

# Validate event parsing
GRID_API_KEY=xxx go run scripts/validate_grid_data.go

# Validate enhanced counter-strategy feature (NEW)
GRID_API_KEY=xxx go run scripts/validate_counter_strategy.go
```

### Validation Results (January 1, 2026)

The enhanced counter-strategy feature has been validated against real GRID API data:

**LoL Analysis (Hanwha Life Esports - LCK):**
- ✅ Class matchup insights generated (tank vs fighter KDA comparison)
- ✅ Draft recommendations with specific champion bans
- ✅ Win conditions with data points (first blood %, dragon control %)
- ✅ Confidence score calculated (100% with 10 games)
- ✅ Target player identification with synergy analysis

**VALORANT Analysis (2GAME eSports - VCT Americas):**
- ✅ Economy analysis (eco/force/full buy win rates)
- ✅ Per-map economy breakdowns
- ✅ Economy insights generated (force buy warnings, eco exploitation)
- ✅ Site-specific analysis with JSONL event data (attack/defense win rates by site)

**Head-to-Head (DRX vs DN FREECS):**
- ✅ Historical record retrieved (4 matches found)
- ✅ Style comparison included (early/mid/late game ratings)
- ✅ Insights generated with confidence scores

**Economy Analyzer (Evil Geniuses):**
- ✅ Eco/force/full buy win rates calculated
- ✅ Per-map statistics (10 maps analyzed)
- ✅ Economy insights generated (6 insights)

---

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.25 | Backend server |
| Router | Chi | HTTP routing |
| Cache | Redis | API response caching |
| Database | PostgreSQL | Report persistence (optional) |
| AI | OpenAI GPT-4 | Executive summaries (optional) |
| API | GRID | Esports data source |

### Key Dependencies

```go
github.com/go-chi/chi/v5      // HTTP router
github.com/redis/go-redis/v9  // Redis client
github.com/google/uuid        // UUID generation
github.com/sashabaranov/go-openai // OpenAI client
```

---

## What Makes Scout9 Different

1. **Data-Backed Insights** - Every recommendation includes the underlying data
2. **Actionable Output** - Not just stats, but specific strategies to exploit
3. **"How to Win" Focus** - The report is structured around winning, not just information
4. **Hackathon-Compliant Format** - Output matches exactly what judges expect
5. **Comprehensive Analysis** - Uses all available GRID API data (Series State + JSONL events)
6. **Synergy Detection** - Identifies player duos and assist networks
7. **Economy Analysis** - VALORANT eco/force/full buy tracking with per-map breakdowns
8. **Role Inference** - Detects roles from champion/agent pools when not provided
9. **Champion Class Matchup Analysis** - Identifies player weaknesses on specific champion classes
10. **Site-Specific Insights** - VALORANT site weakness identification and attack recommendations
11. **Head-to-Head Comparison** - Team style comparison even without historical matches
12. **Confidence Scoring** - Sample size-based confidence with warnings for low-data reports
13. **Native GRID API Integration** - Uses proper `teamId` filter, no hardcoded tournament IDs
14. **Validated with Real Data** - All features tested against live GRID API data (January 2026)

---

## Future Enhancements

1. **Position Heatmaps** - Visual representation of kill positions
2. **Live Match Integration** - Real-time analysis during matches
3. **Historical Trends** - Long-term performance tracking
4. **PDF Export** - Professional report formatting
5. **Slack/Discord Integration** - Report delivery to team channels
6. **Web Dashboard** - Interactive UI for report viewing and comparison
7. **Custom Thresholds** - User-configurable insight thresholds
8. **Multi-Language Support** - Reports in different languages
9. **Enhanced Site Analysis** - More JSONL event data for site-specific insights
10. **Style Comparison** - Deeper team style analysis for head-to-head matchups

---

## Validation Scripts

| Script | Purpose |
|--------|---------|
| `scripts/validate_counter_strategy.go` | Validates enhanced counter-strategy with real GRID data |
| `scripts/test_grid_data_capture.go` | Tests GRID API connection and data capture |
| `scripts/validate_grid_data.go` | Validates event parsing from JSONL files |
| `scripts/test_enhanced_analyzers.go` | Tests enhanced analyzer modules |
| `scripts/verify_full_integration.go` | Full integration verification |

---

*Built for the Cloud9 x JetBrains Hackathon - Category 2: Automated Scouting Report Generator*

*Last validated: January 1, 2026*
