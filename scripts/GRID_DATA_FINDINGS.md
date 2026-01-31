# GRID Data Structure Findings - Real Data Validation

## ✅ VALIDATION COMPLETE - Parser Fixed and Working

**Status**: Event parsers have been completely rewritten to match actual GRID data structure.

## Key Changes Made

### 1. New Event Structure Types (pkg/grid/types.go)
- Added `EventWrapper` - represents each line in JSONL file
- Added `GridEvent` - represents individual events with actor/action/target
- Added `EventEntity` - represents actor or target with state data
- Added helper methods for extracting position, game state, etc.

### 2. Fixed Event Parser (pkg/grid/file_download.go)
- Completely rewrote `ParseLoLEvents()` and `ParseVALEvents()`
- Now correctly parses nested event structure
- Extracts position data from `actor.state.game.position`
- Calculates `GameTime` from event timestamps
- Infers VALORANT site from plant position coordinates

### 3. Enhanced Event Types
- `KillEvent` now includes `KillerPosition`, `VictimPosition`, `GameTime`
- `DragonKillEvent` now includes `Position`, `GameTime`, dragon type from target.id
- `PlantEvent` now includes `Position`, `Site` (inferred), `Agent`, `GameTime`
- `VALKillEvent` includes full context: positions, agents, round number, map

## Actual GRID Event Structure

### Event Wrapper (each line in JSONL)
```json
{
  "id": "uuid",
  "correlationId": "uuid",
  "occurredAt": "2024-06-08T15:20:02.894Z",
  "seriesId": "2692648",
  "sequenceNumber": 123,
  "events": [
    {
      "id": "uuid",
      "action": "killed",
      "actor": { ... },
      "target": { ... },
      "seriesState": { ... }  // optional, full state snapshot
    }
  ]
}
```

### Event Type Construction
The event type is: `{actor.type}-{action}-{target.type}`

Examples:
- `player-killed-player` → Champion/player kill
- `player-killed-ATierNPC` → Dragon/Baron/Herald kill
- `team-destroyed-tower` → Tower destruction
- `team-picked-character` → Draft pick
- `player-completed-plantBomb` → VALORANT spike plant

---

## League of Legends Data Available

### Event Types Found (Series 2692648 - LEC Summer 2024)
| Event Type | Count | Description |
|------------|-------|-------------|
| player-purchased-item | 324 | Item purchases |
| player-used-ability | 210 | Ability usage |
| player-completed-increaseLevel | 150 | Level ups |
| player-lost-item | 139 | Item loss |
| player-killed-player | 24 | Champion kills |
| player-killed-ATierNPC | 13 | Dragon/Baron/Herald |
| player-killed-BTierNPC | 12 | Scuttle crabs |
| player-destroyed-tower | 9 | Tower destruction |
| team-picked-character | 10 | Draft picks |
| team-banned-character | 10 | Draft bans |
| player-completed-slayInfernalDrake | 4 | Specific dragon type |
| player-completed-slayMountainDrake | 1 | Specific dragon type |
| player-completed-slayOceanDrake | 1 | Specific dragon type |
| player-completed-slayBaron | 1 | Baron kill |
| player-completed-slayRiftHerald | 1 | Herald kill |
| player-completed-slayVoidGrub | 6 | Void grub kills |

### Position Data - AVAILABLE! ✅
```json
"actor": {
  "state": {
    "game": {
      "position": {
        "x": 8606,
        "y": 6604
      }
    }
  }
}
```

### Kill Event Structure
```json
{
  "action": "killed",
  "actor": {
    "id": "21542",
    "type": "player",
    "state": {
      "name": "Zwyroo",
      "side": "blue",
      "teamId": "47435",
      "game": {
        "kills": 1,
        "position": { "x": 8606, "y": 6604 },
        "netWorth": 1861,
        "killAssistsReceived": 1,
        "killAssistsReceivedFromPlayer": [
          { "playerId": "22117", "killAssistsReceived": 1 }
        ]
      }
    }
  },
  "target": {
    "id": "21155",
    "type": "player",
    "state": {
      "name": "Jackies",
      "side": "red",
      "teamId": "53168",
      "game": {
        "deaths": 1,
        "position": { "x": 9146, "y": 6630 },
        "respawnClock": { "currentSeconds": 7 }
      }
    }
  }
}
```

### Dragon Kill Structure
```json
{
  "action": "killed",
  "actor": {
    "id": "22716",
    "type": "player",
    "state": { "name": "Juhan", "teamId": "53168" }
  },
  "target": {
    "id": "mountainDrake",  // Dragon type in ID!
    "type": "ATierNPC"
  }
}
```

### Tower Destruction Structure
```json
{
  "action": "destroyed",
  "actor": { "type": "team", "id": "47435" },
  "target": {
    "id": "red-turret-mid-2",  // Lane and tower number in ID!
    "type": "tower",
    "state": { "destroyed": true }
  }
}
```

### Draft Structure
```json
{
  "action": "banned",
  "actor": { "type": "team", "id": "47435" },
  "target": {
    "type": "character",
    "state": { "name": "Lucian" }
  }
}
```

---

## VALORANT Data Available

### Event Types Found (Series 2629390 - VCT Americas)
| Event Type | Count | Description |
|------------|-------|-------------|
| player-used-ability | 1476 | Ability usage |
| player-killed-player | 290 | Player kills |
| game-started-round | 45 | Round starts |
| team-won-round | 44 | Round wins |
| game-ended-round | 44 | Round ends |
| player-completed-plantBomb | 26 | Spike plants |
| player-completed-defuseBomb | 9 | Spike defuses |
| player-completed-explodeBomb | 5 | Spike explosions |
| player-selfkilled-player | 7 | Self-kills |

### Position Data - AVAILABLE! ✅
```json
"actor": {
  "state": {
    "game": {
      "position": {
        "x": 5122.83545,
        "y": 6571.28711
      }
    }
  }
}
```

### Map Data - AVAILABLE! ✅
```json
"actor": {
  "state": {
    "map": {
      "name": "breeze",
      "bounds": {
        "min": { "x": -1650, "y": -6200 },
        "max": { "x": 11450, "y": 7200 }
      }
    }
  }
}
```

### Round Win Structure
```json
{
  "action": "won",
  "actor": { "type": "team", "id": "99" },
  "target": {
    "type": "round",
    "id": "round-1",
    "state": {
      "teams": [
        {
          "id": "99",
          "won": true,
          "winType": "opponentEliminated"  // Win reason!
        }
      ]
    }
  }
}
```

### Spike Plant Structure
```json
{
  "action": "completed",
  "actor": {
    "type": "player",
    "state": {
      "name": "mwzera",
      "character": { "name": "jett" },
      "game": {
        "position": { "x": 7365.79, "y": -4548.02 }  // Plant position!
      }
    }
  },
  "target": {
    "type": "plantBomb"
  }
}
```

### Site Detection
**IMPORTANT**: Site is NOT explicitly provided. Must be inferred from:
1. Plant position coordinates
2. Map bounds
3. Known site locations per map

---

## What We Can Now Build

### LoL Analysis (with real data)
1. **Jungle Pathing** ✅ - Position data available on kills
   - Track jungler kill locations over time
   - Determine early game pathing patterns (bot vs top side)
   
2. **Dragon Priority** ✅ - Dragon type in target.id
   - Track which dragons teams prioritize
   - Calculate dragon soul patterns
   
3. **First Blood Patterns** ✅ - Kill events with timing
   - Analyze first blood rate
   - Identify common first blood locations
   
4. **Champion Performance** ✅ - Full draft and game data
   - Win rates by champion
   - KDA by champion
   - Performance vs specific opponents

### VALORANT Analysis (with real data)
1. **Attack Patterns** ✅ - Plant positions available
   - Cluster plant positions to determine site preference
   - Calculate "fast hit" vs "default" patterns
   
2. **Round Win Types** ✅ - winType field available
   - Track elimination vs defuse vs explosion wins
   - Analyze clutch situations
   
3. **Agent Performance** ✅ - Character data available
   - Win rates by agent
   - Kill patterns by agent
   
4. **Map-Specific Analysis** ✅ - Map name available
   - Site preferences per map
   - Win rates per map

---

## Required Parser Fixes

### Current (BROKEN)
```go
type GameEvent struct {
    Type      string                 `json:"type"`  // WRONG!
    Timestamp time.Time              `json:"timestamp"`
    GameTime  int                    `json:"gameTime"`
    Data      map[string]interface{} `json:"data"`
}
```

### Fixed Structure
```go
type EventWrapper struct {
    ID             string      `json:"id"`
    CorrelationID  string      `json:"correlationId"`
    OccurredAt     time.Time   `json:"occurredAt"`
    SeriesID       string      `json:"seriesId"`
    SequenceNumber int         `json:"sequenceNumber"`
    Events         []GridEvent `json:"events"`
}

type GridEvent struct {
    ID                string                 `json:"id"`
    Action            string                 `json:"action"`
    Actor             *EventEntity           `json:"actor"`
    Target            *EventEntity           `json:"target"`
    SeriesState       map[string]interface{} `json:"seriesState,omitempty"`
    IncludesFullState bool                   `json:"includesFullState"`
}

type EventEntity struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    State     map[string]interface{} `json:"state"`
    StateDelta map[string]interface{} `json:"stateDelta"`
}
```

---

## Priority Actions

### ✅ COMPLETED
1. **Fixed event parser** - Now uses correct nested structure
2. **Added position data extraction** - Can track jungle pathing
3. **Added site inference** - VALORANT site detection from coordinates
4. **Added GameTime calculation** - Timing analysis now possible
5. **Backward compatibility** - Existing analyzers work without changes

### NEXT STEPS
1. **Enhance jungle pathing analysis** - Use position data to track actual paths
2. **Improve site inference** - Refine coordinate boundaries per map
3. **Add champion class matchup analysis** - Track performance vs opponent classes
4. **Build "How to Win" engine** - Generate specific, data-backed recommendations
5. **Create HTML report output** - Digestible format for coaches











curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"GamePlayerState\") { name fields { name type { name kind ofType { name } } } } }"
  }' 2>/dev/null | python3 -m json.tool 2>/dev/null || echo "Raw output:" && curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"GamePlayerState\") { name fields { name type { name kind ofType { name } } } } }"
  }' 








  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"PlayerInventory\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool



curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"GameTeamState\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool


  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"GameState\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool



  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"Objective\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool


  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "query SeriesState { seriesState(id: \"2692648\") { games { id map { name } startedAt clock { currentSeconds } teams { id name objectives { id type completionCount } players { id name killAssistsReceived structuresDestroyed objectives { id type completionCount } inventory { items { id name } } } } } } }"
  }' | python3 -m json.tool 2>/dev/null | head -300


curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"DraftAction\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool


  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "{ __type(name: \"Drafter\") { name fields { name type { name kind ofType { name } } } } }"
  }' | python3 -m json.tool


  curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{
    "query": "query SeriesState { seriesState(id: \"2692648\") { games { draftActions { id type sequenceNumber drafter { id } draftable { id name } } } } }"
  }' | python3 -m json.tool 2>/dev/null | head -100

  GRID API Schema Discovery Summary
Series State API - Available Fields We're NOT Using:
GamePlayerState (Player-level):

killAssistsReceived - assists count ✅ NOW USING
killAssistsGiven - assists given ✅ NOW USING
structuresDestroyed - towers destroyed by player ✅ NOW USING
objectives - array of objectives completed (dragons, barons, etc.) ✅ NOW USING
inventory { items { id name } } - player items ✅ NOW USING
weaponKills - kills by weapon (VALORANT) ✅ NOW USING
multikills - multi-kill stats ✅ NOW USING
money - current money ✅ NOW USING
loadoutValue - loadout value ✅ NOW USING

GameTeamState (Team-level):

objectives - team objectives (dragons, barons, towers, etc.) ✅ NOW USING
structuresDestroyed - total structures destroyed ✅ NOW USING
money - team money ✅ NOW USING

GameState (Game-level):

draftActions { type sequenceNumber drafter { id } draftable { id name } } - full draft data ✅ NOW USING
startedAt - game start time ✅ VERSION RESTRICTED
map { name } - map name ✅ ALREADY HAVE
segments - round-by-round data (VALORANT) ✅ NOW USING

## NEW: Segment (Round) Data for VALORANT

The Series State API provides round-by-round data through the `segments` field:

```graphql
segments {
  id
  sequenceNumber
  type
  finished
  teams {
    id
    name
    side        # "attacker" or "defender"
    won
    kills
    deaths
    objectives {
      id
      type
      completionCount
    }
    players {
      id
      name
      kills
      deaths
      objectives {
        id
        type
        completionCount
      }
    }
  }
}
```

This allows us to:
- Track attack/defense win rates per round
- Identify pistol round performance (rounds 1 and 13)
- Analyze round-by-round team performance
- Calculate side-specific statistics

## NEW: Weapon Kills Data

The Series State API provides weapon kill statistics:

```graphql
weaponKills {
  id
  weaponName
  count
}
```

Example data:
```json
{
  "id": "weaponKills-vandal",
  "weaponName": "vandal",
  "count": 11
}
```

## NEW: Multikills Data

The Series State API provides multikill statistics:

```graphql
multikills {
  id
  numberOfKills  # 2=double, 3=triple, 4=quadra, 5=penta
  count
}
```

Example data:
```json
{
  "id": "multikill-2",
  "numberOfKills": 2,
  "count": 1
}
```

GRID_API_KEY=hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO curl -s -X POST https://api-op.grid.gg/live-data-feed/series-state/graphql \
  -H "Content-Type: application/json" \
  -H "x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO" \
  -d '{"query":"query { seriesState(id: \"2629390\") { id finished games { id map { name } teams { id name players { id name character { name } kills deaths weaponKills { weaponName count } abilities { id } } } } } }"}' | jq '.'