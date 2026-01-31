'use client';

import {
  LineChart as RechartsLineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from 'recharts';

interface LineChartProps {
  data: Record<string, any>[];
  lines: {
    dataKey: string;
    name: string;
    color: string;
    strokeWidth?: number;
  }[];
  xAxisKey: string;
  height?: number;
  showGrid?: boolean;
  showLegend?: boolean;
  valueFormatter?: (value: number) => string;
}

export function LineChart({
  data,
  lines,
  xAxisKey,
  height = 300,
  showGrid = true,
  showLegend = true,
  valueFormatter = (v) => `${v}`,
}: LineChartProps) {
  return (
    <ResponsiveContainer width="100%" height={height}>
      <RechartsLineChart
        data={data}
        margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
      >
        {showGrid && (
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="hsl(var(--border))"
            opacity={0.5}
          />
        )}
        <XAxis
          dataKey={xAxisKey}
          tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
        />
        <YAxis tick={{ fill: 'hsl(var(--muted-foreground))' }} />
        <Tooltip
          formatter={(value, name) => [valueFormatter(value as number), name]}
          contentStyle={{
            backgroundColor: 'hsl(var(--card))',
            border: '1px solid hsl(var(--border))',
            borderRadius: '8px',
          }}
        />
        {showLegend && (
          <Legend
            wrapperStyle={{ paddingTop: '20px' }}
          />
        )}
        {lines.map((line) => (
          <Line
            key={line.dataKey}
            type="monotone"
            dataKey={line.dataKey}
            name={line.name}
            stroke={line.color}
            strokeWidth={line.strokeWidth || 2}
            dot={{ fill: line.color, strokeWidth: 2 }}
            activeDot={{ r: 6 }}
          />
        ))}
      </RechartsLineChart>
    </ResponsiveContainer>
  );
}
