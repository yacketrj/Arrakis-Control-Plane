import { useEffect, useMemo, useState } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { GiveItemAugment, Player } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'
import { AUGMENT_PRESETS, findAugmentPreset } from './augmentPresets'
import { clampFloat, clampInt, parseRollsCsv, presetAugment, toGiveItemPayload } from './giveItemPayload'

type TemplateOption = { id: string; name: string }
type DeliveryMode = 'inventory' | 'live'

type GiveItemDraft = {
  id: number
  template: string
  label: string
  qty: number
  quality: number
  stack_size: number
  augments: GiveItemAugment[]
}

const inputStyle = {
  background: 'var(--color-surface)',
  color: 'var(--color-text)',
  borderColor: '#3a3020',
  outline: 'none',
}

const newGiveItemDraft = (id: number): GiveItemDraft => ({
  id,
  template: '',
  label: '',
  qty: 1,
  quality: 1,
  stack_size: 1,
  augments: [],
})

export default function GiveItemModalAugmented({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [templates, setTemplates] = useState<TemplateOption[]>([])
  const [query, setQuery] = useState('')
  const [augmentQuery, setAugmentQuery] = useState('')
  const [rows, setRows] = useState<GiveItemDraft[]>([newGiveItemDraft(1)])
  const [activeRowId, setActiveRowId] = useState(1)
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [showPayload, setShowPayload] = useState(false)
  const [showAdvancedAugments, setShowAdvancedAugments] = useState(false)
  const [deliveryMode, setDeliveryMode] = useState<DeliveryMode>('inventory')
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  useEffect(() => {
    if (!open) return
    const first = newGiveItemDraft(1)
    setRows([first])
    setActiveRowId(first.id)
    setQuery('')
    setAugmentQuery('')
    setShowPayload(false)
    setShowAdvancedAugments(false)
    setDeliveryMode('inventory')
    setLoading(true)
    api.players.templates()
      .then(setTemplates)
      .catch((e: unknown) => toast.danger(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false))
  }, [open])

  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim()
    if (!q) return templates.slice(0, 80)
    return templates
      .filter(t => t.id.toLowerCase().includes(q) || t.name.toLowerCase().includes(q))
      .slice(0, 120)
  }, [templates, query])

  const augmentTemplates = useMemo(() => {
    return templates
      .filter(t => `${t.id} ${t.name}`.toLowerCase().includes('augment'))
      .sort((a, b) => a.id.localeCompare(b.id))
  }, [templates])

  const filteredAugmentTemplates = useMemo(() => {
    const q = augmentQuery.toLowerCase().trim()
    if (!q) return augmentTemplates.slice(0, 40)
    return augmentTemplates
      .filter(t => t.id.toLowerCase().includes(q) || t.name.toLowerCase().includes(q))
      .slice(0, 80)
  }, [augmentTemplates, augmentQuery])

  const activeRow = rows.find(r => r.id === activeRowId) ?? rows[0]
  const readyRows = useMemo(() => rows.filter(r => r.template.trim() && r.qty > 0 && r.stack_size > 0), [rows])
  const liveIncompatibleRows = useMemo(() => readyRows.filter(r => r.quality > 0 || r.augments.some(a => a.name.trim())), [readyRows])
  const liveModeDisabled = deliveryMode === 'live' && liveIncompatibleRows.length > 0

  const patchRow = (id: number, patch: Partial<GiveItemDraft>) => {
    setRows(prev => prev.map(r => r.id === id ? { ...r, ...patch } : r))
  }

  const patchAugment = (rowId: number, index: number, patch: Partial<GiveItemAugment>) => {
    setRows(prev => prev.map(row => {
      if (row.id !== rowId) return row
      return { ...row, augments: row.augments.map((aug, i) => i === index ? { ...aug, ...patch } : aug) }
    }))
  }

  const addRow = () => {
    const id = Math.max(0, ...rows.map(r => r.id)) + 1
    const row = newGiveItemDraft(id)
    setRows(prev => [...prev, row])
    setActiveRowId(row.id)
    setQuery('')
    setAugmentQuery('')
  }

  const removeRow = (id: number) => {
    setRows(prev => {
      if (prev.length === 1) {
        const row = newGiveItemDraft(1)
        setActiveRowId(row.id)
        setQuery('')
        setAugmentQuery('')
        return [row]
      }
      const next = prev.filter(r => r.id !== id)
      if (activeRowId === id) {
        setActiveRowId(next[0].id)
        setQuery('')
        setAugmentQuery('')
      }
      return next
    })
  }

  const addAugment = (rowId: number, name = '') => {
    setRows(prev => prev.map(row => {
      if (row.id !== rowId || row.augments.length >= 5) return row
      return { ...row, augments: [...row.augments, name ? presetAugment(name) : { name: '', grade: 5, roll: 1, roll_count: 1, effect_indices: [] }] }
    }))
  }

  const applyPreset = (rowId: number, index: number, name: string) => {
    patchAugment(rowId, index, presetAugment(name))
  }

  const removeAugment = (rowId: number, index: number) => {
    setRows(prev => prev.map(row => row.id === rowId ? { ...row, augments: row.augments.filter((_, i) => i !== index) } : row))
  }

  const pick = (template: TemplateOption) => {
    if (!activeRow) return
    patchRow(activeRow.id, { template: template.id, label: template.name ? `${template.id}  —  ${template.name}` : template.id })
    setQuery('')
  }

  const payloadPreview = useMemo(() => {
    if (deliveryMode === 'live') {
      return {
        delivery_mode: 'live_claim_rewards',
        player_controller_id: player.controller_id,
        items: readyRows.map(row => ({
          template: row.template.trim(),
          amount: row.qty * row.stack_size,
        })),
      }
    }
    return { delivery_mode: 'offline_inventory_write', player_id: player.id, items: readyRows.map(toGiveItemPayload) }
  }, [deliveryMode, player.controller_id, player.id, readyRows])

  const handleSubmit = async () => {
    if (readyRows.length === 0) {
      toast.warning('Select at least one item')
      return
    }
    if (deliveryMode === 'live' && liveModeDisabled) {
      toast.warning('Live delivery only supports plain, unaugmented, non-grade item grants. Use Inventory Write for graded or augmented items.')
      return
    }

    const mutationPath = deliveryMode === 'live' ? '/api/v1/players/grant-live' : '/api/v1/players/give-item'
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: mutationPath,
        title: deliveryMode === 'live' ? 'Queue live claim reward' : 'Give items by inventory write',
        summary: deliveryMode === 'live'
          ? `Queue ${readyRows.length} live claim reward row(s) for ${player.name}.`
          : `Write ${readyRows.length} item row(s) directly to ${player.name}'s inventory.`,
        target: `actor:${player.id} · ctrl:${player.controller_id} · acct:${player.account_id}`,
        details: [
          `Delivery mode: ${deliveryMode === 'live' ? 'Live Claim Rewards' : 'Inventory Write'}`,
          `Item rows: ${readyRows.length}`,
          deliveryMode === 'live'
            ? 'Live delivery supports plain template-and-amount grants only.'
            : 'Inventory writes may require the player to relog before the client refreshes inventory.',
        ],
        confirmLabel: deliveryMode === 'live' ? 'Queue reward' : 'Give items',
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

    setSubmitting(true)
    try {
      if (deliveryMode === 'live') {
        for (const row of readyRows) {
          await api.players.grantLive(player.controller_id, row.template, row.qty * row.stack_size, reason)
        }
        toast.success(`Queued ${readyRows.length} live claim reward row(s) for ${player.name}`)
      } else {
        await api.players.giveItems(player.id, readyRows.map(toGiveItemPayload), reason)
        toast.success(`Gave ${readyRows.length} item row(s) to ${player.name}`)
      }
      onClose()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Modal>
        <Modal.Backdrop isOpen={open} onOpenChange={v => !v && onClose()}>
          <Modal.Container size="full">
            <Modal.Dialog style={{ maxHeight: '92vh', display: 'flex', flexDirection: 'column' }}>
              <Modal.CloseTrigger />
              <Modal.Header><Modal.Heading>Give Items &amp; Augments — {player.name}</Modal.Heading></Modal.Header>
              <Modal.Body style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden', padding: '12px 16px' }}>
                {loading ? <div className="flex justify-center py-6"><Spinner size="lg" /></div> : (
                  <div className="flex flex-col gap-3 h-full overflow-hidden">
                    <div className="flex items-center gap-3 shrink-0 flex-wrap">
                      <Button variant="tertiary" size="sm" onPress={onClose}>Cancel</Button>
                      <Button variant="outline" size="sm" onPress={addRow}>Add Item Row</Button>
                      <Button variant="outline" size="sm" onPress={() => setShowPayload(v => !v)}>Payload Preview</Button>
                      <Button variant="outline" size="sm" onPress={() => setShowAdvancedAugments(v => !v)}>{showAdvancedAugments ? 'Hide Advanced Augments' : 'Show Advanced Augments'}</Button>
                      <Button size="sm" onPress={handleSubmit} isDisabled={submitting || readyRows.length === 0 || liveModeDisabled}>{submitting ? <Spinner size="sm" color="current" /> : null}{deliveryMode === 'live' ? 'Queue Live Claim' : 'Give Selected Items'}</Button>
                      <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{readyRows.length} ready / {rows.length} row(s)</span>
                    </div>

                    <div className="rounded-lg p-3 text-xs" style={{ border: '1px solid #2a2418', background: '#0f0d09', color: 'var(--color-text-dim)' }}>
                      <div className="flex gap-3 items-center flex-wrap">
                        <strong style={{ color: 'var(--color-primary)' }}>Delivery mode</strong>
                        <label className="flex gap-1 items-center"><input type="radio" checked={deliveryMode === 'inventory'} onChange={() => setDeliveryMode('inventory')} /> Inventory Write</label>
                        <label className="flex gap-1 items-center"><input type="radio" checked={deliveryMode === 'live'} onChange={() => setDeliveryMode('live')} /> Live Claim Rewards</label>
                      </div>
                      <div className="mt-2">
                        {deliveryMode === 'live'
                          ? 'Live Claim Rewards is for plain items only: template + amount. It cannot attach item grade, weapon augments, armor augments, custom rolls, or direct inventory placement.'
                          : 'Inventory Write is the advanced path. Pick an item, set item grade/stacking, then optionally attach weapon or armor augment template IDs. Online players may need to relog before the client refreshes inventory.'}
                      </div>
                      <div className="mt-2">
                        <strong style={{ color: 'var(--color-primary)' }}>Audit safety:</strong> item grants require an admin reason and are recorded in the Audit tab with mutation safety metadata.
                      </div>
                      <div className="mt-2">
                        <strong style={{ color: 'var(--color-primary)' }}>Augment workflow:</strong> select the item row first, then add augments to that row. Weapon preset buttons are shortcuts only. Armor augments use the same system: add them from the Augment Template Search or paste the armor augment template ID into Custom Augment.
                      </div>
                      {showAdvancedAugments && <div className="mt-2">Advanced rolls: leave explicit rolls blank for normal use. Use comma-separated rolls like <span className="font-mono">1,1,1</span> only when an augment has multiple stat effects that need separate roll strengths.</div>}
                      {liveModeDisabled && <div className="mt-2" style={{ color: '#ffb86b' }}>Live mode blocked because one or more selected rows has item grade &gt; 0 or augments. Switch to Inventory Write or remove grades/augments.</div>}
                    </div>

                    {showPayload && <pre className="text-xs rounded p-3 overflow-auto max-h-36" style={{ background: '#070604', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>{JSON.stringify(payloadPreview, null, 2)}</pre>}

                    <div className="overflow-auto rounded-lg shrink-0" style={{ border: '1px solid #2a2418', maxHeight: '44vh' }}>
                      <table className="w-full text-xs">
                        <thead><tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>{['Item', 'Number of stacks', 'Item grade', 'Stack size', 'Augments attached to this item row', ''].map(h => <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>)}</tr></thead>
                        <tbody>
                          {rows.map((row, i) => (
                            <tr key={row.id} onClick={() => setActiveRowId(row.id)} style={{ borderBottom: '1px solid #1a1610', background: activeRowId === row.id ? '#241e12' : i % 2 === 0 ? '#0d0b07' : '#0f0d09', cursor: 'pointer' }}>
                              <td className="px-3 py-2 font-mono" style={{ color: row.template ? 'var(--color-text)' : 'var(--color-text-dim)' }}>{row.label || 'Select from item search list below...'}<div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Total granted: {(row.qty * row.stack_size).toLocaleString()}</div></td>
                              <td className="px-3 py-2"><input type="number" min={1} max={9999} value={row.qty} onChange={e => patchRow(row.id, { qty: clampInt(e.target.value, 1, 9999, 1) })} className="rounded px-2 py-1 text-sm border w-20" style={inputStyle} /></td>
                              <td className="px-3 py-2"><input type="number" min={0} max={5} value={row.quality} onChange={e => patchRow(row.id, { quality: clampInt(e.target.value, 0, 5, 0) })} className="rounded px-2 py-1 text-sm border w-16" style={inputStyle} disabled={deliveryMode === 'live'} /></td>
                              <td className="px-3 py-2"><input type="number" min={1} max={9999} value={row.stack_size} onChange={e => patchRow(row.id, { stack_size: clampInt(e.target.value, 1, 9999, 1) })} className="rounded px-2 py-1 text-sm border w-24" style={inputStyle} /></td>
                              <td className="px-3 py-2">
                                <div className="flex flex-col gap-2">
                                  {row.augments.length === 0 && <div style={{ color: 'var(--color-text-dim)' }}>No augments attached to this item row.</div>}
                                  {row.augments.map((aug, index) => (
                                    <div key={`${row.id}-${index}`} className="flex gap-1 items-center flex-wrap" style={{ opacity: deliveryMode === 'live' ? 0.55 : 1 }}>
                                      <select className="rounded px-2 py-1 text-xs border w-52" style={inputStyle} value={findAugmentPreset(aug.name)?.name ?? ''} onChange={e => e.target.value && applyPreset(row.id, index, e.target.value)} disabled={deliveryMode === 'live'} title="Optional weapon preset shortcut">
                                        <option value="">Custom / armor augment</option>
                                        {AUGMENT_PRESETS.map(p => <option key={p.name} value={p.name}>{p.label}</option>)}
                                      </select>
                                      <input className="rounded px-2 py-1 text-xs border w-64" style={inputStyle} placeholder="Augment template ID, weapon or armor" value={aug.name} onChange={e => patchAugment(row.id, index, { name: e.target.value })} disabled={deliveryMode === 'live'} />
                                      <label className="flex items-center gap-1" style={{ color: 'var(--color-text-dim)' }}>Aug grade <input type="number" min={1} max={5} className="rounded px-2 py-1 text-xs border w-14" style={inputStyle} value={aug.grade} onChange={e => patchAugment(row.id, index, { grade: clampInt(e.target.value, 1, 5, 5) })} disabled={deliveryMode === 'live'} /></label>
                                      {showAdvancedAugments && <label className="flex items-center gap-1" style={{ color: 'var(--color-text-dim)' }}>Roll <input type="number" min={0} max={1} step={0.01} className="rounded px-2 py-1 text-xs border w-16" style={inputStyle} value={aug.roll ?? 1} onChange={e => patchAugment(row.id, index, { roll: clampFloat(e.target.value, 0, 1, 1), rolls: undefined })} disabled={deliveryMode === 'live'} /></label>}
                                      {showAdvancedAugments && <label className="flex items-center gap-1" style={{ color: 'var(--color-text-dim)' }}>Slots <input type="number" min={1} max={8} className="rounded px-2 py-1 text-xs border w-14" style={inputStyle} value={aug.roll_count ?? 1} onChange={e => patchAugment(row.id, index, { roll_count: clampInt(e.target.value, 1, 8, 1), rolls: undefined })} disabled={deliveryMode === 'live'} /></label>}
                                      {showAdvancedAugments && <input className="rounded px-2 py-1 text-xs border w-48" style={inputStyle} placeholder="Explicit rolls, e.g. 1,1,1" value={(aug.rolls ?? []).join(',')} onChange={e => patchAugment(row.id, index, { rolls: parseRollsCsv(e.target.value) })} title="Optional comma-separated per-effect roll strengths" disabled={deliveryMode === 'live'} />}
                                      <Button size="sm" variant="danger-soft" onPress={() => removeAugment(row.id, index)} isDisabled={deliveryMode === 'live'}>Remove</Button>
                                    </div>
                                  ))}
                                  <div className="flex gap-2 flex-wrap">
                                    <Button size="sm" variant="ghost" onPress={() => addAugment(row.id)} isDisabled={deliveryMode === 'live' || row.augments.length >= 5}>Add Custom / Armor Augment</Button>
                                    {AUGMENT_PRESETS.map(p => <Button key={p.name} size="sm" variant="ghost" onPress={() => addAugment(row.id, p.name)} isDisabled={deliveryMode === 'live' || row.augments.length >= 5}>{p.label}</Button>)}
                                  </div>
                                </div>
                              </td>
                              <td className="px-3 py-2"><Button size="sm" variant="danger-soft" onPress={() => removeRow(row.id)}>Remove</Button></td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-3 min-h-0 flex-1">
                      <div className="flex flex-col min-h-0">
                        <div className="text-xs mb-1" style={{ color: 'var(--color-text-dim)' }}>Selecting item for row: <span className="font-mono" style={{ color: 'var(--color-primary)' }}>{activeRow?.label || 'empty row'}</span></div>
                        <input className="rounded px-3 py-2 text-sm border w-full" style={inputStyle} placeholder="Search item by template ID or item name..." value={query} onChange={e => setQuery(e.target.value)} autoFocus />
                        <div className="flex-1 overflow-y-auto rounded-lg min-h-0 mt-2" style={{ border: '1px solid #2a2418', background: '#0a0806' }}>
                          {filtered.length === 0 ? <div className="flex items-center justify-center h-full py-8 text-xs" style={{ color: 'var(--color-text-dim)' }}>No matching item templates</div> : filtered.map(t => <div key={t.id} className="flex items-baseline gap-3 px-3 py-2 cursor-pointer" style={{ borderBottom: '1px solid #1a1610', background: activeRow?.template === t.id ? '#241e12' : 'transparent' }} onClick={() => pick(t)}><span className="font-mono text-xs shrink-0" style={{ color: activeRow?.template === t.id ? 'var(--color-primary)' : 'var(--color-text)' }}>{t.id}</span>{t.name && <span className="text-xs truncate" style={{ color: 'var(--color-text-dim)' }}>{t.name}</span>}</div>)}
                        </div>
                      </div>

                      <div className="flex flex-col min-h-0">
                        <div className="text-xs mb-1" style={{ color: 'var(--color-text-dim)' }}>Add weapon or armor augment to active row: <span className="font-mono" style={{ color: 'var(--color-primary)' }}>{activeRow?.label || 'empty row'}</span></div>
                        <input className="rounded px-3 py-2 text-sm border w-full" style={inputStyle} placeholder="Search augment templates, e.g. armor, shield, damage, reload..." value={augmentQuery} onChange={e => setAugmentQuery(e.target.value)} disabled={deliveryMode === 'live'} />
                        <div className="flex-1 overflow-y-auto rounded-lg min-h-0 mt-2" style={{ border: '1px solid #2a2418', background: '#0a0806' }}>
                          {!activeRow ? <div className="flex items-center justify-center h-full py-8 text-xs" style={{ color: 'var(--color-text-dim)' }}>Select an item row first</div> : filteredAugmentTemplates.length === 0 ? <div className="flex items-center justify-center h-full py-8 text-xs" style={{ color: 'var(--color-text-dim)' }}>No augment templates found. Use Add Custom / Armor Augment and paste the augment template ID.</div> : filteredAugmentTemplates.map(t => <div key={t.id} className="flex items-baseline gap-3 px-3 py-2 cursor-pointer" style={{ borderBottom: '1px solid #1a1610' }} onClick={() => activeRow && addAugment(activeRow.id, t.id)}><span className="font-mono text-xs shrink-0" style={{ color: 'var(--color-text)' }}>{t.id}</span>{t.name && <span className="text-xs truncate" style={{ color: 'var(--color-text-dim)' }}>{t.name}</span>}</div>)}
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </Modal.Body>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
      {confirmationDialog}
    </>
  )
}
