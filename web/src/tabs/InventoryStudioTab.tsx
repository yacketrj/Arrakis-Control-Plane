import { useEffect, useMemo, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { InventoryItem, ItemTemplate, Player } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

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

type InventorySnapshot = {
  exported_at?: string
  mode?: string
  player?: Player
  item_count?: number
  items?: InventoryItem[]
}

type InventoryDiffRow = {
  key: string
  status: 'added' | 'removed' | 'changed'
  before?: InventoryItem
  after?: InventoryItem
  changes: string[]
}

function snapshotItems(snapshot: InventorySnapshot): InventoryItem[] {
  return Array.isArray(snapshot.items) ? snapshot.items : []
}

function itemComparisonKey(item: InventoryItem): string {
  return String(item.id)
}

function compareInventorySnapshots(before: InventoryItem[], after: InventoryItem[]): InventoryDiffRow[] {
  const beforeMap = new Map(before.map(item => [itemComparisonKey(item), item]))
  const afterMap = new Map(after.map(item => [itemComparisonKey(item), item]))
  const keys = Array.from(new Set([...beforeMap.keys(), ...afterMap.keys()])).sort((a, b) => Number(a) - Number(b))
  const diffs: InventoryDiffRow[] = []

  for (const key of keys) {
    const oldItem = beforeMap.get(key)
    const newItem = afterMap.get(key)
    if (!oldItem && newItem) {
      diffs.push({ key, status: 'added', after: newItem, changes: ['Item row added'] })
      continue
    }
    if (oldItem && !newItem) {
      diffs.push({ key, status: 'removed', before: oldItem, changes: ['Item row removed'] })
      continue
    }
    if (!oldItem || !newItem) continue

    const changes: string[] = []
    if (oldItem.template_id !== newItem.template_id) changes.push(`Template: ${oldItem.template_id} → ${newItem.template_id}`)
    if (oldItem.name !== newItem.name) changes.push(`Name: ${oldItem.name || '—'} → ${newItem.name || '—'}`)
    if (oldItem.stack_size !== newItem.stack_size) changes.push(`Stack: ${oldItem.stack_size} → ${newItem.stack_size}`)
    if (oldItem.quality !== newItem.quality) changes.push(`Quality: ${oldItem.quality} → ${newItem.quality}`)
    if (oldItem.durability !== newItem.durability) changes.push(`Durability: ${oldItem.durability} → ${newItem.durability}`)
    if (oldItem.max_durability !== newItem.max_durability) changes.push(`Max durability: ${oldItem.max_durability} → ${newItem.max_durability}`)
    if (changes.length > 0) diffs.push({ key, status: 'changed', before: oldItem, after: newItem, changes })
  }

  return diffs
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
  const [comparisonSnapshot, setComparisonSnapshot] = useState<InventorySnapshot | null>(null)
  const [comparisonName, setComparisonName] = useState('')
  const [templates, setTemplates] = useState<ItemTemplate[]>([])
  const [templatesLoading, setTemplatesLoading] = useState(false)
  const [templateSearch, setTemplateSearch] = useState('')
  const [selectedTemplateId, setSelectedTemplateId] = useState('')
  const [mutationBusy, setMutationBusy] = useState(false)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

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

  const loadTemplates = async () => {
    setTemplatesLoading(true)
    try {
      setTemplates(await api.players.templates())
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setTemplatesLoading(false)
    }
  }

  const loadInventory = async (player: Player) => {
    setSelectedPlayer(player)
    setItems([])
    setSelectedItemId(null)
    setComparisonSnapshot(null)
    setComparisonName('')
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

  useEffect(() => { loadPlayers(); loadTemplates() }, [])

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

  const filteredTemplates = useMemo(() => {
    const q = templateSearch.toLowerCase().trim()
    if (!q) return templates.slice(0, 80)
    return templates.filter(template =>
      template.id.toLowerCase().includes(q) ||
      template.name.toLowerCase().includes(q)
    ).slice(0, 120)
  }, [templates, templateSearch])

  const selectedItem = useMemo(() => items.find(item => item.id === selectedItemId) ?? null, [items, selectedItemId])
  const selectedTemplate = useMemo(() => templates.find(template => template.id === selectedTemplateId) ?? null, [selectedTemplateId, templates])
  const comparisonItems = useMemo(() => comparisonSnapshot ? snapshotItems(comparisonSnapshot) : [], [comparisonSnapshot])
  const comparisonDiffs = useMemo(() => comparisonSnapshot ? compareInventorySnapshots(comparisonItems, items) : [], [comparisonItems, comparisonSnapshot, items])

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

  const exportBeforeActionSnapshot = (action: string, item: InventoryItem) => {
    if (!selectedPlayer) return
    downloadJson(`inventory-before-${action}-${safeFilenamePart(selectedPlayer.name)}-${selectedPlayer.id}-item-${item.id}.json`, {
      exported_at: new Date().toISOString(),
      mode: 'inventory-studio-before-action-snapshot',
      action,
      player: selectedPlayer,
      target_item: item,
      item_count: items.length,
      items,
    })
  }

  const repairSelectedItem = async () => {
    if (!selectedPlayer || !selectedItem) return

    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: '/api/v1/players/repair-item',
        title: 'Repair selected inventory item',
        summary: `Repair ${itemLabel(selectedItem)} for ${selectedPlayer.name}.`,
        target: `actor:${selectedPlayer.id} · item:${selectedItem.id}`,
        details: [
          `Player: ${selectedPlayer.name}`,
          `Online state: ${selectedPlayer.online_status || 'Offline'}`,
          `Item ID: ${selectedItem.id}`,
          `Template: ${selectedItem.template_id}`,
          `Current durability: ${selectedItem.durability} / ${selectedItem.max_durability}`,
          'A before-action inventory snapshot will be downloaded before the repair request is sent.',
        ],
        confirmLabel: 'Repair item',
        forceReason: true,
      })
    } catch (e: unknown) {
      if (e instanceof Error && e.message === mutationConfirmationCancelledMessage) {
        toast.warning('Cancelled')
        return
      }
      toast.danger(e instanceof Error ? e.message : String(e))
      return
    }

    if (!reason) {
      toast.warning('Cancelled: admin reason is required for this action')
      return
    }

    setMutationBusy(true)
    try {
      exportBeforeActionSnapshot('repair', selectedItem)
      await api.players.repairItem(selectedItem.id, reason)
      toast.success(`Repaired ${itemLabel(selectedItem)}`)
      await loadInventory(selectedPlayer)
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setMutationBusy(false)
    }
  }

  const loadComparisonSnapshot = async (file: File | null) => {
    if (!file) return
    try {
      const text = await file.text()
      const parsed = JSON.parse(text) as InventorySnapshot
      if (!Array.isArray(parsed.items)) throw new Error('Snapshot does not contain an items array')
      setComparisonSnapshot(parsed)
      setComparisonName(file.name)
      toast.success(`Loaded comparison snapshot with ${parsed.items.length} item row(s)`)
    } catch (e: unknown) {
      setComparisonSnapshot(null)
      setComparisonName('')
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  return (
    <>
      <div className="flex flex-col gap-3 h-full min-h-0 overflow-hidden">
        <div className="rounded-lg px-4 py-2 text-xs" style={{ background: '#0f0d09', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
          <strong style={{ color: 'var(--color-primary)' }}>Inventory Studio v2:</strong> player inventory inspection, filtering, selected-item detail, snapshot export, local snapshot comparison, item catalog browsing, and confirmed selected-item repair with before-action snapshot.
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-[320px_1fr_420px] gap-3 flex-1 min-h-0 overflow-hidden">
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
              <div className="text-xs font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>Inspector</div>
              <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Selected item, catalog, and snapshot diff</div>
            </div>
            <div className="overflow-y-auto flex-1 p-3 text-xs min-h-0">
              <div className="rounded p-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}>
                <div className="flex items-center justify-between gap-2 mb-2">
                  <div className="font-semibold" style={{ color: 'var(--color-primary)' }}>Item Catalog</div>
                  <Button size="sm" variant="ghost" onPress={loadTemplates} isDisabled={templatesLoading}>{templatesLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button>
                </div>
                <input
                  className="w-full rounded px-2 py-1.5 text-xs border"
                  style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                  placeholder="Search catalog template ID or name..."
                  value={templateSearch}
                  onChange={event => setTemplateSearch(event.target.value)}
                />
                <div className="mt-2 rounded overflow-y-auto" style={{ border: '1px solid #2a2418', maxHeight: 180 }}>
                  {templatesLoading ? (
                    <div className="flex justify-center py-4"><Spinner size="sm" /></div>
                  ) : filteredTemplates.length === 0 ? (
                    <div className="px-3 py-4 text-center" style={{ color: 'var(--color-text-dim)' }}>No matching templates</div>
                  ) : (
                    filteredTemplates.map(template => (
                      <button
                        key={template.id}
                        type="button"
                        onClick={() => setSelectedTemplateId(template.id)}
                        className="w-full text-left px-2 py-1.5"
                        style={{ borderBottom: '1px solid #1a1610', background: selectedTemplateId === template.id ? '#241e12' : 'transparent', color: 'var(--color-text)', cursor: 'pointer' }}
                      >
                        <div className="font-mono truncate" style={{ color: selectedTemplateId === template.id ? 'var(--color-primary)' : 'var(--color-text)' }}>{template.id}</div>
                        {template.name && <div className="truncate" style={{ color: 'var(--color-text-dim)' }}>{template.name}</div>}
                      </button>
                    ))
                  )}
                </div>
                {selectedTemplate && (
                  <div className="mt-2 grid grid-cols-1 gap-2">
                    <Detail label="Selected Template" value={selectedTemplate.id} />
                    <Detail label="Catalog Name" value={selectedTemplate.name || '—'} />
                  </div>
                )}
              </div>

              <div className="rounded p-3 mt-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}>
                <div className="font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Snapshot Compare</div>
                <input
                  type="file"
                  accept=".json,application/json"
                  className="text-xs"
                  style={{ color: 'var(--color-text-dim)' }}
                  disabled={!selectedPlayer || items.length === 0}
                  onChange={event => loadComparisonSnapshot(event.target.files?.[0] ?? null)}
                />
                {comparisonName && <div className="mt-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{comparisonName}</div>}
                {comparisonSnapshot && (
                  <div className="grid grid-cols-3 gap-2 mt-3">
                    <Detail label="Before" value={String(comparisonItems.length)} />
                    <Detail label="Current" value={String(items.length)} />
                    <Detail label="Diffs" value={String(comparisonDiffs.length)} />
                  </div>
                )}
                {comparisonSnapshot && comparisonDiffs.length === 0 && <div className="mt-3" style={{ color: 'var(--color-success)' }}>No differences detected against loaded snapshot.</div>}
                {comparisonDiffs.length > 0 && (
                  <div className="mt-3 flex flex-col gap-2">
                    {comparisonDiffs.slice(0, 25).map(diff => (
                      <div key={`${diff.status}-${diff.key}`} className="rounded p-2" style={{ background: '#0f0d09', border: '1px solid #2a2418' }}>
                        <div className="flex items-center justify-between gap-2">
                          <span className="font-semibold" style={{ color: diff.status === 'removed' ? '#e88' : diff.status === 'added' ? 'var(--color-success)' : '#f0a830' }}>{diff.status.toUpperCase()}</span>
                          <span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>item:{diff.key}</span>
                        </div>
                        <div className="mt-1" style={{ color: 'var(--color-text)' }}>{itemLabel(diff.after ?? diff.before as InventoryItem)}</div>
                        <ul className="mt-1 pl-4" style={{ color: 'var(--color-text-dim)', listStyle: 'disc' }}>
                          {diff.changes.map(change => <li key={change}>{change}</li>)}
                        </ul>
                      </div>
                    ))}
                    {comparisonDiffs.length > 25 && <div style={{ color: 'var(--color-text-dim)' }}>Showing first 25 of {comparisonDiffs.length} differences.</div>}
                  </div>
                )}
              </div>

              <div className="mt-3">
                {!selectedItem ? (
                  <div className="flex items-center justify-center text-sm p-4 text-center" style={{ color: 'var(--color-text-dim)' }}>
                    Select an inventory row to inspect item details.
                  </div>
                ) : (
                  <>
                    <div className="font-semibold text-sm" style={{ color: 'var(--color-text)' }}>{itemLabel(selectedItem)}</div>
                    <div className="font-mono mt-1" style={{ color: 'var(--color-text-dim)' }}>{selectedItem.template_id}</div>
                    <div className="grid grid-cols-2 gap-2 mt-4">
                      <Detail label="Item ID" value={String(selectedItem.id)} />
                      <Detail label="Stack" value={String(selectedItem.stack_size)} />
                      <Detail label="Quality" value={String(selectedItem.quality)} />
                      <Detail label="Durability" value={`${selectedItem.durability} / ${selectedItem.max_durability}`} />
                    </div>
                    <div className="rounded p-3 mt-4" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                      <div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Confirmed item action</div>
                      <p>Repair selected item downloads a before-action inventory snapshot, then requires shared mutation confirmation and an admin reason.</p>
                      <Button size="sm" className="mt-3" onPress={repairSelectedItem} isDisabled={mutationBusy}>
                        {mutationBusy ? <Spinner size="sm" color="current" /> : null}
                        Repair Selected Item
                      </Button>
                    </div>
                    <pre className="rounded p-3 mt-4 overflow-auto" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                      {JSON.stringify(selectedItem, null, 2)}
                    </pre>
                  </>
                )}
              </div>
            </div>
          </section>
        </div>
      </div>
      {confirmationDialog}
    </>
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
