import { useEffect, useMemo, useState } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { GiveItemAugment, Player } from '../api/client'
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
  const [rows, setRows] = useState<GiveItemDraft[]>([newGiveItemDraft(1)])
  const [activeRowId, setActiveRowId] = useState(1)
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [showPayload, setShowPayload] = useState(false)
  const [deliveryMode, setDeliveryMode] = useState<DeliveryMode>('inventory')

  useEffect(() => {
    if (!open) return
    const first = newGiveItemDraft(1)
    setRows([first])
    setActiveRowId(first.id)
    setQuery('')
    setShowPayload(false)
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
  }

  const removeRow = (id: number) => {
    setRows(prev => {
      if (prev.length === 1) {
        const row = newGiveItemDraft(1)
        setActiveRowId(row.id)
        setQuery('')
        return [row]
      }
      const next = prev.filter(r => r.id !== id)
      if (activeRowId === id) {
        setActiveRowId(next[0].id)
        setQuery('')
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
    setSubmitting(true)
    try {
      if (deliveryMode === 'live') {
        for (const row of readyRows) {
          await api.players.grantLive(player.controller_id, row.template, row.qty * row.stack_size)
        }
        toast.success(`Queued ${readyRows.length} live claim reward row(s) for ${player.name}`)
      } else {
        await api.players.giveItems(player.id, readyRows.map(toGiveItemPayload))
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
                        ? 'Live Claim Rewards is the instant path: the online player should see Claim Rewards without logging out. It supports plain item template + amount only; grades, augments, custom stats, and direct inventory placement require Inventory Write.'
                        : 'Inventory Write modifies the database inventory directly and supports grades, stacks, and augments, but online players may need to relog before the game client refreshes the inventory.'}
                    </div>
                    {liveModeDisabled && <div className="mt-2" style={{ color: '#ffb86b' }}>Live mode blocked because one or more selected rows has item grade &gt; 0 or augments. Switch to Inventory Write or remove grades/augments.</div>}
                  </div>

                  {showPayload && <pre className="text-xs rounded p-3 overflow-auto max-h-36" style={{ background: '#070604', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>{JSON.stringify(payloadPreview, null, 2)}</pre>}

                  <div className="overflow-auto rounded-lg shrink-0" style={{ border: '1px solid #2a2418', maxHeight: '44vh' }}>
                    <table className="w-full text-xs">
                      <thead><tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>{['Item', 'Stacks', 'Grade', 'Stack Size', 'Augments', ''].map(h => <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>)}</tr></thead>
                      <tbody>
                        {rows.map((row, i) => (
                          <tr key={row.id} onClick={() => setActiveRowId(row.id)} style={{ borderBottom: '1px solid #1a1610', background: activeRowId === row.id ? '#241e12' : i % 2 === 0 ? '#0d0b07' : '#0f0d09', cursor: 'pointer' }}>
                            <td className="px-3 py-2 font-mono" style={{ color: row.template ? 'var(--color-text)' : 'var(--color-text-dim)' }}>{row.label || 'Select from search list below...'}<div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Total: {(row.qty * row.stack_size).toLocaleString()}</div></td>
                            <td className="px-3 py-2"><input type="number" min={1} max={9999} value={row.qty} onChange={e => patchRow(row.id, { qty: clampInt(e.target.value, 1, 9999, 1) })} className="rounded px-2 py-1 text-sm border w-20" style={inputStyle} /></td>
                            <td className="px-3 py-2"><input type="number" min={0} max={5} value={row.quality} onChange={e => patchRow(row.id, { quality: clampInt(e.target.value, 0, 5, 0) })} className="rounded px-2 py-1 text-sm border w-16" style={inputStyle} disabled={deliveryMode === 'live'} /></td>
                            <td className="px-3 py-2"><input type="number" min={1} max={9999} value={row.stack_size} onChange={e => patchRow(row.id, { stack_size: clampInt(e.target.value, 1, 9999, 1) })} className="rounded px-2 py-1 text-sm border w-24" style={inputStyle} /></td>
                            <td className="px-3 py-2">
                              <div className="flex flex-col gap-2">
                                {row.augments.map((aug, index) => (
                                  <div key={`${row.id}-${index}`} className="flex gap-1 items-center flex-wrap" style={{ opacity: deliveryMode === 'live' ? 0.55 : 1 }}>
                                    <select className="rounded px-2 py-1 text-xs border w-52" style={inputStyle} value={findAugmentPreset(aug.name)?.name ?? ''} onChange={e => e.target.value && applyPreset(row.id, index, e.target.value)} disabled={deliveryMode === 'live'}>
                                      <option value="">Custom augment</option>
                                      {AUGMENT_PRESETS.map(p => <option key={p.name} value={p.name}>{p.label}</option>)}
                                    </select>
                                    <input className="rounded px-2 py-1 text-xs border w-56" style={inputStyle} placeholder="T6_Augment_Damage1" value={aug.name} onChange={e => patchAugment(row.id, index, { name: e.target.value })} disabled={deliveryMode === 'live'} />
                                    <input type="number" min={1} max={5} className="rounded px-2 py-1 text-xs border w-14" style={inputStyle} value={aug.grade} onChange={e => patchAugment(row.id, index, { grade: clampInt(e.target.value, 1, 5, 5) })} title="Augment grade" disabled={deliveryMode === 'live'} />
                                    <input type="number" min={0} max={1} step={0.01} className="rounded px-2 py-1 text-xs border w-16" style={inputStyle} value={aug.roll ?? 1} onChange={e => patchAugment(row.id, index, { roll: clampFloat(e.target.value, 0, 1, 1), rolls: undefined })} title="Roll value 0.0-1.0" disabled={deliveryMode === 'live'} />
                                    <input type="number" min={1} max={8} className="rounded px-2 py-1 text-xs border w-14" style={inputStyle} value={aug.roll_count ?? 1} onChange={e => patchAugment(row.id, index, { roll_count: clampInt(e.target.value, 1, 8, 1), rolls: undefined })} title="Roll count" disabled={deliveryMode === 'live'} />
                                    <input className="rounded px-2 py-1 text-xs border w-32" style={inputStyle} placeholder="rolls csv" value={(aug.rolls ?? []).join(',')} onChange={e => patchAugment(row.id, index, { rolls: parseRollsCsv(e.target.value) })} title="Explicit rolls, comma separated" disabled={deliveryMode === 'live'} />
                                    <Button size="sm" variant="danger-soft" onPress={() => removeAugment(row.id, index)} isDisabled={deliveryMode === 'live'}>X</Button>
                                  </div>
                                ))}
                                <div className="flex gap-2 flex-wrap">
                                  <Button size="sm" variant="ghost" onPress={() => addAugment(row.id)} isDisabled={deliveryMode === 'live' || row.augments.length >= 5}>Add Custom Augment</Button>
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

                  <div className="shrink-0"><div className="text-xs mb-1" style={{ color: 'var(--color-text-dim)' }}>Selecting item for row: <span className="font-mono" style={{ color: 'var(--color-primary)' }}>{activeRow?.label || 'empty row'}</span></div><input className="rounded px-3 py-2 text-sm border w-full" style={inputStyle} placeholder="Search by template ID or item name..." value={query} onChange={e => setQuery(e.target.value)} autoFocus /></div>

                  <div className="flex-1 overflow-y-auto rounded-lg min-h-0" style={{ border: '1px solid #2a2418', background: '#0a0806' }}>
                    {filtered.length === 0 ? <div className="flex items-center justify-center h-full py-8 text-xs" style={{ color: 'var(--color-text-dim)' }}>No matching templates</div> : filtered.map(t => <div key={t.id} className="flex items-baseline gap-3 px-3 py-2 cursor-pointer" style={{ borderBottom: '1px solid #1a1610', background: activeRow?.template === t.id ? '#241e12' : 'transparent' }} onClick={() => pick(t)}><span className="font-mono text-xs shrink-0" style={{ color: activeRow?.template === t.id ? 'var(--color-primary)' : 'var(--color-text)' }}>{t.id}</span>{t.name && <span className="text-xs truncate" style={{ color: 'var(--color-text-dim)' }}>{t.name}</span>}</div>)}
                  </div>
                </div>
              )}
            </Modal.Body>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal>
  )
}
