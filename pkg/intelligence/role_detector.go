package intelligence

import (
	"fmt"

	"scout9/pkg/grid"
)

// DataDrivenRoleDetector infers roles from player statistics
type DataDrivenRoleDetector struct {
	stats            *StatisticalEngine
	championFallback map[string]string // Champion -> primary role mapping
}

// NewDataDrivenRoleDetector creates a new role detector
func NewDataDrivenRoleDetector() *DataDrivenRoleDetector {
	return &DataDrivenRoleDetector{
		stats:            NewStatisticalEngine(),
		championFallback: buildChampionRoleMap(),
	}
}

// RoleDetection represents a role detection result
type RoleDetection struct {
	Role       string   // "top", "jungle", "mid", "adc", "support"
	Confidence float64  // 0.0 to 1.0
	Method     string   // "statistical", "champion_fallback", "position"
	Evidence   []string // Reasons for classification
}

// PlayerStats contains the statistics needed for role detection
type PlayerStats struct {
	PlayerID          string
	PlayerName        string
	Character         string
	Kills             int
	Deaths            int
	Assists           int
	CS                int     // Creep score (minions killed) - NOT available via GRID API
	NetWorth          int     // Gold earned
	JungleCampsKilled int     // Number of jungle camps killed - NOT available via GRID API
	DamageDealt       int     // Total damage dealt - NOT available via GRID API
	DamageShare       float64 // Percentage of team's damage - NOT available via GRID API
	GameDuration      int     // Game duration in seconds
	Positions         []grid.Position // Kill/death positions - available via JSONL events only
}

// DetectRole determines player role from game statistics
func (d *DataDrivenRoleDetector) DetectRole(player *PlayerStats, teamTotalKills int) *RoleDetection {
	evidence := make([]string, 0)

	// Calculate key metrics
	gameMins := float64(player.GameDuration) / 60.0
	if gameMins < 1 {
		gameMins = 1 // Avoid division by zero
	}

	csPerMin := float64(player.CS) / gameMins

	// Assist ratio (assists / (kills + assists))
	assistRatio := 0.0
	if player.Kills+player.Assists > 0 {
		assistRatio = float64(player.Assists) / float64(player.Kills+player.Assists)
	}

	// 1. Jungle detection: High jungle camp kills relative to CS
	if player.JungleCampsKilled > 0 && player.CS > 0 {
		jungleRatio := float64(player.JungleCampsKilled) / float64(player.CS)
		if jungleRatio > 0.3 || player.JungleCampsKilled > 50 {
			evidence = append(evidence, fmt.Sprintf("High jungle camp ratio: %.1f%% (%d camps)",
				jungleRatio*100, player.JungleCampsKilled))
			return &RoleDetection{
				Role:       "jungle",
				Confidence: 0.9,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}
	}

	// 2. Support detection: Low CS, high assists
	if csPerMin < 3.0 && assistRatio > 0.6 {
		evidence = append(evidence, fmt.Sprintf("Low CS (%.1f/min), high assist ratio (%.1f%%)",
			csPerMin, assistRatio*100))
		return &RoleDetection{
			Role:       "support",
			Confidence: 0.85,
			Method:     "statistical",
			Evidence:   evidence,
		}
	}

	// 3. ADC detection: High CS, high damage share, typically bot lane
	if csPerMin > 7.0 && player.DamageShare > 0.25 {
		evidence = append(evidence, fmt.Sprintf("High CS (%.1f/min), high damage share (%.1f%%)",
			csPerMin, player.DamageShare*100))

		// Check position data if available
		if len(player.Positions) > 0 && d.isBotLanePlayer(player.Positions) {
			evidence = append(evidence, "Primarily bot lane positions")
			return &RoleDetection{
				Role:       "adc",
				Confidence: 0.85,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}

		// Without position data, check if champion is typically ADC
		if d.isADCChampion(player.Character) {
			evidence = append(evidence, fmt.Sprintf("%s is typically played ADC", player.Character))
			return &RoleDetection{
				Role:       "adc",
				Confidence: 0.8,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}
	}

	// 4. Mid detection: High CS, mid lane positions
	if csPerMin > 6.0 {
		if len(player.Positions) > 0 && d.isMidLanePlayer(player.Positions) {
			evidence = append(evidence, fmt.Sprintf("High CS (%.1f/min), mid lane positions", csPerMin))
			return &RoleDetection{
				Role:       "mid",
				Confidence: 0.8,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}
	}

	// 5. Top detection: High CS, top lane positions
	if csPerMin > 6.0 {
		if len(player.Positions) > 0 && d.isTopLanePlayer(player.Positions) {
			evidence = append(evidence, fmt.Sprintf("High CS (%.1f/min), top lane positions", csPerMin))
			return &RoleDetection{
				Role:       "top",
				Confidence: 0.8,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}
	}

	// 6. Fallback to champion-based detection
	if role, ok := d.championFallback[player.Character]; ok {
		evidence = append(evidence, fmt.Sprintf("Champion %s typically played %s", player.Character, role))
		return &RoleDetection{
			Role:       role,
			Confidence: 0.6,
			Method:     "champion_fallback",
			Evidence:   evidence,
		}
	}

	// 7. Last resort: use CS and damage patterns
	if csPerMin > 6.0 {
		// High CS but couldn't determine lane - likely solo laner
		if player.DamageShare > 0.2 {
			evidence = append(evidence, fmt.Sprintf("High CS (%.1f/min), moderate damage share", csPerMin))
			return &RoleDetection{
				Role:       "mid", // Default to mid for high CS unknown
				Confidence: 0.5,
				Method:     "statistical",
				Evidence:   evidence,
			}
		}
	}

	return &RoleDetection{
		Role:       "unknown",
		Confidence: 0.0,
		Method:     "none",
		Evidence:   []string{"Insufficient data for role detection"},
	}
}

// isBotLanePlayer checks if majority of positions are in bot lane
func (d *DataDrivenRoleDetector) isBotLanePlayer(positions []grid.Position) bool {
	if len(positions) == 0 {
		return false
	}

	laneDetector := NewStatisticalLaneDetector()
	botCount := 0

	for _, pos := range positions {
		result := laneDetector.ClassifyPosition(&pos)
		if result.Lane == "bot" {
			botCount++
		}
	}

	return float64(botCount)/float64(len(positions)) > 0.5
}

// isMidLanePlayer checks if majority of positions are in mid lane
func (d *DataDrivenRoleDetector) isMidLanePlayer(positions []grid.Position) bool {
	if len(positions) == 0 {
		return false
	}

	laneDetector := NewStatisticalLaneDetector()
	midCount := 0

	for _, pos := range positions {
		result := laneDetector.ClassifyPosition(&pos)
		if result.Lane == "mid" {
			midCount++
		}
	}

	return float64(midCount)/float64(len(positions)) > 0.5
}

// isTopLanePlayer checks if majority of positions are in top lane
func (d *DataDrivenRoleDetector) isTopLanePlayer(positions []grid.Position) bool {
	if len(positions) == 0 {
		return false
	}

	laneDetector := NewStatisticalLaneDetector()
	topCount := 0

	for _, pos := range positions {
		result := laneDetector.ClassifyPosition(&pos)
		if result.Lane == "top" {
			topCount++
		}
	}

	return float64(topCount)/float64(len(positions)) > 0.5
}

// isADCChampion checks if champion is typically played as ADC
func (d *DataDrivenRoleDetector) isADCChampion(champion string) bool {
	adcChampions := map[string]bool{
		"Jinx": true, "Caitlyn": true, "Ezreal": true, "Kai'Sa": true,
		"Vayne": true, "Jhin": true, "Miss Fortune": true, "Ashe": true,
		"Xayah": true, "Aphelios": true, "Samira": true, "Zeri": true,
		"Draven": true, "Lucian": true, "Tristana": true, "Sivir": true,
		"Kog'Maw": true, "Twitch": true, "Varus": true, "Kalista": true,
		"Nilah": true, "Smolder": true,
	}
	return adcChampions[champion]
}

// DetectTeamRoles detects roles for all players on a team
func (d *DataDrivenRoleDetector) DetectTeamRoles(players []*PlayerStats, teamTotalKills int) map[string]*RoleDetection {
	results := make(map[string]*RoleDetection)

	// First pass: detect roles individually
	for _, player := range players {
		results[player.PlayerID] = d.DetectRole(player, teamTotalKills)
	}

	// Second pass: resolve conflicts (ensure one player per role)
	d.resolveRoleConflicts(players, results)

	return results
}

// resolveRoleConflicts ensures each role is assigned to exactly one player
func (d *DataDrivenRoleDetector) resolveRoleConflicts(players []*PlayerStats, results map[string]*RoleDetection) {
	// Count role assignments
	roleCounts := make(map[string][]string) // role -> player IDs
	for playerID, detection := range results {
		if detection.Role != "unknown" {
			roleCounts[detection.Role] = append(roleCounts[detection.Role], playerID)
		}
	}

	// Resolve duplicates by confidence
	for role, playerIDs := range roleCounts {
		if len(playerIDs) > 1 {
			// Keep the one with highest confidence
			var bestPlayer string
			var bestConfidence float64

			for _, pid := range playerIDs {
				if results[pid].Confidence > bestConfidence {
					bestConfidence = results[pid].Confidence
					bestPlayer = pid
				}
			}

			// Mark others as needing reassignment
			for _, pid := range playerIDs {
				if pid != bestPlayer {
					results[pid] = &RoleDetection{
						Role:       "unknown",
						Confidence: 0.0,
						Method:     "conflict_resolution",
						Evidence:   []string{fmt.Sprintf("Role %s assigned to player with higher confidence", role)},
					}
				}
			}
		}
	}
}

// buildChampionRoleMap creates the fallback champion -> role mapping
func buildChampionRoleMap() map[string]string {
	return map[string]string{
		// Top laners
		"Aatrox": "top", "Camille": "top", "Darius": "top", "Fiora": "top",
		"Gangplank": "top", "Garen": "top", "Gnar": "top", "Gragas": "top",
		"Gwen": "top", "Illaoi": "top", "Irelia": "top", "Jax": "top",
		"Jayce": "top", "K'Sante": "top", "Kayle": "top", "Kennen": "top",
		"Kled": "top", "Malphite": "top", "Mordekaiser": "top", "Nasus": "top",
		"Olaf": "top", "Ornn": "top", "Renekton": "top", "Riven": "top",
		"Rumble": "top", "Sett": "top", "Shen": "top", "Singed": "top",
		"Sion": "top", "Teemo": "top", "Tryndamere": "top", "Urgot": "top",
		"Volibear": "top", "Yorick": "top",

		// Junglers
		"Amumu": "jungle", "Bel'Veth": "jungle", "Briar": "jungle", "Diana": "jungle",
		"Ekko": "jungle", "Elise": "jungle", "Evelynn": "jungle", "Fiddlesticks": "jungle",
		"Graves": "jungle", "Hecarim": "jungle", "Ivern": "jungle", "Jarvan IV": "jungle",
		"Karthus": "jungle", "Kayn": "jungle", "Kha'Zix": "jungle", "Kindred": "jungle",
		"Lee Sin": "jungle", "Lillia": "jungle", "Maokai": "jungle", "Master Yi": "jungle",
		"Nidalee": "jungle", "Nocturne": "jungle", "Nunu & Willump": "jungle",
		"Poppy": "jungle", "Rammus": "jungle", "Rek'Sai": "jungle", "Rengar": "jungle",
		"Sejuani": "jungle", "Shaco": "jungle", "Shyvana": "jungle", "Skarner": "jungle",
		"Taliyah": "jungle", "Trundle": "jungle", "Udyr": "jungle", "Vi": "jungle",
		"Viego": "jungle", "Warwick": "jungle", "Wukong": "jungle", "Xin Zhao": "jungle",
		"Zac": "jungle",

		// Mid laners
		"Ahri": "mid", "Akali": "mid", "Anivia": "mid", "Annie": "mid",
		"Aurelion Sol": "mid", "Azir": "mid", "Cassiopeia": "mid", "Corki": "mid",
		"Fizz": "mid", "Galio": "mid", "Kassadin": "mid", "Katarina": "mid",
		"LeBlanc": "mid", "Lissandra": "mid", "Lux": "mid", "Malzahar": "mid",
		"Naafiri": "mid", "Neeko": "mid", "Orianna": "mid", "Qiyana": "mid",
		"Ryze": "mid", "Syndra": "mid", "Talon": "mid", "Twisted Fate": "mid",
		"Veigar": "mid", "Vel'Koz": "mid", "Vex": "mid", "Viktor": "mid",
		"Vladimir": "mid", "Xerath": "mid", "Yasuo": "mid", "Yone": "mid",
		"Zed": "mid", "Ziggs": "mid", "Zoe": "mid",

		// ADCs
		"Aphelios": "adc", "Ashe": "adc", "Caitlyn": "adc", "Draven": "adc",
		"Ezreal": "adc", "Jhin": "adc", "Jinx": "adc", "Kai'Sa": "adc",
		"Kalista": "adc", "Kog'Maw": "adc", "Lucian": "adc", "Miss Fortune": "adc",
		"Nilah": "adc", "Samira": "adc", "Sivir": "adc", "Smolder": "adc",
		"Tristana": "adc", "Twitch": "adc", "Varus": "adc", "Vayne": "adc",
		"Xayah": "adc", "Zeri": "adc",

		// Supports
		"Alistar": "support", "Bard": "support", "Blitzcrank": "support", "Brand": "support",
		"Braum": "support", "Janna": "support", "Karma": "support", "Leona": "support",
		"Lulu": "support", "Milio": "support", "Morgana": "support", "Nami": "support",
		"Nautilus": "support", "Pyke": "support", "Rakan": "support", "Rell": "support",
		"Renata Glasc": "support", "Senna": "support", "Seraphine": "support", "Sona": "support",
		"Soraka": "support", "Tahm Kench": "support", "Taric": "support", "Thresh": "support",
		"Yuumi": "support", "Zilean": "support", "Zyra": "support",
	}
}
