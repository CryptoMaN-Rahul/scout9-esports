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
import { MapPerformance } from '@/components/report/MapPerformance';
import { getReport } from '@/lib/api';
import { cn } from '@/lib/utils';
import type { ScoutingReport, CommonStrategiesSection } from '@/types';

// Mock VALORANT data for demonstration
const MOCK_REPORT: ScoutingReport = {
  id: 'mock-val-1',
  generatedAt: new Date().toISOString(),
  opponentTeam: { id: '1079', name: 'Sentinels', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d05bdadd653b68669ad96d91d7f86cff' },
  title: 'valorant',
  matchesAnalyzed: 8,
  executiveSummary: 'Sentinels is an aggressive duelist-heavy team that excels on Attack side (62% round win rate) but struggles on Defense (48%). They heavily rely on TenZ\'s Jett for opening picks - when he gets first blood, they win 78% of rounds. Counter by playing anti-flash setups and denying early aggression. Their economy management is poor - they often force-buy and lose anti-eco rounds. Play patient on Defense, stack their A-site preference, and punish their economic mistakes.',
  howToWin: {
    winCondition: 'Deny TenZ first blood opportunities and exploit their weak Defense side. Stack A-site and trade aggressively when they entry.',
    confidenceScore: 72,
    actionableInsights: [
      {
        recommendation: 'Stack A-site with utility',
        dataBacking: 'Sentinels attacks A-site 65% of the time on Bind',
        impact: 'HIGH',
        actionType: 'STRATEGY',
        confidence: 88,
      },
      {
        recommendation: 'Ban Jett or force TenZ onto Chamber',
        dataBacking: 'TenZ has 1.8 KD on Jett vs 1.1 on other duelists',
        impact: 'HIGH',
        actionType: 'BAN',
        confidence: 85,
      },
      {
        recommendation: 'Force their Defense pistol round',
        dataBacking: 'They lose 60% of anti-eco rounds due to force-buying',
        impact: 'MEDIUM',
        actionType: 'STRATEGY',
        confidence: 75,
      },
      {
        recommendation: 'Play retake on B-site',
        dataBacking: 'Their B executes are weak (45% success rate)',
        impact: 'MEDIUM',
        actionType: 'STRATEGY',
        confidence: 70,
      },
    ],
    draftStrategy: {
      priorityBans: [
        { character: 'Jett', reason: 'TenZ signature - 1.8 KD, 28% first blood rate', playerName: 'TenZ', winRate: 0.65 },
        { character: 'Astra', reason: 'Zellsis has 70% WR with team coordination', playerName: 'Zellsis', winRate: 0.70 },
      ],
      targetPicks: [
        { character: 'Chamber', reason: 'Force TenZ onto uncomfortable agents', playerName: 'TenZ' },
      ],
      recommendedPicks: [
        { character: 'Cypher', reason: 'Information advantage counters fast executes', winRate: 0.58 },
        { character: 'Killjoy', reason: 'Lockdown punishes their aggressive Attack', winRate: 0.55 },
        { character: 'Omen', reason: 'Counter-flash and reposition against entries', winRate: 0.52 },
      ],
    },
    inGameStrategy: [
      { phase: 'Pistol Rounds', timing: 'Rd 1 & 13', strategy: 'Play aggressive - their pistol win rate is only 42%', priority: 'HIGH' },
      { phase: 'Attack Side', timing: 'Rd 1-12', strategy: 'Default and look for picks before executing', priority: 'MEDIUM' },
      { phase: 'Defense Side', timing: 'Rd 13-24', strategy: 'Stack A-site, 3-1-1 setup, trade aggressively', priority: 'HIGH' },
      { phase: 'Economy', timing: 'Post-loss', strategy: 'Full save after losses - they force buy and lose bonus', priority: 'HIGH' },
    ],
  },
  teamStrategy: {
    teamId: '2',
    teamName: 'Sentinels',
    title: 'valorant',
    matchesAnalyzed: 8,
    gamesAnalyzed: 24,
    winRate: 0.58,
    recentForm: 'stable',
    strengths: [
      { text: 'Strong Attack side execution (62% round win rate)', importance: 'HIGH' },
      { text: 'TenZ first blood rate of 28% creates early advantages', importance: 'HIGH' },
      { text: 'Excellent clutch rate in 1v1 situations (65%)', importance: 'MEDIUM' },
    ],
    weaknesses: [
      { text: 'Defense side struggles (48% round win rate)', importance: 'HIGH' },
      { text: 'Poor economy management - frequent force buys', importance: 'HIGH' },
      { text: 'Over-reliance on TenZ for opening picks', importance: 'MEDIUM' },
    ],
    valMetrics: {
      attackWinRate: 0.62,
      defenseWinRate: 0.48,
      clutchRate: 0.38,
      pistolWinRate: 0.42,
      attackPistolWinRate: 0.45,
      defensePistolWinRate: 0.39,
      firstBloodRate: 0.55,
      ecoRoundWinRate: 0.35,
      forceBuyWinRate: 0.40,
      fullBuyWinRate: 0.65,
      aggressionScore: 70,
      mapStats: {},
      mapPool: [],
    },
  },
  playerProfiles: [
    {
      playerId: '1',
      nickname: 'TenZ',
      role: 'Duelist',
      teamId: '2',
      gamesPlayed: 24,
      kda: 1.45,
      avgKills: 18.5,
      avgDeaths: 12.8,
      avgAssists: 3.2,
      characterPool: [
        { name: 'Jett', games: 18, wins: 12, losses: 6, winRate: 0.67, avgKDA: 1.8 },
        { name: 'Raze', games: 4, wins: 2, losses: 2, winRate: 0.50, avgKDA: 1.2 },
        { name: 'Chamber', games: 2, wins: 0, losses: 2, winRate: 0.00, avgKDA: 0.9 },
      ],
      signaturePicks: ['Jett', 'Raze'],
      threatLevel: 9,
      threatReason: 'Primary star - creates space and gets opening kills',
      weaknesses: [{ text: 'Performance drops significantly without Jett', importance: 'HIGH' }],
      tendencies: ['Always dashes forward on Attack', 'Updrafts predictably in post-plant'],
    },
    {
      playerId: '2',
      nickname: 'SicK',
      role: 'Flex',
      teamId: '2',
      gamesPlayed: 24,
      kda: 1.15,
      avgKills: 14.2,
      avgDeaths: 12.4,
      avgAssists: 5.8,
      characterPool: [
        { name: 'Sage', games: 12, wins: 7, losses: 5, winRate: 0.58, avgKDA: 1.2 },
        { name: 'Skye', games: 8, wins: 4, losses: 4, winRate: 0.50, avgKDA: 1.1 },
      ],
      signaturePicks: ['Sage', 'Skye'],
      threatLevel: 6,
      threatReason: 'Consistent support player, enables TenZ',
      weaknesses: [{ text: 'Can be exploited when isolated', importance: 'MEDIUM' }],
      tendencies: ['Walls mid on Attack', 'Plays passive angles'],
    },
    {
      playerId: '3',
      nickname: 'Zellsis',
      role: 'Controller',
      teamId: '2',
      gamesPlayed: 24,
      kda: 1.08,
      avgKills: 12.5,
      avgDeaths: 11.6,
      avgAssists: 7.2,
      characterPool: [
        { name: 'Astra', games: 14, wins: 10, losses: 4, winRate: 0.71, avgKDA: 1.1 },
        { name: 'Omen', games: 10, wins: 4, losses: 6, winRate: 0.40, avgKDA: 1.0 },
      ],
      signaturePicks: ['Astra'],
      threatLevel: 7,
      threatReason: 'Excellent utility usage, sets up team plays',
      weaknesses: [{ text: 'Less impactful on Omen', importance: 'MEDIUM' }],
      tendencies: ['Stars B-site early', 'Smokes predictable timings'],
    },
    {
      playerId: '4',
      nickname: 'Dapr',
      role: 'Sentinel',
      teamId: '2',
      gamesPlayed: 24,
      kda: 1.22,
      avgKills: 13.8,
      avgDeaths: 11.3,
      avgAssists: 4.5,
      characterPool: [
        { name: 'Cypher', games: 15, wins: 9, losses: 6, winRate: 0.60, avgKDA: 1.3 },
        { name: 'Killjoy', games: 9, wins: 5, losses: 4, winRate: 0.56, avgKDA: 1.1 },
      ],
      signaturePicks: ['Cypher', 'Killjoy'],
      threatLevel: 7,
      threatReason: 'Strong anchor, holds sites solo',
      weaknesses: [{ text: 'Predictable trap placements', importance: 'MEDIUM' }],
      tendencies: ['Setups A-site', 'Lurks on Attack'],
    },
    {
      playerId: '5',
      nickname: 'Zombs',
      role: 'Controller',
      teamId: '2',
      gamesPlayed: 24,
      kda: 0.95,
      avgKills: 10.2,
      avgDeaths: 10.8,
      avgAssists: 6.5,
      characterPool: [
        { name: 'Viper', games: 20, wins: 12, losses: 8, winRate: 0.60, avgKDA: 1.0 },
        { name: 'Brimstone', games: 4, wins: 2, losses: 2, winRate: 0.50, avgKDA: 0.8 },
      ],
      signaturePicks: ['Viper'],
      threatLevel: 5,
      threatReason: 'Consistent utility, rarely loses fights',
      weaknesses: [{ text: 'Low fragging impact', importance: 'MEDIUM' }],
      tendencies: ['Plays lineups', 'Passive positioning'],
    },
  ],
  compositions: {
    topCompositions: [
      { characters: ['Jett', 'Sage', 'Astra', 'Cypher', 'Viper'], frequency: 0.45, winRate: 0.65, gamesPlayed: 11, archetype: 'Default' },
      { characters: ['Raze', 'Skye', 'Omen', 'Killjoy', 'Viper'], frequency: 0.25, winRate: 0.50, gamesPlayed: 6, archetype: 'Double Controller' },
      { characters: ['Jett', 'Sage', 'Omen', 'Cypher', 'Sova'], frequency: 0.20, winRate: 0.60, gamesPlayed: 5, archetype: 'Info Heavy' },
    ],
    firstPickPriorities: [
      { character: 'Jett', rate: 0.7, winRate: 0.65, gamesPlayed: 14 },
      { character: 'Viper', rate: 0.6, winRate: 0.60, gamesPlayed: 12 },
      { character: 'Astra', rate: 0.5, winRate: 0.55, gamesPlayed: 10 },
      { character: 'Cypher', rate: 0.45, winRate: 0.58, gamesPlayed: 9 },
    ],
    commonBans: ['Chamber', 'Neon', 'Harbor'],
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
      <div className="min-h-screen flex items-center justify-center theme-valorant">
        <Loader2 className="w-8 h-8 animate-spin text-game-primary" />
      </div>
    );
  }

  if (!report) {
    return (
      <div className="min-h-screen flex items-center justify-center theme-valorant">
        <div className="text-center">
          <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold mb-2">Report Not Found</h2>
          <Button onClick={() => router.push('/')}>Go Back</Button>
        </div>
      </div>
    );
  }

  // Build common strategies from team analysis
  const commonStrategies: CommonStrategiesSection = {
    attackPatterns: [
      { text: 'Fast A-site executes with Jett entry', metric: 'attackSuccess', value: 0.65, sampleSize: 50, context: 'Attack' },
    ],
    defenseSetups: [
      { text: '3-1-1 setup with Cypher anchor A', metric: 'defenseHold', value: 0.48, sampleSize: 45, context: 'Defense' },
    ],
    objectivePriorities: [],
    timingPatterns: [],
  };

  return (
    <div className="min-h-screen theme-valorant pb-16">
      {/* Header */}
      <div className="sticky top-16 z-40 bg-background/80 backdrop-blur-xl border-b border-border no-print">
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
            <HowToWin data={report.howToWin} game="valorant" />
          </TabsContent>
          <TabsContent value="strategy">
            <CommonStrategies 
              strategies={commonStrategies}
              valMetrics={report.teamStrategy.valMetrics}
              game="valorant"
            />
            <div className="mt-8">
              <Compositions data={report.compositions} game="valorant" />
            </div>
          </TabsContent>
          <TabsContent value="players">
            <PlayerTendencies players={report.playerProfiles} game="valorant" />
          </TabsContent>
        </Tabs>

        {/* Desktop Layout */}
        <div className="hidden md:block space-y-8">
          {/* How to Win - Full Width Hero */}
          <HowToWin data={report.howToWin} game="valorant" />

          <Separator />

          {/* Two Column Layout for Trend + Map Performance */}
          <div className="grid lg:grid-cols-2 gap-8">
            {/* Trend Analysis */}
            {(report as any).trendAnalysis && (
              <TrendAnalysis data={(report as any).trendAnalysis} game="valorant" />
            )}
            
            {/* Map Performance */}
            {report.teamStrategy?.valMetrics && (
              <MapPerformance 
                mapStats={report.teamStrategy.valMetrics.mapStats}
                mapPool={report.teamStrategy.valMetrics.mapPool}
              />
            )}
          </div>

          {/* Common Strategies */}
          <CommonStrategies 
            strategies={commonStrategies}
            valMetrics={report.teamStrategy.valMetrics}
            game="valorant"
          />

          {/* Player Tendencies */}
          <PlayerTendencies players={report.playerProfiles} game="valorant" />

          {/* Compositions */}
          <Compositions data={report.compositions} game="valorant" />
        </div>
      </div>
    </div>
  );
}
