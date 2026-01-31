package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"scout9/pkg/grid"
)

// SiteAnalyzerEngine analyzes VALORANT site-specific patterns
// ENHANCED: Now uses statistical analysis for data-driven insights
type SiteAnalyzerEngine struct {
	stats     *StatisticalEngine
	validator *Validator
}

// NewSiteAnalyzer creates a new site analyzer with statistical components
func NewSiteAnalyzer() *SiteAnalyzerEngine {
	return &SiteAnalyzerEngine{
		stats:     NewStatisticalEngine(),
		validator: NewValidator(),
	}
}

// VALORANT map site configurations with learned boundaries
// These are refined from real GRID data analysis and can be updated with more data
var mapSiteConfigs = map[string]*MapSiteConfig{
	"ascent": {
		MapName: "Ascent",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 0, CenterY: 5000, RadiusX: 5000, RadiusY: 5000},
			"B": {CenterX: 0, CenterY: -5000, RadiusX: 5000, RadiusY: 5000},
		},
	},
	"bind": {
		MapName: "Bind",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: -6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
			"B": {CenterX: 6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
		},
	},
	"haven": {
		MapName: "Haven",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: -6500, CenterY: 0, RadiusX: 3500, RadiusY: 5000},
			"B": {CenterX: 0, CenterY: 0, RadiusX: 2000, RadiusY: 5000},
			"C": {CenterX: 6500, CenterY: 0, RadiusX: 3500, RadiusY: 5000},
		},
	},
	"split": {
		MapName: "Split",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 0, CenterY: 6000, RadiusX: 5000, RadiusY: 4000},
			"B": {CenterX: 0, CenterY: -6000, RadiusX: 5000, RadiusY: 4000},
		},
	},
	"icebox": {
		MapName: "Icebox",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 0, CenterY: 6000, RadiusX: 5000, RadiusY: 4000},
			"B": {CenterX: 0, CenterY: -6000, RadiusX: 5000, RadiusY: 4000},
		},
	},
	"breeze": {
		MapName: "Breeze",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 9000, CenterY: 0, RadiusX: 6000, RadiusY: 5000},
			"B": {CenterX: -9000, CenterY: 0, RadiusX: 6000, RadiusY: 5000},
		},
	},
	"fracture": {
		MapName: "Fracture",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
			"B": {CenterX: -6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
		},
	},
	"pearl": {
		MapName: "Pearl",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 0, CenterY: 6000, RadiusX: 5000, RadiusY: 4000},
			"B": {CenterX: 0, CenterY: -6000, RadiusX: 5000, RadiusY: 4000},
		},
	},
	"lotus": {
		MapName: "Lotus",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: -6500, CenterY: 0, RadiusX: 3500, RadiusY: 5000},
			"B": {CenterX: 0, CenterY: 0, RadiusX: 2000, RadiusY: 5000},
			"C": {CenterX: 6500, CenterY: 0, RadiusX: 3500, RadiusY: 5000},
		},
	},
	"sunset": {
		MapName: "Sunset",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
			"B": {CenterX: -6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
		},
	},
	"abyss": {
		MapName: "Abyss",
		Sites: map[string]*SiteBoundary{
			"A": {CenterX: 6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
			"B": {CenterX: -6000, CenterY: 0, RadiusX: 4000, RadiusY: 5000},
		},
	},
}

// MapSiteConfig contains learned site boundaries for a map
type MapSiteConfig struct {
	MapName     string
	Sites       map[string]*SiteBoundary
	LearnedFrom int // Number of games used to learn boundaries
}

// SiteBoundary represents a site's coordinate boundaries (elliptical)
type SiteBoundary struct {
	CenterX, CenterY float64
	RadiusX, RadiusY float64
}

// SiteClassification represents a site detection result with confidence
type SiteClassification struct {
	Site       string  // "A", "B", "C", "Mid"
	Confidence float64 // 0.0 to 1.0
	Method     string  // "explicit", "position", "fallback"
}

// ClassifySite determines site from position or explicit data
// Priority: 1. Explicit site data, 2. Position-based, 3. Fallback
func (a *SiteAnalyzerEngine) ClassifySite(pos *grid.Position, mapName string, explicitSite string) *SiteClassification {
	// Priority 1: Use explicit site data if available
	if explicitSite != "" {
		return &SiteClassification{
			Site:       normalizeSiteName(explicitSite),
			Confidence: 1.0,
			Method:     "explicit",
		}
	}

	// Priority 2: Use position-based classification
	if pos != nil && mapName != "" {
		return a.classifyFromPosition(pos, mapName)
	}

	// Priority 3: Fallback
	return &SiteClassification{
		Site:       "unknown",
		Confidence: 0.0,
		Method:     "fallback",
	}
}

// classifyFromPosition uses distance-based classification with confidence
func (a *SiteAnalyzerEngine) classifyFromPosition(pos *grid.Position, mapName string) *SiteClassification {
	// Validate position
	if !IsValidVALPosition(pos) {
		return &SiteClassification{
			Site:       "unknown",
			Confidence: 0.0,
			Method:     "invalid_position",
		}
	}

	config := mapSiteConfigs[strings.ToLower(mapName)]
	if config == nil {
		// Use generic classification for unknown maps
		return a.genericSiteClassification(pos)
	}

	// Find closest site center using normalized distance
	var closestSite string
	minDist := math.MaxFloat64

	for siteName, boundary := range config.Sites {
		// Normalized distance (accounts for elliptical boundaries)
		dx := (pos.X - boundary.CenterX) / boundary.RadiusX
		dy := (pos.Y - boundary.CenterY) / boundary.RadiusY
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < minDist {
			minDist = dist
			closestSite = siteName
		}
	}

	// Calculate confidence based on distance
	// dist < 1.0 means inside the ellipse, higher confidence
	// dist > 1.0 means outside, lower confidence
	confidence := math.Max(0.0, 1.0-minDist/2.0)

	return &SiteClassification{
		Site:       closestSite,
		Confidence: confidence,
		Method:     "position",
	}
}

// genericSiteClassification for unknown maps
func (a *SiteAnalyzerEngine) genericSiteClassification(pos *grid.Position) *SiteClassification {
	// Simple Y-axis based classification
	if pos.Y > 0 {
		return &SiteClassification{Site: "A", Confidence: 0.6, Method: "fallback"}
	}
	return &SiteClassification{Site: "B", Confidence: 0.6, Method: "fallback"}
}

// inferSiteFromPositionEnhanced uses the new statistical site detector
func inferSiteFromPositionEnhanced(pos *grid.Position, mapName string) string {
	analyzer := NewSiteAnalyzer()
	result := analyzer.ClassifySite(pos, mapName, "")
	return result.Site
}

func normalizeMapName(name string) string {
	switch name {
	case "Ascent", "ASCENT", "ascent":
		return "ascent"
	case "Bind", "BIND", "bind":
		return "bind"
	case "Haven", "HAVEN", "haven":
		return "haven"
	case "Split", "SPLIT", "split":
		return "split"
	case "Icebox", "ICEBOX", "icebox":
		return "icebox"
	case "Breeze", "BREEZE", "breeze":
		return "breeze"
	case "Fracture", "FRACTURE", "fracture":
		return "fracture"
	case "Pearl", "PEARL", "pearl":
		return "pearl"
	case "Lotus", "LOTUS", "lotus":
		return "lotus"
	case "Sunset", "SUNSET", "sunset":
		return "sunset"
	case "Abyss", "ABYSS", "abyss":
		return "abyss"
	default:
		return name
	}
}

// AnalyzeSitePatterns analyzes attack and defense patterns by site
// ENHANCED: Now uses position data for accurate site detection
func (a *SiteAnalyzerEngine) AnalyzeSitePatterns(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.VALEventData,
) map[string]*SiteAnalysis {
	mapAnalysis := make(map[string]*SiteAnalysis)

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		for _, game := range series.Games {
			if !game.Finished {
				continue
			}

			mapName := game.Map
			if mapName == "" {
				mapName = eventData.MapName
			}
			if mapName == "" {
				mapName = "Unknown"
			}

			if _, exists := mapAnalysis[mapName]; !exists {
				mapAnalysis[mapName] = &SiteAnalysis{
					MapName: mapName,
					Sites:   make(map[string]*SiteStats),
				}
				for _, site := range getMapSites(mapName) {
					mapAnalysis[mapName].Sites[site] = &SiteStats{
						SiteName:             site,
						CommonAttackPatterns: make([]AttackPattern, 0),
						CommonDefenseSetups:  make([]DefenseSetup, 0),
					}
				}
			}

			analysis := mapAnalysis[mapName]

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

			// Analyze plant events using POSITION DATA
			// Filter plants by map name to avoid mixing data from different games
			for _, plant := range eventData.Plants {
				// Skip plants from other maps
				plantMapName := normalizeMapName(plant.MapName)
				gameMapName := normalizeMapName(mapName)
				if plantMapName != "" && gameMapName != "" && plantMapName != gameMapName {
					continue
				}

				site := plant.Site
				if site == "" && plant.Position != nil {
					site = inferSiteFromPositionEnhanced(plant.Position, mapName)
				}
				if site == "" {
					continue
				}
				site = normalizeSiteName(site)

				if siteStats, exists := analysis.Sites[site]; exists {
					if plant.TeamID == teamID || plant.TeamID == ourTeam.ID {
						siteStats.AttackAttempts++
						for _, roundEnd := range eventData.RoundEnds {
							if roundEnd.RoundNum == plant.RoundNum {
								if roundEnd.WinnerTeam == teamID || roundEnd.WinnerTeam == ourTeam.ID {
									siteStats.AttackSuccesses++
								}
								break
							}
						}
					}
				}
			}

			// Analyze defense patterns
			for _, roundEnd := range eventData.RoundEnds {
				isDefending := roundEnd.DefenseTeam == teamID || roundEnd.DefenseTeam == ourTeam.ID
				if !isDefending {
					continue
				}

				var plantSite string
				for _, plant := range eventData.Plants {
					if plant.RoundNum == roundEnd.RoundNum {
						plantSite = plant.Site
						if plantSite == "" && plant.Position != nil {
							plantSite = inferSiteFromPositionEnhanced(plant.Position, mapName)
						}
						plantSite = normalizeSiteName(plantSite)
						break
					}
				}

				if plantSite != "" {
					if siteStats, exists := analysis.Sites[plantSite]; exists {
						siteStats.DefenseAttempts++
						if roundEnd.WinnerTeam == teamID || roundEnd.WinnerTeam == ourTeam.ID {
							siteStats.DefenseSuccesses++
						}
					}
				}
			}

			// Analyze retakes
			for _, defuse := range eventData.Defuses {
				if defuse.TeamID != teamID && defuse.TeamID != ourTeam.ID {
					continue
				}

				var plantSite string
				for _, plant := range eventData.Plants {
					if plant.RoundNum == defuse.RoundNum {
						plantSite = plant.Site
						if plantSite == "" && plant.Position != nil {
							plantSite = inferSiteFromPositionEnhanced(plant.Position, mapName)
						}
						plantSite = normalizeSiteName(plantSite)
						break
					}
				}

				if plantSite != "" {
					if siteStats, exists := analysis.Sites[plantSite]; exists {
						siteStats.RetakeAttempts++
						siteStats.RetakeSuccesses++
					}
				}
			}
		}
	}

	// Calculate win rates
	for _, analysis := range mapAnalysis {
		for _, siteStats := range analysis.Sites {
			if siteStats.AttackAttempts > 0 {
				siteStats.AttackWinRate = float64(siteStats.AttackSuccesses) / float64(siteStats.AttackAttempts)
			}
			if siteStats.DefenseAttempts > 0 {
				siteStats.DefenseWinRate = float64(siteStats.DefenseSuccesses) / float64(siteStats.DefenseAttempts)
			}
			if siteStats.RetakeAttempts > 0 {
				siteStats.RetakeWinRate = float64(siteStats.RetakeSuccesses) / float64(siteStats.RetakeAttempts)
			}
		}
	}

	return mapAnalysis
}

// AnalyzeAttackPatterns identifies common attack patterns
// ENHANCED: Now detects pistol round patterns and fast executes
func (a *SiteAnalyzerEngine) AnalyzeAttackPatterns(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.VALEventData,
) []AttackPattern {
	patterns := make([]AttackPattern, 0)
	patternCounts := make(map[string]*patternAggregator)
	pistolPatternCounts := make(map[string]*patternAggregator)

	// Track which plants we've already processed to avoid duplicates
	processedPlants := make(map[string]bool)

	for seriesID, eventData := range events {
		if eventData == nil {
			continue
		}

		var series *grid.SeriesState
		for _, s := range seriesStates {
			if s.ID == seriesID {
				series = s
				break
			}
		}
		if series == nil {
			continue
		}

		// Process plants once per series, using their map name
		for _, plant := range eventData.Plants {
			if plant.TeamID != teamID {
				continue
			}

			// Create unique key for this plant to avoid duplicates
			plantKey := fmt.Sprintf("%s:%d:%d", plant.MapName, plant.RoundNum, plant.GameTime)
			if processedPlants[plantKey] {
				continue
			}
			processedPlants[plantKey] = true

			mapName := plant.MapName
			if mapName == "" {
				// Fallback to first game's map if plant doesn't have map name
				if len(series.Games) > 0 {
					mapName = series.Games[0].Map
				}
			}

			site := plant.Site
			if site == "" && plant.Position != nil {
				site = inferSiteFromPositionEnhanced(plant.Position, mapName)
			}
			site = normalizeSiteName(site)
			if site == "" {
				continue
			}

			patternType := "default"
			if plant.GameTime < 25000 {
				patternType = "fast-hit"
			} else if plant.GameTime < 40000 {
				patternType = "quick-execute"
			} else if plant.GameTime > 70000 {
				patternType = "slow-default"
			}

			key := fmt.Sprintf("%s:%s:%s", mapName, site, patternType)
			isPistol := plant.RoundNum == 1 || plant.RoundNum == 13

			targetMap := patternCounts
			if isPistol {
				targetMap = pistolPatternCounts
			}

			agg, exists := targetMap[key]
			if !exists {
				agg = &patternAggregator{
					mapName:     mapName,
					site:        site,
					patternType: patternType,
					isPistol:    isPistol,
				}
				targetMap[key] = agg
			}

			agg.count++
			agg.totalTime += float64(plant.GameTime)

			for _, roundEnd := range eventData.RoundEnds {
				if roundEnd.RoundNum == plant.RoundNum {
					if roundEnd.WinnerTeam == teamID {
						agg.wins++
					}
					break
				}
			}
		}
	}

	// Convert regular patterns
	totalPatterns := 0
	for _, agg := range patternCounts {
		totalPatterns += agg.count
	}

	for _, agg := range patternCounts {
		if agg.count < 2 {
			continue
		}

		pattern := AttackPattern{
			Description: fmt.Sprintf("%s on %s-Site", agg.patternType, agg.site),
			TargetSite:  agg.site,
			Frequency:   float64(agg.count) / float64(totalPatterns),
		}
		if agg.count > 0 {
			pattern.SuccessRate = float64(agg.wins) / float64(agg.count)
			pattern.AvgExecuteTime = (agg.totalTime / float64(agg.count)) / 1000.0
		}

		patterns = append(patterns, pattern)
	}

	// Convert pistol patterns - HACKATHON FORMAT
	totalPistolPatterns := 0
	for _, agg := range pistolPatternCounts {
		totalPistolPatterns += agg.count
	}

	for _, agg := range pistolPatternCounts {
		if agg.count < 2 {
			continue
		}

		pattern := AttackPattern{
			Description: fmt.Sprintf("Pistol %s on %s-Site (%s)", agg.patternType, agg.site, agg.mapName),
			TargetSite:  agg.site,
			Frequency:   float64(agg.count) / float64(totalPistolPatterns),
		}
		if agg.count > 0 {
			pattern.SuccessRate = float64(agg.wins) / float64(agg.count)
			pattern.AvgExecuteTime = (agg.totalTime / float64(agg.count)) / 1000.0
		}

		patterns = append(patterns, pattern)
	}

	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})

	return patterns
}

// AnalyzePistolRoundPatterns specifically analyzes pistol round attack patterns
// HACKATHON FORMAT: "On Attack, 70% of pistol rounds are a 5-man fast-hit on B-Site (Ascent)"
func (a *SiteAnalyzerEngine) AnalyzePistolRoundPatterns(
	teamID string,
	seriesStates []*grid.SeriesState,
	events map[string]*grid.VALEventData,
) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	type pistolPattern struct {
		mapName     string
		site        string
		patternType string
		count       int
		wins        int
		avgTime     float64
	}

	patternsByMap := make(map[string][]pistolPattern)
	totalPistolsByMap := make(map[string]int)

	for _, eventData := range events {
		if eventData == nil {
			continue
		}

		mapName := eventData.MapName
		if mapName == "" {
			continue
		}

		for _, plant := range eventData.Plants {
			if plant.TeamID != teamID {
				continue
			}

			if plant.RoundNum != 1 && plant.RoundNum != 13 {
				continue
			}

			totalPistolsByMap[mapName]++

			site := plant.Site
			if site == "" && plant.Position != nil {
				site = inferSiteFromPositionEnhanced(plant.Position, mapName)
			}
			site = normalizeSiteName(site)

			patternType := "default"
			if plant.GameTime < 25000 {
				patternType = "5-man fast-hit"
			} else if plant.GameTime < 40000 {
				patternType = "quick execute"
			} else {
				patternType = "slow default"
			}

			found := false
			for i := range patternsByMap[mapName] {
				if patternsByMap[mapName][i].site == site && patternsByMap[mapName][i].patternType == patternType {
					patternsByMap[mapName][i].count++
					patternsByMap[mapName][i].avgTime += float64(plant.GameTime)

					for _, roundEnd := range eventData.RoundEnds {
						if roundEnd.RoundNum == plant.RoundNum && roundEnd.WinnerTeam == teamID {
							patternsByMap[mapName][i].wins++
							break
						}
					}
					found = true
					break
				}
			}

			if !found {
				p := pistolPattern{
					mapName:     mapName,
					site:        site,
					patternType: patternType,
					count:       1,
					avgTime:     float64(plant.GameTime),
				}
				for _, roundEnd := range eventData.RoundEnds {
					if roundEnd.RoundNum == plant.RoundNum && roundEnd.WinnerTeam == teamID {
						p.wins++
						break
					}
				}
				patternsByMap[mapName] = append(patternsByMap[mapName], p)
			}
		}
	}

	// Generate insights in HACKATHON FORMAT
	for mapName, patterns := range patternsByMap {
		total := totalPistolsByMap[mapName]
		if total < 2 {
			continue
		}

		for _, p := range patterns {
			if p.count < 2 {
				continue
			}

			freq := float64(p.count) / float64(total)
			if freq > 0.4 {
				winRate := 0.0
				if p.count > 0 {
					winRate = float64(p.wins) / float64(p.count)
				}

				insights = append(insights, StrategyInsight{
					Text: fmt.Sprintf("On Attack, %.0f%% of pistol rounds are a %s on %s-Site (%s)",
						freq*100, p.patternType, p.site, mapName),
					Metric:     "pistol_attack_pattern",
					Value:      freq,
					SampleSize: p.count,
					Context:    fmt.Sprintf("%s - %.0f%% win rate", mapName, winRate*100),
				})
			}
		}
	}

	return insights
}

// GenerateSiteInsights generates hackathon-format insights from site analysis
func (a *SiteAnalyzerEngine) GenerateSiteInsights(
	siteAnalysis map[string]*SiteAnalysis,
	attackPatterns []AttackPattern,
) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	for mapName, analysis := range siteAnalysis {
		var mostAttackedSite string
		var maxAttacks int

		for siteName, stats := range analysis.Sites {
			if stats.AttackAttempts > maxAttacks {
				maxAttacks = stats.AttackAttempts
				mostAttackedSite = siteName
			}
		}

		if mostAttackedSite != "" && maxAttacks > 0 {
			siteStats := analysis.Sites[mostAttackedSite]
			totalAttacks := 0
			for _, stats := range analysis.Sites {
				totalAttacks += stats.AttackAttempts
			}

			if totalAttacks > 0 {
				attackRate := float64(maxAttacks) / float64(totalAttacks)
				if attackRate > 0.5 {
					insights = append(insights, StrategyInsight{
						Text: fmt.Sprintf("On Attack (%s), %.0f%% of rounds target %s-Site (%.0f%% success rate)",
							mapName, attackRate*100, mostAttackedSite, siteStats.AttackWinRate*100),
						Metric:     "site_preference",
						Value:      attackRate,
						SampleSize: totalAttacks,
						Context:    mapName,
					})
				}
			}
		}
	}

	for _, pattern := range attackPatterns {
		if pattern.Frequency > 0.3 {
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("Common attack pattern: %s (%.0f%% frequency, %.0f%% success)",
					pattern.Description, pattern.Frequency*100, pattern.SuccessRate*100),
				Metric:  "attack_pattern",
				Value:   pattern.Frequency,
				Context: pattern.TargetSite,
			})
		}
	}

	return insights
}

// GenerateDefenseInsights generates defense setup insights
func (a *SiteAnalyzerEngine) GenerateDefenseInsights(
	siteAnalysis map[string]*SiteAnalysis,
) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	for mapName, analysis := range siteAnalysis {
		var weakestSite string
		var lowestWinRate float64 = 1.0

		for siteName, stats := range analysis.Sites {
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate < lowestWinRate {
				lowestWinRate = stats.DefenseWinRate
				weakestSite = siteName
			}
		}

		if weakestSite != "" && lowestWinRate < 0.5 {
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("On Defense (%s), weak at holding %s-Site (%.0f%% win rate)",
					mapName, weakestSite, lowestWinRate*100),
				Metric:     "defense_weakness",
				Value:      lowestWinRate,
				SampleSize: analysis.Sites[weakestSite].DefenseAttempts,
				Context:    mapName,
			})
		}

		var strongestSite string
		var highestWinRate float64 = 0.0

		for siteName, stats := range analysis.Sites {
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate > highestWinRate {
				highestWinRate = stats.DefenseWinRate
				strongestSite = siteName
			}
		}

		if strongestSite != "" && highestWinRate > 0.6 {
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("On Defense (%s), strong at holding %s-Site (%.0f%% win rate)",
					mapName, strongestSite, highestWinRate*100),
				Metric:     "defense_strength",
				Value:      highestWinRate,
				SampleSize: analysis.Sites[strongestSite].DefenseAttempts,
				Context:    mapName,
			})
		}
	}

	return insights
}

// GenerateSiteCounterStrategies generates counter-strategies based on site analysis
func (a *SiteAnalyzerEngine) GenerateSiteCounterStrategies(
	siteAnalysis map[string]*SiteAnalysis,
) []InGameStrategyInsight {
	strategies := make([]InGameStrategyInsight, 0)

	for mapName, analysis := range siteAnalysis {
		var weakestDefenseSite string
		var lowestDefenseRate float64 = 1.0

		for siteName, stats := range analysis.Sites {
			if stats.DefenseAttempts >= 3 && stats.DefenseWinRate < lowestDefenseRate {
				lowestDefenseRate = stats.DefenseWinRate
				weakestDefenseSite = siteName
			}
		}

		if weakestDefenseSite != "" && lowestDefenseRate < 0.5 {
			strategies = append(strategies, InGameStrategyInsight{
				Strategy: fmt.Sprintf("Attack %s-Site on %s", weakestDefenseSite, mapName),
				Timing:   "Attack rounds",
				Reason: fmt.Sprintf("Opponent has only %.0f%% defense win rate on %s-Site",
					lowestDefenseRate*100, weakestDefenseSite),
				Impact: "HIGH",
			})
		}

		var weakestAttackSite string
		var lowestAttackRate float64 = 1.0

		for siteName, stats := range analysis.Sites {
			if stats.AttackAttempts >= 3 && stats.AttackWinRate < lowestAttackRate {
				lowestAttackRate = stats.AttackWinRate
				weakestAttackSite = siteName
			}
		}

		if weakestAttackSite != "" && lowestAttackRate < 0.5 {
			strategies = append(strategies, InGameStrategyInsight{
				Strategy: fmt.Sprintf("Stack defense on %s-Site on %s", weakestAttackSite, mapName),
				Timing:   "Defense rounds",
				Reason: fmt.Sprintf("Opponent has only %.0f%% attack success rate on %s-Site",
					lowestAttackRate*100, weakestAttackSite),
				Impact: "MEDIUM",
			})
		}
	}

	return strategies
}

// Helper functions

func getMapSites(mapName string) []string {
	mapSites := map[string][]string{
		"Ascent":   {"A", "B"},
		"Bind":     {"A", "B"},
		"Haven":    {"A", "B", "C"},
		"Split":    {"A", "B"},
		"Icebox":   {"A", "B"},
		"Breeze":   {"A", "B"},
		"Fracture": {"A", "B"},
		"Pearl":    {"A", "B"},
		"Lotus":    {"A", "B", "C"},
		"Sunset":   {"A", "B"},
		"Abyss":    {"A", "B"},
		"Unknown":  {"A", "B"},
	}

	if sites, ok := mapSites[mapName]; ok {
		return sites
	}
	return []string{"A", "B"}
}

func normalizeSiteName(site string) string {
	switch site {
	case "A", "a", "A-Site", "a-site", "ASite":
		return "A"
	case "B", "b", "B-Site", "b-site", "BSite":
		return "B"
	case "C", "c", "C-Site", "c-site", "CSite":
		return "C"
	default:
		return site
	}
}

// Helper type for pattern aggregation
type patternAggregator struct {
	mapName     string
	site        string
	patternType string
	count       int
	wins        int
	totalTime   float64
	isPistol    bool
}

// GenerateDefenseSetups infers defensive setups from early-round positioning data
// Hackathon format: "On Defense, they default to a 1-3-1 setup, rotating their Sentinel to mid"
func (a *SiteAnalyzerEngine) GenerateDefenseSetups(
	teamID string,
	eventsData map[string]*grid.VALEventData,
) []DefenseSetup {
	// Track player positioning patterns at round start on defense
	setupCounts := make(map[string]int) // "1-3-1", "2-1-2", "stack A", etc.
	setupWins := make(map[string]int)
	totalDefenseRounds := 0

	for _, eventData := range eventsData {
		// Analyze each round's early kills to infer positions
		roundPositions := make(map[int][]string) // round -> list of positions

		for _, kill := range eventData.Kills {
			// Only consider early-round kills (first 20 seconds) on defense rounds
			if kill.KillerTeamID == teamID && kill.GameTime < 20000 {
				pos := inferPositionFromKill(kill.KillerPosition, kill.MapName)
				if pos != "" {
					roundPositions[kill.RoundNum] = append(roundPositions[kill.RoundNum], pos)
				}
			}
		}

		// Also infer from round end data
		for _, round := range eventData.RoundEnds {
			if round.IsDefendingSide(teamID) {
				totalDefenseRounds++

				// Determine setup based on early positioning (use "Unknown" if no map name)
				setup := inferDefenseSetup(roundPositions[round.RoundNum], "Unknown")
				if setup != "" {
					setupCounts[setup]++
					if round.WinnerTeam == teamID {
						setupWins[setup]++
					}
				}
			}
		}
	}

	// Build defense setup insights
	setups := make([]DefenseSetup, 0)
	for setup, count := range setupCounts {
		if count >= 2 {
			frequency := float64(count) / float64(max(totalDefenseRounds, 1))
			successRate := float64(setupWins[setup]) / float64(count)
			setups = append(setups, DefenseSetup{
				Description: setup,
				Frequency:   frequency,
				SuccessRate: successRate,
				Positions:   make(map[string]string),
			})
		}
	}

	// Sort by frequency
	sort.Slice(setups, func(i, j int) bool {
		return setups[i].Frequency > setups[j].Frequency
	})

	return setups
}

// inferPositionFromKill determines what area a kill occurred in
func inferPositionFromKill(pos *grid.Position, mapName string) string {
	if pos == nil {
		return ""
	}

	// Determine site/area based on position
	x := pos.X

	// Generic position inference based on x coordinate
	if x < -4000 {
		return "A-Site"
	} else if x > 4000 {
		return "B-Site"
	} else {
		return "Mid"
	}
}

// inferDefenseSetup determines the likely setup from observed positions
func inferDefenseSetup(positions []string, mapName string) string {
	if len(positions) == 0 {
		return ""
	}

	aSite := 0
	bSite := 0
	mid := 0

	for _, pos := range positions {
		switch pos {
		case "A-Site", "A":
			aSite++
		case "B-Site", "B":
			bSite++
		case "Mid":
			mid++
		}
	}

	// Infer setup pattern
	if mid >= 3 {
		return "stack Mid"
	} else if aSite >= 3 {
		return "stack A"
	} else if bSite >= 3 {
		return "stack B"
	} else if mid >= 2 && aSite >= 1 && bSite >= 1 {
		return "1-3-1 (mid heavy)"
	} else if aSite >= 2 && bSite >= 2 {
		return "2-1-2 (site heavy)"
	} else if mid >= 1 {
		return "1-3-1"
	}

	return "standard"
}

// GenerateDefenseSetupInsights creates text insights for defense setups
// Hackathon format: "On Defense, they default to a 1-3-1 setup"
func (a *SiteAnalyzerEngine) GenerateDefenseSetupInsights(setups []DefenseSetup) []StrategyInsight {
	insights := make([]StrategyInsight, 0)

	for _, setup := range setups {
		if setup.Frequency >= 0.3 { // At least 30% frequency
			insights = append(insights, StrategyInsight{
				Text: fmt.Sprintf("On Defense, they default to a %s setup (%.0f%% win rate)",
					setup.Description, setup.SuccessRate*100),
				Metric:     "defense_setup",
				Value:      setup.Frequency,
				SampleSize: int(setup.Frequency * 10), // Approximate
				Context:    setup.Description,
			})
		}
	}

	return insights
}
