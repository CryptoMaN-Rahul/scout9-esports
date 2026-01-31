# GRID API Schema Documentation

## Validated: January 1, 2026

This document contains the verified GRID API schema based on actual API introspection and testing.

---

## Series State API

**Endpoint:** `https://api-op.grid.gg/live-data-feed/series-state/graphql`

### GamePlayerState Fields

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Player ID | ✅ Captured |
| `name` | String! | Player name | ✅ Captured |
| `character` | Character | Champion/Agent | ✅ Captured |
| `roles` | [String]! | Player roles | ❌ Not available via API |
| `participationStatus` | ParticipationStatus! | Participation status | ❌ Not available via API |
| `money` | Int! | Current money | ✅ Captured |
| `loadoutValue` | Int! | Loadout value | ✅ Captured |
| `netWorth` | Int! | Net worth | ✅ Captured |
| `kills` | Int! | Kill count | ✅ Captured |
| `killAssistsReceived` | Int! | Assists received | ✅ Captured |
| `killAssistsGiven` | Int! | Assists given | ✅ Captured |
| `killAssistsReceivedFromPlayer` | [Object]! | Assist details | ✅ Captured (NEW) |
| `weaponKills` | [Object]! | Kills by weapon | ✅ Captured (VALORANT) |
| `teamkills` | Int! | Team kills | ✅ Captured (usually 0) |
| `selfkills` | Int! | Self kills | ✅ Captured (usually 0) |
| `deaths` | Int! | Death count | ✅ Captured |
| `firstKill` | Boolean! | Got first kill | ⚠️ Version restricted (3.10+) |
| `structuresDestroyed` | Int! | Structures destroyed | ✅ Captured |
| `structuresCaptured` | Int! | Structures captured | ❌ Not available via API |
| `inventory` | PlayerInventory! | Player items | ✅ Captured |
| `objectives` | [Objective]! | Objectives completed | ✅ Captured |
| `position` | Coordinates | Current position | ❌ Not available via API |
| `multikills` | [Object]! | Multi-kill stats | ✅ Captured |
| `unitKills` | [Object]! | Unit kills | ❌ Not available via API |
| `abilities` | [Object]! | Abilities used | ✅ Captured (IDs only) |
| `statusEffects` | [Object]! | Status effects | ❌ Not available via API |

### GameTeamState Fields

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Team ID | ✅ Captured |
| `name` | String! | Team name | ✅ Captured |
| `side` | String! | Side (blue/red, attack/defense) | ✅ Captured |
| `won` | Boolean! | Won the game | ✅ Captured |
| `score` | Int! | Score | ✅ Captured |
| `money` | Int! | Team money | ✅ Captured |
| `loadoutValue` | Int! | Team loadout value | ✅ Captured (NEW) |
| `netWorth` | Int! | Team net worth | ✅ Captured |
| `kills` | Int! | Team kills | ✅ Captured |
| `deaths` | Int! | Team deaths | ✅ Captured |
| `structuresDestroyed` | Int! | Structures destroyed | ✅ Captured |
| `objectives` | [Objective]! | Team objectives | ✅ Captured |
| `players` | [GamePlayerState]! | Players | ✅ Captured |

### GameState Fields

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Game ID | ✅ Captured |
| `sequenceNumber` | Int! | Game sequence | ⚠️ Not available (use index) |
| `map` | MapState! | Map info | ✅ Captured |
| `titleVersion` | TitleVersion | Game version | ❌ Not available via API |
| `type` | GameType | Game type | ❌ Not available via API |
| `started` | Boolean! | Game started | ✅ Captured |
| `finished` | Boolean! | Game finished | ✅ Captured |
| `forfeited` | Boolean! | Game forfeited | ❌ Not available via API |
| `paused` | Boolean! | Game paused | ✅ Captured |
| `startedAt` | DateTime | Start time | ⚠️ Version restricted (3.7+) |
| `clock` | ClockState | Game clock | ✅ Captured |
| `structures` | [Object]! | Structures | ❌ Not available via API |
| `nonPlayerCharacters` | [Object]! | NPCs | ❌ Not available via API |
| `teams` | [GameTeamState]! | Teams | ✅ Captured |
| `draftActions` | [DraftAction]! | Draft picks/bans | ✅ Captured |
| `segments` | [Object]! | Game segments/rounds | ✅ Captured (VALORANT) |
| `duration` | Duration! | Game duration | ⚠️ Version restricted (3.15+) |

### DraftAction Fields

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Action ID | ✅ Captured |
| `type` | String! | "ban" or "pick" | ✅ Captured |
| `sequenceNumber` | String! | Draft order | ✅ Captured |
| `drafter` | Drafter! | Team that drafted | ✅ Captured (ID only) |
| `draftable` | Draftable! | Champion/Agent | ✅ Captured |

### Segment Fields (VALORANT Rounds)

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Segment ID | ✅ Captured |
| `sequenceNumber` | Int! | Round number | ✅ Captured |
| `type` | String! | Segment type | ✅ Captured |
| `finished` | Boolean! | Round finished | ✅ Captured |
| `teams` | [SegmentTeam]! | Team states | ✅ Captured |

### Objective Fields

| Field | Type | Description | Status |
|-------|------|-------------|--------|
| `id` | ID! | Objective ID | ✅ Captured |
| `type` | String! | Objective type | ✅ Captured |
| `completedFirst` | Boolean! | First to complete | ⚠️ Version restricted (3.11+) |
| `completionCount` | Int! | Times completed | ✅ Captured |

---

## Newly Captured Fields (January 2026)

### High Impact - Now Available

1. **killAssistsReceivedFromPlayer** - Assist network analysis
   - Shows which teammates assisted on each player's kills
   - Enables synergy analysis between players
   - Available for both LoL and VALORANT

2. **abilities** - Ability/skill usage tracking
   - Returns ability IDs used by each player
   - LoL: Champion ability IDs (e.g., "k-sante-q", "xin-zhao-r")
   - VALORANT: Agent ability IDs (e.g., "poison-cloud", "blade-storm")

3. **loadoutValue** - Economy tracking
   - Team and player loadout values
   - Useful for economy analysis in VALORANT

4. **segments** - Round-by-round data (VALORANT)
   - Per-round team and player stats
   - Enables round-level analysis

### Medium Impact - Now Available

1. **teamkills/selfkills** - Error tracking (usually 0)
2. **weaponKills** - VALORANT weapon usage stats
3. **multikills** - Multi-kill tracking

### Not Available via API

These fields are documented in the schema but return errors or empty data:
- `roles` - Player role assignments
- `participationStatus` - Substitute tracking
- `position` - Real-time positioning
- `unitKills` - CS/farming stats
- `statusEffects` - CC tracking
- `structuresCaptured` - Objective control
- `structures` - Game-level structure states
- `nonPlayerCharacters` - NPC tracking
- `titleVersion` - Patch version

---

## File Download API (Events JSONL)

**Endpoint:** `https://api.grid.gg/file-download/events/grid/series/{seriesId}`

### Event Structure

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
      "actor": { "id": "...", "type": "player", "state": {...} },
      "target": { "id": "...", "type": "player", "state": {...} },
      "seriesState": {...}
    }
  ]
}
```

### Event Type Convention

Event types follow the pattern: `{actor.type}-{action}-{target.type}`

### League of Legends Event Types

| Event Type | Description | Data Available |
|------------|-------------|----------------|
| `player-killed-player` | Champion kill | Killer/victim info, positions, assists |
| `player-killed-ATierNPC` | Dragon/Baron/Herald kill | Objective type, killer info |
| `player-killed-BTierNPC` | Scuttle crab kill | Killer info |
| `player-destroyed-tower` | Tower destruction | Tower ID, lane, killer |
| `team-destroyed-tower` | Tower destruction (team credit) | Tower ID, lane |
| `team-picked-character` | Draft pick | Team, champion |
| `team-banned-character` | Draft ban | Team, champion |
| `player-completed-slayInfernalDrake` | Infernal dragon | Player, timing |
| `player-completed-slayMountainDrake` | Mountain dragon | Player, timing |
| `player-completed-slayOceanDrake` | Ocean dragon | Player, timing |
| `player-completed-slayBaron` | Baron kill | Player, timing |
| `player-completed-slayRiftHerald` | Herald kill | Player, timing |
| `player-completed-slayVoidGrub` | Void grub kill | Player, timing |
| `player-purchased-item` | Item purchase | Player, item |
| `player-used-ability` | Ability usage | Player, ability |
| `player-completed-increaseLevel` | Level up | Player, level |

### VALORANT Event Types

| Event Type | Description | Data Available |
|------------|-------------|----------------|
| `player-killed-player` | Player kill | Killer/victim info, positions, agents |
| `player-selfkilled-player` | Self-kill | Player info |
| `player-completed-plantBomb` | Spike plant | Player, position, site inference |
| `player-completed-defuseBomb` | Spike defuse | Player, position |
| `player-completed-explodeBomb` | Spike explosion | - |
| `player-completed-beginDefuseBomb` | Defuse started | Player |
| `player-completed-stopDefuseBomb` | Defuse stopped | Player |
| `team-won-round` | Round win | Winner, win type |
| `game-started-round` | Round start | Round number, map |
| `game-ended-round` | Round end | - |
| `player-used-ability` | Ability usage | Player, ability |

### Win Types (VALORANT)

- `opponentEliminated` - All enemies killed
- `bombDefused` - Spike defused
- `bombExploded` - Spike detonated
- `timeExpired` - Time ran out

---

## Objective Types

### League of Legends

| Objective Type | Description |
|----------------|-------------|
| `slayInfernalDrake` | Infernal dragon |
| `slayMountainDrake` | Mountain dragon |
| `slayOceanDrake` | Ocean dragon |
| `slayCloudDrake` | Cloud dragon |
| `slayChemtechDrake` | Chemtech dragon |
| `slayHextechDrake` | Hextech dragon |
| `slayElderDrake` | Elder dragon |
| `slayBaron` | Baron Nashor |
| `slayRiftHerald` | Rift Herald |
| `slayVoidGrub` | Void Grub |
| `slayRiftScuttlerTop` | Top scuttle crab |
| `slayRiftScuttlerBot` | Bot scuttle crab |
| `destroyTower` | Tower destruction |
| `destroyFortifier` | Inhibitor destruction |
| `destroyNexus` | Nexus destruction |
| `destroyTurretPlateTop` | Top turret plate |
| `destroyTurretPlateMid` | Mid turret plate |
| `destroyTurretPlateBot` | Bot turret plate |
| `increaseLevel` | Level up |

### VALORANT

| Objective Type | Description |
|----------------|-------------|
| `plantBomb` | Spike plant |
| `defuseBomb` | Spike defuse |
| `explodeBomb` | Spike explosion |
| `beginDefuseBomb` | Defuse started |
| `reachDefuseBombCheckpoint` | Defuse checkpoint |
| `stopDefuseBomb` | Defuse interrupted |

---

## Data Quality Notes

### Position Data

- Available in kill events via `actor.state.game.position` (JSONL files only)
- NOT available via Series State API
- Coordinates are game-specific (LoL uses different scale than VALORANT)
- VALORANT site inference uses position + map bounds

### Timing Data

- `occurredAt` is wall-clock time
- `clock.currentSeconds` is in-game time
- Game time calculation: `occurredAt - gameStartTime` or from `seriesState.games[0].clock.currentSeconds`

### Version Restrictions

Some fields require specific API versions:
- `firstKill`: 3.10+
- `completedFirst`: 3.11+
- `startedAt`: 3.7+
- `duration`: 3.15+

These fields will return errors if queried on older series data.

---

## Implementation Notes

### GraphQL Query (Current)

The following fields are now captured in `pkg/grid/series_state.go`:

```graphql
query SeriesState($seriesId: ID!) {
  seriesState(id: $seriesId) {
    id, started, finished, format
    teams { id, name, won, score, kills, deaths }
    games {
      id, started, finished, paused
      map { name }
      clock { currentSeconds }
      draftActions { id, type, sequenceNumber, drafter { id }, draftable { id, name } }
      segments { id, sequenceNumber, type, finished, teams { ... } }
      teams {
        id, name, side, score, won, kills, deaths
        netWorth, money, loadoutValue, structuresDestroyed
        objectives { id, type, completionCount }
        players {
          id, name, character { id, name }
          kills, deaths, killAssistsReceived, killAssistsGiven
          killAssistsReceivedFromPlayer { playerId, killAssistsReceived }
          teamkills, selfkills, netWorth, money, loadoutValue, structuresDestroyed
          objectives { id, type, completionCount }
          multikills { id, numberOfKills, count }
          weaponKills { id, weaponName, count }
          abilities { id }
          inventory { items { id, name } }
        }
      }
    }
  }
}
```

curl -s -X POST https://api-op.grid.gg/central-data/graphql \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO' \
  -d '{"query": "query SeriesForTeam { allSeries(filter: { teamId: \"47435\" }, orderBy: StartTimeScheduled, orderDirection: DESC, first: 5) { totalCount edges { node { id startTimeScheduled teams { baseInfo { id name } } tournament { id name } title { id name } } } } }"}' | python3 -m json.tool{
    "data": {
        "allSeries": {
            "totalCount": 61,
            "edges": [
                {
                    "node": {
                        "id": "2833788",
                        "startTimeScheduled": "2025-09-05T15:00:00Z",
                        "teams": [
                            {
                                "baseInfo": {
                                    "id": "47435",
                                    "name": "Team Heretics"
                                }
                            },
                            {
                                "baseInfo": {
                                    "id": "47370",
                                    "name": "Team Vitality"
                                }
                            }
                        ],
                        "tournament": {
                            "id": "826911",
                            "name": "LEC - Summer 2025 (Playoffs: Playoffs)"
                        },
                        "title": {
                            "id": "3",
                            "name": "League of Legends"
                        }
                    }
                },
                {
                    "node": {
                        "id": "2833783",
                        "startTimeScheduled": "2025-08-25T17:00:00Z",
                        "teams": [
                            {
                                "baseInfo": {
                                    "id": "52661",
                                    "name": "Shifters"
                                }
                            },
                            {
                                "baseInfo": {
                                    "id": "47435",
                                    "name": "Team Heretics"
                                }
                            }
                        ],
                        "tournament": {
                            "id": "826909",
                            "name": "LEC - Summer 2025 (Groups: Group 2)"
                        },
                        "title": {
                            "id": "3",
                            "name": "League of Legends"
                        }
                    }
                },
                {
                    "node": {
                        "id": "2833776",
                        "startTimeScheduled": "2025-08-16T17:00:00Z",
                        "teams": [
                            {
                                "baseInfo": {
                                    "id": "47435",
                                    "name": "Team Heretics"
                                }
                            },
                            {
                                "baseInfo": {
                                    "id": "353",
                                    "name": "SK Gaming"
                                }
                            }
                        ],
                        "tournament": {
                            "id": "826909",
                            "name": "LEC - Summer 2025 (Groups: Group 2)"
                        },
                        "title": {
                            "id": "3",
                            "name": "League of Legends"
                        }
                    }
                },
                {
                    "node": {
                        "id": "2833771",
                        "startTimeScheduled": "2025-08-09T15:00:00Z",
                        "teams": [
                            {
                                "baseInfo": {
                                    "id": "47380",
                                    "name": "G2 Esports"
                                }
                            },
                            {
                                "baseInfo": {
                                    "id": "47435",
                                    "name": "Team Heretics"
                                }
                            }
                        ],
                        "tournament": {
                            "id": "826909",
                            "name": "LEC - Summer 2025 (Groups: Group 2)"
                        },
                        "title": {
                            "id": "3",
                            "name": "League of Legends"
                        }
                    }
                },
                {
                    "node": {
                        "id": "2833764",
                        "startTimeScheduled": "2025-08-02T17:00:00Z",
                        "teams": [
                            {
                                "baseInfo": {
                                    "id": "47435",
                                    "name": "Team Heretics"
                                }
                            },
                            {
                                "baseInfo": {
                                    "id": "47376",
                                    "name": "Fnatic"
                                }
                            }
                        ],
                        "tournament": {
                            "id": "826909",
                            "name": "LEC - Summer 2025 (Groups: Group 2)"
                        },
                        "title": {
                            "id": "3",
                            "name": "League of Legends"
                        }
                    }
                }
            ]
        }
    }
}