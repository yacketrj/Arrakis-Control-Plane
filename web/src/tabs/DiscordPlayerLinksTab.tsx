import { useEffect, useState, type CSSProperties } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { discordSelfServiceApi, type DiscordPlayerLink } from '../api/discordSelfService'

type LinkForm = {
  discord_id: string
  player_id: string
  player_name: string
  notes: string
}

function fmt(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  return String(value)
}

function fmtDate(value: string): string {
  if (!value) return '—'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

function inputStyle(): CSSProperties {
  return {
    background: '#0d0b07',
    color: 'var(--color-text)',
    borderColor: '#2a2418',
    outline: 'none',
  }
}

export default function DiscordPlayerLinksTab() {
  const [links, setLinks] = useState<DiscordPlayerLink[]>([])
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState<LinkForm>({ discord_id: '', player_id: '', player_name: '', notes: '' })

  const load = async () => {
    setLoading(true)
    try {
      setLinks(await discordSelfServiceApi.listLinks())
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [])

  const save = async () => {
    try {
      const playerId = Number.parseInt(form.player_id, 10)
      await discordSelfServiceApi.upsertLink({
        discord_id: form.discord_id.trim(),
        player_id: playerId,
        player_name: form.player_name.trim() || undefined,
        notes: form.notes.trim() || undefined,
      })
      toast.success('Discord player link saved')
      setForm({ discord_id: '', player_id: '', player_name: '', notes: '' })
      await load()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  const edit = (link: DiscordPlayerLink) => {
    setForm({
      discord_id: link.discord_id,
      player_id: String(link.player_id),
      player_name: link.player_name ?? '',
      notes: link.notes ?? '',
    })
  }

  const remove = async (link: DiscordPlayerLink) => {
    const ok = window.confirm(`Delete Discord player link for ${link.discord_id}?`)
    if (!ok) return
    try {
      await discordSelfServiceApi.deleteLink(link.discord_id)
      toast.success('Discord player link deleted')
      await load()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 className="text-lg font-semibold" style={{ color: 'var(--color-primary)' }}>Discord Player Links</h2>
            <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
              Admin-managed Discord identity to player actor mapping. This enables read-only self-service lookup and does not grant mutation access.
            </p>
          </div>
          <Button variant="outline" size="sm" onPress={load} isDisabled={loading}>
            {loading ? <Spinner size="sm" color="current" /> : null}
            Refresh
          </Button>
        </div>
      </div>

      <section className="rounded-lg p-4 flex flex-col gap-3" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h3 className="font-semibold" style={{ color: 'var(--color-primary)' }}>Create or update link</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
          <input value={form.discord_id} onChange={e => setForm({ ...form, discord_id: e.target.value })} placeholder="Discord user ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
          <input value={form.player_id} onChange={e => setForm({ ...form, player_id: e.target.value })} placeholder="Player actor ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
          <input value={form.player_name} onChange={e => setForm({ ...form, player_name: e.target.value })} placeholder="Player name optional" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
          <textarea value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} placeholder="Notes optional" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
        </div>
        <div className="flex gap-2">
          <Button variant="primary" size="sm" onPress={save}>Save link</Button>
          <Button variant="outline" size="sm" onPress={() => setForm({ discord_id: '', player_id: '', player_name: '', notes: '' })}>Clear</Button>
        </div>
      </section>

      <section className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Current links</h3>
        {loading && links.length === 0 ? <Spinner size="lg" /> : null}
        {!loading && links.length === 0 ? (
          <div className="rounded p-4 text-sm" style={{ background: '#0d0b07', border: '1px solid #1a1610', color: 'var(--color-text-dim)' }}>
            No Discord player links have been created.
          </div>
        ) : null}
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-3">
          {links.map(link => (
            <article key={link.discord_id} className="rounded p-3 text-sm" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
              <div className="flex justify-between gap-3">
                <div>
                  <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{fmt(link.player_name) || `Player ${link.player_id}`}</div>
                  <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Discord {link.discord_id} · Player {link.player_id}</div>
                </div>
                <div className="flex gap-2">
                  <Button size="sm" variant="secondary" onPress={() => edit(link)}>Edit</Button>
                  <Button size="sm" variant="danger-soft" onPress={() => remove(link)}>Delete</Button>
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2 mt-3 text-xs" style={{ color: 'var(--color-text-dim)' }}>
                <div>Linked by: {fmt(link.linked_by_auth_type)} {link.linked_by_discord_id ? `(${link.linked_by_discord_id})` : ''}</div>
                <div>Updated: {fmtDate(link.updated_at)}</div>
                <div>Created: {fmtDate(link.created_at)}</div>
              </div>
              {link.notes ? <p className="mt-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{link.notes}</p> : null}
            </article>
          ))}
        </div>
      </section>
    </div>
  )
}
