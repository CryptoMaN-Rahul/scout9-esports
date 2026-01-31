'use client';

import { useEffect, useState, use } from 'react';
import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Separator } from '@/components/ui/separator';
import { 
  ArrowLeft, 
  Download, 
  Share2, 
  Printer,
  Loader2,
  AlertCircle
} from 'lucide-react';
import { ExecutiveSummary } from '@/components/report/ExecutiveSummary';
import { CommonStrategies } from '@/components/report/CommonStrategies';
import { PlayerTendencies } from '@/components/report/PlayerTendencies';
import { Compositions } from '@/components/report/Compositions';
import { HowToWin } from '@/components/report/HowToWin';
import { TrendAnalysis } from '@/components/report/TrendAnalysis';
import { getReport, ApiError } from '@/lib/api';
import { cn } from '@/lib/utils';
import type { ScoutingReport, CommonStrategiesSection } from '@/types';

// Mock data for demonstration (when API fails)
const MOCK_REPORT: ScoutingReport = {
  id: 'mock-1',
  generatedAt: new Date().toISOString(),
  opponentTeam: { id: '47494', name: 'T1', logoUrl: 'https://cdn.grid.gg/assets/team-logos/1e7311945adc58ac807ffcf10b18d002' },
  title: 'lol',
  matchesAnalyzed: 10,
  executiveSummary: 'T1 is an aggressive early-game focused team with a 65% win rate over the last 10 matches. They prioritize dragon control (72% first dragon rate) and excel at snowballing leads through their bot lane. Key weakness: they struggle in games lasting over 35 minutes (only 45% win rate). Their jungler "Oner" paths bot-side 75% of the time pre-10 minutes. Recommend drafting scaling compositions and contesting their early dragon priority.',
  howToWin: {
    winCondition: 'Survive their early aggression, scale to 35+ minutes, and target their weak late-game teamfighting. Force the game to go late by trading objectives rather than fighting.',
    confidenceScore: 78,
    actionableInsights: [
      {
        recommendation: 'Draft scaling compositions',
        dataBacking: 'T1 has only 45% win rate in games over 35 minutes',
        impact: 'HIGH',
        actionType: 'STRATEGY',
        confidence: 85,
      },
      {
        recommendation: 'Ban Nautilus from Keria',
        dataBacking: '80% pick rate, 75% win rate on this champion',
        impact: 'HIGH',
        actionType: 'BAN',
        confidence: 90,
      },
      {
        recommendation: 'Target top-side early',
        dataBacking: 'Oner paths bot 75% of time, leaving Zeus vulnerable',
        impact: 'MEDIUM',
        actionType: 'STRATEGY',
        confidence: 72,
      },
      {
        recommendation: 'Pick assassins mid',
        dataBacking: 'Faker has 2.1 KDA on control mages vs 8.5 on assassins - force uncomfortable picks',
        impact: 'MEDIUM',
        actionType: 'PICK',
        confidence: 68,
      },
    ],
    draftStrategy: {
      priorityBans: [
        { character: 'Nautilus', reason: '80% pick rate, 75% WR from Keria', playerName: 'Keria', winRate: 0.75 },
        { character: 'Xin Zhao', reason: 'Oner signature - enables early aggression', playerName: 'Oner', winRate: 0.70 },
      ],
      targetPicks: [
        { character: 'Control Mages', reason: 'Force Faker onto uncomfortable champions', playerName: 'Faker' },
      ],
      recommendedPicks: [
        { character: 'Ornn', reason: 'Scale into late game, team utility', winRate: 0.55 },
        { character: 'Azir', reason: 'Late-game insurance, zone control', winRate: 0.52 },
        { character: 'Jinx', reason: 'Hypercarry to exploit weak late-game', winRate: 0.58 },
      ],
    },
    inGameStrategy: [
      { phase: 'Early Game', timing: '0-10 min', strategy: 'Play safe, concede first dragon if needed, focus on farm', priority: 'HIGH' },
      { phase: 'Counter-Jungle', timing: '3-6 min', strategy: 'Invade top-side jungle while Oner is bot', priority: 'MEDIUM' },
      { phase: 'Mid Game', timing: '15-25 min', strategy: 'Trade objectives, avoid 5v5 fights until 3+ items', priority: 'HIGH' },
      { phase: 'Late Game', timing: '35+ min', strategy: 'Force teamfights - their coordination drops significantly', priority: 'HIGH' },
    ],
  },
  teamStrategy: {
    teamId: '1',
    teamName: 'T1',
    title: 'lol',
    matchesAnalyzed: 10,
    gamesAnalyzed: 22,
    winRate: 0.65,
    recentForm: 'improving',
    strengths: [
      { text: 'Excellent early game execution with 72% first dragon rate', importance: 'HIGH' },
      { text: 'Strong bot lane synergy - Gumayusi/Keria average 2.3 KDA', importance: 'HIGH' },
      { text: 'Aggressive jungle pathing creates early leads', importance: 'MEDIUM' },
    ],
    weaknesses: [
      { text: 'Only 45% win rate in games lasting 35+ minutes', importance: 'HIGH' },
      { text: 'Over-reliance on early dragon priority', importance: 'MEDIUM' },
      { text: 'Faker struggles on control mages (2.1 KDA vs 8.5 on assassins)', importance: 'HIGH' },
    ],
    lolMetrics: {
      firstBloodRate: 0.58,
      firstDragonRate: 0.72,
      firstTowerRate: 0.65,
      firstTowerAvgTime: 12.5,
      goldDiff15: 1250,
      dragonControlRate: 0.68,
      heraldControlRate: 0.55,
      baronControlRate: 0.60,
      elderDragonRate: 0.45,
      avgGameDuration: 31,
      earlyGameRating: 82,
      midGameRating: 70,
      lateGameRating: 48,
      aggressionScore: 75,
      winConditions: ['Early Snowball', 'Dragon Soul'],
    },
  },
  playerProfiles: [
    {
      playerId: '1',
      nickname: 'Zeus',
      role: 'Top',
      teamId: '1',
      gamesPlayed: 22,
      kda: 3.8,
      avgKills: 2.5,
      avgDeaths: 2.1,
      avgAssists: 5.5,
      characterPool: [
        { name: 'Ksante', games: 8, wins: 6, losses: 2, winRate: 0.75, avgKDA: 4.2 },
        { name: 'Renekton', games: 5, wins: 3, losses: 2, winRate: 0.60, avgKDA: 3.5 },
      ],
      signaturePicks: ['Ksante', 'Renekton', 'Jax'],
      threatLevel: 7,
      threatReason: 'Consistent performer, can carry from top lane with leads',
      weaknesses: [{ text: 'Vulnerable to early ganks when Oner is bot-side', importance: 'MEDIUM' }],
      tendencies: ['Plays aggressive in lane', 'Teleports for dragon fights'],
    },
    {
      playerId: '2',
      nickname: 'Oner',
      role: 'Jungle',
      teamId: '1',
      gamesPlayed: 22,
      kda: 4.2,
      avgKills: 3.8,
      avgDeaths: 2.5,
      avgAssists: 7.2,
      characterPool: [
        { name: 'Xin Zhao', games: 7, wins: 5, losses: 2, winRate: 0.71, avgKDA: 4.8 },
        { name: 'Lee Sin', games: 6, wins: 4, losses: 2, winRate: 0.67, avgKDA: 4.0 },
      ],
      signaturePicks: ['Xin Zhao', 'Lee Sin', 'Viego'],
      threatLevel: 8,
      threatReason: 'Primary early game catalyst - his aggression enables T1 snowballs',
      weaknesses: [{ text: 'Predictable bot-side pathing (75% of the time)', importance: 'HIGH' }],
      tendencies: ['Paths bot-side pre-10', 'Prioritizes dragon over herald'],
    },
    {
      playerId: '3',
      nickname: 'Faker',
      role: 'Mid',
      teamId: '1',
      gamesPlayed: 22,
      kda: 5.5,
      avgKills: 4.2,
      avgDeaths: 1.8,
      avgAssists: 5.7,
      characterPool: [
        { name: 'Ahri', games: 6, wins: 5, losses: 1, winRate: 0.83, avgKDA: 7.2 },
        { name: 'Azir', games: 5, wins: 2, losses: 3, winRate: 0.40, avgKDA: 2.1 },
      ],
      signaturePicks: ['Ahri', 'LeBlanc', 'Akali'],
      threatLevel: 9,
      threatReason: 'Legendary player - can solo carry on assassins',
      weaknesses: [{ text: 'Significantly weaker on control mages (2.1 KDA vs 8.5)', importance: 'HIGH' }],
      tendencies: ['Roams bot frequently', 'Plays aggressive in lane'],
    },
    {
      playerId: '4',
      nickname: 'Gumayusi',
      role: 'Bot',
      teamId: '1',
      gamesPlayed: 22,
      kda: 4.8,
      avgKills: 5.2,
      avgDeaths: 2.0,
      avgAssists: 4.4,
      characterPool: [
        { name: 'Jinx', games: 7, wins: 5, losses: 2, winRate: 0.71, avgKDA: 5.2 },
        { name: 'Zeri', games: 5, wins: 4, losses: 1, winRate: 0.80, avgKDA: 5.8 },
      ],
      signaturePicks: ['Jinx', 'Zeri', 'Aphelios'],
      threatLevel: 8,
      threatReason: 'Primary carry threat - team plays around him',
      weaknesses: [{ text: 'Positioning can be exploited in late-game fights', importance: 'MEDIUM' }],
      tendencies: ['Farms aggressively', 'Looks for 2v2 kills'],
    },
    {
      playerId: '5',
      nickname: 'Keria',
      role: 'Support',
      teamId: '1',
      gamesPlayed: 22,
      kda: 5.2,
      avgKills: 0.8,
      avgDeaths: 2.2,
      avgAssists: 10.5,
      characterPool: [
        { name: 'Nautilus', games: 9, wins: 7, losses: 2, winRate: 0.78, avgKDA: 5.8 },
        { name: 'Thresh', games: 5, wins: 3, losses: 2, winRate: 0.60, avgKDA: 4.5 },
      ],
      signaturePicks: ['Nautilus', 'Thresh', 'Leona'],
      threatLevel: 9,
      threatReason: 'Best support in the world - engage timing is exceptional',
      weaknesses: [{ text: 'Team is much weaker when his Nautilus is banned', importance: 'HIGH' }],
      tendencies: ['Roams mid after pushing waves', 'Times engages with Oner ganks'],
    },
  ],
  compositions: {
    topCompositions: [
      { characters: ['Ksante', 'Xin Zhao', 'Ahri', 'Jinx', 'Nautilus'], frequency: 0.35, winRate: 0.75, gamesPlayed: 8, archetype: 'Early Aggression' },
      { characters: ['Renekton', 'Lee Sin', 'LeBlanc', 'Zeri', 'Thresh'], frequency: 0.22, winRate: 0.60, gamesPlayed: 5, archetype: 'Pick Comp' },
      { characters: ['Jax', 'Viego', 'Akali', 'Aphelios', 'Leona'], frequency: 0.18, winRate: 0.50, gamesPlayed: 4, archetype: 'Skirmish' },
    ],
    firstPickPriorities: [
      { character: 'Nautilus', rate: 0.6, winRate: 0.7, gamesPlayed: 12 },
      { character: 'Xin Zhao', rate: 0.5, winRate: 0.65, gamesPlayed: 10 },
      { character: 'Ahri', rate: 0.45, winRate: 0.68, gamesPlayed: 9 },
      { character: 'Jinx', rate: 0.4, winRate: 0.72, gamesPlayed: 8 },
    ],
    commonBans: ['Maokai', 'Azir', 'Orianna'],
  },
};

export default function ReportPage({ params }: { params: Promise<{ id: string }> }) {
  const { id: reportId } = use(params);
  const router = useRouter();
  const [report, setReport] = useState<ScoutingReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchReport() {
      try {
        const data = await getReport(reportId);
        setReport(data);
      } catch (err) {
        console.error('Failed to fetch report:', err);
        // Use mock data for demonstration
        setReport(MOCK_REPORT);
        setError('Using demo data - backend connection failed');
      } finally {
        setLoading(false);
      }
    }

    fetchReport();
  }, [reportId]);

  const handlePrint = () => {
    window.print();
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-game-primary" />
      </div>
    );
  }

  if (!report) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold mb-2">Report Not Found</h2>
          <Button onClick={() => router.push('/lol')}>Go Back</Button>
        </div>
      </div>
    );
  }

  // Build common strategies from team analysis
  const commonStrategies: CommonStrategiesSection = {
    attackPatterns: [],
    defenseSetups: [],
    objectivePriorities: report.teamStrategy.lolMetrics ? [
      { text: `Prioritizes first Drake (${Math.round(report.teamStrategy.lolMetrics.firstDragonRate * 100)}% contest rate)`, metric: 'firstDragonRate', value: report.teamStrategy.lolMetrics.firstDragonRate, sampleSize: report.matchesAnalyzed, context: 'Early Game' },
      { text: `First tower at ~${report.teamStrategy.lolMetrics.firstTowerAvgTime?.toFixed(0) || '?'} mins (${Math.round(report.teamStrategy.lolMetrics.firstTowerRate * 100)}% first tower rate)`, metric: 'firstTowerRate', value: report.teamStrategy.lolMetrics.firstTowerRate, sampleSize: report.matchesAnalyzed, context: 'Early Game' },
    ] : [],
    timingPatterns: report.teamStrategy.lolMetrics ? [
      { text: `Average game duration: ${report.teamStrategy.lolMetrics.avgGameDuration?.toFixed(0) || '?'} minutes`, metric: 'avgGameDuration', value: report.teamStrategy.lolMetrics.avgGameDuration || 0, sampleSize: report.matchesAnalyzed, context: 'Game Pace' },
    ] : [],
  };

  return (
    <div className="min-h-screen theme-lol pb-16">
      {/* Header */}
      <div className="sticky top-16 z-40 bg-background/80 backdrop-blur-xl border-b border-border no-print">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Button 
              variant="ghost" 
              onClick={() => router.push('/lol')}
              className="gap-2"
            >
              <ArrowLeft className="w-4 h-4" />
              Back
            </Button>

            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={handlePrint}>
                <Printer className="w-4 h-4" />
                <span className="hidden sm:inline ml-2">Print</span>
              </Button>
              <Button variant="outline" size="sm">
                <Share2 className="w-4 h-4" />
                <span className="hidden sm:inline ml-2">Share</span>
              </Button>
              <Button variant="game" size="sm">
                <Download className="w-4 h-4" />
                <span className="hidden sm:inline ml-2">Export PDF</span>
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Error Banner */}
      {error && (
        <div className="container mx-auto px-4 pt-4">
          <div className="p-3 rounded-lg bg-yellow-500/10 border border-yellow-500/30 text-yellow-500 text-sm flex items-center gap-2">
            <AlertCircle className="w-4 h-4" />
            {error}
          </div>
        </div>
      )}

      {/* Report Content */}
      <div className="container mx-auto px-4 py-8 space-y-8">
        {/* Executive Summary */}
        <ExecutiveSummary
          summary={report.executiveSummary}
          teamAnalysis={report.teamStrategy}
          generatedAt={report.generatedAt}
          matchesAnalyzed={report.matchesAnalyzed}
          logoUrl={report.opponentTeam.logoUrl}
        />

        {/* Navigation Tabs for mobile */}
        <Tabs defaultValue="howtowin" className="md:hidden">
          <TabsList className="w-full">
            <TabsTrigger value="howtowin" className="flex-1">How to Win</TabsTrigger>
            <TabsTrigger value="strategy" className="flex-1">Strategy</TabsTrigger>
            <TabsTrigger value="players" className="flex-1">Players</TabsTrigger>
          </TabsList>
          <TabsContent value="howtowin">
            <HowToWin data={report.howToWin} game="lol" />
          </TabsContent>
          <TabsContent value="strategy">
            <CommonStrategies 
              strategies={commonStrategies}
              lolMetrics={report.teamStrategy.lolMetrics}
              game="lol"
            />
            <div className="mt-8">
              <Compositions data={report.compositions} game="lol" />
            </div>
          </TabsContent>
          <TabsContent value="players">
            <PlayerTendencies players={report.playerProfiles} game="lol" />
          </TabsContent>
        </Tabs>

        {/* Desktop Layout */}
        <div className="hidden md:block space-y-8">
          {/* How to Win - Full Width Hero */}
          <HowToWin data={report.howToWin} game="lol" />

          <Separator />

          {/* Trend Analysis - if available */}
          {(report as any).trendAnalysis && (
            <TrendAnalysis data={(report as any).trendAnalysis} game="lol" />
          )}

          {/* Common Strategies */}
          <CommonStrategies 
            strategies={commonStrategies}
            lolMetrics={report.teamStrategy.lolMetrics}
            game="lol"
          />

          {/* Player Tendencies */}
          <PlayerTendencies players={report.playerProfiles} game="lol" />

          {/* Compositions */}
          <Compositions data={report.compositions} game="lol" />
        </div>
      </div>
    </div>
  );
}
