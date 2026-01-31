// Scout9 Constants
// Game data, colors, icons, and static configuration

export const GAMES = {
  lol: {
    id: '3',
    name: 'League of Legends',
    shortName: 'LoL',
    icon: '/icons/lol.svg',
  },
  valorant: {
    id: '6',
    name: 'VALORANT',
    shortName: 'VAL',
    icon: '/icons/valorant.svg',
  },
} as const;

export const LOL_ROLES = ['Top', 'Jungle', 'Mid', 'Bot', 'Support'] as const;
export const VALORANT_ROLES = ['Duelist', 'Initiator', 'Controller', 'Sentinel'] as const;

// LoL Champion Classes (for matchup analysis)
export const LOL_CHAMPION_CLASSES = {
  assassin: ['Zed', 'LeBlanc', 'Akali', 'Talon', 'Katarina', 'Fizz', 'Ekko', 'Qiyana'],
  mage: ['Azir', 'Orianna', 'Syndra', 'Viktor', 'Ahri', 'Lux', 'Xerath', 'Veigar'],
  tank: ['Ornn', 'Sejuani', 'Maokai', 'Sion', 'Malphite', 'Zac', 'Nautilus'],
  fighter: ['Aatrox', 'Renekton', 'Jax', 'Fiora', 'Camille', 'Irelia', 'Riven'],
  marksman: ['Jinx', 'Kai\'Sa', 'Aphelios', 'Zeri', 'Varus', 'Xayah', 'Ezreal'],
  support: ['Nautilus', 'Thresh', 'Leona', 'Lulu', 'Karma', 'Renata', 'Alistar'],
} as const;

// VALORANT Agents by Role
export const VALORANT_AGENTS = {
  duelist: ['Jett', 'Raze', 'Phoenix', 'Reyna', 'Yoru', 'Neon', 'Iso'],
  initiator: ['Sova', 'Breach', 'Skye', 'KAY/O', 'Fade', 'Gekko'],
  controller: ['Brimstone', 'Omen', 'Viper', 'Astra', 'Harbor', 'Clove'],
  sentinel: ['Sage', 'Cypher', 'Killjoy', 'Chamber', 'Deadlock', 'Vyse'],
} as const;

// VALORANT Maps
export const VALORANT_MAPS = [
  'Ascent',
  'Bind',
  'Haven',
  'Split',
  'Icebox',
  'Breeze',
  'Fracture',
  'Pearl',
  'Lotus',
  'Sunset',
  'Abyss',
] as const;

// Theme Colors
export const THEME_COLORS = {
  lol: {
    primary: '#C89B3C',
    secondary: '#0A1428',
    accent: '#0AC8B9',
    background: '#010A13',
    card: '#1E2328',
    text: '#F0E6D2',
    muted: '#785A28',
    gradient: 'from-amber-500 to-yellow-600',
  },
  valorant: {
    primary: '#FF4655',
    secondary: '#0F1923',
    accent: '#00D4AA',
    background: '#0F1923',
    card: '#1F2326',
    text: '#ECE8E1',
    muted: '#768079',
    gradient: 'from-red-500 to-pink-600',
  },
} as const;

// Impact Level Colors
export const IMPACT_COLORS = {
  HIGH: {
    bg: 'bg-red-500/10',
    border: 'border-red-500/30',
    text: 'text-red-400',
    dot: 'bg-red-500',
  },
  MEDIUM: {
    bg: 'bg-yellow-500/10',
    border: 'border-yellow-500/30',
    text: 'text-yellow-400',
    dot: 'bg-yellow-500',
  },
  LOW: {
    bg: 'bg-green-500/10',
    border: 'border-green-500/30',
    text: 'text-green-400',
    dot: 'bg-green-500',
  },
} as const;

// Action Type Icons (for How to Win section)
export const ACTION_TYPE_ICONS = {
  BAN: '🚫',
  PICK: '✅',
  TARGET_PLAYER: '🎯',
  FORCE_MAP: '🗺️',
  STRATEGY: '⚔️',
} as const;

// Placeholder team logos (for when logoUrl is missing)
export const PLACEHOLDER_TEAM_LOGO = '/images/team-placeholder.png';
export const PLACEHOLDER_CHAMPION_ICON = '/images/champion-placeholder.png';
export const PLACEHOLDER_AGENT_ICON = '/images/agent-placeholder.png';

// API Configuration
export const API_CONFIG = {
  baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  timeout: 120000, // 2 minutes for report generation
} as const;

// LCK Tournament IDs (for demo)
export const DEMO_TOURNAMENTS = {
  lol: [
    { id: '825490', name: 'LCK - Split 2 2025' },
    { id: '774794', name: 'LCK - Summer 2024' },
    { id: '758024', name: 'LCK - Spring 2024' },
  ],
  valorant: [
    { id: '826660', name: 'VCT Americas - Stage 2 2025' },
    { id: '800675', name: 'VCT Americas - Stage 1 2025' },
    { id: '775516', name: 'VCT Americas - Kickoff 2025' },
  ],
} as const;

// Chart Colors
export const CHART_COLORS = {
  primary: '#C89B3C',
  secondary: '#0AC8B9',
  tertiary: '#FF4655',
  success: '#22C55E',
  warning: '#EAB308',
  danger: '#EF4444',
  muted: '#6B7280',
} as const;
