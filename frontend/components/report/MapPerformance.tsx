'use client';

import { motion } from 'framer-motion';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { 
  Map, 
  Swords, 
  Shield, 
  TrendingUp,
  TrendingDown,
  Minus
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { MapPoolEntry as TypedMapPoolEntry } from '@/types';

interface MapStat {
  mapName: string;
  gamesPlayed: number;
  winRate: number;
  attackWinRate: number;
  defenseWinRate: number;
  avgRoundsWon?: number;
  avgRoundsLost?: number;
}

interface MapPerformanceProps {
  mapStats?: Record<string, MapStat>;
  mapPool?: TypedMapPoolEntry[];
  economyByMap?: Record<string, {
    mapName: string;
    ecoWinRate: number;
    ecoRounds: number;
    forceWinRate: number;
    forceRounds: number;
    fullBuyWinRate: number;
    fullBuyRounds: number;
  }>;
}

const mapDisplayNames: Record<string, string> = {
  ascent: 'Ascent',
  bind: 'Bind',
  breeze: 'Breeze',
  fracture: 'Fracture',
  haven: 'Haven',
  icebox: 'Icebox',
  lotus: 'Lotus',
  pearl: 'Pearl',
  split: 'Split',
  sunset: 'Sunset',
  corrode: 'Abyss',
};

const StrengthBadge = ({ strength }: { strength: string }) => {
  const colors = {
    strong: 'bg-green-500/10 text-green-500 border-green-500/30',
    average: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/30',
    weak: 'bg-red-500/10 text-red-500 border-red-500/30',
  };

  const icons = {
    strong: <TrendingUp className="w-3 h-3" />,
    average: <Minus className="w-3 h-3" />,
    weak: <TrendingDown className="w-3 h-3" />,
  };

  return (
    <span className={cn(
      'inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded border',
      colors[strength as keyof typeof colors] || colors.average
    )}>
      {icons[strength as keyof typeof icons]}
      {strength.charAt(0).toUpperCase() + strength.slice(1)}
    </span>
  );
};

const WinRateCircle = ({ rate, label, size = 'sm' }: { rate: number; label: string; size?: 'sm' | 'md' }) => {
  const percentage = rate * 100;
  const color = percentage >= 60 ? 'text-green-500' : percentage >= 45 ? 'text-yellow-500' : 'text-red-500';
  const bgColor = percentage >= 60 ? 'stroke-green-500' : percentage >= 45 ? 'stroke-yellow-500' : 'stroke-red-500';
  
  const sizeClasses = size === 'md' ? 'w-16 h-16' : 'w-12 h-12';
  const textSize = size === 'md' ? 'text-base' : 'text-xs';
  
  return (
    <div className="flex flex-col items-center gap-1">
      <div className={cn('relative', sizeClasses)}>
        <svg className="w-full h-full -rotate-90" viewBox="0 0 36 36">
          <circle
            cx="18"
            cy="18"
            r="15.5"
            fill="none"
            className="stroke-muted"
            strokeWidth="3"
          />
          <circle
            cx="18"
            cy="18"
            r="15.5"
            fill="none"
            className={bgColor}
            strokeWidth="3"
            strokeDasharray={`${percentage} 100`}
            strokeLinecap="round"
          />
        </svg>
        <div className={cn('absolute inset-0 flex items-center justify-center font-bold', textSize, color)}>
          {percentage.toFixed(0)}%
        </div>
      </div>
      <span className="text-xs text-muted-foreground">{label}</span>
    </div>
  );
};

export function MapPerformance({ mapStats, mapPool, economyByMap }: MapPerformanceProps) {
  if (!mapStats && !mapPool) return null;

  // Convert mapStats to array and sort by games played
  const mapStatsArray = mapStats 
    ? Object.values(mapStats)
        .filter(m => m.gamesPlayed > 0)
        .sort((a, b) => b.gamesPlayed - a.gamesPlayed)
    : [];

  // Normalize mapPool to consistent format (handle both mapName/map and gamesPlayed/games)
  const normalizedMapPool = mapPool && mapPool.length > 0 
    ? mapPool.map(m => ({
        mapName: m.mapName || m.map || 'unknown',
        gamesPlayed: m.gamesPlayed ?? m.games ?? 0,
        winRate: m.winRate,
        strength: m.strength || (m.winRate >= 0.6 ? 'strong' : m.winRate >= 0.45 ? 'average' : 'weak'),
      })).filter(m => m.gamesPlayed > 0)
    : [];

  // Use mapPool if available, otherwise derive from mapStats
  const displayMaps = normalizedMapPool.length > 0 
    ? normalizedMapPool
    : mapStatsArray.map(m => ({
        mapName: m.mapName,
        gamesPlayed: m.gamesPlayed,
        winRate: m.winRate,
        strength: m.winRate >= 0.6 ? 'strong' : m.winRate >= 0.45 ? 'average' : 'weak',
      }));

  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2">
          <Map className="w-5 h-5 text-primary" />
          Map Performance
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Map Pool Overview */}
        <div className="grid gap-3">
          {displayMaps.slice(0, 6).map((map, idx) => {
            const fullStats = mapStatsArray.find(m => m.mapName === map.mapName);
            const displayName = mapDisplayNames[map.mapName.toLowerCase()] || map.mapName;
            
            return (
              <motion.div
                key={map.mapName}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: idx * 0.05 }}
                className="p-4 bg-muted/30 rounded-lg"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center">
                      <Map className="w-5 h-5 text-primary" />
                    </div>
                    <div>
                      <div className="font-semibold">{displayName}</div>
                      <div className="text-xs text-muted-foreground">{map.gamesPlayed} games</div>
                    </div>
                  </div>
                  <StrengthBadge strength={map.strength} />
                </div>

                {fullStats && (
                  <div className="flex items-center justify-around">
                    <WinRateCircle rate={fullStats.winRate} label="Overall" size="md" />
                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Swords className="w-4 h-4" />
                    </div>
                    <WinRateCircle rate={fullStats.attackWinRate} label="Attack" />
                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Shield className="w-4 h-4" />
                    </div>
                    <WinRateCircle rate={fullStats.defenseWinRate} label="Defense" />
                  </div>
                )}
                
                {!fullStats && (
                  <div className="flex items-center justify-center">
                    <WinRateCircle rate={map.winRate} label="Win Rate" size="md" />
                  </div>
                )}
              </motion.div>
            );
          })}
        </div>

        {/* Legend */}
        <div className="flex items-center justify-center gap-4 text-xs text-muted-foreground pt-2">
          <div className="flex items-center gap-1">
            <Swords className="w-3 h-3" />
            <span>Attack Side</span>
          </div>
          <div className="flex items-center gap-1">
            <Shield className="w-3 h-3" />
            <span>Defense Side</span>
          </div>
        </div>

        {/* Per-Map Economy (if available) */}
        {economyByMap && Object.keys(economyByMap).length > 0 && (
          <div className="pt-4 border-t border-border">
            <h4 className="text-sm font-medium mb-3 flex items-center gap-2">
              💰 Economy by Map
            </h4>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 text-xs">
              {Object.values(economyByMap).slice(0, 6).map(mapEcon => (
                <div key={mapEcon.mapName} className="p-2 bg-muted/20 rounded">
                  <div className="font-medium mb-1">
                    {mapDisplayNames[mapEcon.mapName.toLowerCase()] || mapEcon.mapName}
                  </div>
                  <div className="space-y-0.5 text-muted-foreground">
                    <div>Eco: {(mapEcon.ecoWinRate * 100).toFixed(0)}% ({mapEcon.ecoRounds})</div>
                    <div>Force: {(mapEcon.forceWinRate * 100).toFixed(0)}% ({mapEcon.forceRounds})</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
