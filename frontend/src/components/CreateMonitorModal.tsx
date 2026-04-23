import { useState } from 'react'
import toast from 'react-hot-toast'
import { monitorsApi } from '@/api/monitors'
import { Modal } from '@/components/ui/Modal'
import { Input } from '@/components/ui/Input'
import { Button } from '@/components/ui/Button'

interface CreateMonitorModalProps {
  onClose: () => void
  onCreated: () => void
}

export function CreateMonitorModal({ onClose, onCreated }: CreateMonitorModalProps) {
  const [name, setName] = useState('')
  const [url, setUrl] = useState('')
  const [interval, setInterval] = useState('5')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!name || !url) {
      toast.error('Заполните все поля')
      return
    }

    let normalizedUrl = url
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      normalizedUrl = 'https://' + url
    }

    setLoading(true)
    try {
      await monitorsApi.create({
        name,
        url: normalizedUrl,
        interval: parseInt(interval) || 5,
      })
      toast.success('Монитор создан!')
      onCreated()
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Ошибка создания')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal isOpen={true} onClose={onClose} title="Новый монитор">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <Input
          label="Название"
          placeholder="Google"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />

        <Input
          label="URL"
          placeholder="https://google.com"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />

        <div>
          <label className="text-sm font-medium text-gray-300 block mb-1.5">
            Интервал проверки (минуты)
          </label>
          <select
            value={interval}
            onChange={(e) => setInterval(e.target.value)}
            className="w-full px-4 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="1">Каждую минуту</option>
            <option value="5">Каждые 5 минут</option>
            <option value="10">Каждые 10 минут</option>
            <option value="15">Каждые 15 минут</option>
            <option value="30">Каждые 30 минут</option>
            <option value="60">Каждый час</option>
          </select>
        </div>

        <div className="flex gap-3 mt-2">
          <Button
            type="button"
            variant="secondary"
            onClick={onClose}
            className="flex-1"
          >
            Отмена
          </Button>
          <Button
            type="submit"
            loading={loading}
            className="flex-1"
          >
            Создать
          </Button>
        </div>
      </form>
    </Modal>
  )
}