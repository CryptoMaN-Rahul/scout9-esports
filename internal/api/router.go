package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"scout9/pkg/cache"
	"scout9/pkg/grid"
	"scout9/pkg/intelligence"
	"scout9/pkg/llm"
	"scout9/pkg/report"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// JSON file for persisting reports (backup storage)
const reportsBackupFile = "data/reports_backup.json"

// splitAndTrim splits a string by separator and trims whitespace from each part
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Server holds the API dependencies
type Server struct {
	gridClient         *grid.Client
	llmService         llm.Service
	cache              *cache.RedisCache
	reportGenerator    *report.Generator
	headToHeadAnalyzer *intelligence.HeadToHeadAnalyzer // Task 10: Added head-to-head analyzer
	lolAnalyzer        *intelligence.LoLAnalyzer        // For team analysis in matchups
	valAnalyzer        *intelligence.VALAnalyzer        // For team analysis in matchups
	// In-memory report storage with JSON file backup
	reports     map[string]*intelligence.ScoutingReport
	reportsLock sync.RWMutex
}

// saveReportsToFile persists reports to JSON file as backup
func (s *Server) saveReportsToFile() {
	s.reportsLock.RLock()
	defer s.reportsLock.RUnlock()

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Warning: Failed to create data directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(s.reports, "", "  ")
	if err != nil {
		log.Printf("Warning: Failed to marshal reports: %v", err)
		return
	}

	if err := os.WriteFile(reportsBackupFile, data, 0644); err != nil {
		log.Printf("Warning: Failed to save reports backup: %v", err)
		return
	}

	log.Printf("Saved %d reports to backup file", len(s.reports))
}

// loadReportsFromFile loads reports from JSON backup file on startup
func (s *Server) loadReportsFromFile() {
	data, err := os.ReadFile(reportsBackupFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No reports backup file found, starting fresh")
			return
		}
		log.Printf("Warning: Failed to read reports backup: %v", err)
		return
	}

	s.reportsLock.Lock()
	defer s.reportsLock.Unlock()

	if err := json.Unmarshal(data, &s.reports); err != nil {
		log.Printf("Warning: Failed to parse reports backup: %v", err)
		return
	}

	log.Printf("Loaded %d reports from backup file", len(s.reports))
}

// NewRouter creates a new API router
func NewRouter(gridClient *grid.Client, llmService llm.Service, cacheClient *cache.RedisCache) http.Handler {
	s := &Server{
		gridClient:         gridClient,
		llmService:         llmService,
		cache:              cacheClient,
		reportGenerator:    report.NewGenerator(gridClient),
		headToHeadAnalyzer: intelligence.NewHeadToHeadAnalyzer(gridClient), // Task 10: Initialize head-to-head analyzer
		lolAnalyzer:        intelligence.NewLoLAnalyzer(),                  // For matchup team analysis
		valAnalyzer:        intelligence.NewVALAnalyzer(),                  // For matchup team analysis
		reports:            make(map[string]*intelligence.ScoutingReport),
	}

	// Load any previously saved reports from backup file
	s.loadReportsFromFile()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(300 * time.Second)) // 5 minutes for report generation

	// CORS - Dynamic origins for local dev and production deployments
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
	}

	// Add production origins from environment variable
	// Set CORS_ORIGINS="https://your-app.vercel.app,https://your-app.pages.dev"
	if envOrigins := os.Getenv("CORS_ORIGINS"); envOrigins != "" {
		for _, origin := range splitAndTrim(envOrigins, ",") {
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	// Also support FRONTEND_URL for single origin
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours preflight cache
	}))

	// Health check
	r.Get("/health", s.healthCheck)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Discovery endpoints
		r.Get("/titles", s.getTitles)
		r.Get("/tournaments", s.getTournaments)
		r.Get("/teams", s.getTeams)
		r.Get("/teams/search", s.searchTeams)
		r.Get("/teams/{teamId}", s.getTeamByID)
		r.Get("/teams/{teamId}/series", s.getSeriesForTeam)

		// Report endpoints
		r.Post("/reports/generate", s.generateReport)
		r.Get("/reports/{reportId}", s.getReport)
		r.Get("/reports", s.listReports)

		// Analysis endpoints
		r.Get("/matchup", s.getMatchup)
		r.Get("/series/{seriesId}/state", s.getSeriesState)
	})

	return r
}

// Response helpers
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "scout9",
	})
}

// getTitles returns available game titles
func (s *Server) getTitles(w http.ResponseWriter, r *http.Request) {
	titles, err := s.gridClient.GetTitles(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch titles: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, titles)
}

// getTournaments returns tournaments for a title
func (s *Server) getTournaments(w http.ResponseWriter, r *http.Request) {
	titleID := r.URL.Query().Get("title")
	if titleID == "" {
		respondError(w, http.StatusBadRequest, "title parameter is required")
		return
	}

	tournaments, err := s.gridClient.GetTournaments(r.Context(), titleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tournaments: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tournaments)
}

// getTeams returns teams in a tournament
func (s *Server) getTeams(w http.ResponseWriter, r *http.Request) {
	tournamentID := r.URL.Query().Get("tournament")
	if tournamentID == "" {
		respondError(w, http.StatusBadRequest, "tournament parameter is required")
		return
	}

	teams, err := s.gridClient.GetTeams(r.Context(), tournamentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch teams: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, teams)
}

// searchTeams searches for teams by name
func (s *Server) searchTeams(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	titleID := r.URL.Query().Get("title")

	teams, err := s.gridClient.SearchTeams(r.Context(), query, titleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to search teams: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, teams)
}

// getTeamByID returns a single team
func (s *Server) getTeamByID(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "teamId")

	team, err := s.gridClient.GetTeamByID(r.Context(), teamID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch team: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, team)
}

// getSeriesForTeam returns recent series for a team
func (s *Server) getSeriesForTeam(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "teamId")

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	series, err := s.gridClient.GetSeriesForTeam(r.Context(), teamID, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch series: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, series)
}

// getSeriesState returns detailed state for a series
func (s *Server) getSeriesState(w http.ResponseWriter, r *http.Request) {
	seriesID := chi.URLParam(r, "seriesId")

	state, err := s.gridClient.GetSeriesState(r.Context(), seriesID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch series state: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, state)
}

// GenerateReportRequest is the request body for report generation
type GenerateReportRequest struct {
	TeamID     string `json:"teamId"`
	TeamName   string `json:"teamName,omitempty"`
	MatchCount int    `json:"matchCount"`
	TitleID    string `json:"titleId"` // "3" for LoL, "6" for VALORANT
}

// generateReport generates a scouting report
func (s *Server) generateReport(w http.ResponseWriter, r *http.Request) {
	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TeamID == "" {
		respondError(w, http.StatusBadRequest, "teamId is required")
		return
	}

	if req.MatchCount <= 0 {
		req.MatchCount = 10
	}

	if req.TitleID == "" {
		req.TitleID = "3" // Default to LoL
	}

	// Generate the report
	genReq := report.GenerateRequest{
		TeamID:     req.TeamID,
		TeamName:   req.TeamName,
		TitleID:    req.TitleID,
		MatchCount: req.MatchCount,
	}

	scoutReport, err := s.reportGenerator.GenerateReport(r.Context(), genReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate report: "+err.Error())
		return
	}

	// Store the report
	s.reportsLock.Lock()
	s.reports[scoutReport.ID] = scoutReport
	s.reportsLock.Unlock()

	// Save to backup file (async to not block response)
	go s.saveReportsToFile()

	respondJSON(w, http.StatusOK, scoutReport)
}

// getReport returns a generated report
func (s *Server) getReport(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "reportId")

	s.reportsLock.RLock()
	report, exists := s.reports[reportID]
	s.reportsLock.RUnlock()

	if !exists {
		respondError(w, http.StatusNotFound, "Report not found")
		return
	}

	respondJSON(w, http.StatusOK, report)
}

// listReports returns all saved reports
func (s *Server) listReports(w http.ResponseWriter, r *http.Request) {
	s.reportsLock.RLock()
	defer s.reportsLock.RUnlock()

	// Return summary list
	summaries := make([]map[string]interface{}, 0, len(s.reports))
	for _, report := range s.reports {
		summaries = append(summaries, map[string]interface{}{
			"id":              report.ID,
			"teamName":        report.OpponentTeam.Name,
			"title":           report.Title,
			"generatedAt":     report.GeneratedAt,
			"matchesAnalyzed": report.MatchesAnalyzed,
		})
	}

	respondJSON(w, http.StatusOK, summaries)
}

// getMatchup returns head-to-head analysis using HeadToHeadAnalyzer (Task 10)
func (s *Server) getMatchup(w http.ResponseWriter, r *http.Request) {
	team1 := r.URL.Query().Get("team1")
	team2 := r.URL.Query().Get("team2")
	titleID := r.URL.Query().Get("title")

	// Task 10.2: Handle invalid team IDs (400)
	if team1 == "" || team2 == "" {
		respondError(w, http.StatusBadRequest, "team1 and team2 parameters are required")
		return
	}

	if titleID == "" {
		titleID = "3" // Default to LoL
	}

	// Determine game title string
	title := "lol"
	if titleID == "6" || titleID == "valorant" {
		title = "valorant"
	}

	matchCount := 20
	if mc := r.URL.Query().Get("matches"); mc != "" {
		if parsed, err := strconv.Atoi(mc); err == nil && parsed > 0 {
			matchCount = parsed
		}
	}

	// Task 10.1: Use HeadToHeadAnalyzer for matchup analysis
	report, err := s.headToHeadAnalyzer.AnalyzeMatchup(r.Context(), team1, team2, titleID, matchCount)
	if err != nil {
		// Task 10.2: Handle team not found (404) vs GRID API errors (502)
		if err.Error() == "team not found" {
			respondError(w, http.StatusNotFound, "Team not found: "+team1+" or "+team2)
			return
		}
		respondError(w, http.StatusBadGateway, "Failed to fetch data from GRID API: "+err.Error())
		return
	}

	// ENHANCED: Generate team analysis and style comparison
	// Fetch series data for both teams to analyze their styles
	team1Series, err1 := s.gridClient.GetSeriesForTeam(r.Context(), team1, matchCount)
	team2Series, err2 := s.gridClient.GetSeriesForTeam(r.Context(), team2, matchCount)

	if err1 == nil && err2 == nil && len(team1Series) > 0 && len(team2Series) > 0 {
		// Get series states for team analysis
		team1SeriesIDs := make([]string, len(team1Series))
		for i, s := range team1Series {
			team1SeriesIDs[i] = s.ID
		}
		team2SeriesIDs := make([]string, len(team2Series))
		for i, s := range team2Series {
			team2SeriesIDs[i] = s.ID
		}

		team1States, err1 := s.gridClient.GetSeriesStates(r.Context(), team1SeriesIDs)
		team2States, err2 := s.gridClient.GetSeriesStates(r.Context(), team2SeriesIDs)

		if err1 == nil && err2 == nil {
			// Run analyzers to get team metrics
			var team1Analysis, team2Analysis *intelligence.TeamAnalysis

			if title == "lol" {
				// Download LoL events for deeper analysis
				lolEvents1 := make(map[string]*grid.LoLEventData)
				lolEvents2 := make(map[string]*grid.LoLEventData)
				for _, seriesID := range team1SeriesIDs[:min(5, len(team1SeriesIDs))] {
					if events, err := s.gridClient.DownloadEvents(r.Context(), seriesID); err == nil {
						if parsed, err := s.gridClient.ParseLoLEvents(events); err == nil {
							lolEvents1[seriesID] = parsed
						}
					}
				}
				for _, seriesID := range team2SeriesIDs[:min(5, len(team2SeriesIDs))] {
					if events, err := s.gridClient.DownloadEvents(r.Context(), seriesID); err == nil {
						if parsed, err := s.gridClient.ParseLoLEvents(events); err == nil {
							lolEvents2[seriesID] = parsed
						}
					}
				}

				team1Analysis, _ = s.lolAnalyzer.AnalyzeTeam(r.Context(), team1, report.Team1Name, team1States, lolEvents1)
				team2Analysis, _ = s.lolAnalyzer.AnalyzeTeam(r.Context(), team2, report.Team2Name, team2States, lolEvents2)
			} else {
				// Download VALORANT events for deeper analysis
				valEvents1 := make(map[string]*grid.VALEventData)
				valEvents2 := make(map[string]*grid.VALEventData)
				for _, seriesID := range team1SeriesIDs[:min(5, len(team1SeriesIDs))] {
					if events, err := s.gridClient.DownloadEvents(r.Context(), seriesID); err == nil {
						if parsed, err := s.gridClient.ParseVALEvents(events); err == nil {
							valEvents1[seriesID] = parsed
						}
					}
				}
				for _, seriesID := range team2SeriesIDs[:min(5, len(team2SeriesIDs))] {
					if events, err := s.gridClient.DownloadEvents(r.Context(), seriesID); err == nil {
						if parsed, err := s.gridClient.ParseVALEvents(events); err == nil {
							valEvents2[seriesID] = parsed
						}
					}
				}

				team1Analysis, _ = s.valAnalyzer.AnalyzeTeam(r.Context(), team1, report.Team1Name, team1States, valEvents1)
				team2Analysis, _ = s.valAnalyzer.AnalyzeTeam(r.Context(), team2, report.Team2Name, team2States, valEvents2)
			}

			// Generate style comparison
			if team1Analysis != nil && team2Analysis != nil {
				styleComp := s.headToHeadAnalyzer.CompareStyles(team1Analysis, team2Analysis)
				report.StyleComparison = styleComp

				// Generate additional insights from style comparison
				additionalInsights := s.headToHeadAnalyzer.GenerateMatchupInsights(report, team1Analysis, team2Analysis)
				report.Insights = append(report.Insights, additionalInsights...)

				// Recalculate confidence with style data
				report.ConfidenceScore = s.calculateEnhancedConfidence(report, team1Analysis, team2Analysis)
			}
		}
	}

	respondJSON(w, http.StatusOK, report)
}

// calculateEnhancedConfidence calculates confidence score with style comparison data
func (s *Server) calculateEnhancedConfidence(report *intelligence.HeadToHeadReport, team1Analysis, team2Analysis *intelligence.TeamAnalysis) float64 {
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
		confidence += 15
	}

	// Boost for team analysis data quality
	if team1Analysis != nil && team2Analysis != nil {
		if team1Analysis.MatchesAnalyzed >= 5 && team2Analysis.MatchesAnalyzed >= 5 {
			confidence += 10
		}
	}

	// Boost for insights
	confidence += float64(len(report.Insights)) * 1.5

	if confidence > 100 {
		confidence = 100
	}

	return confidence
}
