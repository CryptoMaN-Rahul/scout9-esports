'use client';

import { motion } from 'framer-motion';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { 
  Target, 
  Shield, 
  Swords, 
  AlertTriangle,
  Clock,
  Ban,
  UserX,
  Crosshair
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { HowToWinSection, DraftRecommendation, InGameStrategy, PlayerTarget, WeaknessTarget } from '@/types';

interface HowToWinProps {
  data: HowToWinSection;
  game: 'lol' | 'valorant';
}

export function HowToWin({ data, game }: HowToWinProps) {
  const themeClass = game === 'valorant' ? 'theme-valorant' : 'theme-lol';
  
  // Get draft recommendations by type
  const bans = data.draftRecommendations?.filter(d => d.type === 'ban') || [];
  const picks = data.draftRecommendations?.filter(d => d.type === 'pick') || [];
  const targets = data.draftRecommendations?.filter(d => d.type === 'target') || [];
  
  return (
    <section className={cn("space-y-6", themeClass)}>
      {/* Section Header */}
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-lg bg-game-primary/20 flex items-center justify-center">
          <Target className="w-5 h-5 text-game-primary" />
        </div>
        <div>
          <h2 className="text-2xl font-bold">How to Win</h2>
          <p className="text-sm text-muted-foreground">
            Data-backed strategies to defeat this opponent
          </p>
        </div>
      </div>

      {/* Win Condition Banner */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <Card className="border-game-primary/50 bg-gradient-to-r from-game-primary/10 to-transparent">
          <CardContent className="p-6">
            <div className="flex items-start justify-between gap-4">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-2">
                  <Badge variant="game" className="uppercase text-xs">
                    Win Condition
                  </Badge>
                </div>
                <p className="text-xl font-semibold leading-relaxed">
                  {data.winCondition}
                </p>
              </div>
              <div className="flex flex-col items-center">
                <div className="text-3xl font-bold text-game-primary">
                  {data.confidenceScore}%
                </div>
                <span className="text-xs text-muted-foreground">Confidence</span>
              </div>
            </div>
            {data.warnings && data.warnings.length > 0 && (
              <div className="mt-4 flex items-center gap-2 text-sm text-yellow-500">
                <AlertTriangle className="w-4 h-4" />
                {data.warnings.join(', ')}
              </div>
            )}
          </CardContent>
        </Card>
      </motion.div>

      {/* Target Players */}
      {data.targetPlayers && data.targetPlayers.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Crosshair className="w-5 h-5 text-game-primary" />
                Target Players
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                {data.targetPlayers.slice(0, 6).map((player, index) => (
                  <PlayerTargetCard key={index} player={player} />
                ))}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      )}

      {/* Weaknesses */}
      {data.weaknesses && data.weaknesses.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.25 }}
        >
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-yellow-500" />
                Exploitable Weaknesses
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                {data.weaknesses.map((weakness, index) => (
                  <WeaknessCard key={index} weakness={weakness} />
                ))}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      )}

      {/* Draft Strategy */}
      {(bans.length > 0 || picks.length > 0) && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Swords className="w-5 h-5 text-game-primary" />
                Draft Strategy
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Priority Bans */}
              {bans.length > 0 && (
                <div>
                  <h4 className="flex items-center gap-2 text-sm font-medium text-red-400 mb-3">
                    <Ban className="w-4 h-4" />
                    Priority Bans
                  </h4>
                  <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
                    {bans.slice(0, 6).map((draft, index) => (
                      <DraftCard key={index} recommendation={draft} type="ban" />
                    ))}
                  </div>
                </div>
              )}

              {/* Recommended Picks */}
              {picks.length > 0 && (
                <div>
                  {bans.length > 0 && <Separator className="my-4" />}
                  <h4 className="flex items-center gap-2 text-sm font-medium text-green-400 mb-3">
                    <Shield className="w-4 h-4" />
                    Recommended Picks
                  </h4>
                  <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
                    {picks.slice(0, 6).map((draft, index) => (
                      <DraftCard key={index} recommendation={draft} type="pick" />
                    ))}
                  </div>
                </div>
              )}

              {/* Target Picks */}
              {targets.length > 0 && (
                <div>
                  <Separator className="my-4" />
                  <h4 className="flex items-center gap-2 text-sm font-medium text-yellow-400 mb-3">
                    <UserX className="w-4 h-4" />
                    Force Their Picks
                  </h4>
                  <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
                    {targets.slice(0, 6).map((draft, index) => (
                      <DraftCard key={index} recommendation={draft} type="target" />
                    ))}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </motion.div>
      )}

      {/* In-Game Strategy Timeline */}
      {data.inGameStrategies && data.inGameStrategies.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="w-5 h-5 text-game-primary" />
                In-Game Strategy
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {data.inGameStrategies.map((strategy, index) => (
                  <InGameStrategyCard key={index} strategy={strategy} index={index} isLast={index === data.inGameStrategies!.length - 1} />
                ))}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      )}
    </section>
  );
}

function PlayerTargetCard({ player }: { player: PlayerTarget }) {
  const priorityColors: Record<number, string> = {
    1: 'border-red-500/50 bg-red-500/10',
    2: 'border-orange-500/50 bg-orange-500/10',
    3: 'border-yellow-500/50 bg-yellow-500/10',
  };
  
  const color = priorityColors[player.priority] || 'border-muted bg-muted/10';
  
  return (
    <div className={cn("p-4 rounded-lg border", color)}>
      <div className="flex items-start gap-3">
        <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center text-sm font-bold">
          {player.playerName.substring(0, 2).toUpperCase()}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="font-medium">{player.playerName}</span>
            <Badge variant="outline" className="text-xs">
              {player.role}
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground line-clamp-2">{player.reason}</p>
        </div>
      </div>
    </div>
  );
}

function WeaknessCard({ weakness }: { weakness: WeaknessTarget }) {
  return (
    <div className="p-4 rounded-lg border border-yellow-500/30 bg-yellow-500/5">
      <div className="flex items-start gap-3">
        <AlertTriangle className="w-5 h-5 text-yellow-500 mt-0.5" />
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span className="font-medium">{weakness.title}</span>
            <Badge variant="outline" className="text-xs">
              {weakness.impact}% impact
            </Badge>
          </div>
          <p className="text-sm text-muted-foreground">{weakness.description}</p>
          {weakness.evidence && (
            <p className="text-xs text-muted-foreground mt-1 italic">{weakness.evidence}</p>
          )}
        </div>
      </div>
    </div>
  );
}

function DraftCard({ recommendation, type }: { recommendation: DraftRecommendation; type: 'ban' | 'pick' | 'target' }) {
  const colorMap = {
    ban: 'border-red-500/30 bg-red-500/5',
    pick: 'border-green-500/30 bg-green-500/5',
    target: 'border-yellow-500/30 bg-yellow-500/5',
  };
  
  return (
    <div className={cn("p-3 rounded-lg border", colorMap[type])}>
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded bg-muted flex items-center justify-center text-sm font-bold">
          {recommendation.character.substring(0, 2).toUpperCase()}
        </div>
        <div className="flex-1 min-w-0">
          <p className="font-medium truncate">{recommendation.character}</p>
          <p className="text-xs text-muted-foreground">
            Priority {recommendation.priority}
          </p>
        </div>
      </div>
      <p className="text-xs text-muted-foreground mt-2 line-clamp-2">
        {recommendation.reason}
      </p>
    </div>
  );
}

function InGameStrategyCard({ strategy, index, isLast }: { strategy: InGameStrategy; index: number; isLast: boolean }) {
  return (
    <div className="flex gap-4">
      {/* Timeline */}
      <div className="flex flex-col items-center">
        <div className="w-8 h-8 rounded-full bg-game-primary/20 flex items-center justify-center text-sm font-bold text-game-primary">
          {index + 1}
        </div>
        {!isLast && <div className="w-px h-full bg-border mt-2" />}
      </div>
      
      {/* Content */}
      <div className="flex-1 pb-4">
        <div className="flex items-center gap-2 mb-1">
          <span className="font-medium">{strategy.title}</span>
          {strategy.timing && (
            <Badge variant="outline" className="text-xs">
              {strategy.timing}
            </Badge>
          )}
        </div>
        <p className="text-sm text-muted-foreground">{strategy.description}</p>
        {strategy.evidence && (
          <p className="text-xs text-muted-foreground mt-1 italic">📊 {strategy.evidence}</p>
        )}
      </div>
    </div>
  );
}
