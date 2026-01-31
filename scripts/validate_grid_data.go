// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"scout9/pkg/grid"
)

func main() {
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		apiKey = "hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO"
	}

	client := grid.NewClient(apiKey, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("GRID API DATA VALIDATION")
	fmt.Println("=" + strings.Repeat("=", 79))

	// Test LoL Series
	fmt.Println("\n## Testing LoL Series (2692648 - LEC Summer 2024)")
	fmt.Println("-" + strings.Repeat("-", 79))
	
	lolState, err := client.GetSeriesState(ctx, "2692648")
	if err != nil {
		fmt.Printf("ERROR: Failed to get LoL series state: %v\n", err)
	} else {
		validateLoLSeriesState(lolState)
	}

	// Test VALORANT Series
	fmt.Println("\n## Testing VALORANT Series (2629390 - VCT Americas)")
	fmt.Println("-" + strings.Repeat("-", 79))
	
	valState, err := client.GetSeriesState(ctx, "2629390")
	if err != nil {
		fmt.Printf("ERROR: Failed to get VALORANT series state: %v\n", err)
	} else {
		validateVALSeriesState(valState)
	}

	// Test File Download API - LoL Events
	fmt.Println("\n## Testing File Download API - LoL Events")
	fmt.Println("-" + strings.Repeat("-", 79))
	
	lolEvents, err := client.DownloadEvents(ctx, "2692648")
	if err != nil {
		fmt.Printf("ERROR: Failed to download LoL events: %v\n", err)
	} else {
		validateLoLEvents(lolEvents)
	}

	// Test File Download API - VALORANT Events
	fmt.Println("\n## Testing File Download API - VALORANT Events")
	fmt.Println("-" + strings.Repeat("-", 79))
	
	valEvents, err := client.DownloadEvents(ctx, "2629390")
	if err != nil {
		fmt.Printf("ERROR: Failed to download VALORANT events: %v\n", err)
	} else {
		validateVALEvents(valEvents)
	}

	fmt.Println("\n" + "=" + strings.Repeat("=", 79))
	fmt.Println("VALIDATION COMPLETE")
	fmt.Println("=" + strings.Repeat("=", 79))
}

func validateLoLSeriesState(state *grid.SeriesState) {
	fmt.Printf("Series ID: %s\n", state.ID)
	fmt.Printf("Started: %v, Finished: %v\n", state.Started, state.Finished)
	fmt.Printf("Teams: %d\n", len(state.Teams))
	fmt.Printf("Games: %d\n", len(state.Games))

	for _, team := range state.Teams {
		fmt.Printf("  Team: %s (ID: %s) - Won: %v, Score: %d, K/D: %d/%d\n",
			team.Name, team.ID, team.Won, team.Score, team.Kills, team.Deaths)
	}

	for i, game := range state.Games {
		fmt.Printf("\nGame %d (ID: %s):\n", i+1, game.ID)
		fmt.Printf("  Map: %s, Duration: %ds, Finished: %v\n", game.Map, game.Duration, game.Finished)
		fmt.Printf("  Draft Actions: %d\n", len(game.DraftActions))
		
		// Show draft
		bans := 0
		picks := 0
		for _, da := range game.DraftActions {
			if da.Action == "ban" {
				bans++
			} else if da.Action == "pick" {
				picks++
			}
		}
		fmt.Printf("    Bans: %d, Picks: %d\n", bans, picks)

		for _, team := range game.Teams {
			fmt.Printf("  Team %s (%s side):\n", team.Name, team.Side)
			fmt.Printf("    Score: %d, Won: %v, K/D: %d/%d, NetWorth: %d\n",
				team.Score, team.Won, team.Kills, team.Deaths, team.NetWorth)
			fmt.Printf("    Structures Destroyed: %d\n", team.StructuresDestroyed)
			
			// Show objectives
			if len(team.Objectives) > 0 {
				fmt.Printf("    Objectives:\n")
				for _, obj := range team.Objectives {
					fmt.Printf("      - %s: %d\n", obj.Type, obj.CompletionCount)
				}
			}

			fmt.Printf("    Players:\n")
			for _, player := range team.Players {
				fmt.Printf("      %s (%s): K/D/A: %d/%d/%d, NetWorth: %d, Items: %d\n",
					player.Name, player.Character, player.Kills, player.Deaths, player.Assists,
					player.NetWorth, len(player.Items))
				
				// Show player objectives
				if len(player.Objectives) > 0 {
					objTypes := make(map[string]int)
					for _, obj := range player.Objectives {
						objTypes[obj.Type] += obj.CompletionCount
					}
					fmt.Printf("        Objectives: %v\n", objTypes)
				}
			}
		}
	}
}

func validateVALSeriesState(state *grid.SeriesState) {
	fmt.Printf("Series ID: %s\n", state.ID)
	fmt.Printf("Started: %v, Finished: %v\n", state.Started, state.Finished)
	fmt.Printf("Teams: %d\n", len(state.Teams))
	fmt.Printf("Games: %d\n", len(state.Games))

	for _, team := range state.Teams {
		fmt.Printf("  Team: %s (ID: %s) - Won: %v, Score: %d, K/D: %d/%d\n",
			team.Name, team.ID, team.Won, team.Score, team.Kills, team.Deaths)
	}

	for i, game := range state.Games {
		fmt.Printf("\nGame %d (ID: %s):\n", i+1, game.ID)
		fmt.Printf("  Map: %s, Duration: %ds, Finished: %v\n", game.Map, game.Duration, game.Finished)

		for _, team := range game.Teams {
			fmt.Printf("  Team %s (%s side):\n", team.Name, team.Side)
			fmt.Printf("    Score: %d, Won: %v, K/D: %d/%d\n",
				team.Score, team.Won, team.Kills, team.Deaths)
			
			// Show objectives (plants, defuses)
			if len(team.Objectives) > 0 {
				fmt.Printf("    Objectives:\n")
				for _, obj := range team.Objectives {
					fmt.Printf("      - %s: %d\n", obj.Type, obj.CompletionCount)
				}
			}

			fmt.Printf("    Players:\n")
			for _, player := range team.Players {
				fmt.Printf("      %s (%s): K/D/A: %d/%d/%d\n",
					player.Name, player.Character, player.Kills, player.Deaths, player.Assists)
			}
		}
	}
}

func validateLoLEvents(events []grid.EventWrapper) {
	fmt.Printf("Total Event Wrappers: %d\n", len(events))

	// Count event types
	eventCounts := make(map[string]int)
	for _, wrapper := range events {
		for _, event := range wrapper.Events {
			eventType := event.GetEventType()
			eventCounts[eventType]++
		}
	}

	fmt.Printf("Event Types Found:\n")
	for eventType, count := range eventCounts {
		fmt.Printf("  %s: %d\n", eventType, count)
	}

	// Parse LoL events
	lolData, err := grid.ParseLoLEvents(events)
	if err != nil {
		fmt.Printf("ERROR parsing LoL events: %v\n", err)
		return
	}

	fmt.Printf("\nParsed LoL Data:\n")
	fmt.Printf("  Kills: %d\n", len(lolData.Kills))
	fmt.Printf("  Dragon Kills: %d\n", len(lolData.DragonKills))
	fmt.Printf("  Baron Kills: %d\n", len(lolData.BaronKills))
	fmt.Printf("  Herald Kills: %d\n", len(lolData.HeraldKills))
	fmt.Printf("  Void Grub Kills: %d\n", len(lolData.VoidGrubKills))
	fmt.Printf("  Tower Destroys: %d\n", len(lolData.TowerDestroys))
	fmt.Printf("  Draft Actions: %d\n", len(lolData.DraftActions))

	// Show sample kill with position
	if len(lolData.Kills) > 0 {
		kill := lolData.Kills[0]
		fmt.Printf("\nSample Kill Event:\n")
		fmt.Printf("  Killer: %s (Team: %s)\n", kill.KillerName, kill.KillerTeamID)
		fmt.Printf("  Victim: %s (Team: %s)\n", kill.VictimName, kill.VictimTeamID)
		fmt.Printf("  First Blood: %v\n", kill.FirstBlood)
		if kill.KillerPosition != nil {
			fmt.Printf("  Killer Position: (%.0f, %.0f)\n", kill.KillerPosition.X, kill.KillerPosition.Y)
		}
		fmt.Printf("  Game Time: %dms\n", kill.GameTime)
	}

	// Show dragon types
	if len(lolData.DragonKills) > 0 {
		fmt.Printf("\nDragon Kills:\n")
		for _, dragon := range lolData.DragonKills {
			fmt.Printf("  %s dragon by %s at %dms\n", dragon.DragonType, dragon.PlayerName, dragon.GameTime)
		}
	}
}

func validateVALEvents(events []grid.EventWrapper) {
	fmt.Printf("Total Event Wrappers: %d\n", len(events))

	// Count event types
	eventCounts := make(map[string]int)
	for _, wrapper := range events {
		for _, event := range wrapper.Events {
			eventType := event.GetEventType()
			eventCounts[eventType]++
		}
	}

	fmt.Printf("Event Types Found:\n")
	for eventType, count := range eventCounts {
		fmt.Printf("  %s: %d\n", eventType, count)
	}

	// Parse VALORANT events
	valData, err := grid.ParseVALEvents(events)
	if err != nil {
		fmt.Printf("ERROR parsing VALORANT events: %v\n", err)
		return
	}

	fmt.Printf("\nParsed VALORANT Data:\n")
	fmt.Printf("  Map: %s\n", valData.MapName)
	fmt.Printf("  Kills: %d\n", len(valData.Kills))
	fmt.Printf("  Round Ends: %d\n", len(valData.RoundEnds))
	fmt.Printf("  Plants: %d\n", len(valData.Plants))
	fmt.Printf("  Defuses: %d\n", len(valData.Defuses))

	// Show round win types
	winTypes := make(map[string]int)
	for _, round := range valData.RoundEnds {
		winTypes[round.WinType]++
	}
	fmt.Printf("\nRound Win Types:\n")
	for winType, count := range winTypes {
		fmt.Printf("  %s: %d\n", winType, count)
	}

	// Show site preferences
	siteCounts := make(map[string]int)
	for _, plant := range valData.Plants {
		if plant.Site != "" {
			siteCounts[plant.Site]++
		}
	}
	if len(siteCounts) > 0 {
		fmt.Printf("\nSite Preferences:\n")
		total := 0
		for _, count := range siteCounts {
			total += count
		}
		for site, count := range siteCounts {
			fmt.Printf("  %s: %d (%.1f%%)\n", site, count, float64(count)/float64(total)*100)
		}
	}

	// Show sample kill with agent info
	if len(valData.Kills) > 0 {
		kill := valData.Kills[0]
		fmt.Printf("\nSample Kill Event:\n")
		fmt.Printf("  Killer: %s (%s)\n", kill.KillerName, kill.KillerAgent)
		fmt.Printf("  Victim: %s (%s)\n", kill.VictimName, kill.VictimAgent)
		fmt.Printf("  Round: %d, Map: %s\n", kill.RoundNum, kill.MapName)
		if kill.KillerPosition != nil {
			fmt.Printf("  Killer Position: (%.0f, %.0f)\n", kill.KillerPosition.X, kill.KillerPosition.Y)
		}
	}
}
