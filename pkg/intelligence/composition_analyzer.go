package intelligence

import (
	"context"
	"sort"

	"scout9/pkg/grid"
)

// CompositionAnalyzer analyzes team compositions and draft patterns
type CompositionAnalyzer struct{}

// NewCompositionAnalyzer creates a new composition analyzer
func NewCompositionAnalyzer() *CompositionAnalyzer {
	return &CompositionAnalyzer{}
}

// AnalyzeCompositions analyzes team compositions from match data
func (a *CompositionAnalyzer) AnalyzeCompositions(
	ctx context.Context,
	teamID string,
	title string,
	seriesStates []*grid.SeriesState,
) (*CompositionAnalysis, error) {
	analysis := &CompositionAnalysis{
		TeamID:              teamID,
		Title:               title,
		TopCompositions:     make([]Composition, 0),
		FirstPickPriorities: make([]DraftPriority, 0),
		CommonBans:          make([]DraftPriority, 0),
		FlexPicks:           make([]string, 0),
		ArchetypeBreakdown:  make(map[string]float64),
		PlayerSynergies:     make([]Synergy, 0),
	}

	if len(seriesStates) == 0 {
		return analysis, nil
	}

	// Track compositions
	compTracker := make(map[string]*compAggregator)
	// Track character picks
	charPicks := make(map[string]*pickAggregator)
	// Track player-character synergies
	playerChars := make(map[string]map[string]*synergyAggregator)

	totalGames := 0

	for _, series := range seriesStates {
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

			// Build composition key
			chars := make([]string, 0, len(ourTeam.Players))
			for _, player := range ourTeam.Players {
				if player.Character != "" {
					chars = append(chars, player.Character)

					// Track character picks
					if _, exists := charPicks[player.Character]; !exists {
						charPicks[player.Character] = &pickAggregator{name: player.Character}
					}
					charPicks[player.Character].games++
					if ourTeam.Won {
						charPicks[player.Character].wins++
					}

					// Track player-character synergy
					if _, exists := playerChars[player.Name]; !exists {
						playerChars[player.Name] = make(map[string]*synergyAggregator)
					}
					if _, exists := playerChars[player.Name][player.Character]; !exists {
						playerChars[player.Name][player.Character] = &synergyAggregator{
							playerName:      player.Name,
							character:       player.Character,
							assistsReceived: make(map[string]int),
						}
					}
					syn := playerChars[player.Name][player.Character]
					syn.games++
					syn.kills += player.Kills
					syn.deaths += player.Deaths
					syn.assists += player.Assists
					syn.assistsGiven += player.AssistsGiven
					
					// NEW: Track assist network for synergy analysis
					for _, assist := range player.AssistDetails {
						syn.assistsReceived[assist.PlayerID] += assist.AssistsReceived
					}
					
					if ourTeam.Won {
						syn.wins++
					}
				}
			}

			// Sort for consistent key
			sort.Strings(chars)
			compKey := joinChars(chars)

			if _, exists := compTracker[compKey]; !exists {
				compTracker[compKey] = &compAggregator{
					characters: chars,
				}
			}
			compTracker[compKey].games++
			if ourTeam.Won {
				compTracker[compKey].wins++
			}
		}
	}

	// Build top compositions
	for _, agg := range compTracker {
		comp := Composition{
			Characters:  agg.characters,
			GamesPlayed: agg.games,
			Frequency:   float64(agg.games) / float64(totalGames),
		}
		if agg.games > 0 {
			comp.WinRate = float64(agg.wins) / float64(agg.games)
		}
		if title == "lol" {
			comp.Archetype = classifyLoLArchetype(agg.characters)
		}
		analysis.TopCompositions = append(analysis.TopCompositions, comp)
	}

	// Sort by frequency
	sort.Slice(analysis.TopCompositions, func(i, j int) bool {
		return analysis.TopCompositions[i].GamesPlayed > analysis.TopCompositions[j].GamesPlayed
	})

	// Limit to top 5
	if len(analysis.TopCompositions) > 5 {
		analysis.TopCompositions = analysis.TopCompositions[:5]
	}

	// Build first pick priorities (most picked characters)
	for charName, agg := range charPicks {
		priority := DraftPriority{
			Character:   charName,
			Rate:        float64(agg.games) / float64(totalGames),
			GamesPlayed: agg.games,
		}
		if agg.games > 0 {
			priority.WinRate = float64(agg.wins) / float64(agg.games)
		}
		analysis.FirstPickPriorities = append(analysis.FirstPickPriorities, priority)
	}

	// Sort by pick rate
	sort.Slice(analysis.FirstPickPriorities, func(i, j int) bool {
		return analysis.FirstPickPriorities[i].Rate > analysis.FirstPickPriorities[j].Rate
	})

	// Limit to top 10
	if len(analysis.FirstPickPriorities) > 10 {
		analysis.FirstPickPriorities = analysis.FirstPickPriorities[:10]
	}

	// Build player synergies
	for playerName, chars := range playerChars {
		for charName, agg := range chars {
			if agg.games >= 2 { // Only include if played at least twice
				syn := Synergy{
					PlayerName:  playerName,
					Character:   charName,
					GamesPlayed: agg.games,
				}
				if agg.games > 0 {
					syn.WinRate = float64(agg.wins) / float64(agg.games)
					if agg.deaths > 0 {
						syn.KDA = float64(agg.kills+agg.assists) / float64(agg.deaths)
					} else {
						syn.KDA = float64(agg.kills + agg.assists)
					}
				}
				analysis.PlayerSynergies = append(analysis.PlayerSynergies, syn)
			}
		}
	}

	// Sort synergies by games played
	sort.Slice(analysis.PlayerSynergies, func(i, j int) bool {
		return analysis.PlayerSynergies[i].GamesPlayed > analysis.PlayerSynergies[j].GamesPlayed
	})

	// Calculate archetype breakdown for LoL
	if title == "lol" {
		archetypeCounts := make(map[string]int)
		for _, comp := range analysis.TopCompositions {
			if comp.Archetype != "" {
				archetypeCounts[comp.Archetype] += comp.GamesPlayed
			}
		}
		for archetype, count := range archetypeCounts {
			analysis.ArchetypeBreakdown[archetype] = float64(count) / float64(totalGames)
		}
	}

	// Identify flex picks (characters played by multiple players)
	charPlayers := make(map[string]map[string]bool)
	for playerName, chars := range playerChars {
		for charName := range chars {
			if _, exists := charPlayers[charName]; !exists {
				charPlayers[charName] = make(map[string]bool)
			}
			charPlayers[charName][playerName] = true
		}
	}
	for charName, players := range charPlayers {
		if len(players) >= 2 {
			analysis.FlexPicks = append(analysis.FlexPicks, charName)
		}
	}

	return analysis, nil
}

// Helper types
type compAggregator struct {
	characters []string
	games      int
	wins       int
}

type pickAggregator struct {
	name  string
	games int
	wins  int
}

type synergyAggregator struct {
	playerName string
	character  string
	games      int
	wins       int
	kills      int
	deaths     int
	assists    int
	// NEW: Track assist network for this player-character combo
	assistsReceived map[string]int // playerID -> assists received
	assistsGiven    int
}

func joinChars(chars []string) string {
	result := ""
	for i, c := range chars {
		if i > 0 {
			result += "|"
		}
		result += c
	}
	return result
}

// classifyLoLArchetype uses the real champion database to classify composition archetype
func classifyLoLArchetype(characters []string) string {
	// Use the comprehensive champion database for accurate classification
	return ClassifyLoLCompositionArchetype(characters)
}
