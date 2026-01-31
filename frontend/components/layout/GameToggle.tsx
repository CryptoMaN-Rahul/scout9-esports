'use client';

import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';
import type { GameTitle } from '@/types';

interface GameToggleProps {
  value: GameTitle;
  onChange: (game: GameTitle) => void;
  size?: 'sm' | 'md' | 'lg';
}

export function GameToggle({ value, onChange, size = 'md' }: GameToggleProps) {
  const sizeClasses = {
    sm: 'h-10 text-sm',
    md: 'h-12 text-base',
    lg: 'h-14 text-lg',
  };

  const games: { id: GameTitle; label: string; shortLabel: string }[] = [
    { id: 'lol', label: 'League of Legends', shortLabel: 'LoL' },
    { id: 'valorant', label: 'VALORANT', shortLabel: 'VAL' },
  ];

  return (
    <div className={cn(
      "relative flex rounded-xl bg-muted p-1",
      sizeClasses[size]
    )}>
      {games.map((game) => (
        <button
          key={game.id}
          onClick={() => onChange(game.id)}
          className={cn(
            "relative z-10 flex-1 flex items-center justify-center gap-2 rounded-lg font-medium transition-colors",
            value === game.id
              ? "text-black"
              : "text-muted-foreground hover:text-foreground"
          )}
        >
          {/* Game Icon */}
          <span className="text-xl">
            {game.id === 'lol' ? '⚔️' : '🎯'}
          </span>
          <span className="hidden sm:inline">{game.label}</span>
          <span className="sm:hidden">{game.shortLabel}</span>
        </button>
      ))}

      {/* Animated Background */}
      <motion.div
        layout
        className={cn(
          "absolute top-1 bottom-1 rounded-lg",
          value === 'lol' 
            ? "bg-gradient-to-r from-amber-500 to-yellow-500" 
            : "bg-gradient-to-r from-red-500 to-pink-500"
        )}
        initial={false}
        animate={{
          left: value === 'lol' ? '4px' : '50%',
          right: value === 'lol' ? '50%' : '4px',
        }}
        transition={{ type: "spring", bounce: 0.2, duration: 0.5 }}
      />
    </div>
  );
}
