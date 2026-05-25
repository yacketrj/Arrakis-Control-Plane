import { useMemo, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { auditApi, type AdminAuditEvent } from '../api/audit'

function riskColor(risk?: string): string {
  switch (risk) {
    case 'destructive': return '#ff8c7a'
    case 'high': return '#ffb86b'
    case 'medium': return 'var(--color-primary)'
    default: return 'var(--color-text-dim)'
  }
}

function targetSummary(target?: Record<string, string>): string {
  if (!target || Object.keys(target).length === 0) return ''
  return Object.entries(target).map(([key, value]) => `${key}:${value}`).join(' ')
}

function safetySummary(event: AdminAuditEvent): string {
  const parts: string[] = []
  if (event.requires_reason) parts.push('reason')
  if (event.requires_preview) parts.push('preview')
  if (event.destructive) parts.push('destructive')
  return parts.join(', ')
}

function detailsTitle(event: AdminAuditEvent): string {
  return [
    event.rollback_hint ? `Rollback: ${event.rollback_hint}` : '',
    event.recommended_path ? `Recommended: ${event.recommended_path}` : '',
    ...(event.operator_warnings ?? []).map(w => `Warning: ${w}`),
  ].filter(Boolean).join('\n')
}

export default function AuditTab() {
  const [events, setEvents] = useState<AdminAuditEvent[]>([])
  const [filter, setFilter] = useState('')
  const [loading, setLoading] = useState(false)

  const load = async () => {
    setLoading(true)
    try {
      setEvents(await auditApi.events())
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  const filtered = useMemo(() => {
    const q = filter.toLowerCase().trim()
    if (!q) return events
    return events.filter(event => [
      event.timestamp,
      event.method,
      event.path,
      event.action,
      event.result,
      event.risk ?? '',
      event.reason ?? '',
      event.rollback_hint ?? '',
      event.recommended_path ?? '',
      safetySummary(event),
      ...(event.operator_warnings ?? []),
      targetSummary(event.target),
      String(event.status),
    ].some(value => value.toLowerCase().includes(q)))
  }, [events, filter])

  return (
    <div className="flex flex-col gap-4 h-full min-h-0">
      <div className="rounded-lg p-4 shrink-0" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h2 className="text-base font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Admin Action Audit Log</h2>
        <p className="text-sm mb-4" style={{ color: 'var(--color-text-dim)' }}>
          Review protected mutating admin requests captured by the append-only audit log. Mutation Safety metadata flags required reasons, preview expectations, destructive actions, warnings, and rollback hints without storing raw request bodies or secrets.
        </p>
        <div className="flex gap-3 items-end flex-wrap">
          <div className="flex flex-col gap-1">
            <label className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Filter</label>
            <input className="rounded px-3 py-1.5 text-sm border w-80" style={{ background: 'var(--color-background)', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }} value={filter} onChange={e => setFilter(e.target.value)} placeholder="give-item, rollback, preview, destructive..." />
          </div>
          <Button onPress={load} isDisabled={loading} size="sm">{loading ? <Spinner size="sm" color="current" /> : null} Refresh</Button>
          <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{filtered.length} shown / {events.length} loaded</span>
        </div>
      </div>

      <div className="rounded-lg overflow-auto flex-1 min-h-0" style={{ border: '1px solid #2a2418' }}>
        <table className="w-full text-xs">
          <thead>
            <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
              {['Time', 'Result', 'Risk', 'Safety', 'Method', 'Action', 'Target', 'Reason', 'Details', 'Status', 'Duration'].map(h => <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide whitespace-nowrap" style={{ color: 'var(--color-primary)' }}>{h}</th>)}
            </tr>
          </thead>
          <tbody>
            {filtered.map((event, i) => {
              const detailText = detailsTitle(event)
              return (
                <tr key={`${event.timestamp}-${event.action}-${i}`} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text-dim)' }}>{event.timestamp}</td>
                  <td className="px-3 py-2 whitespace-nowrap" style={{ color: event.result === 'success' ? 'var(--color-success)' : '#ff8c7a' }}>{event.result}</td>
                  <td className="px-3 py-2 whitespace-nowrap" style={{ color: riskColor(event.risk) }}>{event.risk || 'unknown'}</td>
                  <td className="px-3 py-2 whitespace-nowrap" style={{ color: event.destructive ? '#ff8c7a' : 'var(--color-text-dim)' }}>{safetySummary(event) || '-'}</td>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text)' }}>{event.method}</td>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text)' }}>{event.action}</td>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text-dim)' }}>{targetSummary(event.target) || '-'}</td>
                  <td className="px-3 py-2 max-w-md truncate" title={event.reason || ''} style={{ color: 'var(--color-text-dim)' }}>{event.reason || '-'}</td>
                  <td className="px-3 py-2 max-w-lg truncate" title={detailText} style={{ color: detailText ? 'var(--color-text-dim)' : '#5a5246' }}>{detailText ? 'hover for details' : '-'}</td>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text)' }}>{event.status}</td>
                  <td className="px-3 py-2 font-mono whitespace-nowrap" style={{ color: 'var(--color-text-dim)' }}>{event.duration_ms} ms</td>
                </tr>
              )
            })}
            {filtered.length === 0 && <tr><td colSpan={11} className="px-3 py-8 text-center" style={{ color: 'var(--color-text-dim)' }}>{events.length === 0 ? 'No audit events loaded. Click Refresh.' : 'No events match the filter.'}</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
