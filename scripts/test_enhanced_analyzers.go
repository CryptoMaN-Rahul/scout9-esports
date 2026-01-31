// +build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"scout9/pkg/grid"
	"scout9/pkg/intelligence"
)

func main() {
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		apiKey = "hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO"
	}

	client := grid.NewClient(apiKey, nil)
	ctx := context.Background()

	fmt.Println("=== ENHANCED ANALYZER TEST ===")
	fmt.Println()

	// Test with LoL series
	fmt.Println("--- Testing LoL Analysis (Series 2692648) ---")
	testLoLAnalysis(ctx, client, "2692648")

	fmt.Println()
	fmt.Println("--- Testing VALORANT Analysis (Series 2629390) ---")
	testVALAnalysis(ctx, client, "2629390")
}

func testLoLAnalysis(ctx context.Context, client *grid.Client, seriesID string) {
	// Download events
	wrappers, err := client.DownloadEvents(ctx, seriesID)
	if err != nil {
		fmt.Printf("Error downloading events: %v\n", err)
		return
	}

	// Parse LoL events
	eventData, err := grid.ParseLoLEvents(wrappers)
	if err != nil {
		fmt.Printf("Error parsing events: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d kills, %d dragons, %d towers\n",
		len(eventData.Kills), len(eventData.DragonKills), len(eventData.TowerDestroys))

	// Get series state for team info
	seriesState, err := client.GetSeriesState(ctx, seriesID)
	if err != nil {
		fmt.Printf("Error getting series state: %v\n", err)
		return
	}

	if len(seriesState.Teams) < 2 {
		fmt.Println("Not enough teams in series")
		return
	}

	teamID := seriesState.Teams[0].ID
	teamName := seriesState.Teams[0].Name
	fmt.Printf("Analyzing team: %s (%s)\n\n", teamName, teamID)

	// Create analyzers
	timingAnalyzer := intelligence.NewTimingAnalyzer()
	matchupAnalyzer := intelligence.NewMatchupAnalyzer()

	// Prepare data
	seriesStates := []*grid.SeriesState{seriesState}
	events := map[string]*grid.LoLEventData{seriesID: eventData}

	// Test Jungle Pathing Analysis
	fmt.Println("=== JUNGLE PATHING ANALYSIS ===")
	junglePathing := timingAnalyzer.AnalyzeJunglePathing(teamID, seriesStates, events)
	fmt.Printf("Path Preference: %s\n", junglePathing.PreSixPathPreference)
	fmt.Printf("Bot Rate: %.1f%%\n", junglePathing.PreSixBotRate*100)
	fmt.Printf("Top Rate: %.1f%%\n", junglePathing.PreSixTopRate*100)
	fmt.Printf("Gank Distribution:\n")
	for lane, rate := range junglePathing.GanksByLane {
		fmt.Printf("  %s: %.1f%%\n", lane, rate*100)
	}

	// Test Objective Timing Analysis
	fmt.Println("\n=== OBJECTIVE TIMING ANALYSIS ===")
	objectiveTimings := timingAnalyzer.AnalyzeObjectiveTimings(teamID, seriesStates, events)
	fmt.Printf("First Dragon Avg Time: %.1f min\n", objectiveTimings.FirstDragonAvgTime)
	fmt.Printf("First Dragon Contest Rate: %.1f%%\n", objectiveTimings.FirstDragonContestRate*100)
	fmt.Printf("First Tower Avg Time: %.1f min\n", objectiveTimings.FirstTowerAvgTime)
	fmt.Printf("First Tower Lane: %s\n", objectiveTimings.FirstTowerLane)
	fmt.Printf("Herald Usage Pattern: %s\n", objectiveTimings.HeraldUsagePattern)

	// Test Timing Insights Generation
	fmt.Println("\n=== TIMING INSIGHTS (HACKATHON FORMAT) ===")
	insights := timingAnalyzer.GenerateTimingInsights(junglePathing, objectiveTimings)
	for _, insight := range insights {
		fmt.Printf("• %s\n", insight.Text)
	}

	// Test Counter Strategies
	fmt.Println("\n=== COUNTER STRATEGIES ===")
	strategies := timingAnalyzer.GenerateTimingCounterStrategies(junglePathing, objectiveTimings)
	for _, strategy := range strategies {
		fmt.Printf("• [%s] %s\n", strategy.Impact, strategy.Strategy)
		fmt.Printf("  Timing: %s\n", strategy.Timing)
		fmt.Printf("  Reason: %s\n", strategy.Reason)
	}

	// Test Player Matchup Analysis
	fmt.Println("\n=== PLAYER MATCHUP ANALYSIS ===")
	if len(seriesState.Games) > 0 && len(seriesState.Games[0].Teams) > 0 {
		for _, player := range seriesState.Games[0].Teams[0].Players {
			classPerf := matchupAnalyzer.AnalyzeLoLPlayerClassPerformance(
				player.ID, player.Name, seriesStates,
			)
			if len(classPerf) > 0 {
				fmt.Printf("\nPlayer: %s\n", player.Name)
				for className, perf := range classPerf {
					fmt.Printf("  %s: %.1f KDA, %.0f%% WR (%d games)\n",
						className, perf.KDA, perf.WinRate*100, perf.GamesPlayed)
				}

				// Generate matchup insights
				matchupInsights := matchupAnalyzer.GenerateMatchupInsights(player.Name, classPerf)
				for _, insight := range matchupInsights {
					fmt.Printf("  → %s\n", insight.Text)
				}
			}
		}
	}
}

func testVALAnalysis(ctx context.Context, client *grid.Client, seriesID string) {
	// Download events
	wrappers, err := client.DownloadEvents(ctx, seriesID)
	if err != nil {
		fmt.Printf("Error downloading events: %v\n", err)
		return
	}

	// Parse VAL events
	eventData, err := grid.ParseVALEvents(wrappers)
	if err != nil {
		fmt.Printf("Error parsing events: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d kills, %d round ends, %d plants\n",
		len(eventData.Kills), len(eventData.RoundEnds), len(eventData.Plants))
	fmt.Printf("Map: %s\n", eventData.MapName)

	// Get series state for team info
	seriesState, err := client.GetSeriesState(ctx, seriesID)
	if err != nil {
		fmt.Printf("Error getting series state: %v\n", err)
		return
	}

	if len(seriesState.Teams) < 2 {
		fmt.Println("Not enough teams in series")
		return
	}

	teamID := seriesState.Teams[0].ID
	teamName := seriesState.Teams[0].Name
	fmt.Printf("Analyzing team: %s (%s)\n\n", teamName, teamID)

	// Create site analyzer
	siteAnalyzer := intelligence.NewSiteAnalyzer()

	// Prepare data
	seriesStates := []*grid.SeriesState{seriesState}
	events := map[string]*grid.VALEventData{seriesID: eventData}

	// Test Site Pattern Analysis
	fmt.Println("=== SITE PATTERN ANALYSIS ===")
	siteAnalysis := siteAnalyzer.AnalyzeSitePatterns(teamID, seriesStates, events)
	for mapName, analysis := range siteAnalysis {
		fmt.Printf("\nMap: %s\n", mapName)
		for siteName, stats := range analysis.Sites {
			if stats.AttackAttempts > 0 || stats.DefenseAttempts > 0 {
				fmt.Printf("  %s-Site:\n", siteName)
				fmt.Printf("    Attack: %d attempts, %.0f%% success\n",
					stats.AttackAttempts, stats.AttackWinRate*100)
				fmt.Printf("    Defense: %d attempts, %.0f%% success\n",
					stats.DefenseAttempts, stats.DefenseWinRate*100)
			}
		}
	}

	// Test Attack Pattern Analysis
	fmt.Println("\n=== ATTACK PATTERNS ===")
	attackPatterns := siteAnalyzer.AnalyzeAttackPatterns(teamID, seriesStates, events)
	for _, pattern := range attackPatterns {
		fmt.Printf("• %s: %.0f%% frequency, %.0f%% success\n",
			pattern.Description, pattern.Frequency*100, pattern.SuccessRate*100)
	}

	// Test Pistol Round Pattern Analysis (HACKATHON FORMAT)
	fmt.Println("\n=== PISTOL ROUND PATTERNS (HACKATHON FORMAT) ===")
	pistolInsights := siteAnalyzer.AnalyzePistolRoundPatterns(teamID, seriesStates, events)
	for _, insight := range pistolInsights {
		fmt.Printf("• %s\n", insight.Text)
	}

	// Test Site Insights
	fmt.Println("\n=== SITE INSIGHTS ===")
	siteInsights := siteAnalyzer.GenerateSiteInsights(siteAnalysis, attackPatterns)
	for _, insight := range siteInsights {
		fmt.Printf("• %s\n", insight.Text)
	}

	// Test Defense Insights
	fmt.Println("\n=== DEFENSE INSIGHTS ===")
	defenseInsights := siteAnalyzer.GenerateDefenseInsights(siteAnalysis)
	for _, insight := range defenseInsights {
		fmt.Printf("• %s\n", insight.Text)
	}

	// Test Counter Strategies
	fmt.Println("\n=== SITE COUNTER STRATEGIES ===")
	strategies := siteAnalyzer.GenerateSiteCounterStrategies(siteAnalysis)
	for _, strategy := range strategies {
		fmt.Printf("• [%s] %s\n", strategy.Impact, strategy.Strategy)
		fmt.Printf("  Reason: %s\n", strategy.Reason)
	}
}
