// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"scout9/pkg/grid"
)

func main() {
	// Load environment variables
	godotenv.Load()

	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		fmt.Println("ERROR: GRID_API_KEY not set")
		os.Exit(1)
	}

	// Create client
	client := grid.NewClient(apiKey, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test with LoL series (LEC Summer 2024)
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("🎮 TESTING LEAGUE OF LEGENDS (Series: 2692648)")
	fmt.Println("=" + string(make([]byte, 60)))
	testSeries(ctx, client, "2692648")

	// Test with VALORANT series (VCT Americas Kickoff 2024)
	fmt.Println("\n" + string(make([]byte, 60)) + "=")
	fmt.Println("🎮 TESTING VALORANT (Series: 2629390)")
	fmt.Println(string(make([]byte, 60)) + "=")
	testSeries(ctx, client, "2629390")
}

func testSeries(ctx context.Context, client *grid.Client, seriesID string) {
	state, err := client.GetSeriesState(ctx, seriesID)
	if err != nil {
		fmt.Printf("ERROR fetching series state: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Series: %s (Started: %v, Finished: %v)\n", state.ID, state.Started, state.Finished)
	fmt.Printf("   Teams: %d, Games: %d\n", len(state.Teams), len(state.Games))

	// Summary counters
	summary := map[string]bool{
		"assistDetails":   false,
		"loadoutValue":    false,
		"abilities":       false,
		"weaponKills":     false,
		"multikills":      false,
		"objectives":      false,
		"segments":        false,
	}

	// Analyze first game only
	if len(state.Games) > 0 {
		game := state.Games[0]
		fmt.Printf("\n📊 Game 1: %s\n", game.ID)
		fmt.Printf("   Map: %s, Duration: %ds, Finished: %v\n", game.Map, game.Duration, game.Finished)

		// Check segments (VALORANT rounds)
		if len(game.Segments) > 0 {
			summary["segments"] = true
			fmt.Printf("   ✅ Segments/Rounds: %d\n", len(game.Segments))
		}

		for _, team := range game.Teams {
			fmt.Printf("\n   👥 Team: %s (Side: %s, Won: %v)\n", team.Name, team.Side, team.Won)
			fmt.Printf("      K/D: %d/%d, NetWorth: %d, LoadoutValue: %d\n", 
				team.Kills, team.Deaths, team.NetWorth, team.LoadoutValue)

			if team.LoadoutValue > 0 {
				summary["loadoutValue"] = true
			}

			// Show first 2 players
			for i, player := range team.Players {
				if i >= 2 {
					fmt.Printf("      ... and %d more players\n", len(team.Players)-2)
					break
				}

				fmt.Printf("\n      🎮 %s (%s): K/D/A %d/%d/%d\n", 
					player.Name, player.Character, player.Kills, player.Deaths, player.Assists)

				// Assist Details
				if len(player.AssistDetails) > 0 {
					summary["assistDetails"] = true
					fmt.Printf("         ✅ Assist Network: %d teammates assisted\n", len(player.AssistDetails))
				}

				// Abilities (VALORANT)
				if len(player.Abilities) > 0 {
					summary["abilities"] = true
					fmt.Printf("         ✅ Abilities: %v\n", getAbilityNames(player.Abilities))
				}

				// Weapon Kills (VALORANT)
				if len(player.WeaponKills) > 0 {
					summary["weaponKills"] = true
					fmt.Printf("         ✅ Weapon Kills: %d weapons used\n", len(player.WeaponKills))
				}

				// Multikills
				if len(player.Multikills) > 0 {
					summary["multikills"] = true
					fmt.Printf("         ✅ Multikills: %v\n", player.Multikills)
				}

				// Objectives
				if len(player.Objectives) > 0 {
					summary["objectives"] = true
					fmt.Printf("         ✅ Objectives: %d types\n", len(player.Objectives))
				}
			}
		}
	}

	// Print summary
	fmt.Println("\n📋 FIELD AVAILABILITY:")
	for field, available := range summary {
		status := "❌"
		if available {
			status = "✅"
		}
		fmt.Printf("   %s %s\n", status, field)
	}

	// Output sample player JSON
	if len(state.Games) > 0 && len(state.Games[0].Teams) > 0 && len(state.Games[0].Teams[0].Players) > 0 {
		player := state.Games[0].Teams[0].Players[0]
		jsonBytes, _ := json.MarshalIndent(player, "   ", "  ")
		fmt.Printf("\n📄 Sample Player JSON:\n%s\n", string(jsonBytes))
	}
}

func getAbilityNames(abilities []grid.AbilityUsage) []string {
	names := make([]string, len(abilities))
	for i, a := range abilities {
		names[i] = a.AbilityName
	}
	return names
}
