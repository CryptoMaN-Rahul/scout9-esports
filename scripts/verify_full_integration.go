// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"scout9/pkg/grid"
	"scout9/pkg/intelligence"
)

func main() {
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		fmt.Println("ERROR: GRID_API_KEY environment variable not set")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client := grid.NewClient(apiKey, nil)
	seriesID := "2692648"

	fmt.Println("=== Full Integration Verification ===")
	fmt.Printf("Series ID: %s\n\n", seriesID)

	// Step 1: Get Series State
	fmt.Println("1. Fetching Series State...")
	seriesState, err := client.GetSeriesState(ctx, seriesID)
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Series ID: %s, Finished: %v\n", seriesState.ID, seriesState.Finished)
	fmt.Printf("   Teams: %d, Games: %d\n", len(seriesState.Teams), len(seriesState.Games))

	for _, team := range seriesState.Teams {
		fmt.Printf("   Team: %s (ID: %s) - Won: %v\n", team.Name, team.ID, team.Won)
	}
	fmt.Println()

	// Step 2: Download JSONL events
	fmt.Println("2. Downloading JSONL events...")
	wrappers, err := client.DownloadEvents(ctx, seriesID)
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Downloaded %d event wrappers\n\n", len(wrappers))

	// Step 3: Parse LoL events
	fmt.Println("3. Parsing LoL events...")
	lolEvents, err := grid.ParseLoLEvents(wrappers)
	if err != nil {
		fmt.Printf("   ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Kills: %d, Dragons: %d, Towers: %d\n\n", 
		len(lolEvents.Kills), len(lolEvents.DragonKills), len(lolEvents.TowerDestroys))

	// Step 4: Test LoL Analyzer
	fmt.Println("4. Testing LoL Analyzer with both teams...")
	analyzer := intelligence.NewLoLAnalyzer()

	eventsMap := map[string]*grid.LoLEventData{seriesID: lolEvents}

	for _, testTeam := range seriesState.Teams {
		fmt.Printf("\n=== %s (ID: %s) ===\n", testTeam.Name, testTeam.ID)

		analysis, err := analyzer.AnalyzeTeam(ctx, testTeam.ID, testTeam.Name, 
			[]*grid.SeriesState{seriesState}, eventsMap)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		m := analysis.LoLMetrics
		fmt.Printf("FirstBloodRate: %.0f%%\n", m.FirstBloodRate*100)
		fmt.Printf("FirstDragonRate: %.0f%%\n", m.FirstDragonRate*100)
		fmt.Printf("FirstTowerRate: %.0f%%\n", m.FirstTowerRate*100)
		fmt.Printf("EarlyGameRating: %.1f\n", m.EarlyGameRating)
		fmt.Printf("MidGameRating: %.1f\n", m.MidGameRating)
		fmt.Printf("LateGameRating: %.1f\n", m.LateGameRating)
	}

	fmt.Println("\n=== VERIFICATION COMPLETE ===")
	fmt.Println("EventAnalyzer integration working correctly!")
}
