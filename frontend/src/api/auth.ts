import client from './client'
import { AuthResponse, User } from '@/types'

export const authApi = {
  register: async (email: string, password: string): Promise<AuthResponse> => {
    const { data } = await client.post('/auth/register', { email, password })
    return data
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const { data } = await client.post('/auth/login', { email, password })
    return data
  },

  me: async (): Promise<User> => {
    const { data } = await client.get('/auth/me')
    return data
  },
}