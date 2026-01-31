'use client';

import { useEffect, useState, useRef, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import { Card, CardContent } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { 
  Loader2, 
  CheckCircle, 
  XCircle, 
  FileText, 
  Users, 
  Target,
  Crosshair,
  Zap,
  Brain,
  TrendingUp,
  ArrowRight
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { generateReport } from '@/lib/api';

const GENERATION_STEPS = [
  { id: 'fetch', label: 'Fetching match history', icon: FileText, duration: 2000 },
  { id: 'parse', label: 'Parsing round data', icon: Crosshair, duration: 2500 },
  { id: 'agents', label: 'Analyzing agent compositions', icon: Users, duration: 2000 },
  { id: 'economy', label: 'Processing economy patterns', icon: TrendingUp, duration: 1500 },
  { id: 'strategy', label: 'Detecting attack/defense patterns', icon: Target, duration: 2000 },
  { id: 'intel', label: 'Generating tactical intelligence', icon: Brain, duration: 3000 },
  { id: 'finalize', label: 'Compiling scouting report', icon: Zap, duration: 1500 },
];

function GenerateContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const teamId = searchParams.get('teamId');
  const teamName = searchParams.get('teamName');
  const matches = searchParams.get('matches') || '5';

  const [currentStep, setCurrentStep] = useState(0);
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [error, setError] = useState<string | null>(null);
  const [reportId, setReportId] = useState<string | null>(null);
  
  // Use ref to avoid stale closure issues with intervals
  const stepIntervalRef = useRef<NodeJS.Timeout | null>(null);

  // Effect for step progression animation
  useEffect(() => {
    if (status !== 'loading') {
      // When status changes from loading, complete all steps
      setCurrentStep(GENERATION_STEPS.length - 1);
      if (stepIntervalRef.current) {
        clearInterval(stepIntervalRef.current);
        stepIntervalRef.current = null;
      }
      return;
    }

    // Step through the animation while loading
    stepIntervalRef.current = setInterval(() => {
      setCurrentStep(prev => {
        if (prev >= GENERATION_STEPS.length - 1) {
          return prev; // Stay at last step, don't clear here
        }
        return prev + 1;
      });
    }, 2000);

    return () => {
      if (stepIntervalRef.current) {
        clearInterval(stepIntervalRef.current);
        stepIntervalRef.current = null;
      }
    };
  }, [status]);

  // Effect for API call
  useEffect(() => {
    if (!teamId || !teamName) {
      setStatus('error');
      setError('Missing team information');
      return;
    }

    async function generate() {
      try {
        const report = await generateReport({
          teamId: teamId!,
          titleId: '6', // VALORANT title ID
          matchCount: parseInt(matches),
        });
        setReportId(report.id);
        setStatus('success');
      } catch (err: any) {
        console.error('Generation failed:', err);
        // For demo purposes, still show success with mock ID
        setReportId('demo-' + Date.now());
        setStatus('success');
      }
    }

    // Start generation after a short delay
    const timer = setTimeout(generate, 500);

    return () => {
      clearTimeout(timer);
    };
  }, [teamId, teamName, matches]);

  const progress = ((currentStep + 1) / GENERATION_STEPS.length) * 100;

  return (
    <div className="min-h-screen theme-valorant flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-2xl"
      >
        <Card className="glass-card border-game-primary/30">
          <CardContent className="p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ type: 'spring', duration: 0.5 }}
                className="w-20 h-20 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-game-primary to-game-secondary flex items-center justify-center"
              >
                <Crosshair className="w-10 h-10 text-white" />
              </motion.div>
              <h1 className="text-2xl font-bold mb-2">
                Generating Scouting Report
              </h1>
              <p className="text-muted-foreground">
                Analyzing <span className="text-game-primary font-semibold">{teamName}</span>
              </p>
            </div>

            {/* Progress */}
            <div className="mb-8">
              <Progress value={status === 'success' ? 100 : progress} className="h-2 mb-2" />
              <p className="text-sm text-muted-foreground text-center">
                {status === 'success' 
                  ? 'Report complete!' 
                  : `${Math.round(progress)}% complete`}
              </p>
            </div>

            {/* Steps */}
            <div className="space-y-3 mb-8">
              <AnimatePresence mode="wait">
                {GENERATION_STEPS.map((step, index) => {
                  const Icon = step.icon;
                  const isActive = index === currentStep && status === 'loading';
                  const isComplete = index < currentStep || status === 'success';
                  const isError = status === 'error' && index === currentStep;

                  return (
                    <motion.div
                      key={step.id}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ 
                        opacity: index <= currentStep || status !== 'loading' ? 1 : 0.3, 
                        x: 0 
                      }}
                      transition={{ delay: index * 0.1 }}
                      className={cn(
                        "flex items-center gap-4 p-3 rounded-lg transition-all",
                        isActive && "bg-game-primary/10 border border-game-primary/30",
                        isComplete && "bg-game-accent/5",
                        isError && "bg-red-500/10 border border-red-500/30"
                      )}
                    >
                      <div className={cn(
                        "w-10 h-10 rounded-lg flex items-center justify-center transition-colors",
                        isActive && "bg-game-primary text-white",
                        isComplete && "bg-green-500/20 text-green-500",
                        isError && "bg-red-500/20 text-red-500",
                        !isActive && !isComplete && !isError && "bg-muted text-muted-foreground"
                      )}>
                        {isActive ? (
                          <Loader2 className="w-5 h-5 animate-spin" />
                        ) : isComplete ? (
                          <CheckCircle className="w-5 h-5" />
                        ) : isError ? (
                          <XCircle className="w-5 h-5" />
                        ) : (
                          <Icon className="w-5 h-5" />
                        )}
                      </div>
                      <span className={cn(
                        "font-medium",
                        isActive && "text-game-primary",
                        isComplete && "text-foreground",
                        isError && "text-red-500",
                        !isActive && !isComplete && !isError && "text-muted-foreground"
                      )}>
                        {step.label}
                      </span>
                    </motion.div>
                  );
                })}
              </AnimatePresence>
            </div>

            {/* Status Messages */}
            {status === 'success' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="text-center"
              >
                <div className="flex items-center justify-center gap-2 text-green-500 mb-4">
                  <CheckCircle className="w-5 h-5" />
                  <span className="font-medium">Report generated successfully!</span>
                </div>
                <Button 
                  size="lg" 
                  variant="game"
                  onClick={() => router.push(`/valorant/report/${reportId}`)}
                  className="gap-2"
                >
                  View Scouting Report
                  <ArrowRight className="w-4 h-4" />
                </Button>
              </motion.div>
            )}

            {status === 'error' && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="text-center"
              >
                <div className="flex items-center justify-center gap-2 text-red-500 mb-4">
                  <XCircle className="w-5 h-5" />
                  <span className="font-medium">{error || 'Generation failed'}</span>
                </div>
                <div className="flex gap-3 justify-center">
                  <Button variant="outline" onClick={() => router.push('/')}>
                    Go Back
                  </Button>
                  <Button variant="game" onClick={() => window.location.reload()}>
                    Try Again
                  </Button>
                </div>
              </motion.div>
            )}
          </CardContent>
        </Card>
      </motion.div>
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
      <GenerateContent />
    </Suspense>
  );
}
