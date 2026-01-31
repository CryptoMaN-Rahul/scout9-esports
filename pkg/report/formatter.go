package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"scout9/pkg/intelligence"
)

// Formatter converts analysis data into hackathon-compliant DigestibleReport format
type Formatter struct{}

// NewFormatter creates a new report formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatDigestibleReport converts a ScoutingReport into the hackathon-compliant format
// This is THE KEY to winning - matching the exact output format judges expect
func (f *Formatter) FormatDigestibleReport(report *intelligence.ScoutingReport) *intelligence.DigestibleReport {
	digestible := &intelligence.DigestibleReport{
		TeamName:        report.OpponentTeam.Name,
		MatchesAnalyzed: report.MatchesAnalyzed,
		GeneratedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	// Generate executive summary (1 paragraph)
	digestible.ExecutiveSummary = f.generateDigestibleSummary(report)

	// Format common strategies section
	digestible.CommonStrategies = f.formatCommonStrategies(report)

	// Format player tendencies
	digestible.PlayerTendencies = f.formatPlayerTendencies(report)

	// Format recent compositions
	digestible.RecentCompositions = f.formatCompositions(report)

	// Format "How to Win" section - THE DIFFERENTIATOR
	digestible.HowToWin = f.formatHowToWin(report)

	return digestible
}

// generateDigestibleSummary creates a concise 1-paragraph summary
func (f *Formatter) generateDigestibleSummary(report *intelligence.ScoutingReport) string {
	teamName := report.OpponentTeam.Name
	winRate := 0.0
	if report.TeamStrategy != nil {
		winRate = report.TeamStrategy.WinRate
	}

	form := "stable"
	if report.TrendAnalysis != nil {
		form = report.TrendAnalysis.FormIndicator
	}

	// Build summary based on game type
	if report.Title == "lol" {
		return f.generateLoLSummary(teamName, winRate, form, report)
	}
	return f.generateVALSummary(teamName, winRate, form, report)
}

func (f *Formatter) generateLoLSummary(teamName string, winRate float64, form string, report *intelligence.ScoutingReport) string {
	m := report.TeamStrategy.LoLMetrics
	if m == nil {
		return fmt.Sprintf("%s has a %.0f%% win rate across %d matches analyzed. Current form: %s.",
			teamName, winRate*100, report.MatchesAnalyzed, form)
	}

	// Determine playstyle
	playstyle := "balanced"
	if m.EarlyGameRating > 60 {
		playstyle = "aggressive early-game"
	} else if m.AvgGameDuration > 35 {
		playstyle = "scaling late-game"
	}

	// Key strength
	keyStrength := ""
	if m.FirstDragonRate > 0.6 {
		keyStrength = fmt.Sprintf("strong dragon control (%.0f%% first dragon)", m.FirstDragonRate*100)
	} else if m.FirstBloodRate > 0.5 {
		keyStrength = fmt.Sprintf("aggressive early game (%.0f%% first blood)", m.FirstBloodRate*100)
	}

	// Key weakness
	keyWeakness := ""
	if m.FirstBloodRate < 0.4 {
		keyWeakness = fmt.Sprintf("passive early game (%.0f%% first blood)", m.FirstBloodRate*100)
	} else if m.FirstDragonRate < 0.4 {
		keyWeakness = fmt.Sprintf("poor dragon control (%.0f%% first dragon)", m.FirstDragonRate*100)
	}

	summary := fmt.Sprintf("%s is a %s team with %.0f%% win rate (%d matches). ",
		teamName, playstyle, winRate*100, report.MatchesAnalyzed)

	if keyStrength != "" {
		summary += fmt.Sprintf("Key strength: %s. ", keyStrength)
	}
	if keyWeakness != "" {
		summary += fmt.Sprintf("Exploitable weakness: %s. ", keyWeakness)
	}
	summary += fmt.Sprintf("Current form: %s.", form)

	return summary
}

func (f *Formatter) generateVALSummary(teamName string, winRate float64, form string, report *intelligence.ScoutingReport) string {
	m := report.TeamStrategy.VALMetrics
	if m == nil {
		return fmt.Sprintf("%s has a %.0f%% win rate across %d matches analyzed. Current form: %s.",
			teamName, winRate*100, report.MatchesAnalyzed, form)
	}

	// Determine side preference
	sidePreference := "balanced"
	if m.AttackWinRate > m.DefenseWinRate+0.1 {
		sidePreference = "attack-favored"
	} else if m.DefenseWinRate > m.AttackWinRate+0.1 {
		sidePreference = "defense-favored"
	}

	// Key strength
	keyStrength := ""
	if m.PistolWinRate > 0.55 {
		keyStrength = fmt.Sprintf("strong pistol rounds (%.0f%% win rate)", m.PistolWinRate*100)
	} else if m.FirstBloodRate > 0.55 {
		keyStrength = fmt.Sprintf("dominant opening duels (%.0f%% first blood)", m.FirstBloodRate*100)
	}

	// Key weakness
	keyWeakness := ""
	if m.AttackWinRate < 0.45 {
		keyWeakness = fmt.Sprintf("weak attack execution (%.0f%% attack win rate)", m.AttackWinRate*100)
	} else if m.DefenseWinRate < 0.45 {
		keyWeakness = fmt.Sprintf("vulnerable defense (%.0f%% defense win rate)", m.DefenseWinRate*100)
	}

	summary := fmt.Sprintf("%s is a %s team with %.0f%% win rate (%d matches). ",
		teamName, sidePreference, winRate*100, report.MatchesAnalyzed)

	if keyStrength != "" {
		summary += fmt.Sprintf("Key strength: %s. ", keyStrength)
	}
	if keyWeakness != "" {
		summary += fmt.Sprintf("Exploitable weakness: %s. ", keyWeakness)
	}
	summary += fmt.Sprintf("Current form: %s.", form)

	return summary
}

// formatCommonStrategies creates hackathon-format strategy insights
func (f *Formatter) formatCommonStrategies(report *intelligence.ScoutingReport) intelligence.CommonStrategiesSection {
	section := intelligence.CommonStrategiesSection{
		AttackPatterns:      make([]intelligence.StrategyInsight, 0),
		DefenseSetups:       make([]intelligence.StrategyInsight, 0),
		ObjectivePriorities: make([]intelligence.StrategyInsight, 0),
		TimingPatterns:      make([]intelligence.StrategyInsight, 0),
	}

	if report.Title == "lol" {
		f.formatLoLStrategies(&section, report)
	} else {
		f.formatVALStrategies(&section, report)
	}

	return section
}

func (f *Formatter) formatLoLStrategies(section *intelligence.CommonStrategiesSection, report *intelligence.ScoutingReport) {
	m := report.TeamStrategy.LoLMetrics
	if m == nil {
		return
	}

	// Objective priorities - hackathon format: "Prioritizes first Drake (82% contest rate)"
	if m.FirstDragonRate > 0.5 {
		section.ObjectivePriorities = append(section.ObjectivePriorities, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("Prioritizes first Drake (%.0f%% contest rate)", m.FirstDragonRate*100),
			Metric:     "first_dragon_rate",
			Value:      m.FirstDragonRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "early game",
		})
	}

	if m.HeraldControlRate > 0.5 {
		section.ObjectivePriorities = append(section.ObjectivePriorities, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("Strong Herald priority (%.0f%% control rate)", m.HeraldControlRate*100),
			Metric:     "herald_control_rate",
			Value:      m.HeraldControlRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "early game",
		})
	}

	if m.BaronControlRate > 0.5 {
		section.ObjectivePriorities = append(section.ObjectivePriorities, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("High Baron priority (%.0f%% control rate)", m.BaronControlRate*100),
			Metric:     "baron_control_rate",
			Value:      m.BaronControlRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "mid-late game",
		})
	}

	// Timing patterns - hackathon format: "4-man group for first tower push, usually in bot lane at ~13 mins"
	if m.FirstTowerAvgTime > 0 {
		section.TimingPatterns = append(section.TimingPatterns, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("Average first tower at ~%.0f mins (%.0f%% first tower rate)", m.FirstTowerAvgTime, m.FirstTowerRate*100),
			Metric:     "first_tower_timing",
			Value:      m.FirstTowerAvgTime,
			SampleSize: report.MatchesAnalyzed,
			Context:    "early game",
		})
	}

	// Game duration pattern
	if m.AvgGameDuration > 0 {
		gameStyle := "standard pace"
		if m.AvgGameDuration < 28 {
			gameStyle = "fast-paced, early game focused"
		} else if m.AvgGameDuration > 35 {
			gameStyle = "slow-paced, scaling focused"
		}
		section.TimingPatterns = append(section.TimingPatterns, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("Average game duration: %.0f mins (%s)", m.AvgGameDuration, gameStyle),
			Metric:     "avg_game_duration",
			Value:      m.AvgGameDuration,
			SampleSize: report.MatchesAnalyzed,
			Context:    "game pace",
		})
	}
}

func (f *Formatter) formatVALStrategies(section *intelligence.CommonStrategiesSection, report *intelligence.ScoutingReport) {
	m := report.TeamStrategy.VALMetrics
	if m == nil {
		return
	}

	// Attack patterns - hackathon format: "On Attack, 70% of pistol rounds are a 5-man fast-hit on B-Site (Ascent)"
	if m.AttackPistolWinRate > 0 {
		section.AttackPatterns = append(section.AttackPatterns, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("On Attack pistol rounds: %.0f%% win rate", m.AttackPistolWinRate*100),
			Metric:     "attack_pistol_win_rate",
			Value:      m.AttackPistolWinRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "pistol rounds",
		})
	}

	if m.AttackWinRate > 0 {
		attackStyle := "balanced"
		if m.AttackWinRate > 0.55 {
			attackStyle = "aggressive"
		} else if m.AttackWinRate < 0.45 {
			attackStyle = "passive"
		}
		section.AttackPatterns = append(section.AttackPatterns, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("On Attack: %.0f%% round win rate (%s style)", m.AttackWinRate*100, attackStyle),
			Metric:     "attack_win_rate",
			Value:      m.AttackWinRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "attack side",
		})
	}

	// Defense setups - hackathon format: "On Defense, they default to a 1-3-1 setup, rotating their Sentinel to mid"
	if m.DefensePistolWinRate > 0 {
		section.DefenseSetups = append(section.DefenseSetups, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("On Defense pistol rounds: %.0f%% win rate", m.DefensePistolWinRate*100),
			Metric:     "defense_pistol_win_rate",
			Value:      m.DefensePistolWinRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "pistol rounds",
		})
	}

	if m.DefenseWinRate > 0 {
		defenseStyle := "balanced"
		if m.DefenseWinRate > 0.55 {
			defenseStyle = "strong hold"
		} else if m.DefenseWinRate < 0.45 {
			defenseStyle = "vulnerable"
		}
		section.DefenseSetups = append(section.DefenseSetups, intelligence.StrategyInsight{
			Text:       fmt.Sprintf("On Defense: %.0f%% round win rate (%s)", m.DefenseWinRate*100, defenseStyle),
			Metric:     "defense_win_rate",
			Value:      m.DefenseWinRate,
			SampleSize: report.MatchesAnalyzed,
			Context:    "defense side",
		})
	}

	// NEW: Economy round insights
	if m.EconomyStats != nil {
		// Eco round performance
		if m.EconomyStats.EcoRounds > 0 {
			ecoStyle := "average"
			if m.EconomyStats.EcoWinRate > 0.2 {
				ecoStyle = "dangerous on eco"
			} else if m.EconomyStats.EcoWinRate < 0.1 {
				ecoStyle = "predictable saves"
			}
			section.TimingPatterns = append(section.TimingPatterns, intelligence.StrategyInsight{
				Text:       fmt.Sprintf("Eco rounds: %.0f%% win rate (%s)", m.EconomyStats.EcoWinRate*100, ecoStyle),
				Metric:     "eco_win_rate",
				Value:      m.EconomyStats.EcoWinRate,
				SampleSize: m.EconomyStats.EcoRounds,
				Context:    "economy",
			})
		}
		
		// Force buy performance
		if m.EconomyStats.ForceRounds > 0 {
			forceStyle := "average"
			if m.EconomyStats.ForceWinRate > 0.4 {
				forceStyle = "strong force buys"
			} else if m.EconomyStats.ForceWinRate < 0.25 {
				forceStyle = "weak force buys"
			}
			section.TimingPatterns = append(section.TimingPatterns, intelligence.StrategyInsight{
				Text:       fmt.Sprintf("Force buy rounds: %.0f%% win rate (%s)", m.EconomyStats.ForceWinRate*100, forceStyle),
				Metric:     "force_buy_win_rate",
				Value:      m.EconomyStats.ForceWinRate,
				SampleSize: m.EconomyStats.ForceRounds,
				Context:    "economy",
			})
		}
		
		// Full buy performance
		if m.EconomyStats.FullBuyRounds > 0 {
			fullBuyStyle := "average"
			if m.EconomyStats.FullBuyWinRate > 0.55 {
				fullBuyStyle = "dominant when full buying"
			} else if m.EconomyStats.FullBuyWinRate < 0.45 {
				fullBuyStyle = "struggles even with full buy"
			}
			section.TimingPatterns = append(section.TimingPatterns, intelligence.StrategyInsight{
				Text:       fmt.Sprintf("Full buy rounds: %.0f%% win rate (%s)", m.EconomyStats.FullBuyWinRate*100, fullBuyStyle),
				Metric:     "full_buy_win_rate",
				Value:      m.EconomyStats.FullBuyWinRate,
				SampleSize: m.EconomyStats.FullBuyRounds,
				Context:    "economy",
			})
		}
	}

	// Map-specific strategies
	for _, mapEntry := range m.MapPool {
		if mapEntry.GamesPlayed >= 3 {
			mapStrength := "average"
			if mapEntry.Strength == "strong" {
				mapStrength = "comfort pick"
			} else if mapEntry.Strength == "weak" {
				mapStrength = "avoid in veto"
			}
			section.ObjectivePriorities = append(section.ObjectivePriorities, intelligence.StrategyInsight{
				Text:       fmt.Sprintf("%s: %.0f%% win rate (%d games) - %s", mapEntry.MapName, mapEntry.WinRate*100, mapEntry.GamesPlayed, mapStrength),
				Metric:     "map_win_rate",
				Value:      mapEntry.WinRate,
				SampleSize: mapEntry.GamesPlayed,
				Context:    mapEntry.MapName,
			})
		}
	}
}

// formatPlayerTendencies creates hackathon-format player insights
// Example: "Player 'Jett' has a 75% first-duel rate with an Operator on A-main defense"
func (f *Formatter) formatPlayerTendencies(report *intelligence.ScoutingReport) []intelligence.PlayerTendencyInsight {
	tendencies := make([]intelligence.PlayerTendencyInsight, 0)

	for _, player := range report.PlayerProfiles {
		// Signature pick tendency
		if len(player.SignaturePicks) > 0 && len(player.CharacterPool) > 0 {
			topPick := player.CharacterPool[0]
			if topPick.GamesPlayed >= 3 {
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' (%s) signature pick: %s (%.0f%% pick rate, %.0f%% win rate)", player.Nickname, player.Role, topPick.Character, topPick.PickRate*100, topPick.WinRate*100),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "signature_pick",
					Value:        topPick.WinRate,
					Context:      topPick.Character,
					SampleSize:   topPick.GamesPlayed,
				})
			}
		}

		// KDA-based tendency
		if player.KDA > 0 {
			kdaLevel := "average"
			if player.KDA > 4.0 {
				kdaLevel = "exceptional"
			} else if player.KDA > 3.0 {
				kdaLevel = "strong"
			} else if player.KDA < 2.0 {
				kdaLevel = "vulnerable"
			}
			tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
				Text:         fmt.Sprintf("Player '%s' has %.1f KDA (%s performer)", player.Nickname, player.KDA, kdaLevel),
				PlayerName:   player.Nickname,
				Role:         player.Role,
				TendencyType: "kda",
				Value:        player.KDA,
				Context:      kdaLevel,
				SampleSize:   player.GamesPlayed,
			})
		}

		// NEW: Multikill stats (LoL) - HIGH VALUE INSIGHT
		if player.MultikillStats != nil && player.MultikillStats.TotalMultikills > 0 {
			avgMultikillsPerGame := float64(player.MultikillStats.TotalMultikills) / float64(player.GamesPlayed)
			clutchLevel := "average"
			if avgMultikillsPerGame > 2.0 {
				clutchLevel = "exceptional teamfight carry"
			} else if avgMultikillsPerGame > 1.0 {
				clutchLevel = "strong teamfight presence"
			} else if avgMultikillsPerGame < 0.5 {
				clutchLevel = "struggles to carry teamfights"
			}
			
			// Build detailed multikill breakdown
			multikillText := fmt.Sprintf("Player '%s' averages %.1f multikills/game (%s)", 
				player.Nickname, avgMultikillsPerGame, clutchLevel)
			
			// Add specific breakdown if notable
			if player.MultikillStats.PentaKills > 0 {
				multikillText += fmt.Sprintf(" - %d PENTA KILLS!", player.MultikillStats.PentaKills)
			} else if player.MultikillStats.QuadraKills > 0 {
				multikillText += fmt.Sprintf(" - %d quadra kills", player.MultikillStats.QuadraKills)
			}
			
			tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
				Text:         multikillText,
				PlayerName:   player.Nickname,
				Role:         player.Role,
				TendencyType: "multikill",
				Value:        avgMultikillsPerGame,
				Context:      clutchLevel,
				SampleSize:   player.GamesPlayed,
			})
		}

		// NEW: Weapon stats (VALORANT) - HIGH VALUE INSIGHT
		if len(player.WeaponStats) > 0 {
			topWeapon := player.WeaponStats[0]
			if topWeapon.KillShare > 0.3 { // Significant weapon preference
				weaponStyle := "versatile"
				if topWeapon.KillShare > 0.5 {
					weaponStyle = "heavily reliant"
				} else if topWeapon.KillShare > 0.4 {
					weaponStyle = "prefers"
				}
				
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' %s %s (%.0f%% of kills)", player.Nickname, weaponStyle, topWeapon.WeaponName, topWeapon.KillShare*100),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "weapon_preference",
					Value:        topWeapon.KillShare,
					Context:      topWeapon.WeaponName,
					SampleSize:   topWeapon.Kills,
				})
				
				// Special callout for Operator dependency
				if topWeapon.WeaponName == "Operator" && topWeapon.KillShare > 0.35 {
					tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
						Text:         fmt.Sprintf("⚠️ Player '%s' is OPERATOR DEPENDENT (%.0f%% of kills) - deny Op rounds!", player.Nickname, topWeapon.KillShare*100),
						PlayerName:   player.Nickname,
						Role:         player.Role,
						TendencyType: "weapon_dependency",
						Value:        topWeapon.KillShare,
						Context:      "Operator",
						SampleSize:   topWeapon.Kills,
					})
				}
			}
		}

		// NEW: Synergy partners - HIGH VALUE INSIGHT
		if len(player.SynergyPartners) > 0 {
			topPartner := player.SynergyPartners[0]
			if topPartner.SynergyScore > 0.3 { // Strong synergy
				partnerName := topPartner.PlayerName
				if partnerName == "" {
					partnerName = topPartner.PlayerID
				}
				
				synergyLevel := "good"
				if topPartner.SynergyScore > 0.5 {
					synergyLevel = "exceptional"
				} else if topPartner.SynergyScore > 0.4 {
					synergyLevel = "strong"
				}
				
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' has %s synergy with %s (%.0f%% of assists)", player.Nickname, synergyLevel, partnerName, topPartner.SynergyScore*100),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "synergy",
					Value:        topPartner.SynergyScore,
					Context:      partnerName,
					SampleSize:   topPartner.AssistCount,
				})
			}
		}

		// NEW: Assist ratio - identifies playmakers vs carries
		if player.AssistRatio > 0 {
			playstyle := "balanced"
			if player.AssistRatio > 1.5 {
				playstyle = "playmaker (sets up teammates)"
			} else if player.AssistRatio > 1.2 {
				playstyle = "team-oriented"
			} else if player.AssistRatio < 0.7 {
				playstyle = "carry (receives setup)"
			} else if player.AssistRatio < 0.5 {
				playstyle = "isolated carry (low team coordination)"
			}
			
			if playstyle != "balanced" {
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' playstyle: %s (assist ratio: %.2f)", player.Nickname, playstyle, player.AssistRatio),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "playstyle",
					Value:        player.AssistRatio,
					Context:      playstyle,
					SampleSize:   player.GamesPlayed,
				})
			}
		}

		// NEW: Objective focus (LoL) - identifies split-pushers vs teamfighters
		if player.ObjectiveFocus != nil && player.ObjectiveFocus.TowersPerGame > 0 {
			if player.ObjectiveFocus.ObjectiveFocusType == "split-pusher" {
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' is a SPLIT-PUSHER (%.1f towers/game)", player.Nickname, player.ObjectiveFocus.TowersPerGame),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "objective_focus",
					Value:        player.ObjectiveFocus.TowersPerGame,
					Context:      "split-pusher",
					SampleSize:   player.GamesPlayed,
				})
			} else if player.ObjectiveFocus.ObjectiveFocusType == "objective-focused" {
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' is OBJECTIVE-FOCUSED (%.1f dragons/game, %d barons)", player.Nickname, player.ObjectiveFocus.DragonsPerGame, player.ObjectiveFocus.BaronsSecured),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "objective_focus",
					Value:        player.ObjectiveFocus.DragonsPerGame,
					Context:      "objective-focused",
					SampleSize:   player.GamesPlayed,
				})
			}
		}

		// NEW: Item build patterns (LoL) - identifies core items
		if len(player.ItemBuilds) >= 3 {
			// Show top 3 core items
			coreItems := make([]string, 0)
			for i, item := range player.ItemBuilds {
				if i >= 3 {
					break
				}
				if item.BuildRate > 0.5 { // Built in >50% of games
					coreItems = append(coreItems, fmt.Sprintf("%s (%.0f%%)", item.ItemName, item.BuildRate*100))
				}
			}
			if len(coreItems) > 0 {
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' core items: %s", player.Nickname, strings.Join(coreItems, ", ")),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "item_build",
					Value:        float64(len(coreItems)),
					Context:      "core_items",
					SampleSize:   player.GamesPlayed,
				})
			}
		}

		// NEW: Ability usage patterns - identifies ability-reliant players
		if len(player.AbilityUsage) > 0 {
			topAbility := player.AbilityUsage[0]
			if topAbility.UsagePerGame > 10 { // High ability usage
				tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
					Text:         fmt.Sprintf("Player '%s' heavily uses %s (%.1f/game)", player.Nickname, topAbility.AbilityName, topAbility.UsagePerGame),
					PlayerName:   player.Nickname,
					Role:         player.Role,
					TendencyType: "ability_usage",
					Value:        topAbility.UsagePerGame,
					Context:      topAbility.AbilityName,
					SampleSize:   player.GamesPlayed,
				})
			}
		}

		// Weakness-based tendency
		for _, weakness := range player.Weaknesses {
			tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
				Text:         fmt.Sprintf("Player '%s' weakness: %s (%.1f)", player.Nickname, weakness.Title, weakness.Value),
				PlayerName:   player.Nickname,
				Role:         player.Role,
				TendencyType: "weakness",
				Value:        weakness.Value,
				Context:      weakness.Description,
				SampleSize:   weakness.SampleSize,
			})
		}

		// Champion/Agent pool depth
		poolDepth := len(player.CharacterPool)
		if poolDepth > 0 {
			poolDescription := "limited"
			if poolDepth >= 5 {
				poolDescription = "deep"
			} else if poolDepth >= 3 {
				poolDescription = "moderate"
			}
			tendencies = append(tendencies, intelligence.PlayerTendencyInsight{
				Text:         fmt.Sprintf("Player '%s' has %s champion pool (%d picks)", player.Nickname, poolDescription, poolDepth),
				PlayerName:   player.Nickname,
				Role:         player.Role,
				TendencyType: "champion_pool",
				Value:        float64(poolDepth),
				Context:      poolDescription,
				SampleSize:   player.GamesPlayed,
			})
		}
	}

	return tendencies
}

// formatCompositions creates hackathon-format composition insights
// Example: "Most-played comp (68% on Split): Jett, Raze, Brimstone, Skye, Cypher"
func (f *Formatter) formatCompositions(report *intelligence.ScoutingReport) []intelligence.CompositionInsight {
	insights := make([]intelligence.CompositionInsight, 0)

	if report.Compositions == nil {
		return insights
	}

	// Sort compositions by frequency
	comps := make([]intelligence.Composition, len(report.Compositions.TopCompositions))
	copy(comps, report.Compositions.TopCompositions)
	sort.Slice(comps, func(i, j int) bool {
		return comps[i].Frequency > comps[j].Frequency
	})

	for i, comp := range comps {
		if i >= 5 { // Top 5 compositions
			break
		}
		if comp.GamesPlayed < 2 {
			continue
		}

		// Format characters list
		charList := strings.Join(comp.Characters, ", ")

		// Build insight text
		text := fmt.Sprintf("Comp #%d (%.0f%% frequency, %.0f%% win rate): %s",
			i+1, comp.Frequency*100, comp.WinRate*100, charList)

		if comp.Archetype != "" {
			text = fmt.Sprintf("%s (%s style)", text, comp.Archetype)
		}

		insights = append(insights, intelligence.CompositionInsight{
			Text:        text,
			Characters:  comp.Characters,
			Frequency:   comp.Frequency,
			WinRate:     comp.WinRate,
			Archetype:   comp.Archetype,
			GamesPlayed: comp.GamesPlayed,
		})
	}

	// Add draft priorities
	if len(report.Compositions.FirstPickPriorities) > 0 {
		topPicks := make([]string, 0)
		for i, pick := range report.Compositions.FirstPickPriorities {
			if i >= 3 {
				break
			}
			topPicks = append(topPicks, fmt.Sprintf("%s (%.0f%%)", pick.Character, pick.Rate*100))
		}
		insights = append(insights, intelligence.CompositionInsight{
			Text:        fmt.Sprintf("First pick priorities: %s", strings.Join(topPicks, ", ")),
			Characters:  []string{},
			Frequency:   0,
			WinRate:     0,
			GamesPlayed: report.MatchesAnalyzed,
		})
	}

	// Add common bans
	if len(report.Compositions.CommonBans) > 0 {
		topBans := make([]string, 0)
		for i, ban := range report.Compositions.CommonBans {
			if i >= 3 {
				break
			}
			topBans = append(topBans, fmt.Sprintf("%s (%.0f%%)", ban.Character, ban.Rate*100))
		}
		insights = append(insights, intelligence.CompositionInsight{
			Text:        fmt.Sprintf("Common bans against: %s", strings.Join(topBans, ", ")),
			Characters:  []string{},
			Frequency:   0,
			WinRate:     0,
			GamesPlayed: report.MatchesAnalyzed,
		})
	}

	return insights
}


// formatHowToWin creates THE KEY DIFFERENTIATOR - specific, actionable, data-backed recommendations
// This is what separates a winning hackathon entry from "just another stats tool"
func (f *Formatter) formatHowToWin(report *intelligence.ScoutingReport) intelligence.HowToWinSection {
	section := intelligence.HowToWinSection{
		ActionableInsights: make([]intelligence.ActionableInsight, 0),
		DraftStrategy: intelligence.DraftStrategySection{
			PriorityBans:     make([]intelligence.DraftInsight, 0),
			RecommendedPicks: make([]intelligence.DraftInsight, 0),
			TargetPicks:      make([]intelligence.DraftInsight, 0),
		},
		InGameStrategy:  make([]intelligence.InGameStrategyInsight, 0),
		ConfidenceScore: 0,
	}

	if report.HowToWin == nil {
		return section
	}

	// Win condition - THE HEADLINE
	section.WinCondition = report.HowToWin.WinCondition
	section.ConfidenceScore = report.HowToWin.ConfidenceScore

	// Convert weaknesses to actionable insights
	for _, weakness := range report.HowToWin.Weaknesses {
		insight := intelligence.ActionableInsight{
			Recommendation: fmt.Sprintf("Exploit: %s", weakness.Title),
			DataBacking:    fmt.Sprintf("%s - %s", weakness.Description, weakness.Evidence),
			Impact:         categorizeImpact(weakness.Impact),
			ActionType:     "STRATEGY",
			Confidence:     weakness.Impact / 100,
		}
		section.ActionableInsights = append(section.ActionableInsights, insight)
	}

	// Convert draft recommendations
	for _, rec := range report.HowToWin.DraftRecommendations {
		draftInsight := intelligence.DraftInsight{
			Text:      fmt.Sprintf("%s %s - %s", strings.ToUpper(rec.Type), rec.Character, rec.Reason),
			Character: rec.Character,
			Priority:  rec.Priority,
		}

		switch rec.Type {
		case "ban":
			section.DraftStrategy.PriorityBans = append(section.DraftStrategy.PriorityBans, draftInsight)
		case "pick":
			section.DraftStrategy.RecommendedPicks = append(section.DraftStrategy.RecommendedPicks, draftInsight)
		case "target":
			section.DraftStrategy.TargetPicks = append(section.DraftStrategy.TargetPicks, draftInsight)
		}
	}

	// Convert in-game strategies
	for _, strat := range report.HowToWin.InGameStrategies {
		section.InGameStrategy = append(section.InGameStrategy, intelligence.InGameStrategyInsight{
			Strategy: strat.Title,
			Timing:   strat.Timing,
			Reason:   fmt.Sprintf("%s - %s", strat.Description, strat.Evidence),
			Impact:   "HIGH",
		})
	}

	// Add target player insights
	for _, target := range report.HowToWin.TargetPlayers {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: fmt.Sprintf("Target %s (%s)", target.PlayerName, target.Role),
			DataBacking:    target.Reason,
			Impact:         "HIGH",
			ActionType:     "TARGET_PLAYER",
			Confidence:     float64(10-target.Priority) / 10,
		})
	}

	// Generate game-specific actionable insights
	if report.Title == "lol" {
		f.addLoLActionableInsights(&section, report)
	} else {
		f.addVALActionableInsights(&section, report)
	}

	return section
}

// addLoLActionableInsights adds LoL-specific actionable recommendations
func (f *Formatter) addLoLActionableInsights(section *intelligence.HowToWinSection, report *intelligence.ScoutingReport) {
	m := report.TeamStrategy.LoLMetrics
	if m == nil {
		return
	}

	// Early game exploitation
	if m.FirstBloodRate < 0.4 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Play aggressive early - invade and force fights",
			DataBacking:    fmt.Sprintf("Opponent has only %.0f%% first blood rate, indicating passive early game", m.FirstBloodRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.8,
		})
	}

	// Dragon control exploitation
	if m.FirstDragonRate < 0.4 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Prioritize dragon control - set up vision and contest every spawn",
			DataBacking:    fmt.Sprintf("Opponent secures first dragon only %.0f%% of games", m.FirstDragonRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.85,
		})
	}

	// Game pace exploitation
	if m.AvgGameDuration < 28 {
		section.InGameStrategy = append(section.InGameStrategy, intelligence.InGameStrategyInsight{
			Strategy: "Draft scaling compositions and survive early",
			Timing:   "Draft phase and 0-15 minutes",
			Reason:   fmt.Sprintf("Opponent wins fast (avg %.0f min) - they may struggle in extended games", m.AvgGameDuration),
			Impact:   "HIGH",
		})
	} else if m.AvgGameDuration > 35 {
		section.InGameStrategy = append(section.InGameStrategy, intelligence.InGameStrategyInsight{
			Strategy: "Draft early-game compositions and force tempo",
			Timing:   "Draft phase and 0-20 minutes",
			Reason:   fmt.Sprintf("Opponent prefers long games (avg %.0f min) - end before they scale", m.AvgGameDuration),
			Impact:   "HIGH",
		})
	}

	// Player-specific insights from character pools
	for _, player := range report.PlayerProfiles {
		// Find weak champions
		for _, char := range player.CharacterPool {
			if char.GamesPlayed >= 3 && char.WinRate < 0.4 {
				section.DraftStrategy.TargetPicks = append(section.DraftStrategy.TargetPicks, intelligence.DraftInsight{
					Text:       fmt.Sprintf("Force %s onto %s - %.0f%% win rate", player.Nickname, char.Character, char.WinRate*100),
					Character:  char.Character,
					PlayerName: player.Nickname,
					WinRate:    char.WinRate,
					SampleSize: char.GamesPlayed,
					Priority:   2,
				})
			}
		}

		// Find strong champions to ban
		for _, char := range player.CharacterPool {
			if char.GamesPlayed >= 3 && char.WinRate > 0.7 {
				section.DraftStrategy.PriorityBans = append(section.DraftStrategy.PriorityBans, intelligence.DraftInsight{
					Text:       fmt.Sprintf("Ban %s from %s - %.0f%% win rate", char.Character, player.Nickname, char.WinRate*100),
					Character:  char.Character,
					PlayerName: player.Nickname,
					WinRate:    char.WinRate,
					SampleSize: char.GamesPlayed,
					Priority:   1,
				})
			}
		}
	}
}

// addVALActionableInsights adds VALORANT-specific actionable recommendations
func (f *Formatter) addVALActionableInsights(section *intelligence.HowToWinSection, report *intelligence.ScoutingReport) {
	m := report.TeamStrategy.VALMetrics
	if m == nil {
		return
	}

	// Attack side exploitation
	if m.AttackWinRate < 0.45 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Play solid defense and force them into attack rounds",
			DataBacking:    fmt.Sprintf("Opponent has only %.0f%% attack round win rate", m.AttackWinRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.85,
		})
	}

	// Defense side exploitation
	if m.DefenseWinRate < 0.45 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Execute quickly on attack to exploit weak defensive setups",
			DataBacking:    fmt.Sprintf("Opponent has only %.0f%% defense round win rate", m.DefenseWinRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.85,
		})
	}

	// Pistol round exploitation
	if m.PistolWinRate < 0.4 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Focus on pistol round preparation - win pistols to build economy advantage",
			DataBacking:    fmt.Sprintf("Opponent has only %.0f%% pistol round win rate", m.PistolWinRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.8,
		})
	}

	// First blood exploitation
	if m.FirstDeathRate > 0.55 {
		section.ActionableInsights = append(section.ActionableInsights, intelligence.ActionableInsight{
			Recommendation: "Take aggressive early peeks to secure first blood",
			DataBacking:    fmt.Sprintf("Opponent gives up first blood %.0f%% of rounds", m.FirstDeathRate*100),
			Impact:         "HIGH",
			ActionType:     "STRATEGY",
			Confidence:     0.8,
		})
	}

	// Map veto recommendations
	for _, mapEntry := range m.MapPool {
		if mapEntry.Strength == "weak" && mapEntry.GamesPlayed >= 3 {
			section.InGameStrategy = append(section.InGameStrategy, intelligence.InGameStrategyInsight{
				Strategy: fmt.Sprintf("Force %s in map veto", mapEntry.MapName),
				Timing:   "Map veto phase",
				Reason:   fmt.Sprintf("Opponent has only %.0f%% win rate on %s (%d games)", mapEntry.WinRate*100, mapEntry.MapName, mapEntry.GamesPlayed),
				Impact:   "HIGH",
			})
		}
	}

	// Agent-specific insights
	for _, player := range report.PlayerProfiles {
		// Find weak agents
		for _, agent := range player.CharacterPool {
			if agent.GamesPlayed >= 3 && agent.WinRate < 0.4 {
				section.DraftStrategy.TargetPicks = append(section.DraftStrategy.TargetPicks, intelligence.DraftInsight{
					Text:       fmt.Sprintf("Force %s onto %s - %.0f%% win rate", player.Nickname, agent.Character, agent.WinRate*100),
					Character:  agent.Character,
					PlayerName: player.Nickname,
					WinRate:    agent.WinRate,
					SampleSize: agent.GamesPlayed,
					Priority:   2,
				})
			}
		}

		// Find strong agents to deny
		for _, agent := range player.CharacterPool {
			if agent.GamesPlayed >= 3 && agent.WinRate > 0.7 {
				section.DraftStrategy.PriorityBans = append(section.DraftStrategy.PriorityBans, intelligence.DraftInsight{
					Text:       fmt.Sprintf("Deny %s from %s - %.0f%% win rate", agent.Character, player.Nickname, agent.WinRate*100),
					Character:  agent.Character,
					PlayerName: player.Nickname,
					WinRate:    agent.WinRate,
					SampleSize: agent.GamesPlayed,
					Priority:   1,
				})
			}
		}
	}
}

// categorizeImpact converts numeric impact to category
func categorizeImpact(impact float64) string {
	if impact >= 70 {
		return "HIGH"
	} else if impact >= 40 {
		return "MEDIUM"
	}
	return "LOW"
}

// FormatTextReport generates a plain text report in hackathon format
func (f *Formatter) FormatTextReport(digestible *intelligence.DigestibleReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("═══════════════════════════════════════════════════════════════\n"))
	sb.WriteString(fmt.Sprintf("  SCOUTING REPORT: %s\n", digestible.TeamName))
	sb.WriteString(fmt.Sprintf("  Generated: %s | Matches Analyzed: %d\n", digestible.GeneratedAt, digestible.MatchesAnalyzed))
	sb.WriteString(fmt.Sprintf("═══════════════════════════════════════════════════════════════\n\n"))

	// Executive Summary
	sb.WriteString("📋 EXECUTIVE SUMMARY\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")
	sb.WriteString(digestible.ExecutiveSummary)
	sb.WriteString("\n\n")

	// Common Strategies
	sb.WriteString("🎯 COMMON STRATEGIES\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")

	if len(digestible.CommonStrategies.AttackPatterns) > 0 {
		sb.WriteString("Attack Patterns:\n")
		for _, pattern := range digestible.CommonStrategies.AttackPatterns {
			sb.WriteString(fmt.Sprintf("  • %s\n", pattern.Text))
		}
	}

	if len(digestible.CommonStrategies.DefenseSetups) > 0 {
		sb.WriteString("Defense Setups:\n")
		for _, setup := range digestible.CommonStrategies.DefenseSetups {
			sb.WriteString(fmt.Sprintf("  • %s\n", setup.Text))
		}
	}

	if len(digestible.CommonStrategies.ObjectivePriorities) > 0 {
		sb.WriteString("Objective Priorities:\n")
		for _, obj := range digestible.CommonStrategies.ObjectivePriorities {
			sb.WriteString(fmt.Sprintf("  • %s\n", obj.Text))
		}
	}

	if len(digestible.CommonStrategies.TimingPatterns) > 0 {
		sb.WriteString("Timing Patterns:\n")
		for _, timing := range digestible.CommonStrategies.TimingPatterns {
			sb.WriteString(fmt.Sprintf("  • %s\n", timing.Text))
		}
	}
	sb.WriteString("\n")

	// Player Tendencies
	sb.WriteString("👤 PLAYER TENDENCIES\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")
	for _, tendency := range digestible.PlayerTendencies {
		sb.WriteString(fmt.Sprintf("  • %s\n", tendency.Text))
	}
	sb.WriteString("\n")

	// Recent Compositions
	sb.WriteString("🎮 RECENT COMPOSITIONS\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")
	for _, comp := range digestible.RecentCompositions {
		sb.WriteString(fmt.Sprintf("  • %s\n", comp.Text))
	}
	sb.WriteString("\n")

	// HOW TO WIN - THE KEY SECTION
	sb.WriteString("🏆 HOW TO WIN\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("WIN CONDITION: %s\n", digestible.HowToWin.WinCondition))
	sb.WriteString(fmt.Sprintf("Confidence: %.0f%%\n\n", digestible.HowToWin.ConfidenceScore))

	if len(digestible.HowToWin.ActionableInsights) > 0 {
		sb.WriteString("Actionable Insights:\n")
		for _, insight := range digestible.HowToWin.ActionableInsights {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", insight.Impact, insight.Recommendation))
			sb.WriteString(fmt.Sprintf("         Data: %s\n", insight.DataBacking))
		}
		sb.WriteString("\n")
	}

	if len(digestible.HowToWin.DraftStrategy.PriorityBans) > 0 {
		sb.WriteString("Draft - Priority Bans:\n")
		for _, ban := range digestible.HowToWin.DraftStrategy.PriorityBans {
			sb.WriteString(fmt.Sprintf("  🚫 %s\n", ban.Text))
		}
		sb.WriteString("\n")
	}

	if len(digestible.HowToWin.DraftStrategy.TargetPicks) > 0 {
		sb.WriteString("Draft - Target Picks (force opponent onto):\n")
		for _, target := range digestible.HowToWin.DraftStrategy.TargetPicks {
			sb.WriteString(fmt.Sprintf("  🎯 %s\n", target.Text))
		}
		sb.WriteString("\n")
	}

	if len(digestible.HowToWin.InGameStrategy) > 0 {
		sb.WriteString("In-Game Strategy:\n")
		for _, strat := range digestible.HowToWin.InGameStrategy {
			sb.WriteString(fmt.Sprintf("  ⚔️  %s (%s)\n", strat.Strategy, strat.Timing))
			sb.WriteString(fmt.Sprintf("      Reason: %s\n", strat.Reason))
		}
	}

	sb.WriteString("\n═══════════════════════════════════════════════════════════════\n")
	sb.WriteString("  Report generated by SCOUT9 - Automated Scouting Report Generator\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")

	return sb.String()
}
