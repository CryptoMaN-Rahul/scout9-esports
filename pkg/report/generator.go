package report

import (
	"context"
	"fmt"
	"time"

	"scout9/pkg/grid"
	"scout9/pkg/intelligence"

	"github.com/google/uuid"
)

// Generator orchestrates the scouting report generation
type Generator struct {
	gridClient          *grid.Client
	lolAnalyzer         *intelligence.LoLAnalyzer
	valAnalyzer         *intelligence.VALAnalyzer
	compositionAnalyzer *intelligence.CompositionAnalyzer
	trendAnalyzer       *intelligence.TrendAnalyzer
	counterEngine       *intelligence.CounterStrategyEngine
	formatter           *Formatter
}

// NewGenerator creates a new report generator
func NewGenerator(gridClient *grid.Client) *Generator {
	return &Generator{
		gridClient:          gridClient,
		lolAnalyzer:         intelligence.NewLoLAnalyzer(),
		valAnalyzer:         intelligence.NewVALAnalyzer(),
		compositionAnalyzer: intelligence.NewCompositionAnalyzer(),
		trendAnalyzer:       intelligence.NewTrendAnalyzer(),
		counterEngine:       intelligence.NewCounterStrategyEngine(),
		formatter:           NewFormatter(),
	}
}

// GenerateRequest contains the parameters for report generation
type GenerateRequest struct {
	TeamID     string
	TeamName   string
	TitleID    string // "3" for LoL, "6" for VALORANT
	MatchCount int
}

// GenerateReport generates a complete scouting report for a team
func (g *Generator) GenerateReport(ctx context.Context, req GenerateRequest) (*intelligence.ScoutingReport, error) {
	// Determine game title
	title := "lol"
	if req.TitleID == "6" {
		title = "valorant"
	}

	// Step 1: Get team info
	team, err := g.gridClient.GetTeamByID(ctx, req.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team info: %w", err)
	}

	teamName := req.TeamName
	if teamName == "" && team != nil {
		teamName = team.Name
	}

	// Step 2: Get recent series for the team
	seriesLimit := req.MatchCount
	if seriesLimit <= 0 {
		seriesLimit = 10
	}

	seriesList, err := g.gridClient.GetSeriesForTeam(ctx, req.TeamID, seriesLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	if len(seriesList) == 0 {
		return nil, fmt.Errorf("no matches found for team %s", req.TeamID)
	}

	// Step 3: Get series states (detailed match data)
	seriesIDs := make([]string, len(seriesList))
	for i, s := range seriesList {
		seriesIDs[i] = s.ID
	}

	seriesStates, err := g.gridClient.GetSeriesStates(ctx, seriesIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get series states: %w", err)
	}

	// Step 4: Download and parse event files for detailed analysis
	var lolEvents map[string]*grid.LoLEventData
	var valEvents map[string]*grid.VALEventData

	if title == "lol" {
		lolEvents = make(map[string]*grid.LoLEventData)
		for _, seriesID := range seriesIDs {
			events, err := g.gridClient.DownloadEvents(ctx, seriesID)
			if err != nil {
				// Log but continue - events are optional enhancement
				continue
			}
			parsed, err := g.gridClient.ParseLoLEvents(events)
			if err != nil {
				continue
			}
			lolEvents[seriesID] = parsed
		}
	} else {
		valEvents = make(map[string]*grid.VALEventData)
		for _, seriesID := range seriesIDs {
			events, err := g.gridClient.DownloadEvents(ctx, seriesID)
			if err != nil {
				continue
			}
			parsed, err := g.gridClient.ParseVALEvents(events)
			if err != nil {
				continue
			}
			valEvents[seriesID] = parsed
		}
	}

	// Step 5: Run all analyzers
	var teamAnalysis *intelligence.TeamAnalysis
	var playerProfiles []*intelligence.PlayerProfile

	if title == "lol" {
		teamAnalysis, err = g.lolAnalyzer.AnalyzeTeam(ctx, req.TeamID, teamName, seriesStates, lolEvents)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze team: %w", err)
		}
		playerProfiles, err = g.lolAnalyzer.AnalyzePlayers(ctx, req.TeamID, seriesStates)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze players: %w", err)
		}
	} else {
		teamAnalysis, err = g.valAnalyzer.AnalyzeTeam(ctx, req.TeamID, teamName, seriesStates, valEvents)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze team: %w", err)
		}
		playerProfiles, err = g.valAnalyzer.AnalyzePlayers(ctx, req.TeamID, seriesStates)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze players: %w", err)
		}
	}

	// Composition analysis
	compositions, err := g.compositionAnalyzer.AnalyzeCompositions(ctx, req.TeamID, title, seriesStates)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze compositions: %w", err)
	}

	// Trend analysis
	trends, err := g.trendAnalyzer.AnalyzeTrends(ctx, req.TeamID, seriesStates, seriesList)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze trends: %w", err)
	}

	// Step 6: Generate counter-strategy (THE KEY DIFFERENTIATOR)
	// Use enhanced method with timing, matchup, and site analysis for hackathon-winning insights
	var counterStrategy *intelligence.CounterStrategy
	if title == "lol" {
		counterStrategy = g.counterEngine.GenerateEnhancedCounterStrategy(
			teamAnalysis, playerProfiles, compositions,
			seriesStates, lolEvents, nil,
		)
	} else {
		counterStrategy = g.counterEngine.GenerateEnhancedCounterStrategy(
			teamAnalysis, playerProfiles, compositions,
			seriesStates, nil, valEvents,
		)
	}

	// Step 7: Build the final report
	report := &intelligence.ScoutingReport{
		ID:          uuid.New().String(),
		GeneratedAt: time.Now(),
		OpponentTeam: intelligence.TeamInfo{
			ID:   req.TeamID,
			Name: teamName,
		},
		Title:           title,
		MatchesAnalyzed: teamAnalysis.MatchesAnalyzed,
		HowToWin:        counterStrategy,
		TeamStrategy:    teamAnalysis,
		PlayerProfiles:  playerProfiles,
		Compositions:    compositions,
		TrendAnalysis:   trends,
	}

	// Add logo if available
	if team != nil {
		report.OpponentTeam.LogoURL = team.LogoURL
	}

	// Generate executive summary
	report.ExecutiveSummary = g.generateExecutiveSummary(report)

	return report, nil
}

// GenerateDigestibleReport generates a hackathon-compliant DigestibleReport
// This is the PRIMARY output format for the hackathon submission
func (g *Generator) GenerateDigestibleReport(ctx context.Context, req GenerateRequest) (*intelligence.DigestibleReport, error) {
	// First generate the full scouting report
	report, err := g.GenerateReport(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to hackathon-compliant format
	return g.formatter.FormatDigestibleReport(report), nil
}

// GenerateTextReport generates a plain text report in hackathon format
func (g *Generator) GenerateTextReport(ctx context.Context, req GenerateRequest) (string, error) {
	digestible, err := g.GenerateDigestibleReport(ctx, req)
	if err != nil {
		return "", err
	}

	return g.formatter.FormatTextReport(digestible), nil
}

// generateExecutiveSummary creates a brief summary of the report
func (g *Generator) generateExecutiveSummary(report *intelligence.ScoutingReport) string {
	summary := fmt.Sprintf("%s Analysis (%d matches analyzed)\n\n", report.OpponentTeam.Name, report.MatchesAnalyzed)

	// Win rate and form
	if report.TeamStrategy != nil {
		summary += fmt.Sprintf("Overall Win Rate: %.0f%%\n", report.TeamStrategy.WinRate*100)
	}

	if report.TrendAnalysis != nil {
		summary += fmt.Sprintf("Current Form: %s (Last 5: %.0f%%)\n\n",
			report.TrendAnalysis.FormIndicator, report.TrendAnalysis.Last5WinRate*100)
	}

	// Key strengths
	if report.TeamStrategy != nil && len(report.TeamStrategy.Strengths) > 0 {
		summary += "Key Strengths:\n"
		for i, s := range report.TeamStrategy.Strengths {
			if i >= 3 {
				break
			}
			summary += fmt.Sprintf("• %s\n", s.Title)
		}
		summary += "\n"
	}

	// Key weaknesses
	if report.TeamStrategy != nil && len(report.TeamStrategy.Weaknesses) > 0 {
		summary += "Key Weaknesses:\n"
		for i, w := range report.TeamStrategy.Weaknesses {
			if i >= 3 {
				break
			}
			summary += fmt.Sprintf("• %s\n", w.Title)
		}
		summary += "\n"
	}

	// How to win
	if report.HowToWin != nil && report.HowToWin.WinCondition != "" {
		summary += fmt.Sprintf("HOW TO WIN: %s\n", report.HowToWin.WinCondition)
	}

	return summary
}

// GenerateMatchupReport generates a head-to-head analysis between two teams
func (g *Generator) GenerateMatchupReport(
	ctx context.Context,
	team1ID, team2ID string,
	titleID string,
	matchCount int,
) (*intelligence.MatchupAnalysis, error) {
	// Get team info
	team1, err := g.gridClient.GetTeamByID(ctx, team1ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team1 info: %w", err)
	}

	team2, err := g.gridClient.GetTeamByID(ctx, team2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team2 info: %w", err)
	}

	analysis := &intelligence.MatchupAnalysis{
		Team1ID:         team1ID,
		Team1Name:       team1.Name,
		Team2ID:         team2ID,
		Team2Name:       team2.Name,
		Patterns:        make([]intelligence.MatchupPattern, 0),
		DraftPatterns:   make([]intelligence.DraftPattern, 0),
		Recommendations: make([]string, 0),
	}

	// Get series for team1
	series1, err := g.gridClient.GetSeriesForTeam(ctx, team1ID, matchCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get team1 series: %w", err)
	}

	// Find head-to-head matches
	for _, s := range series1 {
		for _, t := range s.Teams {
			if t.ID == team2ID {
				// This is a head-to-head match
				analysis.TotalMatches++

				// Get series state to determine winner
				state, err := g.gridClient.GetSeriesState(ctx, s.ID)
				if err != nil {
					continue
				}

				for _, team := range state.Teams {
					if team.ID == team1ID && team.Won {
						analysis.Team1Wins++
					} else if team.ID == team2ID && team.Won {
						analysis.Team2Wins++
					}
				}
			}
		}
	}

	// Generate recommendations based on head-to-head
	if analysis.TotalMatches > 0 {
		team1WinRate := float64(analysis.Team1Wins) / float64(analysis.TotalMatches)
		if team1WinRate > 0.6 {
			analysis.Recommendations = append(analysis.Recommendations,
				fmt.Sprintf("%s has historically dominated this matchup (%.0f%% win rate)", team1.Name, team1WinRate*100))
		} else if team1WinRate < 0.4 {
			analysis.Recommendations = append(analysis.Recommendations,
				fmt.Sprintf("%s has struggled in this matchup (%.0f%% win rate) - consider adjusting strategy", team1.Name, team1WinRate*100))
		} else {
			analysis.Recommendations = append(analysis.Recommendations,
				"This is an even matchup historically - preparation will be key")
		}
	} else {
		analysis.Recommendations = append(analysis.Recommendations,
			"No recent head-to-head matches found - focus on general team analysis")
	}

	return analysis, nil
}
