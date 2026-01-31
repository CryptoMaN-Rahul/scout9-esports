'use client';

import { motion } from 'framer-motion';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  TrendingUp, 
  TrendingDown, 
  Minus,
  Calendar,
  BarChart3
} from 'lucide-react';
import { cn, formatPercentage, formatDateTime } from '@/lib/utils';
import type { TeamAnalysis } from '@/types';

interface ExecutiveSummaryProps {
  summary: string;
  teamAnalysis: TeamAnalysis;
  generatedAt: string;
  matchesAnalyzed: number;
  logoUrl?: string;
}

export function ExecutiveSummary({ 
  summary, 
  teamAnalysis, 
  generatedAt, 
  matchesAnalyzed,
  logoUrl 
}: ExecutiveSummaryProps) {
  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'improving':
        return <TrendingUp className="w-4 h-4 text-green-500" />;
      case 'declining':
        return <TrendingDown className="w-4 h-4 text-red-500" />;
      default:
        return <Minus className="w-4 h-4 text-yellow-500" />;
    }
  };

  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'improving':
        return 'text-green-500';
      case 'declining':
        return 'text-red-500';
      default:
        return 'text-yellow-500';
    }
  };

  const winRateColor = teamAnalysis.winRate >= 0.5 
    ? 'text-green-500' 
    : 'text-red-500';

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
    >
      <Card glass>
        <CardContent className="p-6">
          {/* Header */}
          <div className="flex flex-wrap items-center gap-4 mb-6">
            {/* Team Logo/Name */}
            <div className="flex items-center gap-4">
              {logoUrl ? (
                <img 
                  src={logoUrl} 
                  alt={teamAnalysis.teamName} 
                  className="w-16 h-16 rounded-xl object-contain bg-muted/50 p-1"
                />
              ) : (
                <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-game-primary/20 to-game-accent/20 flex items-center justify-center text-2xl font-bold">
                  {teamAnalysis.teamName.substring(0, 2).toUpperCase()}
                </div>
              )}
              <div>
                <h1 className="text-2xl font-bold">{teamAnalysis.teamName}</h1>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Calendar className="w-4 h-4" />
                  {formatDateTime(generatedAt)}
                </div>
              </div>
            </div>

            {/* Quick Stats */}
            <div className="flex-1 flex flex-wrap justify-end gap-4">
              {/* Win Rate */}
              <div className="text-center px-4 py-2 rounded-lg bg-muted/50">
                <div className={cn("text-2xl font-bold", winRateColor)}>
                  {formatPercentage(teamAnalysis.winRate)}
                </div>
                <div className="text-xs text-muted-foreground">Win Rate</div>
              </div>

              {/* Matches Analyzed */}
              <div className="text-center px-4 py-2 rounded-lg bg-muted/50">
                <div className="text-2xl font-bold text-game-primary">
                  {matchesAnalyzed}
                </div>
                <div className="text-xs text-muted-foreground">Matches</div>
              </div>

              {/* Form */}
              <div className="text-center px-4 py-2 rounded-lg bg-muted/50">
                <div className={cn("text-2xl font-bold flex items-center gap-1", getTrendColor(teamAnalysis.recentForm || 'stable'))}>
                  {getTrendIcon(teamAnalysis.recentForm || 'stable')}
                  {(teamAnalysis.recentForm || 'stable').charAt(0).toUpperCase() + (teamAnalysis.recentForm || 'stable').slice(1)}
                </div>
                <div className="text-xs text-muted-foreground">Form</div>
              </div>
            </div>
          </div>

          {/* Executive Summary Text */}
          <div className="mb-6">
            <h3 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-2">
              <BarChart3 className="w-4 h-4" />
              Executive Summary
            </h3>
            <p className="text-lg leading-relaxed">
              {summary}
            </p>
          </div>

          {/* Strengths & Weaknesses */}
          <div className="grid md:grid-cols-2 gap-4">
            {/* Strengths */}
            <div className="p-4 rounded-lg bg-green-500/5 border border-green-500/20">
              <h4 className="text-sm font-medium text-green-400 mb-3">Key Strengths</h4>
              <ul className="space-y-2">
                {teamAnalysis.strengths.slice(0, 3).map((strength, i) => (
                  <li key={i} className="flex items-start gap-2 text-sm">
                    <span className="text-green-500 mt-0.5">✓</span>
                    <span>{strength.title || strength.text}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* Weaknesses */}
            <div className="p-4 rounded-lg bg-red-500/5 border border-red-500/20">
              <h4 className="text-sm font-medium text-red-400 mb-3">Exploitable Weaknesses</h4>
              <ul className="space-y-2">
                {teamAnalysis.weaknesses.slice(0, 3).map((weakness, i) => (
                  <li key={i} className="flex items-start gap-2 text-sm">
                    <span className="text-red-500 mt-0.5">×</span>
                    <span>{weakness.title || weakness.text}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Game-Specific Metrics Preview */}
          {teamAnalysis.lolMetrics && (
            <div className="mt-6 pt-6 border-t border-border">
              <h4 className="text-sm font-medium text-muted-foreground mb-4">Playstyle Profile</h4>
              <div className="grid grid-cols-3 gap-4">
                <PlaystyleBar 
                  label="Early Game" 
                  value={teamAnalysis.lolMetrics.earlyGameRating} 
                />
                <PlaystyleBar 
                  label="Mid Game" 
                  value={teamAnalysis.lolMetrics.midGameRating} 
                />
                <PlaystyleBar 
                  label="Late Game" 
                  value={teamAnalysis.lolMetrics.lateGameRating} 
                />
              </div>
            </div>
          )}

          {teamAnalysis.valMetrics && (
            <div className="mt-6 pt-6 border-t border-border">
              <h4 className="text-sm font-medium text-muted-foreground mb-4">Side Performance</h4>
              <div className="grid grid-cols-2 gap-4">
                <SidePerformanceBar 
                  label="Attack" 
                  value={teamAnalysis.valMetrics.attackWinRate} 
                  color="text-red-400"
                />
                <SidePerformanceBar 
                  label="Defense" 
                  value={teamAnalysis.valMetrics.defenseWinRate} 
                  color="text-blue-400"
                />
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
}

function PlaystyleBar({ label, value }: { label: string; value: number }) {
  const getColor = (val: number) => {
    if (val >= 70) return 'bg-green-500';
    if (val >= 50) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  return (
    <div>
      <div className="flex justify-between text-sm mb-1">
        <span className="text-muted-foreground">{label}</span>
        <span className="font-medium">{value}</span>
      </div>
      <Progress 
        value={value} 
        className="h-2" 
        indicatorClassName={getColor(value)}
      />
    </div>
  );
}

function SidePerformanceBar({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div className="flex items-center gap-3">
      <span className={cn("text-sm font-medium w-16", color)}>{label}</span>
      <div className="flex-1">
        <Progress 
          value={value * 100} 
          className="h-3" 
          indicatorClassName={color.replace('text-', 'bg-')}
        />
      </div>
      <span className="text-sm font-bold w-12 text-right">
        {formatPercentage(value)}
      </span>
    </div>
  );
}
