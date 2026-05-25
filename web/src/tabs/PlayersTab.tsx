import { useState, useEffect, useMemo } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import GiveItemModalAugmented from './GiveItemModalAugmented'
import type {
  Player,
  CurrencyRow,
  FactionRep,
  SpecTrack,
  OnlineRow,
} from '../api/client'

type Sidebar = 'players' | 'currency' | 'factions' | 'specs' | 'online'

function StatusDot({ status }: { status: string }) {
  const color = status === 'Online' ? '#27ae60' : status === 'LoggingOut' ? '#f0a830' : '#555'
  return (
    <span
      style={{
        display: 'inline-block',
        width: 8,
        height: 8,
        borderRadius: '50%',
        background: color,
        marginRight: 6,
        flexShrink: 0,
      }}
    />
  )
}

function OnlineBadge({ status }: { status: string }) {
  const color = status === 'Online' ? '#27ae60' : status === 'LoggingOut' ? '#f0a830' : '#555'
  const label = status === 'Online' ? 'Online' : status === 'LoggingOut' ? 'LoggingOut' : status || 'Offline'
  return (
    <span className="text-xs px-1.5 py-0.5 rounded font-semibold" style={{ background: color + '22', color, border: `1px solid ${color}44` }}>
      {label}
    </span>
  )
}

export default function PlayersTab() {
  const [active, setActive] = useState<Sidebar>('players')
  const [players, setPlayers] = useState<Player[]>([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>(null)
  const [showGiveItem, setShowGiveItem] = useState(false)
  const [sideLoading, setSideLoading] = useState(false)

  const [currencyData, setCurrencyData] = useState<CurrencyRow[]>([])
  const [factionData, setFactionData] = useState<FactionRep[]>([])
  const [specData, setSpecData] = useState<SpecTrack[]>([])
  const [onlineData, setOnlineData] = useState<OnlineRow[]>([])
  const [sideSearch, setSideSearch] = useState('')

  useEffect(() => { loadPlayers() }, [])

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

  const loadSideData = async (section: Sidebar) => {
    setActive(section)
    setSideSearch('')
    if (section === 'players') return
    setSideLoading(true)
    try {
      if (section === 'online') {
        setOnlineData(await api.players.online())
      } else if (section === 'currency') {
        setCurrencyData(await api.players.currency())
      } else if (section === 'factions') {
        setFactionData(await api.players.factions())
      } else if (section === 'specs') {
        setSpecData(await api.players.specs())
      }
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setSideLoading(false)
    }
  }

  const filtered = useMemo(() => {
    const q = search.toLowerCase()
    return players.filter(p =>
      p.name.toLowerCase().includes(q) || p.class.toLowerCase().includes(q) ||
      p.map.toLowerCase().includes(q) || String(p.id).includes(q)
    )
  }, [players, search])

  const controllerToName = useMemo(() => {
    const m = new Map<number, string>()
    for (const p of players) m.set(p.controller_id, p.name)
    return m
  }, [players])

  const filteredCurrency = useMemo(() => {
    if (!sideSearch) return currencyData
    const q = sideSearch.toLowerCase()
    return currencyData.filter(r => {
      const name = controllerToName.get(r.player_id) ?? ''
      return name.toLowerCase().includes(q) || String(r.player_id).includes(q)
    })
  }, [currencyData, sideSearch, controllerToName])

  const filteredFactions = useMemo(() => {
    if (!sideSearch) return factionData
    const q = sideSearch.toLowerCase()
    return factionData.filter(r => {
      const name = controllerToName.get(r.actor_id) ?? ''
      return name.toLowerCase().includes(q) || String(r.actor_id).includes(q)
    })
  }, [factionData, sideSearch, controllerToName])

  const filteredSpecs = useMemo(() => {
    if (!sideSearch) return specData
    const q = sideSearch.toLowerCase()
    return specData.filter(r => {
      const name = controllerToName.get(r.player_id) ?? ''
      return name.toLowerCase().includes(q) || String(r.player_id).includes(q)
    })
  }, [specData, sideSearch, controllerToName])

  const filteredOnline = useMemo(() => {
    if (!sideSearch) return onlineData
    const q = sideSearch.toLowerCase()
    return onlineData.filter(r =>
      r.name.toLowerCase().includes(q) || String(r.player_id).includes(q)
    )
  }, [onlineData, sideSearch])

  const sidebarItems: { key: Sidebar; label: string }[] = [
    { key: 'players', label: 'Players' },
    { key: 'online', label: 'Online State' },
    { key: 'currency', label: 'Currency' },
    { key: 'factions', label: 'Factions' },
    { key: 'specs', label: 'Specs / XP' },
  ]

  const tableHeader = (cols: string[]) => (
    <thead>
      <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
        {cols.map(h => (
          <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>
        ))}
      </tr>
    </thead>
  )

  return (
    <div className="flex gap-4 h-full">
      <div className="w-40 shrink-0 flex flex-col gap-1 rounded-lg p-2" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        {sidebarItems.map(item => (
          <button
            key={item.key}
            onClick={() => loadSideData(item.key)}
            className="text-left px-3 py-2 rounded text-sm transition-colors"
            style={{ background: active === item.key ? 'var(--color-primary)' : 'transparent', color: active === item.key ? '#fff' : 'var(--color-text)' }}
          >
            {item.label}
          </button>
        ))}
      </div>

      <div className="flex-1 overflow-auto flex flex-col gap-4">
        {active === 'players' && (
          <>
            <div className="flex items-center gap-3">
              <input
                className="rounded px-3 py-1.5 text-sm border w-72"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                placeholder="Search players..."
                value={search}
                onChange={e => setSearch(e.target.value)}
              />
              <Button variant="outline" size="sm" onPress={loadPlayers} isDisabled={loading}>
                {loading ? <Spinner size="sm" color="current" /> : null}
                Refresh
              </Button>
            </div>

            {loading ? (
              <div className="flex justify-center py-12"><Spinner size="lg" /></div>
            ) : (
              <div className="overflow-auto rounded-lg" style={{ border: '1px solid #2a2418' }}>
                <table className="w-full text-xs">
                  <thead>
                    <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
                      {['ID', 'Name', 'Class', 'Map', 'Faction', 'Actions'].map(h => (
                        <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {filtered.map((player, i) => (
                      <tr key={player.id} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                        <td className="px-3 py-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{player.id}</td>
                        <td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>
                          <div className="flex items-center">
                            <StatusDot status={player.online_status} />
                            {player.name}
                          </div>
                        </td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{player.class}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{player.map}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{player.faction_id || '—'}</td>
                        <td className="px-3 py-2">
                          <div className="flex gap-1 flex-wrap">
                            <Button size="sm" variant="ghost">Inventory</Button>
                            <Button size="sm" variant="ghost" onPress={() => { setSelectedPlayer(player); setShowGiveItem(true) }}>Give Item</Button>
                            <Button size="sm" variant="ghost">Actions</Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </>
        )}

        {active === 'currency' && (
          <div className="flex flex-col gap-3 h-full min-h-0">
            <div className="flex items-center gap-2 shrink-0">
              <h3 className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Currency</h3>
              <input
                className="flex-1 rounded px-3 py-1 text-sm border"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                placeholder="Search player..."
                value={sideSearch}
                onChange={e => setSideSearch(e.target.value)}
              />
            </div>
            {sideLoading ? (
              <div className="flex justify-center py-12"><Spinner size="lg" /></div>
            ) : (
              <div className="overflow-auto rounded-lg flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
                <table className="w-full text-xs">
                  {tableHeader(['Player', 'Currency', 'Balance'])}
                  <tbody>
                    {filteredCurrency.map((row, i) => (
                      <tr key={`${row.player_id}-${row.currency_id}`} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                        <td className="px-3 py-2">
                          {controllerToName.get(row.player_id) && <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{controllerToName.get(row.player_id)}</div>}
                          <div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{row.player_id}</div>
                        </td>
                        <td className="px-3 py-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{row.currency_id}</td>
                        <td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>{row.balance.toLocaleString()}</td>
                      </tr>
                    ))}
                    {filteredCurrency.length === 0 && (
                      <tr><td colSpan={3} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>{sideSearch ? 'No matches' : 'No data'}</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {active === 'factions' && (
          <div className="flex flex-col gap-3 h-full min-h-0">
            <div className="flex items-center gap-2 shrink-0">
              <h3 className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Factions</h3>
              <input
                className="flex-1 rounded px-3 py-1 text-sm border"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                placeholder="Search player..."
                value={sideSearch}
                onChange={e => setSideSearch(e.target.value)}
              />
            </div>
            {sideLoading ? (
              <div className="flex justify-center py-12"><Spinner size="lg" /></div>
            ) : (
              <div className="overflow-auto rounded-lg flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
                <table className="w-full text-xs">
                  {tableHeader(['Player', 'Faction', 'Reputation', 'Scrips'])}
                  <tbody>
                    {filteredFactions.map((row, i) => (
                      <tr key={`${row.actor_id}-${row.faction_id}`} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                        <td className="px-3 py-2">
                          {controllerToName.get(row.actor_id) && <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{controllerToName.get(row.actor_id)}</div>}
                          <div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{row.actor_id}</div>
                        </td>
                        <td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>{row.faction_name}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.reputation.toLocaleString()}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.scrips.toLocaleString()}</td>
                      </tr>
                    ))}
                    {filteredFactions.length === 0 && (
                      <tr><td colSpan={4} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>{sideSearch ? 'No matches' : 'No data'}</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {active === 'specs' && (
          <div className="flex flex-col gap-3 h-full min-h-0">
            <div className="flex items-center gap-2 shrink-0">
              <h3 className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Specs / XP</h3>
              <input
                className="flex-1 rounded px-3 py-1 text-sm border"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                placeholder="Search player..."
                value={sideSearch}
                onChange={e => setSideSearch(e.target.value)}
              />
            </div>
            {sideLoading ? (
              <div className="flex justify-center py-12"><Spinner size="lg" /></div>
            ) : (
              <div className="overflow-auto rounded-lg flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
                <table className="w-full text-xs">
                  {tableHeader(['Player', 'Track', 'XP', 'Level'])}
                  <tbody>
                    {filteredSpecs.map((row, i) => (
                      <tr key={`${row.player_id}-${row.track_type}`} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                        <td className="px-3 py-2">
                          {controllerToName.get(row.player_id) && <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{controllerToName.get(row.player_id)}</div>}
                          <div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{row.player_id}</div>
                        </td>
                        <td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>{row.track_type}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.xp.toLocaleString()}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.level}</td>
                      </tr>
                    ))}
                    {filteredSpecs.length === 0 && (
                      <tr><td colSpan={4} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>{sideSearch ? 'No matches' : 'No data'}</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {active === 'online' && (
          <div className="flex flex-col gap-3 h-full min-h-0">
            <div className="flex items-center gap-2 shrink-0">
              <h3 className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Online State</h3>
              <input
                className="flex-1 rounded px-3 py-1 text-sm border"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                placeholder="Search player..."
                value={sideSearch}
                onChange={e => setSideSearch(e.target.value)}
              />
            </div>
            {sideLoading ? (
              <div className="flex justify-center py-12"><Spinner size="lg" /></div>
            ) : (
              <div className="overflow-auto rounded-lg flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
                <table className="w-full text-xs">
                  {tableHeader(['Player', 'Status', 'Last Seen', 'Map'])}
                  <tbody>
                    {filteredOnline.map((row, i) => (
                      <tr key={row.player_id} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                        <td className="px-3 py-2">
                          <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{row.name}</div>
                          <div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{row.player_id}</div>
                        </td>
                        <td className="px-3 py-2"><OnlineBadge status={row.status} /></td>
                        <td className="px-3 py-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{row.last_seen}</td>
                        <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.map}</td>
                      </tr>
                    ))}
                    {filteredOnline.length === 0 && (
                      <tr><td colSpan={4} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>{sideSearch ? 'No matches' : 'No data'}</td></tr>
                    )}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {selectedPlayer && (
        <GiveItemModalAugmented player={selectedPlayer} open={showGiveItem} onClose={() => setShowGiveItem(false)} />
      )}
    </div>
  )
}
