'use client';

import { motion } from 'framer-motion';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { 
  Users, 
  Target,
  Skull,
  Shield,
  Crosshair,
  Zap,
  Users2,
  Award,
  Sword,
  Wand2
} from 'lucide-react';
import { cn, formatPercentage, getThreatColor } from '@/lib/utils';
import type { PlayerProfile } from '@/types';

interface PlayerTendenciesProps {
  players: PlayerProfile[];
  game: 'lol' | 'valorant';
}

export function PlayerTendencies({ players, game }: PlayerTendenciesProps) {
  // Sort by threat level
  const sortedPlayers = [...players].sort((a, b) => b.threatLevel - a.threatLevel);
  
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="w-5 h-5 text-game-primary" />
            Player Tendencies
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {sortedPlayers.map((player, index) => (
              <PlayerCard 
                key={player.playerId} 
                player={player} 
                game={game}
                index={index}
              />
            ))}
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}

function PlayerCard({ player, game, index }: { player: PlayerProfile; game: 'lol' | 'valorant'; index: number }) {
  const threatColor = getThreatColor(player.threatLevel);
  const isHighThreat = player.threatLevel >= 7;
  
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 * index }}
      className={cn(
        "p-4 rounded-lg border transition-all",
        isHighThreat 
          ? "border-red-500/30 bg-red-500/5" 
          : "border-border hover:border-muted-foreground"
      )}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-lg bg-muted flex items-center justify-center text-lg font-bold">
            {player.nickname.substring(0, 2).toUpperCase()}
          </div>
          <div>
            <h4 className="font-semibold">{player.nickname}</h4>
            <Badge variant="outline" className="text-xs">
              {player.role}
            </Badge>
          </div>
        </div>
        
        {/* Threat Level */}
        <div className="text-right">
          <div className={cn("text-xl font-bold", threatColor)}>
            {player.threatLevel}/10
          </div>
          <span className="text-xs text-muted-foreground">Threat</span>
        </div>
      </div>

      {/* KDA */}
      <div className="flex items-center gap-4 mb-3 text-sm">
        <div className="flex items-center gap-1">
          <Crosshair className="w-4 h-4 text-green-500" />
          <span className="font-medium">{player.avgKills.toFixed(1)}</span>
        </div>
        <div className="flex items-center gap-1">
          <Skull className="w-4 h-4 text-red-500" />
          <span className="font-medium">{player.avgDeaths.toFixed(1)}</span>
        </div>
        <div className="flex items-center gap-1">
          <Shield className="w-4 h-4 text-blue-500" />
          <span className="font-medium">{player.avgAssists.toFixed(1)}</span>
        </div>
        <div className="ml-auto">
          <span className="text-muted-foreground">KDA:</span>
          <span className="font-bold ml-1">{player.kda.toFixed(2)}</span>
        </div>
      </div>

      {/* Signature Picks */}
      {player.signaturePicks.length > 0 && (
        <div className="mb-3">
          <span className="text-xs text-muted-foreground block mb-1">
            Signature {game === 'lol' ? 'Champions' : 'Agents'}
          </span>
          <div className="flex flex-wrap gap-1">
            {player.signaturePicks.slice(0, 3).map((pick) => (
              <Badge key={pick} variant="secondary" className="text-xs">
                {pick}
              </Badge>
            ))}
          </div>
        </div>
      )}

      {/* Character Pool Preview */}
      {player.characterPool.length > 0 && (
        <div className="mb-3">
          <span className="text-xs text-muted-foreground block mb-2">
            Top Picks (by win rate)
          </span>
          <div className="space-y-1">
            {player.characterPool.slice(0, 3).map((char) => (
              <div key={char.character || char.name} className="flex items-center gap-2 text-xs">
                <span className="w-20 truncate capitalize">{char.character || char.name}</span>
                <Progress 
                  value={char.winRate * 100} 
                  className="flex-1 h-1.5"
                  indicatorClassName={char.winRate >= 0.5 ? 'bg-green-500' : 'bg-red-500'}
                />
                <span className="w-10 text-right text-muted-foreground">
                  {formatPercentage(char.winRate)}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Threat Reason */}
      {player.threatReason && (
        <div className="pt-3 border-t border-border">
          <div className="flex items-start gap-2 text-xs">
            <Target className={cn("w-3 h-3 mt-0.5 flex-shrink-0", threatColor)} />
            <span className="text-muted-foreground line-clamp-2">
              {player.threatReason}
            </span>
          </div>
        </div>
      )}

      {/* Tendencies */}
      {player.tendencies && player.tendencies.length > 0 && (
        <div className="mt-3 pt-3 border-t border-border">
          <span className="text-xs text-muted-foreground block mb-1">
            Key Tendencies
          </span>
          <ul className="text-xs space-y-1">
            {player.tendencies.slice(0, 2).map((tendency, i) => (
              <li key={i} className="text-muted-foreground">
                • {tendency}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Multikill Stats */}
      {player.multikillStats && (
        (player.multikillStats.doubles || player.multikillStats.doubleKills || 0) > 0 || 
        (player.multikillStats.triples || player.multikillStats.tripleKills || 0) > 0
      ) && (
        <div className="mt-3 pt-3 border-t border-border">
          <div className="flex items-center gap-1 mb-2">
            <Zap className="w-3 h-3 text-yellow-500" />
            <span className="text-xs text-muted-foreground">Multikills</span>
          </div>
          <div className="flex flex-wrap gap-2 text-xs">
            {(player.multikillStats.doubles || player.multikillStats.doubleKills || 0) > 0 && (
              <span className="px-2 py-0.5 bg-yellow-500/10 text-yellow-500 rounded">
                2K: {player.multikillStats.doubles || player.multikillStats.doubleKills}
              </span>
            )}
            {(player.multikillStats.triples || player.multikillStats.tripleKills || 0) > 0 && (
              <span className="px-2 py-0.5 bg-orange-500/10 text-orange-500 rounded">
                3K: {player.multikillStats.triples || player.multikillStats.tripleKills}
              </span>
            )}
            {(player.multikillStats.quadras || player.multikillStats.quadraKills || 0) > 0 && (
              <span className="px-2 py-0.5 bg-red-500/10 text-red-500 rounded">
                4K: {player.multikillStats.quadras || player.multikillStats.quadraKills}
              </span>
            )}
            {(player.multikillStats.pentas || player.multikillStats.pentaKills || 0) > 0 && (
              <span className="px-2 py-0.5 bg-purple-500/10 text-purple-500 rounded font-bold">
                ACE: {player.multikillStats.pentas || player.multikillStats.pentaKills}
              </span>
            )}
          </div>
        </div>
      )}

      {/* Synergy Partners */}
      {player.synergyPartners && player.synergyPartners.length > 0 && (
        <div className="mt-3 pt-3 border-t border-border">
          <div className="flex items-center gap-1 mb-2">
            <Users2 className="w-3 h-3 text-blue-400" />
            <span className="text-xs text-muted-foreground">Best Synergy</span>
          </div>
          <div className="space-y-1">
            {player.synergyPartners.slice(0, 2).map((partner, i) => (
              <div key={partner.playerId} className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground">
                  {partner.playerName || `Partner ${i + 1}`}
                </span>
                <div className="flex items-center gap-1">
                  <Award className="w-3 h-3 text-green-500" />
                  <span className="text-green-500 font-medium">
                    {(partner.synergyScore * 100).toFixed(0)}%
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Weapon Stats - VALORANT specific */}
      {game === 'valorant' && player.weaponStats && player.weaponStats.length > 0 && (
        <div className="mt-3 pt-3 border-t border-border">
          <div className="flex items-center gap-1 mb-2">
            <Sword className="w-3 h-3 text-orange-400" />
            <span className="text-xs text-muted-foreground">Top Weapons</span>
          </div>
          <div className="space-y-1.5">
            {player.weaponStats.slice(0, 3).map((weapon) => {
              const weaponName = weapon.weaponName || weapon.weapon || 'Unknown';
              const killShare = weapon.killShare || weapon.percentage || 0;
              return (
                <div key={weaponName} className="flex items-center gap-2 text-xs">
                  <span className="w-16 truncate capitalize font-medium">
                    {weaponName}
                  </span>
                  <Progress 
                    value={killShare * 100} 
                    className="flex-1 h-1.5"
                    indicatorClassName="bg-orange-500"
                  />
                  <span className="w-14 text-right text-muted-foreground">
                    {weapon.kills} kills
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Ability Usage - VALORANT specific */}
      {game === 'valorant' && player.abilityUsage && player.abilityUsage.length > 0 && (
        <div className="mt-3 pt-3 border-t border-border">
          <div className="flex items-center gap-1 mb-2">
            <Wand2 className="w-3 h-3 text-purple-400" />
            <span className="text-xs text-muted-foreground">Ability Usage</span>
          </div>
          <div className="flex flex-wrap gap-1">
            {player.abilityUsage.slice(0, 4).map((ability) => {
              const abilityName = ability.abilityName || ability.ability || ability.abilityId || 'Unknown';
              const usageCount = ability.usageCount || ability.uses || 0;
              return (
                <div 
                  key={abilityName}
                  className="px-2 py-1 bg-purple-500/10 rounded text-xs"
                  title={`${usageCount} total uses (${(ability.usagePerGame || 0).toFixed(2)}/game)`}
                >
                  <span className="text-purple-400 capitalize">{abilityName.replace(/-/g, ' ')}</span>
                  <span className="text-muted-foreground ml-1">×{usageCount}</span>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </motion.div>
  );
}
