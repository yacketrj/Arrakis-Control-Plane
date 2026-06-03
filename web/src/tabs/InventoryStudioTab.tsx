import { useEffect, useMemo, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { InventoryItem, ItemTemplate, Player } from '../api/client'
import { inventoryStudioMutations } from '../api/inventoryStudioMutations'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

type InventorySnapshot = { exported_at?: string; mode?: string; player?: Player; item_count?: number; items?: InventoryItem[] }
type InventoryDiffRow = { key: string; status: 'added' | 'removed' | 'changed'; before?: InventoryItem; after?: InventoryItem; changes: string[] }
type ActionDiff = { action: string; target: string; beforeCount: number; afterCount: number; diffs: InventoryDiffRow[]; checkedAt: string }
type LastActionDiff = ActionDiff | null

function itemLabel(item: InventoryItem): string { return item.name || item.template_id || `item:${item.id}` }
function safeFilenamePart(value: string): string { return value.trim().replace(/[^a-zA-Z0-9._-]+/g, '-').replace(/^-+|-+$/g, '') || 'unknown' }
function clampInt(value: number, min: number, max: number): number { return Number.isFinite(value) ? Math.max(min, Math.min(max, Math.trunc(value))) : min }
function snapshotItems(snapshot: InventorySnapshot): InventoryItem[] { return Array.isArray(snapshot.items) ? snapshot.items : [] }

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

function compareInventorySnapshots(before: InventoryItem[], after: InventoryItem[]): InventoryDiffRow[] {
  const beforeMap = new Map(before.map(item => [String(item.id), item]))
  const afterMap = new Map(after.map(item => [String(item.id), item]))
  const keys = Array.from(new Set([...beforeMap.keys(), ...afterMap.keys()])).sort((a, b) => Number(a) - Number(b))
  const diffs: InventoryDiffRow[] = []

  for (const key of keys) {
    const oldItem = beforeMap.get(key)
    const newItem = afterMap.get(key)
    if (!oldItem && newItem) { diffs.push({ key, status: 'added', after: newItem, changes: ['Item row added'] }); continue }
    if (oldItem && !newItem) { diffs.push({ key, status: 'removed', before: oldItem, changes: ['Item row removed'] }); continue }
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
  const [templates, setTemplates] = useState<ItemTemplate[]>([])
  const [templatesLoading, setTemplatesLoading] = useState(false)
  const [templateSearch, setTemplateSearch] = useState('')
  const [selectedTemplateId, setSelectedTemplateId] = useState('')
  const [addQty, setAddQty] = useState(1)
  const [addQuality, setAddQuality] = useState(1)
  const [stackSizeDraft, setStackSizeDraft] = useState(1)
  const [comparisonSnapshot, setComparisonSnapshot] = useState<InventorySnapshot | null>(null)
  const [comparisonName, setComparisonName] = useState('')
  const [lastActionDiff, setLastActionDiff] = useState<LastActionDiff>(null)
  const [actionDiffHistory, setActionDiffHistory] = useState<ActionDiff[]>([])
  const [mutationBusy, setMutationBusy] = useState(false)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  const selectedItem = useMemo(() => items.find(item => item.id === selectedItemId) ?? null, [items, selectedItemId])
  const selectedTemplate = useMemo(() => templates.find(template => template.id === selectedTemplateId) ?? null, [selectedTemplateId, templates])
  const comparisonItems = useMemo(() => comparisonSnapshot ? snapshotItems(comparisonSnapshot) : [], [comparisonSnapshot])
  const comparisonDiffs = useMemo(() => comparisonSnapshot ? compareInventorySnapshots(comparisonItems, items) : [], [comparisonItems, comparisonSnapshot, items])

  useEffect(() => { if (selectedItem) setStackSizeDraft(clampInt(selectedItem.stack_size, 1, 9999)) }, [selectedItem])

  const filteredPlayers = useMemo(() => {
    const q = playerSearch.toLowerCase().trim()
    if (!q) return players
    return players.filter(p => p.name.toLowerCase().includes(q) || p.class.toLowerCase().includes(q) || p.map.toLowerCase().includes(q) || String(p.id).includes(q) || String(p.account_id).includes(q) || String(p.controller_id).includes(q))
  }, [players, playerSearch])

  const filteredItems = useMemo(() => {
    const q = itemSearch.toLowerCase().trim()
    if (!q) return items
    return items.filter(i => i.template_id.toLowerCase().includes(q) || i.name.toLowerCase().includes(q) || String(i.id).includes(q) || String(i.quality).includes(q))
  }, [items, itemSearch])

  const filteredTemplates = useMemo(() => {
    const q = templateSearch.toLowerCase().trim()
    const rows = q ? templates.filter(t => t.id.toLowerCase().includes(q) || t.name.toLowerCase().includes(q)) : templates
    return rows.slice(0, q ? 120 : 80)
  }, [templates, templateSearch])

  const loadPlayers = async () => {
    setPlayersLoading(true)
    try { setPlayers(await api.players.list()) } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setPlayersLoading(false) }
  }

  const loadTemplates = async () => {
    setTemplatesLoading(true)
    try { setTemplates(await api.players.templates()) } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setTemplatesLoading(false) }
  }

  const reloadInventory = async (player: Player, reset = false): Promise<InventoryItem[]> => {
    if (reset) {
      setSelectedPlayer(player)
      setItems([])
      setSelectedItemId(null)
      setComparisonSnapshot(null)
      setComparisonName('')
      setLastActionDiff(null)
      setActionDiffHistory([])
    }
    setItemsLoading(true)
    try {
      const rows = await api.players.inventory(player.id)
      setItems(rows)
      setSelectedItemId(prev => rows.some(item => item.id === prev) ? prev : rows[0]?.id ?? null)
      return rows
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
      return []
    } finally { setItemsLoading(false) }
  }

  useEffect(() => { loadPlayers(); loadTemplates() }, [])

  const exportSnapshot = () => {
    if (!selectedPlayer) return
    downloadJson(`inventory-snapshot-${safeFilenamePart(selectedPlayer.name)}-${selectedPlayer.id}.json`, { exported_at: new Date().toISOString(), mode: 'inventory-studio-snapshot', player: selectedPlayer, item_count: items.length, items })
    toast.success('Inventory snapshot exported')
  }

  const exportActionHistory = () => {
    if (!selectedPlayer || actionDiffHistory.length === 0) return
    downloadJson(`inventory-action-history-${safeFilenamePart(selectedPlayer.name)}-${selectedPlayer.id}-${safeFilenamePart(new Date().toISOString())}.json`, { exported_at: new Date().toISOString(), mode: 'inventory-studio-action-history', player: selectedPlayer, action_count: actionDiffHistory.length, actions: actionDiffHistory })
    toast.success('Inventory action history exported')
  }

  const exportBeforeActionSnapshot = (action: string, beforeItems: InventoryItem[], details: Record<string, unknown>) => {
    if (!selectedPlayer) return
    downloadJson(`inventory-before-${action}-${safeFilenamePart(selectedPlayer.name)}-${selectedPlayer.id}-${safeFilenamePart(new Date().toISOString())}.json`, { exported_at: new Date().toISOString(), mode: 'inventory-studio-before-action-snapshot', action, player: selectedPlayer, ...details, item_count: beforeItems.length, items: beforeItems })
  }

  const setPostActionDiff = (action: string, target: string, beforeItems: InventoryItem[], afterItems: InventoryItem[]) => {
    const next: ActionDiff = { action, target, beforeCount: beforeItems.length, afterCount: afterItems.length, diffs: compareInventorySnapshots(beforeItems, afterItems), checkedAt: new Date().toISOString() }
    setLastActionDiff(next)
    setActionDiffHistory(history => [next, ...history].slice(0, 10))
  }

  const handleMutationCancel = (e: unknown) => {
    if (e instanceof Error && e.message === mutationConfirmationCancelledMessage) toast.warning('Cancelled')
    else toast.danger(e instanceof Error ? e.message : String(e))
  }

  const addSelectedTemplateItem = async () => {
    if (!selectedPlayer || !selectedTemplate) return
    const qty = clampInt(addQty, 1, 9999)
    const quality = clampInt(addQuality, 0, 5)
    let reason: string | undefined
    try {
      reason = await confirmMutation({ method: 'POST', path: '/api/v1/players/give-item', title: 'Add catalog item to inventory', summary: `Add ${qty}× ${selectedTemplate.id} to ${selectedPlayer.name}.`, target: `actor:${selectedPlayer.id} · template:${selectedTemplate.id}`, details: [`Player: ${selectedPlayer.name}`, `Online state: ${selectedPlayer.online_status || 'Offline'}`, `Template: ${selectedTemplate.id}`, `Catalog name: ${selectedTemplate.name || '—'}`, `Quantity: ${qty}`, `Quality: ${quality}`, 'This uses the direct inventory write path for the selected player.', 'A before-action inventory snapshot will be downloaded before the add request is sent.'], confirmLabel: 'Add item', forceReason: true })
    } catch (e: unknown) { handleMutationCancel(e); return }
    if (!reason) { toast.warning('Cancelled: admin reason is required for this action'); return }

    setMutationBusy(true)
    const beforeItems = [...items]
    try {
      exportBeforeActionSnapshot('add', beforeItems, { target_template: selectedTemplate, quantity: qty, quality })
      await api.players.giveItem(selectedPlayer.id, selectedTemplate.id, qty, quality, reason)
      const afterItems = await reloadInventory(selectedPlayer)
      setPostActionDiff('add', selectedTemplate.id, beforeItems, afterItems)
      toast.success(`Added ${qty}× ${selectedTemplate.id}`)
    } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setMutationBusy(false) }
  }

  const repairSelectedItem = async () => {
    if (!selectedPlayer || !selectedItem) return
    let reason: string | undefined
    try {
      reason = await confirmMutation({ method: 'POST', path: '/api/v1/players/repair-item', title: 'Repair selected inventory item', summary: `Repair ${itemLabel(selectedItem)} for ${selectedPlayer.name}.`, target: `actor:${selectedPlayer.id} · item:${selectedItem.id}`, details: [`Player: ${selectedPlayer.name}`, `Online state: ${selectedPlayer.online_status || 'Offline'}`, `Item ID: ${selectedItem.id}`, `Template: ${selectedItem.template_id}`, `Current durability: ${selectedItem.durability} / ${selectedItem.max_durability}`, 'A before-action inventory snapshot will be downloaded before the repair request is sent.'], confirmLabel: 'Repair item', forceReason: true })
    } catch (e: unknown) { handleMutationCancel(e); return }
    if (!reason) { toast.warning('Cancelled: admin reason is required for this action'); return }

    setMutationBusy(true)
    const beforeItems = [...items]
    try {
      exportBeforeActionSnapshot('repair', beforeItems, { target_item: selectedItem })
      await api.players.repairItem(selectedItem.id, reason)
      const afterItems = await reloadInventory(selectedPlayer)
      setPostActionDiff('repair', itemLabel(selectedItem), beforeItems, afterItems)
      toast.success(`Repaired ${itemLabel(selectedItem)}`)
    } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setMutationBusy(false) }
  }

  const deleteSelectedItem = async () => {
    if (!selectedPlayer || !selectedItem) return
    let reason: string | undefined
    try {
      reason = await confirmMutation({ method: 'DELETE', path: `/api/v1/players/item/${selectedItem.id}`, title: 'Remove selected inventory item', summary: `Remove ${itemLabel(selectedItem)} from ${selectedPlayer.name}.`, target: `actor:${selectedPlayer.id} · item:${selectedItem.id}`, details: [`Player: ${selectedPlayer.name}`, `Online state: ${selectedPlayer.online_status || 'Offline'}`, `Item ID: ${selectedItem.id}`, `Template: ${selectedItem.template_id}`, `Stack size: ${selectedItem.stack_size}`, `Quality: ${selectedItem.quality}`, 'This removes a persisted inventory row for the selected player.', 'A before-action inventory snapshot will be downloaded before the delete request is sent.'], confirmLabel: 'Remove item', forceReason: true })
    } catch (e: unknown) { handleMutationCancel(e); return }
    if (!reason) { toast.warning('Cancelled: admin reason is required for this action'); return }

    setMutationBusy(true)
    const beforeItems = [...items]
    try {
      exportBeforeActionSnapshot('delete', beforeItems, { target_item: selectedItem })
      await api.players.deleteItem(selectedItem.id, reason)
      const afterItems = await reloadInventory(selectedPlayer)
      setPostActionDiff('delete', itemLabel(selectedItem), beforeItems, afterItems)
      toast.success(`Removed ${itemLabel(selectedItem)}`)
    } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setMutationBusy(false) }
  }

  const setSelectedItemStackSize = async () => {
    if (!selectedPlayer || !selectedItem) return
    const nextStack = clampInt(stackSizeDraft, 1, 9999)
    if (nextStack === selectedItem.stack_size) { toast.warning('Stack size is unchanged'); return }
    let reason: string | undefined
    try {
      reason = await confirmMutation({ method: 'POST', path: '/api/v1/players/item/stack-size', title: 'Change selected item stack size', summary: `Change ${itemLabel(selectedItem)} stack from ${selectedItem.stack_size} to ${nextStack}.`, target: `actor:${selectedPlayer.id} · item:${selectedItem.id}`, details: [`Player: ${selectedPlayer.name}`, `Online state: ${selectedPlayer.online_status || 'Offline'}`, `Item ID: ${selectedItem.id}`, `Template: ${selectedItem.template_id}`, `Current stack: ${selectedItem.stack_size}`, `New stack: ${nextStack}`, 'A before-action inventory snapshot will be downloaded before the stack-size request is sent.'], confirmLabel: 'Set stack size', forceReason: true })
    } catch (e: unknown) { handleMutationCancel(e); return }
    if (!reason) { toast.warning('Cancelled: admin reason is required for this action'); return }

    setMutationBusy(true)
    const beforeItems = [...items]
    try {
      exportBeforeActionSnapshot('stack-size', beforeItems, { target_item: selectedItem, previous_stack_size: selectedItem.stack_size, next_stack_size: nextStack })
      await inventoryStudioMutations.setItemStackSize(selectedItem.id, nextStack, reason)
      const afterItems = await reloadInventory(selectedPlayer)
      setPostActionDiff('stack-size', itemLabel(selectedItem), beforeItems, afterItems)
      toast.success(`Set stack size to ${nextStack}`)
    } catch (e: unknown) { toast.danger(e instanceof Error ? e.message : String(e)) } finally { setMutationBusy(false) }
  }

  const loadComparisonSnapshot = async (file: File | null) => {
    if (!file) return
    try {
      const parsed = JSON.parse(await file.text()) as InventorySnapshot
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
          <strong style={{ color: 'var(--color-primary)' }}>Inventory Studio v2:</strong> inventory inspection, snapshot export/compare, catalog browsing, confirmed add/repair/removal, stack-size edits, post-action diff review, and browser-session action history.
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-[320px_1fr_420px] gap-3 flex-1 min-h-0 overflow-hidden">
          <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
            <Header title="Players" right={<Button size="sm" variant="ghost" onPress={loadPlayers} isDisabled={playersLoading}>{playersLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button>} />
            <div className="p-2 shrink-0"><TextInput placeholder="Search player, actor, account, controller..." value={playerSearch} onChange={setPlayerSearch} /></div>
            <div className="overflow-y-auto flex-1 min-h-0">
              {filteredPlayers.map(player => <button key={player.id} type="button" onClick={() => reloadInventory(player, true)} className="w-full text-left px-3 py-2 text-xs" style={{ borderBottom: '1px solid #1a1610', borderLeft: selectedPlayer?.id === player.id ? '2px solid var(--color-primary)' : '2px solid transparent', background: selectedPlayer?.id === player.id ? '#241e12' : 'transparent', color: 'var(--color-text)', cursor: 'pointer' }}><div className="flex items-center justify-between gap-2"><span className="font-semibold truncate">{player.name}</span><span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>#{player.id}</span></div><div className="mt-0.5 truncate" style={{ color: 'var(--color-text-dim)' }}>{player.class || 'Unknown'} · {player.map || 'unknown'} · {player.online_status || 'Offline'}</div></button>)}
              {filteredPlayers.length === 0 && <EmptyText text="No matching players" />}
            </div>
          </section>

          <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
            <Header title="Inventory" subtitle={selectedPlayer ? `${selectedPlayer.name} · ${items.length} item row(s)` : 'Select a player'} right={<div className="flex gap-2"><Button size="sm" variant="ghost" isDisabled={!selectedPlayer || itemsLoading} onPress={() => selectedPlayer && reloadInventory(selectedPlayer)}>{itemsLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button><Button size="sm" variant="outline" isDisabled={!selectedPlayer || items.length === 0} onPress={exportSnapshot}>Export Snapshot</Button></div>} />
            <div className="p-2 shrink-0"><TextInput placeholder="Search item template, name, ID, quality..." value={itemSearch} onChange={setItemSearch} disabled={!selectedPlayer} /></div>
            {itemsLoading ? <div className="flex justify-center py-12"><Spinner size="lg" /></div> : !selectedPlayer ? <EmptyText text="Select a player to load inventory." /> : filteredItems.length === 0 ? <EmptyText text={items.length === 0 ? 'No inventory rows found.' : 'No matching items.'} /> : <InventoryTable items={filteredItems} selectedItemId={selectedItemId} onSelect={setSelectedItemId} />}
          </section>

          <section className="rounded-lg flex flex-col min-h-0 overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
            <Header title="Inspector" subtitle="Catalog, selected item, and diffs" />
            <div className="overflow-y-auto flex-1 p-3 text-xs min-h-0">
              <CatalogPanel templatesLoading={templatesLoading} templates={filteredTemplates} selectedTemplate={selectedTemplate} selectedTemplateId={selectedTemplateId} templateSearch={templateSearch} addQty={addQty} addQuality={addQuality} mutationBusy={mutationBusy} selectedPlayer={selectedPlayer} onRefresh={loadTemplates} onSearch={setTemplateSearch} onSelectTemplate={setSelectedTemplateId} onQty={value => setAddQty(clampInt(value, 1, 9999))} onQuality={value => setAddQuality(clampInt(value, 0, 5))} onAdd={addSelectedTemplateItem} />
              <DiffPanel title="Post-action Diff" emptyText="No confirmed action has been completed in this session." diff={lastActionDiff} />
              <ActionHistoryPanel history={actionDiffHistory} onClear={() => { setActionDiffHistory([]); setLastActionDiff(null) }} onExport={exportActionHistory} />
              <SnapshotCompare comparisonName={comparisonName} comparisonSnapshot={comparisonSnapshot} beforeCount={comparisonItems.length} currentCount={items.length} diffs={comparisonDiffs} disabled={!selectedPlayer || items.length === 0} onLoad={loadComparisonSnapshot} />
              <SelectedItemPanel item={selectedItem} mutationBusy={mutationBusy} stackSizeDraft={stackSizeDraft} onStackSizeDraft={value => setStackSizeDraft(clampInt(value, 1, 9999))} onSetStackSize={setSelectedItemStackSize} onRepair={repairSelectedItem} onDelete={deleteSelectedItem} />
            </div>
          </section>
        </div>
      </div>
      {confirmationDialog}
    </>
  )
}

function Header({ title, subtitle, right }: { title: string; subtitle?: string; right?: React.ReactNode }) {
  return <div className="flex items-center justify-between gap-3 px-3 py-2 shrink-0" style={{ borderBottom: '1px solid #2a2418' }}><div><div className="text-xs font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{title}</div>{subtitle && <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{subtitle}</div>}</div>{right}</div>
}

function TextInput({ value, onChange, placeholder, disabled }: { value: string; onChange: (value: string) => void; placeholder: string; disabled?: boolean }) {
  return <input className="w-full rounded px-2 py-1.5 text-xs border" style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }} placeholder={placeholder} value={value} onChange={event => onChange(event.target.value)} disabled={disabled} />
}

function EmptyText({ text }: { text: string }) { return <div className="flex items-center justify-center flex-1 text-sm p-4 text-center" style={{ color: 'var(--color-text-dim)' }}>{text}</div> }
function Detail({ label, value }: { label: string; value: string }) { return <div className="rounded p-2" style={{ background: '#0a0806', border: '1px solid #2a2418' }}><div className="text-[10px] uppercase tracking-wide" style={{ color: 'var(--color-text-dim)' }}>{label}</div><div className="font-mono mt-1" style={{ color: 'var(--color-text)' }}>{value || '—'}</div></div> }

function InventoryTable({ items, selectedItemId, onSelect }: { items: InventoryItem[]; selectedItemId: number | null; onSelect: (id: number) => void }) {
  return <div className="overflow-auto flex-1 min-h-0"><table className="w-full text-xs"><thead><tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418', position: 'sticky', top: 0 }}>{['Item', 'Stack', 'Quality', 'Durability'].map(h => <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>)}</tr></thead><tbody>{items.map((item, index) => <tr key={item.id} onClick={() => onSelect(item.id)} style={{ borderBottom: '1px solid #1a1610', background: selectedItemId === item.id ? '#241e12' : index % 2 === 0 ? '#0d0b07' : '#0f0d09', cursor: 'pointer' }}><td className="px-3 py-2"><div className="font-semibold" style={{ color: 'var(--color-text)' }}>{itemLabel(item)}</div><div className="font-mono" style={{ color: 'var(--color-text-dim)' }}>{item.template_id}</div></td><td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.stack_size}</td><td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.quality}</td><td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{item.durability} / {item.max_durability}</td></tr>)}</tbody></table></div>
}

function CatalogPanel(props: { templatesLoading: boolean; templates: ItemTemplate[]; selectedTemplate: ItemTemplate | null; selectedTemplateId: string; templateSearch: string; addQty: number; addQuality: number; mutationBusy: boolean; selectedPlayer: Player | null; onRefresh: () => void; onSearch: (value: string) => void; onSelectTemplate: (id: string) => void; onQty: (value: number) => void; onQuality: (value: number) => void; onAdd: () => void }) {
  return <div className="rounded p-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}><div className="flex items-center justify-between gap-2 mb-2"><div className="font-semibold" style={{ color: 'var(--color-primary)' }}>Item Catalog</div><Button size="sm" variant="ghost" onPress={props.onRefresh} isDisabled={props.templatesLoading}>{props.templatesLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button></div><TextInput placeholder="Search catalog template ID or name..." value={props.templateSearch} onChange={props.onSearch} /><div className="mt-2 rounded overflow-y-auto" style={{ border: '1px solid #2a2418', maxHeight: 180 }}>{props.templatesLoading ? <div className="flex justify-center py-4"><Spinner size="sm" /></div> : props.templates.length === 0 ? <div className="px-3 py-4 text-center" style={{ color: 'var(--color-text-dim)' }}>No matching templates</div> : props.templates.map(template => <button key={template.id} type="button" onClick={() => props.onSelectTemplate(template.id)} className="w-full text-left px-2 py-1.5" style={{ borderBottom: '1px solid #1a1610', background: props.selectedTemplateId === template.id ? '#241e12' : 'transparent', color: 'var(--color-text)', cursor: 'pointer' }}><div className="font-mono truncate" style={{ color: props.selectedTemplateId === template.id ? 'var(--color-primary)' : 'var(--color-text)' }}>{template.id}</div>{template.name && <div className="truncate" style={{ color: 'var(--color-text-dim)' }}>{template.name}</div>}</button>)}</div>{props.selectedTemplate && <div className="mt-2 grid grid-cols-1 gap-2"><Detail label="Selected Template" value={props.selectedTemplate.id} /><Detail label="Catalog Name" value={props.selectedTemplate.name || '—'} /><div className="grid grid-cols-2 gap-2"><NumberField label="Quantity" value={props.addQty} min={1} max={9999} onChange={props.onQty} /><NumberField label="Quality" value={props.addQuality} min={0} max={5} onChange={props.onQuality} /></div><Button size="sm" onPress={props.onAdd} isDisabled={!props.selectedPlayer || props.mutationBusy}>{props.mutationBusy ? <Spinner size="sm" color="current" /> : null}Add Selected Catalog Item</Button></div>}</div>
}

function NumberField({ label, value, min, max, onChange, disabled }: { label: string; value: number; min: number; max: number; onChange: (value: number) => void; disabled?: boolean }) {
  return <label className="flex flex-col gap-1" style={{ color: 'var(--color-text-dim)' }}>{label}<input type="number" min={min} max={max} value={value} onChange={event => onChange(Number(event.target.value))} className="rounded px-2 py-1 text-xs border" style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }} disabled={disabled} /></label>
}

function SnapshotCompare({ comparisonName, comparisonSnapshot, beforeCount, currentCount, diffs, disabled, onLoad }: { comparisonName: string; comparisonSnapshot: InventorySnapshot | null; beforeCount: number; currentCount: number; diffs: InventoryDiffRow[]; disabled: boolean; onLoad: (file: File | null) => void }) {
  return <div className="rounded p-3 mt-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}><div className="font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Snapshot Compare</div><input type="file" accept=".json,application/json" className="text-xs" style={{ color: 'var(--color-text-dim)' }} disabled={disabled} onChange={event => onLoad(event.target.files?.[0] ?? null)} />{comparisonName && <div className="mt-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{comparisonName}</div>}{comparisonSnapshot && <div className="grid grid-cols-3 gap-2 mt-3"><Detail label="Before" value={String(beforeCount)} /><Detail label="Current" value={String(currentCount)} /><Detail label="Diffs" value={String(diffs.length)} /></div>}{comparisonSnapshot && <DiffList diffs={diffs} emptyText="No differences detected against loaded snapshot." />}</div>
}

function DiffPanel({ title, emptyText, diff }: { title: string; emptyText: string; diff: LastActionDiff }) {
  return <div className="rounded p-3 mt-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}><div className="font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>{title}</div>{!diff ? <div style={{ color: 'var(--color-text-dim)' }}>{emptyText}</div> : <><div className="grid grid-cols-3 gap-2"><Detail label="Action" value={diff.action} /><Detail label="Before/After" value={`${diff.beforeCount} / ${diff.afterCount}`} /><Detail label="Diffs" value={String(diff.diffs.length)} /></div><div className="font-mono mt-2" style={{ color: 'var(--color-text-dim)' }}>{diff.target} · {diff.checkedAt}</div><DiffList diffs={diff.diffs} emptyText="No row-level differences detected after reload." /></>}</div>
}

function ActionHistoryPanel({ history, onClear, onExport }: { history: ActionDiff[]; onClear: () => void; onExport: () => void }) {
  return <div className="rounded p-3 mt-3" style={{ background: '#0a0806', border: '1px solid #2a2418' }}><div className="flex items-center justify-between gap-2 mb-2"><div className="font-semibold" style={{ color: 'var(--color-primary)' }}>Action History</div><div className="flex gap-2"><Button size="sm" variant="outline" onPress={onExport} isDisabled={history.length === 0}>Export</Button><Button size="sm" variant="ghost" onPress={onClear} isDisabled={history.length === 0}>Clear</Button></div></div>{history.length === 0 ? <div style={{ color: 'var(--color-text-dim)' }}>No completed Inventory Studio actions in this browser session.</div> : <div className="flex flex-col gap-2">{history.map((entry, index) => <div key={`${entry.checkedAt}-${entry.action}-${entry.target}`} className="rounded p-2" style={{ background: index === 0 ? '#241e12' : '#0f0d09', border: '1px solid #2a2418' }}><div className="flex items-center justify-between gap-2"><span className="font-semibold" style={{ color: 'var(--color-text)' }}>{entry.action.toUpperCase()} · {entry.target}</span><span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>{entry.diffs.length} diff(s)</span></div><div className="mt-1 font-mono" style={{ color: 'var(--color-text-dim)' }}>{entry.checkedAt} · {entry.beforeCount} → {entry.afterCount} rows</div>{entry.diffs.length > 0 && <div className="mt-1" style={{ color: 'var(--color-text-dim)' }}>{entry.diffs.slice(0, 3).map(diff => `${diff.status}: item:${diff.key}`).join(' · ')}{entry.diffs.length > 3 ? ` · +${entry.diffs.length - 3} more` : ''}</div>}</div>)}</div>}</div>
}

function DiffList({ diffs, emptyText }: { diffs: InventoryDiffRow[]; emptyText: string }) {
  if (diffs.length === 0) return <div className="mt-3" style={{ color: 'var(--color-success)' }}>{emptyText}</div>
  return <div className="mt-3 flex flex-col gap-2">{diffs.slice(0, 25).map(diff => <div key={`${diff.status}-${diff.key}`} className="rounded p-2" style={{ background: '#0f0d09', border: '1px solid #2a2418' }}><div className="flex items-center justify-between gap-2"><span className="font-semibold" style={{ color: diff.status === 'removed' ? '#e88' : diff.status === 'added' ? 'var(--color-success)' : '#f0a830' }}>{diff.status.toUpperCase()}</span><span className="font-mono" style={{ color: 'var(--color-text-dim)' }}>item:{diff.key}</span></div><div className="mt-1" style={{ color: 'var(--color-text)' }}>{itemLabel((diff.after ?? diff.before) as InventoryItem)}</div><ul className="mt-1 pl-4" style={{ color: 'var(--color-text-dim)', listStyle: 'disc' }}>{diff.changes.map(change => <li key={change}>{change}</li>)}</ul></div>)}{diffs.length > 25 && <div style={{ color: 'var(--color-text-dim)' }}>Showing first 25 of {diffs.length} differences.</div>}</div>
}

function SelectedItemPanel({ item, mutationBusy, stackSizeDraft, onStackSizeDraft, onSetStackSize, onRepair, onDelete }: { item: InventoryItem | null; mutationBusy: boolean; stackSizeDraft: number; onStackSizeDraft: (value: number) => void; onSetStackSize: () => void; onRepair: () => void; onDelete: () => void }) {
  if (!item) return <div className="flex items-center justify-center text-sm p-4 text-center mt-3" style={{ color: 'var(--color-text-dim)' }}>Select an inventory row to inspect item details.</div>
  return <div className="mt-3"><div className="font-semibold text-sm" style={{ color: 'var(--color-text)' }}>{itemLabel(item)}</div><div className="font-mono mt-1" style={{ color: 'var(--color-text-dim)' }}>{item.template_id}</div><div className="grid grid-cols-2 gap-2 mt-4"><Detail label="Item ID" value={String(item.id)} /><Detail label="Stack" value={String(item.stack_size)} /><Detail label="Quality" value={String(item.quality)} /><Detail label="Durability" value={`${item.durability} / ${item.max_durability}`} /></div><div className="rounded p-3 mt-4" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}><div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Confirmed stack-size edit</div><p>Stack-size edits download a before-action inventory snapshot, require shared mutation confirmation, require an admin reason, reload inventory, and update post-action diff/history.</p><div className="grid grid-cols-[1fr_auto] gap-2 mt-3 items-end"><NumberField label="New stack size" value={stackSizeDraft} min={1} max={9999} onChange={onStackSizeDraft} disabled={mutationBusy} /><Button size="sm" onPress={onSetStackSize} isDisabled={mutationBusy || stackSizeDraft === item.stack_size}>{mutationBusy ? <Spinner size="sm" color="current" /> : null}Set Stack</Button></div></div><div className="rounded p-3 mt-4" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}><div className="font-semibold mb-1" style={{ color: 'var(--color-primary)' }}>Confirmed item actions</div><p>Repair and remove actions download a before-action inventory snapshot, then require shared mutation confirmation and an admin reason.</p><div className="flex flex-wrap gap-2 mt-3"><Button size="sm" onPress={onRepair} isDisabled={mutationBusy}>{mutationBusy ? <Spinner size="sm" color="current" /> : null}Repair Selected Item</Button><Button size="sm" variant="danger-soft" onPress={onDelete} isDisabled={mutationBusy}>Remove Selected Item</Button></div></div><pre className="rounded p-3 mt-4 overflow-auto" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>{JSON.stringify(item, null, 2)}</pre></div>
}
