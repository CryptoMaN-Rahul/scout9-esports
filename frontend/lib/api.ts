// Scout9 API Client
// Handles all communication with the Go backend on Azure

import type {
  Title,
  Tournament,
  Team,
  ScoutingReport,
  ReportListItem,
  GenerateReportRequest,
  HeadToHeadAnalysis,
  GameTitle,
  Series,
} from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://scout9-api.agreeableocean-9805fd33.eastus.azurecontainerapps.io';

// Default timeout for API requests (5 minutes for report generation)
const DEFAULT_TIMEOUT = 300000;

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function fetchApi<T>(endpoint: string, options?: RequestInit & { timeout?: number }): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const timeout = options?.timeout ?? DEFAULT_TIMEOUT;
  
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);
  
  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new ApiError(response.status, errorData.error || errorData.message || 'Request failed');
    }

    return response.json();
  } catch (err) {
    clearTimeout(timeoutId);
    if (err instanceof Error && err.name === 'AbortError') {
      throw new ApiError(408, 'Request timed out');
    }
    throw err;
  }
}

// === Health Check ===
export async function checkHealth(): Promise<{ status: string; service: string }> {
  return fetchApi('/health');
}

// === Titles (Games) ===
export async function getTitles(): Promise<Title[]> {
  return fetchApi('/api/titles');
}

// === Tournaments ===
export async function getTournaments(title: GameTitle): Promise<Tournament[]> {
  return fetchApi(`/api/tournaments?title=${title}`);
}

// === Teams ===
export async function getTeamsByTournament(tournamentId: string): Promise<Team[]> {
  return fetchApi(`/api/teams?tournament=${tournamentId}`);
}

export async function searchTeams(query: string, title?: GameTitle): Promise<Team[]> {
  const params = new URLSearchParams({ q: query });
  if (title) params.append('title', title);
  return fetchApi(`/api/teams/search?${params}`);
}

export async function getTeamById(teamId: string): Promise<Team> {
  return fetchApi(`/api/teams/${teamId}`);
}

export async function getTeamSeries(teamId: string, limit = 10): Promise<Series[]> {
  return fetchApi(`/api/teams/${teamId}/series?limit=${limit}`);
}

// === Reports ===
export async function generateReport(request: GenerateReportRequest): Promise<ScoutingReport> {
  return fetchApi('/api/reports/generate', {
    method: 'POST',
    body: JSON.stringify({
      teamId: request.teamId,
      teamName: request.teamName,
      matchCount: request.matchCount || 10,
      titleId: request.titleId || '3', // Default to LoL
    }),
  });
}

export async function getReport(reportId: string): Promise<ScoutingReport> {
  return fetchApi(`/api/reports/${reportId}`);
}

export async function listReports(): Promise<ReportListItem[]> {
  return fetchApi('/api/reports');
}

// === Matchups (Head-to-Head) ===
export async function getMatchup(
  team1Id: string,
  team2Id: string,
  title: GameTitle,
  matches = 20
): Promise<HeadToHeadAnalysis> {
  const params = new URLSearchParams({
    team1: team1Id,
    team2: team2Id,
    title,
    matches: matches.toString(),
  });
  return fetchApi(`/api/matchup?${params}`);
}

// === Series State ===
export async function getSeriesState(seriesId: string): Promise<unknown> {
  return fetchApi(`/api/series/${seriesId}`);
}

// Export the error class for type checking
export { ApiError };
