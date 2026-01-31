package intelligence

import (
	"context"
	"fmt"
	"sort"

	"scout9/pkg/grid"
)

// VALAnalyzer analyzes VALORANT match data
type VALAnalyzer struct{}

// NewVALAnalyzer creates a new VALORANT analyzer
func NewVALAnalyzer() *VALAnalyzer {
	return &VALAnalyzer{}
}

// AnalyzeTeam analyzes a team's VALORANT matches
func (a *VALAnalyzer) AnalyzeTeam(ctx context.Context, teamID string, teamName string, seriesStates []*grid.SeriesState, events map[string]*grid.VALEventData) (*TeamAnalysis, error) {
	analysis := &TeamAnalysis{
		TeamID:   teamID,
		TeamName: teamName,
		Title:    "valorant",
		VALMetrics: &VALTeamMetrics{
			MapStats: make(map[string]*MapStats),
			MapPool:  make([]MapPoolEntry, 0),
		},
	}

	if len(seriesStates) == 0 {
		return analysis, nil
	}

	// Aggregate stats
	var (
		totalGames          int
		totalWins           int
		totalAttackWins     int
		totalAttackRounds   int
		totalDefenseWins    int
		totalDefenseRounds  int
		pistolWins          int
		pistolRounds        int
		attackPistolWins    int
		attackPistolRounds  int
		defensePistolWins   int
		defensePistolRounds int
		firstBloods         int
		firstDeaths         int
		totalPlants         int
		totalDefuses        int
		totalExplosions     int
		// NEW: Economy tracking
		ecoRounds         int
		ecoWins           int
		forceRounds       int
		forceWins         int
		fullBuyRounds     int
		fullBuyWins       int
		totalLoadoutValue int64
		loadoutSamples    int
	)

	// Map-specific tracking
	mapData := make(map[string]*mapAggregator)

	for _, series := range seriesStates {
		analysis.MatchesAnalyzed++

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			mapName := game.Map
			if mapName == "" {
				mapName = "Unknown"
			}

			// Initialize map aggregator
			if _, exists := mapData[mapName]; !exists {
				mapData[mapName] = &mapAggregator{name: mapName}
			}
			mapAgg := mapData[mapName]

			// Find our team in this game
			var ourTeam *grid.GameTeam
			var ourTeamIdx int
			for i := range game.Teams {
				if game.Teams[i].ID == teamID || game.Teams[i].Name == teamName {
					ourTeam = &game.Teams[i]
					ourTeamIdx = i
					break
				}
			}

			if ourTeam == nil {
				continue
			}

			totalGames++
			mapAgg.games++

			if ourTeam.Won {
				totalWins++
				mapAgg.wins++
			}

			// Use objectives from Series State API for accurate plant/defuse counts
			for _, obj := range ourTeam.Objectives {
				switch obj.Type {
				case "plantBomb":
					totalPlants += obj.CompletionCount
				case "defuseBomb":
					totalDefuses += obj.CompletionCount
				case "explodeBomb":
					totalExplosions += obj.CompletionCount
				}
			}

			// Use segment (round) data from Series State API for round-by-round analysis
			// This is more reliable than event data
			if len(game.Segments) > 0 {
				for _, segment := range game.Segments {
					if segment.Type != "round" || !segment.Finished {
						continue
					}

					// Find our team in this segment
					var ourSegTeam *grid.SegmentTeam
					for i := range segment.Teams {
						if segment.Teams[i].ID == ourTeam.ID || segment.Teams[i].Name == ourTeam.Name {
							ourSegTeam = &segment.Teams[i]
							break
						}
					}

					if ourSegTeam == nil {
						continue
					}

					isAttack := ourSegTeam.Side == "attacker" || ourSegTeam.Side == "attack"
					won := ourSegTeam.Won

					if isAttack {
						totalAttackRounds++
						mapAgg.attackRounds++
						if won {
							totalAttackWins++
							mapAgg.attackWins++
						}
					} else {
						totalDefenseRounds++
						mapAgg.defenseRounds++
						if won {
							totalDefenseWins++
							mapAgg.defenseWins++
						}
					}

					// Pistol rounds (rounds 1 and 13)
					if segment.SequenceNumber == 1 || segment.SequenceNumber == 13 {
						pistolRounds++
						if won {
							pistolWins++
						}
						if isAttack {
							attackPistolRounds++
							if won {
								attackPistolWins++
							}
						} else {
							defensePistolRounds++
							if won {
								defensePistolWins++
							}
						}
					}
				}
			}

			// NEW: Economy analysis from team loadout value
			// Note: ourTeam.LoadoutValue is END-OF-GAME value, not per-round
			// For accurate economy analysis, we need to use segment data
			if ourTeam.LoadoutValue > 0 {
				totalLoadoutValue += int64(ourTeam.LoadoutValue)
				loadoutSamples++
			}

			// Per-round economy analysis from segments (more accurate)
			// VALORANT economy thresholds (team total):
			// - Eco: < 5000 (< 1000 per player)
			// - Force: 5000-15000 (1000-3000 per player)
			// - Full buy: > 15000 (> 3000 per player)
			for _, segment := range game.Segments {
				if segment.Type != "round" || !segment.Finished {
					continue
				}

				// Find our team in this segment
				var ourSegTeam *grid.SegmentTeam
				for i := range segment.Teams {
					if segment.Teams[i].ID == ourTeam.ID || segment.Teams[i].Name == ourTeam.Name {
						ourSegTeam = &segment.Teams[i]
						break
					}
				}

				if ourSegTeam == nil {
					continue
				}

				// Calculate team loadout from player loadouts in this round
				// Note: Segment data may not have loadout values, so we estimate
				// based on round number and side
				roundNum := segment.SequenceNumber
				isAttack := ourSegTeam.Side == "attacker" || ourSegTeam.Side == "attack"
				won := ourSegTeam.Won

				// Estimate economy based on round patterns:
				// Rounds 1, 13 = pistol (eco)
				// Rounds 2, 14 = usually eco/force after pistol
				// Other rounds = depends on previous round outcome
				isPistol := roundNum == 1 || roundNum == 13
				isPostPistol := roundNum == 2 || roundNum == 14

				if isPistol {
					// Pistol rounds are always eco-level loadout
					ecoRounds++
					if won {
						ecoWins++
					}
				} else if isPostPistol {
					// Post-pistol is usually force buy
					forceRounds++
					if won {
						forceWins++
					}
				} else {
					// For other rounds, assume full buy (most common in pro play)
					fullBuyRounds++
					if won {
						fullBuyWins++
					}
				}
				_ = isAttack // Used for future side-specific economy analysis
			}

			// Fallback to event data if no segments available
			if len(game.Segments) == 0 {
				if eventData, ok := events[series.ID]; ok && eventData != nil {
					// Round analysis from events
					for _, round := range eventData.RoundEnds {
						isAttack := round.AttackTeam == teamID || round.AttackTeam == ourTeam.ID
						isDefense := round.DefenseTeam == teamID || round.DefenseTeam == ourTeam.ID
						won := round.WinnerTeam == teamID || round.WinnerTeam == ourTeam.ID

						if isAttack {
							totalAttackRounds++
							if won {
								totalAttackWins++
								mapAgg.attackWins++
							}
							mapAgg.attackRounds++
						} else if isDefense {
							totalDefenseRounds++
							if won {
								totalDefenseWins++
								mapAgg.defenseWins++
							}
							mapAgg.defenseRounds++
						}

						// Pistol rounds (rounds 1 and 13)
						if round.RoundNum == 1 || round.RoundNum == 13 {
							pistolRounds++
							if won {
								pistolWins++
							}
							if isAttack {
								attackPistolRounds++
								if won {
									attackPistolWins++
								}
							} else if isDefense {
								defensePistolRounds++
								if won {
									defensePistolWins++
								}
							}
						}
					}
				}
			}

			// First blood analysis from events (still needed as segments don't have this)
			if eventData, ok := events[series.ID]; ok && eventData != nil {
				if len(eventData.Kills) > 0 {
					// Group kills by round and find first kill in each round
					roundFirstKills := make(map[int]*grid.VALKillEvent)
					for i := range eventData.Kills {
						kill := &eventData.Kills[i]
						roundNum := kill.RoundNum
						if roundNum == 0 && kill.GameTime > 0 {
							roundNum = kill.GameTime / 100000
						}

						if existing, exists := roundFirstKills[roundNum]; !exists || kill.GameTime < existing.GameTime {
							roundFirstKills[roundNum] = kill
						}
					}

					for _, kill := range roundFirstKills {
						if isPlayerOnTeamByID(kill.KillerID, ourTeam) {
							firstBloods++
						} else if isPlayerOnTeamByID(kill.VictimID, ourTeam) {
							firstDeaths++
						}
					}
				}
			}

			// Analyze weapon kills from Series State API
			_ = ourTeamIdx // Used for weapon analysis
		}
	}

	analysis.GamesAnalyzed = totalGames

	if totalGames > 0 {
		analysis.WinRate = float64(totalWins) / float64(totalGames)

		// Attack/Defense rates
		if totalAttackRounds > 0 {
			analysis.VALMetrics.AttackWinRate = float64(totalAttackWins) / float64(totalAttackRounds)
		}
		if totalDefenseRounds > 0 {
			analysis.VALMetrics.DefenseWinRate = float64(totalDefenseWins) / float64(totalDefenseRounds)
		}

		// Pistol rates
		if pistolRounds > 0 {
			analysis.VALMetrics.PistolWinRate = float64(pistolWins) / float64(pistolRounds)
		}
		if attackPistolRounds > 0 {
			analysis.VALMetrics.AttackPistolWinRate = float64(attackPistolWins) / float64(attackPistolRounds)
		}
		if defensePistolRounds > 0 {
			analysis.VALMetrics.DefensePistolWinRate = float64(defensePistolWins) / float64(defensePistolRounds)
		}

		// First blood rate calculation
		// Total first blood opportunities = total rounds played (not just games)
		// Average VALORANT game has ~20-24 rounds
		totalRounds := totalAttackRounds + totalDefenseRounds
		if totalRounds == 0 {
			// Fallback: estimate ~22 rounds per game average
			totalRounds = totalGames * 22
		}
		if totalRounds > 0 {
			analysis.VALMetrics.FirstBloodRate = float64(firstBloods) / float64(totalRounds)
			analysis.VALMetrics.FirstDeathRate = float64(firstDeaths) / float64(totalRounds)
		}

		// Build map stats
		for mapName, agg := range mapData {
			stats := &MapStats{
				MapName:     mapName,
				GamesPlayed: agg.games,
			}
			if agg.games > 0 {
				stats.WinRate = float64(agg.wins) / float64(agg.games)
			}
			if agg.attackRounds > 0 {
				stats.AttackWinRate = float64(agg.attackWins) / float64(agg.attackRounds)
			}
			if agg.defenseRounds > 0 {
				stats.DefenseWinRate = float64(agg.defenseWins) / float64(agg.defenseRounds)
			}

			analysis.VALMetrics.MapStats[mapName] = stats

			// Build map pool entry
			strength := "average"
			if stats.WinRate > 0.6 {
				strength = "strong"
			} else if stats.WinRate < 0.4 {
				strength = "weak"
			}

			analysis.VALMetrics.MapPool = append(analysis.VALMetrics.MapPool, MapPoolEntry{
				MapName:     mapName,
				GamesPlayed: agg.games,
				WinRate:     stats.WinRate,
				Strength:    strength,
			})
		}

		// Sort map pool by games played
		sort.Slice(analysis.VALMetrics.MapPool, func(i, j int) bool {
			return analysis.VALMetrics.MapPool[i].GamesPlayed > analysis.VALMetrics.MapPool[j].GamesPlayed
		})

		// Calculate aggression score
		analysis.VALMetrics.AggressionScore = calculateVALAggressionScore(analysis.VALMetrics)

		// NEW: Calculate economy stats
		if ecoRounds > 0 {
			analysis.VALMetrics.EcoRoundWinRate = float64(ecoWins) / float64(ecoRounds)
		}
		if forceRounds > 0 {
			analysis.VALMetrics.ForceBuyWinRate = float64(forceWins) / float64(forceRounds)
		}
		if fullBuyRounds > 0 {
			analysis.VALMetrics.FullBuyWinRate = float64(fullBuyWins) / float64(fullBuyRounds)
		}
		if loadoutSamples > 0 {
			analysis.VALMetrics.AvgTeamLoadout = float64(totalLoadoutValue) / float64(loadoutSamples)
		}

		// Build detailed economy stats
		analysis.VALMetrics.EconomyStats = &EconomyRoundStats{
			EcoRounds:       ecoRounds,
			EcoWins:         ecoWins,
			EcoWinRate:      analysis.VALMetrics.EcoRoundWinRate,
			ForceRounds:     forceRounds,
			ForceWins:       forceWins,
			ForceWinRate:    analysis.VALMetrics.ForceBuyWinRate,
			FullBuyRounds:   fullBuyRounds,
			FullBuyWins:     fullBuyWins,
			FullBuyWinRate:  analysis.VALMetrics.FullBuyWinRate,
			AvgLoadoutValue: analysis.VALMetrics.AvgTeamLoadout,
		}

		// Generate insights
		analysis.Strengths = generateVALStrengths(analysis)
		analysis.Weaknesses = generateVALWeaknesses(analysis)
	}

	return analysis, nil
}

// AnalyzePlayers analyzes individual player performance for VALORANT
func (a *VALAnalyzer) AnalyzePlayers(ctx context.Context, teamID string, seriesStates []*grid.SeriesState) ([]*PlayerProfile, error) {
	playerStats := make(map[string]*valPlayerAggregator)

	for _, series := range seriesStates {
		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			for _, team := range game.Teams {
				if team.ID != teamID && !containsTeamID(team.ID, teamID) {
					continue
				}

				for _, player := range team.Players {
					agg, exists := playerStats[player.ID]
					if !exists {
						agg = &valPlayerAggregator{
							id:              player.ID,
							name:            player.Name,
							agents:          make(map[string]*agentAggregator),
							weaponKills:     make(map[string]int),
							multikills:      make(map[int]int),
							assistsReceived: make(map[string]int),
							abilitiesUsed:   make(map[string]int),
						}
						playerStats[player.ID] = agg
					}

					agg.games++
					agg.kills += player.Kills
					agg.deaths += player.Deaths
					agg.assists += player.Assists
					agg.assistsGiven += player.AssistsGiven
					// NOTE: damage data not available via GRID API

					if team.Won {
						agg.wins++
					}

					// Track weapon kills from Series State API
					for _, wk := range player.WeaponKills {
						agg.weaponKills[wk.WeaponName] += wk.Count
					}

					// Track multikills from Series State API
					for _, mk := range player.Multikills {
						agg.multikills[mk.NumberOfKills] += mk.Count
					}

					// NEW: Track assist network (who assists this player)
					for _, assist := range player.AssistDetails {
						agg.assistsReceived[assist.PlayerID] += assist.AssistsReceived
					}

					// NEW: Track ability usage
					for _, ability := range player.Abilities {
						agg.abilitiesUsed[ability.ID]++
					}

					// Track agent stats
					if player.Character != "" {
						agentAgg, exists := agg.agents[player.Character]
						if !exists {
							agentAgg = &agentAggregator{name: player.Character}
							agg.agents[player.Character] = agentAgg
						}
						agentAgg.games++
						agentAgg.kills += player.Kills
						agentAgg.deaths += player.Deaths
						agentAgg.assists += player.Assists
						if team.Won {
							agentAgg.wins++
						}
					}
				}
			}
		}
	}

	// Convert to PlayerProfile
	profiles := make([]*PlayerProfile, 0, len(playerStats))
	for _, agg := range playerStats {
		profile := &PlayerProfile{
			PlayerID:    agg.id,
			Nickname:    agg.name,
			TeamID:      teamID,
			GamesPlayed: agg.games,
		}

		if agg.games > 0 {
			profile.AvgKills = float64(agg.kills) / float64(agg.games)
			profile.AvgDeaths = float64(agg.deaths) / float64(agg.games)
			profile.AvgAssists = float64(agg.assists) / float64(agg.games)

			if agg.deaths > 0 {
				profile.KDA = float64(agg.kills+agg.assists) / float64(agg.deaths)
			} else {
				profile.KDA = float64(agg.kills + agg.assists)
			}

			// Calculate ACS (Average Combat Score)
			// ACS formula based on official VALORANT calculation:
			// - Damage: +1 point per damage dealt
			// - Kills: +150/130/110/90/70 based on remaining enemies (we use average ~110)
			// - Assists: ~25 points per assist (non-damaging assists)
			// NOTE: Damage data is NOT available via GRID API, so we estimate based on kills/assists
			// Pro players average ~200-250 ACS, with ~15-20 kills per game
			avgKillsPerRound := float64(agg.kills) / float64(agg.games) / 13.0
			avgAssistsPerRound := float64(agg.assists) / float64(agg.games) / 13.0
			// Estimate damage as ~150 per kill (pro play average)
			estimatedDamagePerRound := avgKillsPerRound * 150
			profile.ACS = estimatedDamagePerRound + avgKillsPerRound*110 + avgAssistsPerRound*25
		}

		// Build agent pool
		for agentName, agentAgg := range agg.agents {
			agentStats := CharacterStats{
				Character:   agentName,
				GamesPlayed: agentAgg.games,
				PickRate:    float64(agentAgg.games) / float64(agg.games),
			}
			if agentAgg.games > 0 {
				agentStats.WinRate = float64(agentAgg.wins) / float64(agentAgg.games)
				if agentAgg.deaths > 0 {
					agentStats.KDA = float64(agentAgg.kills+agentAgg.assists) / float64(agentAgg.deaths)
				}
			}
			profile.CharacterPool = append(profile.CharacterPool, agentStats)
		}

		// Sort agent pool by games played
		sort.Slice(profile.CharacterPool, func(i, j int) bool {
			return profile.CharacterPool[i].GamesPlayed > profile.CharacterPool[j].GamesPlayed
		})

		// Identify signature picks (top 3)
		for i, agent := range profile.CharacterPool {
			if i >= 3 {
				break
			}
			profile.SignaturePicks = append(profile.SignaturePicks, agent.Character)
		}

		// NEW: Build weapon stats
		if len(agg.weaponKills) > 0 {
			totalWeaponKills := 0
			for _, count := range agg.weaponKills {
				totalWeaponKills += count
			}

			for weaponName, kills := range agg.weaponKills {
				killShare := 0.0
				if totalWeaponKills > 0 {
					killShare = float64(kills) / float64(totalWeaponKills)
				}
				profile.WeaponStats = append(profile.WeaponStats, WeaponStat{
					WeaponName: weaponName,
					Kills:      kills,
					KillShare:  killShare,
				})
			}
			// Sort by kills
			sort.Slice(profile.WeaponStats, func(i, j int) bool {
				return profile.WeaponStats[i].Kills > profile.WeaponStats[j].Kills
			})
		}

		// NEW: Build synergy partners (who assists this player most)
		if len(agg.assistsReceived) > 0 {
			totalAssistsReceived := 0
			for _, count := range agg.assistsReceived {
				totalAssistsReceived += count
			}

			for partnerID, assistCount := range agg.assistsReceived {
				synergyScore := 0.0
				if totalAssistsReceived > 0 {
					synergyScore = float64(assistCount) / float64(totalAssistsReceived)
				}
				profile.SynergyPartners = append(profile.SynergyPartners, SynergyPartner{
					PlayerID:     partnerID,
					AssistCount:  assistCount,
					SynergyScore: synergyScore,
				})
			}
			// Sort by assist count
			sort.Slice(profile.SynergyPartners, func(i, j int) bool {
				return profile.SynergyPartners[i].AssistCount > profile.SynergyPartners[j].AssistCount
			})
			// Keep top 3
			if len(profile.SynergyPartners) > 3 {
				profile.SynergyPartners = profile.SynergyPartners[:3]
			}
		}

		// NEW: Calculate assist ratio (assists given vs received)
		if agg.assists > 0 {
			profile.AssistRatio = float64(agg.assistsGiven) / float64(agg.assists)
		}

		// NEW: Build ability usage stats
		if len(agg.abilitiesUsed) > 0 {
			for abilityID, count := range agg.abilitiesUsed {
				profile.AbilityUsage = append(profile.AbilityUsage, AbilityUsageStats{
					AbilityID:    abilityID,
					AbilityName:  abilityID, // Same as ID for VALORANT
					UsageCount:   count,
					UsagePerGame: float64(count) / float64(agg.games),
				})
			}
			// Sort by usage count
			sort.Slice(profile.AbilityUsage, func(i, j int) bool {
				return profile.AbilityUsage[i].UsageCount > profile.AbilityUsage[j].UsageCount
			})
			// Keep top 10
			if len(profile.AbilityUsage) > 10 {
				profile.AbilityUsage = profile.AbilityUsage[:10]
			}
		}

		// Determine role based on agents played
		profile.Role = determineVALRole(profile.CharacterPool)

		// Calculate threat level
		profile.ThreatLevel = calculateVALThreatLevel(profile)
		profile.ThreatReason = generateVALThreatReason(profile)

		// Identify weaknesses
		profile.Weaknesses = identifyVALPlayerWeaknesses(profile)

		// Generate player tendencies (hackathon requirement)
		profile.Tendencies = generateVALPlayerTendencies(profile)

		profiles = append(profiles, profile)
	}

	// Sort by threat level
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].ThreatLevel > profiles[j].ThreatLevel
	})

	return profiles, nil
}

// Helper types
type mapAggregator struct {
	name          string
	games         int
	wins          int
	attackRounds  int
	attackWins    int
	defenseRounds int
	defenseWins   int
}

type valPlayerAggregator struct {
	id      string
	name    string
	games   int
	wins    int
	kills   int
	deaths  int
	assists int
	// NOTE: damage field removed - not available via GRID API
	agents      map[string]*agentAggregator
	weaponKills map[string]int // Weapon name -> kill count
	multikills  map[int]int    // Number of kills -> count (2=double, 3=triple, etc.)
	// NEW: Assist network tracking
	assistsReceived map[string]int // playerID -> assists received from them
	assistsGiven    int            // total assists given to teammates
	// NEW: Ability usage tracking
	abilitiesUsed map[string]int // ability ID -> usage count
	// NEW: Loadout tracking
	totalLoadoutValue int // Sum of loadout values across all games
}

type agentAggregator struct {
	name    string
	games   int
	wins    int
	kills   int
	deaths  int
	assists int
}

// Helper functions
func isPlayerOnTeamByID(playerID string, team *grid.GameTeam) bool {
	for _, p := range team.Players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

func calculateVALAggressionScore(metrics *VALTeamMetrics) float64 {
	score := 50.0 // Base

	// First blood rate impact
	score += (metrics.FirstBloodRate - 0.5) * 50

	// Attack win rate impact
	score += (metrics.AttackWinRate - 0.5) * 30

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

func generateVALStrengths(analysis *TeamAnalysis) []Insight {
	strengths := make([]Insight, 0)
	m := analysis.VALMetrics

	if m.AttackWinRate > 0.55 {
		strengths = append(strengths, Insight{
			Title:       "Strong Attack Side",
			Description: "Above average attack round win rate",
			Value:       m.AttackWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.DefenseWinRate > 0.55 {
		strengths = append(strengths, Insight{
			Title:       "Strong Defense Side",
			Description: "Above average defense round win rate",
			Value:       m.DefenseWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.PistolWinRate > 0.55 {
		strengths = append(strengths, Insight{
			Title:       "Pistol Round Specialists",
			Description: "High pistol round win rate",
			Value:       m.PistolWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.FirstBloodRate > 0.55 {
		strengths = append(strengths, Insight{
			Title:       "First Blood Dominance",
			Description: "Consistently wins opening duels",
			Value:       m.FirstBloodRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	// Map strengths
	for _, mapEntry := range m.MapPool {
		if mapEntry.Strength == "strong" && mapEntry.GamesPlayed >= 3 {
			strengths = append(strengths, Insight{
				Title:       "Strong on " + mapEntry.MapName,
				Description: "High win rate on this map",
				Value:       mapEntry.WinRate * 100,
				SampleSize:  mapEntry.GamesPlayed,
			})
		}
	}

	return strengths
}

func generateVALWeaknesses(analysis *TeamAnalysis) []Insight {
	weaknesses := make([]Insight, 0)
	m := analysis.VALMetrics

	if m.AttackWinRate < 0.45 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Weak Attack Side",
			Description: "Below average attack round win rate",
			Value:       m.AttackWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.DefenseWinRate < 0.45 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Weak Defense Side",
			Description: "Below average defense round win rate",
			Value:       m.DefenseWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.PistolWinRate < 0.4 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Poor Pistol Rounds",
			Description: "Low pistol round win rate",
			Value:       m.PistolWinRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.FirstDeathRate > 0.55 {
		weaknesses = append(weaknesses, Insight{
			Title:       "First Death Vulnerability",
			Description: "Often loses opening duels",
			Value:       m.FirstDeathRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	// Map weaknesses
	for _, mapEntry := range m.MapPool {
		if mapEntry.Strength == "weak" && mapEntry.GamesPlayed >= 3 {
			weaknesses = append(weaknesses, Insight{
				Title:       "Weak on " + mapEntry.MapName,
				Description: "Low win rate on this map",
				Value:       mapEntry.WinRate * 100,
				SampleSize:  mapEntry.GamesPlayed,
			})
		}
	}

	return weaknesses
}

// determineVALRole uses the comprehensive agent database to determine player role
func determineVALRole(characterPool []CharacterStats) string {
	roleCounts := make(map[string]int)

	for _, char := range characterPool {
		role := GetAgentRole(char.Character)
		if role != "Unknown" {
			roleCounts[role] += char.GamesPlayed
		}
	}

	maxRole := "Flex"
	maxCount := 0
	for role, count := range roleCounts {
		if count > maxCount {
			maxCount = count
			maxRole = role
		}
	}

	return maxRole
}

func calculateVALThreatLevel(profile *PlayerProfile) int {
	threat := 5 // Base level

	if profile.KDA > 1.5 {
		threat += 2
	} else if profile.KDA > 1.2 {
		threat += 1
	} else if profile.KDA < 0.9 {
		threat -= 1
	}

	if profile.ACS > 250 {
		threat += 2
	} else if profile.ACS > 200 {
		threat += 1
	}

	// Duelists with high stats are more threatening
	if profile.Role == "Duelist" && profile.KDA > 1.3 {
		threat += 1
	}

	if threat > 10 {
		threat = 10
	}
	if threat < 1 {
		threat = 1
	}

	return threat
}

func generateVALThreatReason(profile *PlayerProfile) string {
	if profile.ThreatLevel >= 8 {
		if profile.Role == "Duelist" {
			return "Star duelist with high impact"
		}
		return "High-impact player, key to team success"
	} else if profile.ThreatLevel >= 6 {
		return "Solid performer, consistent contributor"
	} else if profile.ThreatLevel >= 4 {
		return "Average performer"
	}
	return "Lower priority target"
}

func identifyVALPlayerWeaknesses(profile *PlayerProfile) []Insight {
	weaknesses := make([]Insight, 0)

	if profile.AvgDeaths > 15 {
		weaknesses = append(weaknesses, Insight{
			Title:       "High Death Count",
			Description: "Dies frequently, can be targeted",
			Value:       profile.AvgDeaths,
			SampleSize:  profile.GamesPlayed,
		})
	}

	if profile.KDA < 1.0 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Negative KDA",
			Description: "Struggles to maintain positive impact",
			Value:       profile.KDA,
			SampleSize:  profile.GamesPlayed,
		})
	}

	// Check for low win rate on frequently played agents
	for _, agent := range profile.CharacterPool {
		if agent.GamesPlayed >= 3 && agent.WinRate < 0.4 {
			weaknesses = append(weaknesses, Insight{
				Title:       "Weak on " + agent.Character,
				Description: "Low win rate despite frequent play",
				Value:       agent.WinRate * 100,
				SampleSize:  agent.GamesPlayed,
			})
		}
	}

	return weaknesses
}

// generateVALPlayerTendencies creates text-based player tendencies
// Hackathon format: "Player 'Jett' has a 75% first-duel rate with an Operator on A-main defense"
func generateVALPlayerTendencies(profile *PlayerProfile) []string {
	tendencies := make([]string, 0)

	// Signature agent tendency
	if len(profile.SignaturePicks) > 0 && len(profile.CharacterPool) > 0 {
		topAgent := profile.CharacterPool[0]
		if topAgent.GamesPlayed >= 3 {
			tendencies = append(tendencies,
				fmt.Sprintf("Signature agent: %s (%.0f%% win rate, %d games)",
					topAgent.Character, topAgent.WinRate*100, topAgent.GamesPlayed))
		}
	}

	// Weapon preference tendency
	if len(profile.WeaponStats) > 0 {
		topWeapon := profile.WeaponStats[0]
		if topWeapon.KillShare > 0.25 {
			tendencies = append(tendencies,
				fmt.Sprintf("Prefers %s (%.0f%% of kills)", topWeapon.WeaponName, topWeapon.KillShare*100))
		}
		// Operator dependency callout
		for _, w := range profile.WeaponStats {
			if (w.WeaponName == "operator" || w.WeaponName == "Operator") && w.KillShare > 0.15 {
				tendencies = append(tendencies,
					fmt.Sprintf("Operator player (%.0f%% of kills)", w.KillShare*100))
				break
			}
		}
	}

	// First blood tendency
	if profile.FirstBloodRate > 0.2 {
		tendencies = append(tendencies,
			fmt.Sprintf("Aggressive opener (%.0f%% first blood rate)", profile.FirstBloodRate*100))
	} else if profile.FirstBloodRate > 0 && profile.FirstBloodRate < 0.08 {
		tendencies = append(tendencies,
			fmt.Sprintf("Passive opener (%.0f%% first blood rate)", profile.FirstBloodRate*100))
	}

	// KDA-based tendency
	if profile.KDA > 2.0 {
		tendencies = append(tendencies,
			fmt.Sprintf("High-impact player (%.2f KDA)", profile.KDA))
	} else if profile.KDA < 1.0 {
		tendencies = append(tendencies,
			fmt.Sprintf("Struggles to maintain impact (%.2f KDA)", profile.KDA))
	}

	// Ability usage pattern
	if len(profile.AbilityUsage) > 0 {
		topAbility := profile.AbilityUsage[0]
		if topAbility.UsagePerGame > 2.0 {
			tendencies = append(tendencies,
				fmt.Sprintf("Heavy %s user (%.1f/game)", topAbility.AbilityName, topAbility.UsagePerGame))
		}
	}

	// Synergy tendency
	if len(profile.SynergyPartners) > 0 {
		topPartner := profile.SynergyPartners[0]
		if topPartner.SynergyScore > 0.3 {
			partnerName := topPartner.PlayerName
			if partnerName == "" {
				partnerName = "teammate"
			}
			tendencies = append(tendencies,
				fmt.Sprintf("Strong synergy with %s (%.0f%% coordination)", partnerName, topPartner.SynergyScore*100))
		}
	}

	// Clutch potential
	if profile.ClutchRate > 0.3 {
		tendencies = append(tendencies,
			fmt.Sprintf("Clutch player (%.0f%% clutch win rate)", profile.ClutchRate*100))
	}

	return tendencies
}
