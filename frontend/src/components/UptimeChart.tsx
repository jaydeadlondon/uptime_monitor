import {
    AreaChart,
    Area,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
  } from 'recharts'
  import { ChartPoint } from '@/types'
  import { format } from 'date-fns'
  
  interface UptimeChartProps {
    data: ChartPoint[]
    period: '24h' | '7d' | '30d'
  }
  
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (!active || !payload?.length) return null
  
    return (
      <div className="bg-gray-800 border border-gray-700 rounded-lg p-3 shadow-xl">
        <p className="text-gray-400 text-xs mb-2">{label}</p>
        <p className="text-emerald-400 text-sm font-medium">
          Uptime: {payload[0]?.value?.toFixed(1)}%
        </p>
        <p className="text-blue-400 text-sm">
          Ответ: {payload[1]?.value?.toFixed(0)} мс
        </p>
      </div>
    )
  }
  
  export function UptimeChart({ data, period }: UptimeChartProps) {
    const formatted = data.map((point) => ({
      ...point,
      label: format(
        new Date(point.time),
        period === '30d' ? 'dd MMM' : 'HH:mm'
      ),
    }))
  
    if (data.length === 0) {
      return (
        <div className="flex items-center justify-center h-48 text-gray-500">
          Недостаточно данных для графика
        </div>
      )
    }
  
    return (
      <ResponsiveContainer width="100%" height={220}>
        <AreaChart data={formatted} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
          <defs>
            <linearGradient id="uptimeGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#10b981" stopOpacity={0.2} />
              <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="responseGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.2} />
              <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
            </linearGradient>
          </defs>
  
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
  
          <XAxis
            dataKey="label"
            tick={{ fill: '#6b7280', fontSize: 12 }}
            tickLine={false}
            axisLine={false}
          />
  
          <YAxis
            tick={{ fill: '#6b7280', fontSize: 12 }}
            tickLine={false}
            axisLine={false}
          />
  
          <Tooltip content={<CustomTooltip />} />
  
          <Area
            type="monotone"
            dataKey="uptime_percent"
            stroke="#10b981"
            strokeWidth={2}
            fill="url(#uptimeGradient)"
            name="Uptime %"
          />
  
          <Area
            type="monotone"
            dataKey="avg_response_ms"
            stroke="#3b82f6"
            strokeWidth={2}
            fill="url(#responseGradient)"
            name="Response ms"
          />
        </AreaChart>
      </ResponsiveContainer>
    )
  }