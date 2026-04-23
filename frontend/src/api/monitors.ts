import client from './client'
import { Monitor, MonitorCheck, MonitorStats } from '@/types'

export const monitorsApi = {
  getAll: async (): Promise<{ data: Monitor[]; total: number }> => {
    const { data } = await client.get('/monitors')
    return data
  },

  getOne: async (id: number): Promise<Monitor> => {
    const { data } = await client.get(`/monitors/${id}`)
    return data
  },

  create: async (payload: {
    name: string
    url: string
    interval: number
  }): Promise<Monitor> => {
    const { data } = await client.post('/monitors', payload)
    return data
  },

  update: async (
    id: number,
    payload: Partial<{
      name: string
      url: string
      interval: number
      is_active: boolean
    }>
  ): Promise<Monitor> => {
    const { data } = await client.put(`/monitors/${id}`, payload)
    return data
  },

  delete: async (id: number): Promise<void> => {
    await client.delete(`/monitors/${id}`)
  },

  getChecks: async (id: number): Promise<{ data: MonitorCheck[]; total: number }> => {
    const { data } = await client.get(`/monitors/${id}/checks`)
    return data
  },

  getStats: async (id: number): Promise<MonitorStats> => {
    const { data } = await client.get(`/monitors/${id}/stats`)
    return data
  },
}