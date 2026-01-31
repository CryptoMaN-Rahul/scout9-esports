package intelligence

import (
	"fmt"

	"scout9/pkg/grid"
)

// EconomyAnalyzer analyzes VALORANT economy round performance
type EconomyAnalyzer struct{}

// NewEconomyAnalyzer creates a new economy analyzer
func NewEconomyAnalyzer() *EconomyAnalyzer {
	return &EconomyAnalyzer{}
}

// AnalyzeEconomyRounds analyzes team performance by economy state
// Returns detailed economy analysis with per-map breakdowns
func (e *EconomyAnalyzer) AnalyzeEconomyRounds(
	teamID string,
	teamName string,
	seriesStates []*grid.SeriesState,
) *EconomyAnalysis {
	analysis := &EconomyAnalysis{
		TeamID: teamID,
		ByMap:  make(map[string]*MapEconomyStats),
	}

	// Aggregate economy stats
	var (
		ecoRounds     int
		ecoWins       int
		forceRounds   int
		forceWins     int
		fullBuyRounds int
		fullBuyWins   int
	)

	// Per-map tracking
	mapEcoData := make(map[string]*mapEconomyAggregator)

	for _, series := range seriesStates {
		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			mapName := game.Map
			if mapName == "" {
				mapName = "Unknown"
			}

			// Initialize map aggregator
			if _, exists := mapEcoData[mapName]; !exists {
				mapEcoData[mapName] = &mapEconomyAggregator{mapName: mapName}
			}
			mapAgg := mapEcoData[mapName]

			// Find our team in this game
			var ourTeam *grid.GameTeam
			for i := range game.Teams {
				if game.Teams[i].ID == teamID || game.Teams[i].Name == teamName {
					ourTeam = &game.Teams[i]
					break
				}
			}

			if ourTeam == nil {
				continue
			}

			// Analyze rounds from segments
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

				roundNum := segment.SequenceNumber
				won := ourSegTeam.Won

				// Classify round economy based on round number patterns
				// In pro VALORANT:
				// - Rounds 1, 13 = Pistol (eco-level)
				// - Rounds 2, 14 = Usually eco/force after pistol
				// - Rounds 3, 15 = Usually force if lost pistol, full buy if won
				// - Other rounds = Typically full buy in pro play
				economyType := classifyRoundEconomy(roundNum)

				switch economyType {
				case "eco":
					ecoRounds++
					mapAgg.ecoRounds++
					if won {
						ecoWins++
						mapAgg.ecoWins++
					}
				case "force":
					forceRounds++
					mapAgg.forceRounds++
					if won {
						forceWins++
						mapAgg.forceWins++
					}
				case "full":
					fullBuyRounds++
					mapAgg.fullBuyRounds++
					if won {
						fullBuyWins++
						mapAgg.fullBuyWins++
					}
				}
			}
		}
	}

	// Calculate overall rates
	if ecoRounds > 0 {
		analysis.EcoRoundWinRate = float64(ecoWins) / float64(ecoRounds)
	}
	analysis.EcoRounds = ecoRounds

	if forceRounds > 0 {
		analysis.ForceWinRate = float64(forceWins) / float64(forceRounds)
	}
	analysis.ForceRounds = forceRounds

	if fullBuyRounds > 0 {
		analysis.FullBuyWinRate = float64(fullBuyWins) / float64(fullBuyRounds)
	}
	analysis.FullBuyRounds = fullBuyRounds

	// Build per-map stats
	for mapName, agg := range mapEcoData {
		mapStats := &MapEconomyStats{
			MapName:   mapName,
			EcoRounds: agg.ecoRounds,
			ForceRounds: agg.forceRounds,
			FullBuyRounds: agg.fullBuyRounds,
		}

		if agg.ecoRounds > 0 {
			mapStats.EcoWinRate = float64(agg.ecoWins) / float64(agg.ecoRounds)
		}
		if agg.forceRounds > 0 {
			mapStats.ForceWinRate = float64(agg.forceWins) / float64(agg.forceRounds)
		}
		if agg.fullBuyRounds > 0 {
			mapStats.FullBuyWinRate = float64(agg.fullBuyWins) / float64(agg.fullBuyRounds)
		}

		analysis.ByMap[mapName] = mapStats
	}

	return analysis
}

// GenerateEconomyInsights generates hackathon-format economy insights
// Example: "They only win 15% of eco rounds - play aggressive on their saves"
func (e *EconomyAnalyzer) GenerateEconomyInsights(analysis *EconomyAnalysis) []EconomyInsight {
	insights := make([]EconomyInsight, 0)

	// Check eco round weakness (< 15% win rate)
	if analysis.EcoRounds >= 3 && analysis.EcoRoundWinRate < 0.15 {
		insights = append(insights, EconomyInsight{
			Text: fmt.Sprintf("They only win %.0f%% of eco rounds - play aggressive on their saves",
				analysis.EcoRoundWinRate*100),
			Type:       "eco_weak",
			Value:      analysis.EcoRoundWinRate,
			SampleSize: analysis.EcoRounds,
			Impact:     "HIGH",
		})
	}

	// Check force buy strength (> 45% win rate)
	if analysis.ForceRounds >= 3 && analysis.ForceWinRate > 0.45 {
		insights = append(insights, EconomyInsight{
			Text: fmt.Sprintf("They win %.0f%% of force buys - respect their force rounds",
				analysis.ForceWinRate*100),
			Type:       "force_strong",
			Value:      analysis.ForceWinRate,
			SampleSize: analysis.ForceRounds,
			Impact:     "MEDIUM",
		})
	}

	// Check force buy weakness (< 25% win rate)
	if analysis.ForceRounds >= 3 && analysis.ForceWinRate < 0.25 {
		insights = append(insights, EconomyInsight{
			Text: fmt.Sprintf("They only win %.0f%% of force buys - punish their force rounds",
				analysis.ForceWinRate*100),
			Type:       "force_weak",
			Value:      analysis.ForceWinRate,
			SampleSize: analysis.ForceRounds,
			Impact:     "HIGH",
		})
	}

	// Check full buy weakness (< 45% win rate - below expected)
	if analysis.FullBuyRounds >= 5 && analysis.FullBuyWinRate < 0.45 {
		insights = append(insights, EconomyInsight{
			Text: fmt.Sprintf("They only win %.0f%% of full buy rounds - they struggle in even fights",
				analysis.FullBuyWinRate*100),
			Type:       "full_buy_weak",
			Value:      analysis.FullBuyWinRate,
			SampleSize: analysis.FullBuyRounds,
			Impact:     "HIGH",
		})
	}

	// Check full buy strength (> 60% win rate)
	if analysis.FullBuyRounds >= 5 && analysis.FullBuyWinRate > 0.60 {
		insights = append(insights, EconomyInsight{
			Text: fmt.Sprintf("They win %.0f%% of full buy rounds - avoid even economy fights",
				analysis.FullBuyWinRate*100),
			Type:       "full_buy_strong",
			Value:      analysis.FullBuyWinRate,
			SampleSize: analysis.FullBuyRounds,
			Impact:     "MEDIUM",
		})
	}

	// Per-map insights for significant differences
	for mapName, mapStats := range analysis.ByMap {
		// Map-specific eco weakness
		if mapStats.EcoRounds >= 2 && mapStats.EcoWinRate < 0.10 {
			insights = append(insights, EconomyInsight{
				Text: fmt.Sprintf("On %s, they win only %.0f%% of eco rounds - exploit their saves",
					mapName, mapStats.EcoWinRate*100),
				Type:       "eco_weak",
				Value:      mapStats.EcoWinRate,
				SampleSize: mapStats.EcoRounds,
				Impact:     "HIGH",
				MapContext: mapName,
			})
		}

		// Map-specific force strength
		if mapStats.ForceRounds >= 2 && mapStats.ForceWinRate > 0.50 {
			insights = append(insights, EconomyInsight{
				Text: fmt.Sprintf("On %s, they win %.0f%% of force buys - be careful on their forces",
					mapName, mapStats.ForceWinRate*100),
				Type:       "force_strong",
				Value:      mapStats.ForceWinRate,
				SampleSize: mapStats.ForceRounds,
				Impact:     "MEDIUM",
				MapContext: mapName,
			})
		}
	}

	return insights
}

// Helper types
type mapEconomyAggregator struct {
	mapName       string
	ecoRounds     int
	ecoWins       int
	forceRounds   int
	forceWins     int
	fullBuyRounds int
	fullBuyWins   int
}

// classifyRoundEconomy classifies a round's economy type based on round number
// In pro VALORANT, economy follows predictable patterns based on round number
func classifyRoundEconomy(roundNum int) string {
	// Pistol rounds (rounds 1 and 13)
	if roundNum == 1 || roundNum == 13 {
		return "eco"
	}

	// Post-pistol rounds (rounds 2, 3, 14, 15) - typically eco or force
	if roundNum == 2 || roundNum == 14 {
		return "eco" // Usually full eco after pistol loss
	}

	if roundNum == 3 || roundNum == 15 {
		return "force" // Usually force buy or bonus round
	}

	// All other rounds are typically full buys in pro play
	return "full"
}


// =============================================================================
// SITE WEAKNESS IDENTIFICATION
// =============================================================================

// SiteWeakness represents an identified site weakness
type SiteWeakness struct {
	MapName     string  `json:"mapName"`
	SiteName    string  `json:"siteName"`
	WinRate     float64 `json:"winRate"`
	Attempts    int     `json:"attempts"`
	Type        string  `json:"type"` // "defense", "attack"
	Insight     string  `json:"insight"`
	Exploitable bool    `json:"exploitable"`
}

// IdentifySiteWeaknesses identifies sites with <40% win rate and ≥3 attempts
func IdentifySiteWeaknesses(siteAnalysis map[string]*SiteAnalysis) []SiteWeakness {
	weaknesses := make([]SiteWeakness, 0)

	for mapName, analysis := range siteAnalysis {
		for siteName, stats := range analysis.Sites {
			// Defense weakness: <40% win rate with ≥3 attempts
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate < 0.40 {
				weaknesses = append(weaknesses, SiteWeakness{
					MapName:     mapName,
					SiteName:    siteName,
					WinRate:     stats.DefenseWinRate,
					Attempts:    stats.DefenseAttempts,
					Type:        "defense",
					Insight:     fmt.Sprintf("Weak defense on %s-Site (%s) - only %.0f%% win rate", siteName, mapName, stats.DefenseWinRate*100),
					Exploitable: true,
				})
			}

			// Attack weakness: <40% win rate with ≥3 attempts
			if stats.AttackAttempts >= 3 && stats.AttackWinRate < 0.40 {
				weaknesses = append(weaknesses, SiteWeakness{
					MapName:     mapName,
					SiteName:    siteName,
					WinRate:     stats.AttackWinRate,
					Attempts:    stats.AttackAttempts,
					Type:        "attack",
					Insight:     fmt.Sprintf("Weak attack on %s-Site (%s) - only %.0f%% success rate", siteName, mapName, stats.AttackWinRate*100),
					Exploitable: false, // Their attack weakness is not directly exploitable
				})
			}
		}
	}

	return weaknesses
}

// GenerateSiteWeaknessInsights generates hackathon-format insights from site weaknesses
func GenerateSiteWeaknessInsights(weaknesses []SiteWeakness) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	for _, w := range weaknesses {
		if w.Type == "defense" && w.Exploitable {
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("Attack %s-Site on %s - they only hold it %.0f%% of the time",
					w.SiteName, w.MapName, w.WinRate*100),
				Metric:     "site_weakness",
				Value:      w.WinRate,
				SampleSize: w.Attempts,
				Context:    w.MapName,
			})
		}
	}

	return insights
}

// GenerateSiteRecommendations generates in-game strategy recommendations from site analysis
func GenerateSiteRecommendations(siteAnalysis map[string]*SiteAnalysis) []InGameStrategyInsight {
	recommendations := make([]InGameStrategyInsight, 0)

	weaknesses := IdentifySiteWeaknesses(siteAnalysis)

	for _, w := range weaknesses {
		if w.Type == "defense" && w.Exploitable {
			recommendations = append(recommendations, InGameStrategyInsight{
				Strategy: fmt.Sprintf("Target %s-Site on %s", w.SiteName, w.MapName),
				Timing:   "Attack rounds",
				Reason:   fmt.Sprintf("Opponent defense win rate is only %.0f%% on this site", w.WinRate*100),
				Impact:   "HIGH",
			})
		}
	}

	return recommendations
}
