// ── Give Item Modal ────────────────────────────────────────────────────────────

type GiveItemDraft = {
  id: number
  template: string
  label: string
  qty: number
  quality: number
  stack_size: number
}

const newGiveItemDraft = (id: number): GiveItemDraft => ({
  id,
  template: '',
  label: '',
  qty: 1,
  quality: 1,
  stack_size: 1,
})

function GiveItemModal({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [templates, setTemplates] = useState<{id: string; name: string}[]>([])
  const [query, setQuery] = useState('')
  const [rows, setRows] = useState<GiveItemDraft[]>([newGiveItemDraft(1)])
  const [activeRowId, setActiveRowId] = useState(1)
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!open) return
    const first = newGiveItemDraft(1)
    setRows([first])
    setActiveRowId(first.id)
    setQuery('')
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
  const readyRows = rows.filter(r => r.template.trim() && r.qty > 0 && r.stack_size > 0)

  const patchRow = (id: number, patch: Partial<GiveItemDraft>) => {
    setRows(prev => prev.map(r => r.id === id ? { ...r, ...patch } : r))
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

  const pick = (t: {id: string; name: string}) => {
    if (!activeRow) return
    patchRow(activeRow.id, {
      template: t.id,
      label: t.name ? `${t.id}  —  ${t.name}` : t.id,
    })
    setQuery('')
  }

  const handleSubmit = async () => {
    if (readyRows.length === 0) {
      toast.warning('Select at least one item')
      return
    }

    setSubmitting(true)
    try {
      await api.players.giveItems(player.id, readyRows.map(r => ({
        template: r.template.trim(),
        qty: r.qty,
        quality: r.quality,
        stack_size: r.stack_size,
      })))
      toast.success(`Gave ${readyRows.length} item row(s) to ${player.name}`)
      onClose()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setSubmitting(false)
    }
  }

  const inputStyle = {
    background: 'var(--color-surface)',
    color: 'var(--color-text)',
    borderColor: '#3a3020',
    outline: 'none',
  }

  return (
    <Modal>
      <Modal.Backdrop isOpen={open} onOpenChange={v => !v && onClose()}>
        <Modal.Container size="full">
          <Modal.Dialog style={{ maxHeight: '88vh', display: 'flex', flexDirection: 'column' }}>
            <Modal.CloseTrigger />
            <Modal.Header>
              <Modal.Heading>Give Items — {player.name}</Modal.Heading>
            </Modal.Header>
            <Modal.Body style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden', padding: '12px 16px' }}>
              {loading ? (
                <div className="flex justify-center py-6"><Spinner size="lg" /></div>
              ) : (
                <div className="flex flex-col gap-3 h-full overflow-hidden">
                  <div className="flex items-center gap-3 shrink-0">
                    <Button variant="tertiary" size="sm" onPress={onClose}>Cancel</Button>
                    <Button variant="outline" size="sm" onPress={addRow}>Add Item Row</Button>
                    <Button size="sm" onPress={handleSubmit} isDisabled={submitting || readyRows.length === 0}>
                      {submitting ? <Spinner size="sm" color="current" /> : null}
                      Give Selected Items
                    </Button>
                    <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>
                      {readyRows.length} ready / {rows.length} row(s)
                    </span>
                  </div>

                  <div className="overflow-auto rounded-lg shrink-0" style={{ border: '1px solid #2a2418', maxHeight: '34vh' }}>
                    <table className="w-full text-xs">
                      <thead>
                        <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
                          {['Item', 'Quantity', 'Grade', 'Stack Size', 'Total', ''].map(h => (
                            <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {rows.map((row, i) => (
                          <tr
                            key={row.id}
                            onClick={() => setActiveRowId(row.id)}
                            style={{
                              borderBottom: '1px solid #1a1610',
                              background: activeRowId === row.id ? '#241e12' : i % 2 === 0 ? '#0d0b07' : '#0f0d09',
                              cursor: 'pointer',
                            }}
                          >
                            <td className="px-3 py-2 font-mono" style={{ color: row.template ? 'var(--color-text)' : 'var(--color-text-dim)' }}>
                              {row.label || 'Select from search list below...'}
                            </td>
                            <td className="px-3 py-2">
                              <input
                                type="number"
                                min={1}
                                max={9999}
                                value={row.qty}
                                onChange={e => patchRow(row.id, { qty: Math.max(1, parseInt(e.target.value) || 1) })}
                                className="rounded px-2 py-1 text-sm border w-24"
                                style={inputStyle}
                              />
                            </td>
                            <td className="px-3 py-2">
                              <input
                                type="number"
                                min={0}
                                max={5}
                                value={row.quality}
                                onChange={e => patchRow(row.id, { quality: Math.max(0, Math.min(5, parseInt(e.target.value) || 0)) })}
                                className="rounded px-2 py-1 text-sm border w-20"
                                style={inputStyle}
                              />
                            </td>
                            <td className="px-3 py-2">
                              <input
                                type="number"
                                min={1}
                                max={9999}
                                value={row.stack_size}
                                onChange={e => patchRow(row.id, { stack_size: Math.max(1, parseInt(e.target.value) || 1) })}
                                className="rounded px-2 py-1 text-sm border w-24"
                                style={inputStyle}
                              />
                            </td>
                            <td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>
                              {(row.qty * row.stack_size).toLocaleString()}
                            </td>
                            <td className="px-3 py-2">
                              <Button size="sm" variant="danger-soft" onPress={() => removeRow(row.id)}>Remove</Button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>

                  <div className="shrink-0">
                    <div className="text-xs mb-1" style={{ color: 'var(--color-text-dim)' }}>
                      Selecting item for row:{' '}
                      <span className="font-mono" style={{ color: 'var(--color-primary)' }}>{activeRow?.label || 'empty row'}</span>
                    </div>
                    <input
                      className="rounded px-3 py-2 text-sm border w-full"
                      style={inputStyle}
                      placeholder="Search by template ID or item name..."
                      value={query}
                      onChange={e => setQuery(e.target.value)}
                      autoFocus
                    />
                  </div>

                  <div className="flex-1 overflow-y-auto rounded-lg min-h-0" style={{ border: '1px solid #2a2418', background: '#0a0806' }}>
                    {filtered.length === 0 ? (
                      <div className="flex items-center justify-center h-full py-8 text-xs" style={{ color: 'var(--color-text-dim)' }}>
                        No matching templates
                      </div>
                    ) : (
                      filtered.map(t => (
                        <div
                          key={t.id}
                          className="flex items-baseline gap-3 px-3 py-2 cursor-pointer"
                          style={{ borderBottom: '1px solid #1a1610', background: activeRow?.template === t.id ? '#241e12' : 'transparent' }}
                          onMouseEnter={e => { if (activeRow?.template !== t.id) e.currentTarget.style.background = '#161208' }}
                          onMouseLeave={e => { if (activeRow?.template !== t.id) e.currentTarget.style.background = 'transparent' }}
                          onClick={() => pick(t)}
                        >
                          <span className="font-mono text-xs shrink-0" style={{ color: activeRow?.template === t.id ? 'var(--color-primary)' : 'var(--color-text)' }}>{t.id}</span>
                          {t.name && <span className="text-xs truncate" style={{ color: 'var(--color-text-dim)' }}>{t.name}</span>}
                        </div>
                      ))
                    )}
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
