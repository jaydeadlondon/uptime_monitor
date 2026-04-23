import client from './client'
import { Overview, ChartPoint, Incident } from '@/types'

export const analyticsApi = {
  getOverview: async (): Promise<Overview> => {
    const { data } = await client.get('/analytics/overview')
    return data
  },

  getChart: async (
    id: number,
    period: '24h' | '7d' | '30d' = '24h'
  ): Promise<{ data: ChartPoint[]; period: string }> => {
    const { data } = await client.get(`/analytics/monitors/${id}/chart`, {
      params: { period },
    })
    return data
  },

  getIncidents: async (
    id: number
  ): Promise<{ data: Incident[]; total: number }> => {
    const { data } = await client.get(`/analytics/monitors/${id}/incidents`)
    return data
  },
}