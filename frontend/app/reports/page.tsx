'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  FileText, 
  Calendar, 
  Trophy,
  ArrowRight,
  Loader2,
  AlertCircle,
  Gamepad2
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface ReportListItem {
  id: string;
  teamName: string;
  title: 'lol' | 'valorant';
  generatedAt: string;
  matchesAnalyzed: number;
}

export default function ReportsPage() {
  const [reports, setReports] = useState<ReportListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchReports() {
      try {
        const res = await fetch('http://localhost:8080/api/reports');
        if (!res.ok) throw new Error('Failed to fetch reports');
        const data = await res.json();
        setReports(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load reports');
      } finally {
        setLoading(false);
      }
    }
    fetchReports();
  }, []);

  if (loading) {
    return (
      <div className="container mx-auto py-12 flex items-center justify-center min-h-[60vh]">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-game-primary" />
          <p className="text-muted-foreground">Loading reports...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto py-12 flex items-center justify-center min-h-[60vh]">
        <div className="flex flex-col items-center gap-4 text-center">
          <AlertCircle className="w-12 h-12 text-red-500" />
          <h2 className="text-xl font-semibold">Error Loading Reports</h2>
          <p className="text-muted-foreground">{error}</p>
          <Button onClick={() => window.location.reload()}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Scouting Reports</h1>
        <p className="text-muted-foreground">
          View all generated scouting reports for League of Legends and VALORANT teams
        </p>
      </div>

      {reports.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="py-12 flex flex-col items-center justify-center text-center">
            <FileText className="w-12 h-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No Reports Yet</h3>
            <p className="text-muted-foreground mb-4">
              Generate your first scouting report to see it here
            </p>
            <div className="flex gap-3">
              <Button asChild variant="outline">
                <Link href="/lol/generate">
                  <Gamepad2 className="w-4 h-4 mr-2" />
                  Generate LoL Report
                </Link>
              </Button>
              <Button asChild>
                <Link href="/valorant/generate">
                  <Gamepad2 className="w-4 h-4 mr-2" />
                  Generate VALORANT Report
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
          {reports.map((report, index) => (
            <ReportCard key={report.id} report={report} index={index} />
          ))}
        </div>
      )}
    </div>
  );
}

function ReportCard({ report, index }: { report: ReportListItem; index: number }) {
  const gameRoute = report.title === 'lol' ? 'lol' : 'valorant';
  const gameLabel = report.title === 'lol' ? 'League of Legends' : 'VALORANT';
  const gameColor = report.title === 'lol' ? 'text-blue-400' : 'text-red-400';
  const gameBg = report.title === 'lol' ? 'bg-blue-500/10' : 'bg-red-500/10';
  
  const formattedDate = new Date(report.generatedAt).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05 }}
    >
      <Link href={`/${gameRoute}/report/${report.id}`}>
        <Card className="h-full hover:border-game-primary/50 transition-all cursor-pointer group">
          <CardHeader className="pb-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <div className={cn(
                  "w-12 h-12 rounded-lg flex items-center justify-center font-bold text-lg",
                  gameBg, gameColor
                )}>
                  {report.teamName.substring(0, 2).toUpperCase()}
                </div>
                <div>
                  <CardTitle className="text-lg group-hover:text-game-primary transition-colors">
                    {report.teamName}
                  </CardTitle>
                  <Badge variant="outline" className={cn("text-xs", gameColor)}>
                    {gameLabel}
                  </Badge>
                </div>
              </div>
              <ArrowRight className="w-5 h-5 text-muted-foreground group-hover:text-game-primary group-hover:translate-x-1 transition-all" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-1">
                <Trophy className="w-4 h-4" />
                <span>{report.matchesAnalyzed} matches</span>
              </div>
              <div className="flex items-center gap-1">
                <Calendar className="w-4 h-4" />
                <span>{formattedDate}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </Link>
    </motion.div>
  );
}
