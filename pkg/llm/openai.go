package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService implements LLM service using OpenAI API
type OpenAIService struct {
	client *openai.Client
	model  string
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
		model:  openai.GPT4o,
	}
}

// GenerateExecutiveSummary creates an executive summary from analysis data
func (s *OpenAIService) GenerateExecutiveSummary(ctx context.Context, data *AnalysisData) (string, error) {
	gameType := "League of Legends"
	if data.Title == "valorant" || data.Title == "25" {
		gameType = "VALORANT"
	}

	prompt := fmt.Sprintf(`You are an expert esports analyst creating a scouting report executive summary.

Team: %s
Game: %s
Matches Analyzed: %d
Win Rate: %.1f%%

Strengths:
%s

Weaknesses:
%s

Key Statistics:
%s

Player Highlights:
%s

Write a concise executive summary (2-3 paragraphs) that:
1. Summarizes the team's overall playstyle and current form
2. Highlights their key strengths and how they typically win
3. Identifies their main vulnerabilities

Use professional coaching terminology. Be specific and actionable.`,
		data.TeamName,
		gameType,
		data.MatchesAnalyzed,
		data.WinRate*100,
		formatStrings(data.Strengths),
		formatStrings(data.Weaknesses),
		formatKeyStats(data.KeyStats),
		formatPlayerHighlights(data.PlayerHighlights),
	)

	return s.complete(ctx, prompt)
}

// GenerateCounterStrategyNarrative creates a narrative for counter-strategies
func (s *OpenAIService) GenerateCounterStrategyNarrative(ctx context.Context, strategy *CounterStrategyData) (string, error) {
	gameType := "League of Legends"
	if strategy.Title == "valorant" || strategy.Title == "25" {
		gameType = "VALORANT"
	}

	prompt := fmt.Sprintf(`You are an expert esports coach creating a "How to Win" strategy guide.

Opponent: %s
Game: %s

Exploitable Weaknesses:
%s

Draft Recommendations:
%s

In-Game Strategies:
%s

Write a compelling "How to Win" narrative (2-3 paragraphs) that:
1. Clearly states the win condition against this opponent
2. Explains the key strategic approach
3. Provides specific, actionable recommendations

Be direct and confident. Use coaching language that motivates the team.`,
		strategy.OpponentName,
		gameType,
		formatWeaknesses(strategy.Weaknesses),
		formatDraftTargets(strategy.DraftTargets),
		formatStrings(strategy.InGameStrategies),
	)

	return s.complete(ctx, prompt)
}

// ExplainStatisticalPattern explains a statistical pattern in natural language
func (s *OpenAIService) ExplainStatisticalPattern(ctx context.Context, pattern *PatternData) (string, error) {
	prompt := fmt.Sprintf(`You are an esports analyst explaining a statistical finding.

Pattern: %s
Description: %s
Value: %.2f
Comparison: %s
Context: %s

Write a brief (1-2 sentences) explanation of what this statistic means for coaching strategy. Be specific and actionable.`,
		pattern.Title,
		pattern.Description,
		pattern.Value,
		pattern.Comparison,
		pattern.Context,
	)

	return s.complete(ctx, prompt)
}

// GenerateWinCondition creates a concise win condition statement
func (s *OpenAIService) GenerateWinCondition(ctx context.Context, data *AnalysisData) (string, error) {
	prompt := fmt.Sprintf(`You are an expert esports coach. Based on this opponent analysis, write a single, powerful sentence that captures the win condition.

Team: %s
Weaknesses: %s
Key Stats: %s

Write ONE sentence starting with "To beat %s," that captures the essential strategy. Be specific and actionable.`,
		data.TeamName,
		formatStrings(data.Weaknesses),
		formatKeyStats(data.KeyStats),
		data.TeamName,
	)

	return s.complete(ctx, prompt)
}

// complete sends a completion request to OpenAI
func (s *OpenAIService) complete(ctx context.Context, prompt string) (string, error) {
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert esports analyst and coach specializing in League of Legends and VALORANT. You provide concise, actionable insights for professional teams.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   500,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("openai completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// Helper functions for formatting
func formatStrings(items []string) string {
	if len(items) == 0 {
		return "- None identified"
	}
	var sb strings.Builder
	for _, item := range items {
		sb.WriteString("- ")
		sb.WriteString(item)
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatKeyStats(stats map[string]interface{}) string {
	if len(stats) == 0 {
		return "- No key stats available"
	}
	var sb strings.Builder
	for key, value := range stats {
		sb.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}
	return sb.String()
}

func formatPlayerHighlights(players []PlayerHighlight) string {
	if len(players) == 0 {
		return "- No player highlights"
	}
	var sb strings.Builder
	for _, p := range players {
		sb.WriteString(fmt.Sprintf("- %s (%s): Threat Level %d/10, Top picks: %s. %s\n",
			p.Name, p.Role, p.ThreatLevel, strings.Join(p.TopPicks, ", "), p.KeyStat))
	}
	return sb.String()
}

func formatWeaknesses(weaknesses []WeaknessData) string {
	if len(weaknesses) == 0 {
		return "- No clear weaknesses identified"
	}
	var sb strings.Builder
	for _, w := range weaknesses {
		sb.WriteString(fmt.Sprintf("- %s (Impact: %.0f%%)\n  Evidence: %s\n",
			w.Description, w.Impact*100, w.Evidence))
	}
	return sb.String()
}

func formatDraftTargets(targets []DraftTarget) string {
	if len(targets) == 0 {
		return "- No specific draft recommendations"
	}
	var sb strings.Builder
	for _, t := range targets {
		sb.WriteString(fmt.Sprintf("- [%s] %s (Priority %d): %s\n",
			strings.ToUpper(t.Type), t.Character, t.Priority, t.Reason))
	}
	return sb.String()
}
