package intelligence

import (
	"context"
	"fmt"

	"scout9/pkg/grid"
)

// HeadToHeadAnalyzer compares two teams
type HeadToHeadAnalyzer struct {
	gridClient *grid.Client
}

// NewHeadToHeadAnalyzer creates a new head-to-head analyzer
func NewHeadToHeadAnalyzer(client *grid.Client) *HeadToHeadAnalyzer {
	return &HeadToHeadAnalyzer{
		gridClient: client,
	}
}

// AnalyzeMatchup analyzes historical matchup between two teams
func (h *HeadToHeadAnalyzer) AnalyzeMatchup(
	ctx context.Context,
	team1ID string,
	team2ID string,
	titleID string,
	matchCount int,
) (*HeadToHeadReport, error) {
	if matchCount <= 0 {
		matchCount = 20
	}

	// Get team info
	team1, err := h.gridClient.GetTeamByID(ctx, team1ID)
	if err != nil {
		return nil, fmt.Errorf("get team1: %w", err)
	}

	team2, err := h.gridClient.GetTeamByID(ctx, team2ID)
	if err != nil {
		return nil, fmt.Errorf("get team2: %w", err)
	}

	report := &HeadToHeadReport{
		Team1ID:   team1ID,
		Team1Name: team1.Name,
		Team2ID:   team2ID,
		Team2Name: team2.Name,
		Title:     titleID,
		Insights:  make([]HeadToHeadInsight, 0),
		Warnings:  make([]string, 0),
	}

	// Get series for both teams
	team1Series, err := h.gridClient.GetSeriesForTeam(ctx, team1ID, matchCount)
	if err != nil {
		return nil, fmt.Errorf("get team1 series: %w", err)
	}

	team2Series, err := h.gridClient.GetSeriesForTeam(ctx, team2ID, matchCount)
	if err != nil {
		return nil, fmt.Errorf("get team2 series: %w", err)
	}

	// Find head-to-head matches (now fetches winner info from Series State API)
	h2hMatches := h.findHeadToHeadMatches(ctx, team1Series, team2Series, team1ID, team2ID)

	report.TotalMatches = len(h2hMatches)

	// Count wins
	for _, match := range h2hMatches {
		if match.winnerID == team1ID {
			report.Team1Wins++
		} else if match.winnerID == team2ID {
			report.Team2Wins++
		}
	}

	// Generate historical insight
	if report.TotalMatches > 0 {
		report.Insights = append(report.Insights, HeadToHeadInsight{
			Text: fmt.Sprintf("Historical record: %s leads %d-%d against %s",
				getLeader(report), report.Team1Wins, report.Team2Wins, getTrailer(report)),
			Type:       "historical",
			Confidence: calculateHistoricalConfidence(report.TotalMatches),
		})
	} else {
		report.Warnings = append(report.Warnings, "No historical matches found between these teams")
	}

	// Calculate confidence score
	report.ConfidenceScore = h.calculateReportConfidence(report)

	return report, nil
}

// CompareStyles compares team styles even without historical matches
func (h *HeadToHeadAnalyzer) CompareStyles(
	team1Analysis *TeamAnalysis,
	team2Analysis *TeamAnalysis,
) *StyleComparison {
	comparison := &StyleComparison{}

	// Compare based on title
	if team1Analysis.Title == "lol" && team1Analysis.LoLMetrics != nil && team2Analysis.LoLMetrics != nil {
		comparison = h.compareLoLStyles(team1Analysis, team2Analysis)
	} else if team1Analysis.Title == "valorant" && team1Analysis.VALMetrics != nil && team2Analysis.VALMetrics != nil {
		comparison = h.compareVALStyles(team1Analysis, team2Analysis)
	}

	return comparison
}

// compareLoLStyles compares LoL team styles
func (h *HeadToHeadAnalyzer) compareLoLStyles(team1 *TeamAnalysis, team2 *TeamAnalysis) *StyleComparison {
	m1 := team1.LoLMetrics
	m2 := team2.LoLMetrics

	comparison := &StyleComparison{
		Team1EarlyGameRating: m1.EarlyGameRating,
		Team2EarlyGameRating: m2.EarlyGameRating,
		Team1MidGameRating:   m1.MidGameRating,
		Team2MidGameRating:   m2.MidGameRating,
		Team1LateGameRating:  m1.LateGameRating,
		Team2LateGameRating:  m2.LateGameRating,
		Team1Aggression:      m1.AggressionScore,
		Team2Aggression:      m2.AggressionScore,
	}

	// Determine advantages
	comparison.EarlyGameAdvantage = determineAdvantage(m1.EarlyGameRating, m2.EarlyGameRating)
	comparison.MidGameAdvantage = determineAdvantage(m1.MidGameRating, m2.MidGameRating)
	comparison.LateGameAdvantage = determineAdvantage(m1.LateGameRating, m2.LateGameRating)

	// Generate insights
	comparison.EarlyGameInsight = fmt.Sprintf("%s has %.0f early game rating vs %s's %.0f",
		team1.TeamName, m1.EarlyGameRating, team2.TeamName, m2.EarlyGameRating)
	comparison.MidGameInsight = fmt.Sprintf("%s has %.0f mid game rating vs %s's %.0f",
		team1.TeamName, m1.MidGameRating, team2.TeamName, m2.MidGameRating)
	comparison.LateGameInsight = fmt.Sprintf("%s has %.0f late game rating vs %s's %.0f",
		team1.TeamName, m1.LateGameRating, team2.TeamName, m2.LateGameRating)

	// Overall style insight
	comparison.StyleInsight = generateLoLStyleInsight(team1, team2, comparison)

	return comparison
}

// compareVALStyles compares VALORANT team styles
func (h *HeadToHeadAnalyzer) compareVALStyles(team1 *TeamAnalysis, team2 *TeamAnalysis) *StyleComparison {
	m1 := team1.VALMetrics
	m2 := team2.VALMetrics

	// For VALORANT, use attack/defense rates as game phase proxies
	comparison := &StyleComparison{
		Team1EarlyGameRating: m1.PistolWinRate * 100,
		Team2EarlyGameRating: m2.PistolWinRate * 100,
		Team1MidGameRating:   m1.AttackWinRate * 100,
		Team2MidGameRating:   m2.AttackWinRate * 100,
		Team1LateGameRating:  m1.DefenseWinRate * 100,
		Team2LateGameRating:  m2.DefenseWinRate * 100,
		Team1Aggression:      m1.AggressionScore,
		Team2Aggression:      m2.AggressionScore,
	}

	// Determine advantages
	comparison.EarlyGameAdvantage = determineAdvantage(m1.PistolWinRate, m2.PistolWinRate)
	comparison.MidGameAdvantage = determineAdvantage(m1.AttackWinRate, m2.AttackWinRate)
	comparison.LateGameAdvantage = determineAdvantage(m1.DefenseWinRate, m2.DefenseWinRate)

	// Generate insights
	comparison.EarlyGameInsight = fmt.Sprintf("%s wins %.0f%% of pistol rounds vs %s's %.0f%%",
		team1.TeamName, m1.PistolWinRate*100, team2.TeamName, m2.PistolWinRate*100)
	comparison.MidGameInsight = fmt.Sprintf("%s has %.0f%% attack win rate vs %s's %.0f%%",
		team1.TeamName, m1.AttackWinRate*100, team2.TeamName, m2.AttackWinRate*100)
	comparison.LateGameInsight = fmt.Sprintf("%s has %.0f%% defense win rate vs %s's %.0f%%",
		team1.TeamName, m1.DefenseWinRate*100, team2.TeamName, m2.DefenseWinRate*100)

	// Overall style insight
	comparison.StyleInsight = generateVALStyleInsight(team1, team2, comparison)

	return comparison
}

// GenerateMatchupInsights generates insights from style comparison
func (h *HeadToHeadAnalyzer) GenerateMatchupInsights(
	report *HeadToHeadReport,
	team1Analysis *TeamAnalysis,
	team2Analysis *TeamAnalysis,
) []HeadToHeadInsight {
	insights := make([]HeadToHeadInsight, 0)

	// Add style comparison
	styleComp := h.CompareStyles(team1Analysis, team2Analysis)
	report.StyleComparison = styleComp

	// Early game advantage insight
	if styleComp.EarlyGameAdvantage != "even" {
		advantageTeam := report.Team1Name
		if styleComp.EarlyGameAdvantage == "team2" {
			advantageTeam = report.Team2Name
		}
		insights = append(insights, HeadToHeadInsight{
			Text:       fmt.Sprintf("%s has the early game advantage", advantageTeam),
			Type:       "early_game",
			Confidence: 0.7,
		})
	}

	// Aggression mismatch insight
	aggressionDiff := styleComp.Team1Aggression - styleComp.Team2Aggression
	if aggressionDiff > 15 {
		insights = append(insights, HeadToHeadInsight{
			Text:       fmt.Sprintf("%s plays more aggressively (%.0f vs %.0f)", report.Team1Name, styleComp.Team1Aggression, styleComp.Team2Aggression),
			Type:       "style",
			Confidence: 0.6,
		})
	} else if aggressionDiff < -15 {
		insights = append(insights, HeadToHeadInsight{
			Text:       fmt.Sprintf("%s plays more aggressively (%.0f vs %.0f)", report.Team2Name, styleComp.Team2Aggression, styleComp.Team1Aggression),
			Type:       "style",
			Confidence: 0.6,
		})
	}

	return insights
}

// Helper types and functions

type h2hMatch struct {
	seriesID string
	winnerID string
}

func (h *HeadToHeadAnalyzer) findHeadToHeadMatches(ctx context.Context, team1Series, team2Series []grid.Series, team1ID, team2ID string) []h2hMatch {
	matches := make([]h2hMatch, 0)

	// Create a map of team2's series for quick lookup
	team2SeriesMap := make(map[string]bool)
	for _, s := range team2Series {
		team2SeriesMap[s.ID] = true
	}

	// Collect series IDs where both teams played
	h2hSeriesIDs := make([]string, 0)
	for _, s := range team1Series {
		if !team2SeriesMap[s.ID] {
			continue
		}

		// Check if both teams are in this series
		hasTeam1 := false
		hasTeam2 := false
		for _, t := range s.Teams {
			if t.ID == team1ID {
				hasTeam1 = true
			}
			if t.ID == team2ID {
				hasTeam2 = true
			}
		}

		if hasTeam1 && hasTeam2 {
			h2hSeriesIDs = append(h2hSeriesIDs, s.ID)
		}
	}

	// Fetch series states to get winners
	if len(h2hSeriesIDs) > 0 {
		seriesStates, err := h.gridClient.GetSeriesStates(ctx, h2hSeriesIDs)
		if err == nil {
			for _, state := range seriesStates {
				if state == nil || !state.Finished {
					continue
				}
				winnerID := ""
				for _, team := range state.Teams {
					if team.Won {
						winnerID = team.ID
						break
					}
				}
				matches = append(matches, h2hMatch{
					seriesID: state.ID,
					winnerID: winnerID,
				})
			}
		} else {
			// Fallback: add matches without winner info if API fails
			for _, seriesID := range h2hSeriesIDs {
				matches = append(matches, h2hMatch{
					seriesID: seriesID,
					winnerID: "",
				})
			}
		}
	}

	return matches
}

func (h *HeadToHeadAnalyzer) calculateReportConfidence(report *HeadToHeadReport) float64 {
	confidence := 50.0 // Base confidence

	// Boost for historical matches
	if report.TotalMatches >= 5 {
		confidence += 25
	} else if report.TotalMatches >= 3 {
		confidence += 15
	} else if report.TotalMatches >= 1 {
		confidence += 5
	}

	// Boost for style comparison
	if report.StyleComparison != nil {
		confidence += 10
	}

	// Boost for insights
	confidence += float64(len(report.Insights)) * 2

	if confidence > 100 {
		confidence = 100
	}

	return confidence
}

func calculateHistoricalConfidence(matchCount int) float64 {
	if matchCount >= 5 {
		return 0.9
	} else if matchCount >= 3 {
		return 0.7
	} else if matchCount >= 1 {
		return 0.5
	}
	return 0.3
}

func determineAdvantage(rating1, rating2 float64) string {
	diff := rating1 - rating2
	if diff > 5 {
		return "team1"
	} else if diff < -5 {
		return "team2"
	}
	return "even"
}

func getLeader(report *HeadToHeadReport) string {
	if report.Team1Wins > report.Team2Wins {
		return report.Team1Name
	} else if report.Team2Wins > report.Team1Wins {
		return report.Team2Name
	}
	return report.Team1Name // Tie goes to team1
}

func getTrailer(report *HeadToHeadReport) string {
	if report.Team1Wins > report.Team2Wins {
		return report.Team2Name
	} else if report.Team2Wins > report.Team1Wins {
		return report.Team1Name
	}
	return report.Team2Name
}

func generateLoLStyleInsight(team1, team2 *TeamAnalysis, comp *StyleComparison) string {
	// Determine dominant phase for each team
	team1Phase := getDominantPhase(comp.Team1EarlyGameRating, comp.Team1MidGameRating, comp.Team1LateGameRating)
	team2Phase := getDominantPhase(comp.Team2EarlyGameRating, comp.Team2MidGameRating, comp.Team2LateGameRating)

	if team1Phase == team2Phase {
		return fmt.Sprintf("Both teams excel in %s - expect a close match in that phase", team1Phase)
	}

	return fmt.Sprintf("%s is stronger in %s while %s excels in %s",
		team1.TeamName, team1Phase, team2.TeamName, team2Phase)
}

func generateVALStyleInsight(team1, team2 *TeamAnalysis, comp *StyleComparison) string {
	m1 := team1.VALMetrics
	m2 := team2.VALMetrics

	// Compare attack vs defense focus
	team1AttackFocused := m1.AttackWinRate > m1.DefenseWinRate
	team2AttackFocused := m2.AttackWinRate > m2.DefenseWinRate

	if team1AttackFocused && !team2AttackFocused {
		return fmt.Sprintf("%s is attack-focused (%.0f%% attack WR) while %s is defense-focused (%.0f%% defense WR)",
			team1.TeamName, m1.AttackWinRate*100, team2.TeamName, m2.DefenseWinRate*100)
	} else if !team1AttackFocused && team2AttackFocused {
		return fmt.Sprintf("%s is defense-focused (%.0f%% defense WR) while %s is attack-focused (%.0f%% attack WR)",
			team1.TeamName, m1.DefenseWinRate*100, team2.TeamName, m2.AttackWinRate*100)
	}

	return fmt.Sprintf("Both teams have similar playstyles - expect a balanced match")
}

func getDominantPhase(early, mid, late float64) string {
	if early >= mid && early >= late {
		return "early game"
	} else if mid >= early && mid >= late {
		return "mid game"
	}
	return "late game"
}
