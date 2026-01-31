'use client';

import {
  PieChart as RechartsPieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from 'recharts';

interface PieChartProps {
  data: {
    name: string;
    value: number;
    color?: string;
  }[];
  colors?: string[];
  height?: number;
  innerRadius?: number;
  outerRadius?: number;
  showLabels?: boolean;
  showLegend?: boolean;
  valueFormatter?: (value: number) => string;
}

const DEFAULT_COLORS = [
  'hsl(var(--game-primary))',
  'hsl(var(--game-secondary))',
  'hsl(var(--game-accent))',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#06b6d4',
];

export function PieChart({
  data,
  colors = DEFAULT_COLORS,
  height = 300,
  innerRadius = 60,
  outerRadius = 100,
  showLabels = true,
  showLegend = true,
  valueFormatter = (v) => `${v}`,
}: PieChartProps) {
  const renderCustomizedLabel = ({
    cx,
    cy,
    midAngle,
    innerRadius,
    outerRadius,
    percent,
    name,
  }: any) => {
    const RADIAN = Math.PI / 180;
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * RADIAN);
    const y = cy + radius * Math.sin(-midAngle * RADIAN);

    return percent > 0.05 ? (
      <text
        x={x}
        y={y}
        fill="white"
        textAnchor={x > cx ? 'start' : 'end'}
        dominantBaseline="central"
        fontSize={12}
        fontWeight="bold"
      >
        {`${(percent * 100).toFixed(0)}%`}
      </text>
    ) : null;
  };

  return (
    <ResponsiveContainer width="100%" height={height}>
      <RechartsPieChart>
        <Pie
          data={data}
          cx="50%"
          cy="50%"
          labelLine={false}
          label={showLabels ? renderCustomizedLabel : undefined}
          outerRadius={outerRadius}
          innerRadius={innerRadius}
          fill="#8884d8"
          dataKey="value"
          animationBegin={0}
          animationDuration={800}
        >
          {data.map((entry, index) => (
            <Cell
              key={`cell-${index}`}
              fill={entry.color || colors[index % colors.length]}
            />
          ))}
        </Pie>
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
            layout="horizontal"
            verticalAlign="bottom"
            align="center"
          />
        )}
      </RechartsPieChart>
    </ResponsiveContainer>
  );
}

interface DonutStatProps {
  value: number;
  max: number;
  label: string;
  color?: string;
  size?: number;
  thickness?: number;
}

export function DonutStat({
  value,
  max,
  label,
  color = 'hsl(var(--game-primary))',
  size = 120,
  thickness = 12,
}: DonutStatProps) {
  const percentage = (value / max) * 100;
  const circumference = 2 * Math.PI * ((size - thickness) / 2);
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  return (
    <div className="flex flex-col items-center">
      <div className="relative" style={{ width: size, height: size }}>
        <svg
          width={size}
          height={size}
          className="transform -rotate-90"
        >
          {/* Background circle */}
          <circle
            cx={size / 2}
            cy={size / 2}
            r={(size - thickness) / 2}
            fill="none"
            stroke="hsl(var(--muted))"
            strokeWidth={thickness}
          />
          {/* Progress circle */}
          <circle
            cx={size / 2}
            cy={size / 2}
            r={(size - thickness) / 2}
            fill="none"
            stroke={color}
            strokeWidth={thickness}
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            className="transition-all duration-1000 ease-out"
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-2xl font-bold">{Math.round(percentage)}%</span>
        </div>
      </div>
      <span className="mt-2 text-sm text-muted-foreground text-center">{label}</span>
    </div>
  );
}
