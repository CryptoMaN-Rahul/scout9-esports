// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"scout9/pkg/grid"
	"scout9/pkg/intelligence"
)

// DataCapture holds all captured data for analysis
type DataCapture struct {
	CapturedAt      time.Time                      `json:"capturedAt"`
	LoLSeriesID     string                         `json:"lolSeriesId"`
	VALSeriesID     string                         `json:"valSeriesId"`
	LoLSeriesState  *grid.SeriesState              `json:"lolSeriesState,omitempty"`
	VALSeriesState  *grid.SeriesState              `json:"valSeriesState,omitempty"`
	LoLTeamAnalysis *intelligence.TeamAnalysis     `json:"lolTeamAnalysis,omitempty"`
	VALTeamAnalysis *intelligence.TeamAnalysis     `json:"valTeamAnalysis,omitempty"`
	LoLPlayers      []*intelligence.PlayerProfile  `json:"lolPlayers,omitempty"`
	VALPlayers      []*intelligence.PlayerProfile  `json:"valPlayers,omitempty"`
	FieldAvailability map[string]FieldStatus       `json:"fieldAvailability"`
	Errors          []string                       `json:"errors,omitempty"`
}

// FieldStatus tracks whether a field has data
type FieldStatus struct {
	Available   bool   `json:"available"`
	SampleValue string `json:"sampleValue,omitempty"`
	Count       int    `json:"count,omitempty"`
}

func main() {
	// Get API key
	apiKey := os.Getenv("GRID_API_KEY")
	if apiKey == "" {
		apiKey = "hmadti95q6W8muoklBmMh5V7YGwH4rB0It3kKhOO"
	}

	// Create client
	client := grid.NewClient(apiKey, nil)
	ctx := context.Background()

	// Test series IDs
	lolSeriesID := "2692648"  // LEC Summer 2024 - Team Heretics vs GIANTX
	valSeriesID := "2629390"  // VCT Americas - FURIA vs NRG

	capture := &DataCapture{
		CapturedAt:        time.Now(),
		LoLSeriesID:       lolSeriesID,
		VALSeriesID:       valSeriesID,
		FieldAvailability: make(map[string]FieldStatus),
		Errors:            make([]string, 0),
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  SCOUT9 GRID Data Capture Test")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("  Timestamp: %s\n", capture.CapturedAt.Format(time.RFC3339))
	fmt.Println()

	// ═══════════════════════════════════════════════════════════════
	// PART 1: Fetch LoL Series State
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("📊 Fetching LoL Series State...")
	lolState, err := client.GetSeriesState(ctx, lolSeriesID)
	if err != nil {
		capture.Errors = append(capture.Errors, fmt.Sprintf("LoL series fetch error: %v", err))
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		capture.LoLSeriesState = lolState
		fmt.Printf("   ✅ Series: %s (Finished: %v)\n", lolSeriesID, lolState.Finished)
		analyzeLoLSeriesState(lolState, capture)
	}

	// ═══════════════════════════════════════════════════════════════
	// PART 2: Fetch VALORANT Series State
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("\n📊 Fetching VALORANT Series State...")
	valState, err := client.GetSeriesState(ctx, valSeriesID)
	if err != nil {
		capture.Errors = append(capture.Errors, fmt.Sprintf("VAL series fetch error: %v", err))
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		capture.VALSeriesState = valState
		fmt.Printf("   ✅ Series: %s (Finished: %v)\n", valSeriesID, valState.Finished)
		analyzeVALSeriesState(valState, capture)
	}

	// ═══════════════════════════════════════════════════════════════
	// PART 3: Run Intelligence Analysis on LoL
	// ═══════════════════════════════════════════════════════════════
	if lolState != nil && len(lolState.Teams) > 0 {
		fmt.Println("\n🧠 Running LoL Intelligence Analysis...")
		lolAnalyzer := intelligence.NewLoLAnalyzer()
		
		teamID := lolState.Teams[0].ID
		teamName := lolState.Teams[0].Name
		
		teamAnalysis, err := lolAnalyzer.AnalyzeTeam(ctx, teamID, teamName, []*grid.SeriesState{lolState}, nil)
		if err != nil {
			capture.Errors = append(capture.Errors, fmt.Sprintf("LoL team analysis error: %v", err))
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			capture.LoLTeamAnalysis = teamAnalysis
			fmt.Printf("   ✅ Team: %s (Win Rate: %.1f%%)\n", teamName, teamAnalysis.WinRate*100)
			printLoLTeamAnalysis(teamAnalysis)
		}

		players, err := lolAnalyzer.AnalyzePlayers(ctx, teamID, []*grid.SeriesState{lolState})
		if err != nil {
			capture.Errors = append(capture.Errors, fmt.Sprintf("LoL player analysis error: %v", err))
		} else {
			capture.LoLPlayers = players
			fmt.Printf("   ✅ Analyzed %d players\n", len(players))
			printPlayerProfiles("LoL", players)
		}
	}

	// ═══════════════════════════════════════════════════════════════
	// PART 4: Run Intelligence Analysis on VALORANT
	// ═══════════════════════════════════════════════════════════════
	if valState != nil && len(valState.Teams) > 0 {
		fmt.Println("\n🧠 Running VALORANT Intelligence Analysis...")
		valAnalyzer := intelligence.NewVALAnalyzer()
		
		teamID := valState.Teams[0].ID
		teamName := valState.Teams[0].Name
		
		teamAnalysis, err := valAnalyzer.AnalyzeTeam(ctx, teamID, teamName, []*grid.SeriesState{valState}, nil)
		if err != nil {
			capture.Errors = append(capture.Errors, fmt.Sprintf("VAL team analysis error: %v", err))
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			capture.VALTeamAnalysis = teamAnalysis
			fmt.Printf("   ✅ Team: %s (Win Rate: %.1f%%)\n", teamName, teamAnalysis.WinRate*100)
			printVALTeamAnalysis(teamAnalysis)
		}

		players, err := valAnalyzer.AnalyzePlayers(ctx, teamID, []*grid.SeriesState{valState})
		if err != nil {
			capture.Errors = append(capture.Errors, fmt.Sprintf("VAL player analysis error: %v", err))
		} else {
			capture.VALPlayers = players
			fmt.Printf("   ✅ Analyzed %d players\n", len(players))
			printPlayerProfiles("VALORANT", players)
		}
	}

	// ═══════════════════════════════════════════════════════════════
	// PART 5: Save captured data to file
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("\n💾 Saving captured data...")
	
	outputFile := "scripts/grid_data_capture.json"
	data, err := json.MarshalIndent(capture, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
	
	fmt.Printf("   ✅ Saved to %s (%d bytes)\n", outputFile, len(data))

	// ═══════════════════════════════════════════════════════════════
	// PART 6: Print Field Availability Summary
	// ═══════════════════════════════════════════════════════════════
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("  FIELD AVAILABILITY SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	
	for field, status := range capture.FieldAvailability {
		icon := "❌"
		if status.Available {
			icon = "✅"
		}
		fmt.Printf("  %s %s: ", icon, field)
		if status.Available {
			if status.SampleValue != "" {
				fmt.Printf("%s (count: %d)\n", status.SampleValue, status.Count)
			} else {
				fmt.Printf("count: %d\n", status.Count)
			}
		} else {
			fmt.Println("NOT AVAILABLE")
		}
	}

	// Print errors if any
	if len(capture.Errors) > 0 {
		fmt.Println("\n⚠️  ERRORS:")
		for _, e := range capture.Errors {
			fmt.Printf("   - %s\n", e)
		}
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("  Data capture complete!")
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

func analyzeLoLSeriesState(state *grid.SeriesState, capture *DataCapture) {
	fmt.Println("\n   📋 LoL Series Data Analysis:")
	
	for _, game := range state.Games {
		fmt.Printf("      Game %d: Duration=%ds, Finished=%v\n", game.Sequence, game.Duration, game.Finished)
		
		// Check draft actions
		if len(game.DraftActions) > 0 {
			capture.FieldAvailability["lol_draftActions"] = FieldStatus{
				Available:   true,
				Count:       len(game.DraftActions),
				SampleValue: fmt.Sprintf("%s %s", game.DraftActions[0].Action, game.DraftActions[0].CharacterName),
			}
			fmt.Printf("      Draft Actions: %d\n", len(game.DraftActions))
		} else {
			capture.FieldAvailability["lol_draftActions"] = FieldStatus{Available: false}
		}

		for _, team := range game.Teams {
			fmt.Printf("      Team: %s (Won: %v, Kills: %d)\n", team.Name, team.Won, team.Kills)
			
			// Check team objectives
			if len(team.Objectives) > 0 {
				capture.FieldAvailability["lol_teamObjectives"] = FieldStatus{
					Available:   true,
					Count:       len(team.Objectives),
					SampleValue: team.Objectives[0].Type,
				}
			}
			
			// Check loadout value
			if team.LoadoutValue > 0 {
				capture.FieldAvailability["lol_teamLoadoutValue"] = FieldStatus{
					Available:   true,
					SampleValue: fmt.Sprintf("%d", team.LoadoutValue),
				}
			}

			for _, player := range team.Players {
				// Check multikills
				if len(player.Multikills) > 0 {
					total := 0
					for _, mk := range player.Multikills {
						total += mk.Count
					}
					capture.FieldAvailability["lol_multikills"] = FieldStatus{
						Available:   true,
						Count:       total,
						SampleValue: fmt.Sprintf("%d-kills: %d", player.Multikills[0].NumberOfKills, player.Multikills[0].Count),
					}
				}

				// Check assist details
				if len(player.AssistDetails) > 0 {
					capture.FieldAvailability["lol_assistDetails"] = FieldStatus{
						Available:   true,
						Count:       len(player.AssistDetails),
						SampleValue: fmt.Sprintf("from %s: %d", player.AssistDetails[0].PlayerID, player.AssistDetails[0].AssistsReceived),
					}
				}

				// Check abilities
				if len(player.Abilities) > 0 {
					capture.FieldAvailability["lol_abilities"] = FieldStatus{
						Available:   true,
						Count:       len(player.Abilities),
						SampleValue: player.Abilities[0].ID,
					}
				}

				// Check items
				if len(player.Items) > 0 {
					capture.FieldAvailability["lol_items"] = FieldStatus{
						Available:   true,
						Count:       len(player.Items),
						SampleValue: player.Items[0].Name,
					}
				}

				// Check player objectives
				if len(player.Objectives) > 0 {
					capture.FieldAvailability["lol_playerObjectives"] = FieldStatus{
						Available:   true,
						Count:       len(player.Objectives),
						SampleValue: player.Objectives[0].Type,
					}
				}

				// Check structures destroyed
				if player.StructuresDestroyed > 0 {
					capture.FieldAvailability["lol_structuresDestroyed"] = FieldStatus{
						Available:   true,
						Count:       player.StructuresDestroyed,
					}
				}

				// Check net worth
				if player.NetWorth > 0 {
					capture.FieldAvailability["lol_netWorth"] = FieldStatus{
						Available:   true,
						SampleValue: fmt.Sprintf("%d", player.NetWorth),
					}
				}
			}
		}
	}
}

func analyzeVALSeriesState(state *grid.SeriesState, capture *DataCapture) {
	fmt.Println("\n   📋 VALORANT Series Data Analysis:")
	
	for _, game := range state.Games {
		fmt.Printf("      Game %d: Map=%s, Duration=%ds, Finished=%v\n", game.Sequence, game.Map, game.Duration, game.Finished)
		
		// Check segments (rounds)
		if len(game.Segments) > 0 {
			capture.FieldAvailability["val_segments"] = FieldStatus{
				Available:   true,
				Count:       len(game.Segments),
				SampleValue: fmt.Sprintf("Round %d", game.Segments[0].SequenceNumber),
			}
			fmt.Printf("      Segments (Rounds): %d\n", len(game.Segments))
		} else {
			capture.FieldAvailability["val_segments"] = FieldStatus{Available: false}
		}

		for _, team := range game.Teams {
			fmt.Printf("      Team: %s (Won: %v, Score: %d)\n", team.Name, team.Won, team.Score)
			
			// Check team objectives
			if len(team.Objectives) > 0 {
				capture.FieldAvailability["val_teamObjectives"] = FieldStatus{
					Available:   true,
					Count:       len(team.Objectives),
					SampleValue: team.Objectives[0].Type,
				}
			}
			
			// Check loadout value
			if team.LoadoutValue > 0 {
				capture.FieldAvailability["val_teamLoadoutValue"] = FieldStatus{
					Available:   true,
					SampleValue: fmt.Sprintf("%d", team.LoadoutValue),
				}
			}

			for _, player := range team.Players {
				// Check weapon kills
				if len(player.WeaponKills) > 0 {
					capture.FieldAvailability["val_weaponKills"] = FieldStatus{
						Available:   true,
						Count:       len(player.WeaponKills),
						SampleValue: fmt.Sprintf("%s: %d", player.WeaponKills[0].WeaponName, player.WeaponKills[0].Count),
					}
				}

				// Check assist details
				if len(player.AssistDetails) > 0 {
					capture.FieldAvailability["val_assistDetails"] = FieldStatus{
						Available:   true,
						Count:       len(player.AssistDetails),
						SampleValue: fmt.Sprintf("from %s: %d", player.AssistDetails[0].PlayerID, player.AssistDetails[0].AssistsReceived),
					}
				}

				// Check abilities
				if len(player.Abilities) > 0 {
					capture.FieldAvailability["val_abilities"] = FieldStatus{
						Available:   true,
						Count:       len(player.Abilities),
						SampleValue: player.Abilities[0].ID,
					}
				}

				// Check player objectives
				if len(player.Objectives) > 0 {
					capture.FieldAvailability["val_playerObjectives"] = FieldStatus{
						Available:   true,
						Count:       len(player.Objectives),
						SampleValue: player.Objectives[0].Type,
					}
				}

				// Check money
				if player.Money > 0 {
					capture.FieldAvailability["val_money"] = FieldStatus{
						Available:   true,
						SampleValue: fmt.Sprintf("%d", player.Money),
					}
				}

				// Check loadout value
				if player.LoadoutValue > 0 {
					capture.FieldAvailability["val_playerLoadoutValue"] = FieldStatus{
						Available:   true,
						SampleValue: fmt.Sprintf("%d", player.LoadoutValue),
					}
				}
			}
		}
	}
}

func printLoLTeamAnalysis(analysis *intelligence.TeamAnalysis) {
	if analysis.LoLMetrics == nil {
		return
	}
	m := analysis.LoLMetrics
	fmt.Println("\n   📈 LoL Team Metrics:")
	fmt.Printf("      First Blood Rate: %.1f%%\n", m.FirstBloodRate*100)
	fmt.Printf("      First Dragon Rate: %.1f%%\n", m.FirstDragonRate*100)
	fmt.Printf("      First Tower Rate: %.1f%%\n", m.FirstTowerRate*100)
	fmt.Printf("      Baron Control Rate: %.1f%%\n", m.BaronControlRate*100)
	fmt.Printf("      Avg Game Duration: %.1f min\n", m.AvgGameDuration)
	fmt.Printf("      Early Game Rating: %.1f\n", m.EarlyGameRating)
	fmt.Printf("      Aggression Score: %.1f\n", m.AggressionScore)
}

func printVALTeamAnalysis(analysis *intelligence.TeamAnalysis) {
	if analysis.VALMetrics == nil {
		return
	}
	m := analysis.VALMetrics
	fmt.Println("\n   📈 VALORANT Team Metrics:")
	fmt.Printf("      Attack Win Rate: %.1f%%\n", m.AttackWinRate*100)
	fmt.Printf("      Defense Win Rate: %.1f%%\n", m.DefenseWinRate*100)
	fmt.Printf("      Pistol Win Rate: %.1f%%\n", m.PistolWinRate*100)
	fmt.Printf("      First Blood Rate: %.1f%%\n", m.FirstBloodRate*100)
	fmt.Printf("      Eco Round Win Rate: %.1f%%\n", m.EcoRoundWinRate*100)
	fmt.Printf("      Force Buy Win Rate: %.1f%%\n", m.ForceBuyWinRate*100)
	fmt.Printf("      Full Buy Win Rate: %.1f%%\n", m.FullBuyWinRate*100)
	fmt.Printf("      Avg Team Loadout: %.0f\n", m.AvgTeamLoadout)
	fmt.Printf("      Aggression Score: %.1f\n", m.AggressionScore)
	
	if len(m.MapPool) > 0 {
		fmt.Println("      Map Pool:")
		for _, mp := range m.MapPool {
			fmt.Printf("         %s: %.1f%% WR (%d games) - %s\n", mp.MapName, mp.WinRate*100, mp.GamesPlayed, mp.Strength)
		}
	}
}

func printPlayerProfiles(game string, players []*intelligence.PlayerProfile) {
	fmt.Printf("\n   👤 %s Player Profiles:\n", game)
	
	for _, p := range players {
		fmt.Printf("\n      %s (%s) - Threat: %d/10\n", p.Nickname, p.Role, p.ThreatLevel)
		fmt.Printf("         KDA: %.2f (%.1f/%.1f/%.1f)\n", p.KDA, p.AvgKills, p.AvgDeaths, p.AvgAssists)
		
		if len(p.SignaturePicks) > 0 {
			fmt.Printf("         Signature Picks: %v\n", p.SignaturePicks)
		}
		
		// Multikill stats
		if p.MultikillStats != nil && p.MultikillStats.TotalMultikills > 0 {
			fmt.Printf("         Multikills: %d total (D:%d T:%d Q:%d P:%d)\n", 
				p.MultikillStats.TotalMultikills,
				p.MultikillStats.DoubleKills,
				p.MultikillStats.TripleKills,
				p.MultikillStats.QuadraKills,
				p.MultikillStats.PentaKills)
		}
		
		// Weapon stats (VALORANT)
		if len(p.WeaponStats) > 0 {
			fmt.Printf("         Top Weapons: ")
			for i, ws := range p.WeaponStats {
				if i >= 3 {
					break
				}
				fmt.Printf("%s(%.0f%%) ", ws.WeaponName, ws.KillShare*100)
			}
			fmt.Println()
		}
		
		// Synergy partners
		if len(p.SynergyPartners) > 0 {
			fmt.Printf("         Top Synergy: %s (%.0f%% of assists)\n", 
				p.SynergyPartners[0].PlayerID, p.SynergyPartners[0].SynergyScore*100)
		}
		
		// Assist ratio
		if p.AssistRatio > 0 {
			fmt.Printf("         Assist Ratio: %.2f\n", p.AssistRatio)
		}
		
		// Objective focus (LoL)
		if p.ObjectiveFocus != nil {
			fmt.Printf("         Objective Focus: %s (%.1f towers/game)\n", 
				p.ObjectiveFocus.ObjectiveFocusType, p.ObjectiveFocus.TowersPerGame)
		}
		
		// Item builds (LoL)
		if len(p.ItemBuilds) > 0 {
			fmt.Printf("         Core Items: ")
			for i, item := range p.ItemBuilds {
				if i >= 3 {
					break
				}
				fmt.Printf("%s(%.0f%%) ", item.ItemName, item.BuildRate*100)
			}
			fmt.Println()
		}
		
		// Ability usage
		if len(p.AbilityUsage) > 0 {
			fmt.Printf("         Top Abilities: ")
			for i, ab := range p.AbilityUsage {
				if i >= 3 {
					break
				}
				fmt.Printf("%s(%.1f/g) ", ab.AbilityID, ab.UsagePerGame)
			}
			fmt.Println()
		}
	}
}
