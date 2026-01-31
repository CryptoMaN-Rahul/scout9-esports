import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatPercentage(value: number, decimals = 0): string {
  return `${(value * 100).toFixed(decimals)}%`;
}

export function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

export function formatDateTime(dateString: string): string {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function formatDuration(minutes: number): string {
  const mins = Math.floor(minutes);
  const secs = Math.round((minutes - mins) * 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function getImpactColor(impact: 'HIGH' | 'MEDIUM' | 'LOW'): string {
  switch (impact) {
    case 'HIGH':
      return 'text-red-500 bg-red-500/10 border-red-500/30';
    case 'MEDIUM':
      return 'text-yellow-500 bg-yellow-500/10 border-yellow-500/30';
    case 'LOW':
      return 'text-green-500 bg-green-500/10 border-green-500/30';
  }
}

export function getThreatColor(level: number): string {
  if (level >= 8) return 'text-red-500';
  if (level >= 6) return 'text-orange-500';
  if (level >= 4) return 'text-yellow-500';
  return 'text-green-500';
}

export function getConfidenceColor(score: number): string {
  if (score >= 80) return 'text-green-500';
  if (score >= 60) return 'text-yellow-500';
  if (score >= 40) return 'text-orange-500';
  return 'text-red-500';
}

export function getWinRateColor(rate: number): string {
  if (rate >= 0.6) return 'text-green-500';
  if (rate >= 0.5) return 'text-yellow-500';
  if (rate >= 0.4) return 'text-orange-500';
  return 'text-red-500';
}

export function truncate(str: string, length: number): string {
  if (str.length <= length) return str;
  return str.slice(0, length) + '...';
}

export function getGameThemeClass(game: 'lol' | 'valorant'): string {
  return game === 'lol' ? 'theme-lol' : 'theme-valorant';
}

export function getTitleId(game: 'lol' | 'valorant'): string {
  return game === 'lol' ? '3' : '6';
}

export function getGameFromTitleId(titleId: string): 'lol' | 'valorant' {
  return titleId === '3' ? 'lol' : 'valorant';
}
