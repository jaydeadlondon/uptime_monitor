// src/pages/Dashboard.tsx
import { useState } from 'react'
import { Plus, Activity, AlertTriangle } from 'lucide-react'
import { monitorsApi } from '@/api/monitors'
import { analyticsApi } from '@/api/analytics'
import { Card } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { CreateMonitorModal } from '@/components/CreateMonitorModal'
import { MonitorCard } from '@/components/MonitorCard'
import { useQuery } from '@tanstack/react-query'

export function Dashboard() {
  const [showCreateModal, setShowCreateModal] = useState(false)

  const { data: monitorsData, isLoading: monitorsLoading, refetch: refetchMonitors } = useQuery({
    queryKey: ['monitors'],
    queryFn: () => monitorsApi.getAll(),
    refetchInterval: 30000,
  })

  const { data: overview } = useQuery({
    queryKey: ['analytics-overview'],
    queryFn: () => analyticsApi.getOverview(),
    refetchInterval: 60000,
  })

  const monitors = monitorsData?.data || []

  return (
    <div className="flex flex-col gap-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Dashboard</h1>
          <p className="text-gray-400 mt-1">
            {monitors.length} {monitors.length === 1 ? 'монитор' : 'мониторов'}
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus size={18} className="mr-2" />
          Добавить монитор
        </Button>
      </div>

      {/* Stats Cards */}
      {overview && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card padding="sm">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-emerald-500/10 rounded-lg flex items-center justify-center">
                <Activity size={20} className="text-emerald-400" />
              </div>
              <div>
                <p className="text-2xl font-bold text-white">{overview.up}</p>
                <p className="text-sm text-gray-400">Online</p>
              </div>
            </div>
          </Card>

          <Card padding="sm">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-red-500/10 rounded-lg flex items-center justify-center">
                <AlertTriangle size={20} className="text-red-400" />
              </div>
              <div>
                <p className="text-2xl font-bold text-white">{overview.down}</p>
                <p className="text-sm text-gray-400">Offline</p>
              </div>
            </div>
          </Card>

          <Card padding="sm">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-blue-500/10 rounded-lg flex items-center justify-center">
                <Activity size={20} className="text-blue-400" />
              </div>
              <div>
                <p className="text-2xl font-bold text-white">
                  {overview.overall_uptime.toFixed(1)}%
                </p>
                <p className="text-sm text-gray-400">Uptime 24h</p>
              </div>
            </div>
          </Card>

          <Card padding="sm">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-yellow-500/10 rounded-lg flex items-center justify-center">
                <AlertTriangle size={20} className="text-yellow-400" />
              </div>
              <div>
                <p className="text-2xl font-bold text-white">{overview.active_incidents}</p>
                <p className="text-sm text-gray-400">Инциденты</p>
              </div>
            </div>
          </Card>
        </div>
      )}

      {/* Loading */}
      {monitorsLoading && (
        <div className="flex items-center justify-center py-20">
          <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full" />
        </div>
      )}

      {/* Empty state */}
      {!monitorsLoading && monitors.length === 0 && (
        <div className="flex flex-col items-center justify-center py-20">
          <Activity size={48} className="text-gray-600 mb-4" />
          <h3 className="text-lg font-medium text-gray-400">Нет мониторов</h3>
          <p className="text-gray-500 mt-1 mb-4">Добавьте первый сайт для мониторинга</p>
          <Button onClick={() => setShowCreateModal(true)}>
            <Plus size={18} className="mr-2" />
            Добавить монитор
          </Button>
        </div>
      )}

      {/* Monitors Grid */}
      {!monitorsLoading && monitors.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {monitors.map((monitor) => (
            <MonitorCard
              key={monitor.id}
              monitor={monitor}
              onUpdate={refetchMonitors}
            />
          ))}
        </div>
      )}

      {/* Create Modal */}
      {showCreateModal && (
        <CreateMonitorModal
          onClose={() => setShowCreateModal(false)}
          onCreated={() => {
            setShowCreateModal(false)
            refetchMonitors()
          }}
        />
      )}
    </div>
  )
}