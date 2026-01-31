'use client';

import {
  RadarChart as RechartsRadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
  ResponsiveContainer,
  Tooltip,
} from 'recharts';

interface RadarChartProps {
  data: {
    label: string;
    value: number;
    fullMark?: number;
  }[];
  color?: string;
  secondaryData?: {
    label: string;
    value: number;
  }[];
  secondaryColor?: string;
}

export function RadarChart({ 
  data, 
  color = 'hsl(var(--game-primary))',
  secondaryData,
  secondaryColor = 'hsl(var(--game-secondary))',
}: RadarChartProps) {
  // Transform data for recharts
  const chartData = data.map((item, index) => ({
    subject: item.label,
    A: item.value,
    B: secondaryData?.[index]?.value || 0,
    fullMark: item.fullMark || 100,
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <RechartsRadarChart cx="50%" cy="50%" outerRadius="80%" data={chartData}>
        <PolarGrid stroke="hsl(var(--border))" />
        <PolarAngleAxis
          dataKey="subject"
          tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
        />
        <PolarRadiusAxis
          angle={30}
          domain={[0, 100]}
          tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 10 }}
        />
        <Radar
          name="Primary"
          dataKey="A"
          stroke={color}
          fill={color}
          fillOpacity={0.3}
          strokeWidth={2}
        />
        {secondaryData && (
          <Radar
            name="Secondary"
            dataKey="B"
            stroke={secondaryColor}
            fill={secondaryColor}
            fillOpacity={0.2}
            strokeWidth={2}
          />
        )}
        <Tooltip
          contentStyle={{
            backgroundColor: 'hsl(var(--card))',
            border: '1px solid hsl(var(--border))',
            borderRadius: '8px',
          }}
        />
      </RechartsRadarChart>
    </ResponsiveContainer>
  );
}
