package intelligence

import (
	"fmt"
	"strings"

	"scout9/pkg/grid"
)

// Validator provides data validation for GRID API data
type Validator struct {
	// Known event types for validation
	knownLoLEventTypes map[string]bool
	knownVALEventTypes map[string]bool
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		knownLoLEventTypes: buildKnownLoLEventTypes(),
		knownVALEventTypes: buildKnownVALEventTypes(),
	}
}

// ValidationResult contains the result of a validation check
type ValidationResult struct {
	IsValid  bool
	Warnings []string
	Errors   []string
}

// PositionBounds defines valid coordinate ranges for a game type
type PositionBounds struct {
	MinX, MaxX float64
	MinY, MaxY float64
}

// Game-specific position bounds
var (
	LoLPositionBounds = PositionBounds{
		MinX: 0, MaxX: 14000,
		MinY: 0, MaxY: 14000,
	}
	VALPositionBounds = PositionBounds{
		MinX: -20000, MaxX: 20000,
		MinY: -20000, MaxY: 20000,
	}
)

// ValidatePosition checks if a position is within expected map bounds
func (v *Validator) ValidatePosition(pos *grid.Position, gameType string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if pos == nil {
		result.Warnings = append(result.Warnings, "Position is nil")
		return result
	}

	var bounds PositionBounds
	switch strings.ToLower(gameType) {
	case "lol", "league-of-legends", "3":
		bounds = LoLPositionBounds
	case "val", "valorant", "25":
		bounds = VALPositionBounds
	default:
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown game type: %s, using VAL bounds", gameType))
		bounds = VALPositionBounds
	}

	if pos.X < bounds.MinX || pos.X > bounds.MaxX {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("X coordinate %.2f out of bounds [%.0f, %.0f]", pos.X, bounds.MinX, bounds.MaxX))
	}

	if pos.Y < bounds.MinY || pos.Y > bounds.MaxY {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Y coordinate %.2f out of bounds [%.0f, %.0f]", pos.Y, bounds.MinY, bounds.MaxY))
	}

	return result
}

// ValidateGameTime checks if a game time value is valid
func (v *Validator) ValidateGameTime(gameTimeMs int) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	// Game time should be positive
	if gameTimeMs < 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Negative game time: %d ms", gameTimeMs))
		return result
	}

	// Game time should be less than 2 hours (7200000 ms)
	maxGameTime := 7200000 // 2 hours in milliseconds
	if gameTimeMs > maxGameTime {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Game time %d ms exceeds maximum %d ms (2 hours)", gameTimeMs, maxGameTime))
	}

	// Warning for very long games (over 1 hour)
	if gameTimeMs > 3600000 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Unusually long game time: %d ms (%.1f minutes)", gameTimeMs, float64(gameTimeMs)/60000))
	}

	return result
}

// ValidateEventType checks if an event type is known
func (v *Validator) ValidateEventType(eventType string, gameType string) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if eventType == "" {
		result.Warnings = append(result.Warnings, "Empty event type")
		return result
	}

	var knownTypes map[string]bool
	switch strings.ToLower(gameType) {
	case "lol", "league-of-legends", "3":
		knownTypes = v.knownLoLEventTypes
	case "val", "valorant", "25":
		knownTypes = v.knownVALEventTypes
	default:
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown game type: %s", gameType))
		return result
	}

	if !knownTypes[eventType] {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown event type: %s", eventType))
	}

	return result
}


// ValidateLoLEventData validates a complete LoL event data structure
func (v *Validator) ValidateLoLEventData(data *grid.LoLEventData) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if data == nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "LoL event data is nil")
		return result
	}

	// Validate kill events
	for i, kill := range data.Kills {
		// Validate game time
		timeResult := v.ValidateGameTime(kill.GameTime)
		if !timeResult.IsValid {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Kill %d: %s", i, strings.Join(timeResult.Errors, ", ")))
		}

		// Validate positions
		if kill.KillerPosition != nil {
			posResult := v.ValidatePosition(kill.KillerPosition, "lol")
			if !posResult.IsValid {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Kill %d killer position: %s", i, strings.Join(posResult.Errors, ", ")))
			}
		}
		if kill.VictimPosition != nil {
			posResult := v.ValidatePosition(kill.VictimPosition, "lol")
			if !posResult.IsValid {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Kill %d victim position: %s", i, strings.Join(posResult.Errors, ", ")))
			}
		}
	}

	// Validate dragon kills
	for i, dragon := range data.DragonKills {
		timeResult := v.ValidateGameTime(dragon.GameTime)
		if !timeResult.IsValid {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Dragon %d: %s", i, strings.Join(timeResult.Errors, ", ")))
		}
	}

	// Validate tower destroys
	for i, tower := range data.TowerDestroys {
		timeResult := v.ValidateGameTime(tower.GameTime)
		if !timeResult.IsValid {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Tower %d: %s", i, strings.Join(timeResult.Errors, ", ")))
		}
	}

	return result
}

// ValidateVALEventData validates a complete VALORANT event data structure
func (v *Validator) ValidateVALEventData(data *grid.VALEventData) *ValidationResult {
	result := &ValidationResult{IsValid: true}

	if data == nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "VAL event data is nil")
		return result
	}

	// Validate plant events
	for i, plant := range data.Plants {
		// Validate game time (round time, typically < 2 minutes)
		if plant.GameTime < 0 {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Plant %d: negative game time %d", i, plant.GameTime))
		}
		if plant.GameTime > 180000 { // 3 minutes max for a round
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Plant %d: unusually long round time %d ms", i, plant.GameTime))
		}

		// Validate position
		if plant.Position != nil {
			posResult := v.ValidatePosition(plant.Position, "val")
			if !posResult.IsValid {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Plant %d position: %s", i, strings.Join(posResult.Errors, ", ")))
			}
		}

		// Validate round number
		if plant.RoundNum < 1 || plant.RoundNum > 30 {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Plant %d: unusual round number %d", i, plant.RoundNum))
		}
	}

	// Validate kill events
	for i, kill := range data.Kills {
		if kill.KillerPosition != nil {
			posResult := v.ValidatePosition(kill.KillerPosition, "val")
			if !posResult.IsValid {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Kill %d killer position: %s", i, strings.Join(posResult.Errors, ", ")))
			}
		}
		if kill.VictimPosition != nil {
			posResult := v.ValidatePosition(kill.VictimPosition, "val")
			if !posResult.IsValid {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Kill %d victim position: %s", i, strings.Join(posResult.Errors, ", ")))
			}
		}
	}

	return result
}

// DataQualityReport contains statistics about parsed data quality
type DataQualityReport struct {
	TotalEvents       int
	ValidEvents       int
	EventsWithWarnings int
	EventsWithErrors  int
	PositionCoverage  float64 // Percentage of events with valid position data
	TimingCoverage    float64 // Percentage of events with valid timing data
	Warnings          []string
	Errors            []string
}

// GenerateLoLDataQualityReport generates a quality report for LoL event data
func (v *Validator) GenerateLoLDataQualityReport(data *grid.LoLEventData) *DataQualityReport {
	report := &DataQualityReport{}

	if data == nil {
		report.Errors = append(report.Errors, "No data provided")
		return report
	}

	// Count events
	report.TotalEvents = len(data.Kills) + len(data.DragonKills) + len(data.BaronKills) +
		len(data.HeraldKills) + len(data.TowerDestroys)

	// Track position and timing coverage
	eventsWithPosition := 0
	eventsWithTiming := 0

	// Validate kills
	for _, kill := range data.Kills {
		if kill.KillerPosition != nil || kill.VictimPosition != nil {
			eventsWithPosition++
		}
		if kill.GameTime > 0 {
			eventsWithTiming++
		}
	}

	// Validate objectives
	for _, dragon := range data.DragonKills {
		if dragon.Position != nil {
			eventsWithPosition++
		}
		if dragon.GameTime > 0 {
			eventsWithTiming++
		}
	}

	for _, baron := range data.BaronKills {
		if baron.Position != nil {
			eventsWithPosition++
		}
		if baron.GameTime > 0 {
			eventsWithTiming++
		}
	}

	for _, herald := range data.HeraldKills {
		if herald.Position != nil {
			eventsWithPosition++
		}
		if herald.GameTime > 0 {
			eventsWithTiming++
		}
	}

	for _, tower := range data.TowerDestroys {
		if tower.GameTime > 0 {
			eventsWithTiming++
		}
	}

	// Calculate coverage
	if report.TotalEvents > 0 {
		report.PositionCoverage = float64(eventsWithPosition) / float64(report.TotalEvents) * 100
		report.TimingCoverage = float64(eventsWithTiming) / float64(report.TotalEvents) * 100
	}

	report.ValidEvents = report.TotalEvents - report.EventsWithErrors

	return report
}

// buildKnownLoLEventTypes returns known LoL event types
func buildKnownLoLEventTypes() map[string]bool {
	return map[string]bool{
		"player-killed-player":     true,
		"player-killed-ATierNPC":   true,
		"player-killed-BTierNPC":   true,
		"team-destroyed-tower":     true,
		"team-destroyed-inhibitor": true,
		"player-picked-character":  true,
		"player-banned-character":  true,
		"series-started-game":      true,
		"game-ended":               true,
		"player-used-ability":      true,
		"team-completed-destroyTurretPlateTop": true,
		"team-completed-destroyTurretPlateMid": true,
		"team-completed-destroyTurretPlateBot": true,
		"game-started-npcRespawnClock":         true,
		"game-respawned-BTierNpc":              true,
	}
}

// buildKnownVALEventTypes returns known VALORANT event types
func buildKnownVALEventTypes() map[string]bool {
	return map[string]bool{
		"player-killed-player":      true,
		"player-completed-plantBomb": true,
		"player-completed-defuseBomb": true,
		"team-won-round":            true,
		"game-started-round":        true,
		"series-started-game":       true,
		"game-ended":                true,
		"player-used-ability":       true,
	}
}

// IsValidLoLPosition checks if a position is valid for LoL
func IsValidLoLPosition(pos *grid.Position) bool {
	if pos == nil {
		return false
	}
	return pos.X >= 0 && pos.X <= 14000 && pos.Y >= 0 && pos.Y <= 14000
}

// IsValidVALPosition checks if a position is valid for VALORANT
func IsValidVALPosition(pos *grid.Position) bool {
	if pos == nil {
		return false
	}
	return pos.X >= -20000 && pos.X <= 20000 && pos.Y >= -20000 && pos.Y <= 20000
}

// IsValidGameTime checks if a game time is valid
func IsValidGameTime(gameTimeMs int) bool {
	return gameTimeMs >= 0 && gameTimeMs <= 7200000
}
