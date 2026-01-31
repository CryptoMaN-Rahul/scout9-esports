// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"scout9/pkg/grid"
)

func main() {
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		apiKey = "hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO"
	}

	client := grid.NewClient(apiKey, nil)
	ctx := context.Background()

	// Download events for LoL series
	fmt.Println("=== DEBUGGING LOL EVENTS (Series 2692648) ===")
	wrappers, err := client.DownloadEvents(ctx, "2692648")
	if err != nil {
		fmt.Printf("Error downloading events: %v\n", err)
		return
	}

	fmt.Printf("Total event wrappers: %d\n\n", len(wrappers))

	// Look for game start events
	fmt.Println("=== LOOKING FOR GAME START EVENTS ===")
	gameStartFound := false
	for i, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			if event.Action == "started" {
				targetType := ""
				if event.Target != nil {
					targetType = event.Target.Type
				}
				actorType := ""
				if event.Actor != nil {
					actorType = event.Actor.Type
				}
				fmt.Printf("Wrapper %d: action=started, actor.type=%s, target.type=%s, time=%v\n",
					i, actorType, targetType, wrapper.OccurredAt)
				if targetType == "game" {
					gameStartFound = true
					fmt.Println("  ^^^ GAME START FOUND!")
				}
			}
		}
		if i > 50 && gameStartFound {
			break
		}
	}

	// Look at first few kill events with clock data
	fmt.Println("\n=== FIRST 5 KILL EVENTS WITH CLOCK DATA ===")
	killCount := 0
	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			actorType := ""
			if event.Actor != nil {
				actorType = event.Actor.Type
			}
			targetType := ""
			if event.Target != nil {
				targetType = event.Target.Type
			}
			if actorType == "player" && event.Action == "killed" && targetType == "player" {
				killCount++
				fmt.Printf("Kill %d at %v\n", killCount, wrapper.OccurredAt)
				
				// Check for position data
				if event.Actor != nil {
					pos := event.GetActorPosition()
					if pos != nil {
						fmt.Printf("  Killer position: x=%.0f, y=%.0f\n", pos.X, pos.Y)
					} else {
						fmt.Println("  Killer position: NOT AVAILABLE")
					}
				}
				if event.Target != nil {
					pos := event.GetTargetPosition()
					if pos != nil {
						fmt.Printf("  Victim position: x=%.0f, y=%.0f\n", pos.X, pos.Y)
					} else {
						fmt.Println("  Victim position: NOT AVAILABLE")
					}
				}
				
				// Check for clock data in seriesState
				if event.SeriesState != nil {
					if games, ok := event.SeriesState["games"].([]interface{}); ok && len(games) > 0 {
						if game, ok := games[0].(map[string]interface{}); ok {
							if clock, ok := game["clock"].(map[string]interface{}); ok {
								if currentSeconds, ok := clock["currentSeconds"].(float64); ok {
									fmt.Printf("  In-game clock: %.0f seconds (%.1f min)\n", currentSeconds, currentSeconds/60)
								}
							}
						}
					}
				} else {
					fmt.Println("  SeriesState: NOT AVAILABLE")
				}
				
				if killCount >= 5 {
					break
				}
			}
		}
		if killCount >= 5 {
			break
		}
	}

	// Look at first dragon event with clock
	fmt.Println("\n=== FIRST DRAGON EVENT WITH CLOCK ===")
	for _, wrapper := range wrappers {
		for _, event := range wrapper.Events {
			if event.Action == "killed" {
				targetType := ""
				targetID := ""
				if event.Target != nil {
					targetType = event.Target.Type
					targetID = event.Target.ID
				}
				if targetType == "ATierNPC" && (contains(targetID, "drake") || contains(targetID, "dragon")) {
					fmt.Printf("Dragon killed at %v\n", wrapper.OccurredAt)
					fmt.Printf("  Target ID: %s\n", targetID)
					
					// Check for clock data
					if event.SeriesState != nil {
						if games, ok := event.SeriesState["games"].([]interface{}); ok && len(games) > 0 {
							if game, ok := games[0].(map[string]interface{}); ok {
								if clock, ok := game["clock"].(map[string]interface{}); ok {
									if currentSeconds, ok := clock["currentSeconds"].(float64); ok {
										fmt.Printf("  In-game clock: %.0f seconds (%.1f min)\n", currentSeconds, currentSeconds/60)
									}
								}
							}
						}
					}
					
					// Print partial event for debugging
					eventJSON, _ := json.MarshalIndent(event.SeriesState, "  ", "  ")
					if len(eventJSON) > 500 {
						eventJSON = eventJSON[:500]
					}
					fmt.Printf("  SeriesState (partial):\n%s...\n", string(eventJSON))
					return
				}
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsLower(s, substr))
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if toLower(s[i:i+len(substr)]) == toLower(substr) {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
