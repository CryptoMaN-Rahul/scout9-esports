package intelligence

import (
	"math"

	"scout9/pkg/grid"
)

// StatisticalLaneDetector uses data-driven lane classification for LoL
type StatisticalLaneDetector struct {
	stats *StatisticalEngine

	// Learned thresholds from data (can be updated with more samples)
	midLaneThreshold float64 // Distance from diagonal to be considered mid
	riverRadius      float64 // Radius around river objectives
}

// NewStatisticalLaneDetector creates a new lane detector with default thresholds
func NewStatisticalLaneDetector() *StatisticalLaneDetector {
	return &StatisticalLaneDetector{
		stats:            NewStatisticalEngine(),
		midLaneThreshold: 2500.0, // Default, can be learned from data
		riverRadius:      2000.0,
	}
}

// LaneClassification represents a lane detection result with confidence
type LaneClassification struct {
	Lane       string  // "top", "mid", "bot", "jungle", "river", "base", "unknown"
	SubRegion  string  // More specific: "top_jungle", "bot_jungle", "dragon_pit", "baron_pit"
	Confidence float64 // 0.0 to 1.0
	Method     string  // "statistical", "fallback"
}

// Summoner's Rift constants
const (
	// Map dimensions (approximately 14000x14000 units)
	lolMapSize   = 14000.0
	lolMapCenter = 7000.0

	// Base positions
	blueBaseX = 1500.0
	blueBaseY = 1500.0
	redBaseX  = 12500.0
	redBaseY  = 12500.0

	// River objective positions
	dragonPitX = 9866.0
	dragonPitY = 4414.0
	baronPitX  = 4414.0
	baronPitY  = 9866.0

	// Base radius for spawn area
	baseRadius = 1500.0
)

// ClassifyPosition determines lane from position using statistical analysis
func (d *StatisticalLaneDetector) ClassifyPosition(pos *grid.Position) *LaneClassification {
	if pos == nil {
		return &LaneClassification{Lane: "unknown", Confidence: 0.0, Method: "none"}
	}

	// Validate position is within map bounds
	if !d.isValidPosition(pos) {
		return &LaneClassification{Lane: "unknown", Confidence: 0.0, Method: "invalid"}
	}

	// Check if in base area first
	if d.isInBase(pos) {
		return &LaneClassification{
			Lane:       "base",
			SubRegion:  d.getBaseSide(pos),
			Confidence: 0.95,
			Method:     "statistical",
		}
	}

	// Check river objectives (dragon pit, baron pit)
	if riverResult := d.checkRiverObjectives(pos); riverResult != nil {
		return riverResult
	}

	// Calculate diagonal distance (key insight: mid lane is where y ≈ x)
	// Positive = above diagonal (top side), Negative = below diagonal (bot side)
	diagonalDist := pos.Y - pos.X

	// Classification based on diagonal distance
	absDiag := math.Abs(diagonalDist)

	// Mid lane: near the diagonal
	if absDiag < d.midLaneThreshold {
		// Check if it's river (perpendicular to mid lane, near center)
		if d.isRiverPosition(pos) {
			return &LaneClassification{
				Lane:       "river",
				SubRegion:  d.getRiverSubRegion(pos),
				Confidence: d.calculateConfidence(absDiag, "river"),
				Method:     "statistical",
			}
		}
		return &LaneClassification{
			Lane:       "mid",
			Confidence: d.calculateConfidence(absDiag, "mid"),
			Method:     "statistical",
		}
	}

	// Top side (above diagonal)
	if diagonalDist > d.midLaneThreshold {
		if d.isLanePosition(pos, "top") {
			return &LaneClassification{
				Lane:       "top",
				Confidence: d.calculateConfidence(diagonalDist, "top"),
				Method:     "statistical",
			}
		}
		return &LaneClassification{
			Lane:       "jungle",
			SubRegion:  d.getJungleSubRegion(pos, "top"),
			Confidence: d.calculateConfidence(diagonalDist, "top_jungle"),
			Method:     "statistical",
		}
	}

	// Bot side (below diagonal)
	if d.isLanePosition(pos, "bot") {
		return &LaneClassification{
			Lane:       "bot",
			Confidence: d.calculateConfidence(-diagonalDist, "bot"),
			Method:     "statistical",
		}
	}
	return &LaneClassification{
		Lane:       "jungle",
		SubRegion:  d.getJungleSubRegion(pos, "bot"),
		Confidence: d.calculateConfidence(-diagonalDist, "bot_jungle"),
		Method:     "statistical",
	}
}


// isValidPosition checks if position is within Summoner's Rift bounds
func (d *StatisticalLaneDetector) isValidPosition(pos *grid.Position) bool {
	return pos.X >= 0 && pos.X <= lolMapSize && pos.Y >= 0 && pos.Y <= lolMapSize
}

// isInBase checks if position is in either team's base
func (d *StatisticalLaneDetector) isInBase(pos *grid.Position) bool {
	// Blue base (bottom-left)
	blueDist := math.Sqrt(math.Pow(pos.X-blueBaseX, 2) + math.Pow(pos.Y-blueBaseY, 2))
	if blueDist < baseRadius {
		return true
	}

	// Red base (top-right)
	redDist := math.Sqrt(math.Pow(pos.X-redBaseX, 2) + math.Pow(pos.Y-redBaseY, 2))
	return redDist < baseRadius
}

// getBaseSide returns which team's base the position is in
func (d *StatisticalLaneDetector) getBaseSide(pos *grid.Position) string {
	blueDist := math.Sqrt(math.Pow(pos.X-blueBaseX, 2) + math.Pow(pos.Y-blueBaseY, 2))
	redDist := math.Sqrt(math.Pow(pos.X-redBaseX, 2) + math.Pow(pos.Y-redBaseY, 2))

	if blueDist < redDist {
		return "blue_base"
	}
	return "red_base"
}

// checkRiverObjectives checks if position is near dragon or baron pit
func (d *StatisticalLaneDetector) checkRiverObjectives(pos *grid.Position) *LaneClassification {
	// Dragon pit check
	dragonDist := math.Sqrt(math.Pow(pos.X-dragonPitX, 2) + math.Pow(pos.Y-dragonPitY, 2))
	if dragonDist < d.riverRadius {
		confidence := 1.0 - (dragonDist / d.riverRadius) // Higher confidence closer to center
		return &LaneClassification{
			Lane:       "river",
			SubRegion:  "dragon_pit",
			Confidence: math.Max(0.7, confidence),
			Method:     "statistical",
		}
	}

	// Baron pit check
	baronDist := math.Sqrt(math.Pow(pos.X-baronPitX, 2) + math.Pow(pos.Y-baronPitY, 2))
	if baronDist < d.riverRadius {
		confidence := 1.0 - (baronDist / d.riverRadius)
		return &LaneClassification{
			Lane:       "river",
			SubRegion:  "baron_pit",
			Confidence: math.Max(0.7, confidence),
			Method:     "statistical",
		}
	}

	return nil
}

// isRiverPosition checks if position is in the river (perpendicular to mid lane)
func (d *StatisticalLaneDetector) isRiverPosition(pos *grid.Position) bool {
	// River runs perpendicular to mid lane, roughly from dragon to baron
	// Check if position is close to the river line

	// River center line goes from approximately (9866, 4414) to (4414, 9866)
	// This is perpendicular to the mid lane diagonal

	// Calculate distance to river center line
	// River line: y = -x + 14280 (approximately)
	// Distance from point to line ax + by + c = 0: |ax + by + c| / sqrt(a² + b²)
	// For y = -x + 14280: x + y - 14280 = 0
	riverDist := math.Abs(pos.X+pos.Y-14280) / math.Sqrt(2)

	return riverDist < 1500 // Within 1500 units of river center
}

// getRiverSubRegion determines which part of the river
func (d *StatisticalLaneDetector) getRiverSubRegion(pos *grid.Position) string {
	// Check proximity to objectives
	dragonDist := math.Sqrt(math.Pow(pos.X-dragonPitX, 2) + math.Pow(pos.Y-dragonPitY, 2))
	baronDist := math.Sqrt(math.Pow(pos.X-baronPitX, 2) + math.Pow(pos.Y-baronPitY, 2))

	if dragonDist < 2500 {
		return "dragon_side"
	}
	if baronDist < 2500 {
		return "baron_side"
	}
	return "mid_river"
}

// isLanePosition checks if position is in a lane (vs jungle)
func (d *StatisticalLaneDetector) isLanePosition(pos *grid.Position, lane string) bool {
	// Lanes are along the edges of the map
	// Top lane: high Y, low-to-mid X (upper-left edge)
	// Bot lane: low Y, mid-to-high X (lower-right edge)

	switch lane {
	case "top":
		// Top lane runs along the top-left edge
		// High Y values (>10000) or low X values (<4000) with high Y
		if pos.Y > 10000 {
			return true
		}
		if pos.X < 4000 && pos.Y > 8000 {
			return true
		}
		return false

	case "bot":
		// Bot lane runs along the bottom-right edge
		// Low Y values (<4000) or high X values (>10000) with low Y
		if pos.Y < 4000 {
			return true
		}
		if pos.X > 10000 && pos.Y < 6000 {
			return true
		}
		return false
	}

	return false
}

// getJungleSubRegion determines which quadrant of jungle
func (d *StatisticalLaneDetector) getJungleSubRegion(pos *grid.Position, side string) string {
	// Divide jungle into quadrants based on position relative to center
	if pos.X < lolMapCenter {
		if pos.Y > lolMapCenter {
			return "blue_top_jungle" // Top-left quadrant
		}
		return "blue_bot_jungle" // Bottom-left quadrant
	}
	if pos.Y > lolMapCenter {
		return "red_top_jungle" // Top-right quadrant
	}
	return "red_bot_jungle" // Bottom-right quadrant
}

// calculateConfidence calculates confidence based on distance from threshold
func (d *StatisticalLaneDetector) calculateConfidence(distance float64, region string) float64 {
	// Base confidence starts at 0.6
	baseConfidence := 0.6

	// Increase confidence based on how far from threshold
	switch region {
	case "mid":
		// Closer to diagonal = higher confidence
		if distance < d.midLaneThreshold*0.5 {
			return 0.9
		}
		return baseConfidence + 0.2*(1-distance/d.midLaneThreshold)

	case "top", "bot":
		// Further from diagonal = higher confidence
		if distance > d.midLaneThreshold*2 {
			return 0.9
		}
		return baseConfidence + 0.2*(distance/d.midLaneThreshold-1)

	case "river":
		return 0.8

	default:
		return baseConfidence
	}
}

// UpdateThresholds allows updating thresholds based on learned data
func (d *StatisticalLaneDetector) UpdateThresholds(midThreshold, riverRadius float64) {
	if midThreshold > 0 {
		d.midLaneThreshold = midThreshold
	}
	if riverRadius > 0 {
		d.riverRadius = riverRadius
	}
}

// LearnThresholdsFromData calculates optimal thresholds from position data
func (d *StatisticalLaneDetector) LearnThresholdsFromData(positions []grid.Position, lanes []string) {
	if len(positions) != len(lanes) || len(positions) < 10 {
		return // Not enough data
	}

	// Collect diagonal distances for each lane
	midDistances := make([]float64, 0)
	topDistances := make([]float64, 0)
	botDistances := make([]float64, 0)

	for i, pos := range positions {
		diag := pos.Y - pos.X
		switch lanes[i] {
		case "mid":
			midDistances = append(midDistances, math.Abs(diag))
		case "top":
			topDistances = append(topDistances, diag)
		case "bot":
			botDistances = append(botDistances, -diag)
		}
	}

	// Calculate optimal threshold as the point that best separates mid from side lanes
	if len(midDistances) > 0 && (len(topDistances) > 0 || len(botDistances) > 0) {
		midDist := d.stats.CalculateDistribution(midDistances)

		// Use 90th percentile of mid lane distances as threshold
		if midDist.SampleSize >= 5 {
			d.midLaneThreshold = midDist.Percentiles[90]
		}
	}
}

// GetJungleSide determines which side of the jungle a position is on relative to a team
func (d *StatisticalLaneDetector) GetJungleSide(pos *grid.Position, teamSide string) string {
	if pos == nil {
		return "unknown"
	}

	// Blue side jungle is bottom-left, Red side jungle is top-right
	if teamSide == "blue" {
		if pos.X < lolMapCenter && pos.Y < lolMapCenter {
			return "own_bot" // Own bot-side jungle
		} else if pos.X < lolMapCenter && pos.Y > lolMapCenter {
			return "own_top" // Own top-side jungle
		} else if pos.X > lolMapCenter && pos.Y < lolMapCenter {
			return "enemy_bot" // Enemy bot-side jungle
		}
		return "enemy_top" // Enemy top-side jungle
	}

	// Red team's jungle (inverted)
	if pos.X > lolMapCenter && pos.Y > lolMapCenter {
		return "own_top" // Own top-side jungle
	} else if pos.X > lolMapCenter && pos.Y < lolMapCenter {
		return "own_bot" // Own bot-side jungle
	} else if pos.X < lolMapCenter && pos.Y > lolMapCenter {
		return "enemy_top" // Enemy top-side jungle
	}
	return "enemy_bot" // Enemy bot-side jungle
}
