package llm

import (
	"context"
	"fmt"
	"strings"
)

// TemplateService provides template-based text generation as fallback
type TemplateService struct{}

// NewTemplateService creates a new template service
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// GenerateExecutiveSummary creates an executive summary using templates
func (s *TemplateService) GenerateExecutiveSummary(ctx context.Context, data *AnalysisData) (string, error) {
	gameType := "League of Legends"
	if data.Title == "valorant" || data.Title == "25" {
		gameType = "VALORANT"
	}

	// Determine form description
	formDesc := "average"
	if data.WinRate >= 0.7 {
		formDesc = "excellent"
	} else if data.WinRate >= 0.55 {
		formDesc = "solid"
	} else if data.WinRate < 0.4 {
		formDesc = "struggling"
	}

	var sb strings.Builder

	// Opening paragraph
	sb.WriteString(fmt.Sprintf("%s has shown %s form in their recent %d %s matches, posting a %.0f%% win rate. ",
		data.TeamName, formDesc, data.MatchesAnalyzed, gameType, data.WinRate*100))

	// Strengths
	if len(data.Strengths) > 0 {
		sb.WriteString(fmt.Sprintf("Their key strengths include %s. ", strings.Join(data.Strengths[:min(3, len(data.Strengths))], ", ")))
	}

	sb.WriteString("\n\n")

	// Weaknesses paragraph
	if len(data.Weaknesses) > 0 {
		sb.WriteString(fmt.Sprintf("However, the team shows vulnerabilities in %s. ",
			strings.Join(data.Weaknesses[:min(3, len(data.Weaknesses))], ", ")))
		sb.WriteString("These areas present opportunities for exploitation with proper preparation.")
	}

	sb.WriteString("\n\n")

	// Player highlights
	if len(data.PlayerHighlights) > 0 {
		sb.WriteString("Key players to watch: ")
		var highlights []string
		for _, p := range data.PlayerHighlights[:min(3, len(data.PlayerHighlights))] {
			highlights = append(highlights, fmt.Sprintf("%s (%s, Threat: %d/10)", p.Name, p.Role, p.ThreatLevel))
		}
		sb.WriteString(strings.Join(highlights, ", "))
		sb.WriteString(".")
	}

	return sb.String(), nil
}

// GenerateCounterStrategyNarrative creates a counter-strategy narrative using templates
func (s *TemplateService) GenerateCounterStrategyNarrative(ctx context.Context, strategy *CounterStrategyData) (string, error) {
	var sb strings.Builder

	// Win condition opening
	sb.WriteString(fmt.Sprintf("To defeat %s, focus on exploiting their identified weaknesses. ", strategy.OpponentName))

	// Weaknesses
	if len(strategy.Weaknesses) > 0 {
		sb.WriteString("\n\nKey vulnerabilities to target:\n")
		for i, w := range strategy.Weaknesses[:min(3, len(strategy.Weaknesses))] {
			sb.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, w.Description, w.Evidence))
		}
	}

	// Draft recommendations
	if len(strategy.DraftTargets) > 0 {
		sb.WriteString("\nDraft priorities:\n")
		for _, t := range strategy.DraftTargets[:min(5, len(strategy.DraftTargets))] {
			sb.WriteString(fmt.Sprintf("- %s %s: %s\n", strings.ToUpper(t.Type), t.Character, t.Reason))
		}
	}

	// In-game strategies
	if len(strategy.InGameStrategies) > 0 {
		sb.WriteString("\nIn-game approach:\n")
		for _, strat := range strategy.InGameStrategies[:min(3, len(strategy.InGameStrategies))] {
			sb.WriteString(fmt.Sprintf("- %s\n", strat))
		}
	}

	return sb.String(), nil
}

// ExplainStatisticalPattern explains a pattern using templates
func (s *TemplateService) ExplainStatisticalPattern(ctx context.Context, pattern *PatternData) (string, error) {
	return fmt.Sprintf("%s: %.1f%% (%s). %s",
		pattern.Title, pattern.Value*100, pattern.Comparison, pattern.Context), nil
}

// GenerateWinCondition creates a win condition statement using templates
func (s *TemplateService) GenerateWinCondition(ctx context.Context, data *AnalysisData) (string, error) {
	if len(data.Weaknesses) == 0 {
		return fmt.Sprintf("To beat %s, maintain consistent execution and capitalize on any mistakes.", data.TeamName), nil
	}

	return fmt.Sprintf("To beat %s, exploit their %s while neutralizing their %s.",
		data.TeamName,
		data.Weaknesses[0],
		func() string {
			if len(data.Strengths) > 0 {
				return data.Strengths[0]
			}
			return "key players"
		}()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
