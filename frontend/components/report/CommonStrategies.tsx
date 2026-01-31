'use client';

import { motion } from 'framer-motion';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Compass,
  Swords,
  Shield,
  Target,
  Clock,
  Flame
} from 'lucide-react';
import { cn, formatPercentage } from '@/lib/utils';
import type { CommonStrategiesSection, LoLTeamMetrics, VALTeamMetrics, EconomyRoundStats } from '@/types';

interface CommonStrategiesProps {
  strategies: CommonStrategiesSection;
  lolMetrics?: LoLTeamMetrics;
  valMetrics?: VALTeamMetrics;
  game: 'lol' | 'valorant';
}

export function CommonStrategies({ strategies, lolMetrics, valMetrics, game }: CommonStrategiesProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Compass className="w-5 h-5 text-game-primary" />
            Common Strategies
          </CardTitle>
        </CardHeader>
        <CardContent>
          {game === 'lol' && lolMetrics && (
            <LoLStrategies metrics={lolMetrics} strategies={strategies} />
          )}
          {game === 'valorant' && valMetrics && (
            <VALStrategies metrics={valMetrics} strategies={strategies} />
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
}

function LoLStrategies({ metrics, strategies }: { metrics: LoLTeamMetrics; strategies: CommonStrategiesSection }) {
  return (
    <div className="space-y-6">
      {/* Objective Control */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          <Target className="w-4 h-4" />
          Objective Control
        </h4>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <MetricCard 
            label="First Dragon" 
            value={metrics.firstDragonRate} 
            icon="🐉"
          />
          <MetricCard 
            label="First Herald" 
            value={metrics.heraldControlRate} 
            icon="👁️"
          />
          <MetricCard 
            label="First Tower" 
            value={metrics.firstTowerRate} 
            icon="🏰"
            subtitle={`~${metrics.firstTowerAvgTime?.toFixed(0) || '?'}min`}
          />
          <MetricCard 
            label="Baron Control" 
            value={metrics.baronControlRate} 
            icon="👑"
          />
        </div>
      </div>

      {/* Early Game */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          <Flame className="w-4 h-4" />
          Early Game Aggression
        </h4>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          <MetricCard 
            label="First Blood" 
            value={metrics.firstBloodRate} 
            icon="🩸"
          />
          <MetricCard 
            label="Gold Diff @15" 
            value={metrics.goldDiff15} 
            isGold
            icon="💰"
          />
          <MetricCard 
            label="Aggression" 
            value={metrics.aggressionScore / 100} 
            icon="⚔️"
          />
        </div>
      </div>

      {/* Timing */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          <Clock className="w-4 h-4" />
          Game Timing
        </h4>
        <div className="grid grid-cols-2 gap-4">
          <div className="p-4 rounded-lg bg-muted/30">
            <div className="text-2xl font-bold text-game-primary">
              {metrics.avgGameDuration?.toFixed(0) || '?'} min
            </div>
            <div className="text-sm text-muted-foreground">Avg Game Duration</div>
          </div>
          <div className="p-4 rounded-lg bg-muted/30">
            <div className="flex gap-2">
              {metrics.winConditions?.map((wc, i) => (
                <Badge key={i} variant="game">{wc}</Badge>
              )) || <span className="text-muted-foreground">N/A</span>}
            </div>
            <div className="text-sm text-muted-foreground mt-2">Win Conditions</div>
          </div>
        </div>
      </div>

      {/* Strategy Insights */}
      {strategies.objectivePriorities.length > 0 && (
        <div className="pt-4 border-t border-border">
          <h4 className="text-sm font-medium text-muted-foreground mb-3">Key Patterns</h4>
          <ul className="space-y-2">
            {strategies.objectivePriorities.slice(0, 4).map((insight, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <span className="text-game-primary">•</span>
                <span>{insight.text}</span>
                {insight.value > 0 && (
                  <Badge variant="outline" className="ml-auto">
                    {formatPercentage(insight.value)}
                  </Badge>
                )}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

function VALStrategies({ metrics, strategies }: { metrics: VALTeamMetrics; strategies: CommonStrategiesSection }) {
  return (
    <div className="space-y-6">
      {/* Team Performance Metrics */}
      {((metrics.aggressionScore ?? 0) > 0 || (metrics.clutchRate ?? 0) > 0 || (metrics.firstDeathRate ?? 0) > 0) && (
        <div>
          <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
            <Flame className="w-4 h-4 text-orange-400" />
            Team Performance Metrics
          </h4>
          <div className="grid grid-cols-3 gap-4">
            {(metrics.aggressionScore ?? 0) > 0 && (
              <div className="p-3 rounded-lg bg-orange-500/5 border border-orange-500/20 text-center">
                <div className="text-2xl font-bold text-orange-400">
                  {(metrics.aggressionScore ?? 0).toFixed(1)}
                </div>
                <div className="text-xs text-muted-foreground">Aggression Score</div>
              </div>
            )}
            {(metrics.clutchRate ?? 0) > 0 && (
              <div className="p-3 rounded-lg bg-purple-500/5 border border-purple-500/20 text-center">
                <div className="text-2xl font-bold text-purple-400">
                  {formatPercentage(metrics.clutchRate ?? 0)}
                </div>
                <div className="text-xs text-muted-foreground">Clutch Rate</div>
              </div>
            )}
            {(metrics.firstDeathRate ?? 0) > 0 && (
              <div className="p-3 rounded-lg bg-red-500/5 border border-red-500/20 text-center">
                <div className="text-2xl font-bold text-red-400">
                  {formatPercentage(metrics.firstDeathRate ?? 0)}
                </div>
                <div className="text-xs text-muted-foreground">First Death Rate</div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Side Performance */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          <Swords className="w-4 h-4 text-red-400" />
          <Shield className="w-4 h-4 text-blue-400" />
          Side Performance
        </h4>
        <div className="grid grid-cols-2 gap-4">
          <div className="p-4 rounded-lg bg-red-500/5 border border-red-500/20">
            <div className="flex items-center justify-between mb-2">
              <span className="text-red-400 font-medium">Attack</span>
              <span className="text-2xl font-bold text-red-400">
                {formatPercentage(metrics.attackWinRate)}
              </span>
            </div>
            <Progress 
              value={metrics.attackWinRate * 100} 
              className="h-2"
              indicatorClassName="bg-red-500"
            />
          </div>
          <div className="p-4 rounded-lg bg-blue-500/5 border border-blue-500/20">
            <div className="flex items-center justify-between mb-2">
              <span className="text-blue-400 font-medium">Defense</span>
              <span className="text-2xl font-bold text-blue-400">
                {formatPercentage(metrics.defenseWinRate)}
              </span>
            </div>
            <Progress 
              value={metrics.defenseWinRate * 100} 
              className="h-2"
              indicatorClassName="bg-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Economy */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          💰 Economy Performance
        </h4>
        <div className="grid grid-cols-3 gap-4">
          <EconomyMetricCard 
            label="Eco Rounds" 
            winRate={metrics.ecoRoundWinRate}
            stats={metrics.economyStats}
            type="eco"
            icon="📉"
          />
          <EconomyMetricCard 
            label="Force Buys" 
            winRate={metrics.forceBuyWinRate}
            stats={metrics.economyStats}
            type="force"
            icon="⚡"
          />
          <EconomyMetricCard 
            label="Full Buys" 
            winRate={metrics.fullBuyWinRate}
            stats={metrics.economyStats}
            type="fullBuy"
            icon="💪"
          />
        </div>
        {/* Average Team Loadout */}
        {(metrics.avgTeamLoadout || metrics.economyStats?.avgLoadoutValue) && (
          <div className="mt-3 p-3 rounded-lg bg-yellow-500/5 border border-yellow-500/20 text-center">
            <div className="text-sm text-muted-foreground">Average Team Loadout</div>
            <div className="text-2xl font-bold text-yellow-500">
              {Math.round(metrics.avgTeamLoadout || metrics.economyStats?.avgLoadoutValue || 0).toLocaleString()}
              <span className="text-sm font-normal text-muted-foreground ml-1">credits</span>
            </div>
          </div>
        )}
      </div>

      {/* Pistol Rounds */}
      <div>
        <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
          🔫 Pistol Rounds
        </h4>
        <div className="grid grid-cols-3 gap-4">
          <MetricCard 
            label="Overall" 
            value={metrics.pistolWinRate} 
            icon="🎯"
          />
          <MetricCard 
            label="Attack Pistol" 
            value={metrics.attackPistolWinRate} 
            icon="⚔️"
          />
          <MetricCard 
            label="Defense Pistol" 
            value={metrics.defensePistolWinRate} 
            icon="🛡️"
          />
        </div>
      </div>

      {/* Map Pool */}
      {metrics.mapPool && metrics.mapPool.length > 0 && (
        <div>
          <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-4">
            🗺️ Map Pool
          </h4>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
            {metrics.mapPool.slice(0, 6).map((map) => (
              <div 
                key={map.map}
                className={cn(
                  "p-3 rounded-lg border",
                  map.winRate >= 0.5 
                    ? "border-green-500/30 bg-green-500/5" 
                    : "border-red-500/30 bg-red-500/5"
                )}
              >
                <div className="font-medium">{map.map}</div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">{map.games} games</span>
                  <span className={map.winRate >= 0.5 ? "text-green-500" : "text-red-500"}>
                    {formatPercentage(map.winRate)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Attack Patterns */}
      {strategies.attackPatterns.length > 0 && (
        <div className="pt-4 border-t border-border">
          <h4 className="text-sm font-medium text-red-400 mb-3">Attack Patterns</h4>
          <ul className="space-y-2">
            {strategies.attackPatterns.slice(0, 3).map((insight, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <span className="text-red-400">•</span>
                <span>{insight.text}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Defense Setups */}
      {strategies.defenseSetups.length > 0 && (
        <div className="pt-4 border-t border-border">
          <h4 className="text-sm font-medium text-blue-400 mb-3">Defense Setups</h4>
          <ul className="space-y-2">
            {strategies.defenseSetups.slice(0, 3).map((insight, i) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <span className="text-blue-400">•</span>
                <span>{insight.text}</span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

function MetricCard({ 
  label, 
  value, 
  icon, 
  subtitle,
  isGold = false 
}: { 
  label: string; 
  value: number; 
  icon: string;
  subtitle?: string;
  isGold?: boolean;
}) {
  const displayValue = isGold 
    ? (value > 0 ? `+${value}` : value.toString())
    : formatPercentage(value);
  
  const valueColor = isGold
    ? (value > 0 ? 'text-green-500' : value < 0 ? 'text-red-500' : 'text-muted-foreground')
    : (value >= 0.5 ? 'text-green-500' : value >= 0.4 ? 'text-yellow-500' : 'text-red-500');

  return (
    <div className="p-3 rounded-lg bg-muted/30 text-center">
      <span className="text-2xl">{icon}</span>
      <div className={cn("text-xl font-bold mt-1", valueColor)}>
        {displayValue}
      </div>
      <div className="text-xs text-muted-foreground">{label}</div>
      {subtitle && (
        <div className="text-xs text-muted-foreground/70">{subtitle}</div>
      )}
    </div>
  );
}

function EconomyMetricCard({ 
  label, 
  winRate,
  stats,
  type,
  icon 
}: { 
  label: string; 
  winRate: number;
  stats?: EconomyRoundStats;
  type: 'eco' | 'force' | 'fullBuy';
  icon: string;
}) {
  // Get detailed stats if available
  let rounds = 0;
  let wins = 0;
  
  if (stats) {
    if (type === 'eco') {
      rounds = typeof stats.ecoRounds === 'number' ? stats.ecoRounds : stats.ecoRounds?.total || 0;
      wins = stats.ecoWins || (typeof stats.ecoRounds === 'object' ? stats.ecoRounds?.won : 0) || 0;
    } else if (type === 'force') {
      rounds = typeof stats.forceRounds === 'number' ? stats.forceRounds : stats.forceRounds?.total || 0;
      wins = stats.forceWins || (typeof stats.forceRounds === 'object' ? stats.forceRounds?.won : 0) || 0;
    } else if (type === 'fullBuy') {
      rounds = typeof stats.fullBuyRounds === 'number' ? stats.fullBuyRounds : stats.fullBuyRounds?.total || 0;
      wins = stats.fullBuyWins || (typeof stats.fullBuyRounds === 'object' ? stats.fullBuyRounds?.won : 0) || 0;
    }
  }
  
  const valueColor = winRate >= 0.5 ? 'text-green-500' : winRate >= 0.4 ? 'text-yellow-500' : 'text-red-500';

  return (
    <div className="p-3 rounded-lg bg-muted/30 text-center">
      <span className="text-2xl">{icon}</span>
      <div className={cn("text-xl font-bold mt-1", valueColor)}>
        {formatPercentage(winRate)}
      </div>
      <div className="text-xs text-muted-foreground">{label}</div>
      {rounds > 0 && (
        <div className="text-xs text-muted-foreground/70 mt-0.5">
          {wins}W / {rounds} rounds
        </div>
      )}
    </div>
  );
}
