import { useCallback, useState, type ReactNode } from 'react'
import { api, type MutationSafetyClass } from '../api/client'

export type MutationMethod = 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export type MutationConfirmationRequest = {
  method: MutationMethod
  path: string
  title: string
  summary: string
  target?: string
  details?: string[]
  confirmLabel?: string
  cancelLabel?: string
  defaultReason?: string
  forceConfirmation?: boolean
  forceReason?: boolean
}

type PendingMutationConfirmation = {
  request: MutationConfirmationRequest
  safety: MutationSafetyClass
  requiresReason: boolean
  resolve: (reason: string | undefined) => void
  reject: (reason?: unknown) => void
}

export const mutationConfirmationCancelledMessage = 'Mutation cancelled'

function fallbackMutationSafety(method: MutationMethod, path: string): MutationSafetyClass {
  const lower = path.toLowerCase()
  const destructive = method === 'DELETE' || lower.includes('/wipe') || lower.includes('/delete') || lower.includes('/reset') || lower.includes('/kick')
  const high = destructive || lower.includes('/give') || lower.includes('/award') || lower.includes('/grant') || lower.includes('/teleport') || lower.includes('/journey/') || lower.includes('/set-faction') || lower.includes('/repair')

  return {
    action: `${method.toLowerCase()}:${path.replace(/^\//, '').replaceAll('/', '.') || 'root'}`,
    risk: destructive ? 'destructive' : high ? 'high' : 'medium',
    requires_reason: high,
    reason_enforcement_enabled: false,
    requires_preview: high,
    destructive,
    operator_warnings: high ? ['This action changes player or server state and should be performed only with a support reason.'] : [],
  }
}

function shouldShowConfirmation(request: MutationConfirmationRequest, safety: MutationSafetyClass): boolean {
  return Boolean(
    request.forceConfirmation ||
    request.forceReason ||
    safety.requires_preview ||
    safety.requires_reason ||
    safety.operator_warnings?.length ||
    safety.recommended_path ||
    safety.rollback_hint
  )
}

function riskLabel(safety: MutationSafetyClass): string {
  return safety.risk ? safety.risk.toUpperCase() : 'UNKNOWN'
}

export function useMutationConfirmation() {
  const [pending, setPending] = useState<PendingMutationConfirmation | null>(null)
  const [reason, setReason] = useState('')

  const confirmMutation = useCallback(async (request: MutationConfirmationRequest): Promise<string | undefined> => {
    const safety = await api.mutationSafety.classify(request.method, request.path)
      .catch(() => fallbackMutationSafety(request.method, request.path))

    if (!shouldShowConfirmation(request, safety)) return undefined

    const requiresReason = Boolean(request.forceReason || safety.requires_reason)

    return new Promise<string | undefined>((resolve, reject) => {
      setReason(request.defaultReason ?? '')
      setPending({ request, safety, requiresReason, resolve, reject })
    })
  }, [])

  const closePending = useCallback(() => {
    setPending(null)
    setReason('')
  }, [])

  const cancel = useCallback(() => {
    if (!pending) return
    pending.reject(new Error(mutationConfirmationCancelledMessage))
    closePending()
  }, [closePending, pending])

  const confirm = useCallback(() => {
    if (!pending) return
    const finalReason = reason.trim()
    if (pending.requiresReason && !finalReason) return
    pending.resolve(finalReason || undefined)
    closePending()
  }, [closePending, pending, reason])

  const confirmationDialog: ReactNode = pending ? (
    <div
      role="presentation"
      onClick={event => { if (event.target === event.currentTarget) cancel() }}
      style={{ position: 'fixed', inset: 0, zIndex: 200, background: 'rgba(0,0,0,0.72)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 16 }}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="mutation-confirmation-title"
        className="rounded-lg p-4"
        style={{ width: 'min(640px, 100%)', maxHeight: '88vh', overflow: 'auto', background: '#0d0b07', border: '1px solid #5a3a10', boxShadow: '0 20px 60px rgba(0,0,0,0.55)' }}
      >
        <div className="flex items-start justify-between gap-4">
          <div>
            <div className="text-xs font-semibold uppercase tracking-widest" style={{ color: pending.safety.destructive ? '#e88' : '#f0a830' }}>
              {riskLabel(pending.safety)} mutation
            </div>
            <h2 id="mutation-confirmation-title" className="text-lg font-semibold mt-1" style={{ color: 'var(--color-primary)' }}>
              {pending.request.title}
            </h2>
          </div>
          <button
            type="button"
            aria-label="Cancel mutation"
            onClick={cancel}
            style={{ background: 'transparent', border: 'none', color: 'var(--color-text-dim)', cursor: 'pointer', fontSize: 18, lineHeight: 1 }}
          >
            ×
          </button>
        </div>

        <p className="text-sm mt-3" style={{ color: 'var(--color-text)' }}>
          {pending.request.summary}
        </p>

        {pending.request.target && (
          <div className="mt-3 rounded p-2 text-xs font-mono" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
            Target: {pending.request.target}
          </div>
        )}

        {pending.request.details && pending.request.details.length > 0 && (
          <ul className="mt-3 pl-5 text-xs" style={{ color: 'var(--color-text-dim)', listStyle: 'disc' }}>
            {pending.request.details.map((detail, index) => <li key={`${detail}-${index}`}>{detail}</li>)}
          </ul>
        )}

        {(pending.safety.operator_warnings?.length || pending.safety.recommended_path || pending.safety.rollback_hint) && (
          <div className="mt-3 rounded p-3 text-xs" style={{ background: '#1a1200', border: '1px solid #5a3a10', color: '#f0a830' }}>
            {pending.safety.operator_warnings?.map((warning, index) => <div key={`${warning}-${index}`}>Warning: {warning}</div>)}
            {pending.safety.recommended_path && <div>Recommended path: {pending.safety.recommended_path}</div>}
            {pending.safety.rollback_hint && <div>Rollback hint: {pending.safety.rollback_hint}</div>}
          </div>
        )}

        <div className="mt-3 grid grid-cols-1 md:grid-cols-3 gap-2 text-xs">
          <div className="rounded p-2" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
            Action: <span className="font-mono" style={{ color: 'var(--color-text)' }}>{pending.safety.action}</span>
          </div>
          <div className="rounded p-2" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
            Reason required: <span style={{ color: pending.requiresReason ? '#f0a830' : 'var(--color-text)' }}>{pending.requiresReason ? 'Yes' : 'No'}</span>
          </div>
          <div className="rounded p-2" style={{ background: '#0a0806', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
            Backend enforcement: <span style={{ color: pending.safety.reason_enforcement_enabled ? '#f0a830' : 'var(--color-text)' }}>{pending.safety.reason_enforcement_enabled ? 'Enabled' : 'Disabled'}</span>
          </div>
        </div>

        {!pending.safety.reason_enforcement_enabled && pending.requiresReason && (
          <div className="mt-3 rounded p-2 text-xs" style={{ background: '#160f06', border: '1px solid #5a3a10', color: '#f0a830' }}>
            Frontend reason capture is active, but backend reason enforcement is disabled. Set ADMIN_REQUIRE_REASON=true to reject high-risk requests that bypass the UI without a reason.
          </div>
        )}

        <label className="block mt-4 text-xs font-semibold" style={{ color: 'var(--color-text-dim)' }}>
          Admin reason {pending.requiresReason ? '(required)' : '(optional)'}
        </label>
        <textarea
          value={reason}
          onChange={event => setReason(event.target.value)}
          rows={4}
          placeholder="Example: Support ticket, player-reported issue, or operator remediation reason."
          className="mt-1 w-full rounded px-3 py-2 text-sm border"
          style={{ background: '#0a0806', color: 'var(--color-text)', borderColor: pending.requiresReason && !reason.trim() ? '#8a4a18' : '#2a2418', outline: 'none', resize: 'vertical' }}
        />

        <div className="mt-4 flex justify-end gap-2">
          <button
            type="button"
            onClick={cancel}
            className="rounded px-3 py-2 text-sm"
            style={{ background: 'transparent', border: '1px solid #2a2418', color: 'var(--color-text-dim)', cursor: 'pointer' }}
          >
            {pending.request.cancelLabel ?? 'Cancel'}
          </button>
          <button
            type="button"
            onClick={confirm}
            disabled={pending.requiresReason && !reason.trim()}
            className="rounded px-3 py-2 text-sm font-semibold"
            style={{ background: pending.requiresReason && !reason.trim() ? '#3a3020' : 'var(--color-primary)', border: 'none', color: '#fff', cursor: pending.requiresReason && !reason.trim() ? 'not-allowed' : 'pointer' }}
          >
            {pending.request.confirmLabel ?? 'Confirm mutation'}
          </button>
        </div>
      </div>
    </div>
  ) : null

  return { confirmMutation, confirmationDialog }
}
