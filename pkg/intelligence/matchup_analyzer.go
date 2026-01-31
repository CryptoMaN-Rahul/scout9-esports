package intelligence

import (
	"fmt"
	"sort"

	"scout9/pkg/grid"
)

// MatchupAnalyzer analyzes player performance against specific champion/agent classes
type MatchupAnalyzer struct{}

// NewMatchupAnalyzer creates a new matchup analyzer
func NewMatchupAnalyzer() *MatchupAnalyzer {
	return &MatchupAnalyzer{}
}

// AnalyzeLoLPlayerMatchups analyzes a LoL player's performance against different champion classes
// ENHANCED: Now tracks performance AGAINST opponent classes for hackathon-winning insights
func (a *MatchupAnalyzer) AnalyzeLoLPlayerMatchups(
	playerID string,
	playerName string,
	role string,
	seriesStates []*grid.SeriesState,
) *PlayerMatchupProfile {
	profile := &PlayerMatchupProfile{
		PlayerID:      playerID,
		PlayerName:    playerName,
		Role:          role,
		Matchups:      make([]MatchupStats, 0),
		StrongAgainst: make([]MatchupStats, 0),
		WeakAgainst:   make([]MatchupStats, 0),
	}

	// Track performance AGAINST opponent classes (not just on our own classes)
	// Key: "vsClass" -> aggregated stats (how player performs AGAINST this class)
	vsClassData := make(map[string]*matchupAggregator)
	
	// Also track by our champion class for comparison
	onClassData := make(map[string]*matchupAggregator)

	for _, series := range seriesStates {
		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			// Find our player and their opponent
			var ourPlayer *grid.GamePlayer
			var ourTeam *grid.GameTeam
			var enemyTeam *grid.GameTeam

			for i := range game.Teams {
				for j := range game.Teams[i].Players {
					if game.Teams[i].Players[j].ID == playerID {
						ourPlayer = &game.Teams[i].Players[j]
						ourTeam = &game.Teams[i]
						break
					}
				}
				if ourPlayer != nil {
					break
				}
			}

			if ourPlayer == nil || ourPlayer.Character == "" {
				continue
			}

			// Find enemy team
			for i := range game.Teams {
				if &game.Teams[i] != ourTeam {
					enemyTeam = &game.Teams[i]
					break
				}
			}

			if enemyTeam == nil {
				continue
			}
			
			// Track performance ON our champion class
			ourClass := GetClassCategory(ourPlayer.Character)
			if ourClass != "unknown" {
				agg, exists := onClassData[ourClass]
				if !exists {
					agg = &matchupAggregator{
						playedCharacter: ourClass,
						vsClass:         "all",
					}
					onClassData[ourClass] = agg
				}
				agg.games++
				agg.kills += ourPlayer.Kills
				agg.deaths += ourPlayer.Deaths
				agg.assists += ourPlayer.Assists
				if ourTeam.Won {
					agg.wins++
				}
			}

			// Find lane opponent (same role position)
			// For mid laners, find enemy mid; for top, find enemy top, etc.
			var laneOpponent *grid.GamePlayer
			ourPlayerIndex := -1
			for j := range ourTeam.Players {
				if ourTeam.Players[j].ID == playerID {
					ourPlayerIndex = j
					break
				}
			}
			
			// Match by position index (assumes same order: Top, Jungle, Mid, ADC, Support)
			if ourPlayerIndex >= 0 && ourPlayerIndex < len(enemyTeam.Players) {
				laneOpponent = &enemyTeam.Players[ourPlayerIndex]
			}
			
			// If we found a lane opponent, track performance AGAINST their class
			if laneOpponent != nil && laneOpponent.Character != "" {
				opponentClass := GetClassCategory(laneOpponent.Character)
				if opponentClass != "unknown" {
					key := opponentClass
					agg, exists := vsClassData[key]
					if !exists {
						agg = &matchupAggregator{
							playedCharacter: "all",
							vsClass:         opponentClass,
						}
						vsClassData[key] = agg
					}
					agg.games++
					agg.kills += ourPlayer.Kills
					agg.deaths += ourPlayer.Deaths
					agg.assists += ourPlayer.Assists
					if ourTeam.Won {
						agg.wins++
					}
				}
			}
		}
	}

	// Convert VS class data to MatchupStats
	// This answers: "How does this player perform AGAINST assassins vs AGAINST mages?"
	for vsClass, agg := range vsClassData {
		if agg.games < 2 {
			continue
		}

		stats := MatchupStats{
			PlayedCharacter: "all",
			VsCharacter:     vsClass,
			GamesPlayed:     agg.games,
			AvgKills:        float64(agg.kills) / float64(agg.games),
			AvgDeaths:       float64(agg.deaths) / float64(agg.games),
		}

		if agg.games > 0 {
			stats.WinRate = float64(agg.wins) / float64(agg.games)
		}
		if agg.deaths > 0 {
			stats.KDA = float64(agg.kills+agg.assists) / float64(agg.deaths)
		} else {
			stats.KDA = float64(agg.kills + agg.assists)
		}

		// Classify matchup type
		if stats.WinRate > 0.6 && stats.KDA > 3.0 {
			stats.MatchupType = "favorable"
			profile.StrongAgainst = append(profile.StrongAgainst, stats)
		} else if stats.WinRate < 0.4 || stats.KDA < 2.0 {
			stats.MatchupType = "unfavorable"
			profile.WeakAgainst = append(profile.WeakAgainst, stats)
		} else {
			stats.MatchupType = "even"
		}

		profile.Matchups = append(profile.Matchups, stats)
	}

	// Sort by games played
	sort.Slice(profile.Matchups, func(i, j int) bool {
		return profile.Matchups[i].GamesPlayed > profile.Matchups[j].GamesPlayed
	})
	sort.Slice(profile.StrongAgainst, func(i, j int) bool {
		return profile.StrongAgainst[i].KDA > profile.StrongAgainst[j].KDA
	})
	sort.Slice(profile.WeakAgainst, func(i, j int) bool {
		return profile.WeakAgainst[i].KDA < profile.WeakAgainst[j].KDA
	})

	return profile
}

// AnalyzeLoLPlayerClassPerformance analyzes a player's performance ON different champion classes
// This answers: "How does this player perform on assassins vs control mages?"
func (a *MatchupAnalyzer) AnalyzeLoLPlayerClassPerformance(
	playerID string,
	playerName string,
	seriesStates []*grid.SeriesState,
) map[string]*ClassPerformance {
	// Track performance by champion class played
	classData := make(map[string]*classAggregator)

	for _, series := range seriesStates {
		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			// Find our player
			var ourPlayer *grid.GamePlayer
			var ourTeam *grid.GameTeam

			for i := range game.Teams {
				for j := range game.Teams[i].Players {
					if game.Teams[i].Players[j].ID == playerID {
						ourPlayer = &game.Teams[i].Players[j]
						ourTeam = &game.Teams[i]
						break
					}
				}
				if ourPlayer != nil {
					break
				}
			}

			if ourPlayer == nil || ourPlayer.Character == "" {
				continue
			}

			// Get champion class
			champClass := GetClassCategory(ourPlayer.Character)
			if champClass == "unknown" {
				continue
			}

			agg, exists := classData[champClass]
			if !exists {
				agg = &classAggregator{className: champClass}
				classData[champClass] = agg
			}

			agg.games++
			agg.kills += ourPlayer.Kills
			agg.deaths += ourPlayer.Deaths
			agg.assists += ourPlayer.Assists
			if ourTeam.Won {
				agg.wins++
			}
		}
	}

	// Convert to ClassPerformance
	result := make(map[string]*ClassPerformance)
	for className, agg := range classData {
		if agg.games < 2 {
			continue
		}

		perf := &ClassPerformance{
			ClassName:   className,
			GamesPlayed: agg.games,
			AvgKills:    float64(agg.kills) / float64(agg.games),
			AvgDeaths:   float64(agg.deaths) / float64(agg.games),
		}

		if agg.games > 0 {
			perf.WinRate = float64(agg.wins) / float64(agg.games)
		}
		if agg.deaths > 0 {
			perf.KDA = float64(agg.kills+agg.assists) / float64(agg.deaths)
		} else {
			perf.KDA = float64(agg.kills + agg.assists)
		}

		result[className] = perf
	}

	return result
}

// GenerateMatchupInsights generates hackathon-format insights from matchup data
// HACKATHON FORMAT: "Mid laner has a KDA of 2.1 on control mages vs. 8.5 on assassins"
func (a *MatchupAnalyzer) GenerateMatchupInsights(
	playerName string,
	classPerformance map[string]*ClassPerformance,
) []PlayerTendencyInsight {
	insights := make([]PlayerTendencyInsight, 0)

	// Find best and worst class performance
	var bestClass, worstClass *ClassPerformance
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 2 {
			continue
		}
		if bestClass == nil || perf.KDA > bestClass.KDA {
			bestClass = perf
		}
		if worstClass == nil || perf.KDA < worstClass.KDA {
			worstClass = perf
		}
	}

	// Generate insight if there's a significant difference
	if bestClass != nil && worstClass != nil && bestClass.ClassName != worstClass.ClassName {
		kdaDiff := bestClass.KDA - worstClass.KDA
		if kdaDiff > 1.0 { // Significant difference
			// HACKATHON FORMAT: "Mid laner has a KDA of 2.1 on control mages vs. 8.5 on assassins"
			insight := PlayerTendencyInsight{
				Text: fmt.Sprintf("%s has %.1f KDA on %ss vs. %.1f on %ss",
					playerName, worstClass.KDA, worstClass.ClassName, bestClass.KDA, bestClass.ClassName),
				PlayerName:   playerName,
				TendencyType: "class_performance",
				Value:        kdaDiff,
				Context:      fmt.Sprintf("Strong on %ss, weak on %ss", bestClass.ClassName, worstClass.ClassName),
				SampleSize:   bestClass.GamesPlayed + worstClass.GamesPlayed,
			}
			insights = append(insights, insight)
		}
	}
	
	// Generate insights for each class with significant sample size
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 3 {
			continue
		}
		
		// Highlight very strong or very weak performance
		if perf.KDA > 4.0 && perf.WinRate > 0.6 {
			insights = append(insights, PlayerTendencyInsight{
				Text: fmt.Sprintf("%s excels on %ss (%.1f KDA, %.0f%% win rate)",
					playerName, perf.ClassName, perf.KDA, perf.WinRate*100),
				PlayerName:   playerName,
				TendencyType: "class_strength",
				Value:        perf.KDA,
				Context:      perf.ClassName,
				SampleSize:   perf.GamesPlayed,
			})
		} else if perf.KDA < 2.0 || perf.WinRate < 0.35 {
			insights = append(insights, PlayerTendencyInsight{
				Text: fmt.Sprintf("%s struggles on %ss (%.1f KDA, %.0f%% win rate)",
					playerName, perf.ClassName, perf.KDA, perf.WinRate*100),
				PlayerName:   playerName,
				TendencyType: "class_weakness",
				Value:        perf.KDA,
				Context:      perf.ClassName,
				SampleSize:   perf.GamesPlayed,
			})
		}
	}

	return insights
}

// GenerateVsClassInsights generates insights about performance AGAINST specific classes
// HACKATHON FORMAT: "Mid laner has 15% win rate when facing assassins"
func (a *MatchupAnalyzer) GenerateVsClassInsights(
	playerName string,
	matchupProfile *PlayerMatchupProfile,
) []PlayerTendencyInsight {
	insights := make([]PlayerTendencyInsight, 0)
	
	// Find best and worst matchups (vs opponent classes)
	var bestMatchup, worstMatchup *MatchupStats
	for i := range matchupProfile.Matchups {
		m := &matchupProfile.Matchups[i]
		if m.GamesPlayed < 2 {
			continue
		}
		if bestMatchup == nil || m.KDA > bestMatchup.KDA {
			bestMatchup = m
		}
		if worstMatchup == nil || m.KDA < worstMatchup.KDA {
			worstMatchup = m
		}
	}
	
	// Generate insight if there's a significant difference
	if bestMatchup != nil && worstMatchup != nil && bestMatchup.VsCharacter != worstMatchup.VsCharacter {
		kdaDiff := bestMatchup.KDA - worstMatchup.KDA
		if kdaDiff > 1.5 {
			// HACKATHON FORMAT
			insights = append(insights, PlayerTendencyInsight{
				Text: fmt.Sprintf("%s has %.1f KDA vs %ss but only %.1f vs %ss",
					playerName, bestMatchup.KDA, bestMatchup.VsCharacter, worstMatchup.KDA, worstMatchup.VsCharacter),
				PlayerName:   playerName,
				TendencyType: "matchup",
				Value:        kdaDiff,
				Context:      fmt.Sprintf("Strong vs %ss, weak vs %ss", bestMatchup.VsCharacter, worstMatchup.VsCharacter),
				SampleSize:   bestMatchup.GamesPlayed + worstMatchup.GamesPlayed,
			})
		}
	}
	
	// Highlight specific weak matchups for targeting
	for _, m := range matchupProfile.WeakAgainst {
		if m.GamesPlayed >= 2 && m.WinRate < 0.4 {
			insights = append(insights, PlayerTendencyInsight{
				Text: fmt.Sprintf("%s has %.0f%% win rate when facing %ss",
					playerName, m.WinRate*100, m.VsCharacter),
				PlayerName:   playerName,
				TendencyType: "weak_matchup",
				Value:        m.WinRate,
				Context:      fmt.Sprintf("Target with %ss", m.VsCharacter),
				SampleSize:   m.GamesPlayed,
			})
		}
	}
	
	return insights
}

// GenerateDraftRecommendations generates draft recommendations based on matchup weaknesses
// HACKATHON FORMAT: "Recommend picking LeBlanc and Zed - opponent mid has 15% win rate vs assassins"
func (a *MatchupAnalyzer) GenerateDraftRecommendations(
	playerName string,
	role string,
	classPerformance map[string]*ClassPerformance,
) []DraftInsight {
	recommendations := make([]DraftInsight, 0)

	// Find weak class
	var weakestClass *ClassPerformance
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 2 {
			continue
		}
		if weakestClass == nil || perf.KDA < weakestClass.KDA {
			weakestClass = perf
		}
	}

	if weakestClass != nil && weakestClass.KDA < 2.5 {
		// Generate recommendation to pick champions that counter their weak class
		// HACKATHON FORMAT: "Recommend picking LeBlanc and Zed"
		var recommendedPicks string
		var specificChamps string
		switch weakestClass.ClassName {
		case "mage":
			recommendedPicks = "assassins"
			specificChamps = "Zed, LeBlanc, Akali"
		case "assassin":
			recommendedPicks = "tanks or fighters"
			specificChamps = "Ornn, Aatrox, Renekton"
		case "tank":
			recommendedPicks = "sustained damage dealers"
			specificChamps = "Vayne, Kog'Maw, Cassiopeia"
		case "fighter":
			recommendedPicks = "control mages with CC"
			specificChamps = "Orianna, Syndra, Viktor"
		case "marksman":
			recommendedPicks = "assassins or divers"
			specificChamps = "Zed, Diana, Talon"
		case "support":
			recommendedPicks = "engage supports"
			specificChamps = "Nautilus, Leona, Thresh"
		default:
			recommendedPicks = "champions that counter their playstyle"
			specificChamps = ""
		}

		rec := DraftInsight{
			Text: fmt.Sprintf("Pick %s vs %s - they have %.1f KDA on %ss (%s)",
				recommendedPicks, playerName, weakestClass.KDA, weakestClass.ClassName, specificChamps),
			PlayerName: playerName,
			Character:  specificChamps,
			WinRate:    weakestClass.WinRate,
			SampleSize: weakestClass.GamesPlayed,
			Priority:   1,
		}
		recommendations = append(recommendations, rec)
	}
	
	// Also recommend banning their strong class champions
	var strongestClass *ClassPerformance
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 3 {
			continue
		}
		if strongestClass == nil || perf.KDA > strongestClass.KDA {
			strongestClass = perf
		}
	}
	
	if strongestClass != nil && strongestClass.KDA > 4.0 && strongestClass.WinRate > 0.6 {
		rec := DraftInsight{
			Text: fmt.Sprintf("Consider banning %s %ss - %s has %.1f KDA (%.0f%% WR)",
				strongestClass.ClassName, "champions", playerName, strongestClass.KDA, strongestClass.WinRate*100),
			PlayerName: playerName,
			Character:  strongestClass.ClassName,
			WinRate:    strongestClass.WinRate,
			SampleSize: strongestClass.GamesPlayed,
			Priority:   2,
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

// ClassPerformance tracks performance on a specific champion class
type ClassPerformance struct {
	ClassName   string
	GamesPlayed int
	WinRate     float64
	KDA         float64
	AvgKills    float64
	AvgDeaths   float64
}

// Helper types for aggregation
type matchupAggregator struct {
	playedCharacter string
	vsClass         string
	games           int
	wins            int
	kills           int
	deaths          int
	assists         int
}

type classAggregator struct {
	className string
	games     int
	wins      int
	kills     int
	deaths    int
	assists   int
}


// =============================================================================
// ENHANCED COUNTER-STRATEGY METHODS
// =============================================================================

// CounterClassMap maps weak classes to counter picks with specific champion names
var CounterClassMap = map[string]CounterClassRecommendation{
	"mage": {
		WeakClass:    "mage",
		CounterClass: "assassin",
		Champions:    []string{"Zed", "LeBlanc", "Akali", "Talon", "Fizz", "Katarina"},
	},
	"assassin": {
		WeakClass:    "assassin",
		CounterClass: "tank",
		Champions:    []string{"Ornn", "Malphite", "Galio", "Aatrox", "Renekton", "Sett"},
	},
	"tank": {
		WeakClass:    "tank",
		CounterClass: "marksman",
		Champions:    []string{"Vayne", "Kog'Maw", "Cassiopeia", "Kai'Sa", "Jinx"},
	},
	"fighter": {
		WeakClass:    "fighter",
		CounterClass: "mage",
		Champions:    []string{"Orianna", "Syndra", "Viktor", "Azir", "Ryze"},
	},
	"marksman": {
		WeakClass:    "marksman",
		CounterClass: "assassin",
		Champions:    []string{"Zed", "Diana", "Talon", "Akali", "Qiyana"},
	},
	"support": {
		WeakClass:    "support",
		CounterClass: "engage",
		Champions:    []string{"Nautilus", "Leona", "Thresh", "Alistar", "Rakan"},
	},
}

// GenerateClassMatchupInsights generates ClassMatchupInsight structs for hackathon-format output
// Example output: "Faker has 2.1 KDA on control mages vs 8.5 on assassins"
func (a *MatchupAnalyzer) GenerateClassMatchupInsights(
	playerName string,
	role string,
	classPerformance map[string]*ClassPerformance,
) []ClassMatchupInsight {
	insights := make([]ClassMatchupInsight, 0)

	// Find best and worst class performance
	var bestClass, worstClass *ClassPerformance
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 2 {
			continue
		}
		if bestClass == nil || perf.KDA > bestClass.KDA {
			bestClass = perf
		}
		if worstClass == nil || perf.KDA < worstClass.KDA {
			worstClass = perf
		}
	}

	// Generate insight if there's a significant difference (>1.5 KDA diff)
	if bestClass != nil && worstClass != nil && bestClass.ClassName != worstClass.ClassName {
		kdaDiff := bestClass.KDA - worstClass.KDA
		if kdaDiff > 1.5 {
			// Calculate confidence based on sample size
			totalGames := bestClass.GamesPlayed + worstClass.GamesPlayed
			confidence := calculateInsightConfidence(totalGames)

			insight := ClassMatchupInsight{
				PlayerName:     playerName,
				Role:           role,
				WeakClass:      worstClass.ClassName,
				WeakClassKDA:   worstClass.KDA,
				WeakClassWR:    worstClass.WinRate,
				StrongClass:    bestClass.ClassName,
				StrongClassKDA: bestClass.KDA,
				StrongClassWR:  bestClass.WinRate,
				Text: fmt.Sprintf("%s has %.1f KDA on %ss vs %.1f on %ss - pick %ss!",
					playerName, worstClass.KDA, worstClass.ClassName,
					bestClass.KDA, bestClass.ClassName,
					getCounterClass(worstClass.ClassName)),
				SampleSize: totalGames,
				Confidence: confidence,
			}
			insights = append(insights, insight)
		}
	}

	// Also generate insights for individual weak classes
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 2 {
			continue
		}

		// Weak class: KDA < 2.5
		if perf.KDA < 2.5 {
			confidence := calculateInsightConfidence(perf.GamesPlayed)
			counterClass := getCounterClass(perf.ClassName)

			insight := ClassMatchupInsight{
				PlayerName:   playerName,
				Role:         role,
				WeakClass:    perf.ClassName,
				WeakClassKDA: perf.KDA,
				WeakClassWR:  perf.WinRate,
				Text: fmt.Sprintf("%s struggles on %ss (%.1f KDA, %.0f%% WR) - counter with %ss",
					playerName, perf.ClassName, perf.KDA, perf.WinRate*100, counterClass),
				SampleSize: perf.GamesPlayed,
				Confidence: confidence,
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

// GenerateSpecificDraftRecommendations generates draft recommendations with specific champion names
// Example: "Pick Zed, LeBlanc, Akali vs Faker - they have 2.1 KDA on mages"
func (a *MatchupAnalyzer) GenerateSpecificDraftRecommendations(
	playerName string,
	role string,
	classPerformance map[string]*ClassPerformance,
) []SpecificDraftRecommendation {
	recommendations := make([]SpecificDraftRecommendation, 0)

	// Find weak classes (KDA < 2.5 with ≥2 games)
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 2 || perf.KDA >= 2.5 {
			continue
		}

		// Get counter class and specific champions
		counterRec, exists := CounterClassMap[perf.ClassName]
		if !exists {
			continue
		}

		confidence := calculateInsightConfidence(perf.GamesPlayed)

		rec := SpecificDraftRecommendation{
			Type:         "pick",
			TargetPlayer: playerName,
			TargetRole:   role,
			TargetClass:  perf.ClassName,
			Champions:    counterRec.Champions,
			Reason: fmt.Sprintf("Pick %s vs %s - they have %.1f KDA on %ss (%.0f%% WR)",
				formatChampionList(counterRec.Champions[:min(3, len(counterRec.Champions))]),
				playerName, perf.KDA, perf.ClassName, perf.WinRate*100),
			WinRate:    perf.WinRate,
			KDA:        perf.KDA,
			SampleSize: perf.GamesPlayed,
			Priority:   calculatePriority(perf.KDA, perf.GamesPlayed),
			Confidence: confidence,
		}
		recommendations = append(recommendations, rec)
	}

	// Find strong classes to ban (KDA > 4.0 and WR > 60% with ≥3 games)
	for _, perf := range classPerformance {
		if perf.GamesPlayed < 3 || perf.KDA <= 4.0 || perf.WinRate <= 0.6 {
			continue
		}

		// Get champions from this class to ban
		banChamps := getChampionsForClass(perf.ClassName)
		if len(banChamps) == 0 {
			continue
		}

		confidence := calculateInsightConfidence(perf.GamesPlayed)

		rec := SpecificDraftRecommendation{
			Type:         "ban",
			TargetPlayer: playerName,
			TargetRole:   role,
			TargetClass:  perf.ClassName,
			Champions:    banChamps,
			Reason: fmt.Sprintf("Ban %s %ss - %s has %.1f KDA (%.0f%% WR)",
				perf.ClassName, "champions", playerName, perf.KDA, perf.WinRate*100),
			WinRate:    perf.WinRate,
			KDA:        perf.KDA,
			SampleSize: perf.GamesPlayed,
			Priority:   1, // High priority for bans
			Confidence: confidence,
		}
		recommendations = append(recommendations, rec)
	}

	// Sort by priority
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority < recommendations[j].Priority
	})

	return recommendations
}

// Helper functions

func calculateInsightConfidence(sampleSize int) float64 {
	if sampleSize >= 15 {
		return 0.9
	} else if sampleSize >= 10 {
		return 0.75
	} else if sampleSize >= 5 {
		return 0.6
	} else if sampleSize >= 2 {
		return 0.4
	}
	return 0.2
}

func calculatePriority(kda float64, games int) int {
	// Lower KDA = higher priority (1 is highest)
	if kda < 1.5 && games >= 3 {
		return 1
	} else if kda < 2.0 && games >= 2 {
		return 2
	} else if kda < 2.5 {
		return 3
	}
	return 4
}

func getCounterClass(weakClass string) string {
	if rec, exists := CounterClassMap[weakClass]; exists {
		return rec.CounterClass
	}
	return "counter"
}

func getChampionsForClass(className string) []string {
	// Return popular champions for each class
	classChampions := map[string][]string{
		"assassin": {"Zed", "LeBlanc", "Akali", "Talon", "Fizz", "Katarina", "Qiyana"},
		"mage":     {"Orianna", "Syndra", "Viktor", "Azir", "Ryze", "Cassiopeia", "Ahri"},
		"tank":     {"Ornn", "Malphite", "Sion", "Maokai", "Sejuani", "Zac", "Cho'Gath"},
		"fighter":  {"Aatrox", "Renekton", "Jax", "Fiora", "Camille", "Irelia", "Riven"},
		"marksman": {"Jinx", "Kai'Sa", "Zeri", "Aphelios", "Varus", "Ezreal", "Caitlyn"},
		"support":  {"Thresh", "Nautilus", "Leona", "Lulu", "Yuumi", "Renata", "Karma"},
	}

	if champs, exists := classChampions[className]; exists {
		return champs
	}
	return []string{}
}

func formatChampionList(champions []string) string {
	if len(champions) == 0 {
		return ""
	}
	if len(champions) == 1 {
		return champions[0]
	}
	if len(champions) == 2 {
		return champions[0] + " or " + champions[1]
	}
	return champions[0] + ", " + champions[1] + ", or " + champions[2]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
