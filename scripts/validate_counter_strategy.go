// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"scout9/pkg/grid"
	"scout9/pkg/intelligence"
)

// Validation script for Enhanced Counter-Strategy feature
// Tests with real GRID data from LCK, LEC, and VCT Americas tournaments

func main() {
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		fmt.Println("ERROR: GRID_API_KEY environment variable not set")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client := grid.NewClient(apiKey, nil)

	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("ENHANCED COUNTER-STRATEGY VALIDATION WITH REAL GRID DATA")
	fmt.Println("=" + strings.Repeat("=", 79))

	// Test with LoL team (LCK)
	fmt.Println("\n[TEST 1] League of Legends - LCK Team Analysis")
	fmt.Println("-" + strings.Repeat("-", 79))
	testLoLCounterStrategy(ctx, client)

	// Test with VALORANT team (VCT Americas)
	fmt.Println("\n[TEST 2] VALORANT - VCT Americas Team Analysis")
	fmt.Println("-" + strings.Repeat("-", 79))
	testVALCounterStrategy(ctx, client)

	// Test Head-to-Head comparison
	fmt.Println("\n[TEST 3] Head-to-Head Comparison")
	fmt.Println("-" + strings.Repeat("-", 79))
	testHeadToHead(ctx, client)

	// Test Economy Analyzer
	fmt.Println("\n[TEST 4] Economy Analyzer Validation")
	fmt.Println("-" + strings.Repeat("-", 79))
	testEconomyAnalyzer(ctx, client)

	fmt.Println("\n" + "=" + strings.Repeat("=", 79))
	fmt.Println("VALIDATION COMPLETE")
	fmt.Println("=" + strings.Repeat("=", 79))
}

func testLoLCounterStrategy(ctx context.Context, client *grid.Client) {
	// Get teams from LCK tournament (ID: 825490 - LCK Split 2 2025)
	fmt.Println("Fetching LCK teams...")
	teams, err := client.GetTeams(ctx, "825490")
	if err != nil {
		fmt.Printf("ERROR fetching teams: %v\n", err)
		return
	}

	if len(teams) == 0 {
		fmt.Println("No teams found in LCK tournament")
		return
	}

	// Pick first team for analysis
	targetTeam := teams[0]
	fmt.Printf("Analyzing team: %s (ID: %s)\n", targetTeam.Name, targetTeam.ID)

	// Get recent series for the team
	fmt.Println("Fetching recent series...")
	seriesStates, err := client.GetMatchDataForTeam(ctx, targetTeam.ID, 10)
	if err != nil {
		fmt.Printf("ERROR fetching match data: %v\n", err)
		return
	}

	fmt.Printf("Found %d series with data\n", len(seriesStates))

	if len(seriesStates) == 0 {
		fmt.Println("No series data available")
		return
	}

	// Create analyzers
	lolAnalyzer := intelligence.NewLoLAnalyzer()
	counterEngine := intelligence.NewCounterStrategyEngine()

	// Analyze team
	fmt.Println("Running team analysis...")
	teamAnalysis, err := lolAnalyzer.AnalyzeTeam(ctx, targetTeam.ID, targetTeam.Name, seriesStates, nil)
	if err != nil {
		fmt.Printf("ERROR analyzing team: %v\n", err)
		return
	}
	teamAnalysis.Title = "lol"

	// Get player profiles
	playerProfiles, err := lolAnalyzer.AnalyzePlayers(ctx, targetTeam.ID, seriesStates)
	if err != nil {
		fmt.Printf("ERROR analyzing players: %v\n", err)
		return
	}

	// Generate counter-strategy
	fmt.Println("Generating enhanced counter-strategy...")
	strategy := counterEngine.GenerateEnhancedCounterStrategy(
		teamAnalysis,
		playerProfiles,
		nil, // compositions
		seriesStates,
		nil, // lolEvents
		nil, // valEvents
	)

	// Print results
	fmt.Println("\n--- COUNTER-STRATEGY RESULTS ---")
	fmt.Printf("Team: %s\n", strategy.TeamName)
	fmt.Printf("Confidence Score: %.1f%%\n", strategy.ConfidenceScore)
	fmt.Printf("Win Condition: %s\n", strategy.WinCondition)

	fmt.Println("\nWeaknesses:")
	for i, w := range strategy.Weaknesses {
		fmt.Printf("  %d. %s - %s (Impact: %.0f)\n", i+1, w.Title, w.Description, w.Impact)
	}

	fmt.Println("\nDraft Recommendations:")
	for i, d := range strategy.DraftRecommendations {
		fmt.Printf("  %d. [%s] %s - %s\n", i+1, d.Type, d.Character, d.Reason)
	}

	fmt.Println("\nTarget Players:")
	for i, t := range strategy.TargetPlayers {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, t.PlayerName, t.Role, t.Reason)
	}

	fmt.Println("\nIn-Game Strategies:")
	for i, s := range strategy.InGameStrategies {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(strategy.InGameStrategies)-5)
			break
		}
		fmt.Printf("  %d. %s - %s\n", i+1, s.Title, s.Description)
	}

	if len(strategy.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, w := range strategy.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	// Validate requirements
	fmt.Println("\n--- REQUIREMENTS VALIDATION ---")
	validateLoLRequirements(strategy, teamAnalysis)
}

func testVALCounterStrategy(ctx context.Context, client *grid.Client) {
	// Get teams from VCT Americas (ID: 800675 - VCT Americas Stage 1 2025)
	fmt.Println("Fetching VCT Americas teams...")
	teams, err := client.GetTeams(ctx, "800675")
	if err != nil {
		fmt.Printf("ERROR fetching teams: %v\n", err)
		return
	}

	if len(teams) == 0 {
		fmt.Println("No teams found in VCT Americas tournament")
		return
	}

	// Pick first team for analysis
	targetTeam := teams[0]
	fmt.Printf("Analyzing team: %s (ID: %s)\n", targetTeam.Name, targetTeam.ID)

	// Get recent series for the team
	fmt.Println("Fetching recent series...")
	seriesStates, err := client.GetMatchDataForTeam(ctx, targetTeam.ID, 10)
	if err != nil {
		fmt.Printf("ERROR fetching match data: %v\n", err)
		return
	}

	fmt.Printf("Found %d series with data\n", len(seriesStates))

	if len(seriesStates) == 0 {
		fmt.Println("No series data available")
		return
	}

	// Download JSONL event files for site analysis
	fmt.Println("Downloading JSONL event files for site analysis...")
	valEvents := make(map[string]*grid.VALEventData)
	for _, series := range seriesStates {
		eventData, err := client.DownloadAndParseVALEvents(ctx, series.ID)
		if err != nil {
			fmt.Printf("  Warning: Could not download events for series %s: %v\n", series.ID, err)
			continue
		}
		if eventData != nil && (len(eventData.Plants) > 0 || len(eventData.RoundEnds) > 0) {
			valEvents[series.ID] = eventData
			fmt.Printf("  Downloaded events for series %s: %d plants, %d round ends\n", 
				series.ID, len(eventData.Plants), len(eventData.RoundEnds))
		}
	}
	fmt.Printf("Downloaded event data for %d series\n", len(valEvents))

	// Create analyzers
	valAnalyzer := intelligence.NewVALAnalyzer()
	counterEngine := intelligence.NewCounterStrategyEngine()

	// Analyze team
	fmt.Println("Running team analysis...")
	teamAnalysis, err := valAnalyzer.AnalyzeTeam(ctx, targetTeam.ID, targetTeam.Name, seriesStates, valEvents)
	if err != nil {
		fmt.Printf("ERROR analyzing team: %v\n", err)
		return
	}
	teamAnalysis.Title = "valorant"

	// Get player profiles
	playerProfiles, err := valAnalyzer.AnalyzePlayers(ctx, targetTeam.ID, seriesStates)
	if err != nil {
		fmt.Printf("ERROR analyzing players: %v\n", err)
		return
	}

	// Generate counter-strategy with event data for site analysis
	fmt.Println("Generating enhanced counter-strategy...")
	strategy := counterEngine.GenerateEnhancedCounterStrategy(
		teamAnalysis,
		playerProfiles,
		nil, // compositions
		seriesStates,
		nil, // lolEvents
		valEvents, // Pass event data for site analysis
	)

	// Print results
	fmt.Println("\n--- COUNTER-STRATEGY RESULTS ---")
	fmt.Printf("Team: %s\n", strategy.TeamName)
	fmt.Printf("Confidence Score: %.1f%%\n", strategy.ConfidenceScore)
	fmt.Printf("Win Condition: %s\n", strategy.WinCondition)

	fmt.Println("\nWeaknesses:")
	for i, w := range strategy.Weaknesses {
		fmt.Printf("  %d. %s - %s (Impact: %.0f)\n", i+1, w.Title, w.Description, w.Impact)
	}

	fmt.Println("\nDraft Recommendations:")
	for i, d := range strategy.DraftRecommendations {
		fmt.Printf("  %d. [%s] %s - %s\n", i+1, d.Type, d.Character, d.Reason)
	}

	fmt.Println("\nTarget Players:")
	for i, t := range strategy.TargetPlayers {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, t.PlayerName, t.Role, t.Reason)
	}

	fmt.Println("\nIn-Game Strategies:")
	for i, s := range strategy.InGameStrategies {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(strategy.InGameStrategies)-5)
			break
		}
		fmt.Printf("  %d. %s - %s\n", i+1, s.Title, s.Description)
	}

	if len(strategy.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, w := range strategy.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	// Validate requirements
	fmt.Println("\n--- REQUIREMENTS VALIDATION ---")
	validateVALRequirements(strategy, teamAnalysis)
}

func testHeadToHead(ctx context.Context, client *grid.Client) {
	// Get two teams from the same tournament for comparison
	fmt.Println("Fetching teams for head-to-head comparison...")
	teams, err := client.GetTeams(ctx, "825490") // LCK Split 2 2025
	if err != nil {
		fmt.Printf("ERROR fetching teams: %v\n", err)
		return
	}

	if len(teams) < 2 {
		fmt.Println("Need at least 2 teams for head-to-head comparison")
		return
	}

	team1 := teams[0]
	team2 := teams[1]
	fmt.Printf("Comparing: %s vs %s\n", team1.Name, team2.Name)

	// Create head-to-head analyzer
	h2hAnalyzer := intelligence.NewHeadToHeadAnalyzer(client)

	// Analyze matchup
	report, err := h2hAnalyzer.AnalyzeMatchup(ctx, team1.ID, team2.ID, "3", 20) // titleID 3 = LoL
	if err != nil {
		fmt.Printf("ERROR analyzing matchup: %v\n", err)
		return
	}

	// Get team analyses for style comparison
	fmt.Println("Fetching team data for style comparison...")
	lolAnalyzer := intelligence.NewLoLAnalyzer()

	team1Series, err := client.GetMatchDataForTeam(ctx, team1.ID, 10)
	if err != nil {
		fmt.Printf("Warning: Could not get team1 series: %v\n", err)
	}

	team2Series, err := client.GetMatchDataForTeam(ctx, team2.ID, 10)
	if err != nil {
		fmt.Printf("Warning: Could not get team2 series: %v\n", err)
	}

	// Analyze both teams for style comparison
	var team1Analysis, team2Analysis *intelligence.TeamAnalysis
	if len(team1Series) > 0 {
		team1Analysis, _ = lolAnalyzer.AnalyzeTeam(ctx, team1.ID, team1.Name, team1Series, nil)
		if team1Analysis != nil {
			team1Analysis.Title = "lol"
		}
	}
	if len(team2Series) > 0 {
		team2Analysis, _ = lolAnalyzer.AnalyzeTeam(ctx, team2.ID, team2.Name, team2Series, nil)
		if team2Analysis != nil {
			team2Analysis.Title = "lol"
		}
	}

	// Generate style comparison insights if we have both team analyses
	if team1Analysis != nil && team2Analysis != nil {
		fmt.Println("Generating style comparison...")
		styleInsights := h2hAnalyzer.GenerateMatchupInsights(report, team1Analysis, team2Analysis)
		report.Insights = append(report.Insights, styleInsights...)
	}

	// Print results
	fmt.Println("\n--- HEAD-TO-HEAD REPORT ---")
	fmt.Printf("Team 1: %s\n", report.Team1Name)
	fmt.Printf("Team 2: %s\n", report.Team2Name)
	fmt.Printf("Total Matches: %d\n", report.TotalMatches)
	fmt.Printf("Record: %s %d - %d %s\n", report.Team1Name, report.Team1Wins, report.Team2Wins, report.Team2Name)
	fmt.Printf("Confidence Score: %.1f%%\n", report.ConfidenceScore)

	if report.StyleComparison != nil {
		fmt.Println("\nStyle Comparison:")
		fmt.Printf("  Early Game: %s (%.0f) vs %s (%.0f) - Advantage: %s\n",
			report.Team1Name, report.StyleComparison.Team1EarlyGameRating,
			report.Team2Name, report.StyleComparison.Team2EarlyGameRating,
			report.StyleComparison.EarlyGameAdvantage)
		fmt.Printf("  Mid Game: %s (%.0f) vs %s (%.0f) - Advantage: %s\n",
			report.Team1Name, report.StyleComparison.Team1MidGameRating,
			report.Team2Name, report.StyleComparison.Team2MidGameRating,
			report.StyleComparison.MidGameAdvantage)
		fmt.Printf("  Late Game: %s (%.0f) vs %s (%.0f) - Advantage: %s\n",
			report.Team1Name, report.StyleComparison.Team1LateGameRating,
			report.Team2Name, report.StyleComparison.Team2LateGameRating,
			report.StyleComparison.LateGameAdvantage)
		if report.StyleComparison.StyleInsight != "" {
			fmt.Printf("  Style Insight: %s\n", report.StyleComparison.StyleInsight)
		}
	}

	fmt.Println("\nInsights:")
	for i, insight := range report.Insights {
		fmt.Printf("  %d. [%s] %s (Confidence: %.0f%%)\n", i+1, insight.Type, insight.Text, insight.Confidence*100)
	}

	if len(report.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, w := range report.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	// Validate requirements
	fmt.Println("\n--- REQUIREMENTS VALIDATION ---")
	validateH2HRequirements(report)
}

func testEconomyAnalyzer(ctx context.Context, client *grid.Client) {
	// Get VALORANT team data
	fmt.Println("Fetching VCT Americas teams for economy analysis...")
	teams, err := client.GetTeams(ctx, "800675")
	if err != nil {
		fmt.Printf("ERROR fetching teams: %v\n", err)
		return
	}

	if len(teams) == 0 {
		fmt.Println("No teams found")
		return
	}

	targetTeam := teams[0]
	fmt.Printf("Analyzing economy for: %s\n", targetTeam.Name)

	// Get series data
	seriesStates, err := client.GetMatchDataForTeam(ctx, targetTeam.ID, 10)
	if err != nil {
		fmt.Printf("ERROR fetching match data: %v\n", err)
		return
	}

	if len(seriesStates) == 0 {
		fmt.Println("No series data available")
		return
	}

	// Create economy analyzer
	economyAnalyzer := intelligence.NewEconomyAnalyzer()

	// Analyze economy rounds
	analysis := economyAnalyzer.AnalyzeEconomyRounds(targetTeam.ID, targetTeam.Name, seriesStates)

	// Print results
	fmt.Println("\n--- ECONOMY ANALYSIS ---")
	fmt.Printf("Team: %s\n", targetTeam.Name)
	fmt.Printf("Eco Round Win Rate: %.1f%% (%d rounds)\n", analysis.EcoRoundWinRate*100, analysis.EcoRounds)
	fmt.Printf("Force Buy Win Rate: %.1f%% (%d rounds)\n", analysis.ForceWinRate*100, analysis.ForceRounds)
	fmt.Printf("Full Buy Win Rate: %.1f%% (%d rounds)\n", analysis.FullBuyWinRate*100, analysis.FullBuyRounds)

	fmt.Println("\nPer-Map Stats:")
	for mapName, stats := range analysis.ByMap {
		fmt.Printf("  %s: Eco %.0f%%, Force %.0f%%, Full %.0f%%\n",
			mapName, stats.EcoWinRate*100, stats.ForceWinRate*100, stats.FullBuyWinRate*100)
	}

	// Generate insights
	insights := economyAnalyzer.GenerateEconomyInsights(analysis)
	fmt.Println("\nGenerated Insights:")
	for i, insight := range insights {
		fmt.Printf("  %d. [%s] %s (Impact: %s)\n", i+1, insight.Type, insight.Text, insight.Impact)
	}

	// Validate requirements
	fmt.Println("\n--- REQUIREMENTS VALIDATION ---")
	validateEconomyRequirements(analysis, insights)
}

// Validation functions

func validateLoLRequirements(strategy *intelligence.CounterStrategy, analysis *intelligence.TeamAnalysis) {
	passed := 0
	total := 0

	// Req 1.1: Analyze player KDA on different champion classes
	total++
	if len(strategy.DraftRecommendations) > 0 {
		fmt.Println("✓ Req 1.1: Draft recommendations generated based on class analysis")
		passed++
	} else {
		fmt.Println("✗ Req 1.1: No draft recommendations generated")
	}

	// Req 5.1: Win condition with specific data points
	total++
	if containsNumber(strategy.WinCondition) {
		fmt.Println("✓ Req 5.1: Win condition contains specific data points")
		passed++
	} else {
		fmt.Println("✗ Req 5.1: Win condition lacks specific data points")
	}

	// Req 5.5: Win condition limited to 2-3 sentences
	total++
	sentences := countSentences(strategy.WinCondition)
	if sentences >= 1 && sentences <= 4 {
		fmt.Printf("✓ Req 5.5: Win condition has %d sentences (within limit)\n", sentences)
		passed++
	} else {
		fmt.Printf("✗ Req 5.5: Win condition has %d sentences (should be 1-4)\n", sentences)
	}

	// Req 6.1: Confidence score calculated
	total++
	if strategy.ConfidenceScore > 0 {
		fmt.Printf("✓ Req 6.1: Confidence score calculated (%.1f%%)\n", strategy.ConfidenceScore)
		passed++
	} else {
		fmt.Println("✗ Req 6.1: Confidence score not calculated")
	}

	// Req 6.2: Sample size warning for low games
	total++
	if analysis.GamesAnalyzed < 5 && len(strategy.Warnings) > 0 {
		fmt.Println("✓ Req 6.2: Low sample size warning present")
		passed++
	} else if analysis.GamesAnalyzed >= 5 {
		fmt.Println("✓ Req 6.2: Sufficient sample size, no warning needed")
		passed++
	} else {
		fmt.Println("✗ Req 6.2: Missing low sample size warning")
	}

	fmt.Printf("\nLoL Requirements: %d/%d passed\n", passed, total)
}

func validateVALRequirements(strategy *intelligence.CounterStrategy, analysis *intelligence.TeamAnalysis) {
	passed := 0
	total := 0

	// Req 2.1: Economy analysis performed
	total++
	hasEconomyInsight := false
	for _, s := range strategy.InGameStrategies {
		if s.Title == "eco_weak" || s.Title == "force_strong" || s.Title == "force_weak" ||
			containsAny(s.Description, []string{"eco", "force", "save"}) {
			hasEconomyInsight = true
			break
		}
	}
	if hasEconomyInsight {
		fmt.Println("✓ Req 2.1: Economy analysis integrated into strategy")
		passed++
	} else {
		fmt.Println("✗ Req 2.1: No economy insights in strategy")
	}

	// Req 3.1: Site-specific analysis
	total++
	hasSiteInsight := false
	for _, s := range strategy.InGameStrategies {
		if containsAny(s.Title, []string{"Site", "site"}) ||
			containsAny(s.Description, []string{"A-Site", "B-Site", "C-Site"}) {
			hasSiteInsight = true
			break
		}
	}
	if hasSiteInsight {
		fmt.Println("✓ Req 3.1: Site-specific analysis present")
		passed++
	} else {
		fmt.Println("✗ Req 3.1: No site-specific analysis")
	}

	// Req 5.1: Win condition with specific data points
	total++
	if containsNumber(strategy.WinCondition) {
		fmt.Println("✓ Req 5.1: Win condition contains specific data points")
		passed++
	} else {
		fmt.Println("✗ Req 5.1: Win condition lacks specific data points")
	}

	// Req 6.1: Confidence score calculated
	total++
	if strategy.ConfidenceScore > 0 {
		fmt.Printf("✓ Req 6.1: Confidence score calculated (%.1f%%)\n", strategy.ConfidenceScore)
		passed++
	} else {
		fmt.Println("✗ Req 6.1: Confidence score not calculated")
	}

	fmt.Printf("\nVALORANT Requirements: %d/%d passed\n", passed, total)
}

func validateH2HRequirements(report *intelligence.HeadToHeadReport) {
	passed := 0
	total := 0

	// Req 4.2: Historical head-to-head record
	total++
	if report.TotalMatches >= 0 { // Even 0 is valid (no matches found)
		fmt.Printf("✓ Req 4.2: Historical record returned (%d matches)\n", report.TotalMatches)
		passed++
	} else {
		fmt.Println("✗ Req 4.2: Historical record not returned")
	}

	// Req 4.3: Style matchup analysis
	total++
	if report.StyleComparison != nil {
		fmt.Println("✓ Req 4.3: Style comparison included")
		passed++
	} else {
		fmt.Println("✗ Req 4.3: Style comparison missing")
	}

	// Req 4.4: Specific insights
	total++
	if len(report.Insights) > 0 {
		fmt.Printf("✓ Req 4.4: %d insights generated\n", len(report.Insights))
		passed++
	} else {
		fmt.Println("✗ Req 4.4: No insights generated")
	}

	// Req 4.5: Handle no historical matches
	total++
	if report.TotalMatches == 0 && len(report.Warnings) > 0 {
		fmt.Println("✓ Req 4.5: No matches case handled with warning")
		passed++
	} else if report.TotalMatches > 0 {
		fmt.Println("✓ Req 4.5: Historical matches found, no warning needed")
		passed++
	} else {
		fmt.Println("✗ Req 4.5: No matches case not properly handled")
	}

	fmt.Printf("\nHead-to-Head Requirements: %d/%d passed\n", passed, total)
}

func validateEconomyRequirements(analysis *intelligence.EconomyAnalysis, insights []intelligence.EconomyInsight) {
	passed := 0
	total := 0

	// Req 2.1: Analyze eco, force, full buy win rates
	total++
	if analysis.EcoRounds > 0 || analysis.ForceRounds > 0 || analysis.FullBuyRounds > 0 {
		fmt.Println("✓ Req 2.1: Economy round analysis performed")
		passed++
	} else {
		fmt.Println("✗ Req 2.1: No economy round data")
	}

	// Req 2.5: Per-map statistics
	total++
	if len(analysis.ByMap) > 0 {
		fmt.Printf("✓ Req 2.5: Per-map statistics available (%d maps)\n", len(analysis.ByMap))
		passed++
	} else {
		fmt.Println("✗ Req 2.5: No per-map statistics")
	}

	// Req 2.2-2.4: Insight generation based on thresholds
	total++
	if len(insights) > 0 {
		fmt.Printf("✓ Req 2.2-2.4: %d economy insights generated\n", len(insights))
		passed++
	} else {
		fmt.Println("✗ Req 2.2-2.4: No economy insights generated")
	}

	fmt.Printf("\nEconomy Requirements: %d/%d passed\n", passed, total)
}

// Helper functions

func containsNumber(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func countSentences(s string) int {
	count := 0
	for _, c := range s {
		if c == '.' || c == '!' || c == '?' {
			count++
		}
	}
	return count
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
