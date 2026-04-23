import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { ExternalLink, Link, Unlink, User, Shield } from 'lucide-react'
import toast from 'react-hot-toast'
import { telegramApi } from '@/api/telegram'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/store/authStore'
import { Card } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'

export function Settings() {
  const { user, setUser } = useAuthStore()
  const queryClient = useQueryClient()
  const [linkLoading, setLinkLoading] = useState(false)
  const [unlinkLoading, setUnlinkLoading] = useState(false)

  const { data: tgStatus, refetch: refetchTgStatus } = useQuery({
    queryKey: ['telegram-status'],
    queryFn: () => telegramApi.getStatus(),
  })

  const { refetch: refetchMe } = useQuery({
    queryKey: ['me'],
    queryFn: async () => {
      const me = await authApi.me()
      setUser(me)
      return me
    },
    enabled: false,
  })

  const handleLinkTelegram = async () => {
    setLinkLoading(true)
    try {
      const { link } = await telegramApi.generateLink()
      window.open(link, '_blank')
      toast.success('Откройте бота в Telegram и нажмите Start')

      setTimeout(async () => {
        await refetchTgStatus()
        await refetchMe()
      }, 5000)
    } catch {
      toast.error('Ошибка генерации ссылки')
    } finally {
      setLinkLoading(false)
    }
  }

  const handleUnlinkTelegram = async () => {
    if (!confirm('Отвязать Telegram?')) return
    setUnlinkLoading(true)
    try {
      await telegramApi.unlink()
      await refetchTgStatus()
      await refetchMe()
      toast.success('Telegram отвязан')
    } catch {
      toast.error('Ошибка')
    } finally {
      setUnlinkLoading(false)
    }
  }

  return (
    <div className="flex flex-col gap-6 max-w-2xl">
      <div>
        <h1 className="text-2xl font-bold text-white">Настройки</h1>
        <p className="text-gray-400 mt-1">Управление аккаунтом и уведомлениями</p>
      </div>

      {/* Профиль */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <div className="w-8 h-8 bg-indigo-600/20 rounded-lg flex items-center justify-center">
            <User size={16} className="text-indigo-400" />
          </div>
          <h2 className="text-lg font-semibold text-white">Профиль</h2>
        </div>

        <div className="flex flex-col gap-3">
          <div className="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg">
            <div>
              <p className="text-sm text-gray-400">Email</p>
              <p className="text-white font-medium">{user?.email}</p>
            </div>
          </div>

          <div className="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg">
            <div>
              <p className="text-sm text-gray-400">ID аккаунта</p>
              <p className="text-white font-medium">#{user?.id}</p>
            </div>
          </div>

          <div className="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg">
            <div>
              <p className="text-sm text-gray-400">Публичная статус-страница</p>
              <p className="text-indigo-400 font-medium text-sm mt-0.5">
                {window.location.origin}/status/{user?.id}
              </p>
            </div>
            <a
              href={`/status/${user?.id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-600 transition-colors"
            >
              <ExternalLink size={16} />
            </a>
          </div>
        </div>
      </Card>

      {/* Telegram */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <div className="w-8 h-8 bg-blue-600/20 rounded-lg flex items-center justify-center">
            <Shield size={16} className="text-blue-400" />
          </div>
          <h2 className="text-lg font-semibold text-white">Telegram уведомления</h2>
        </div>

        {tgStatus?.linked ? (
          <div className="flex flex-col gap-4">
            {/* Статус привязан */}
            <div className="flex items-center gap-3 p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-lg">
              <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              <div>
                <p className="text-emerald-400 font-medium">Telegram подключён</p>
                <p className="text-sm text-gray-400 mt-0.5">
                  Chat ID: {tgStatus.telegram_chat_id}
                </p>
              </div>
            </div>

            <p className="text-sm text-gray-400">
              Вы будете получать уведомления когда сайт упадёт или восстановится.
            </p>

            <Button
              variant="danger"
              onClick={handleUnlinkTelegram}
              loading={unlinkLoading}
              className="w-fit"
            >
              <Unlink size={16} className="mr-2" />
              Отвязать Telegram
            </Button>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {/* Статус не привязан */}
            <div className="flex items-center gap-3 p-4 bg-gray-700/50 border border-gray-600 rounded-lg">
              <div className="w-2 h-2 rounded-full bg-gray-400" />
              <p className="text-gray-400">Telegram не подключён</p>
            </div>

            <div className="flex flex-col gap-2">
              <p className="text-sm text-gray-400">
                Подключите Telegram чтобы получать уведомления о падении сайтов.
              </p>
              <ul className="text-sm text-gray-500 list-disc list-inside space-y-1">
                <li>Нажмите кнопку ниже</li>
                <li>Откроется бот в Telegram</li>
                <li>Нажмите Start в боте</li>
              </ul>
            </div>

            <Button
              onClick={handleLinkTelegram}
              loading={linkLoading}
              className="w-fit"
            >
              <Link size={16} className="mr-2" />
              Привязать Telegram
            </Button>
          </div>
        )}
      </Card>

      {/* Публичная ссылка */}
      <Card>
        <div className="flex items-center gap-3 mb-4">
          <div className="w-8 h-8 bg-purple-600/20 rounded-lg flex items-center justify-center">
            <ExternalLink size={16} className="text-purple-400" />
          </div>
          <h2 className="text-lg font-semibold text-white">Статус-страница</h2>
        </div>

        <p className="text-sm text-gray-400 mb-4">
          Публичная страница с состоянием всех ваших мониторов.
          Можно поделиться с командой или клиентами — авторизация не нужна.
        </p>

        <div className="flex items-center gap-3 p-3 bg-gray-700/50 rounded-lg">
          <code className="text-indigo-400 text-sm flex-1 truncate">
            {window.location.origin}/status/{user?.id}
          </code>
          <button
            onClick={() => {
              navigator.clipboard.writeText(
                `${window.location.origin}/status/${user?.id}`
              )
              toast.success('Ссылка скопирована!')
            }}
            className="text-xs text-gray-400 hover:text-white bg-gray-600 hover:bg-gray-500 px-3 py-1.5 rounded-lg transition-colors shrink-0"
          >
            Скопировать
          </button>
        </div>
      </Card>
    </div>
  )
}