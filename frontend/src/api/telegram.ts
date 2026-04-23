import client from './client'

export const telegramApi = {
  generateLink: async (): Promise<{ link: string; token: string }> => {
    const { data } = await client.post('/telegram/link')
    return data
  },

  unlink: async (): Promise<void> => {
    await client.delete('/telegram/unlink')
  },

  getStatus: async (): Promise<{ linked: boolean; telegram_chat_id: number | null }> => {
    const { data } = await client.get('/telegram/status')
    return data
  },
}