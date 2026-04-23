export type MonitorStatus = 'up' | 'down' | 'pending'
export type CheckStatus = 'up' | 'down'

export interface User {
  id: number
  email: string
  telegram_chat_id: number | null
  created_at: string
  updated_at: string
}

export interface Monitor {
  id: number
  user_id: number
  name: string
  url: string
  interval: number
  is_active: boolean
  current_status: MonitorStatus
  last_checked_at: string | null
  created_at: string
  updated_at: string
}

export interface MonitorCheck {
  id: number
  monitor_id: number
  status: CheckStatus
  status_code: number | null
  response_time_ms: number | null
  error_message: string | null
  checked_at: string
}

export interface Incident {
  id: number
  monitor_id: number
  started_at: string
  resolved_at: string | null
  duration_seconds: number | null
}

export interface MonitorStats {
  uptime_24h: number
  uptime_7d: number
  uptime_30d: number
  avg_response_time: number
}

export interface Overview {
  total_monitors: number
  up: number
  down: number
  pending: number
  overall_uptime: number
  active_incidents: number
}

export interface ChartPoint {
  time: string
  avg_response_ms: number
  uptime_percent: number
  total_checks: number
}

export interface AuthResponse {
  token: string
  user: User
}