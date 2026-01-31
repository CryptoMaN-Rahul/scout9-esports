'use client';

import { useState, useEffect, Suspense, useMemo } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Separator } from '@/components/ui/separator';
import { 
  Search, 
  ArrowLeft, 
  Swords, 
  Trophy,
  TrendingUp,
  TrendingDown,
  Minus,
  Users,
  Target,
  Zap,
  Loader2,
  ChevronRight,
  Flame,
  Shield
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { getMatchup, ApiError } from '@/lib/api';
import { ALL_TEAMS, FEATURED_TEAMS, searchTeamsLocally } from '@/lib/teams-data';
import type { Team, HeadToHeadData } from '@/types';

// Mock H2H data for fallback
const MOCK_H2H: HeadToHeadData = {
  team1: { id: '47494', name: 'T1', logoUrl: 'https://cdn.grid.gg/assets/team-logos/1e7311945adc58ac807ffcf10b18d002' },
  team2: { id: '47558', name: 'Gen.G Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d2eded1af01ce76afb9540de0ef8b1d8' },
  matchHistory: [
    {
      matchId: 'm1',
      date: '2024-01-15',
      winner: '1',
      score: '2-1',
      tournamentName: 'LCK Spring 2024',
      games: [
        { gameNumber: 1, winnerId: '1', duration: 32 },
        { gameNumber: 2, winnerId: '2', duration: 28 },
        { gameNumber: 3, winnerId: '1', duration: 35 },
      ],
    },
    {
      matchId: 'm2',
      date: '2024-01-02',
      winner: '2',
      score: '2-0',
      tournamentName: 'LCK Winter Finals',
      games: [
        { gameNumber: 1, winnerId: '2', duration: 29 },
        { gameNumber: 2, winnerId: '2', duration: 31 },
      ],
    },
  ],
  stats: {
    totalMatches: 8,
    team1Wins: 5,
    team2Wins: 3,
    avgGameDuration: 32,
    commonPicks: {
      team1: ['Nautilus', 'Xin Zhao', 'Ahri'],
      team2: ['Thresh', 'Lee Sin', 'Azir'],
    },
    keyMatchups: [
      {
        player1: { playerId: '1', nickname: 'Faker', role: 'Mid', winRate: 0.62, avgKDA: 5.2 },
        player2: { playerId: '2', nickname: 'Chovy', role: 'Mid', winRate: 0.38, avgKDA: 4.8 },
        gamesPlayed: 8,
        significance: 'The mid lane matchup often determines the game outcome',
      },
      {
        player1: { playerId: '3', nickname: 'Keria', role: 'Support', winRate: 0.70, avgKDA: 6.1 },
        player2: { playerId: '4', nickname: 'Lehends', role: 'Support', winRate: 0.30, avgKDA: 3.5 },
        gamesPlayed: 8,
        significance: 'Keria dominates lane phase and roam timings',
      },
    ],
  },
};

function TeamSearchCard({
  title,
  selectedTeam,
  onSelectTeam,
  otherTeamId,
  game,
}: {
  title: string;
  selectedTeam: Team | null;
  onSelectTeam: (team: Team) => void;
  otherTeamId: string | null;
  game: 'lol' | 'valorant';
}) {
  const [query, setQuery] = useState('');
  const [showResults, setShowResults] = useState(false);

  // Use local search - instant, no API calls needed
  const results = useMemo(() => {
    if (query.length < 2) return [];
    return searchTeamsLocally(query, game).filter(t => t.id !== otherTeamId);
  }, [query, game, otherTeamId]);

  // Featured teams for initial display
  const featuredTeams = useMemo(() => {
    return FEATURED_TEAMS[game].filter(t => t.id !== otherTeamId).slice(0, 8);
  }, [game, otherTeamId]);

  const handleSelect = (team: Team) => {
    onSelectTeam(team);
    setQuery('');
    setShowResults(false);
  };

  return (
    <Card className="glass-card flex-1">
      <CardHeader>
        <CardTitle className="text-lg">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {selectedTeam ? (
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              {selectedTeam.logoUrl ? (
                <img 
                  src={selectedTeam.logoUrl} 
                  alt={selectedTeam.name} 
                  className="w-12 h-12 rounded-lg object-contain bg-muted p-1"
                />
              ) : (
                <div className="w-12 h-12 rounded-lg bg-game-primary/20 flex items-center justify-center text-game-primary font-bold text-xl">
                  {selectedTeam.name.charAt(0)}
                </div>
              )}
              <div>
                <p className="font-bold text-lg">{selectedTeam.name}</p>
              </div>
            </div>
            <Button variant="ghost" size="sm" onClick={() => onSelectTeam(null as any)}>
              Change
            </Button>
          </div>
        ) : (
          <div className="relative">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                placeholder="Search team..."
                value={query}
                onChange={(e) => {
                  setQuery(e.target.value);
                  setShowResults(e.target.value.length >= 2);
                }}
                className="pl-10"
              />
            </div>

            <AnimatePresence>
              {showResults && results.length > 0 && (
                <motion.div
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  className="absolute w-full mt-2 bg-card border border-border rounded-lg shadow-xl z-10 overflow-hidden max-h-64 overflow-y-auto"
                >
                  {results.map((team) => (
                    <button
                      key={team.id}
                      className="w-full flex items-center gap-3 p-3 hover:bg-muted transition-colors text-left"
                      onClick={() => handleSelect(team)}
                    >
                      {team.logoUrl ? (
                        <img 
                          src={team.logoUrl} 
                          alt={team.name} 
                          className="w-8 h-8 rounded object-contain bg-muted p-0.5"
                        />
                      ) : (
                        <div className="w-8 h-8 rounded bg-game-primary/20 flex items-center justify-center text-game-primary font-semibold">
                          {team.name.charAt(0)}
                        </div>
                      )}
                      <span>{team.name}</span>
                    </button>
                  ))}
                </motion.div>
              )}
            </AnimatePresence>

            {/* Featured teams for quick select */}
            {!showResults && (
              <div className="mt-4">
                <p className="text-sm text-muted-foreground mb-2">Select a Team:</p>
                <div className="grid grid-cols-2 gap-2">
                  {featuredTeams.map((team) => (
                    <button
                      key={team.id}
                      onClick={() => handleSelect(team)}
                      className="flex items-center gap-2 p-2 rounded-lg border border-border hover:border-muted-foreground hover:bg-muted/50 transition-all text-left text-sm"
                    >
                      {team.logoUrl ? (
                        <img 
                          src={team.logoUrl} 
                          alt={team.name} 
                          className="w-6 h-6 rounded object-contain bg-muted p-0.5"
                        />
                      ) : (
                        <div className="w-6 h-6 rounded bg-muted flex items-center justify-center text-xs font-bold">
                          {team.name.substring(0, 2).toUpperCase()}
                        </div>
                      )}
                      <span className="font-medium truncate">{team.name}</span>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function ComparisonContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [game, setGame] = useState<'lol' | 'valorant'>('lol');
  const [team1, setTeam1] = useState<Team | null>(null);
  const [team2, setTeam2] = useState<Team | null>(null);
  const [h2hData, setH2hData] = useState<HeadToHeadData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Check URL params for pre-selected teams
  useEffect(() => {
    const gameParam = searchParams.get('game');
    if (gameParam === 'lol' || gameParam === 'valorant') {
      setGame(gameParam);
    }
  }, [searchParams]);

  const canCompare = team1 && team2;

  const handleCompare = async () => {
    if (!team1 || !team2) return;

    setLoading(true);
    setError(null);

    try {
      const apiData = await getMatchup(team1.id, team2.id, game);
      // Transform API response to HeadToHeadData format
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const rawData = apiData as any;
      const transformedData: HeadToHeadData = {
        team1,
        team2,
        stats: {
          totalMatches: rawData.totalMatches || 0,
          team1Wins: rawData.team1Wins || 0,
          team2Wins: rawData.team2Wins || 0,
          avgGameDuration: rawData.avgGameDuration,
        },
        styleComparison: rawData.styleComparison,
        insights: rawData.insights?.map((i: { text: string }) => i.text) || [],
        warnings: rawData.warnings || [],
        confidenceScore: rawData.confidenceScore || 0,
      };
      setH2hData(transformedData);
    } catch (err) {
      console.error('Failed to get matchup:', err);
      // Use mock data with team names
      setH2hData({
        ...MOCK_H2H,
        team1,
        team2,
      });
      setError('Using demo data - backend connection failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={cn("min-h-screen pb-16", game === 'lol' ? 'theme-lol' : 'theme-valorant')}>
      {/* Header */}
      <div className="sticky top-16 z-40 bg-background/80 backdrop-blur-xl border-b border-border">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Button 
              variant="ghost" 
              onClick={() => router.push('/')}
              className="gap-2"
            >
              <ArrowLeft className="w-4 h-4" />
              Back
            </Button>

            {/* Game Toggle */}
            <div className="flex gap-2">
              <Button
                variant={game === 'lol' ? 'game' : 'outline'}
                size="sm"
                onClick={() => { setGame('lol'); setH2hData(null); }}
              >
                League of Legends
              </Button>
              <Button
                variant={game === 'valorant' ? 'game' : 'outline'}
                size="sm"
                onClick={() => { setGame('valorant'); setH2hData(null); }}
              >
                VALORANT
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Title */}
        <div className="text-center mb-8">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-game-primary to-game-secondary mb-4"
          >
            <Swords className="w-8 h-8 text-white" />
          </motion.div>
          <h1 className="text-3xl font-bold mb-2">Head-to-Head Analysis</h1>
          <p className="text-muted-foreground">
            Compare two teams to discover historical matchup data and insights
          </p>
        </div>

        {/* Team Selection */}
        <div className="grid md:grid-cols-2 gap-6 mb-8">
          <TeamSearchCard
            title="Team 1"
            selectedTeam={team1}
            onSelectTeam={setTeam1}
            otherTeamId={team2?.id || null}
            game={game}
          />
          <TeamSearchCard
            title="Team 2"
            selectedTeam={team2}
            onSelectTeam={setTeam2}
            otherTeamId={team1?.id || null}
            game={game}
          />
        </div>

        {/* Compare Button */}
        <div className="text-center mb-8">
          <Button
            size="lg"
            variant="game"
            disabled={!canCompare || loading}
            onClick={handleCompare}
            className="gap-2"
          >
            {loading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Analyzing...
              </>
            ) : (
              <>
                <Swords className="w-5 h-5" />
                Compare Teams
              </>
            )}
          </Button>
        </div>

        {/* Error Banner */}
        {error && (
          <div className="p-3 rounded-lg bg-yellow-500/10 border border-yellow-500/30 text-yellow-500 text-sm flex items-center gap-2 mb-8">
            {error}
          </div>
        )}

        {/* Results */}
        <AnimatePresence mode="wait">
          {h2hData && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="space-y-8"
            >
              {/* Overall Stats */}
              <Card className="glass-card overflow-hidden">
                <div className="relative">
                  {/* Background gradient showing advantage */}
                  <div 
                    className="absolute inset-0 opacity-20"
                    style={{
                      background: `linear-gradient(90deg, 
                        hsl(var(--game-primary)) ${h2hData.stats?.totalMatches ? (h2hData.stats.team1Wins / h2hData.stats.totalMatches) * 100 : 50}%, 
                        hsl(var(--game-secondary)) ${h2hData.stats?.totalMatches ? (h2hData.stats.team1Wins / h2hData.stats.totalMatches) * 100 : 50}%)`
                    }}
                  />
                  
                  <CardContent className="relative p-8">
                    <div className="flex items-center justify-between">
                      {/* Team 1 */}
                      <div className="text-center flex-1">
                        {h2hData.team1.logoUrl ? (
                          <img 
                            src={h2hData.team1.logoUrl} 
                            alt={h2hData.team1.name} 
                            className="w-20 h-20 mx-auto rounded-2xl object-contain bg-muted p-2 mb-3"
                          />
                        ) : (
                          <div className="w-20 h-20 mx-auto rounded-2xl bg-game-primary/20 flex items-center justify-center text-game-primary font-bold text-3xl mb-3">
                            {h2hData.team1.name.charAt(0)}
                          </div>
                        )}
                        <h3 className="text-xl font-bold">{h2hData.team1.name}</h3>
                        <div className="text-4xl font-bold text-game-primary mt-2">
                          {h2hData.stats?.team1Wins ?? 0}
                        </div>
                        <p className="text-sm text-muted-foreground">Wins</p>
                      </div>

                      {/* VS */}
                      <div className="px-8">
                        <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
                          <span className="text-xl font-bold text-muted-foreground">VS</span>
                        </div>
                        <p className="text-center text-sm text-muted-foreground mt-2">
                          {h2hData.stats?.totalMatches ?? 0} matches
                        </p>
                      </div>

                      {/* Team 2 */}
                      <div className="text-center flex-1">
                        {h2hData.team2.logoUrl ? (
                          <img 
                            src={h2hData.team2.logoUrl} 
                            alt={h2hData.team2.name} 
                            className="w-20 h-20 mx-auto rounded-2xl object-contain bg-muted p-2 mb-3"
                          />
                        ) : (
                          <div className="w-20 h-20 mx-auto rounded-2xl bg-game-secondary/20 flex items-center justify-center text-game-secondary font-bold text-3xl mb-3">
                            {h2hData.team2.name.charAt(0)}
                          </div>
                        )}
                        <h3 className="text-xl font-bold">{h2hData.team2.name}</h3>
                        <div className="text-4xl font-bold text-game-secondary mt-2">
                          {h2hData.stats?.team2Wins ?? 0}
                        </div>
                        <p className="text-sm text-muted-foreground">Wins</p>
                      </div>
                    </div>

                    {/* Win rate bar */}
                    <div className="mt-8">
                      <div className="h-3 rounded-full bg-muted overflow-hidden flex">
                        <div 
                          className="bg-game-primary transition-all duration-500"
                          style={{ width: `${h2hData.stats?.totalMatches ? (h2hData.stats.team1Wins / h2hData.stats.totalMatches) * 100 : 50}%` }}
                        />
                        <div 
                          className="bg-game-secondary transition-all duration-500"
                          style={{ width: `${h2hData.stats?.totalMatches ? (h2hData.stats.team2Wins / h2hData.stats.totalMatches) * 100 : 50}%` }}
                        />
                      </div>
                      <div className="flex justify-between mt-2 text-sm text-muted-foreground">
                        <span>{h2hData.stats?.totalMatches ? Math.round((h2hData.stats.team1Wins / h2hData.stats.totalMatches) * 100) : 50}%</span>
                        <span>{h2hData.stats?.totalMatches ? Math.round((h2hData.stats.team2Wins / h2hData.stats.totalMatches) * 100) : 50}%</span>
                      </div>
                    </div>

                    {/* Confidence Score */}
                    {h2hData.confidenceScore && (
                      <div className="mt-4 text-center">
                        <Badge variant="secondary" className="text-sm">
                          Analysis Confidence: {Math.round(h2hData.confidenceScore)}%
                        </Badge>
                      </div>
                    )}
                  </CardContent>
                </div>
              </Card>

              {/* Style Comparison - Game Phase Analysis */}
              {h2hData.styleComparison && (
                <Card className="glass-card">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Zap className="w-5 h-5 text-game-primary" />
                      Style Comparison
                    </CardTitle>
                    {h2hData.styleComparison.styleInsight && (
                      <p className="text-sm text-muted-foreground mt-1">
                        {h2hData.styleComparison.styleInsight}
                      </p>
                    )}
                  </CardHeader>
                  <CardContent className="space-y-6">
                    {/* Early Game */}
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-medium flex items-center gap-2">
                          <Target className="w-4 h-4" />
                          {game === 'lol' ? 'Early Game' : 'Pistol Rounds'}
                        </span>
                        <Badge 
                          variant={h2hData.styleComparison.earlyGameAdvantage === 'team1' ? 'default' : 
                                   h2hData.styleComparison.earlyGameAdvantage === 'team2' ? 'secondary' : 'outline'}
                          className={h2hData.styleComparison.earlyGameAdvantage === 'team1' ? 'bg-game-primary' : 
                                     h2hData.styleComparison.earlyGameAdvantage === 'team2' ? 'bg-game-secondary' : ''}
                        >
                          {h2hData.styleComparison.earlyGameAdvantage === 'team1' ? h2hData.team1.name :
                           h2hData.styleComparison.earlyGameAdvantage === 'team2' ? h2hData.team2.name : 'Even'}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-sm w-16 text-right font-semibold text-game-primary">
                          {Math.round(h2hData.styleComparison.team1EarlyGameRating || 0)}
                        </span>
                        <div className="flex-1 h-3 rounded-full bg-muted overflow-hidden flex">
                          <div 
                            className="bg-game-primary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team1EarlyGameRating || 0) / 
                                       ((h2hData.styleComparison.team1EarlyGameRating || 0) + (h2hData.styleComparison.team2EarlyGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                          <div 
                            className="bg-game-secondary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team2EarlyGameRating || 0) / 
                                       ((h2hData.styleComparison.team1EarlyGameRating || 0) + (h2hData.styleComparison.team2EarlyGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                        </div>
                        <span className="text-sm w-16 font-semibold text-game-secondary">
                          {Math.round(h2hData.styleComparison.team2EarlyGameRating || 0)}
                        </span>
                      </div>
                      {h2hData.styleComparison.earlyGameInsight && (
                        <p className="text-xs text-muted-foreground mt-1">{h2hData.styleComparison.earlyGameInsight}</p>
                      )}
                    </div>

                    {/* Mid Game */}
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-medium flex items-center gap-2">
                          <Swords className="w-4 h-4" />
                          {game === 'lol' ? 'Mid Game' : 'Attack Side'}
                        </span>
                        <Badge 
                          variant={h2hData.styleComparison.midGameAdvantage === 'team1' ? 'default' : 
                                   h2hData.styleComparison.midGameAdvantage === 'team2' ? 'secondary' : 'outline'}
                          className={h2hData.styleComparison.midGameAdvantage === 'team1' ? 'bg-game-primary' : 
                                     h2hData.styleComparison.midGameAdvantage === 'team2' ? 'bg-game-secondary' : ''}
                        >
                          {h2hData.styleComparison.midGameAdvantage === 'team1' ? h2hData.team1.name :
                           h2hData.styleComparison.midGameAdvantage === 'team2' ? h2hData.team2.name : 'Even'}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-sm w-16 text-right font-semibold text-game-primary">
                          {Math.round(h2hData.styleComparison.team1MidGameRating || 0)}
                        </span>
                        <div className="flex-1 h-3 rounded-full bg-muted overflow-hidden flex">
                          <div 
                            className="bg-game-primary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team1MidGameRating || 0) / 
                                       ((h2hData.styleComparison.team1MidGameRating || 0) + (h2hData.styleComparison.team2MidGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                          <div 
                            className="bg-game-secondary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team2MidGameRating || 0) / 
                                       ((h2hData.styleComparison.team1MidGameRating || 0) + (h2hData.styleComparison.team2MidGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                        </div>
                        <span className="text-sm w-16 font-semibold text-game-secondary">
                          {Math.round(h2hData.styleComparison.team2MidGameRating || 0)}
                        </span>
                      </div>
                      {h2hData.styleComparison.midGameInsight && (
                        <p className="text-xs text-muted-foreground mt-1">{h2hData.styleComparison.midGameInsight}</p>
                      )}
                    </div>

                    {/* Late Game */}
                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-medium flex items-center gap-2">
                          <Trophy className="w-4 h-4" />
                          {game === 'lol' ? 'Late Game' : 'Defense Side'}
                        </span>
                        <Badge 
                          variant={h2hData.styleComparison.lateGameAdvantage === 'team1' ? 'default' : 
                                   h2hData.styleComparison.lateGameAdvantage === 'team2' ? 'secondary' : 'outline'}
                          className={h2hData.styleComparison.lateGameAdvantage === 'team1' ? 'bg-game-primary' : 
                                     h2hData.styleComparison.lateGameAdvantage === 'team2' ? 'bg-game-secondary' : ''}
                        >
                          {h2hData.styleComparison.lateGameAdvantage === 'team1' ? h2hData.team1.name :
                           h2hData.styleComparison.lateGameAdvantage === 'team2' ? h2hData.team2.name : 'Even'}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-4">
                        <span className="text-sm w-16 text-right font-semibold text-game-primary">
                          {Math.round(h2hData.styleComparison.team1LateGameRating || 0)}
                        </span>
                        <div className="flex-1 h-3 rounded-full bg-muted overflow-hidden flex">
                          <div 
                            className="bg-game-primary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team1LateGameRating || 0) / 
                                       ((h2hData.styleComparison.team1LateGameRating || 0) + (h2hData.styleComparison.team2LateGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                          <div 
                            className="bg-game-secondary transition-all duration-500"
                            style={{ 
                              width: `${((h2hData.styleComparison.team2LateGameRating || 0) / 
                                       ((h2hData.styleComparison.team1LateGameRating || 0) + (h2hData.styleComparison.team2LateGameRating || 0) || 1)) * 100}%` 
                            }}
                          />
                        </div>
                        <span className="text-sm w-16 font-semibold text-game-secondary">
                          {Math.round(h2hData.styleComparison.team2LateGameRating || 0)}
                        </span>
                      </div>
                      {h2hData.styleComparison.lateGameInsight && (
                        <p className="text-xs text-muted-foreground mt-1">{h2hData.styleComparison.lateGameInsight}</p>
                      )}
                    </div>

                    {/* Aggression Score */}
                    {(h2hData.styleComparison.team1Aggression !== undefined || h2hData.styleComparison.team2Aggression !== undefined) && (
                      <div>
                        <Separator className="my-4" />
                        <div className="flex items-center justify-between mb-2">
                          <span className="font-medium flex items-center gap-2">
                            <Flame className="w-4 h-4" />
                            Aggression Level
                          </span>
                        </div>
                        <div className="flex items-center gap-4">
                          <div className="flex-1">
                            <div className="flex items-center justify-between text-sm mb-1">
                              <span>{h2hData.team1.name}</span>
                              <span className="font-semibold text-game-primary">{Math.round(h2hData.styleComparison.team1Aggression || 0)}</span>
                            </div>
                            <Progress value={h2hData.styleComparison.team1Aggression || 0} className="h-2" />
                          </div>
                          <div className="flex-1">
                            <div className="flex items-center justify-between text-sm mb-1">
                              <span>{h2hData.team2.name}</span>
                              <span className="font-semibold text-game-secondary">{Math.round(h2hData.styleComparison.team2Aggression || 0)}</span>
                            </div>
                            <Progress value={h2hData.styleComparison.team2Aggression || 0} className="h-2" />
                          </div>
                        </div>
                      </div>
                    )}
                  </CardContent>
                </Card>
              )}

              {/* Insights */}
              {h2hData.insights && h2hData.insights.length > 0 && (
                <Card className="glass-card">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <TrendingUp className="w-5 h-5 text-game-primary" />
                      Key Insights
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <ul className="space-y-2">
                      {h2hData.insights.map((insight, i) => (
                        <li key={i} className="flex items-start gap-2">
                          <ChevronRight className="w-4 h-4 text-game-primary mt-0.5 flex-shrink-0" />
                          <span className="text-sm">{insight}</span>
                        </li>
                      ))}
                    </ul>
                  </CardContent>
                </Card>
              )}
              {h2hData.stats?.keyMatchups && h2hData.stats.keyMatchups.length > 0 && (
              <div>
                <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                  <Users className="w-5 h-5 text-game-primary" />
                  Key Player Matchups
                </h2>
                <div className="grid gap-4">
                  {h2hData.stats.keyMatchups.map((matchup, i) => (
                    <Card key={i} className="glass-card">
                      <CardContent className="p-6">
                        <div className="flex items-center justify-between">
                          {/* Player 1 */}
                          <div className="flex-1">
                            <div className="flex items-center gap-3">
                              <div className="w-12 h-12 rounded-lg bg-game-primary/20 flex items-center justify-center text-game-primary font-bold">
                                {matchup.player1.nickname.charAt(0)}
                              </div>
                              <div>
                                <p className="font-bold">{matchup.player1.nickname}</p>
                                <p className="text-sm text-muted-foreground">{matchup.player1.role}</p>
                              </div>
                            </div>
                            <div className="mt-3 grid grid-cols-2 gap-2 text-sm">
                              <div>
                                <span className="text-muted-foreground">Win Rate:</span>
                                <span className={cn(
                                  "ml-2 font-semibold",
                                  matchup.player1.winRate > matchup.player2.winRate ? "text-green-500" : "text-muted-foreground"
                                )}>
                                  {Math.round(matchup.player1.winRate * 100)}%
                                </span>
                              </div>
                              <div>
                                <span className="text-muted-foreground">Avg KDA:</span>
                                <span className="ml-2 font-semibold">{matchup.player1.avgKDA.toFixed(1)}</span>
                              </div>
                            </div>
                          </div>

                          {/* VS indicator */}
                          <div className="px-4">
                            <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
                              {matchup.player1.winRate > matchup.player2.winRate ? (
                                <ChevronRight className="w-5 h-5 text-game-primary" />
                              ) : matchup.player1.winRate < matchup.player2.winRate ? (
                                <ChevronRight className="w-5 h-5 text-game-secondary rotate-180" />
                              ) : (
                                <Minus className="w-5 h-5" />
                              )}
                            </div>
                          </div>

                          {/* Player 2 */}
                          <div className="flex-1 text-right">
                            <div className="flex items-center gap-3 justify-end">
                              <div>
                                <p className="font-bold">{matchup.player2.nickname}</p>
                                <p className="text-sm text-muted-foreground">{matchup.player2.role}</p>
                              </div>
                              <div className="w-12 h-12 rounded-lg bg-game-secondary/20 flex items-center justify-center text-game-secondary font-bold">
                                {matchup.player2.nickname.charAt(0)}
                              </div>
                            </div>
                            <div className="mt-3 grid grid-cols-2 gap-2 text-sm">
                              <div>
                                <span className="text-muted-foreground">Win Rate:</span>
                                <span className={cn(
                                  "ml-2 font-semibold",
                                  matchup.player2.winRate > matchup.player1.winRate ? "text-green-500" : "text-muted-foreground"
                                )}>
                                  {Math.round(matchup.player2.winRate * 100)}%
                                </span>
                              </div>
                              <div>
                                <span className="text-muted-foreground">Avg KDA:</span>
                                <span className="ml-2 font-semibold">{matchup.player2.avgKDA.toFixed(1)}</span>
                              </div>
                            </div>
                          </div>
                        </div>

                        {matchup.significance && (
                          <p className="text-sm text-muted-foreground mt-4 text-center italic">
                            &quot;{matchup.significance}&quot;
                          </p>
                        )}
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
              )}

              {/* Match History */}
              {h2hData.matchHistory && h2hData.matchHistory.length > 0 && (
              <div>
                <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                  <Trophy className="w-5 h-5 text-game-primary" />
                  Recent Match History
                </h2>
                <div className="space-y-3">
                  {h2hData.matchHistory.map((match) => {
                    const team1Won = match.winner === h2hData.team1.id;
                    return (
                      <Card key={match.matchId} className="glass-card">
                        <CardContent className="p-4">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-4">
                              <div className={cn(
                                "w-10 h-10 rounded-lg flex items-center justify-center font-bold",
                                team1Won ? "bg-game-primary/20 text-game-primary" : "bg-game-secondary/20 text-game-secondary"
                              )}>
                                {team1Won ? h2hData.team1.name.charAt(0) : h2hData.team2.name.charAt(0)}
                              </div>
                              <div>
                                <p className="font-semibold">
                                  {team1Won ? h2hData.team1.name : h2hData.team2.name} wins
                                </p>
                                <p className="text-sm text-muted-foreground">{match.tournamentName}</p>
                              </div>
                            </div>
                            <div className="text-right">
                              <p className="text-2xl font-bold">{match.score}</p>
                              <p className="text-sm text-muted-foreground">
                                {new Date(match.date).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })}
                </div>
              </div>
              )}

              {/* Common Picks */}
              {h2hData.stats?.commonPicks && (
              <div className="grid md:grid-cols-2 gap-6">
                <Card className="glass-card">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Flame className="w-5 h-5 text-game-primary" />
                      {h2hData.team1.name}&apos;s Go-To Picks
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-wrap gap-2">
                      {h2hData.stats.commonPicks.team1.map((pick) => (
                        <Badge key={pick} variant="secondary" className="text-sm">
                          {pick}
                        </Badge>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                <Card className="glass-card">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Shield className="w-5 h-5 text-game-secondary" />
                      {h2hData.team2.name}&apos;s Go-To Picks
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-wrap gap-2">
                      {h2hData.stats.commonPicks.team2.map((pick) => (
                        <Badge key={pick} variant="secondary" className="text-sm">
                          {pick}
                        </Badge>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

export default function ComparePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    }>
      <ComparisonContent />
    </Suspense>
  );
}
