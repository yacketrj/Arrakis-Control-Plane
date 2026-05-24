import { useState, useEffect } from 'react'
import { api, getAdminToken } from '../api/client'
import type { Status } from '../api/client'

export function useStatus() {
  const [status, setStatus] = useState<Status | null>(null)

  useEffect(() => {
    if (!getAdminToken()) {
      setStatus(null)
      return
    }

    let cancelled = false
    const poll = async () => {
      try {
        const s = await api.status()
        if (!cancelled) setStatus(s)
      } catch {
        if (!cancelled) setStatus(null)
      }
    }
    poll()
    const id = setInterval(poll, 5000)
    return () => {
      cancelled = true
      clearInterval(id)
    }
  }, [])

  return status
}
