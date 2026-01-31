package grid

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/machinebox/graphql"
)

// SeriesStateCacheTTL is the cache duration for series state data
const SeriesStateCacheTTL = 24 * time.Hour

// GetSeriesState fetches detailed match state for a series
func (c *Client) GetSeriesState(ctx context.Context, seriesID string) (*SeriesState, error) {
	cacheKey := fmt.Sprintf("series:state:%s", seriesID)

	// Check cache
	var state SeriesState
	if c.getCached(ctx, cacheKey, &state) {
		return &state, nil
	}

	// Enhanced query with ALL available fields for rich data analysis
	// Verified available fields:
	// - killAssistsReceivedFromPlayer: Assist network analysis (who assisted on kills)
	// - loadoutValue: Team/player economy tracking
	// - structuresDestroyed: Objective control stats
	// - abilities: VALORANT agent ability usage (IDs only)
	// - teamkills/selfkills: Error tracking (usually 0)
	// - segments: Round-by-round data for VALORANT
	req := graphql.NewRequest(`
		query SeriesState($seriesId: ID!) {
			seriesState(id: $seriesId) {
				id
				started
				finished
				format
				teams {
					id
					name
					won
					score
					kills
					deaths
				}
				games {
					id
					started
					finished
					paused
					map {
						name
					}
					clock {
						currentSeconds
					}
					draftActions {
						id
						type
						sequenceNumber
						drafter {
							id
						}
						draftable {
							id
							name
						}
					}
					segments {
						id
						sequenceNumber
						type
						finished
						teams {
							id
							name
							side
							won
							kills
							deaths
							objectives {
								id
								type
								completionCount
							}
							players {
								id
								name
								kills
								deaths
								objectives {
									id
									type
									completionCount
								}
							}
						}
					}
					teams {
						id
						name
						side
						score
						won
						kills
						deaths
						netWorth
						money
						loadoutValue
						structuresDestroyed
						objectives {
							id
							type
							completionCount
						}
						players {
							id
							name
							character {
								id
								name
							}
							kills
							deaths
							killAssistsReceived
							killAssistsGiven
							killAssistsReceivedFromPlayer {
								playerId
								killAssistsReceived
							}
							teamkills
							selfkills
							netWorth
							money
							loadoutValue
							structuresDestroyed
							objectives {
								id
								type
								completionCount
							}
							multikills {
								id
								numberOfKills
								count
							}
							weaponKills {
								id
								weaponName
								count
							}
							abilities {
								id
							}
							inventory {
								items {
									id
									name
								}
							}
						}
					}
				}
			}
		}
	`)
	req.Var("seriesId", seriesID)

	var resp struct {
		SeriesState struct {
			ID       string `json:"id"`
			Started  bool   `json:"started"`
			Finished bool   `json:"finished"`
			Format   string `json:"format"`
			Teams    []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Won    bool   `json:"won"`
				Score  int    `json:"score"`
				Kills  int    `json:"kills"`
				Deaths int    `json:"deaths"`
			} `json:"teams"`
			Games []struct {
				ID       string `json:"id"`
				Started  bool   `json:"started"`
				Finished bool   `json:"finished"`
				Paused   bool   `json:"paused"`
				Map      struct {
					Name string `json:"name"`
				} `json:"map"`
				Clock struct {
					CurrentSeconds int `json:"currentSeconds"`
				} `json:"clock"`
				DraftActions []struct {
					ID             string `json:"id"`
					Type           string `json:"type"`
					SequenceNumber string `json:"sequenceNumber"`
					Drafter        struct {
						ID string `json:"id"`
					} `json:"drafter"`
					Draftable struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"draftable"`
				} `json:"draftActions"`
				Segments []struct {
					ID             string `json:"id"`
					SequenceNumber int    `json:"sequenceNumber"`
					Type           string `json:"type"`
					Finished       bool   `json:"finished"`
					Teams          []struct {
						ID     string `json:"id"`
						Name   string `json:"name"`
						Side   string `json:"side"`
						Won    bool   `json:"won"`
						Kills  int    `json:"kills"`
						Deaths int    `json:"deaths"`
						Objectives []struct {
							ID              string `json:"id"`
							Type            string `json:"type"`
							CompletionCount int    `json:"completionCount"`
						} `json:"objectives"`
						Players []struct {
							ID     string `json:"id"`
							Name   string `json:"name"`
							Kills  int    `json:"kills"`
							Deaths int    `json:"deaths"`
							Objectives []struct {
								ID              string `json:"id"`
								Type            string `json:"type"`
								CompletionCount int    `json:"completionCount"`
							} `json:"objectives"`
						} `json:"players"`
					} `json:"teams"`
				} `json:"segments"`
				Teams []struct {
					ID                  string `json:"id"`
					Name                string `json:"name"`
					Side                string `json:"side"`
					Score               int    `json:"score"`
					Won                 bool   `json:"won"`
					Kills               int    `json:"kills"`
					Deaths              int    `json:"deaths"`
					NetWorth            int    `json:"netWorth"`
					Money               int    `json:"money"`
					LoadoutValue        int    `json:"loadoutValue"`
					StructuresDestroyed int    `json:"structuresDestroyed"`
					Objectives          []struct {
						ID              string `json:"id"`
						Type            string `json:"type"`
						CompletionCount int    `json:"completionCount"`
					} `json:"objectives"`
					Players []struct {
						ID        string `json:"id"`
						Name      string `json:"name"`
						Character struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"character"`
						Kills               int `json:"kills"`
						Deaths              int `json:"deaths"`
						KillAssistsReceived int `json:"killAssistsReceived"`
						KillAssistsGiven    int `json:"killAssistsGiven"`
						KillAssistsReceivedFromPlayer []struct {
							PlayerID            string `json:"playerId"`
							KillAssistsReceived int    `json:"killAssistsReceived"`
						} `json:"killAssistsReceivedFromPlayer"`
						Teamkills           int `json:"teamkills"`
						Selfkills           int `json:"selfkills"`
						NetWorth            int `json:"netWorth"`
						Money               int `json:"money"`
						LoadoutValue        int `json:"loadoutValue"`
						StructuresDestroyed int `json:"structuresDestroyed"`
						Objectives          []struct {
							ID              string `json:"id"`
							Type            string `json:"type"`
							CompletionCount int    `json:"completionCount"`
						} `json:"objectives"`
						Multikills []struct {
							ID            string `json:"id"`
							NumberOfKills int    `json:"numberOfKills"`
							Count         int    `json:"count"`
						} `json:"multikills"`
						WeaponKills []struct {
							ID         string `json:"id"`
							WeaponName string `json:"weaponName"`
							Count      int    `json:"count"`
						} `json:"weaponKills"`
						// VALORANT-specific: abilities used by player
						Abilities []struct {
							ID string `json:"id"`
						} `json:"abilities"`
						Inventory struct {
							Items []struct {
								ID   string `json:"id"`
								Name string `json:"name"`
							} `json:"items"`
						} `json:"inventory"`
					} `json:"players"`
				} `json:"teams"`
			} `json:"games"`
		} `json:"seriesState"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runSeriesStateQuery(ctx, seriesID, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get series state: %w", err)
	}

	// Convert response to SeriesState
	state = SeriesState{
		ID:       resp.SeriesState.ID,
		Started:  resp.SeriesState.Started,
		Finished: resp.SeriesState.Finished,
	}

	// Convert teams
	for _, t := range resp.SeriesState.Teams {
		state.Teams = append(state.Teams, TeamState{
			ID:     t.ID,
			Name:   t.Name,
			Won:    t.Won,
			Score:  t.Score,
			Kills:  t.Kills,
			Deaths: t.Deaths,
		})
	}

	// Convert games
	for i, g := range resp.SeriesState.Games {
		game := Game{
			ID:       g.ID,
			Sequence: i + 1, // Use index as sequence since sequenceNumber is not available
			Map:      g.Map.Name,
			Duration: g.Clock.CurrentSeconds,
			Finished: g.Finished,
			Started:  g.Started,
			Paused:   g.Paused,
		}

		// Convert draft actions
		for _, da := range g.DraftActions {
			game.DraftActions = append(game.DraftActions, DraftAction{
				TeamID:        da.Drafter.ID,
				Action:        da.Type,
				CharacterName: da.Draftable.Name,
				CharacterID:   da.Draftable.ID,
			})
		}

		// Convert segments (rounds for VALORANT)
		for _, seg := range g.Segments {
			segment := Segment{
				ID:             seg.ID,
				SequenceNumber: seg.SequenceNumber,
				Type:           seg.Type,
				Finished:       seg.Finished,
			}

			for _, st := range seg.Teams {
				segTeam := SegmentTeam{
					ID:     st.ID,
					Name:   st.Name,
					Side:   st.Side,
					Won:    st.Won,
					Kills:  st.Kills,
					Deaths: st.Deaths,
				}

				// Convert segment team objectives
				for _, obj := range st.Objectives {
					segTeam.Objectives = append(segTeam.Objectives, TeamObjective{
						ID:              obj.ID,
						Type:            obj.Type,
						CompletionCount: obj.CompletionCount,
					})
				}

				// Convert segment players
				for _, sp := range st.Players {
					segPlayer := SegmentPlayer{
						ID:     sp.ID,
						Name:   sp.Name,
						Kills:  sp.Kills,
						Deaths: sp.Deaths,
					}

					// Convert segment player objectives
					for _, obj := range sp.Objectives {
						segPlayer.Objectives = append(segPlayer.Objectives, PlayerObjective{
							ID:              obj.ID,
							Type:            obj.Type,
							CompletionCount: obj.CompletionCount,
						})
					}

					segTeam.Players = append(segTeam.Players, segPlayer)
				}

				segment.Teams = append(segment.Teams, segTeam)
			}

			game.Segments = append(game.Segments, segment)
		}

		for _, gt := range g.Teams {
			gameTeam := GameTeam{
				ID:                  gt.ID,
				Name:                gt.Name,
				Side:                gt.Side,
				Score:               gt.Score,
				Won:                 gt.Won,
				Kills:               gt.Kills,
				Deaths:              gt.Deaths,
				NetWorth:            gt.NetWorth,
				Money:               gt.Money,
				LoadoutValue:        gt.LoadoutValue,
				StructuresDestroyed: gt.StructuresDestroyed,
			}

			// Convert team objectives
			for _, obj := range gt.Objectives {
				gameTeam.Objectives = append(gameTeam.Objectives, TeamObjective{
					ID:              obj.ID,
					Type:            obj.Type,
					CompletionCount: obj.CompletionCount,
				})
			}

			for _, p := range gt.Players {
				player := GamePlayer{
					ID:                  p.ID,
					Name:                p.Name,
					Character:           p.Character.Name,
					Kills:               p.Kills,
					Deaths:              p.Deaths,
					Assists:             p.KillAssistsReceived,
					AssistsGiven:        p.KillAssistsGiven,
					Teamkills:           p.Teamkills,
					Selfkills:           p.Selfkills,
					NetWorth:            p.NetWorth,
					Money:               p.Money,
					LoadoutValue:        p.LoadoutValue,
					StructuresDestroyed: p.StructuresDestroyed,
				}

				// Convert assist details (who assisted on kills)
				for _, assist := range p.KillAssistsReceivedFromPlayer {
					player.AssistDetails = append(player.AssistDetails, AssistDetail{
						PlayerID:        assist.PlayerID,
						AssistsReceived: assist.KillAssistsReceived,
					})
				}

				// Convert player items
				for _, item := range p.Inventory.Items {
					player.Items = append(player.Items, Item{
						ID:   item.ID,
						Name: item.Name,
					})
				}

				// Convert player objectives
				for _, obj := range p.Objectives {
					player.Objectives = append(player.Objectives, PlayerObjective{
						ID:              obj.ID,
						Type:            obj.Type,
						CompletionCount: obj.CompletionCount,
					})
				}

				// Convert multikills
				for _, mk := range p.Multikills {
					player.Multikills = append(player.Multikills, Multikill{
						ID:            mk.ID,
						NumberOfKills: mk.NumberOfKills,
						Count:         mk.Count,
					})
				}

				// Convert weapon kills
				for _, wk := range p.WeaponKills {
					player.WeaponKills = append(player.WeaponKills, WeaponKill{
						ID:         wk.ID,
						WeaponName: wk.WeaponName,
						Count:      wk.Count,
					})
				}

				// Convert abilities (VALORANT - just IDs)
				for _, ab := range p.Abilities {
					player.Abilities = append(player.Abilities, AbilityUsage{
						ID:          ab.ID,
						AbilityName: ab.ID, // Use ID as name since that's all we get
					})
				}

				gameTeam.Players = append(gameTeam.Players, player)
			}

			game.Teams = append(game.Teams, gameTeam)
		}

		state.Games = append(state.Games, game)
	}

	// Cache result
	c.setCache(ctx, cacheKey, state, SeriesStateCacheTTL)

	return &state, nil
}

// GetSeriesStates fetches multiple series states concurrently
func (c *Client) GetSeriesStates(ctx context.Context, seriesIDs []string) ([]*SeriesState, error) {
	results := make([]*SeriesState, len(seriesIDs))
	errors := make([]error, len(seriesIDs))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Limit concurrent requests

	for i, id := range seriesIDs {
		wg.Add(1)
		go func(idx int, seriesID string) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			state, err := c.GetSeriesState(ctx, seriesID)
			if err != nil {
				errors[idx] = err
				return
			}
			results[idx] = state
		}(i, id)
	}

	wg.Wait()

	// Collect non-nil results (skip failed fetches)
	var validResults []*SeriesState
	for i, result := range results {
		if result != nil {
			validResults = append(validResults, result)
		} else if errors[i] != nil {
			// Log error but continue with available data
			fmt.Printf("Warning: failed to fetch series %s: %v\n", seriesIDs[i], errors[i])
		}
	}

	return validResults, nil
}

// GetMatchDataForTeam fetches all match data for a team's recent series
func (c *Client) GetMatchDataForTeam(ctx context.Context, teamID string, matchCount int) ([]*SeriesState, error) {
	// First get the series IDs
	series, err := c.GetSeriesForTeam(ctx, teamID, matchCount)
	if err != nil {
		return nil, fmt.Errorf("get series for team: %w", err)
	}

	// Extract series IDs
	seriesIDs := make([]string, len(series))
	for i, s := range series {
		seriesIDs[i] = s.ID
	}

	// Fetch all series states concurrently
	return c.GetSeriesStates(ctx, seriesIDs)
}
