import { useEffect, useMemo, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { InventoryItem, Player } from '../api/client'

function itemLabel(item: InventoryItem): string {
  return item.name || item.template_id || `item:${item.id}`
}

function downloadJson(filename: string, data: unknown): void {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

function safeFilenamePart(value: string): string {
  return value.trim().replace(/[^a-zA-Z0-9._-]+/g, '-').replace(/^-+|-+$/g, '') || 'unknown'
}

export default function InventoryStudioTab() {
  const [players, setPlayers] = useState<Player[]>([])
  const [playersLoading, setPlayersLoading] = useState(false)
  const [playerSearch, setPlayerSearch] = useState('')
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>(null)
  const [items, setItems] = useState<InventoryItem[]>([])
  const [itemsLoading, setItemsLoading] = useState(false)
  const [itemSearch, setItemSearch] = useState('')
  const [selectedItemId, setSelectedItemId] = useState<number | null>(null)

  const loadPlayers = async () => {
    setPlayersLoading(true)
    try {
      setPlayers(await api.players.list())
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setPlayersLoading(false)
    }
  }

  const loadInventory = async (player: Player) => {
    setSelectedPlayer(player)
    setItems([])
    setSelectedItemId(null)
    setItemsLoading(true)
    try {
      const rows = await api.players.inventory(player.id)
      setItems(rows)
      setSelectedItemId(rows[0]?.id ?? null)
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setItemsLoading(false)
    }
  }

  useEffect(() => { loadPlayers() }, [])

  const filteredPlayers = useMemo(() => {
    const q = playerSearch.toLowerCase().trim()
    if (!q) return players
    return players.filter(player =>
      player.name.toLowerCase().includes(q) ||
      player.class.toLowerCase().includes(q) ||
      player.map.toLowerCase().includes(q) ||
      String(player.id).includes(q) ||
      String(player.account_id).includes(q) ||
      String(player.controller_id).includes(q)
    )
  }, [players, playerSearch])

  const filteredItems = useMemo(() => {
    const q = itemSearch.toLowerCase().trim()
    if (!q) return items
    return items.filter(item =>
      item.template_id.toLowerCase().includes(q) ||
      item.name.toLowerCase().includes(q) ||
      String(item.id).includes(q) ||
      String(item.quality).includes(q)
    )
  }, [items, itemSearch])

  const selectedItem = useMemo(() => items.find(item => item.id === selectedItemId) ?? null, [items, selectedItemId])

  const exportSnapshot = () => {
    if (!selectedPlayer) return
    downloadJson(`inventory-snapshot-${safeFilenamePart(selectedPlayer.name)}-${selectedPlayer.id}.json`, {
      exported_at: new Date().toISOString(),
      mode: 'read-only-inventory-studio-v2-foundation',
      player: selectedPlayer,
      item_count: items.length,
      items,
    })
    toast.success('Inventory snapshot exported')
  }

  return (
    <div className="flex flex-col gap-3 h-full min-h-0 overflow-hidden">
      <div className="rounded-lg px-4 py-2 text-xs" style={{ background: '#0f0d09', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
        <strong style={{ color: 'var(--color-primary)' }}>Inventory Studio v2 foundation:</strong> read-only player inventory inspection, filtering, selected-item detail, and snapshot export. Item edits will be added later as confirmed workflows with before/after snapshots.
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-[320px_1fr_360px] gap-3 flex-1 min-h-0 overflow-hidden">
        <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
          <div className="flex items-center justify-between px-3 py-2 shrink-0" style={{ borderBottom: '1px solid #2a2418' }}>
            <span className="text-xs font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>Players</span>
            <Button size="sm" variant="ghost" onPress={loadPlayers} isDisabled={playersLoading}>{playersLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button>
          </div>
          <div className="p-2 shrink-0">
            <input
              className="w-full rounded px-2 py-1.5 text-xs border"
              style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
              placeholder="Search player, actor, account, controller..."
              value={playerSearch}
              onChange={e => setPlayerSearch(e.target.value)}
            />
          </div>
          <div className="overflow-y-auto flex-1 min-h-0">
            {filteredPlayers.map(player => (
              <button
                key={player.id}
                type="button"
                onClick={() => loadInventory(player)}
                className="w-full text-left px-3 py-2 text-xs"
                style={{
                  borderBottom: '1px solid #1a1610',
                  borderLeft: selectedPlayer?.id === player.id ? '2px solid var(--color-primary)' : '2px solid transparent',
                  background: selectedPlayer?.id === player.id ? '#241e12' : 'transparent',
                  color: 'var(--color-text)',
                  cursor: 'pointer',
                }}
              >
                <div className="flex items-center justify-between gap-2">
                  <span className="font-semibold truncate">{player.name}</span>
                  <span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{player.id}</span>
                </div>
                <div className="mt-0.5 truncate" style={{ color: 'var(--color-text-dim)' }}>
                  {player.class || 'Unknown'} · {player.map || 'unknown'} · {player.online_status || 'Offline'}
                </div>
              </button>
            ))}
            {filteredPlayers.length === 0 && (
              <div className="px-3 py-8 text-center text-xs" style={{ color: 'var(--color-text-dim)' }}>No matching players</div>
            )}
          </div>
        </section>

        <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
          <div className="flex items-center justify-between gap-3 px-3 py-2 shrink-0" style={{ borderBottom: '1px solid #2a2418' }}>
            <div>
              <div className="text-xs font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>Inventory</div>
              <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>
                {selectedPlayer ? `${selectedPlayer.name} · ${items.length} item row(s)` : 'Select a player'}
              </div>
            </div>
            <div className="flex gap-2">
              <Button size="sm" variant="ghost" isDisabled={!selectedPlayer || itemsLoading} onPress={() => selectedPlayer && loadInventory(selectedPlayer)}>{itemsLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button>
              <Button size="sm" variant="outline" isDisabled={!selectedPlayer || items.length === 0} onPress={exportSnapshot}>Export Snapshot</Button>
            </div>
          </div>
          <div className="p-2 shrink-0">
            <input
              className="w-full rounded px-2 py-1.5 text-xs border"
              style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
              placeholder="Search item template, name, ID, quality..."
              value={itemSearch}
              onChange={e => setItemSearch(e.target.value)}
              disabled={!selectedPlayer}
            />
          </div>
          {itemsLoading ? (
            <div className="flex justify-center py-12"><Spinner size="lg" /></div>
          ) : !selectedPlayer ? (
            <div className="flex items-center justify-center flex-1 text-sm" style={{ color: 'var(--color-text-dim)' }}>Select a player to load inventory.</div>
          ) : filteredItems.length === 0 ? (
            <div className="flex items-center justify-center flex-1 text-sm" style={{ color: 'var(--color-text-dim)' }}>{items.length === 0 ? 'No inventory rows found.' : 'No matching items.'}</div>
          ) : (
            <div className="overflow-auto flex-1 min-h-0">
              <table className="w-full text-xs">
                <thead>
                  <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418', position: 'sticky', top: 0 }}>
                    {['Item', 'Stack', 'Quality', 'Durability'].map(header => (
                      <th key={header} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{header}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {filteredItems.map((item, index) => (
                    <tr
                      key={item.id}
                      onClick={() => setSelectedItemId(item.id)}
                      style={{
                        borderBottom: '1px solid #1a1610',
                        background: selectedItemId === item.id ? '#241e12' : index % 2 === 0 ? '#0d0b07' : '#0f0d09',
                        cursor: 'pointer',
                      }}
                    >
                      <td className="px-3 py-2">
                        <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{itemLabel(item)}</div>
                        <div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>{item.template_id}</div>
                      </td>
                      <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.stack_size}</td>
                      <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.quality}</td>
                      <td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.durability} / {item.max_durability}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

        <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
          <div className="px-3 py-2 shrink-0" style={{ borderBottom: '1px solid #2a2418' }}>
            <div className="text-xs font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>Selected Item</div>
            <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Read-only details</div>
          </div>
          {!selectedItem ? (
            <div className="flex items-center justify-center flex-1 text-sm p-4 text-center" style={{ color: 'var(--color-text-dim)' }}>
              Select an inventory row to inspect item details.
            </div>
          ) : (
            <div className="overflow-y-auto flex-1 p-3 text-xs min-h-0">
              <div className="font-semibold text-sm" style={{ color: 'var(--color-text)' }}>{itemLabel(selectedItem)}</div>
              <div className="font-mono mt-1" style={{ color: 'var(--color-text-dim)' }}>{selectedItem.template_id}</div>
              <div className="grid grid-cols-2 gap-2 mt-4">
                <Detail label="Item ID" value={String(selectedItem.id)} />
                <Detail label="Stack" value={String(selectedItem.stack_size)} />
                <Detail label="Quality" value={String(selectedItem.quality)} />
                <Detail label="Durability" value={`${selectedItem.durability} / ${selectedItem.max_durability}`} />
              </div>
              <div className="rounded p-3 mt-4" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                <div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Next editing phase</div>
                <p>Future edits will require a before snapshot, after preview, shared mutation confirmation, admin reason capture, and audit visibility before any write is sent.</p>
              </div>
              <pre className="rounded p-3 mt-4 overflow-auto" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                {JSON.stringify(selectedItem, null, 2)}
              </pre>
            </div>
          )}
        </section>
      </div>
    </div>
  )
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded p-2" style={{ background: '#0a0806', border: '1px solid #2a2418' }}>
      <div className="text-[10px] uppercase tracking-wide" style={{ color: 'var(--color-text-dim)' }}>{label}</div>
      <div className="font-mono mt-1" style={{ color: 'var(--color-text)' }}>{value || '—'}</div>
    </div>
  )
}
