'use client';

import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { GameToggle } from '@/components/layout/GameToggle';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Slider } from '@/components/ui/slider';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { 
  Search, 
  Zap, 
  Target, 
  TrendingUp, 
  Users,
  ChevronRight,
  Loader2
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { GameTitle, Team } from '@/types';
import { ALL_TEAMS, FEATURED_TEAMS, searchTeamsLocally } from '@/lib/teams-data';

const features = [
  {
    icon: Target,
    title: 'How to Win',
    description: 'Data-backed strategies to exploit opponent weaknesses',
  },
  {
    icon: Users,
    title: 'Player Analysis',
    description: 'Champion pools, KDAs, and threat assessments',
  },
  {
    icon: TrendingUp,
    title: 'Trend Detection',
    description: 'Recent form, composition patterns, and meta picks',
  },
  {
    icon: Zap,
    title: 'Draft Insights',
    description: 'Priority bans and recommended counter-picks',
  },
];

export default function HomePage() {
  const router = useRouter();
  const [game, setGame] = useState<GameTitle>('lol');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
  const [matchCount, setMatchCount] = useState([10]);
  const [isGenerating, setIsGenerating] = useState(false);

  // Use local search - no API calls needed
  const searchResults = useMemo(() => {
    return searchTeamsLocally(searchQuery, game);
  }, [searchQuery, game]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    if (selectedTeam && query !== selectedTeam.name) {
      setSelectedTeam(null);
    }
  };

  const handleTeamSelect = (team: Team) => {
    setSelectedTeam(team);
    setSearchQuery(team.name);
    // Search results are derived from searchQuery via useMemo, so clearing happens automatically
  };

  const handleGenerate = async () => {
    if (!selectedTeam) return;
    
    setIsGenerating(true);
    // Navigate to the report generation page
    router.push(`/${game}/generate?teamId=${selectedTeam.id}&teamName=${encodeURIComponent(selectedTeam.name)}&matches=${matchCount[0]}`);
  };

  // Show search results if searching, otherwise show featured teams
  const displayTeams = searchQuery.length >= 2 ? searchResults : FEATURED_TEAMS[game];

  return (
    <div className={cn("min-h-screen", game === 'valorant' ? 'theme-valorant' : 'theme-lol')}>
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Background Effects */}
        <div className="absolute inset-0 bg-grid-pattern opacity-30" />
        <div className={cn(
          "absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[800px] rounded-full blur-3xl opacity-20",
          game === 'lol' ? 'bg-amber-500' : 'bg-red-500'
        )} />
        
        <div className="container mx-auto px-4 py-16 md:py-24 relative z-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="text-center mb-12"
          >
            <Badge variant="game" className="mb-4">
              Cloud9 x JetBrains Hackathon
            </Badge>
            <h1 className="text-4xl md:text-6xl font-bold mb-4 tracking-tight">
              <span className="text-game-primary">Scout9</span>
              <br />
              <span className="text-foreground">Automated Scouting Reports</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Generate data-driven scouting reports for any opponent. 
              Know their strategies, exploit their weaknesses, win your matches.
            </p>
          </motion.div>

          {/* Game Toggle */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="max-w-md mx-auto mb-8"
          >
            <GameToggle value={game} onChange={setGame} size="lg" />
          </motion.div>

          {/* Main Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="max-w-2xl mx-auto"
          >
            <Card glass className="p-8">
              <CardContent className="p-0 space-y-6">
                {/* Team Search */}
                <div className="space-y-2">
                  <label className="text-sm font-medium text-muted-foreground">
                    Select Opponent Team
                  </label>
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                    <Input
                      placeholder="Search for a team..."
                      value={searchQuery}
                      onChange={(e) => handleSearch(e.target.value)}
                      className="pl-10 h-12 text-lg"
                    />
                  </div>

                  {/* Search Results / Featured Teams */}
                  {displayTeams.length > 0 && (
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mt-4">
                      {displayTeams.slice(0, 12).map((team) => (
                        <button
                          key={team.id}
                          onClick={() => handleTeamSelect(team)}
                          className={cn(
                            "flex items-center gap-2 p-3 rounded-lg border transition-all text-left",
                            selectedTeam?.id === team.id
                              ? "border-game-primary bg-game-primary/10"
                              : "border-border hover:border-muted-foreground hover:bg-muted/50"
                          )}
                        >
                          {team.logoUrl ? (
                            <img 
                              src={team.logoUrl} 
                              alt={team.name} 
                              className="w-8 h-8 rounded object-contain bg-muted p-0.5"
                            />
                          ) : (
                            <div className="w-8 h-8 rounded bg-muted flex items-center justify-center text-xs font-bold">
                              {team.name.substring(0, 2).toUpperCase()}
                            </div>
                          )}
                          <span className="font-medium truncate">{team.name}</span>
                        </button>
                      ))}
                    </div>
                  )}
                </div>

                {/* Match Count Slider */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <label className="text-sm font-medium text-muted-foreground">
                      Matches to Analyze
                    </label>
                    <span className="text-lg font-bold text-game-primary">
                      {matchCount[0]}
                    </span>
                  </div>
                  <Slider
                    value={matchCount}
                    onValueChange={setMatchCount}
                    min={5}
                    max={50}
                    step={5}
                    className="w-full"
                  />
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>5 matches</span>
                    <span>50 matches</span>
                  </div>
                </div>

                {/* Generate Button */}
                <Button
                  variant="game"
                  size="xl"
                  className="w-full"
                  disabled={!selectedTeam || isGenerating}
                  onClick={handleGenerate}
                >
                  {isGenerating ? (
                    <>
                      <Loader2 className="w-5 h-5 animate-spin" />
                      Generating Report...
                    </>
                  ) : (
                    <>
                      Generate Scouting Report
                      <ChevronRight className="w-5 h-5" />
                    </>
                  )}
                </Button>
              </CardContent>
            </Card>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 border-t border-border/40">
        <div className="container mx-auto px-4">
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            className="text-center mb-12"
          >
            <h2 className="text-3xl font-bold mb-4">
              Everything You Need to Win
            </h2>
            <p className="text-muted-foreground max-w-xl mx-auto">
              Scout9 analyzes official GRID data to give you actionable insights
            </p>
          </motion.div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature, index) => {
              const Icon = feature.icon;
              return (
                <motion.div
                  key={feature.title}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.1 }}
                >
                  <Card className="h-full p-6 hover:border-game-primary/50 transition-colors">
                    <div className={cn(
                      "w-12 h-12 rounded-lg flex items-center justify-center mb-4",
                      "bg-game-primary/10"
                    )}>
                      <Icon className="w-6 h-6 text-game-primary" />
                    </div>
                    <h3 className="font-semibold mb-2">{feature.title}</h3>
                    <p className="text-sm text-muted-foreground">
                      {feature.description}
                    </p>
                  </Card>
                </motion.div>
              );
            })}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 border-t border-border/40">
        <div className="container mx-auto px-4 text-center">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            className="max-w-2xl mx-auto"
          >
            <h2 className="text-3xl font-bold mb-4">
              Ready to Scout Your Next Opponent?
            </h2>
            <p className="text-muted-foreground mb-8">
              Join professional esports teams using data-driven analysis to gain the competitive edge.
            </p>
            <Button
              variant="game"
              size="lg"
              onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
            >
              Get Started
              <ChevronRight className="w-5 h-5" />
            </Button>
          </motion.div>
        </div>
      </section>
    </div>
  );
}