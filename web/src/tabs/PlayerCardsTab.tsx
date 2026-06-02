import { useEffect, useMemo, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api, type Player } from '../api/client'
import { getPlayerProfile, type PlayerProfile } from '../api/playerProfile'

const player360EventName = 'dune-admin-open-player-360'
const player360StorageKey = 'dune_admin_player_360_id'

function fmt(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  if (typeof value === 'number') return value.toLocaleString()
  return String(value)
}

function pct(numerator: number, denominator: number): number {
  if (!denominator || denominator <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((numerator / denominator) * 100)))
}

function ProgressBar({ value }: { value: number }) {
  return (
    <div className="h-2 rounded-full overflow-hidden" style={{ background: '#1a1610', border: '1px solid #2a2418' }}>
      <div className="h-full" style={{ width: `${value}%`, background: 'var(--color-primary)' }} />
    </div>
  )
}

function openPlayer360(playerId: number) {
  localStorage.setItem(player360StorageKey, String(playerId))
  window.dispatchEvent(new CustomEvent(player360EventName, { detail: { playerId } }))
}

function PlayerProgressSummary({ profile }: { profile: PlayerProfile }) {
  const journeyComplete = pct(profile.journey_summary.complete_nodes, profile.journey_summary.total_nodes)
  const revealed = pct(profile.journey_summary.revealed_nodes, profile.journey_summary.total_nodes)
  const topSpecs = profile.specializations
    .slice()
    .sort((a, b) => b.xp - a.xp)
    .slice(0, 3)

  return (
    <div className="mt-3 flex flex-col gap-3">
      <div className="grid grid-cols-3 gap-2 text-xs">
        <div className="rounded p-2" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
          <div style={{ color: 'var(--color-text-dim)' }}>Level</div>
          <div className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>{fmt(profile.character_xp?.level)}</div>
        </div>
        <div className="rounded p-2" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
          <div style={{ color: 'var(--color-text-dim)' }}>Inventory Rows</div>
          <div className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>{fmt(profile.inventory_summary.total_items)}</div>
        </div>
        <div className="rounded p-2" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
          <div style={{ color: 'var(--color-text-dim)' }}>Vehicles</div>
          <div className="text-lg font-semibold" style={{ color: 'var(--color-text)' }}>{fmt(profile.vehicles.length)}</div>
        </div>
      </div>

      <div className="text-xs flex flex-col gap-1">
        <div className="flex justify-between"><span style={{ color: 'var(--color-text-dim)' }}>Journey complete</span><span className="font-mono">{journeyComplete}%</span></div>
        <ProgressBar value={journeyComplete} />
        <div className="flex justify-between mt-1"><span style={{ color: 'var(--color-text-dim)' }}>Journey revealed</span><span className="font-mono">{revealed}%</span></div>
        <ProgressBar value={revealed} />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs">
        <div>
          <div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Top Specs</div>
          {topSpecs.length === 0 ? <div style={{ color: 'var(--color-text-dim)' }}>No spec rows.</div> : topSpecs.map(spec => (
            <div key={`${spec.track_type}-${spec.level}`} className="flex justify-between gap-3 py-0.5" style={{ borderBottom: '1px solid #1a1610' }}>
              <span style={{ color: 'var(--color-text)' }}>{spec.track_type}</span>
              <span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>Lv {fmt(spec.level)} · {fmt(spec.xp)} XP</span>
            </div>
          ))}
        </div>
        <div>
          <div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Factions</div>
          {profile.factions.length === 0 ? <div style={{ color: 'var(--color-text-dim)' }}>No faction rows.</div> : profile.factions.slice(0, 3).map(faction => (
            <div key={`${faction.faction_id}-${faction.actor_id}`} className="flex justify-between gap-3 py-0.5" style={{ borderBottom: '1px solid #1a1610' }}>
              <span style={{ color: 'var(--color-text)' }}>{faction.faction_name}</span>
              <span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>{fmt(faction.reputation)}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function PlayerCard({ player }: { player: Player }) {
  const [profile, setProfile] = useState<PlayerProfile | null>(null)
  const [loading, setLoading] = useState(false)

  const loadProgress = async () => {
    setLoading(true)
    try {
      setProfile(await getPlayerProfile(player.id))
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <article className="rounded-lg p-4 flex flex-col gap-3" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <h3 className="text-base font-semibold truncate" style={{ color: 'var(--color-primary)' }}>{player.name || `Player ${player.id}`}</h3>
          <div className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
            Actor {player.id} · Account {player.account_id || '—'} · Controller {player.controller_id || '—'}
          </div>
        </div>
        <span className="rounded px-2 py-1 text-xs uppercase" style={{ border: '1px solid #2a2418', color: player.online_status === 'online' ? 'var(--color-primary)' : 'var(--color-text-dim)' }}>
          {player.online_status || 'unknown'}
        </span>
      </div>

      <div className="grid grid-cols-2 gap-2 text-xs">
        <div><span style={{ color: 'var(--color-text-dim)' }}>Class</span><div style={{ color: 'var(--color-text)' }}>{fmt(player.class)}</div></div>
        <div><span style={{ color: 'var(--color-text-dim)' }}>Map</span><div style={{ color: 'var(--color-text)' }}>{fmt(player.map)}</div></div>
        <div><span style={{ color: 'var(--color-text-dim)' }}>Faction</span><div style={{ color: 'var(--color-text)' }}>{fmt(player.faction_id)}</div></div>
        <div><span style={{ color: 'var(--color-text-dim)' }}>FLS</span><div className="truncate" style={{ color: 'var(--color-text)' }}>{fmt(player.fls_id)}</div></div>
      </div>

      {profile && <PlayerProgressSummary profile={profile} />}

      <div className="flex gap-2 mt-auto">
        <Button size="sm" variant="outline" onPress={loadProgress} isDisabled={loading}>
          {loading ? <Spinner size="sm" color="current" /> : null}
          {profile ? 'Refresh Progress' : 'Load Progress'}
        </Button>
        <Button size="sm" variant="secondary" onPress={() => openPlayer360(player.id)}>
          Open 360
        </Button>
      </div>
    </article>
  )
}

export default function PlayerCardsTab() {
  const [players, setPlayers] = useState<Player[]>([])
  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | 'online' | 'offline'>('all')

  const loadPlayers = async () => {
    setLoading(true)
    try {
      setPlayers(await api.players.list())
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void loadPlayers() }, [])

  const filtered = useMemo(() => {
    const term = query.trim().toLowerCase()
    return players.filter(player => {
      const status = (player.online_status || '').toLowerCase()
      if (statusFilter !== 'all' && status !== statusFilter) return false
      if (!term) return true
      return [player.name, player.class, player.map, player.fls_id, String(player.id), String(player.account_id), String(player.controller_id)]
        .some(value => (value ?? '').toLowerCase().includes(term))
    })
  }, [players, query, statusFilter])

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 className="text-lg font-semibold" style={{ color: 'var(--color-primary)' }}>Player Cards</h2>
            <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
              Read-only progress cards for player support and self-service views. Load progress on demand to summarize level, journey, inventory, specs, factions, vehicles, and location.
            </p>
          </div>
          <div className="flex flex-wrap gap-2">
            <input
              value={query}
              onChange={event => setQuery(event.target.value)}
              placeholder="Search player, ID, map, class"
              className="rounded px-3 py-1.5 text-sm border"
              style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none', minWidth: 240 }}
            />
            <select
              value={statusFilter}
              onChange={event => setStatusFilter(event.target.value as 'all' | 'online' | 'offline')}
              className="rounded px-3 py-1.5 text-sm border"
              style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
            >
              <option value="all">All</option>
              <option value="online">Online</option>
              <option value="offline">Offline</option>
            </select>
            <Button variant="outline" size="sm" onPress={loadPlayers} isDisabled={loading}>
              {loading ? <Spinner size="sm" color="current" /> : null}
              Refresh
            </Button>
          </div>
        </div>
      </div>

      {loading && players.length === 0 ? <div className="flex justify-center py-12"><Spinner size="lg" /></div> : null}

      {!loading && filtered.length === 0 ? (
        <div className="rounded-lg p-6 text-sm" style={{ background: '#0d0b07', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
          No players match the current filters.
        </div>
      ) : (
        <div className="grid grid-cols-1 xl:grid-cols-2 2xl:grid-cols-3 gap-3">
          {filtered.map(player => <PlayerCard key={player.id} player={player} />)}
        </div>
      )}
    </div>
  )
}
