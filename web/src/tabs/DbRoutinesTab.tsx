import { useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { dbFunctionApi, type DBFunctionInspection, type DBFunctionRow } from '../api/dbFunctions'

export default function DbRoutinesTab() {
  const [term, setTerm] = useState('')
  const [rows, setRows] = useState<DBFunctionRow[]>([])
  const [selected, setSelected] = useState<DBFunctionInspection | null>(null)
  const [loading, setLoading] = useState(false)

  const discover = async () => {
    setLoading(true)
    setSelected(null)
    try {
      setRows(await dbFunctionApi.list(term, ''))
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  const inspect = async (oid: string) => {
    setLoading(true)
    try {
      setSelected(await dbFunctionApi.inspect(oid))
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-col gap-4 h-full min-h-0">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h2 className="text-base font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Database Routine Discovery</h2>
        <p className="text-sm mb-4" style={{ color: 'var(--color-text-dim)' }}>Find candidate database routines for item delivery, inventories, rewards, notifications, movement, guilds, factions, currency, and progression.</p>
        <div className="flex gap-3 items-end">
          <div className="flex flex-col gap-1">
            <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Search</label>
            <input className="rounded px-3 py-1.5 text-sm border w-72" style={{ background: 'var(--color-background)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }} value={term} onChange={e => setTerm(e.target.value)} placeholder="item, inventory, reward" />
          </div>
          <Button size="sm" onPress={discover} isDisabled={loading}>{loading ? <Spinner size="sm" color="current" /> : null}Discover</Button>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 flex-1 min-h-0 overflow-hidden">
        <div className="rounded-lg overflow-auto" style={{ border: '1px solid #2a2418' }}>
          <table className="w-full text-xs">
            <thead><tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>{['Routine', 'Category', 'References', ''].map(h => <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>)}</tr></thead>
            <tbody>
              {rows.map((row, i) => <tr key={row.oid} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}><td className="px-3 py-2 font-mono" style={{ color: 'var(--color-text)' }}>{row.schema}.{row.name}<div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{row.arguments}</div></td><td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.category}</td><td className="px-3 py-2" style={{ color: 'var(--color-text-dim)' }}>{row.references.join(', ')}</td><td className="px-3 py-2"><Button size="sm" variant="outline" onPress={() => inspect(row.oid)}>Inspect</Button></td></tr>)}
              {rows.length === 0 && <tr><td colSpan={4} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>No routines loaded.</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="rounded-lg p-4 overflow-auto" style={{ background: '#0a0806', border: '1px solid #2a2418' }}>
          {!selected ? <p className="text-sm" style={{ color: 'var(--color-text-dim)' }}>Select a routine to inspect.</p> : <div className="flex flex-col gap-3"><h3 className="font-semibold" style={{ color: 'var(--color-primary)' }}>{selected.schema}.{selected.name}</h3><div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Risk: {selected.risk}</div><ul className="list-disc pl-5 text-xs" style={{ color: 'var(--color-text-dim)' }}>{selected.notes.map(note => <li key={note}>{note}</li>)}</ul><pre className="text-xs font-mono whitespace-pre-wrap rounded p-3 overflow-auto" style={{ background: '#070604', color: '#e8dcc8', border: '1px solid #2a2418' }}>{selected.definition}</pre></div>}
        </div>
      </div>
    </div>
  )
}
