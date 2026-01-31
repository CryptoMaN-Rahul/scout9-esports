// test_new_parser.go - Test the fixed GRID event parser
package main

import (
	"context"
	"fmt"
	"strings"

	"scout9/pkg/grid"
)

const APIKey = "hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO"

func main() {
	ctx := context.Background()

	// Create client without cache for testing
	client := grid.NewClient(APIKey, nil)

	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("TESTING NEW GRID EVENT PARSER")
	fmt.Println("=" + strings.Repeat("=", 79))

	// Test LoL series
	fmt.Println("\n📊 Testing LoL Series (2692648 - LEC Summer 2024)")
	fmt.Println(strings.Repeat("-", 60))

	lolWrappers, err := client.DownloadEvents(ctx, "2692648")
	if err != nil {
		fmt.Printf("❌ Error downloading LoL events: %v\n", err)
	} else {
		fmt.Printf("✅ Downloaded %d event wrappers\n", len(lolWrappers))

		lolData, err := grid.ParseLoLEvents(lolWrappers)
		if err != nil {
			fmt.Printf("❌ Error parsing LoL events: %v\n", err)
		} else {
			fmt.Printf("\n📈 LoL Event Summary:\n")
			fmt.Printf("   Kills: %d\n", len(lolData.Kills))
			fmt.Printf("   Dragon Kills: %d\n", len(lolData.DragonKills))
			fmt.Printf("   Baron Kills: %d\n", len(lolData.BaronKills))
			fmt.Printf("   Herald Kills: %d\n", len(lolData.HeraldKills))
			fmt.Printf("   Void Grub Kills: %d\n", len(lolData.VoidGrubKills))
			fmt.Printf("   Tower Destroys: %d\n", len(lolData.TowerDestroys))
			fmt.Printf("   Draft Actions: %d\n", len(lolData.DraftActions))

			// Show sample kills with position
			fmt.Println("\n🎯 Sample Kills with Position:")
			for i, kill := range lolData.Kills {
				if i >= 5 {
					break
				}
				posStr := "no position"
				if kill.KillerPosition != nil {
					posStr = fmt.Sprintf("(%.0f, %.0f)", kill.KillerPosition.X, kill.KillerPosition.Y)
				}
				fmt.Printf("   %s killed %s at %s (First Blood: %v)\n",
					kill.KillerName, kill.VictimName, posStr, kill.FirstBlood)
			}

			// Show dragon kills
			fmt.Println("\n🐉 Dragon Kills:")
			for _, dragon := range lolData.DragonKills {
				fmt.Printf("   %s killed %s dragon\n", dragon.PlayerName, dragon.DragonType)
			}

			// Show draft
			fmt.Println("\n📋 Draft Actions:")
			for _, draft := range lolData.DraftActions {
				fmt.Printf("   %s %s %s\n", draft.TeamName, draft.Action, draft.CharacterName)
			}

			// Show tower destroys
			fmt.Println("\n🏰 Tower Destroys:")
			for _, tower := range lolData.TowerDestroys {
				fmt.Printf("   %s destroyed %s (lane: %s)\n", tower.TeamName, tower.TowerID, tower.Lane)
			}
		}
	}

	// Test VALORANT series
	fmt.Println("\n\n📊 Testing VALORANT Series (2629390 - VCT Americas)")
	fmt.Println(strings.Repeat("-", 60))

	valWrappers, err := client.DownloadEvents(ctx, "2629390")
	if err != nil {
		fmt.Printf("❌ Error downloading VALORANT events: %v\n", err)
	} else {
		fmt.Printf("✅ Downloaded %d event wrappers\n", len(valWrappers))

		valData, err := grid.ParseVALEvents(valWrappers)
		if err != nil {
			fmt.Printf("❌ Error parsing VALORANT events: %v\n", err)
		} else {
			fmt.Printf("\n📈 VALORANT Event Summary:\n")
			fmt.Printf("   Map: %s\n", valData.MapName)
			fmt.Printf("   Kills: %d\n", len(valData.Kills))
			fmt.Printf("   Round Ends: %d\n", len(valData.RoundEnds))
			fmt.Printf("   Plants: %d\n", len(valData.Plants))
			fmt.Printf("   Defuses: %d\n", len(valData.Defuses))

			// Show sample kills with position
			fmt.Println("\n🎯 Sample Kills with Position:")
			for i, kill := range valData.Kills {
				if i >= 5 {
					break
				}
				posStr := "no position"
				if kill.KillerPosition != nil {
					posStr = fmt.Sprintf("(%.0f, %.0f)", kill.KillerPosition.X, kill.KillerPosition.Y)
				}
				fmt.Printf("   %s (%s) killed %s (%s) at %s [Round %d]\n",
					kill.KillerName, kill.KillerAgent, kill.VictimName, kill.VictimAgent, posStr, kill.RoundNum)
			}

			// Show round results
			fmt.Println("\n🏆 Round Results:")
			for i, round := range valData.RoundEnds {
				if i >= 10 {
					fmt.Printf("   ... and %d more rounds\n", len(valData.RoundEnds)-10)
					break
				}
				fmt.Printf("   Round %d: %s won by %s\n", round.RoundNum, round.WinnerName, round.WinType)
			}

			// Show plants with site inference
			fmt.Println("\n💣 Spike Plants:")
			for i, plant := range valData.Plants {
				if i >= 10 {
					fmt.Printf("   ... and %d more plants\n", len(valData.Plants)-10)
					break
				}
				posStr := "no position"
				if plant.Position != nil {
					posStr = fmt.Sprintf("(%.0f, %.0f)", plant.Position.X, plant.Position.Y)
				}
				fmt.Printf("   Round %d: %s (%s) planted at %s - Site: %s\n",
					plant.RoundNum, plant.PlayerName, plant.Agent, posStr, plant.Site)
			}

			// Calculate site preference
			siteCount := make(map[string]int)
			for _, plant := range valData.Plants {
				if plant.Site != "" {
					siteCount[plant.Site]++
				}
			}
			fmt.Println("\n📊 Site Preference:")
			total := len(valData.Plants)
			for site, count := range siteCount {
				pct := float64(count) / float64(total) * 100
				fmt.Printf("   Site %s: %d plants (%.1f%%)\n", site, count, pct)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ Parser validation complete!")
	fmt.Println(strings.Repeat("=", 80))
}
