// src/pages/StatusPage.tsx
import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import axios from 'axios'
import { format } from 'date-fns'
import { Activity, CheckCircle, XCircle, Clock } from 'lucide-react'
import { MonitorStatus } from '@/types'

interface PublicMonitor {
  name: string
  url: string
  current_status: MonitorStatus
  uptime_24h: number
  uptime_7d: number
  avg_response_ms: number
  last_checked_at: string | null
}

interface StatusPageData {
  system_status: 'operational' | 'degraded'
  monitors: PublicMonitor[]
  total: number
}

export function StatusPage() {
  const { userID } = useParams<{ userID: string }>()

  const { data, isLoading, isError } = useQuery({
    queryKey: ['status-page', userID],
    queryFn: async (): Promise<StatusPageData> => {
      const { data } = await axios.get(`/status/${userID}`)
      return data
    },
    refetchInterval: 60000,
  })

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full" />
      </div>
    )
  }

  if (isError || !data) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 text-lg">Страница не найдена</p>
        </div>
      </div>
    )
  }

  const isOperational = data.system_status === 'operational'

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-3xl mx-auto px-4 py-12">

        {/* Header */}
        <div className="flex items-center gap-3 mb-10">
          <div className="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center">
            <Activity size={20} className="text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">Статус сервисов</h1>
            <p className="text-gray-400 text-sm">
              Обновлено: {format(new Date(), 'dd.MM.yyyy HH:mm')}
            </p>
          </div>
        </div>

        {/* System Status Banner */}
        <div className={`flex items-center gap-4 p-5 rounded-xl mb-8 ${
          isOperational
            ? 'bg-emerald-500/10 border border-emerald-500/20'
            : 'bg-red-500/10 border border-red-500/20'
        }`}>
          {isOperational ? (
            <CheckCircle size={28} className="text-emerald-400 shrink-0" />
          ) : (
            <XCircle size={28} className="text-red-400 shrink-0" />
          )}
          <div>
            <p className={`text-lg font-semibold ${
              isOperational ? 'text-emerald-400' : 'text-red-400'
            }`}>
              {isOperational ? 'Все системы работают нормально' : 'Обнаружены проблемы'}
            </p>
            <p className="text-sm text-gray-400 mt-0.5">
              {data.total} {data.total === 1 ? 'сервис' : 'сервисов'} отслеживается
            </p>
          </div>
        </div>

        {/* Monitors List */}
        <div className="flex flex-col gap-3">
          {data.monitors.map((monitor, index) => (
            <div
              key={index}
              className="bg-gray-800 border border-gray-700 rounded-xl p-5"
            >
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h3 className="font-semibold text-white">{monitor.name}</h3>
                  <p className="text-sm text-gray-400 mt-0.5">{monitor.url}</p>
                </div>

                {/* Status */}
                <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium ${
                  monitor.current_status === 'up'
                    ? 'bg-emerald-500/10 text-emerald-400'
                    : monitor.current_status === 'down'
                    ? 'bg-red-500/10 text-red-400'
                    : 'bg-gray-500/10 text-gray-400'
                }`}>
                  <span className={`w-1.5 h-1.5 rounded-full ${
                    monitor.current_status === 'up'
                      ? 'bg-emerald-400 animate-pulse'
                      : monitor.current_status === 'down'
                      ? 'bg-red-400 animate-pulse'
                      : 'bg-gray-400'
                  }`} />
                  {monitor.current_status === 'up'
                    ? 'Online'
                    : monitor.current_status === 'down'
                    ? 'Offline'
                    : 'Pending'}
                </span>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <p className="text-xs text-gray-500 mb-1">Uptime 24ч</p>
                  <p className={`text-sm font-semibold ${
                    monitor.uptime_24h >= 99
                      ? 'text-emerald-400'
                      : monitor.uptime_24h >= 90
                      ? 'text-yellow-400'
                      : 'text-red-400'
                  }`}>
                    {monitor.uptime_24h.toFixed(2)}%
                  </p>
                </div>

                <div>
                  <p className="text-xs text-gray-500 mb-1">Uptime 7д</p>
                  <p className={`text-sm font-semibold ${
                    monitor.uptime_7d >= 99
                      ? 'text-emerald-400'
                      : monitor.uptime_7d >= 90
                      ? 'text-yellow-400'
                      : 'text-red-400'
                  }`}>
                    {monitor.uptime_7d.toFixed(2)}%
                  </p>
                </div>

                <div>
                  <p className="text-xs text-gray-500 mb-1">Avg Response</p>
                  <p className="text-sm font-semibold text-blue-400">
                    {monitor.avg_response_ms.toFixed(0)}мс
                  </p>
                </div>
              </div>

              {/* Last checked */}
              {monitor.last_checked_at && (
                <div className="flex items-center gap-1.5 mt-3 text-xs text-gray-500">
                  <Clock size={12} />
                  Последняя проверка: {format(new Date(monitor.last_checked_at), 'HH:mm:ss')}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Footer */}
        <div className="mt-10 text-center">
          <p className="text-gray-600 text-sm">
            Powered by{' '}
            <span className="text-indigo-400">UptimeMonitor</span>
          </p>
        </div>

      </div>
    </div>
  )
}