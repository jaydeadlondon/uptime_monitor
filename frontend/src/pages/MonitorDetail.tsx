// src/pages/MonitorDetail.tsx
import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { format } from 'date-fns'
import {
  ArrowLeft,
  ExternalLink,
  Clock,
  CheckCircle,
  XCircle,
} from 'lucide-react'
import { monitorsApi } from '@/api/monitors'
import { analyticsApi } from '@/api/analytics'
import { Card } from '@/components/ui/Card'
import { StatusBadge } from '@/components/ui/Badge'
import { UptimeChart } from '@/components/UptimeChart'

type Period = '24h' | '7d' | '30d'

export function MonitorDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const monitorId = parseInt(id!)
  const [period, setPeriod] = useState<Period>('24h')

  const { data: monitor, isLoading } = useQuery({
    queryKey: ['monitor', monitorId],
    queryFn: () => monitorsApi.getOne(monitorId),
    refetchInterval: 30000,
  })

  const { data: stats } = useQuery({
    queryKey: ['monitor-stats', monitorId],
    queryFn: () => monitorsApi.getStats(monitorId),
    refetchInterval: 60000,
  })

  const { data: chartData } = useQuery({
    queryKey: ['monitor-chart', monitorId, period],
    queryFn: () => analyticsApi.getChart(monitorId, period),
    refetchInterval: 60000,
  })

  const { data: checksData } = useQuery({
    queryKey: ['monitor-checks', monitorId],
    queryFn: () => monitorsApi.getChecks(monitorId),
    refetchInterval: 30000,
  })

  const { data: incidentsData } = useQuery({
    queryKey: ['monitor-incidents', monitorId],
    queryFn: () => analyticsApi.getIncidents(monitorId),
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full" />
      </div>
    )
  }

  if (!monitor) {
    return (
      <div className="text-center py-20 text-gray-400">
        Монитор не найден
      </div>
    )
  }

  const checks = checksData?.data || []
  const incidents = incidentsData?.data || []

  return (
    <div className="flex flex-col gap-6">

      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate('/dashboard')}
          className="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
        >
          <ArrowLeft size={20} />
        </button>

        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-white">{monitor.name}</h1>
            <StatusBadge status={monitor.current_status} />
          </div>
          <div className="flex items-center gap-2 text-gray-400 mt-1">
            <ExternalLink size={14} />
            <a
              href={monitor.url}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-white transition-colors"
              onClick={(e) => e.stopPropagation()}
            >
              {monitor.url}
            </a>
          </div>
        </div>
      </div>

      {/* Stats */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card padding="sm">
            <p className="text-sm text-gray-400 mb-1">Uptime 24ч</p>
            <p className="text-2xl font-bold text-white">
              {stats.uptime_24h.toFixed(1)}%
            </p>
          </Card>
          <Card padding="sm">
            <p className="text-sm text-gray-400 mb-1">Uptime 7д</p>
            <p className="text-2xl font-bold text-white">
              {stats.uptime_7d.toFixed(1)}%
            </p>
          </Card>
          <Card padding="sm">
            <p className="text-sm text-gray-400 mb-1">Uptime 30д</p>
            <p className="text-2xl font-bold text-white">
              {stats.uptime_30d.toFixed(1)}%
            </p>
          </Card>
          <Card padding="sm">
            <p className="text-sm text-gray-400 mb-1">Avg Response</p>
            <p className="text-2xl font-bold text-white">
              {stats.avg_response_time.toFixed(0)}мс
            </p>
          </Card>
        </div>
      )}

      {/* Chart */}
      <Card>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-white">График</h2>
          <div className="flex gap-1">
            {(['24h', '7d', '30d'] as Period[]).map((p) => (
              <button
                key={p}
                onClick={() => setPeriod(p)}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                  period === p
                    ? 'bg-indigo-600 text-white'
                    : 'text-gray-400 hover:text-white hover:bg-gray-700'
                }`}
              >
                {p}
              </button>
            ))}
          </div>
        </div>

        <UptimeChart
          data={chartData?.data || []}
          period={period}
        />

        {/* Legend */}
        <div className="flex items-center gap-6 mt-4 justify-center">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-emerald-400" />
            <span className="text-sm text-gray-400">Uptime %</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-blue-400" />
            <span className="text-sm text-gray-400">Время ответа (мс)</span>
          </div>
        </div>
      </Card>

      {/* Incidents */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">
          Инциденты
          {incidents.length > 0 && (
            <span className="ml-2 text-sm text-gray-400 font-normal">
              ({incidents.length})
            </span>
          )}
        </h2>

        {incidents.length === 0 ? (
          <div className="flex items-center gap-2 text-emerald-400 py-4">
            <CheckCircle size={18} />
            <span className="text-sm">Инцидентов не было</span>
          </div>
        ) : (
          <div className="flex flex-col gap-3">
            {incidents.map((incident) => (
              <div
                key={incident.id}
                className="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  {incident.resolved_at ? (
                    <CheckCircle size={16} className="text-emerald-400 shrink-0" />
                  ) : (
                    <XCircle size={16} className="text-red-400 shrink-0 animate-pulse" />
                  )}
                  <div>
                    <p className="text-sm text-white">
                      {incident.resolved_at ? 'Устранён' : '🔴 Активный инцидент'}
                    </p>
                    <p className="text-xs text-gray-400 mt-0.5">
                      Начался: {format(new Date(incident.started_at), 'dd.MM.yyyy HH:mm')}
                    </p>
                  </div>
                </div>

                <div className="text-right">
                  {incident.duration_seconds && (
                    <p className="text-sm text-gray-300">
                      {incident.duration_seconds < 60
                        ? `${incident.duration_seconds}с`
                        : `${Math.floor(incident.duration_seconds / 60)}м`}
                    </p>
                  )}
                  {!incident.resolved_at && (
                    <span className="text-xs text-red-400">Идёт сейчас</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Last Checks */}
      <Card>
        <h2 className="text-lg font-semibold text-white mb-4">
          Последние проверки
        </h2>

        {checks.length === 0 ? (
          <p className="text-gray-400 text-sm py-4">Проверок пока нет</p>
        ) : (
          <div className="flex flex-col gap-2">
            {checks.slice(0, 20).map((check) => (
              <div
                key={check.id}
                className="flex items-center justify-between py-2 border-b border-gray-700/50 last:border-0"
              >
                <div className="flex items-center gap-3">
                  {check.status === 'up' ? (
                    <CheckCircle size={16} className="text-emerald-400 shrink-0" />
                  ) : (
                    <XCircle size={16} className="text-red-400 shrink-0" />
                  )}
                  <div>
                    <span className={`text-sm font-medium ${
                      check.status === 'up' ? 'text-emerald-400' : 'text-red-400'
                    }`}>
                      {check.status === 'up' ? 'Online' : 'Offline'}
                    </span>
                    {check.status_code && (
                      <span className="text-xs text-gray-500 ml-2">
                        HTTP {check.status_code}
                      </span>
                    )}
                    {check.error_message && (
                      <p className="text-xs text-gray-500 mt-0.5 truncate max-w-xs">
                        {check.error_message}
                      </p>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-4 text-right shrink-0">
                  {check.response_time_ms && (
                    <span className="flex items-center gap-1 text-sm text-gray-400">
                      <Clock size={12} />
                      {check.response_time_ms}мс
                    </span>
                  )}
                  <span className="text-xs text-gray-500">
                    {format(new Date(check.checked_at), 'HH:mm:ss')}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

    </div>
  )
}