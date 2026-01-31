package llm

import (
	"context"
)

// Service defines the interface for LLM-powered insights
type Service interface {
	GenerateExecutiveSummary(ctx context.Context, data *AnalysisData) (string, error)
	GenerateCounterStrategyNarrative(ctx context.Context, strategy *CounterStrategyData) (string, error)
	ExplainStatisticalPattern(ctx context.Context, pattern *PatternData) (string, error)
	GenerateWinCondition(ctx context.Context, data *AnalysisData) (string, error)
}

// AnalysisData contains the data needed for generating insights
type AnalysisData struct {
	TeamName        string
	Title           string // "lol" or "valorant"
	MatchesAnalyzed int
	WinRate         float64
	Strengths       []string
	Weaknesses      []string
	KeyStats        map[string]interface{}
	PlayerHighlights []PlayerHighlight
}

// PlayerHighlight contains notable player information
type PlayerHighlight struct {
	Name        string
	Role        string
	ThreatLevel int
	TopPicks    []string
	KeyStat     string
}

// CounterStrategyData contains data for counter-strategy generation
type CounterStrategyData struct {
	OpponentName     string
	Title            string
	Weaknesses       []WeaknessData
	DraftTargets     []DraftTarget
	InGameStrategies []string
}

// WeaknessData describes an exploitable weakness
type WeaknessData struct {
	Description string
	Evidence    string
	Impact      float64
}

// DraftTarget represents a draft recommendation
type DraftTarget struct {
	Type      string // "ban", "pick", "target"
	Character string
	Reason    string
	Priority  int
}

// PatternData contains a statistical pattern to explain
type PatternData struct {
	Title       string
	Description string
	Value       float64
	Comparison  string
	Context     string
}
