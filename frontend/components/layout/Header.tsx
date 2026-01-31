'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';
import { 
  Search, 
  FileText, 
  GitCompare, 
  Menu,
  X 
} from 'lucide-react';
import { useState } from 'react';

const navItems = [
  { href: '/', label: 'Home', icon: Search },
  { href: '/reports', label: 'Reports', icon: FileText },
  { href: '/compare', label: 'Compare', icon: GitCompare },
];

export function Header() {
  const pathname = usePathname();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  // Determine game from path
  const isValorant = pathname.startsWith('/valorant');
  const gameClass = isValorant ? 'theme-valorant' : 'theme-lol';

  return (
    <header className={cn("sticky top-0 z-50 w-full border-b border-border/40 bg-background/80 backdrop-blur-xl", gameClass)}>
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-3 group">
          <div className="relative">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-game-primary to-game-accent flex items-center justify-center font-bold text-black text-xl">
              S9
            </div>
            <div className="absolute inset-0 rounded-lg bg-game-primary/20 blur-lg group-hover:blur-xl transition-all" />
          </div>
          <div className="flex flex-col">
            <span className="font-bold text-lg tracking-tight">Scout9</span>
            <span className="text-xs text-muted-foreground">Scouting Reports</span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center gap-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href || 
              (item.href !== '/' && pathname.startsWith(item.href));
            const Icon = item.icon;
            
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "relative flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors",
                  isActive 
                    ? "text-foreground" 
                    : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
                )}
              >
                <Icon className="w-4 h-4" />
                {item.label}
                {isActive && (
                  <motion.div
                    layoutId="activeNav"
                    className="absolute inset-0 bg-muted rounded-lg -z-10"
                    transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
                  />
                )}
              </Link>
            );
          })}
        </nav>

        {/* Game Badge */}
        <div className="hidden md:flex items-center gap-4">
          <div className={cn(
            "px-3 py-1 rounded-full text-xs font-semibold border",
            isValorant 
              ? "border-red-500/30 bg-red-500/10 text-red-400" 
              : "border-amber-500/30 bg-amber-500/10 text-amber-400"
          )}>
            {isValorant ? 'VALORANT' : 'League of Legends'}
          </div>
        </div>

        {/* Mobile Menu Button */}
        <button
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          className="md:hidden p-2 rounded-lg hover:bg-muted"
        >
          {mobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
        </button>
      </div>

      {/* Mobile Menu */}
      {mobileMenuOpen && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -10 }}
          className="md:hidden border-t border-border bg-background"
        >
          <nav className="container mx-auto px-4 py-4 flex flex-col gap-2">
            {navItems.map((item) => {
              const isActive = pathname === item.href;
              const Icon = item.icon;
              
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setMobileMenuOpen(false)}
                  className={cn(
                    "flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors",
                    isActive 
                      ? "bg-muted text-foreground" 
                      : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
                  )}
                >
                  <Icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </motion.div>
      )}
    </header>
  );
}
