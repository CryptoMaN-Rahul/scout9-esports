import Link from 'next/link';
import { Github, ExternalLink } from 'lucide-react';

export function Footer() {
  return (
    <footer className="border-t border-border/40 bg-background/50 backdrop-blur-sm">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          {/* Branding */}
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-game-primary to-game-accent flex items-center justify-center font-bold text-black text-sm">
              S9
            </div>
            <div className="text-sm text-muted-foreground">
              <span className="font-semibold text-foreground">Scout9</span>
              {' · '}
              Cloud9 x JetBrains Hackathon
            </div>
          </div>

          {/* Links */}
          <div className="flex items-center gap-6">
            <Link
              href="https://github.com/CryptoMaN-Rahul/scout9-esports
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              <Github className="w-4 h-4" />
              GitHub
            </Link>
            <Link
              href="https://docs.grid.gg"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              <ExternalLink className="w-4 h-4" />
              GRID API
            </Link>
          </div>

          {/* Copyright */}
          <div className="text-sm text-muted-foreground">
            Powered by{' '}
            <span className="font-semibold text-game-primary">GRID</span>
            {' '}Esports Data
          </div>
        </div>
      </div>
    </footer>
  );
}
