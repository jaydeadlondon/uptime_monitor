import { useNavigate } from 'react-router-dom'
import { ExternalLink, Clock, Trash2, Power } from 'lucide-react'
import { Monitor } from '@/types'
import { Card } from '@/components/ui/Card'
import { StatusBadge } from '@/components/ui/Badge'
import { monitorsApi } from '@/api/monitors'
import toast from 'react-hot-toast'
import { useState } from 'react'

interface MonitorCardProps {
  monitor: Monitor
  onUpdate: () => void
}

export function MonitorCard({ monitor, onUpdate }: MonitorCardProps) {
  const navigate = useNavigate()
  const [toggling, setToggling] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const handleToggle = async (e: React.MouseEvent) => {
    e.stopPropagation()
    setToggling(true)
    try {
      await monitorsApi.update(monitor.id, { is_active: !monitor.is_active })
      toast.success(monitor.is_active ? 'Монитор отключен' : 'Монитор включен')
      onUpdate()
    } catch {
      toast.error('Ошибка')
    } finally {
      setToggling(false)
    }
  }

  const handleDelete = async (e: React.MouseEvent) => {
    e.stopPropagation()
    if (!confirm('Удалить монитор?')) return
    setDeleting(true)
    try {
      await monitorsApi.delete(monitor.id)
      toast.success('Монитор удален')
      onUpdate()
    } catch {
      toast.error('Ошибка удаления')
    } finally {
      setDeleting(false)
    }
  }

  const timeSinceLastCheck = () => {
    if (!monitor.last_checked_at) return 'Еще не проверялся'
    const diff = Date.now() - new Date(monitor.last_checked_at).getTime()
    const minutes = Math.floor(diff / 60000)
    if (minutes < 1) return 'Только что'
    if (minutes < 60) return `${minutes} мин назад`
    const hours = Math.floor(minutes / 60)
    return `${hours} ч назад`
  }

  return (
    <Card
      padding="md"
      className="hover:border-gray-600 transition-colors cursor-pointer relative group"
      onClick={() => navigate(`/monitors/${monitor.id}`)}
    >
      <div className="flex items-start justify-between">
        {/* Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 mb-2">
            <h3 className="font-semibold text-white truncate">
              {monitor.name}
            </h3>
            <StatusBadge status={monitor.current_status} />
            {!monitor.is_active && (
              <span className="text-xs text-gray-500 bg-gray-700 px-2 py-0.5 rounded">
                Выключен
              </span>
            )}
          </div>

          <div className="flex items-center gap-2 text-sm text-gray-400">
            <ExternalLink size={14} />
            <span className="truncate">{monitor.url}</span>
          </div>

          <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
            <span className="flex items-center gap-1">
              <Clock size={12} />
              {timeSinceLastCheck()}
            </span>
            <span>Интервал: {monitor.interval} мин</span>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-1 ml-4 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={handleToggle}
            disabled={toggling}
            className="p-2 rounded-lg hover:bg-gray-700 text-gray-400 hover:text-yellow-400 transition-colors"
            title={monitor.is_active ? 'Выключить' : 'Включить'}
          >
            <Power size={16} className={toggling ? 'animate-spin' : ''} />
          </button>
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="p-2 rounded-lg hover:bg-gray-700 text-gray-400 hover:text-red-400 transition-colors"
            title="Удалить"
          >
            <Trash2 size={16} className={deleting ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>
    </Card>
  )
}