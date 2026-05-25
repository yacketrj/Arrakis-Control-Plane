import { useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

type Section = 'tables' | 'describe' | 'sample' | 'search' | 'sql'

type TableData = { headers: string[]; rows: string[][] }

function ResultTable({ headers, rows }: TableData) {
  if (rows.length === 0) return <p className="text-sm" style={{ color: 'var(--color-text-dim)' }}>No results.</p>
  return (
    <div className="overflow-auto rounded-lg flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
      <table className="w-full text-xs">
        <thead>
          <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
            {headers.map(h => (
              <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide whitespace-nowrap" style={{ color: 'var(--color-primary)' }}>{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
              {row.map((cell, j) => (
                <td key={j} className="px-3 py-1.5 font-mono whitespace-nowrap" style={{ color: 'var(--color-text)' }}>{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function sqlPreview(sql: string): string {
  const compact = sql.trim().replace(/\s+/g, ' ')
  if (compact.length <= 240) return compact
  return `${compact.slice(0, 240)}…`
}

function sqlRiskDetails(sql: string): string[] {
  const lower = sql.toLowerCase()
  const details = [
    `SQL preview: ${sqlPreview(sql)}`,
    'Run SQL is a privileged database operation and is recorded in the audit log.',
  ]

  if (/\b(delete|update|insert|truncate|drop|alter|create|grant|revoke)\b/.test(lower)) {
    details.push('The query appears to contain mutating or schema-changing SQL keywords.')
  } else {
    details.push('The query does not obviously contain common mutating SQL keywords, but it still runs through the protected SQL endpoint.')
  }

  details.push('Verify the target database, table names, predicates, and expected blast radius before continuing.')
  return details
}

export default function DatabaseTab() {
  const [active, setActive] = useState<Section>('tables')
  const [tableInput, setTableInput] = useState('')
  const [limitInput, setLimitInput] = useState('20')
  const [searchInput, setSearchInput] = useState('')
  const [sqlInput, setSqlInput] = useState('')
  const [result, setResult] = useState<TableData | null>(null)
  const [sqlResult, setSqlResult] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  const run = async () => {
    setLoading(true)
    setResult(null)
    setSqlResult(null)
    setError(null)
    try {
      if (active === 'tables') {
        const rows = await api.database.tables()
        setResult({
          headers: ['Table', 'Rows'],
          rows: rows.map(r => [r.name, String(r.row_count)]),
        })
      } else if (active === 'describe') {
        if (!tableInput.trim()) { toast.warning('Enter a table name'); return }
        const r = await api.database.describe(tableInput.trim())
        setResult({
          headers: ['Column', 'Type', 'Nullable'],
          rows: r.columns.map(c => [c.name, c.data_type, c.nullable]),
        })
      } else if (active === 'sample') {
        if (!tableInput.trim()) { toast.warning('Enter a table name'); return }
        const r = await api.database.sample(tableInput.trim(), Number(limitInput) || 20)
        setResult({ headers: r.headers, rows: r.rows })
      } else if (active === 'search') {
        if (!searchInput.trim()) { toast.warning('Enter a search term'); return }
        const r = await api.database.search(searchInput.trim())
        setResult({ headers: r.headers, rows: r.rows })
      } else {
        const sql = sqlInput.trim()
        if (!sql) { toast.warning('Enter a SQL query'); return }

        let reason: string | undefined
        try {
          reason = await confirmMutation({
            method: 'POST',
            path: '/api/v1/database/sql',
            title: 'Run database SQL',
            summary: 'Run the SQL query through the protected database SQL endpoint.',
            target: 'database:game-management',
            details: sqlRiskDetails(sql),
            confirmLabel: 'Run SQL',
            forceReason: true,
          })
        } catch (e: unknown) {
          if (e instanceof Error && e.message === mutationConfirmationCancelledMessage) {
            toast.warning('Cancelled')
            return
          }
          throw e
        }

        if (!reason) {
          toast.warning('Cancelled: admin reason is required for this action')
          return
        }

        const r = await api.database.sql(sql, reason)
        setSqlResult(r.result)
      }
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      setError(msg)
      toast.danger(`Failed: ${msg}`)
    } finally {
      setLoading(false)
    }
  }

  const sections: { key: Section; label: string }[] = [
    { key: 'tables', label: 'Tables' },
    { key: 'describe', label: 'Describe' },
    { key: 'sample', label: 'Sample' },
    { key: 'search', label: 'Search Columns' },
    { key: 'sql', label: 'Run SQL' },
  ]

  return (
    <>
      <div className="flex gap-4 h-full min-h-0">
        <div
          className="w-40 shrink-0 flex flex-col gap-1 rounded-lg p-2"
          style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}
        >
          {sections.map(s => (
            <button
              key={s.key}
              onClick={() => { setActive(s.key); setResult(null); setSqlResult(null); setError(null) }}
              className="text-left px-3 py-2 rounded text-sm transition-colors"
              style={{
                background: active === s.key ? 'var(--color-primary)' : 'transparent',
                color: active === s.key ? '#fff' : 'var(--color-text)',
              }}
            >
              {s.label}
            </button>
          ))}
        </div>

        <div className="flex-1 overflow-auto flex flex-col gap-4 min-h-0">
          <h2 className="text-base font-semibold shrink-0" style={{ color: 'var(--color-primary)' }}>
            {sections.find(s => s.key === active)?.label}
          </h2>

          {(active === 'describe' || active === 'sample') && (
            <div className="flex gap-3 items-end shrink-0">
              <div className="flex flex-col gap-1">
                <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Table Name</label>
                <input
                  className="rounded px-3 py-1.5 text-sm font-mono border"
                  style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                  value={tableInput}
                  onChange={e => setTableInput(e.target.value)}
                  placeholder="dune.actors"
                  onKeyDown={e => e.key === 'Enter' && run()}
                />
              </div>
              {active === 'sample' && (
                <div className="flex flex-col gap-1">
                  <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Limit</label>
                  <input
                    className="rounded px-3 py-1.5 text-sm w-20 border"
                    style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                    value={limitInput}
                    onChange={e => setLimitInput(e.target.value)}
                    type="number" min={1} max={1000}
                  />
                </div>
              )}
              <Button onPress={run} isDisabled={loading} size="sm">
                {loading ? <Spinner size="sm" color="current" /> : null} Run
              </Button>
            </div>
          )}

          {active === 'search' && (
            <div className="flex gap-3 items-end shrink-0">
              <div className="flex flex-col gap-1">
                <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Column or Table Name</label>
                <input
                  className="rounded px-3 py-1.5 text-sm border w-72"
                  style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                  value={searchInput}
                  onChange={e => setSearchInput(e.target.value)}
                  placeholder="player_id, faction..."
                  onKeyDown={e => e.key === 'Enter' && run()}
                />
              </div>
              <Button onPress={run} isDisabled={loading} size="sm">
                {loading ? <Spinner size="sm" color="current" /> : null} Search
              </Button>
            </div>
          )}

          {active === 'tables' && (
            <div className="shrink-0">
              <Button onPress={run} isDisabled={loading} size="sm" variant="outline">
                {loading ? <Spinner size="sm" color="current" /> : null} List Tables
              </Button>
            </div>
          )}

          {active === 'sql' && (
            <div className="flex flex-col gap-2 shrink-0">
              <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>SQL Query</label>
              <textarea
                value={sqlInput}
                onChange={e => setSqlInput(e.target.value)}
                placeholder="SELECT * FROM dune.actors LIMIT 10;"
                rows={5}
                className="rounded px-3 py-2 text-sm font-mono border w-full resize-y"
                style={{ background: 'var(--color-surface)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
                onKeyDown={e => { if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) run() }}
              />
              <p className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Cmd/Ctrl+Enter to run. A confirmation and admin reason are required before SQL is submitted.</p>
              <div>
                <Button onPress={run} isDisabled={loading} size="sm">
                  {loading ? <Spinner size="sm" color="current" /> : null} Run Query
                </Button>
              </div>
            </div>
          )}

          {loading && (
            <div className="flex justify-center py-8 shrink-0">
              <Spinner size="lg" />
            </div>
          )}

          {error && !loading && (
            <div className="rounded-lg p-4" style={{ background: '#1a0808', border: '1px solid #5a1010', color: '#e88' }}>
              <strong>Error:</strong> {error}
            </div>
          )}

          {result && !loading && !error && (
            <div className="flex-1 min-h-0 flex flex-col">
              <ResultTable headers={result.headers} rows={result.rows} />
            </div>
          )}

          {sqlResult !== null && !loading && !error && (
            <div
              className="rounded-lg p-4 overflow-auto flex-1 min-h-0"
              style={{ background: '#0a0806', border: '1px solid #2a2418' }}
            >
              <pre className="text-sm font-mono whitespace-pre-wrap" style={{ color: '#e8dcc8', margin: 0 }}>
                {sqlResult || '(empty result)'}
              </pre>
            </div>
          )}
        </div>
      </div>
      {confirmationDialog}
    </>
  )
}
