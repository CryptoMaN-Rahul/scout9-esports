package intelligence

import (
	"fmt"
	"sort"

	"scout9/pkg/grid"
)

// TimingAnalyzerEngine analyzes timing-based patterns for LoL
// ENHANCED: Now uses statistical analysis engine for data-driven insights
type TimingAnalyzerEngine struct {
	stats        *StatisticalEngine
	laneDetector *StatisticalLaneDetector
	validator    *Validator
}

// NewTimingAnalyzer creates a new timing analyzer with statistical components
func NewTimingAnalyzer() *TimingAnalyzerEngine {
	return &TimingAnalyzerEngine{
		stats:        NewStatisticalEngine(),
		laneDetector: NewStatisticalLaneDetector(),
		validator:    NewValidator(),
	}
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// determineLaneFromPosition uses the statistical lane detector for accurate classification
// This replaces the old hardcoded approach with data-driven analysis
func determineLaneFromPosition(pos *grid.Position) string {
	if pos == nil {
		return "unknown"
	}

	detector := NewStatisticalLaneDetector()
	result := detector.ClassifyPosition(pos)

	// Map detailed classification to simple lane names
	switch result.Lane {
	case "top":
		return "top"
	case "mid":
		return "mid"
	case "bot":
		return "bot"
	case "jungle":
		if result.SubRegion != "" {
			return result.SubRegion
		}
		return "jungle"
	case "river":
		return "river"
	case "base":
		return "base"
	default:
		return "unknown"
	}
}

// determineJungleSide determines which side of the jungle a position is on
func determineJungleSide(pos *grid.Position, teamSide string) string {
	if pos == nil {
		return "unknown"
	}

	detector := NewStatisticalLaneDetector()
	return detector.GetJungleSide(pos, teamSide)
}

// AnalyzeJunglePathing analyzes jungle pathing patterns from event data
// ENHANCED: Now uses actual position data from kills!
func (a *TimingAnalyzerEngine) AnalyzeJunglePathing(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) *JunglePathingAnalysis {
	analysis := &JunglePathingAnalysis{
		GanksByLane:        make(map[string]float64),
		FirstClearPatterns: make([]ClearPattern, 0),
	}

	var (
		totalGames       int
		preTenBotGanks   int
		preTenTopGanks   int
		preTenMidGanks   int
		totalPreTenGanks int
	)

	// Lane gank tracking with position data
	laneGanks := map[string]int{
		"top": 0,
		"mid": 0,
		"bot": 0,
	}
	totalGanks := 0
	
	// First blood location tracking
	firstBloodLocations := map[string]int{
		"top": 0,
		"mid": 0,
		"bot": 0,
		"jungle": 0,
	}
	totalFirstBloods := 0

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		// Find the series state to get team info
		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			// Find our team and jungler
			var ourTeam *grid.GameTeam
			var junglerID string
			for i := range game.Teams {
				if game.Teams[i].ID == teamID {
					ourTeam = &game.Teams[i]
					// Find jungler by character
					for _, player := range ourTeam.Players {
						if isJungleChampion(player.Character) {
							junglerID = player.ID
							break
						}
					}
					break
				}
			}

			if ourTeam == nil || junglerID == "" {
				continue
			}

			totalGames++

			// Analyze kills for gank patterns using POSITION DATA
			for _, kill := range eventData.Kills {
				// Check if our jungler was involved
				junglerInvolved := kill.KillerID == junglerID
				if !junglerInvolved {
					for _, assistID := range kill.AssistIDs {
						if assistID == junglerID {
							junglerInvolved = true
							break
						}
					}
				}

				// Track first blood location (regardless of jungler involvement)
				if kill.FirstBlood {
					totalFirstBloods++
					// Use position data if available
					pos := kill.KillerPosition
					if pos == nil {
						pos = kill.VictimPosition
					}
					if pos != nil {
						lane := determineLaneFromPosition(pos)
						switch lane {
						case "top", "top_jungle":
							firstBloodLocations["top"]++
						case "mid":
							firstBloodLocations["mid"]++
						case "bot", "bot_jungle":
							firstBloodLocations["bot"]++
						default:
							firstBloodLocations["jungle"]++
						}
					}
				}

				if !junglerInvolved {
					continue
				}

				// Pre-10 minute ganks (600000ms)
				if kill.GameTime < 600000 {
					totalGanks++
					
					// USE POSITION DATA for lane detection!
					pos := kill.KillerPosition
					if pos == nil {
						pos = kill.VictimPosition
					}
					
					if pos != nil {
						lane := determineLaneFromPosition(pos)
						
						switch lane {
						case "top", "top_jungle":
							laneGanks["top"]++
							preTenTopGanks++
							totalPreTenGanks++
						case "mid":
							laneGanks["mid"]++
							preTenMidGanks++
							totalPreTenGanks++
						case "bot", "bot_jungle":
							laneGanks["bot"]++
							preTenBotGanks++
							totalPreTenGanks++
						default:
							// Unknown position, use timing heuristic as fallback
							if kill.GameTime < 360000 { // Pre-6 min
								laneGanks["bot"]++
								preTenBotGanks++
							} else {
								laneGanks["mid"]++
								preTenMidGanks++
							}
							totalPreTenGanks++
						}
					} else {
						// No position data, use timing heuristic
						if kill.GameTime < 360000 {
							laneGanks["bot"]++
							preTenBotGanks++
						} else {
							laneGanks["mid"]++
							preTenMidGanks++
						}
						totalPreTenGanks++
					}
				}
			}
			
			// Analyze dragon kills for jungle side preference
			for _, dragon := range eventData.DragonKills {
				if dragon.TeamID == teamID && dragon.Position != nil {
					// Dragon is always bot-side, but we can track timing
					// Early dragon (pre-10 min) indicates bot-side focus
					if dragon.GameTime < 600000 {
						preTenBotGanks++ // Proxy for bot-side presence
					}
				}
			}
		}
	}

	// Calculate rates with actual data
	if totalGames > 0 {
		if totalPreTenGanks > 0 {
			analysis.PreSixBotRate = float64(preTenBotGanks) / float64(totalPreTenGanks)
			analysis.PreSixTopRate = float64(preTenTopGanks) / float64(totalPreTenGanks)
		}

		// Determine path preference based on actual gank locations
		if analysis.PreSixBotRate > 0.6 {
			analysis.PreSixPathPreference = "bot-focused"
		} else if analysis.PreSixTopRate > 0.6 {
			analysis.PreSixPathPreference = "top-focused"
		} else if preTenMidGanks > preTenBotGanks && preTenMidGanks > preTenTopGanks {
			analysis.PreSixPathPreference = "mid-focused"
		} else {
			analysis.PreSixPathPreference = "balanced"
		}

		// Calculate gank distribution
		if totalGanks > 0 {
			for lane, count := range laneGanks {
				analysis.GanksByLane[lane] = float64(count) / float64(totalGanks)
			}
		}
		
		// Counter-jungle rate estimation from enemy jungle kills
		// (This would need more detailed event parsing for accuracy)
		analysis.CounterJungleRate = 0.0 // TODO: Implement with more event data
	}

	return analysis
}

// AnalyzeObjectiveTimings analyzes objective timing patterns
// ENHANCED: Now tracks dragon types and provides more detailed insights
func (a *TimingAnalyzerEngine) AnalyzeObjectiveTimings(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) *ObjectiveTimingAnalysis {
	analysis := &ObjectiveTimingAnalysis{
		BaronAttemptTimings: make([]float64, 0),
	}

	var (
		totalGames          int
		firstDragonTimes    []float64
		firstTowerTimes     []float64
		firstHeraldTimes    []float64
		towerLaneCounts     = map[string]int{"top": 0, "mid": 0, "bot": 0}
		dragonContests      int
		
		// Dragon type tracking for priority analysis
		dragonTypeCounts    = map[string]int{}
		totalDragons        int
		
		// Herald usage tracking
		heraldLaneCounts    = map[string]int{"top": 0, "mid": 0, "bot": 0}
		totalHeralds        int
		
		// Void grub tracking
		voidGrubTakes       int
		voidGrubGames       int
	)

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		// Find the series state
		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			// Find our team
			var ourTeam *grid.GameTeam
			for i := range game.Teams {
				if game.Teams[i].ID == teamID {
					ourTeam = &game.Teams[i]
					break
				}
			}

			if ourTeam == nil {
				continue
			}

			totalGames++
			gameHadFirstDragon := false
			gameHadFirstTower := false
			gameHadFirstHerald := false
			
			// Track if team got the FIRST dragon of the game (not just their first dragon)
			// This is important for contest rate calculation
			firstDragonOfGame := true

			// Track dragon kills with type information
			for _, dragon := range eventData.DragonKills {
				// Check if this is the first dragon of the game
				if firstDragonOfGame {
					// If our team got the first dragon, count it as a contest win
					if dragon.TeamID == teamID || isTeamID(dragon.TeamID, ourTeam) {
						dragonContests++
					}
					firstDragonOfGame = false
				}
				
				if dragon.TeamID == teamID || isTeamID(dragon.TeamID, ourTeam) {
					totalDragons++
					
					// Track dragon type
					if dragon.DragonType != "" && dragon.DragonType != "unknown" {
						dragonTypeCounts[dragon.DragonType]++
					}
					
					// First dragon timing for our team
					if !gameHadFirstDragon {
						timeMinutes := float64(dragon.GameTime) / 60000.0
						firstDragonTimes = append(firstDragonTimes, timeMinutes)
						gameHadFirstDragon = true
					}
				}
			}

			// Track first tower timing and lane with actual lane data
			for _, tower := range eventData.TowerDestroys {
				if (tower.TeamID == teamID || isTeamID(tower.TeamID, ourTeam)) && !gameHadFirstTower {
					timeMinutes := float64(tower.GameTime) / 60000.0
					firstTowerTimes = append(firstTowerTimes, timeMinutes)
					
					// Use actual lane data from event
					lane := tower.Lane
					if lane == "" {
						lane = "bot" // Default assumption
					}
					towerLaneCounts[lane]++
					gameHadFirstTower = true
				}
			}

			// Track herald timing and usage
			for _, herald := range eventData.HeraldKills {
				if herald.TeamID == teamID || isTeamID(herald.TeamID, ourTeam) {
					totalHeralds++
					
					if !gameHadFirstHerald {
						timeMinutes := float64(herald.GameTime) / 60000.0
						firstHeraldTimes = append(firstHeraldTimes, timeMinutes)
						gameHadFirstHerald = true
					}
					
					// Track where herald was used (based on subsequent tower destroy)
					// Look for tower destroy within 60 seconds of herald kill
					heraldTime := herald.GameTime
					for _, tower := range eventData.TowerDestroys {
						if tower.GameTime > heraldTime && tower.GameTime < heraldTime+60000 {
							if tower.TeamID == teamID || isTeamID(tower.TeamID, ourTeam) {
								lane := tower.Lane
								if lane != "" {
									heraldLaneCounts[lane]++
								}
								break
							}
						}
					}
				}
			}

			// Track void grub takes
			if len(eventData.VoidGrubKills) > 0 {
				voidGrubGames++
				for _, grub := range eventData.VoidGrubKills {
					if grub.TeamID == teamID || isTeamID(grub.TeamID, ourTeam) {
						voidGrubTakes++
					}
				}
			}

			// Track baron timings
			for _, baron := range eventData.BaronKills {
				if baron.TeamID == teamID || isTeamID(baron.TeamID, ourTeam) {
					timeMinutes := float64(baron.GameTime) / 60000.0
					analysis.BaronAttemptTimings = append(analysis.BaronAttemptTimings, timeMinutes)
				}
			}
		}
	}

	// Calculate averages
	if len(firstDragonTimes) > 0 {
		analysis.FirstDragonAvgTime = average(firstDragonTimes)
	}
	if totalGames > 0 {
		analysis.FirstDragonContestRate = float64(dragonContests) / float64(totalGames)
	}
	if len(firstTowerTimes) > 0 {
		analysis.FirstTowerAvgTime = average(firstTowerTimes)
	}
	if len(firstHeraldTimes) > 0 {
		analysis.FirstHeraldAvgTime = average(firstHeraldTimes)
	}

	// Determine most common first tower lane
	maxLane := "bot"
	maxCount := 0
	for lane, count := range towerLaneCounts {
		if count > maxCount {
			maxCount = count
			maxLane = lane
		}
	}
	analysis.FirstTowerLane = maxLane

	// Determine herald usage pattern from actual data
	maxHeraldLane := "mid"
	maxHeraldCount := 0
	for lane, count := range heraldLaneCounts {
		if count > maxHeraldCount {
			maxHeraldCount = count
			maxHeraldLane = lane
		}
	}
	analysis.HeraldUsagePattern = maxHeraldLane

	return analysis
}

// AnalyzeDragonPriority analyzes which dragon types the team prioritizes
func (a *TimingAnalyzerEngine) AnalyzeDragonPriority(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) map[string]float64 {
	dragonCounts := make(map[string]int)
	totalDragons := 0

	for _, eventData := range events {
		if eventData == nil {
			continue
		}

		for _, dragon := range eventData.DragonKills {
			if dragon.TeamID == teamID {
				totalDragons++
				if dragon.DragonType != "" && dragon.DragonType != "unknown" {
					dragonCounts[dragon.DragonType]++
				}
			}
		}
	}

	// Convert to percentages
	priority := make(map[string]float64)
	if totalDragons > 0 {
		for dragonType, count := range dragonCounts {
			priority[dragonType] = float64(count) / float64(totalDragons)
		}
	}

	return priority
}

// AnalyzeFirstBloodPatterns analyzes first blood timing and location patterns
func (a *TimingAnalyzerEngine) AnalyzeFirstBloodPatterns(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) *FirstBloodAnalysis {
	analysis := &FirstBloodAnalysis{
		LocationBreakdown: make(map[string]float64),
	}

	var (
		totalGames      int
		firstBloods     int
		firstBloodTimes []float64
		locationCounts  = map[string]int{"top": 0, "mid": 0, "bot": 0, "jungle": 0}
	)

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			var ourTeam *grid.GameTeam
			for i := range game.Teams {
				if game.Teams[i].ID == teamID {
					ourTeam = &game.Teams[i]
					break
				}
			}

			if ourTeam == nil {
				continue
			}

			totalGames++

			// Find first blood
			for _, kill := range eventData.Kills {
				if kill.FirstBlood {
					// Check if our team got first blood
					if isPlayerOnTeam(kill.KillerID, ourTeam) {
						firstBloods++
						firstBloodTimes = append(firstBloodTimes, float64(kill.GameTime)/60000.0)
						
						// Determine location from position
						pos := kill.KillerPosition
						if pos == nil {
							pos = kill.VictimPosition
						}
						if pos != nil {
							lane := determineLaneFromPosition(pos)
							switch lane {
							case "top", "top_jungle":
								locationCounts["top"]++
							case "mid":
								locationCounts["mid"]++
							case "bot", "bot_jungle":
								locationCounts["bot"]++
							default:
								locationCounts["jungle"]++
							}
						}
					}
					break // Only one first blood per game
				}
			}
		}
	}

	if totalGames > 0 {
		analysis.Rate = float64(firstBloods) / float64(totalGames)
	}
	if len(firstBloodTimes) > 0 {
		analysis.AvgTime = average(firstBloodTimes)
	}
	
	totalLocations := 0
	for _, count := range locationCounts {
		totalLocations += count
	}
	if totalLocations > 0 {
		for loc, count := range locationCounts {
			analysis.LocationBreakdown[loc] = float64(count) / float64(totalLocations)
		}
	}

	return analysis
}

// FirstBloodAnalysis contains first blood pattern data
type FirstBloodAnalysis struct {
	Rate              float64            `json:"rate"`
	AvgTime           float64            `json:"avgTime"` // minutes
	LocationBreakdown map[string]float64 `json:"locationBreakdown"`
}

// GenerateTimingInsights generates hackathon-format insights from timing data
// ENHANCED: Now produces specific, data-backed insights in hackathon format
func (a *TimingAnalyzerEngine) GenerateTimingInsights(
	junglePathing *JunglePathingAnalysis,
	objectiveTimings *ObjectiveTimingAnalysis,
) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	// Jungle pathing insight - HACKATHON FORMAT
	// "Their jungler paths to bot 75% of the time pre-10 mins"
	if junglePathing.PreSixBotRate > 0.5 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Their jungler paths to bot %.0f%% of the time pre-10 mins",
				junglePathing.PreSixBotRate*100),
			Metric:  "jungle_pathing",
			Value:   junglePathing.PreSixBotRate,
			Context: "early game",
		})
	} else if junglePathing.PreSixTopRate > 0.5 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Their jungler paths to top %.0f%% of the time pre-10 mins",
				junglePathing.PreSixTopRate*100),
			Metric:  "jungle_pathing",
			Value:   junglePathing.PreSixTopRate,
			Context: "early game",
		})
	}
	
	// Gank distribution insight
	if len(junglePathing.GanksByLane) > 0 {
		// Find most ganked lane
		maxLane := ""
		maxRate := 0.0
		for lane, rate := range junglePathing.GanksByLane {
			if rate > maxRate {
				maxRate = rate
				maxLane = lane
			}
		}
		if maxRate > 0.4 { // Significant preference
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("%.0f%% of early ganks target %s lane",
					maxRate*100, maxLane),
				Metric:  "gank_distribution",
				Value:   maxRate,
				Context: maxLane,
			})
		}
	}

	// First tower timing insight - HACKATHON FORMAT
	// "4-man group for first tower push, usually in bot lane at ~13 mins"
	if objectiveTimings.FirstTowerAvgTime > 0 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("First tower typically falls at ~%.0f mins (usually %s lane)",
				objectiveTimings.FirstTowerAvgTime, objectiveTimings.FirstTowerLane),
			Metric:  "first_tower_timing",
			Value:   objectiveTimings.FirstTowerAvgTime,
			Context: objectiveTimings.FirstTowerLane,
		})
	}

	// Dragon priority insight - HACKATHON FORMAT
	// "Prioritizes first Drake (82% contest rate)"
	if objectiveTimings.FirstDragonContestRate > 0.6 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Prioritizes first Drake (%.0f%% contest rate)",
				objectiveTimings.FirstDragonContestRate*100),
			Metric:  "dragon_priority",
			Value:   objectiveTimings.FirstDragonContestRate,
			Context: "early game",
		})
	} else if objectiveTimings.FirstDragonContestRate < 0.4 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Low dragon priority (only %.0f%% first drake contest rate)",
				objectiveTimings.FirstDragonContestRate*100),
			Metric:  "dragon_priority",
			Value:   objectiveTimings.FirstDragonContestRate,
			Context: "early game",
		})
	}
	
	// Herald usage insight
	if objectiveTimings.FirstHeraldAvgTime > 0 {
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Takes first Herald at ~%.0f mins, typically uses it %s",
				objectiveTimings.FirstHeraldAvgTime, objectiveTimings.HeraldUsagePattern),
			Metric:  "herald_usage",
			Value:   objectiveTimings.FirstHeraldAvgTime,
			Context: objectiveTimings.HeraldUsagePattern,
		})
	}
	
	// Baron timing pattern
	if len(objectiveTimings.BaronAttemptTimings) >= 3 {
		avgBaronTime := average(objectiveTimings.BaronAttemptTimings)
		insights = append(insights, StrategyInsight{
			Text: fmt.Sprintf("Baron attempts typically around ~%.0f mins",
				avgBaronTime),
			Metric:  "baron_timing",
			Value:   avgBaronTime,
			Context: "late game",
		})
	}

	return insights
}

// GenerateTimingCounterStrategies generates counter-strategies based on timing patterns
// ENHANCED: Now produces specific, actionable recommendations in hackathon format
func (a *TimingAnalyzerEngine) GenerateTimingCounterStrategies(
	junglePathing *JunglePathingAnalysis,
	objectiveTimings *ObjectiveTimingAnalysis,
) []InGameStrategyInsight {
	strategies := make([]InGameStrategyInsight, 0)

	// Counter jungle pathing - HACKATHON FORMAT
	// "Recommend aggressive counter-jungling on their top side"
	if junglePathing.PreSixBotRate > 0.6 {
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Aggressive counter-jungling on their top side",
			Timing:   "Pre-10 minutes",
			Reason:   fmt.Sprintf("Their jungler paths to bot %.0f%% of the time, leaving top side vulnerable",
				junglePathing.PreSixBotRate*100),
			Impact: "HIGH",
		})
		
		// Additional strategy: ward their bot-side jungle
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Deep ward their bot-side jungle entrances",
			Timing:   "3-4 minutes",
			Reason:   "Track their predictable bot-side pathing for lane safety",
			Impact:   "MEDIUM",
		})
	} else if junglePathing.PreSixTopRate > 0.6 {
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Aggressive counter-jungling on their bot side",
			Timing:   "Pre-10 minutes",
			Reason:   fmt.Sprintf("Their jungler paths to top %.0f%% of the time, leaving bot side vulnerable",
				junglePathing.PreSixTopRate*100),
			Impact: "HIGH",
		})
		
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Deep ward their top-side jungle entrances",
			Timing:   "3-4 minutes",
			Reason:   "Track their predictable top-side pathing for lane safety",
			Impact:   "MEDIUM",
		})
	}
	
	// Gank lane counter-strategy
	if len(junglePathing.GanksByLane) > 0 {
		maxLane := ""
		maxRate := 0.0
		for lane, rate := range junglePathing.GanksByLane {
			if rate > maxRate {
				maxRate = rate
				maxLane = lane
			}
		}
		if maxRate > 0.5 {
			strategies = append(strategies, InGameStrategyInsight{
				Strategy: fmt.Sprintf("Play safe in %s lane, expect jungle pressure", maxLane),
				Timing:   "Pre-10 minutes",
				Reason:   fmt.Sprintf("%.0f%% of their early ganks target %s lane", maxRate*100, maxLane),
				Impact:   "HIGH",
			})
		}
	}

	// Dragon contest strategy
	if objectiveTimings.FirstDragonContestRate < 0.5 {
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Prioritize dragon control - set up vision and contest every spawn",
			Timing:   "5-25 minutes",
			Reason:   fmt.Sprintf("Opponent only contests first dragon %.0f%% of games - free drakes available",
				objectiveTimings.FirstDragonContestRate*100),
			Impact: "HIGH",
		})
	} else if objectiveTimings.FirstDragonContestRate > 0.8 {
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: "Trade objectives - take Herald while they focus dragon",
			Timing:   "5-14 minutes",
			Reason:   fmt.Sprintf("They heavily prioritize dragon (%.0f%% contest rate) - Herald often free",
				objectiveTimings.FirstDragonContestRate*100),
			Impact: "MEDIUM",
		})
	}

	// First tower strategy
	if objectiveTimings.FirstTowerLane != "" && objectiveTimings.FirstTowerAvgTime > 0 {
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: fmt.Sprintf("Defend %s tower aggressively around %.0f minutes",
				objectiveTimings.FirstTowerLane, objectiveTimings.FirstTowerAvgTime),
			Timing:   fmt.Sprintf("~%.0f minutes", objectiveTimings.FirstTowerAvgTime),
			Reason:   fmt.Sprintf("Opponent typically takes first tower in %s lane at this timing",
				objectiveTimings.FirstTowerLane),
			Impact:   "MEDIUM",
		})
	}
	
	// Baron timing strategy
	if len(objectiveTimings.BaronAttemptTimings) >= 2 {
		avgBaronTime := average(objectiveTimings.BaronAttemptTimings)
		strategies = append(strategies, InGameStrategyInsight{
			Strategy: fmt.Sprintf("Set up Baron vision and control around %.0f minutes",
				avgBaronTime-2),
			Timing:   fmt.Sprintf("%.0f-%.0f minutes", avgBaronTime-2, avgBaronTime),
			Reason:   fmt.Sprintf("Opponent typically attempts Baron around %.0f minutes",
				avgBaronTime),
			Impact:   "HIGH",
		})
	}

	return strategies
}

// Helper functions

func isJungleChampion(champion string) bool {
	jungleChampions := map[string]bool{
		"Lee Sin": true, "Elise": true, "Nidalee": true, "Graves": true,
		"Kindred": true, "Rek'Sai": true, "Jarvan IV": true, "Vi": true,
		"Xin Zhao": true, "Nocturne": true, "Hecarim": true, "Sejuani": true,
		"Zac": true, "Amumu": true, "Rammus": true, "Volibear": true,
		"Warwick": true, "Trundle": true, "Olaf": true, "Udyr": true,
		"Kha'Zix": true, "Rengar": true, "Evelynn": true, "Shaco": true,
		"Ekko": true, "Diana": true, "Viego": true, "Lillia": true,
		"Fiddlesticks": true, "Karthus": true, "Taliyah": true,
		"Wukong": true, "Poppy": true, "Maokai": true, "Ivern": true,
		"Bel'Veth": true, "Briar": true,
	}
	return jungleChampions[champion]
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Sort helper for timing data
func sortTimings(timings []float64) {
	sort.Float64s(timings)
}

// =============================================================================
// ENHANCED STATISTICAL TIMING ANALYSIS
// =============================================================================

// EnhancedObjectiveAnalysis contains percentile-based timing analysis
type EnhancedObjectiveAnalysis struct {
	FirstDragon *TimingStats `json:"firstDragon"`
	FirstTower  *TimingStats `json:"firstTower"`
	FirstHerald *TimingStats `json:"firstHerald"`
	Baron       *TimingStats `json:"baron"`
}

// TimingStats contains distribution and insights for a timing metric
type TimingStats struct {
	Distribution *Distribution    `json:"distribution"`
	Insights     []TimingInsight  `json:"insights"`
}

// TimingInsight represents a timing-based insight with context
type TimingInsight struct {
	Text        string  `json:"text"`
	Metric      string  `json:"metric"`
	Value       float64 `json:"value"`
	Percentile  int     `json:"percentile"`
	SampleSize  int     `json:"sampleSize"`
	Confidence  float64 `json:"confidence"`
	Context     string  `json:"context"` // "early", "average", "late"
}

// AnalyzeObjectiveTimingsEnhanced provides percentile-based timing analysis
func (a *TimingAnalyzerEngine) AnalyzeObjectiveTimingsEnhanced(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.LoLEventData,
) *EnhancedObjectiveAnalysis {
	// Collect all timing data
	var firstDragonTimes []float64
	var firstTowerTimes []float64
	var firstHeraldTimes []float64
	var baronTimes []float64

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			var ourTeam *grid.GameTeam
			for i := range game.Teams {
				if game.Teams[i].ID == teamID {
					ourTeam = &game.Teams[i]
					break
				}
			}
			if ourTeam == nil {
				continue
			}

			// Track first dragon timing
			gameHadDragon := false
			for _, dragon := range eventData.DragonKills {
				if !gameHadDragon && (dragon.TeamID == teamID || isTeamID(dragon.TeamID, ourTeam)) {
					if IsValidGameTime(dragon.GameTime) {
						timeMinutes := float64(dragon.GameTime) / 60000.0
						firstDragonTimes = append(firstDragonTimes, timeMinutes)
						gameHadDragon = true
					}
				}
			}

			// Track first tower timing
			gameHadTower := false
			for _, tower := range eventData.TowerDestroys {
				if !gameHadTower && (tower.TeamID == teamID || isTeamID(tower.TeamID, ourTeam)) {
					if IsValidGameTime(tower.GameTime) {
						timeMinutes := float64(tower.GameTime) / 60000.0
						firstTowerTimes = append(firstTowerTimes, timeMinutes)
						gameHadTower = true
					}
				}
			}

			// Track first herald timing
			gameHadHerald := false
			for _, herald := range eventData.HeraldKills {
				if !gameHadHerald && (herald.TeamID == teamID || isTeamID(herald.TeamID, ourTeam)) {
					if IsValidGameTime(herald.GameTime) {
						timeMinutes := float64(herald.GameTime) / 60000.0
						firstHeraldTimes = append(firstHeraldTimes, timeMinutes)
						gameHadHerald = true
					}
				}
			}

			// Track baron timings
			for _, baron := range eventData.BaronKills {
				if baron.TeamID == teamID || isTeamID(baron.TeamID, ourTeam) {
					if IsValidGameTime(baron.GameTime) {
						timeMinutes := float64(baron.GameTime) / 60000.0
						baronTimes = append(baronTimes, timeMinutes)
					}
				}
			}
		}
	}

	// Calculate distributions
	dragonDist := a.stats.CalculateDistribution(firstDragonTimes)
	towerDist := a.stats.CalculateDistribution(firstTowerTimes)
	heraldDist := a.stats.CalculateDistribution(firstHeraldTimes)
	baronDist := a.stats.CalculateDistribution(baronTimes)

	return &EnhancedObjectiveAnalysis{
		FirstDragon: &TimingStats{
			Distribution: dragonDist,
			Insights:     a.generateTimingInsightsFromDist("first dragon", dragonDist),
		},
		FirstTower: &TimingStats{
			Distribution: towerDist,
			Insights:     a.generateTimingInsightsFromDist("first tower", towerDist),
		},
		FirstHerald: &TimingStats{
			Distribution: heraldDist,
			Insights:     a.generateTimingInsightsFromDist("first herald", heraldDist),
		},
		Baron: &TimingStats{
			Distribution: baronDist,
			Insights:     a.generateTimingInsightsFromDist("baron", baronDist),
		},
	}
}

// generateTimingInsightsFromDist creates contextual insights from distribution
func (a *TimingAnalyzerEngine) generateTimingInsightsFromDist(metric string, dist *Distribution) []TimingInsight {
	insights := make([]TimingInsight, 0)

	if dist.SampleSize < 3 {
		insights = append(insights, TimingInsight{
			Text:       fmt.Sprintf("Insufficient data for %s analysis (n=%d)", metric, dist.SampleSize),
			Metric:     metric,
			SampleSize: dist.SampleSize,
			Confidence: 0.0,
		})
		return insights
	}

	// Generate insight based on distribution
	p50 := dist.Percentiles[50]
	confidence := a.stats.CalculateConfidence(dist)

	insights = append(insights, TimingInsight{
		Text: fmt.Sprintf("Takes %s at ~%.1f minutes (median, n=%d)",
			metric, p50, dist.SampleSize),
		Metric:     metric,
		Value:      p50,
		Percentile: 50,
		SampleSize: dist.SampleSize,
		Confidence: confidence,
		Context:    "median",
	})

	// Add range insight
	p25 := dist.Percentiles[25]
	p75 := dist.Percentiles[75]
	insights = append(insights, TimingInsight{
		Text: fmt.Sprintf("%s timing range: %.1f-%.1f minutes (25th-75th percentile)",
			metric, p25, p75),
		Metric:     metric,
		Value:      p75 - p25,
		SampleSize: dist.SampleSize,
		Confidence: confidence,
		Context:    "range",
	})

	// Add variance insight if significant
	if dist.StdDev > 2.0 {
		insights = append(insights, TimingInsight{
			Text: fmt.Sprintf("High variance in %s timing (±%.1f min)",
				metric, dist.StdDev),
			Metric:     metric,
			Value:      dist.StdDev,
			SampleSize: dist.SampleSize,
			Confidence: confidence,
			Context:    "variance",
		})
	}

	return insights
}

// ClassifyTimingValue determines if a timing is early/average/late
func (a *TimingAnalyzerEngine) ClassifyTimingValue(value float64, dist *Distribution) string {
	return a.stats.ClassifyTiming(value, dist)
}

// CompareTimingToBaseline compares a team's timing to a baseline distribution
func (a *TimingAnalyzerEngine) CompareTimingToBaseline(teamDist, baselineDist *Distribution) *DistributionComparison {
	return a.stats.CompareDistributions(teamDist, baselineDist)
}
