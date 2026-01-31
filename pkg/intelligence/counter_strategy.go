package intelligence

import (
	"fmt"
	"sort"

	"scout9/pkg/grid"
)

// CounterStrategyEngine generates "How to Win" recommendations
type CounterStrategyEngine struct {
	matchupAnalyzer  *MatchupAnalyzer
	timingAnalyzer   *TimingAnalyzerEngine
	siteAnalyzer     *SiteAnalyzerEngine
	economyAnalyzer  *EconomyAnalyzer // Task 9.1: Added economy analyzer
}

// NewCounterStrategyEngine creates a new counter-strategy engine
func NewCounterStrategyEngine() *CounterStrategyEngine {
	return &CounterStrategyEngine{
		matchupAnalyzer:  NewMatchupAnalyzer(),
		timingAnalyzer:   NewTimingAnalyzer(),
		siteAnalyzer:     NewSiteAnalyzer(),
		economyAnalyzer:  NewEconomyAnalyzer(), // Task 9.1: Initialize economy analyzer
	}
}

// GenerateCounterStrategy creates counter-strategy recommendations
func (e *CounterStrategyEngine) GenerateCounterStrategy(
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	compositions *CompositionAnalysis,
) *CounterStrategy {
	strategy := &CounterStrategy{
		TeamID:               teamAnalysis.TeamID,
		TeamName:             teamAnalysis.TeamName,
		Weaknesses:           make([]WeaknessTarget, 0),
		DraftRecommendations: make([]DraftRecommendation, 0),
		InGameStrategies:     make([]Strategy, 0),
		TargetPlayers:        make([]PlayerTarget, 0),
	}

	// Generate based on game type
	if teamAnalysis.Title == "lol" {
		e.generateLoLCounterStrategy(strategy, teamAnalysis, playerProfiles, compositions)
	} else {
		e.generateVALCounterStrategy(strategy, teamAnalysis, playerProfiles, compositions)
	}

	// Calculate confidence score
	strategy.ConfidenceScore = e.calculateConfidence(teamAnalysis, strategy)

	// Generate win condition statement
	strategy.WinCondition = e.generateWinCondition(strategy, teamAnalysis)

	return strategy
}

func (e *CounterStrategyEngine) generateLoLCounterStrategy(
	strategy *CounterStrategy,
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	compositions *CompositionAnalysis,
) {
	m := teamAnalysis.LoLMetrics
	if m == nil {
		return
	}

	// Identify weaknesses from team analysis
	for _, weakness := range teamAnalysis.Weaknesses {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       weakness.Title,
			Description: weakness.Description,
			Evidence:    fmt.Sprintf("%.1f%% (n=%d)", weakness.Value, weakness.SampleSize),
			Impact:      calculateWeaknessImpact(weakness),
		})
	}

	// Early game weaknesses
	if m.FirstBloodRate < 0.4 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Passive Early Game",
			Description: "Team rarely secures first blood, vulnerable to early aggression",
			Evidence:    fmt.Sprintf("%.0f%% first blood rate", m.FirstBloodRate*100),
			Impact:      70,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Early Aggression",
			Description: "Play aggressively in lanes and invade jungle early",
			Timing:      "0-10 minutes",
			Evidence:    "Low first blood rate indicates passive early game",
		})
	}

	// Dragon control weakness
	if m.FirstDragonRate < 0.4 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Poor Dragon Control",
			Description: "Team struggles to secure first dragon",
			Evidence:    fmt.Sprintf("%.0f%% first dragon rate", m.FirstDragonRate*100),
			Impact:      65,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Dragon Priority",
			Description: "Set up vision and contest every dragon spawn",
			Timing:      "5-25 minutes",
			Evidence:    "Low dragon control rate",
		})
	}

	// Late game weakness (if they prefer early)
	if m.AvgGameDuration < 28 && teamAnalysis.WinRate > 0.5 {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Scale to Late Game",
			Description: "Draft scaling compositions and avoid early fights",
			Timing:      "Draft and 0-20 minutes",
			Evidence:    "Team wins quickly, may struggle in extended games",
		})
	}

	// Generate draft recommendations from player weaknesses
	for _, player := range playerProfiles {
		for _, weakness := range player.Weaknesses {
			if len(weakness.Title) > 7 && weakness.Title[:7] == "Weak on" {
				// Extract champion name
				champName := weakness.Title[8:]
				strategy.DraftRecommendations = append(strategy.DraftRecommendations, DraftRecommendation{
					Type:      "target",
					Character: champName,
					Reason:    fmt.Sprintf("%s has %.0f%% win rate on %s - force this pick", player.Nickname, weakness.Value, champName),
					Priority:  3,
				})
			}
		}

		// Ban signature picks with high win rates
		for _, char := range player.CharacterPool {
			if char.GamesPlayed >= 3 && char.WinRate > 0.65 {
				strategy.DraftRecommendations = append(strategy.DraftRecommendations, DraftRecommendation{
					Type:      "ban",
					Character: char.Character,
					Reason:    fmt.Sprintf("%s's %s has %.0f%% win rate", player.Nickname, char.Character, char.WinRate*100),
					Priority:  1,
				})
			}
		}

		// NEW: Multikill-based targeting - identify players who can't carry teamfights
		if player.MultikillStats != nil && player.GamesPlayed >= 3 {
			avgMultikills := float64(player.MultikillStats.TotalMultikills) / float64(player.GamesPlayed)
			
			// Player with low multikills = can't carry teamfights
			if avgMultikills < 0.5 && (player.Role == "ADC" || player.Role == "Mid") {
				strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
					Title:       fmt.Sprintf("Force teamfights - %s can't carry", player.Nickname),
					Description: fmt.Sprintf("%s averages only %.1f multikills/game - they struggle to carry teamfights", player.Nickname, avgMultikills),
					Timing:      "Mid-late game teamfights",
					Evidence:    fmt.Sprintf("%.1f multikills/game (low for %s role)", avgMultikills, player.Role),
				})
			}
			
			// Player with high multikills = must shut down
			if avgMultikills > 2.0 {
				strategy.TargetPlayers = append(strategy.TargetPlayers, PlayerTarget{
					PlayerName: player.Nickname,
					Role:       player.Role,
					Reason:     fmt.Sprintf("TEAMFIGHT CARRY - %.1f multikills/game, %d penta kills. Must shut down early!", avgMultikills, player.MultikillStats.PentaKills),
					Priority:   1, // Highest priority
				})
			}
		}

		// NEW: Synergy-based targeting - identify isolated players
		if len(player.SynergyPartners) > 0 {
			topSynergy := player.SynergyPartners[0]
			
			// Low synergy score = isolated player, easier to target
			if topSynergy.SynergyScore < 0.25 && player.GamesPlayed >= 3 {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Isolated Player: %s", player.Nickname),
					Description: fmt.Sprintf("%s has low team coordination (%.0f%% synergy) - easier to pick off", player.Nickname, topSynergy.SynergyScore*100),
					Evidence:    fmt.Sprintf("Top synergy partner only %.0f%% of assists", topSynergy.SynergyScore*100),
					Impact:      60,
				})
			}
		}

		// NEW: Assist ratio analysis - identify playmakers to shut down
		if player.AssistRatio > 1.5 && player.GamesPlayed >= 3 {
			// High assist ratio = playmaker, shutting them down disrupts team
			strategy.TargetPlayers = append(strategy.TargetPlayers, PlayerTarget{
				PlayerName: player.Nickname,
				Role:       player.Role,
				Reason:     fmt.Sprintf("PLAYMAKER - assist ratio %.2f. Shutting down %s disrupts team coordination", player.AssistRatio, player.Nickname),
				Priority:   2,
			})
		}
	}

	// NEW: Team synergy analysis - find duo to break
	e.analyzeTeamSynergy(strategy, playerProfiles)

	// Target weak players
	for _, player := range playerProfiles {
		if player.ThreatLevel <= 4 || len(player.Weaknesses) > 0 {
			// Check if already added
			alreadyAdded := false
			for _, target := range strategy.TargetPlayers {
				if target.PlayerName == player.Nickname {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				strategy.TargetPlayers = append(strategy.TargetPlayers, PlayerTarget{
					PlayerName: player.Nickname,
					Role:       player.Role,
					Reason:     player.ThreatReason,
					Priority:   10 - player.ThreatLevel,
				})
			}
		}
	}

	// Sort recommendations by priority
	sort.Slice(strategy.DraftRecommendations, func(i, j int) bool {
		return strategy.DraftRecommendations[i].Priority < strategy.DraftRecommendations[j].Priority
	})

	// Sort target players by priority
	sort.Slice(strategy.TargetPlayers, func(i, j int) bool {
		return strategy.TargetPlayers[i].Priority < strategy.TargetPlayers[j].Priority
	})

	// Limit to top recommendations
	if len(strategy.DraftRecommendations) > 5 {
		strategy.DraftRecommendations = strategy.DraftRecommendations[:5]
	}
	if len(strategy.TargetPlayers) > 3 {
		strategy.TargetPlayers = strategy.TargetPlayers[:3]
	}
}

func (e *CounterStrategyEngine) generateVALCounterStrategy(
	strategy *CounterStrategy,
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	compositions *CompositionAnalysis,
) {
	m := teamAnalysis.VALMetrics
	if m == nil {
		return
	}

	// Identify weaknesses from team analysis
	for _, weakness := range teamAnalysis.Weaknesses {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       weakness.Title,
			Description: weakness.Description,
			Evidence:    fmt.Sprintf("%.1f%% (n=%d)", weakness.Value, weakness.SampleSize),
			Impact:      calculateWeaknessImpact(weakness),
		})
	}

	// Attack side weakness
	if m.AttackWinRate < 0.45 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Weak Attack Execution",
			Description: "Team struggles on attack side",
			Evidence:    fmt.Sprintf("%.0f%% attack round win rate", m.AttackWinRate*100),
			Impact:      75,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Force Attack Rounds",
			Description: "Play aggressive defense to build economy advantage, then hold on attack",
			Timing:      "Throughout match",
			Evidence:    "Low attack win rate",
		})
	}

	// Defense side weakness
	if m.DefenseWinRate < 0.45 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Weak Defense Setup",
			Description: "Team struggles on defense side",
			Evidence:    fmt.Sprintf("%.0f%% defense round win rate", m.DefenseWinRate*100),
			Impact:      75,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Aggressive Attack Executes",
			Description: "Execute quickly on attack to exploit weak defensive setups",
			Timing:      "Attack rounds",
			Evidence:    "Low defense win rate",
		})
	}

	// Pistol round weakness
	if m.PistolWinRate < 0.4 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Poor Pistol Rounds",
			Description: "Team loses pistol rounds frequently",
			Evidence:    fmt.Sprintf("%.0f%% pistol win rate", m.PistolWinRate*100),
			Impact:      60,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Pistol Round Focus",
			Description: "Prioritize pistol round preparation and execution",
			Timing:      "Rounds 1 and 13",
			Evidence:    "Low pistol win rate gives economy advantage",
		})
	}

	// First blood vulnerability
	if m.FirstDeathRate > 0.55 {
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       "Opening Duel Vulnerability",
			Description: "Team often loses first engagement",
			Evidence:    fmt.Sprintf("%.0f%% first death rate", m.FirstDeathRate*100),
			Impact:      65,
		})

		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Aggressive Opening Duels",
			Description: "Take aggressive early peeks to secure first blood",
			Timing:      "Round start",
			Evidence:    "High first death rate",
		})
	}

	// NEW: Economy analysis - eco/force buy patterns
	if m.EcoRoundWinRate > 0 || m.ForceBuyWinRate > 0 {
		if m.EcoRoundWinRate < 0.15 {
			strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
				Title:       "Poor Eco Round Execution",
				Description: "Team rarely wins eco rounds - predictable save patterns",
				Evidence:    fmt.Sprintf("%.0f%% eco round win rate", m.EcoRoundWinRate*100),
				Impact:      50,
			})
		}
		
		if m.ForceBuyWinRate > 0.45 {
			strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
				Title:       "Respect Force Buys",
				Description: fmt.Sprintf("Opponent has %.0f%% force buy win rate - don't underestimate their eco rounds", m.ForceBuyWinRate*100),
				Timing:      "Post-loss rounds",
				Evidence:    "High force buy success rate",
			})
		} else if m.ForceBuyWinRate < 0.25 {
			strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
				Title:       "Punish Force Buys",
				Description: fmt.Sprintf("Opponent only has %.0f%% force buy win rate - play aggressive on their eco rounds", m.ForceBuyWinRate*100),
				Timing:      "Post-loss rounds",
				Evidence:    "Low force buy success rate",
			})
		}
	}

	// Map-specific strategies
	for _, mapEntry := range m.MapPool {
		if mapEntry.Strength == "weak" && mapEntry.GamesPlayed >= 3 {
			strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
				Title:       fmt.Sprintf("Force %s", mapEntry.MapName),
				Description: fmt.Sprintf("Pick %s in map veto - opponent has %.0f%% win rate", mapEntry.MapName, mapEntry.WinRate*100),
				Timing:      "Map veto",
				Evidence:    fmt.Sprintf("%.0f%% win rate on %s", mapEntry.WinRate*100, mapEntry.MapName),
			})
		}
	}

	// NEW: Weapon-based strategies
	for _, player := range playerProfiles {
		if len(player.WeaponStats) > 0 {
			topWeapon := player.WeaponStats[0]
			
			// Operator dependency - HIGH VALUE INSIGHT
			if topWeapon.WeaponName == "Operator" && topWeapon.KillShare > 0.35 {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Operator Dependency: %s", player.Nickname),
					Description: fmt.Sprintf("%s gets %.0f%% of kills with Operator - deny Op economy!", player.Nickname, topWeapon.KillShare*100),
					Evidence:    fmt.Sprintf("%.0f%% kill share with Operator (%d kills)", topWeapon.KillShare*100, topWeapon.Kills),
					Impact:      80,
				})
				
				strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
					Title:       fmt.Sprintf("Deny Operator from %s", player.Nickname),
					Description: "Force eco rounds, use utility to flush Op angles, close distance quickly",
					Timing:      "Full buy rounds",
					Evidence:    fmt.Sprintf("%s is %.0f%% Operator dependent", player.Nickname, topWeapon.KillShare*100),
				})
			}
			
			// Rifle-only player - can exploit at range
			if (topWeapon.WeaponName == "Vandal" || topWeapon.WeaponName == "Phantom") && topWeapon.KillShare > 0.6 {
				// Check if they have low Op kills
				hasOp := false
				for _, ws := range player.WeaponStats {
					if ws.WeaponName == "Operator" && ws.KillShare > 0.1 {
						hasOp = true
						break
					}
				}
				if !hasOp {
					strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
						Title:       fmt.Sprintf("Use long angles vs %s", player.Nickname),
						Description: fmt.Sprintf("%s doesn't use Operator (%.0f%% rifle kills) - take long-range duels", player.Nickname, topWeapon.KillShare*100),
						Timing:      "Defense setup",
						Evidence:    "No Operator usage in weapon stats",
					})
				}
			}
		}

		// NEW: Synergy-based targeting for VALORANT
		if len(player.SynergyPartners) > 0 {
			topSynergy := player.SynergyPartners[0]
			
			// Low synergy = isolated player
			if topSynergy.SynergyScore < 0.25 && player.GamesPlayed >= 3 {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Isolated Player: %s", player.Nickname),
					Description: fmt.Sprintf("%s plays isolated (%.0f%% synergy) - easier to pick off in clutches", player.Nickname, topSynergy.SynergyScore*100),
					Evidence:    fmt.Sprintf("Low assist connection with teammates"),
					Impact:      55,
				})
			}
			
			// High synergy duo = break them up
			if topSynergy.SynergyScore > 0.5 {
				partnerName := topSynergy.PlayerName
				if partnerName == "" {
					partnerName = topSynergy.PlayerID
				}
				strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
					Title:       fmt.Sprintf("Break %s + %s duo", player.Nickname, partnerName),
					Description: fmt.Sprintf("Strong synergy (%.0f%%) - use utility to separate them", topSynergy.SynergyScore*100),
					Timing:      "Site executes and retakes",
					Evidence:    fmt.Sprintf("%.0f%% of %s's assists come from %s", topSynergy.SynergyScore*100, player.Nickname, partnerName),
				})
			}
		}
	}

	// NEW: Team synergy analysis
	e.analyzeTeamSynergy(strategy, playerProfiles)

	// Target weak players
	for _, player := range playerProfiles {
		if player.ThreatLevel <= 4 || len(player.Weaknesses) > 0 {
			reason := player.ThreatReason
			if len(player.Weaknesses) > 0 {
				reason = player.Weaknesses[0].Description
			}
			
			// Check if already added
			alreadyAdded := false
			for _, target := range strategy.TargetPlayers {
				if target.PlayerName == player.Nickname {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				strategy.TargetPlayers = append(strategy.TargetPlayers, PlayerTarget{
					PlayerName: player.Nickname,
					Role:       player.Role,
					Reason:     reason,
					Priority:   10 - player.ThreatLevel,
				})
			}
		}
	}

	// Agent bans based on player performance
	for _, player := range playerProfiles {
		for _, agent := range player.CharacterPool {
			if agent.GamesPlayed >= 3 && agent.WinRate > 0.65 {
				strategy.DraftRecommendations = append(strategy.DraftRecommendations, DraftRecommendation{
					Type:      "ban",
					Character: agent.Character,
					Reason:    fmt.Sprintf("%s's %s has %.0f%% win rate", player.Nickname, agent.Character, agent.WinRate*100),
					Priority:  1,
				})
			}
		}
	}

	// Sort and limit
	sort.Slice(strategy.DraftRecommendations, func(i, j int) bool {
		return strategy.DraftRecommendations[i].Priority < strategy.DraftRecommendations[j].Priority
	})
	sort.Slice(strategy.TargetPlayers, func(i, j int) bool {
		return strategy.TargetPlayers[i].Priority < strategy.TargetPlayers[j].Priority
	})

	if len(strategy.DraftRecommendations) > 5 {
		strategy.DraftRecommendations = strategy.DraftRecommendations[:5]
	}
	if len(strategy.TargetPlayers) > 3 {
		strategy.TargetPlayers = strategy.TargetPlayers[:3]
	}
}

func calculateWeaknessImpact(weakness Insight) float64 {
	// Base impact on how far from average the value is
	impact := 50.0

	// Adjust based on sample size
	if weakness.SampleSize >= 10 {
		impact += 20
	} else if weakness.SampleSize >= 5 {
		impact += 10
	}

	// Adjust based on severity (assuming lower is worse for most metrics)
	if weakness.Value < 30 {
		impact += 20
	} else if weakness.Value < 40 {
		impact += 10
	}

	if impact > 100 {
		impact = 100
	}

	return impact
}

func (e *CounterStrategyEngine) calculateConfidence(teamAnalysis *TeamAnalysis, strategy *CounterStrategy) float64 {
	confidence := 50.0

	// More games analyzed = higher confidence (Task 8.1)
	if teamAnalysis.GamesAnalyzed >= 15 {
		confidence += 25 // Boost by 25 points when ≥15 games
	} else if teamAnalysis.GamesAnalyzed >= 10 {
		confidence += 15 // Boost by 15 points when ≥10 games
	} else if teamAnalysis.GamesAnalyzed >= 5 {
		confidence += 5
	}

	// More weaknesses identified = higher confidence
	if len(strategy.Weaknesses) >= 3 {
		confidence += 15
	} else if len(strategy.Weaknesses) >= 1 {
		confidence += 5
	}

	// Draft recommendations increase confidence
	if len(strategy.DraftRecommendations) >= 3 {
		confidence += 10
	}

	// Target players increase confidence
	if len(strategy.TargetPlayers) >= 2 {
		confidence += 5
	}

	if confidence > 100 {
		confidence = 100
	}

	return confidence
}

// calculateInsightConfidence calculates confidence for individual insights (Task 8.1)
func (e *CounterStrategyEngine) calculateInsightConfidence(sampleSize int, baseConfidence float64) float64 {
	confidence := baseConfidence

	if sampleSize >= 15 {
		confidence += 25
	} else if sampleSize >= 10 {
		confidence += 15
	} else if sampleSize >= 5 {
		confidence += 5
	} else if sampleSize < 3 {
		confidence -= 10 // Reduce confidence for very small samples
	}

	if confidence > 100 {
		confidence = 100
	}
	if confidence < 0 {
		confidence = 0
	}

	return confidence
}

// generateSampleSizeWarning generates warning for low sample sizes (Task 8.2)
func (e *CounterStrategyEngine) generateSampleSizeWarning(gamesAnalyzed int) string {
	if gamesAnalyzed < 5 {
		return fmt.Sprintf("Warning: Only %d games analyzed - insights may be less reliable due to low sample size.", gamesAnalyzed)
	}
	return ""
}

func (e *CounterStrategyEngine) generateWinCondition(strategy *CounterStrategy, teamAnalysis *TeamAnalysis) string {
	if len(strategy.Weaknesses) == 0 {
		return fmt.Sprintf("To beat %s, maintain consistent execution and capitalize on any mistakes.", teamAnalysis.TeamName)
	}

	// Find highest impact weakness
	var topWeakness WeaknessTarget
	for _, w := range strategy.Weaknesses {
		if w.Impact > topWeakness.Impact {
			topWeakness = w
		}
	}

	// Generate based on game type
	if teamAnalysis.Title == "lol" {
		return e.generateLoLWinCondition(strategy, teamAnalysis, topWeakness)
	}
	return e.generateVALWinCondition(strategy, teamAnalysis, topWeakness)
}

func (e *CounterStrategyEngine) generateLoLWinCondition(strategy *CounterStrategy, teamAnalysis *TeamAnalysis, topWeakness WeaknessTarget) string {
	m := teamAnalysis.LoLMetrics
	if m == nil {
		return fmt.Sprintf("To beat %s, exploit their %s.", teamAnalysis.TeamName, topWeakness.Title)
	}

	// Build win condition based on weaknesses
	if m.FirstBloodRate < 0.4 && m.FirstDragonRate < 0.4 {
		return fmt.Sprintf("To beat %s, dominate the early game with aggressive plays and secure dragon control. Their passive early game (%.0f%% first blood) leaves them vulnerable to snowballing.",
			teamAnalysis.TeamName, m.FirstBloodRate*100)
	}

	if m.AvgGameDuration < 28 {
		return fmt.Sprintf("To beat %s, draft scaling compositions and survive the early game. They win fast (avg %.0f min) but may struggle in extended games.",
			teamAnalysis.TeamName, m.AvgGameDuration)
	}

	if len(strategy.TargetPlayers) > 0 {
		return fmt.Sprintf("To beat %s, focus pressure on %s (%s) while exploiting their %s.",
			teamAnalysis.TeamName, strategy.TargetPlayers[0].PlayerName, strategy.TargetPlayers[0].Role, topWeakness.Title)
	}

	return fmt.Sprintf("To beat %s, exploit their %s and maintain objective control.",
		teamAnalysis.TeamName, topWeakness.Title)
}

func (e *CounterStrategyEngine) generateVALWinCondition(strategy *CounterStrategy, teamAnalysis *TeamAnalysis, topWeakness WeaknessTarget) string {
	m := teamAnalysis.VALMetrics
	if m == nil {
		return fmt.Sprintf("To beat %s, exploit their %s.", teamAnalysis.TeamName, topWeakness.Title)
	}

	// Build win condition based on weaknesses
	if m.AttackWinRate < 0.45 {
		return fmt.Sprintf("To beat %s, play solid defense and force them into attack rounds. Their %.0f%% attack win rate is exploitable.",
			teamAnalysis.TeamName, m.AttackWinRate*100)
	}

	if m.DefenseWinRate < 0.45 {
		return fmt.Sprintf("To beat %s, execute quickly on attack to exploit their weak defensive setups (%.0f%% defense win rate).",
			teamAnalysis.TeamName, m.DefenseWinRate*100)
	}

	if m.PistolWinRate < 0.4 {
		return fmt.Sprintf("To beat %s, win pistol rounds to build economy advantage. Their %.0f%% pistol win rate gives you a head start.",
			teamAnalysis.TeamName, m.PistolWinRate*100)
	}

	// Map-based win condition
	for _, mapEntry := range m.MapPool {
		if mapEntry.Strength == "weak" && mapEntry.GamesPlayed >= 3 {
			return fmt.Sprintf("To beat %s, force %s in map veto. They have only %.0f%% win rate on this map.",
				teamAnalysis.TeamName, mapEntry.MapName, mapEntry.WinRate*100)
		}
	}

	return fmt.Sprintf("To beat %s, exploit their %s and maintain round-by-round discipline.",
		teamAnalysis.TeamName, topWeakness.Title)
}

// =============================================================================
// HACKATHON-WINNING ENHANCED METHODS
// =============================================================================

// GenerateEnhancedCounterStrategy creates counter-strategy with timing, matchup, and site analysis
// This is the HACKATHON-WINNING version that produces specific, actionable insights
func (e *CounterStrategyEngine) GenerateEnhancedCounterStrategy(
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	compositions *CompositionAnalysis,
	seriesStates []*grid.SeriesState,
	lolEvents map[string]*grid.LoLEventData,
	valEvents map[string]*grid.VALEventData,
) *CounterStrategy {
	// Start with base counter-strategy
	strategy := e.GenerateCounterStrategy(teamAnalysis, playerProfiles, compositions)

	// Track insights for enhanced win condition
	var classInsights []ClassMatchupInsight
	var economyInsights []EconomyInsight
	var siteInsights []StrategyInsight

	// Enhance based on game type
	if teamAnalysis.Title == "lol" {
		classInsights = e.enhanceLoLStrategyWithInsights(strategy, teamAnalysis, playerProfiles, seriesStates, lolEvents)
	} else {
		economyInsights, siteInsights = e.enhanceVALStrategyWithInsights(strategy, teamAnalysis, playerProfiles, seriesStates, valEvents)
	}

	// Generate enhanced win condition (Task 7)
	strategy.WinCondition = e.generateEnhancedWinCondition(strategy, teamAnalysis, classInsights, economyInsights, siteInsights)

	// Add sample size warning if needed (Task 8.2)
	warning := e.generateSampleSizeWarning(teamAnalysis.GamesAnalyzed)
	if warning != "" {
		strategy.Warnings = append(strategy.Warnings, warning)
	}

	return strategy
}

// enhanceLoLStrategyWithInsights adds timing and matchup analysis to LoL counter-strategy
// Returns class insights for win condition generation
func (e *CounterStrategyEngine) enhanceLoLStrategyWithInsights(
	strategy *CounterStrategy,
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) []ClassMatchupInsight {
	var allClassInsights []ClassMatchupInsight

	// Analyze jungle pathing
	junglePathing := e.timingAnalyzer.AnalyzeJunglePathing(teamAnalysis.TeamID, seriesStates, events)

	// Analyze objective timings
	objectiveTimings := e.timingAnalyzer.AnalyzeObjectiveTimings(teamAnalysis.TeamID, seriesStates, events)

	// Generate timing-based counter-strategies
	timingStrategies := e.timingAnalyzer.GenerateTimingCounterStrategies(junglePathing, objectiveTimings)
	for _, ts := range timingStrategies {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       ts.Strategy,
			Description: ts.Reason,
			Timing:      ts.Timing,
			Evidence:    ts.Reason,
		})
	}

	// Analyze player matchups for draft recommendations (Task 9.3)
	for _, player := range playerProfiles {
		classPerf := e.matchupAnalyzer.AnalyzeLoLPlayerClassPerformance(
			player.PlayerID,
			player.Nickname,
			seriesStates,
		)

		// Generate class matchup insights (Task 1.2)
		classInsights := e.matchupAnalyzer.GenerateClassMatchupInsights(player.Nickname, player.Role, classPerf)
		allClassInsights = append(allClassInsights, classInsights...)

		// Generate specific draft recommendations (Task 1.3)
		specificDraftRecs := e.matchupAnalyzer.GenerateSpecificDraftRecommendations(player.Nickname, player.Role, classPerf)
		for _, rec := range specificDraftRecs {
			strategy.DraftRecommendations = append(strategy.DraftRecommendations, DraftRecommendation{
				Type:      rec.Type,
				Character: rec.Champions[0], // Use first champion
				Reason:    rec.Reason,
				Priority:  rec.Priority,
			})
		}

		// Generate matchup-based draft recommendations
		draftRecs := e.matchupAnalyzer.GenerateDraftRecommendations(player.Nickname, player.Role, classPerf)
		for _, rec := range draftRecs {
			strategy.DraftRecommendations = append(strategy.DraftRecommendations, DraftRecommendation{
				Type:      "pick",
				Character: rec.Character,
				Reason:    rec.Text,
				Priority:  rec.Priority,
			})
		}

		// Add matchup insights to target players
		matchupInsights := e.matchupAnalyzer.GenerateMatchupInsights(player.Nickname, classPerf)
		for _, insight := range matchupInsights {
			found := false
			for i, target := range strategy.TargetPlayers {
				if target.PlayerName == player.Nickname {
					strategy.TargetPlayers[i].Reason = fmt.Sprintf("%s. %s", target.Reason, insight.Text)
					found = true
					break
				}
			}
			if !found && insight.Value > 1.5 {
				strategy.TargetPlayers = append(strategy.TargetPlayers, PlayerTarget{
					PlayerName: player.Nickname,
					Role:       player.Role,
					Reason:     insight.Text,
					Priority:   3,
				})
			}
		}
	}

	// Update win condition with timing data
	if junglePathing.PreSixBotRate > 0.6 {
		strategy.WinCondition = fmt.Sprintf("%s Their jungler paths to bot %.0f%% pre-10 mins - counter-jungle top side.",
			strategy.WinCondition, junglePathing.PreSixBotRate*100)
	} else if junglePathing.PreSixTopRate > 0.6 {
		strategy.WinCondition = fmt.Sprintf("%s Their jungler paths to top %.0f%% pre-10 mins - counter-jungle bot side.",
			strategy.WinCondition, junglePathing.PreSixTopRate*100)
	}

	// Re-sort and limit
	sort.Slice(strategy.DraftRecommendations, func(i, j int) bool {
		return strategy.DraftRecommendations[i].Priority < strategy.DraftRecommendations[j].Priority
	})
	if len(strategy.DraftRecommendations) > 7 {
		strategy.DraftRecommendations = strategy.DraftRecommendations[:7]
	}

	return allClassInsights
}

// enhanceVALStrategyWithInsights adds site and economy analysis to VALORANT counter-strategy
// Returns economy and site insights for win condition generation
func (e *CounterStrategyEngine) enhanceVALStrategyWithInsights(
	strategy *CounterStrategy,
	teamAnalysis *TeamAnalysis,
	playerProfiles []*PlayerProfile,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.VALEventData,
) ([]EconomyInsight, []StrategyInsight) {
	// Analyze economy rounds (Task 9.4)
	economyAnalysis := e.economyAnalyzer.AnalyzeEconomyRounds(teamAnalysis.TeamID, teamAnalysis.TeamName, seriesStates)
	economyInsights := e.economyAnalyzer.GenerateEconomyInsights(economyAnalysis)

	// Add economy insights to strategy
	for _, insight := range economyInsights {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       insight.Type,
			Description: insight.Text,
			Timing:      "Economy rounds",
			Evidence:    fmt.Sprintf("%.0f%% (n=%d)", insight.Value*100, insight.SampleSize),
		})
	}

	// Analyze site patterns
	siteAnalysis := e.siteAnalyzer.AnalyzeSitePatterns(teamAnalysis.TeamID, seriesStates, events)

	// Identify site weaknesses (Task 4.2)
	siteWeaknesses := IdentifySiteWeaknesses(siteAnalysis)
	siteInsights := GenerateSiteWeaknessInsights(siteWeaknesses)

	// Add site weakness insights
	for _, insight := range siteInsights {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       "Site Weakness",
			Description: insight.Text,
			Timing:      "Attack rounds",
			Evidence:    fmt.Sprintf("%.0f%% win rate (n=%d)", insight.Value*100, insight.SampleSize),
		})
	}

	// Generate site recommendations
	siteRecommendations := GenerateSiteRecommendations(siteAnalysis)
	for _, rec := range siteRecommendations {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       rec.Strategy,
			Description: rec.Reason,
			Timing:      rec.Timing,
			Evidence:    rec.Reason,
		})
	}

	// Analyze attack patterns
	attackPatterns := e.siteAnalyzer.AnalyzeAttackPatterns(teamAnalysis.TeamID, seriesStates, events)

	// Generate site-based counter-strategies
	siteStrategies := e.siteAnalyzer.GenerateSiteCounterStrategies(siteAnalysis)
	for _, ss := range siteStrategies {
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       ss.Strategy,
			Description: ss.Reason,
			Timing:      ss.Timing,
			Evidence:    ss.Reason,
		})
	}

	// Add site insights to weaknesses
	for mapName, analysis := range siteAnalysis {
		for siteName, stats := range analysis.Sites {
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate < 0.4 {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Weak %s-Site Defense on %s", siteName, mapName),
					Description: fmt.Sprintf("Only %.0f%% defense win rate on %s-Site", stats.DefenseWinRate*100, siteName),
					Evidence:    fmt.Sprintf("%.0f%% win rate (n=%d)", stats.DefenseWinRate*100, stats.DefenseAttempts),
					Impact:      75,
				})
			}

			if stats.AttackAttempts >= 3 && stats.AttackWinRate < 0.4 {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Weak %s-Site Attack on %s", siteName, mapName),
					Description: fmt.Sprintf("Only %.0f%% attack success rate on %s-Site", stats.AttackWinRate*100, siteName),
					Evidence:    fmt.Sprintf("%.0f%% win rate (n=%d)", stats.AttackWinRate*100, stats.AttackAttempts),
					Impact:      60,
				})
			}
		}
	}

	// Add attack pattern insights
	for _, pattern := range attackPatterns {
		if pattern.Frequency > 0.4 {
			strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
				Title:       fmt.Sprintf("Counter %s pattern", pattern.Description),
				Description: fmt.Sprintf("Opponent uses %s %.0f%% of the time - prepare counter-setup", pattern.Description, pattern.Frequency*100),
				Timing:      "Defense rounds",
				Evidence:    fmt.Sprintf("%.0f%% frequency, %.0f%% success rate", pattern.Frequency*100, pattern.SuccessRate*100),
			})
		}
	}

	// Update win condition with site data
	for mapName, analysis := range siteAnalysis {
		for siteName, stats := range analysis.Sites {
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate < 0.4 {
				strategy.WinCondition = fmt.Sprintf("%s Attack %s-Site on %s (%.0f%% defense win rate).",
					strategy.WinCondition, siteName, mapName, stats.DefenseWinRate*100)
				break
			}
		}
	}

	return economyInsights, siteInsights
}


// =============================================================================
// ENHANCED WIN CONDITION GENERATION (Task 7)
// =============================================================================

// generateEnhancedWinCondition creates a specific, data-backed win condition statement
// Hackathon format: Include specific data points (percentages, player names), 2-3 sentences max
func (e *CounterStrategyEngine) generateEnhancedWinCondition(
	strategy *CounterStrategy,
	teamAnalysis *TeamAnalysis,
	classInsights []ClassMatchupInsight,
	economyInsights []EconomyInsight,
	siteInsights []StrategyInsight,
) string {
	var parts []string

	// Start with team name
	teamName := teamAnalysis.TeamName

	// Add class matchup insight if available (LoL)
	if len(classInsights) > 0 && teamAnalysis.Title == "lol" {
		insight := classInsights[0]
		parts = append(parts, fmt.Sprintf("Target %s's %s (%.1f KDA) with %s picks.",
			insight.PlayerName, insight.WeakClass, insight.WeakClassKDA, insight.StrongClass))
	}

	// Add economy insight if available (VALORANT)
	if len(economyInsights) > 0 && teamAnalysis.Title == "valorant" {
		for _, insight := range economyInsights {
			if insight.Impact == "HIGH" {
				parts = append(parts, insight.Text)
				break
			}
		}
	}

	// Add site insight if available (VALORANT)
	if len(siteInsights) > 0 && teamAnalysis.Title == "valorant" {
		parts = append(parts, siteInsights[0].Text)
	}

	// Add timing-based insight for LoL
	if teamAnalysis.Title == "lol" && teamAnalysis.LoLMetrics != nil {
		m := teamAnalysis.LoLMetrics
		if m.FirstBloodRate < 0.35 {
			parts = append(parts, fmt.Sprintf("%s only gets first blood %.0f%% - play aggressive early.", teamName, m.FirstBloodRate*100))
		} else if m.FirstDragonRate < 0.35 {
			parts = append(parts, fmt.Sprintf("%s only secures first dragon %.0f%% - prioritize dragon control.", teamName, m.FirstDragonRate*100))
		}
	}

	// Add target player if available
	if len(strategy.TargetPlayers) > 0 {
		target := strategy.TargetPlayers[0]
		parts = append(parts, fmt.Sprintf("Focus pressure on %s (%s).", target.PlayerName, target.Role))
	}

	// Combine parts, limit to 3 sentences
	if len(parts) == 0 {
		return fmt.Sprintf("To beat %s, maintain consistent execution and capitalize on their weaknesses.", teamName)
	}

	if len(parts) > 3 {
		parts = parts[:3]
	}

	result := fmt.Sprintf("To beat %s: ", teamName)
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}

	return result
}

// analyzeTeamSynergy identifies strong duos and isolated players across the team
func (e *CounterStrategyEngine) analyzeTeamSynergy(strategy *CounterStrategy, playerProfiles []*PlayerProfile) {
	// Build synergy matrix
	type duoSynergy struct {
		player1     string
		player2     string
		totalAssists int
		synergyScore float64
	}
	
	duos := make([]duoSynergy, 0)
	
	for _, player := range playerProfiles {
		for _, partner := range player.SynergyPartners {
			// Find partner name
			partnerName := partner.PlayerName
			if partnerName == "" {
				// Look up partner name from profiles
				for _, p := range playerProfiles {
					if p.PlayerID == partner.PlayerID {
						partnerName = p.Nickname
						break
					}
				}
			}
			if partnerName == "" {
				partnerName = partner.PlayerID
			}
			
			// Avoid duplicates (A->B and B->A)
			if player.Nickname < partnerName {
				duos = append(duos, duoSynergy{
					player1:      player.Nickname,
					player2:      partnerName,
					totalAssists: partner.AssistCount,
					synergyScore: partner.SynergyScore,
				})
			}
		}
	}
	
	// Sort by synergy score
	sort.Slice(duos, func(i, j int) bool {
		return duos[i].synergyScore > duos[j].synergyScore
	})
	
	// Add top duo as a weakness to exploit
	if len(duos) > 0 && duos[0].synergyScore > 0.4 {
		topDuo := duos[0]
		strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
			Title:       fmt.Sprintf("Strong Duo: %s + %s", topDuo.player1, topDuo.player2),
			Description: fmt.Sprintf("This duo has %.0f%% synergy - breaking them up disrupts team coordination", topDuo.synergyScore*100),
			Evidence:    fmt.Sprintf("%d combined assists between them", topDuo.totalAssists),
			Impact:      65,
		})
		
		strategy.InGameStrategies = append(strategy.InGameStrategies, Strategy{
			Title:       fmt.Sprintf("Split %s and %s", topDuo.player1, topDuo.player2),
			Description: "Use utility and positioning to prevent this duo from playing together",
			Timing:      "Throughout the game",
			Evidence:    fmt.Sprintf("%.0f%% synergy score - they rely on each other", topDuo.synergyScore*100),
		})
	}
	
	// Identify isolated players (no strong synergy with anyone)
	for _, player := range playerProfiles {
		if len(player.SynergyPartners) == 0 {
			continue
		}
		
		// Check if player has any strong synergy
		hasStrongSynergy := false
		for _, partner := range player.SynergyPartners {
			if partner.SynergyScore > 0.3 {
				hasStrongSynergy = true
				break
			}
		}
		
		if !hasStrongSynergy && player.GamesPlayed >= 3 {
			// Check if already added as weakness
			alreadyAdded := false
			for _, w := range strategy.Weaknesses {
				if w.Title == fmt.Sprintf("Isolated Player: %s", player.Nickname) {
					alreadyAdded = true
					break
				}
			}
			
			if !alreadyAdded {
				strategy.Weaknesses = append(strategy.Weaknesses, WeaknessTarget{
					Title:       fmt.Sprintf("Isolated Player: %s", player.Nickname),
					Description: fmt.Sprintf("%s has no strong synergy with teammates - plays independently", player.Nickname),
					Evidence:    "No synergy partner above 30%",
					Impact:      50,
				})
			}
		}
	}
}
