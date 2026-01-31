'use client';

import {
  BarChart as RechartsBarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  LabelList,
} from 'recharts';
import { cn } from '@/lib/utils';

interface BarChartProps {
  data: {
    name: string;
    value: number;
    label?: string;
    color?: string;
  }[];
  horizontal?: boolean;
  showLabels?: boolean;
  showGrid?: boolean;
  height?: number;
  color?: string;
  valueFormatter?: (value: number) => string;
}

export function BarChart({
  data,
  horizontal = false,
  showLabels = true,
  showGrid = true,
  height = 300,
  color = 'hsl(var(--game-primary))',
  valueFormatter = (v) => `${v}`,
}: BarChartProps) {
  return (
    <ResponsiveContainer width="100%" height={height}>
      <RechartsBarChart
        data={data}
        layout={horizontal ? 'vertical' : 'horizontal'}
        margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
      >
        {showGrid && (
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="hsl(var(--border))"
            opacity={0.5}
          />
        )}
        {horizontal ? (
          <>
            <XAxis type="number" tick={{ fill: 'hsl(var(--muted-foreground))' }} />
            <YAxis
              dataKey="name"
              type="category"
              tick={{ fill: 'hsl(var(--muted-foreground))' }}
              width={100}
            />
          </>
        ) : (
          <>
            <XAxis
              dataKey="name"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
            />
            <YAxis tick={{ fill: 'hsl(var(--muted-foreground))' }} />
          </>
        )}
        <Tooltip
          formatter={(value) => [valueFormatter(value as number), 'Value']}
          contentStyle={{
            backgroundColor: 'hsl(var(--card))',
            border: '1px solid hsl(var(--border))',
            borderRadius: '8px',
          }}
        />
        <Bar dataKey="value" radius={[4, 4, 0, 0]}>
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color || color} />
          ))}
          {showLabels && (
            <LabelList
              dataKey="value"
              position="top"
              formatter={(value) => valueFormatter(value as number)}
              style={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
            />
          )}
        </Bar>
      </RechartsBarChart>
    </ResponsiveContainer>
  );
}

interface ComparisonBarChartProps {
  data: {
    name: string;
    value1: number;
    value2: number;
    label1?: string;
    label2?: string;
  }[];
  color1?: string;
  color2?: string;
  label1?: string;
  label2?: string;
  height?: number;
  valueFormatter?: (value: number) => string;
}

export function ComparisonBarChart({
  data,
  color1 = 'hsl(var(--game-primary))',
  color2 = 'hsl(var(--game-secondary))',
  label1 = 'Team 1',
  label2 = 'Team 2',
  height = 300,
  valueFormatter = (v) => `${v}`,
}: ComparisonBarChartProps) {
  return (
    <ResponsiveContainer width="100%" height={height}>
      <RechartsBarChart
        data={data}
        margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid
          strokeDasharray="3 3"
          stroke="hsl(var(--border))"
          opacity={0.5}
        />
        <XAxis
          dataKey="name"
          tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
        />
        <YAxis tick={{ fill: 'hsl(var(--muted-foreground))' }} />
        <Tooltip
          formatter={(value, name) => [
            valueFormatter(value as number),
            name === 'value1' ? label1 : label2,
          ]}
          contentStyle={{
            backgroundColor: 'hsl(var(--card))',
            border: '1px solid hsl(var(--border))',
            borderRadius: '8px',
          }}
        />
        <Bar dataKey="value1" name={label1} fill={color1} radius={[4, 4, 0, 0]} />
        <Bar dataKey="value2" name={label2} fill={color2} radius={[4, 4, 0, 0]} />
      </RechartsBarChart>
    </ResponsiveContainer>
  );
}
