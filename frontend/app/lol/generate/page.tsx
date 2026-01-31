'use client';

import { useEffect, useState, useRef, Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Loader2, 
  CheckCircle2, 
  XCircle, 
  ArrowRight,
  RefreshCw 
} from 'lucide-react';
import { generateReport, ApiError } from '@/lib/api';
import type { ScoutingReport } from '@/types';
import { cn } from '@/lib/utils';

const loadingSteps = [
  { id: 'fetch', label: 'Fetching match data from GRID API...', duration: 3000 },
  { id: 'analyze', label: 'Analyzing team strategies...', duration: 4000 },
  { id: 'players', label: 'Profiling player tendencies...', duration: 3000 },
  { id: 'compositions', label: 'Identifying composition patterns...', duration: 2000 },
  { id: 'counter', label: 'Generating counter-strategies...', duration: 3000 },
  { id: 'report', label: 'Compiling scouting report...', duration: 2000 },
];

function GeneratePageContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  
  const teamId = searchParams.get('teamId');
  const teamName = searchParams.get('teamName') || 'Unknown Team';
  const matches = parseInt(searchParams.get('matches') || '10');
  
  const [currentStep, setCurrentStep] = useState(0);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [report, setReport] = useState<ScoutingReport | null>(null);
  const [error, setError] = useState<string | null>(null);
  
  // Use refs to avoid stale closure issues
  const stepTimerRef = useRef<NodeJS.Timeout | null>(null);
  const progressTimerRef = useRef<NodeJS.Timeout | null>(null);

  // Effect for step/progress animation
  useEffect(() => {
    if (status !== 'loading') {
      // When status changes from loading, complete everything
      setCurrentStep(loadingSteps.length - 1);
      setProgress(100);
      if (stepTimerRef.current) {
        clearInterval(stepTimerRef.current);
        stepTimerRef.current = null;
      }
      if (progressTimerRef.current) {
        clearInterval(progressTimerRef.current);
        progressTimerRef.current = null;
      }
      return;
    }

    // Simulate loading steps while waiting for API
    stepTimerRef.current = setInterval(() => {
      setCurrentStep((prev) => {
        if (prev < loadingSteps.length - 1) return prev + 1;
        return prev;
      });
    }, 3000);

    // Progress animation
    progressTimerRef.current = setInterval(() => {
      setProgress((prev) => {
        if (prev < 95) return prev + 1;
        return prev;
      });
    }, 200);

    return () => {
      if (stepTimerRef.current) {
        clearInterval(stepTimerRef.current);
        stepTimerRef.current = null;
      }
      if (progressTimerRef.current) {
        clearInterval(progressTimerRef.current);
        progressTimerRef.current = null;
      }
    };
  }, [status]);

  // Effect for API call
  useEffect(() => {
    if (!teamId) {
      router.push('/lol');
      return;
    }

    const fetchReport = async () => {
      try {
        const result = await generateReport({
          teamId,
          teamName: decodeURIComponent(teamName),
          matchCount: matches,
          titleId: '3', // LoL
        });
        
        setReport(result);
        setStatus('success');
      } catch (err) {
        console.error('Failed to generate report:', err);
        setError(err instanceof ApiError ? err.message : 'Failed to generate report');
        setStatus('error');
      }
    };

    fetchReport();
  }, [teamId, teamName, matches, router]);

  const handleViewReport = () => {
    if (report) {
      router.push(`/lol/report/${report.id}`);
    }
  };

  const handleRetry = () => {
    setStatus('loading');
    setProgress(0);
    setCurrentStep(0);
    setError(null);
    window.location.reload();
  };

  return (
    <div className="min-h-screen theme-lol">
      <div className="container mx-auto px-4 py-16 flex items-center justify-center">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="w-full max-w-xl"
        >
          <Card glass className="p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold mb-2">
                {status === 'loading' && 'Generating Report'}
                {status === 'success' && 'Report Ready!'}
                {status === 'error' && 'Generation Failed'}
              </h1>
              <p className="text-muted-foreground">
                {decodeURIComponent(teamName)} · {matches} matches
              </p>
            </div>

            {/* Status Icon */}
            <div className="flex justify-center mb-8">
              {status === 'loading' && (
                <motion.div
                  animate={{ rotate: 360 }}
                  transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                  className="w-20 h-20 rounded-full bg-game-primary/10 flex items-center justify-center"
                >
                  <Loader2 className="w-10 h-10 text-game-primary" />
                </motion.div>
              )}
              {status === 'success' && (
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="w-20 h-20 rounded-full bg-green-500/10 flex items-center justify-center"
                >
                  <CheckCircle2 className="w-10 h-10 text-green-500" />
                </motion.div>
              )}
              {status === 'error' && (
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="w-20 h-20 rounded-full bg-red-500/10 flex items-center justify-center"
                >
                  <XCircle className="w-10 h-10 text-red-500" />
                </motion.div>
              )}
            </div>

            {/* Progress Bar */}
            {status === 'loading' && (
              <div className="mb-8">
                <Progress value={progress} className="h-2 mb-4" />
                <div className="space-y-2">
                  {loadingSteps.map((step, index) => (
                    <div
                      key={step.id}
                      className={cn(
                        "flex items-center gap-3 text-sm transition-all",
                        index === currentStep && "text-game-primary font-medium",
                        index < currentStep && "text-muted-foreground",
                        index > currentStep && "text-muted-foreground/50"
                      )}
                    >
                      {index < currentStep && (
                        <CheckCircle2 className="w-4 h-4 text-green-500" />
                      )}
                      {index === currentStep && (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      )}
                      {index > currentStep && (
                        <div className="w-4 h-4 rounded-full border border-muted-foreground/30" />
                      )}
                      {step.label}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Error Message */}
            {status === 'error' && (
              <div className="mb-8 p-4 rounded-lg bg-red-500/10 border border-red-500/30 text-center">
                <p className="text-red-400">{error}</p>
              </div>
            )}

            {/* Actions */}
            <div className="flex gap-4">
              {status === 'success' && (
                <Button
                  variant="game"
                  size="lg"
                  className="flex-1"
                  onClick={handleViewReport}
                >
                  View Report
                  <ArrowRight className="w-5 h-5" />
                </Button>
              )}
              {status === 'error' && (
                <>
                  <Button
                    variant="outline"
                    size="lg"
                    className="flex-1"
                    onClick={() => router.push('/lol')}
                  >
                    Go Back
                  </Button>
                  <Button
                    variant="game"
                    size="lg"
                    className="flex-1"
                    onClick={handleRetry}
                  >
                    <RefreshCw className="w-5 h-5" />
                    Retry
                  </Button>
                </>
              )}
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}

export default function GeneratePage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-game-primary" />
      </div>
    }>
      <GeneratePageContent />
    </Suspense>
  );
}
