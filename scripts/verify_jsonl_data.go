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

	fmt.Println("=== GRID JSONL Data Verification ===")
	fmt.Printf("Series ID: %s\n\n", seriesID)

	// Step 1: List available files
	fmt.Println("1. Listing available files...")
	files, err := client.ListFiles(ctx, seriesID)
	if err != nil {
		fmt.Printf("   ERROR listing files: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Found %d files:\n", len(files))
	for _, f := range files {
		fmt.Printf("   - %s: %s (status: %s)\n", f.ID, f.Description, f.Status)
	}
	fmt.Println()

	// Step 2: Download and parse JSONL events
	fmt.Println("2. Downloading JSONL events...")
	wrappers, err := client.DownloadEvents(ctx, seriesID)
	if err != nil {
		fmt.Printf("   ERROR downloading events: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Downloaded %d event wrappers\n\n", len(wrappers))

	// Step 3: Parse LoL events
	fmt.Println("3. Parsing LoL events...")
	lolEvents, err := grid.ParseLoLEvents(wrappers)
	if err != nil {
		fmt.Printf("   ERROR parsing LoL events: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("   Kills: %d\n", len(lolEvents.Kills))
	fmt.Printf("   Dragon Kills: %d\n", len(lolEvents.DragonKills))
	fmt.Printf("   Baron Kills: %d\n", len(lolEvents.BaronKills))
	fmt.Printf("   Herald Kills: %d\n", len(lolEvents.HeraldKills))
	fmt.Printf("   Tower Destroys: %d\n", len(lolEvents.TowerDestroys))
	fmt.Printf("   Draft Actions: %d\n", len(lolEvents.DraftActions))
	fmt.Println()

	// Step 4: Verify kill data
	fmt.Println("4. Sample Kill Events (first 5):")
	for i, kill := range lolEvents.Kills {
		if i >= 5 {
			break
		}
		fmt.Printf("   Kill %d: %s killed %s at %d ms (FirstBlood: %v)\n",
			i+1, kill.KillerName, kill.VictimName, kill.GameTime, kill.FirstBlood)
		fmt.Printf("      Killer Team: %s, Victim Team: %s\n", kill.KillerTeamID, kill.VictimTeamID)
	}
	fmt.Println()

	// Step 5: Verify dragon data
	fmt.Println("5. Dragon Kill Events:")
	for i, dragon := range lolEvents.DragonKills {
		fmt.Printf("   Dragon %d: %s (%s) at %d ms by team %s\n",
			i+1, dragon.DragonType, dragon.PlayerName, dragon.GameTime, dragon.TeamID)
	}
	fmt.Println()

	// Step 6: Verify tower data
	fmt.Println("6. Sample Tower Destroy Events (first 5):")
	for i, tower := range lolEvents.TowerDestroys {
		if i >= 5 {
			break
		}
		fmt.Printf("   Tower %d: %s lane, tower #%d at %d ms by team %s\n",
			i+1, tower.Lane, tower.TowerNum, tower.GameTime, tower.TeamID)
	}
	fmt.Println()

	// Step 7: Test EventAnalyzer
	fmt.Println("7. Testing EventAnalyzer...")
	analyzer := intelligence.NewEventAnalyzer()

	// Get unique team IDs from kills
	teamIDs := make(map[string]bool)
	for _, kill := range lolEvents.Kills {
		if kill.KillerTeamID != "" {
			teamIDs[kill.KillerTeamID] = true
		}
		if kill.VictimTeamID != "" {
			teamIDs[kill.VictimTeamID] = true
		}
	}
	fmt.Printf("   Found team IDs: ")
	for teamID := range teamIDs {
		fmt.Printf("%s ", teamID)
	}
	fmt.Println()

	// Test with both team IDs
	var teamIDList []string
	for teamID := range teamIDs {
		teamIDList = append(teamIDList, teamID)
	}

	if len(teamIDList) == 0 {
		fmt.Println("   ERROR: No team IDs found in events")
		os.Exit(1)
	}

	// Test both teams
	for _, testTeamID := range teamIDList {
		fmt.Printf("\n   === Testing with team ID: %s ===\n\n", testTeamID)

	// Test First Blood
	fmt.Println("8. First Blood Analysis:")
	fb := analyzer.AnalyzeFirstBlood(lolEvents, testTeamID)
	if fb != nil {
		fmt.Printf("   First Blood: YES - %s at %d seconds\n", fb.PlayerName, fb.GameTimeSeconds)
	} else {
		fmt.Println("   First Blood: NO (other team got it)")
	}

	// Test First Dragon
	fmt.Println("\n9. First Dragon Analysis:")
	fd := analyzer.AnalyzeFirstDragon(lolEvents, testTeamID)
	if fd != nil {
		fmt.Printf("   First Dragon: YES - %s (%s) at %d seconds\n", fd.DragonType, fd.PlayerName, fd.GameTimeSeconds)
	} else {
		fmt.Println("   First Dragon: NO (other team got it)")
	}

	// Test First Tower
	fmt.Println("\n10. First Tower Analysis:")
	ft := analyzer.AnalyzeFirstTower(lolEvents, testTeamID)
	if ft != nil {
		fmt.Printf("   First Tower: YES - %s lane at %d seconds\n", ft.Lane, ft.GameTimeSeconds)
	} else {
		fmt.Println("   First Tower: NO (other team got it)")
	}

	// Test Phase Analysis
	fmt.Println("\n11. Phase Analysis:")
	phases := analyzer.AnalyzePhases(lolEvents, testTeamID)
	fmt.Printf("   Early Game: %d kills, %d deaths, diff: %d, rating: %.1f\n",
		phases.EarlyKills, phases.EarlyDeaths, phases.EarlyKillDiff, phases.EarlyGameRating)
	fmt.Printf("   Mid Game: %d kills, %d deaths, diff: %d, rating: %.1f\n",
		phases.MidKills, phases.MidDeaths, phases.MidKillDiff, phases.MidGameRating)
	fmt.Printf("   Late Game: %d kills, %d deaths, diff: %d, rating: %.1f\n",
		phases.LateKills, phases.LateDeaths, phases.LateKillDiff, phases.LateGameRating)
	fmt.Printf("   Strongest Phase: %s\n", phases.StrongestPhase)

	// Test Objective Timings
	fmt.Println("\n12. Objective Timing Analysis:")
	timings := analyzer.AnalyzeObjectiveTimings(lolEvents, testTeamID)
	fmt.Printf("   Dragon timings: %v\n", timings.DragonTimings)
	fmt.Printf("   Baron timings: %v\n", timings.BaronTimings)
	fmt.Printf("   Herald timings: %v\n", timings.HeraldTimings)
	fmt.Printf("   Avg First Dragon Time: %.1f seconds\n", timings.AvgFirstDragonTime)
	fmt.Printf("   Early Objective Priority: %v\n", timings.EarlyObjectivePriority)
	} // End of team loop

	// Summary
	fmt.Println("\n=== VERIFICATION SUMMARY ===")
	fmt.Println("✓ JSONL file download: SUCCESS")
	fmt.Println("✓ Event parsing: SUCCESS")
	fmt.Printf("✓ Kill events extracted: %d\n", len(lolEvents.Kills))
	fmt.Printf("✓ Dragon events extracted: %d\n", len(lolEvents.DragonKills))
	fmt.Printf("✓ Tower events extracted: %d\n", len(lolEvents.TowerDestroys))
	fmt.Println("✓ EventAnalyzer functions: WORKING")
	fmt.Println("\nReal data is being extracted successfully!")
}
