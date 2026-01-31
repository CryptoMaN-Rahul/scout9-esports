package intelligence

import (
	"context"
	"fmt"
	"sort"

	"scout9/pkg/grid"
)

// LoLAnalyzer analyzes League of Legends match data
type LoLAnalyzer struct{}

// NewLoLAnalyzer creates a new LoL analyzer
func NewLoLAnalyzer() *LoLAnalyzer {
	return &LoLAnalyzer{}
}

// AnalyzeTeam analyzes a team's LoL matches
func (a *LoLAnalyzer) AnalyzeTeam(ctx context.Context, teamID string, teamName string, seriesStates []*grid.SeriesState, events map[string]*grid.LoLEventData) (*TeamAnalysis, error) {
	analysis := &TeamAnalysis{
		TeamID:   teamID,
		TeamName: teamName,
		Title:    "lol",
		LoLMetrics: &LoLTeamMetrics{
			WinConditions: make([]string, 0),
		},
	}

	if len(seriesStates) == 0 {
		return analysis, nil
	}

	// Create EventAnalyzer for proper first objective detection
	eventAnalyzer := NewEventAnalyzer()

	// Aggregate stats across all games
	var (
		totalGames      int
		totalWins       int
		totalDragons    int
		totalBarons     int
		totalHeralds    int
		totalVoidGrubs  int
		totalGoldDiff15 float64
		goldDiff15Count int
		totalDuration   float64
	)

	// Collect all event data for EventAnalyzer
	var allEventData []*grid.LoLEventData
	var allPhaseAnalyses []*PhaseAnalysis
	var allObjectiveTimings []*EventObjectiveTimings

	for _, series := range seriesStates {
		analysis.MatchesAnalyzed++

		// Collect event data for this series
		if eventData, ok := events[series.ID]; ok && eventData != nil {
			allEventData = append(allEventData, eventData)

			// Analyze phases for this game
			phases := eventAnalyzer.AnalyzePhases(eventData, teamID)
			allPhaseAnalyses = append(allPhaseAnalyses, phases)

			// Analyze objective timings for this game
			timings := eventAnalyzer.AnalyzeObjectiveTimings(eventData, teamID)
			allObjectiveTimings = append(allObjectiveTimings, timings)
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			// Find our team in this game
			var ourTeam *grid.GameTeam
			var enemyTeam *grid.GameTeam
			for i := range game.Teams {
				if game.Teams[i].ID == teamID || game.Teams[i].Name == teamName {
					ourTeam = &game.Teams[i]
				} else {
					enemyTeam = &game.Teams[i]
				}
			}

			if ourTeam == nil {
				continue
			}

			totalGames++
			if ourTeam.Won {
				totalWins++
			}

			totalDuration += float64(game.Duration) / 60.0 // Convert to minutes

			// Use objectives from Series State API (more reliable for total counts)
			for _, obj := range ourTeam.Objectives {
				switch obj.Type {
				case "slayInfernalDrake", "slayMountainDrake", "slayOceanDrake",
					"slayCloudDrake", "slayChemtechDrake", "slayHextechDrake", "slayElderDrake":
					totalDragons += obj.CompletionCount
				case "slayBaron":
					totalBarons += obj.CompletionCount
				case "slayRiftHerald":
					totalHeralds += obj.CompletionCount
				case "slayVoidGrub":
					totalVoidGrubs += obj.CompletionCount
				}
			}

			// Gold diff at 15 (if we have netWorth data)
			if ourTeam.NetWorth > 0 && enemyTeam != nil && enemyTeam.NetWorth > 0 {
				// This is end-game netWorth, we'd need event data for 15-min snapshot
				// For now, use as proxy
				goldDiff15Count++
				totalGoldDiff15 += float64(ourTeam.NetWorth - enemyTeam.NetWorth)
			}
		}
	}

	// Use EventAnalyzer for accurate first objective detection
	// This properly sorts events by GameTime to find the true "first" objective
	firstsAnalysis := eventAnalyzer.AnalyzeFirsts(allEventData, teamID)

	analysis.GamesAnalyzed = totalGames

	if totalGames > 0 {
		analysis.WinRate = float64(totalWins) / float64(totalGames)

		// Use EventAnalyzer results for first objective rates (accurate time-based detection)
		analysis.LoLMetrics.FirstBloodRate = firstsAnalysis.FirstBloodRate
		analysis.LoLMetrics.FirstDragonRate = firstsAnalysis.FirstDragonRate
		analysis.LoLMetrics.FirstTowerRate = firstsAnalysis.FirstTowerRate
		analysis.LoLMetrics.FirstTowerAvgTime = firstsAnalysis.AvgFirstTowerTime / 60.0 // Convert to minutes

		// Herald control rate from first herald analysis
		if firstsAnalysis.FirstBloodGames > 0 {
			// Use first herald count from objective timings
			firstHeralds := 0
			for _, timings := range allObjectiveTimings {
				if len(timings.HeraldTimings) > 0 {
					firstHeralds++
				}
			}
			if len(allObjectiveTimings) > 0 {
				analysis.LoLMetrics.HeraldControlRate = float64(firstHeralds) / float64(len(allObjectiveTimings))
			}
		}

		analysis.LoLMetrics.AvgGameDuration = totalDuration / float64(totalGames)

		// Baron control rate (barons typically spawn after 20 min, so not all games have them)
		if totalBarons > 0 || totalGames > 0 {
			// Estimate baron opportunities (games longer than 20 min)
			baronOpportunities := 0
			for _, series := range seriesStates {
				for _, game := range series.Games {
					if game.Duration > 1200 { // 20 minutes in seconds
						baronOpportunities++
					}
				}
			}
			if baronOpportunities > 0 {
				analysis.LoLMetrics.BaronControlRate = float64(totalBarons) / float64(baronOpportunities)
			}
		}

		if goldDiff15Count > 0 {
			analysis.LoLMetrics.GoldDiff15 = totalGoldDiff15 / float64(goldDiff15Count)
		}

		// Aggregate phase analysis for game phase ratings
		if len(allPhaseAnalyses) > 0 {
			aggregatedPhases := eventAnalyzer.AggregatePhases(allPhaseAnalyses)
			analysis.LoLMetrics.EarlyGameRating = aggregatedPhases.EarlyGameRating
			analysis.LoLMetrics.MidGameRating = aggregatedPhases.MidGameRating
			analysis.LoLMetrics.LateGameRating = aggregatedPhases.LateGameRating
		} else {
			// Fallback to old calculation if no event data
			analysis.LoLMetrics.EarlyGameRating = calculateEarlyGameRating(analysis.LoLMetrics)
		}

		analysis.LoLMetrics.AggressionScore = calculateAggressionScore(analysis.LoLMetrics)

		// Determine win conditions
		analysis.LoLMetrics.WinConditions = determineWinConditions(analysis.LoLMetrics)

		// Generate insights
		analysis.Strengths = generateLoLStrengths(analysis)
		analysis.Weaknesses = generateLoLWeaknesses(analysis)
	}

	return analysis, nil
}

// AnalyzePlayers analyzes individual player performance
func (a *LoLAnalyzer) AnalyzePlayers(ctx context.Context, teamID string, seriesStates []*grid.SeriesState) ([]*PlayerProfile, error) {
	playerStats := make(map[string]*playerAggregator)

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
						agg = &playerAggregator{
							id:              player.ID,
							name:            player.Name,
							characters:      make(map[string]*characterAggregator),
							objectives:      make(map[string]int),
							multikills:      make(map[int]int),
							assistsReceived: make(map[string]int),
							abilitiesUsed:   make(map[string]int),
							itemsBuilt:      make(map[string]int),
						}
						playerStats[player.ID] = agg
					}

					agg.games++
					agg.kills += player.Kills
					agg.deaths += player.Deaths
					agg.assists += player.Assists
					agg.netWorth += player.NetWorth
					agg.structuresDestroyed += player.StructuresDestroyed
					agg.assistsGiven += player.AssistsGiven

					// Track player objectives (dragons, barons, etc.)
					for _, obj := range player.Objectives {
						agg.objectives[obj.Type] += obj.CompletionCount
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

					// NEW: Track item builds
					for _, item := range player.Items {
						if item.Name != "" {
							agg.itemsBuilt[item.Name]++
						}
					}

					if team.Won {
						agg.wins++
					}

					// Track character stats
					if player.Character != "" {
						charAgg, exists := agg.characters[player.Character]
						if !exists {
							charAgg = &characterAggregator{name: player.Character}
							agg.characters[player.Character] = charAgg
						}
						charAgg.games++
						charAgg.kills += player.Kills
						charAgg.deaths += player.Deaths
						charAgg.assists += player.Assists
						if team.Won {
							charAgg.wins++
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

			profile.GoldPerMin = float64(agg.netWorth) / float64(agg.games) / 30.0 // Rough estimate
		}

		// Build character pool
		for charName, charAgg := range agg.characters {
			charStats := CharacterStats{
				Character:   charName,
				GamesPlayed: charAgg.games,
				PickRate:    float64(charAgg.games) / float64(agg.games),
			}
			if charAgg.games > 0 {
				charStats.WinRate = float64(charAgg.wins) / float64(charAgg.games)
				if charAgg.deaths > 0 {
					charStats.KDA = float64(charAgg.kills+charAgg.assists) / float64(charAgg.deaths)
				}
			}
			profile.CharacterPool = append(profile.CharacterPool, charStats)
		}

		// Sort character pool by games played
		sort.Slice(profile.CharacterPool, func(i, j int) bool {
			return profile.CharacterPool[i].GamesPlayed > profile.CharacterPool[j].GamesPlayed
		})

		// Identify signature picks (top 3)
		for i, char := range profile.CharacterPool {
			if i >= 3 {
				break
			}
			profile.SignaturePicks = append(profile.SignaturePicks, char.Character)
		}

		// NEW: Build multikill stats
		if len(agg.multikills) > 0 {
			profile.MultikillStats = &MultikillStats{
				DoubleKills: agg.multikills[2],
				TripleKills: agg.multikills[3],
				QuadraKills: agg.multikills[4],
				PentaKills:  agg.multikills[5],
			}
			profile.MultikillStats.TotalMultikills = profile.MultikillStats.DoubleKills +
				profile.MultikillStats.TripleKills +
				profile.MultikillStats.QuadraKills +
				profile.MultikillStats.PentaKills
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
					AbilityName:  abilityID, // Same as ID for now
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

		// NEW: Build item build stats
		if len(agg.itemsBuilt) > 0 {
			for itemName, count := range agg.itemsBuilt {
				profile.ItemBuilds = append(profile.ItemBuilds, ItemBuildStats{
					ItemName:   itemName,
					ItemID:     itemName,
					BuildCount: count,
					BuildRate:  float64(count) / float64(agg.games),
				})
			}
			// Sort by build count
			sort.Slice(profile.ItemBuilds, func(i, j int) bool {
				return profile.ItemBuilds[i].BuildCount > profile.ItemBuilds[j].BuildCount
			})
			// Keep top 10
			if len(profile.ItemBuilds) > 10 {
				profile.ItemBuilds = profile.ItemBuilds[:10]
			}
		}

		// NEW: Build objective focus stats
		towersDestroyed := agg.structuresDestroyed
		dragonsSecured := 0
		baronsSecured := 0
		heraldsSecured := 0
		for objType, count := range agg.objectives {
			switch objType {
			case "slayInfernalDrake", "slayMountainDrake", "slayOceanDrake",
				"slayCloudDrake", "slayChemtechDrake", "slayHextechDrake", "slayElderDrake":
				dragonsSecured += count
			case "slayBaron":
				baronsSecured += count
			case "slayRiftHerald":
				heraldsSecured += count
			}
		}

		towersPerGame := 0.0
		dragonsPerGame := 0.0
		if agg.games > 0 {
			towersPerGame = float64(towersDestroyed) / float64(agg.games)
			dragonsPerGame = float64(dragonsSecured) / float64(agg.games)
		}

		objectiveFocused := towersPerGame > 1.5 || dragonsPerGame > 1.0
		focusType := "teamfighter"
		if towersPerGame > 2.0 {
			focusType = "split-pusher"
		} else if dragonsPerGame > 1.0 || baronsSecured > 0 {
			focusType = "objective-focused"
		}

		profile.ObjectiveFocus = &ObjectiveFocusStats{
			TowersDestroyed:    towersDestroyed,
			TowersPerGame:      towersPerGame,
			DragonsSecured:     dragonsSecured,
			DragonsPerGame:     dragonsPerGame,
			BaronsSecured:      baronsSecured,
			HeraldsSecured:     heraldsSecured,
			ObjectiveFocused:   objectiveFocused,
			ObjectiveFocusType: focusType,
		}

		// Calculate threat level
		profile.ThreatLevel = calculateThreatLevel(profile)
		profile.ThreatReason = generateThreatReason(profile)

		// Determine role from champion pool
		profile.Role = determineLoLRole(profile.CharacterPool)

		// Identify weaknesses
		profile.Weaknesses = identifyPlayerWeaknesses(profile)

		// Generate player tendencies (hackathon requirement)
		profile.Tendencies = generateLoLPlayerTendencies(profile)

		profiles = append(profiles, profile)
	}

	// Sort by threat level
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].ThreatLevel > profiles[j].ThreatLevel
	})

	return profiles, nil
}

// Helper types for aggregation
type playerAggregator struct {
	id                  string
	name                string
	games               int
	wins                int
	kills               int
	deaths              int
	assists             int
	netWorth            int
	structuresDestroyed int
	objectives          map[string]int // objective type -> count
	characters          map[string]*characterAggregator
	multikills          map[int]int // Number of kills -> count (2=double, 3=triple, etc.)
	// NEW: Assist network tracking
	assistsReceived map[string]int // playerID -> assists received from them
	assistsGiven    int            // total assists given to teammates
	// NEW: Ability usage tracking
	abilitiesUsed map[string]int // ability ID -> usage count
	// NEW: Item build tracking
	itemsBuilt map[string]int // item name -> build count
}

type characterAggregator struct {
	name    string
	games   int
	wins    int
	kills   int
	deaths  int
	assists int
}

// Helper functions
func isPlayerOnTeam(playerID string, team *grid.GameTeam) bool {
	for _, p := range team.Players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

func isTeamID(id string, team *grid.GameTeam) bool {
	return team.ID == id
}

func containsTeamID(id1, id2 string) bool {
	return id1 == id2
}

func calculateEarlyGameRating(metrics *LoLTeamMetrics) float64 {
	// Weight: First blood 30%, First dragon 35%, First tower 35%
	rating := metrics.FirstBloodRate*30 + metrics.FirstDragonRate*35 + metrics.FirstTowerRate*35
	return rating
}

func calculateAggressionScore(metrics *LoLTeamMetrics) float64 {
	// Based on first blood rate and early game metrics
	score := metrics.FirstBloodRate * 50
	if metrics.GoldDiff15 > 0 {
		score += 25
	}
	score += metrics.FirstDragonRate * 25
	return score
}

func determineWinConditions(metrics *LoLTeamMetrics) []string {
	conditions := make([]string, 0)

	if metrics.FirstDragonRate > 0.6 {
		conditions = append(conditions, "Dragon control and soul priority")
	}
	if metrics.FirstBloodRate > 0.5 {
		conditions = append(conditions, "Early aggression and snowballing")
	}
	if metrics.AvgGameDuration > 35 {
		conditions = append(conditions, "Late game scaling and teamfighting")
	} else if metrics.AvgGameDuration < 28 {
		conditions = append(conditions, "Fast-paced early game dominance")
	}

	if len(conditions) == 0 {
		conditions = append(conditions, "Balanced playstyle")
	}

	return conditions
}

func generateLoLStrengths(analysis *TeamAnalysis) []Insight {
	strengths := make([]Insight, 0)
	m := analysis.LoLMetrics

	if m.FirstBloodRate > 0.5 {
		strengths = append(strengths, Insight{
			Title:       "Strong Early Aggression",
			Description: "High first blood rate indicates aggressive early game",
			Value:       m.FirstBloodRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.FirstDragonRate > 0.6 {
		strengths = append(strengths, Insight{
			Title:       "Dragon Control",
			Description: "Excellent at securing first dragon",
			Value:       m.FirstDragonRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.GoldDiff15 > 1000 {
		strengths = append(strengths, Insight{
			Title:       "Early Gold Lead",
			Description: "Consistently builds gold advantages",
			Value:       m.GoldDiff15,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	return strengths
}

func generateLoLWeaknesses(analysis *TeamAnalysis) []Insight {
	weaknesses := make([]Insight, 0)
	m := analysis.LoLMetrics

	if m.FirstBloodRate < 0.4 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Weak Early Game",
			Description: "Low first blood rate suggests passive early game",
			Value:       m.FirstBloodRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.FirstDragonRate < 0.4 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Poor Dragon Control",
			Description: "Struggles to secure first dragon",
			Value:       m.FirstDragonRate * 100,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	if m.GoldDiff15 < -1000 {
		weaknesses = append(weaknesses, Insight{
			Title:       "Early Gold Deficit",
			Description: "Often falls behind in gold early",
			Value:       m.GoldDiff15,
			SampleSize:  analysis.GamesAnalyzed,
		})
	}

	return weaknesses
}

func calculateThreatLevel(profile *PlayerProfile) int {
	// Base threat on KDA and win rate
	threat := 5 // Base level

	if profile.KDA > 4.0 {
		threat += 2
	} else if profile.KDA > 3.0 {
		threat += 1
	} else if profile.KDA < 2.0 {
		threat -= 1
	}

	// Check for high win rate on signature picks
	for _, char := range profile.CharacterPool {
		if char.GamesPlayed >= 3 && char.WinRate > 0.7 {
			threat += 1
			break
		}
	}

	if threat > 10 {
		threat = 10
	}
	if threat < 1 {
		threat = 1
	}

	return threat
}

func generateThreatReason(profile *PlayerProfile) string {
	if profile.ThreatLevel >= 8 {
		return "High-impact player with excellent stats"
	} else if profile.ThreatLevel >= 6 {
		return "Solid performer, key to team success"
	} else if profile.ThreatLevel >= 4 {
		return "Average performer"
	}
	return "Lower priority target"
}

func identifyPlayerWeaknesses(profile *PlayerProfile) []Insight {
	weaknesses := make([]Insight, 0)

	if profile.AvgDeaths > 4.0 {
		weaknesses = append(weaknesses, Insight{
			Title:       "High Death Count",
			Description: "Dies frequently, can be targeted",
			Value:       profile.AvgDeaths,
			SampleSize:  profile.GamesPlayed,
		})
	}

	// Check for low win rate on frequently played characters
	for _, char := range profile.CharacterPool {
		if char.GamesPlayed >= 3 && char.WinRate < 0.4 {
			weaknesses = append(weaknesses, Insight{
				Title:       "Weak on " + char.Character,
				Description: "Low win rate despite frequent play",
				Value:       char.WinRate * 100,
				SampleSize:  char.GamesPlayed,
			})
		}
	}

	return weaknesses
}

// generateLoLPlayerTendencies creates text-based player tendencies
// Hackathon format: "Top laner has a 90% pick/ban rate on Renekton"
func generateLoLPlayerTendencies(profile *PlayerProfile) []string {
	tendencies := make([]string, 0)

	// Signature champion tendency
	if len(profile.SignaturePicks) > 0 && len(profile.CharacterPool) > 0 {
		topChamp := profile.CharacterPool[0]
		if topChamp.GamesPlayed >= 3 {
			tendencies = append(tendencies,
				fmt.Sprintf("Signature champion: %s (%.0f%% win rate, %d games)",
					topChamp.Character, topChamp.WinRate*100, topChamp.GamesPlayed))
		}
	}

	// High pick rate champions
	for _, champ := range profile.CharacterPool {
		if champ.PickRate > 0.5 && champ.GamesPlayed >= 3 {
			tendencies = append(tendencies,
				fmt.Sprintf("%.0f%% pick rate on %s", champ.PickRate*100, champ.Character))
			break // Only show top one
		}
	}

	// KDA-based tendency
	if profile.KDA > 4.0 {
		tendencies = append(tendencies,
			fmt.Sprintf("Exceptional performer (%.2f KDA)", profile.KDA))
	} else if profile.KDA > 3.0 {
		tendencies = append(tendencies,
			fmt.Sprintf("Strong performer (%.2f KDA)", profile.KDA))
	} else if profile.KDA < 2.0 {
		tendencies = append(tendencies,
			fmt.Sprintf("Vulnerable player (%.2f KDA)", profile.KDA))
	}

	// Multikill tendency
	if profile.MultikillStats != nil {
		avgMultikills := float64(profile.MultikillStats.TotalMultikills) / float64(max(profile.GamesPlayed, 1))
		if avgMultikills > 1.5 {
			tendencies = append(tendencies,
				fmt.Sprintf("Teamfight carry (%.1f multikills/game)", avgMultikills))
		}
		if profile.MultikillStats.PentaKills > 0 {
			tendencies = append(tendencies,
				fmt.Sprintf("Penta kill threat (%d pentas)", profile.MultikillStats.PentaKills))
		}
	}

	// Objective focus
	if profile.ObjectiveFocus != nil {
		if profile.ObjectiveFocus.ObjectiveFocusType == "split-pusher" {
			tendencies = append(tendencies,
				fmt.Sprintf("Split-pusher (%.1f towers/game)", profile.ObjectiveFocus.TowersPerGame))
		} else if profile.ObjectiveFocus.DragonsPerGame > 1.5 {
			tendencies = append(tendencies,
				fmt.Sprintf("Dragon focused (%.1f dragons/game)", profile.ObjectiveFocus.DragonsPerGame))
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
				fmt.Sprintf("Strong duo with %s", partnerName))
		}
	}

	// Gold efficiency (if available)
	if profile.GoldPerMin > 0 {
		if profile.GoldPerMin > 450 {
			tendencies = append(tendencies,
				fmt.Sprintf("High economy (%.0f gold/min)", profile.GoldPerMin))
		} else if profile.GoldPerMin < 350 && profile.Role != "Support" {
			tendencies = append(tendencies,
				fmt.Sprintf("Low economy (%.0f gold/min)", profile.GoldPerMin))
		}
	}

	return tendencies
}

// determineLoLRole uses the comprehensive champion database to determine player role
func determineLoLRole(characterPool []CharacterStats) string {
	roleCounts := make(map[string]int)

	for _, char := range characterPool {
		role := GetChampionRole(char.Character)
		if role != "" {
			roleCounts[role] += char.GamesPlayed
		}
	}

	maxRole := ""
	maxCount := 0
	for role, count := range roleCounts {
		if count > maxCount {
			maxCount = count
			maxRole = role
		}
	}

	return maxRole
}
