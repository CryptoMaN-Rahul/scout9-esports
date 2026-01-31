package grid

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
)

// Cache TTLs
const (
	TitlesCacheTTL      = 24 * time.Hour
	TournamentsCacheTTL = 24 * time.Hour
	TeamsCacheTTL       = 12 * time.Hour
	SeriesCacheTTL      = 6 * time.Hour
)

// GetTitles fetches available game titles (LoL, VALORANT)
func (c *Client) GetTitles(ctx context.Context) ([]Title, error) {
	cacheKey := "titles:all"

	// Check cache
	var titles []Title
	if c.getCached(ctx, cacheKey, &titles) {
		return titles, nil
	}

	// GraphQL query
	req := graphql.NewRequest(`
		query Titles {
			titles {
				id
				name
			}
		}
	`)

	var resp struct {
		Titles []Title `json:"titles"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get titles: %w", err)
	}

	// Cache result
	c.setCache(ctx, cacheKey, resp.Titles, TitlesCacheTTL)

	return resp.Titles, nil
}

// GetTournaments fetches tournaments for a specific title
func (c *Client) GetTournaments(ctx context.Context, titleID string) ([]Tournament, error) {
	cacheKey := fmt.Sprintf("tournaments:%s", titleID)

	// Check cache
	var tournaments []Tournament
	if c.getCached(ctx, cacheKey, &tournaments) {
		return tournaments, nil
	}

	// GraphQL query
	req := graphql.NewRequest(`
		query Tournaments($titleId: [ID!]) {
			tournaments(filter: { title: { id: { in: $titleId } } }, first: 100) {
				totalCount
				edges {
					node {
						id
						name
						logoUrl
						startDate
						endDate
					}
				}
			}
		}
	`)
	req.Var("titleId", []string{titleID})

	var resp struct {
		Tournaments struct {
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node struct {
					ID        string `json:"id"`
					Name      string `json:"name"`
					LogoURL   string `json:"logoUrl"`
					StartDate string `json:"startDate"`
					EndDate   string `json:"endDate"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"tournaments"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get tournaments: %w", err)
	}

	// Convert response to Tournament slice
	tournaments = make([]Tournament, 0, len(resp.Tournaments.Edges))
	for _, edge := range resp.Tournaments.Edges {
		t := Tournament{
			ID:      edge.Node.ID,
			Name:    edge.Node.Name,
			LogoURL: edge.Node.LogoURL,
			TitleID: titleID,
		}
		if edge.Node.StartDate != "" {
			t.StartDate, _ = time.Parse(time.RFC3339, edge.Node.StartDate)
		}
		if edge.Node.EndDate != "" {
			t.EndDate, _ = time.Parse(time.RFC3339, edge.Node.EndDate)
		}
		tournaments = append(tournaments, t)
	}

	// Cache result
	c.setCache(ctx, cacheKey, tournaments, TournamentsCacheTTL)

	return tournaments, nil
}

// GetTeams fetches teams that have participated in a tournament
func (c *Client) GetTeams(ctx context.Context, tournamentID string) ([]Team, error) {
	cacheKey := fmt.Sprintf("teams:tournament:%s", tournamentID)

	// Check cache
	var teams []Team
	if c.getCached(ctx, cacheKey, &teams) {
		return teams, nil
	}

	// First get series for the tournament, then extract unique teams
	// Note: GRID API expects tournament ID as [ID!] array
	req := graphql.NewRequest(`
		query TeamsInTournament($tournamentId: [ID!]!) {
			allSeries(
				filter: { tournament: { id: { in: $tournamentId }, includeChildren: { equals: true } } }
				first: 50
			) {
				edges {
					node {
						teams {
							baseInfo {
								id
								name
								logoUrl
							}
						}
					}
				}
			}
		}
	`)
	req.Var("tournamentId", []string{tournamentID})

	var resp struct {
		AllSeries struct {
			Edges []struct {
				Node struct {
					Teams []struct {
						BaseInfo struct {
							ID      string `json:"id"`
							Name    string `json:"name"`
							LogoURL string `json:"logoUrl"`
						} `json:"baseInfo"`
					} `json:"teams"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"allSeries"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get teams: %w", err)
	}

	// Extract unique teams
	teamMap := make(map[string]Team)
	for _, edge := range resp.AllSeries.Edges {
		for _, t := range edge.Node.Teams {
			if t.BaseInfo.ID != "" {
				teamMap[t.BaseInfo.ID] = Team{
					ID:      t.BaseInfo.ID,
					Name:    t.BaseInfo.Name,
					LogoURL: t.BaseInfo.LogoURL,
				}
			}
		}
	}

	teams = make([]Team, 0, len(teamMap))
	for _, team := range teamMap {
		teams = append(teams, team)
	}

	// Cache result
	c.setCache(ctx, cacheKey, teams, TeamsCacheTTL)

	return teams, nil
}

// GetSeriesForTeam fetches recent series for a specific team
// Uses GRID API's native teamId filter for direct team filtering
func (c *Client) GetSeriesForTeam(ctx context.Context, teamID string, limit int) ([]Series, error) {
	if limit <= 0 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("series:team:%s:limit:%d", teamID, limit)

	// Check cache
	var series []Series
	if c.getCached(ctx, cacheKey, &series) {
		return series, nil
	}

	// Query series directly using GRID API's teamId filter
	// This is the proper approach - no hardcoded tournament IDs needed
	req := graphql.NewRequest(`
		query SeriesForTeam($teamId: ID!, $limit: Int!) {
			allSeries(
				filter: { teamId: $teamId }
				orderBy: StartTimeScheduled
				orderDirection: DESC
				first: $limit
			) {
				totalCount
				edges {
					node {
						id
						startTimeScheduled
						format {
							nameShortened
						}
						teams {
							baseInfo {
								id
								name
								logoUrl
							}
						}
						tournament {
							id
							name
						}
						title {
							id
							name
						}
					}
				}
			}
		}
	`)
	req.Var("teamId", teamID)
	req.Var("limit", limit)

	var resp struct {
		AllSeries struct {
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node struct {
					ID                 string `json:"id"`
					StartTimeScheduled string `json:"startTimeScheduled"`
					Format             struct {
						NameShortened string `json:"nameShortened"`
					} `json:"format"`
					Teams []struct {
						BaseInfo struct {
							ID      string `json:"id"`
							Name    string `json:"name"`
							LogoURL string `json:"logoUrl"`
						} `json:"baseInfo"`
					} `json:"teams"`
					Tournament struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"tournament"`
					Title struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"title"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"allSeries"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get series for team: %w", err)
	}

	// Convert response to Series slice
	series = make([]Series, 0, len(resp.AllSeries.Edges))
	for _, edge := range resp.AllSeries.Edges {
		s := Series{
			ID:           edge.Node.ID,
			TournamentID: edge.Node.Tournament.ID,
			Format:       edge.Node.Format.NameShortened,
			TitleID:      edge.Node.Title.ID,
		}
		if edge.Node.StartTimeScheduled != "" {
			s.StartTime, _ = time.Parse(time.RFC3339, edge.Node.StartTimeScheduled)
		}
		for _, t := range edge.Node.Teams {
			s.Teams = append(s.Teams, Team{
				ID:      t.BaseInfo.ID,
				Name:    t.BaseInfo.Name,
				LogoURL: t.BaseInfo.LogoURL,
			})
		}
		series = append(series, s)
	}

	// Cache result
	c.setCache(ctx, cacheKey, series, SeriesCacheTTL)

	return series, nil
}

// SearchTeams searches for teams by name across all tournaments
func (c *Client) SearchTeams(ctx context.Context, query string, titleID string) ([]Team, error) {
	req := graphql.NewRequest(`
		query SearchTeams($query: String!, $titleId: ID) {
			teams(
				filter: { 
					name: { contains: $query }
					titleId: $titleId
				}
				first: 20
			) {
				edges {
					node {
						id
						name
						logoUrl
					}
				}
			}
		}
	`)
	req.Var("query", query)
	if titleID != "" {
		req.Var("titleId", titleID)
	}

	var resp struct {
		Teams struct {
			Edges []struct {
				Node struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					LogoURL string `json:"logoUrl"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"teams"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("search teams: %w", err)
	}

	teams := make([]Team, 0, len(resp.Teams.Edges))
	for _, edge := range resp.Teams.Edges {
		teams = append(teams, Team{
			ID:      edge.Node.ID,
			Name:    edge.Node.Name,
			LogoURL: edge.Node.LogoURL,
		})
	}

	return teams, nil
}

// GetTeamByID fetches a single team by ID
// Uses allSeries query since direct team(id:) query is not available in hackathon API
func (c *Client) GetTeamByID(ctx context.Context, teamID string) (*Team, error) {
	cacheKey := fmt.Sprintf("team:%s", teamID)

	// Check cache
	var team Team
	if c.getCached(ctx, cacheKey, &team) {
		return &team, nil
	}

	// Use allSeries with teamId filter to get team info
	// This is the approved approach for the hackathon API
	req := graphql.NewRequest(`
		query TeamFromSeries($teamId: ID!) {
			allSeries(
				filter: { teamId: $teamId }
				first: 1
			) {
				edges {
					node {
						teams {
							baseInfo {
								id
								name
								logoUrl
							}
						}
					}
				}
			}
		}
	`)
	req.Var("teamId", teamID)

	var resp struct {
		AllSeries struct {
			Edges []struct {
				Node struct {
					Teams []struct {
						BaseInfo struct {
							ID      string `json:"id"`
							Name    string `json:"name"`
							LogoURL string `json:"logoUrl"`
						} `json:"baseInfo"`
					} `json:"teams"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"allSeries"`
	}

	err := withRetry(ctx, 3, func() error {
		return c.runCentralDataQuery(ctx, req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("get team from series: %w", err)
	}

	// Extract team info from series response
	for _, edge := range resp.AllSeries.Edges {
		for _, t := range edge.Node.Teams {
			if t.BaseInfo.ID == teamID {
				team = Team{
					ID:      t.BaseInfo.ID,
					Name:    t.BaseInfo.Name,
					LogoURL: t.BaseInfo.LogoURL,
				}
				// Cache result
				c.setCache(ctx, cacheKey, team, TeamsCacheTTL)
				return &team, nil
			}
		}
	}

	// If team not found in series, return a basic team struct with just the ID
	// This allows the report to proceed with the team name from the request
	return &Team{ID: teamID}, nil
}
