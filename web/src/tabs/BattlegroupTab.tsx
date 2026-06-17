import { useState, useEffect, useCallback } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api, type BGHealthSection } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

type RuntimeMode = 'kubernetes' | 'docker' | string
type StatusRow = { name: string; detail1: string; status: string; detail2: string; detail3: string }
type BattlegroupView = 'pods' | 'health'

function splitRuntimeTableLine(line: string): string[] {
  const tabParts = line.trim().split(/\t+/).map(p => p.trim()).filter(Boolean)
  if (tabParts.length > 1) return tabParts
  return line.trim().split(/\s{2,}/).map(p => p.trim()).filter(Boolean)
}

function parseKubernetesStatusOutput(raw: string): StatusRow[] {
  const lines = raw.trim().split('\n').filter(Boolean)
  if (lines.length < 2) return []
  return lines.slice(1).map(line => {
    const parts = line.trim().split(/\s+/)
    return {
      name: parts[0] ?? '',
      detail1: parts[1] ?? '',
      status: parts[2] ?? '',
      detail2: parts[3] ?? '',
      detail3: parts[4] ?? '',
    }
  })
}

function parseDockerStatusOutput(raw: string): StatusRow[] {
  const lines = raw.trim().split('\n').filter(Boolean)
  if (lines.length < 2) return []
  return lines.slice(1).map(line => {
    const parts = splitRuntimeTableLine(line)
    return {
      name: parts[0] ?? '',
      detail1: parts[1] ?? '',
      status: parts[3] ?? parts[2] ?? '',
      detail2: parts[2] ?? '',
      detail3: parts.slice(4).join(' ') || '',
    }
  }).filter(row => row.name)
}

function parseRuntimeStatusOutput(runtime: RuntimeMode, raw: string): StatusRow[] {
  if (runtime === 'docker') return parseDockerStatusOutput(raw)
  return parseKubernetesStatusOutput(raw)
}

function runtimeNoun(runtime: RuntimeMode): string {
  return runtime === 'docker' ? 'Containers' : 'Pods'
}

function safeFilenamePart(value: string): string {
  return value.trim().replace(/[^a-zA-Z0-9._-]+/g, '-').replace(/^-+|-+$/g, '') || 'unknown'
}

function buildHealthBundleText(namespace: string, checkedAt: string, sections: BGHealthSection[]): string {
  const lines = [
    'Dune Admin Battlegroup Health Diagnostics',
    `Namespace: ${namespace || 'unknown'}`,
    `Checked At: ${checkedAt || new Date().toISOString()}`,
    '',
    'This support bundle is generated from protected, read-only diagnostics.',
    'Review and redact environment-specific details before sharing externally.',
    '',
  ]

  for (const section of sections) {
    lines.push('='.repeat(80))
    lines.push(section.name)
    lines.push('='.repeat(80))
    lines.push(section.description)
    lines.push('')
    lines.push(`Command: ${section.command}`)
    if (section.error) lines.push(`Error: ${section.error}`)
    lines.push('')
    lines.push(section.output || '(no output)')
    lines.push('')
  }

  return lines.join('\n')
}

function redactHealthBundleText(content: string): string {
  return content
    .replace(/\b(?:\d{1,3}\.){3}\d{1,3}\b/g, '<redacted-ipv4>')
    .replace(/\b(?:[0-9a-fA-F]{1,4}:){2,}[0-9a-fA-F]{1,4}\b/g, '<redacted-ipv6>')
    .replace(/\b[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\b/gi, '<redacted-uuid>')
    .replace(/\bip-\d+-\d+-\d+-\d+(?:\.[a-z0-9.-]+)?\b/gi, '<redacted-host>')
    .replace(/\b[a-z0-9-]+\.compute(?:-internal)?\.[a-z0-9.-]+\b/gi, '<redacted-host>')
}

function downloadTextFile(filename: string, content: string): void {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

const STATUS_COLOR: Record<string, string> = {
  Running: '#27ae60',
  Up: '#27ae60',
  Pending: '#f0a830',
  CrashLoopBackOff: '#c0392b',
  Error: '#c0392b',
  Exited: '#c0392b',
  Terminating: '#c0392b',
  Completed: '#8a7a60',
}

const ACTIONS = [
  { label: 'Start',   cmd: 'start',   danger: false, msg: 'Start the battlegroup server?' },
  { label: 'Stop',    cmd: 'stop',    danger: true,  msg: 'Stop the server? All players will be disconnected.' },
  { label: 'Restart', cmd: 'restart', danger: false, msg: 'Restart the server? Players will be briefly disconnected.' },
  { label: 'Update',  cmd: 'update',  danger: false, msg: 'Run a server update? This takes the server offline briefly.' },
  { label: 'Backup',  cmd: 'backup',  danger: false, msg: 'Create a database backup?' },
  { label: 'Restore', cmd: 'restore', danger: true,  msg: 'Restore from backup? This overwrites current data.' },
]

type ActionDef = typeof ACTIONS[0]

function battlegroupActionDetails(action: ActionDef, namespace: string): string[] {
  const details = [
    `Command: ${action.cmd}`,
    `Namespace: ${namespace || 'unknown'}`,
    action.msg,
    'Battlegroup Exec is a privileged server-control operation and is recorded in the audit log.',
  ]

  if (action.cmd === 'stop' || action.cmd === 'restart' || action.cmd === 'update' || action.cmd === 'restore') {
    details.push('This command can disrupt connected players or change active server state.')
  }
  if (action.cmd === 'restore') {
    details.push('Restore can overwrite current data. Verify the intended backup source and maintenance window before continuing.')
  }
  if (action.cmd === 'backup') {
    details.push('Verify available disk space and expected backup destination before continuing.')
  }

  return details
}

export default function BattlegroupTab() {
  const [view, setView] = useState<BattlegroupView>('pods')
  const [rows, setRows] = useState<StatusRow[]>([])
  const [runtime, setRuntime] = useState<RuntimeMode>('kubernetes')
  const [healthSections, setHealthSections] = useState<BGHealthSection[]>([])
  const [healthNamespace, setHealthNamespace] = useState('')
  const [healthCheckedAt, setHealthCheckedAt] = useState('')
  const [statusLoading, setStatusLoading] = useState(false)
  const [healthLoading, setHealthLoading] = useState(false)
  const [runningCmd, setRunningCmd] = useState<string | null>(null)
  const [cmdOutput, setCmdOutput] = useState<string | null>(null)
  const [cmdDone, setCmdDone] = useState(false)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  const fetchStatus = useCallback(async () => {
    setStatusLoading(true)
    try {
      const res = await api.battlegroup.status()
      const detectedRuntime = ((res as { runtime?: RuntimeMode }).runtime || 'kubernetes')
      setRuntime(detectedRuntime)
      setRows(parseRuntimeStatusOutput(detectedRuntime, res.output))
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      toast.danger(`Status failed: ${msg}`)
    } finally {
      setStatusLoading(false)
    }
  }, [])

  const fetchHealth = useCallback(async () => {
    setHealthLoading(true)
    try {
      const res = await api.battlegroup.health()
      const detectedRuntime = ((res as { runtime?: RuntimeMode }).runtime || runtime)
      setRuntime(detectedRuntime)
      setHealthSections(res.sections)
      setHealthNamespace(res.namespace)
      setHealthCheckedAt(res.checked_at)
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      toast.danger(`Health diagnostics failed: ${msg}`)
    } finally {
      setHealthLoading(false)
    }
  }, [runtime])

  useEffect(() => { fetchStatus() }, [fetchStatus])

  const openHealth = () => {
    setView('health')
    if (healthSections.length === 0) fetchHealth()
  }

  const exportHealthBundle = (redacted: boolean) => {
    const checkedAt = healthCheckedAt || new Date().toISOString()
    const rawContent = buildHealthBundleText(healthNamespace, checkedAt, healthSections)
    const content = redacted ? redactHealthBundleText(rawContent) : rawContent
    const suffix = redacted ? '-redacted' : ''
    const filename = `dune-admin-health-${safeFilenamePart(healthNamespace)}-${safeFilenamePart(checkedAt)}${suffix}.txt`
    downloadTextFile(filename, content)
    toast.success(redacted ? 'Redacted health diagnostics bundle exported' : 'Health diagnostics bundle exported')
  }

  const runCmd = async (action: ActionDef) => {
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: '/api/v1/battlegroup/exec',
        title: `${action.label} battlegroup server`,
        summary: action.msg,
        target: `battlegroup:${healthNamespace || runtime}`,
        details: battlegroupActionDetails(action, healthNamespace || runtime),
        confirmLabel: `Confirm ${action.label}`,
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

    setRunningCmd(action.label)
    setCmdOutput(null)
    setCmdDone(false)
    try {
      const res = await api.battlegroup.exec(action.cmd, reason)
      setCmdOutput(res.output || '(no output)')
      setCmdDone(true)
      toast.success(`${action.label} completed`)
      fetchStatus()
      if (healthSections.length > 0) fetchHealth()
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      setCmdOutput(`Error: ${msg}`)
      setCmdDone(true)
      toast.danger(`${action.label} failed: ${msg}`)
    }
  }

  const noun = runtimeNoun(runtime)

  return (
    <>
      <div style={{ display: 'flex', flexDirection: 'column', height: '100%', padding: '16px', gap: '0' }}>
        <div style={{ flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <div className="flex items-center gap-3 flex-wrap">
            <h2 className="text-base font-semibold" style={{ color: 'var(--color-primary)' }}>
              Battlegroup Status
            </h2>
            <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>runtime: {runtime}</span>
            <Button size="sm" variant={view === 'pods' ? 'primary' : 'ghost'} onPress={() => setView('pods')}>{noun}</Button>
            <Button size="sm" variant={view === 'health' ? 'primary' : 'ghost'} onPress={openHealth}>Health Diagnostics</Button>
            {view === 'pods' ? (
              <Button size="sm" variant="ghost" onPress={fetchStatus} isDisabled={statusLoading}>
                {statusLoading ? <Spinner size="sm" color="current" /> : '↻ Refresh'}
              </Button>
            ) : (
              <>
                <Button size="sm" variant="ghost" onPress={fetchHealth} isDisabled={healthLoading}>
                  {healthLoading ? <Spinner size="sm" color="current" /> : '↻ Run Diagnostics'}
                </Button>
                <Button size="sm" variant="outline" onPress={() => exportHealthBundle(false)} isDisabled={healthSections.length === 0 || healthLoading}>
                  Export Raw Bundle
                </Button>
                <Button size="sm" variant="outline" onPress={() => exportHealthBundle(true)} isDisabled={healthSections.length === 0 || healthLoading}>
                  Export Redacted Bundle
                </Button>
              </>
            )}
            {view === 'health' && healthNamespace && (
              <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>
                target: {healthNamespace}{healthCheckedAt ? ` · ${healthCheckedAt}` : ''}
              </span>
            )}
          </div>

          <div style={{ flex: 1, minHeight: 0, overflowY: 'auto' }}>
            {view === 'pods' && (
              <RuntimeStatusTable runtime={runtime} rows={rows} loading={statusLoading} />
            )}
            {view === 'health' && (
              <HealthDiagnostics sections={healthSections} loading={healthLoading} />
            )}
          </div>
        </div>

        <div
          className="shrink-0"
          style={{ borderTop: '1px solid #2a2418', paddingTop: '12px', marginTop: '12px' }}
        >
          <h2 className="text-base font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>
            Server Control
          </h2>
          <div className="flex flex-wrap gap-2">
            {ACTIONS.map(action => (
              <Button
                key={action.cmd}
                variant={action.danger ? 'danger-soft' : 'outline'}
                onPress={() => runCmd(action)}
                isDisabled={runningCmd !== null || runtime === 'docker'}
                size="sm"
              >
                {action.label}
              </Button>
            ))}
          </div>
          {runtime === 'docker' && (
            <p className="text-xs mt-2" style={{ color: 'var(--color-text-dim)' }}>
              Docker runtime detected. Read-only status and diagnostics are available; battlegroup script controls remain disabled.
            </p>
          )}
        </div>

        <Modal>
          <Modal.Backdrop
            isOpen={runningCmd !== null}
            onOpenChange={v => { if (!v && cmdDone) { setRunningCmd(null); setCmdOutput(null) } }}
          >
            <Modal.Container>
              <Modal.Dialog>
                <Modal.Header>
                  <Modal.Heading>{runningCmd ?? ''}</Modal.Heading>
                </Modal.Header>
                <Modal.Body>
                  {!cmdDone ? (
                    <div className="flex flex-col items-center gap-4 py-6">
                      <Spinner size="lg" />
                      <p className="text-sm" style={{ color: 'var(--color-text-dim)' }}>
                        Running {runningCmd?.toLowerCase() ?? ''}... this may take a moment.
                      </p>
                    </div>
                  ) : (
                    <div
                      className="rounded-lg p-3 font-mono text-xs overflow-auto max-h-60"
                      style={{ background: '#0a0806', color: '#a8d8a8', border: '1px solid #2a2418' }}
                    >
                      <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{cmdOutput}</pre>
                    </div>
                  )}
                </Modal.Body>
                {cmdDone && (
                  <Modal.Footer>
                    <Button onPress={() => { setRunningCmd(null); setCmdOutput(null) }}>
                      Close
                    </Button>
                  </Modal.Footer>
                )}
              </Modal.Dialog>
            </Modal.Container>
          </Modal.Backdrop>
        </Modal>
      </div>
      {confirmationDialog}
    </>
  )
}

function RuntimeStatusTable({ runtime, rows, loading }: { runtime: RuntimeMode; rows: StatusRow[]; loading: boolean }) {
  const noun = runtimeNoun(runtime).toLowerCase()
  const headers = runtime === 'docker'
    ? ['Name', 'ID', 'Status', 'Image', 'Ports']
    : ['Name', 'Ready', 'Status', 'Restarts', 'Age']

  if (loading && rows.length === 0) {
    return (
      <div className="flex items-center gap-2 py-4" style={{ color: 'var(--color-text-dim)' }}>
        <Spinner size="sm" color="current" />
        <span className="text-sm">Loading {noun} status...</span>
      </div>
    )
  }
  if (rows.length === 0) {
    return <p className="text-sm" style={{ color: 'var(--color-text-dim)' }}>No {noun} found. Click Refresh to try again.</p>
  }
  return (
    <div className="overflow-auto rounded-lg" style={{ border: '1px solid #2a2418' }}>
      <table className="w-full text-sm">
        <thead>
          <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
            {headers.map(h => (
              <th key={h} className="text-left px-4 py-2 font-semibold text-xs uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={`${row.name}-${i}`} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#111009' }}>
              <td className="px-4 py-2 font-mono text-xs" style={{ color: 'var(--color-text)' }}>{row.name}</td>
              <td className="px-4 py-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{row.detail1}</td>
              <td className="px-4 py-2 text-xs font-semibold" style={{ color: STATUS_COLOR[row.status.split(/\s+/)[0]] ?? 'var(--color-text)' }}>{row.status}</td>
              <td className="px-4 py-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{row.detail2}</td>
              <td className="px-4 py-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{row.detail3}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function HealthDiagnostics({ sections, loading }: { sections: BGHealthSection[]; loading: boolean }) {
  if (loading && sections.length === 0) {
    return (
      <div className="flex items-center gap-2 py-4" style={{ color: 'var(--color-text-dim)' }}>
        <Spinner size="sm" color="current" />
        <span className="text-sm">Running read-only health diagnostics...</span>
      </div>
    )
  }
  if (sections.length === 0) {
    return <p className="text-sm" style={{ color: 'var(--color-text-dim)' }}>No diagnostics loaded. Click Run Diagnostics.</p>
  }
  return (
    <div className="flex flex-col gap-3">
      {sections.map(section => (
        <HealthSectionCard key={section.name} section={section} />
      ))}
    </div>
  )
}

function HealthSectionCard({ section }: { section: BGHealthSection }) {
  return (
    <div className="rounded-lg overflow-hidden" style={{ border: '1px solid #2a2418', background: 'var(--color-surface)' }}>
      <div className="px-3 py-2" style={{ borderBottom: '1px solid #2a2418', background: '#1a1610' }}>
        <div className="flex items-center gap-2 flex-wrap">
          <span className="font-semibold text-sm" style={{ color: 'var(--color-primary)' }}>{section.name}</span>
          {section.error && <span className="text-xs" style={{ color: '#ff8c7a' }}>command returned an error</span>}
        </div>
        <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>{section.description}</p>
        <code className="text-[10px] block mt-1" style={{ color: 'var(--color-text-dim)' }}>{section.command}</code>
      </div>
      {section.error && (
        <div className="px-3 py-2 text-xs" style={{ color: '#ff8c7a', borderBottom: '1px solid #2a2418' }}>
          {section.error}
        </div>
      )}
      <pre className="p-3 text-xs overflow-auto" style={{ margin: 0, color: 'var(--color-text)', background: '#0a0806', maxHeight: '260px', whiteSpace: 'pre-wrap' }}>
        {section.output || '(no output)'}
      </pre>
    </div>
  )
}
