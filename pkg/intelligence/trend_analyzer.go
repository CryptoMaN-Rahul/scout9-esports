package intelligence

import (
	"context"
	"sort"
	"time"

	"scout9/pkg/grid"
)

// TrendAnalyzer analyzes performance trends over time
type TrendAnalyzer struct{}

// NewTrendAnalyzer creates a new trend analyzer
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{}
}

// AnalyzeTrends analyzes team performance trends
func (a *TrendAnalyzer) AnalyzeTrends(
	ctx context.Context,
	teamID string,
	seriesStates []*grid.SeriesState,
	seriesInfo []grid.Series,
) (*TrendAnalysis, error) {
	analysis := &TrendAnalysis{
		TeamID:       teamID,
		TrendChanges: make([]TrendChange, 0),
		MetricTrends: make([]MetricTrend, 0),
	}

	if len(seriesStates) == 0 {
		return analysis, nil
	}

	// Sort series by time (need series info for timestamps)
	type seriesWithTime struct {
		state *grid.SeriesState
		time  time.Time
	}

	seriesMap := make(map[string]time.Time)
	for _, s := range seriesInfo {
		seriesMap[s.ID] = s.StartTime
	}

	sortedSeries := make([]seriesWithTime, 0, len(seriesStates))
	for _, state := range seriesStates {
		t := seriesMap[state.ID]
		if t.IsZero() {
			t = time.Now() // Fallback
		}
		sortedSeries = append(sortedSeries, seriesWithTime{state: state, time: t})
	}

	sort.Slice(sortedSeries, func(i, j int) bool {
		return sortedSeries[i].time.Before(sortedSeries[j].time)
	})

	// Calculate win rates over time
	var (
		totalGames    int
		totalWins     int
		last5Games    []bool
		last10Games   []bool
		winRateValues []float64
	)

	for _, sw := range sortedSeries {
		for _, game := range sw.state.Games {
			if !game.Finished {
				continue
			}

			// Find our team
			var won bool
			for _, team := range game.Teams {
				if team.ID == teamID {
					won = team.Won
					break
				}
			}

			totalGames++
			if won {
				totalWins++
			}

			// Track recent games
			last5Games = append(last5Games, won)
			if len(last5Games) > 5 {
				last5Games = last5Games[1:]
			}

			last10Games = append(last10Games, won)
			if len(last10Games) > 10 {
				last10Games = last10Games[1:]
			}

			// Track win rate progression
			if totalGames > 0 {
				winRateValues = append(winRateValues, float64(totalWins)/float64(totalGames))
			}
		}
	}

	// Calculate rates
	if totalGames > 0 {
		analysis.OverallWinRate = float64(totalWins) / float64(totalGames)
	}

	if len(last5Games) > 0 {
		wins := 0
		for _, w := range last5Games {
			if w {
				wins++
			}
		}
		analysis.Last5WinRate = float64(wins) / float64(len(last5Games))
	}

	if len(last10Games) > 0 {
		wins := 0
		for _, w := range last10Games {
			if w {
				wins++
			}
		}
		analysis.Last10WinRate = float64(wins) / float64(len(last10Games))
	}

	// Determine form indicator
	analysis.FormIndicator, analysis.FormScore = calculateForm(analysis.Last5WinRate, analysis.OverallWinRate)

	// Add win rate trend
	if len(winRateValues) > 0 {
		analysis.MetricTrends = append(analysis.MetricTrends, MetricTrend{
			Metric:    "Win Rate",
			Direction: calculateTrendDirection(winRateValues),
			Values:    winRateValues,
		})
	}

	// Detect significant trend changes
	if len(winRateValues) >= 5 {
		// Compare first half to second half
		mid := len(winRateValues) / 2
		firstHalfAvg := averageFloat64(winRateValues[:mid])
		secondHalfAvg := averageFloat64(winRateValues[mid:])

		diff := secondHalfAvg - firstHalfAvg
		if diff > 0.15 {
			analysis.TrendChanges = append(analysis.TrendChanges, TrendChange{
				Metric:      "Win Rate",
				OldValue:    firstHalfAvg,
				NewValue:    secondHalfAvg,
				Explanation: "Significant improvement in recent matches",
			})
		} else if diff < -0.15 {
			analysis.TrendChanges = append(analysis.TrendChanges, TrendChange{
				Metric:      "Win Rate",
				OldValue:    firstHalfAvg,
				NewValue:    secondHalfAvg,
				Explanation: "Performance decline in recent matches",
			})
		}
	}

	return analysis, nil
}

func calculateForm(recentRate, overallRate float64) (string, float64) {
	diff := recentRate - overallRate

	// Form score: -100 to 100
	formScore := diff * 200 // Scale to -100 to 100 range

	if formScore > 100 {
		formScore = 100
	}
	if formScore < -100 {
		formScore = -100
	}

	// Form indicator
	if recentRate >= 0.7 || diff > 0.15 {
		return "hot", formScore
	} else if recentRate <= 0.3 || diff < -0.15 {
		return "cold", formScore
	}
	return "stable", formScore
}

func calculateTrendDirection(values []float64) string {
	if len(values) < 3 {
		return "stable"
	}

	// Simple linear regression slope
	n := float64(len(values))
	var sumX, sumY, sumXY, sumX2 float64

	for i, v := range values {
		x := float64(i)
		sumX += x
		sumY += v
		sumXY += x * v
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	if slope > 0.01 {
		return "improving"
	} else if slope < -0.01 {
		return "declining"
	}
	return "stable"
}

func averageFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
