'use client';

import { motion } from 'framer-motion';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { 
  TrendingUp, 
  TrendingDown, 
  Minus, 
  Flame, 
  Snowflake, 
  Activity,
  Target,
  BarChart3
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { TrendAnalysis as TrendAnalysisType } from '@/types';

interface TrendAnalysisProps {
  data?: {
    teamId?: string;
    last5WinRate?: number;
    last10WinRate?: number;
    overallWinRate?: number;
    formIndicator?: string; // "hot", "stable", "cold"
    formScore?: number; // -100 to 100
    metricTrends?: Array<{
      metric: string;
      direction: string;
      values: number[];
    }>;
    trendChanges?: Array<{
      metric: string;
      oldValue: number;
      newValue: number;
      changeDate: string;
      explanation: string;
    }>;
  };
  game?: 'lol' | 'valorant';
}

const FormIcon = ({ form }: { form: string }) => {
  switch (form) {
    case 'hot':
      return <Flame className="w-6 h-6 text-orange-500" />;
    case 'cold':
      return <Snowflake className="w-6 h-6 text-blue-400" />;
    default:
      return <Activity className="w-6 h-6 text-yellow-500" />;
  }
};

const FormBadge = ({ form, score }: { form: string; score: number }) => {
  const colors = {
    hot: 'bg-orange-500/10 border-orange-500/30 text-orange-500',
    stable: 'bg-yellow-500/10 border-yellow-500/30 text-yellow-500',
    cold: 'bg-blue-400/10 border-blue-400/30 text-blue-400',
  };

  const labels = {
    hot: 'On Fire',
    stable: 'Stable',
    cold: 'Cold Streak',
  };

  return (
    <div className={cn(
      'flex items-center gap-2 px-4 py-2 rounded-lg border',
      colors[form as keyof typeof colors] || colors.stable
    )}>
      <FormIcon form={form} />
      <div>
        <div className="font-semibold">{labels[form as keyof typeof labels] || 'Stable'}</div>
        <div className="text-xs opacity-75">
          Momentum: {score > 0 ? '+' : ''}{score.toFixed(0)}
        </div>
      </div>
    </div>
  );
};

const WinRateBar = ({ label, value, isRecent = false }: { label: string; value: number; isRecent?: boolean }) => {
  const percentage = value * 100;
  const color = percentage >= 60 ? 'bg-green-500' : percentage >= 45 ? 'bg-yellow-500' : 'bg-red-500';
  
  return (
    <div className="space-y-1">
      <div className="flex justify-between text-sm">
        <span className={cn(isRecent && 'font-medium')}>{label}</span>
        <span className={cn('font-mono', isRecent && 'font-bold')}>{percentage.toFixed(0)}%</span>
      </div>
      <div className="h-2 bg-muted rounded-full overflow-hidden">
        <motion.div 
          className={cn('h-full rounded-full', color)}
          initial={{ width: 0 }}
          animate={{ width: `${percentage}%` }}
          transition={{ duration: 0.5, delay: 0.1 }}
        />
      </div>
    </div>
  );
};

const TrendDirection = ({ direction }: { direction: string }) => {
  switch (direction) {
    case 'improving':
      return (
        <div className="flex items-center gap-1 text-green-500 text-sm">
          <TrendingUp className="w-4 h-4" />
          <span>Improving</span>
        </div>
      );
    case 'declining':
      return (
        <div className="flex items-center gap-1 text-red-500 text-sm">
          <TrendingDown className="w-4 h-4" />
          <span>Declining</span>
        </div>
      );
    default:
      return (
        <div className="flex items-center gap-1 text-muted-foreground text-sm">
          <Minus className="w-4 h-4" />
          <span>Stable</span>
        </div>
      );
  }
};

const MiniSparkline = ({ values }: { values: number[] }) => {
  if (!values || values.length < 2) return null;
  
  // Take last 10 values
  const displayValues = values.slice(-10);
  const max = Math.max(...displayValues);
  const min = Math.min(...displayValues);
  const range = max - min || 1;
  
  const points = displayValues.map((v, i) => {
    const x = (i / (displayValues.length - 1)) * 100;
    const y = 100 - ((v - min) / range) * 100;
    return `${x},${y}`;
  }).join(' ');
  
  const isImproving = displayValues[displayValues.length - 1] > displayValues[0];
  
  return (
    <svg viewBox="0 0 100 100" className="w-20 h-8" preserveAspectRatio="none">
      <polyline
        points={points}
        fill="none"
        stroke={isImproving ? '#22c55e' : '#ef4444'}
        strokeWidth="3"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export function TrendAnalysis({ data, game = 'lol' }: TrendAnalysisProps) {
  if (!data) return null;

  const {
    last5WinRate = 0,
    last10WinRate = 0,
    overallWinRate = 0,
    formIndicator = 'stable',
    formScore = 0,
    metricTrends = [],
    trendChanges = [],
  } = data;

  const winRateTrend = metricTrends.find(t => t.metric === 'Win Rate');

  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2">
          <BarChart3 className="w-5 h-5 text-primary" />
          Trend Analysis
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Form Status */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          <FormBadge form={formIndicator} score={formScore} />
          
          {/* Mini Sparkline */}
          {winRateTrend && winRateTrend.values.length > 2 && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Recent Form:</span>
              <MiniSparkline values={winRateTrend.values} />
            </div>
          )}
        </div>

        {/* Win Rate Comparison */}
        <div className="space-y-4">
          <h4 className="text-sm font-medium text-muted-foreground">Win Rate Breakdown</h4>
          <div className="space-y-3">
            <WinRateBar label="Last 5 Matches" value={last5WinRate} isRecent />
            <WinRateBar label="Last 10 Matches" value={last10WinRate} />
            <WinRateBar label="Overall" value={overallWinRate} />
          </div>
        </div>

        {/* Metric Trends */}
        {metricTrends.length > 0 && (
          <div className="space-y-3">
            <h4 className="text-sm font-medium text-muted-foreground">Metric Trends</h4>
            <div className="grid gap-3">
              {metricTrends.slice(0, 3).map((trend, idx) => (
                <div 
                  key={idx}
                  className="flex items-center justify-between p-3 bg-muted/30 rounded-lg"
                >
                  <div className="flex items-center gap-2">
                    <Target className="w-4 h-4 text-muted-foreground" />
                    <span className="text-sm font-medium">{trend.metric}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    {trend.values.length > 2 && (
                      <MiniSparkline values={trend.values} />
                    )}
                    <TrendDirection direction={trend.direction} />
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Trend Changes */}
        {trendChanges.length > 0 && (
          <div className="space-y-3">
            <h4 className="text-sm font-medium text-muted-foreground">Recent Changes</h4>
            <div className="space-y-2">
              {trendChanges.slice(0, 3).map((change, idx) => {
                const isPositive = change.newValue > change.oldValue;
                return (
                  <div 
                    key={idx}
                    className="p-3 bg-muted/30 rounded-lg text-sm"
                  >
                    <div className="flex items-center gap-2">
                      {isPositive ? (
                        <TrendingUp className="w-4 h-4 text-green-500" />
                      ) : (
                        <TrendingDown className="w-4 h-4 text-red-500" />
                      )}
                      <span className="font-medium">{change.metric}</span>
                      <span className={cn(
                        'text-xs px-2 py-0.5 rounded',
                        isPositive ? 'bg-green-500/10 text-green-500' : 'bg-red-500/10 text-red-500'
                      )}>
                        {(change.oldValue * 100).toFixed(0)}% → {(change.newValue * 100).toFixed(0)}%
                      </span>
                    </div>
                    {change.explanation && (
                      <p className="text-muted-foreground mt-1 ml-6">{change.explanation}</p>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
