import { clsx } from 'clsx'
import { MonitorStatus } from '@/types'

interface BadgeProps {
  status: MonitorStatus
  pulse?: boolean
}

export function StatusBadge({ status, pulse = true }: BadgeProps) {
  return (
    <span className={clsx(
      'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium',
      {
        'bg-emerald-500/10 text-emerald-400': status === 'up',
        'bg-red-500/10 text-red-400': status === 'down',
        'bg-gray-500/10 text-gray-400': status === 'pending',
      }
    )}>
      <span className={clsx(
        'w-1.5 h-1.5 rounded-full',
        {
          'bg-emerald-400': status === 'up',
          'bg-red-400': status === 'down',
          'bg-gray-400': status === 'pending',
          'animate-pulse': pulse && status !== 'pending',
        }
      )} />
      {status === 'up' ? 'Online' : status === 'down' ? 'Offline' : 'Pending'}
    </span>
  )
}