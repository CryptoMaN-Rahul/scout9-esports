'use client';

import { motion } from 'framer-motion';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Layers, 
  TrendingUp,
  Ban
} from 'lucide-react';
import { cn, formatPercentage } from '@/lib/utils';
import type { CompositionAnalysis, CompositionInsight } from '@/types';

interface CompositionsProps {
  data: CompositionAnalysis;
  game: 'lol' | 'valorant';
}

export function Compositions({ data, game }: CompositionsProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Layers className="w-5 h-5 text-game-primary" />
            Recent Compositions
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Top Compositions */}
          <div>
            <h4 className="text-sm font-medium text-muted-foreground mb-3">
              Most Played Compositions
            </h4>
            <div className="space-y-3">
              {data.topCompositions?.slice(0, 5).map((comp, index) => (
                <CompositionRow 
                  key={index} 
                  composition={comp} 
                  rank={index + 1}
                  game={game}
                />
              ))}
              {(!data.topCompositions || data.topCompositions.length === 0) && (
                <span className="text-sm text-muted-foreground">
                  No composition data available
                </span>
              )}
            </div>
          </div>

          {/* First Pick Priorities & Bans */}
          <div className="grid md:grid-cols-2 gap-6 pt-4 border-t border-border">
            {/* First Pick Priorities */}
            <div>
              <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-3">
                <TrendingUp className="w-4 h-4" />
                First Pick Priorities
              </h4>
              <div className="flex flex-wrap gap-2">
                {data.firstPickPriorities?.slice(0, 5).map((pick, index) => (
                  <Badge 
                    key={typeof pick === 'string' ? pick : pick.character} 
                    variant={index === 0 ? "game" : "secondary"}
                    className="text-sm"
                  >
                    {index + 1}. {typeof pick === 'string' ? pick : pick.character}
                    {typeof pick !== 'string' && pick.winRate !== undefined && (
                      <span className="ml-1 text-xs opacity-70">
                        ({(pick.winRate * 100).toFixed(0)}%)
                      </span>
                    )}
                  </Badge>
                ))}
                {(!data.firstPickPriorities || data.firstPickPriorities.length === 0) && (
                  <span className="text-sm text-muted-foreground">
                    No data available
                  </span>
                )}
              </div>
            </div>

            {/* Common Bans */}
            <div>
              <h4 className="flex items-center gap-2 text-sm font-medium text-muted-foreground mb-3">
                <Ban className="w-4 h-4" />
                Common Bans (Against Them)
              </h4>
              <div className="flex flex-wrap gap-2">
                {data.commonBans?.slice(0, 5).map((ban, index) => (
                  <Badge 
                    key={ban} 
                    variant={index === 0 ? "high" : "outline"}
                    className="text-sm"
                  >
                    {ban}
                  </Badge>
                ))}
                {(!data.commonBans || data.commonBans.length === 0) && (
                  <span className="text-sm text-muted-foreground">
                    No data available
                  </span>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}

function CompositionRow({ 
  composition, 
  rank,
  game 
}: { 
  composition: CompositionInsight; 
  rank: number;
  game: 'lol' | 'valorant';
}) {
  const winRateColor = composition.winRate >= 0.5 ? 'text-green-500' : 'text-red-500';
  const frequencyWidth = Math.min(composition.frequency * 100, 100);
  
  return (
    <div className="p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors">
      <div className="flex items-center gap-3">
        {/* Rank */}
        <div className="w-8 h-8 rounded-full bg-game-primary/10 flex items-center justify-center text-sm font-bold text-game-primary">
          #{rank}
        </div>

        {/* Characters */}
        <div className="flex-1 min-w-0">
          <div className="flex flex-wrap gap-1 mb-1">
            {composition.characters.map((char, i) => (
              <span key={i} className="text-sm">
                {char}
                {i < composition.characters.length - 1 && (
                  <span className="text-muted-foreground mx-1">·</span>
                )}
              </span>
            ))}
          </div>
          
          {/* Archetype */}
          {composition.archetype && (
            <Badge variant="outline" className="text-xs">
              {composition.archetype}
            </Badge>
          )}
        </div>

        {/* Stats */}
        <div className="flex items-center gap-4 text-sm">
          {/* Frequency */}
          <div className="text-right w-20">
            <div className="text-muted-foreground text-xs">Played</div>
            <div className="font-medium">
              {formatPercentage(composition.frequency)} ({composition.gamesPlayed})
            </div>
          </div>

          {/* Win Rate */}
          <div className="text-right w-16">
            <div className="text-muted-foreground text-xs">Win Rate</div>
            <div className={cn("font-bold", winRateColor)}>
              {formatPercentage(composition.winRate)}
            </div>
          </div>
        </div>
      </div>

      {/* Frequency Bar */}
      <div className="mt-2">
        <Progress 
          value={frequencyWidth} 
          className="h-1" 
          indicatorClassName={composition.winRate >= 0.5 ? 'bg-green-500/50' : 'bg-red-500/50'}
        />
      </div>
    </div>
  );
}
