import { useEffect, useMemo, useState, type ReactNode } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { Player, SpecTrack } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

type ActionSection = 'resources' | 'specs'

const ACTION_SECTIONS: { key: ActionSection; label: string }[] = [
  { key: 'resources', label: 'Stats' },
  { key: 'specs', label: 'Specs' },
]

const XP_TRACKS = ['Combat', 'Crafting', 'Gathering', 'Exploration', 'Sabotage']
const FACTIONS = [{ id: 1, name: 'Atreides' }, { id: 2, name: 'Harkonnen' }, { id: 3, name: 'None' }, { id: 4, name: 'Smuggler' }]

type ConfirmedMutation = {
  path: string
  title: string
  summary: string
  details?: string[]
  confirmLabel: string
  successLabel: string
  run: (reason: string) => Promise<unknown>
  after?: () => void | Promise<void>
}

function inputStyle() {
  return { background: '#0d0b07', color: 'var(--color-text)', borderColor: '#3a3020', outline: 'none' }
}

export default function PlayerActionsModalConfirmed({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [section, setSection] = useState<ActionSection>('resources')
  const [busy, setBusy] = useState(false)
  const [currency, setCurrency] = useState(100)
  const [scrip, setScrip] = useState(100)
  const [intel, setIntel] = useState(100)
  const [charXP, setCharXP] = useState(1000)
  const [charXPCurrent, setCharXPCurrent] = useState<{xp: number; level: number} | null>(null)
  const [factionId, setFactionId] = useState(player.faction_id || 1)
  const [repDelta, setRepDelta] = useState(100)
  const [playerSpecs, setPlayerSpecs] = useState<SpecTrack[]>([])
  const [specsLoaded, setSpecsLoaded] = useState(false)
  const [specsLoading, setSpecsLoading] = useState(false)
  const [specXPInputs, setSpecXPInputs] = useState<Record<string, number>>({})
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  useEffect(() => {
    if (!open) {
      setSection('resources')
      setPlayerSpecs([])
      setSpecsLoaded(false)
      setSpecXPInputs({})
      setCharXPCurrent(null)
      return
    }
    setFactionId(player.faction_id > 0 ? player.faction_id : 1)
    api.players.charXPCurrent(player.id).then(setCharXPCurrent).catch(() => {})
  }, [open, player.faction_id, player.id])

  useEffect(() => {
    if (section !== 'specs' || specsLoaded || !open) return
    setSpecsLoading(true)
    api.players.specs_for(player.controller_id)
      .then(rows => {
        setPlayerSpecs(rows)
        setSpecsLoaded(true)
        const inputs: Record<string, number> = {}
        XP_TRACKS.forEach(track => {
          const found = rows.find(row => row.track_type === track)
          inputs[track] = found ? found.xp : 0
        })
        setSpecXPInputs(inputs)
      })
      .catch((e: unknown) => toast.danger(e instanceof Error ? e.message : String(e)))
      .finally(() => setSpecsLoading(false))
  }, [section, specsLoaded, open, player.controller_id])

  const target = `actor:${player.id} · ctrl:${player.controller_id} · acct:${player.account_id}`

  const runConfirmed = async (mutation: ConfirmedMutation) => {
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: mutation.path,
        title: mutation.title,
        summary: mutation.summary,
        target,
        details: [`Player: ${player.name}`, `Online state: ${player.online_status || 'Offline'}`, ...(mutation.details ?? [])],
        confirmLabel: mutation.confirmLabel,
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

    setBusy(true)
    try {
      await mutation.run(reason)
      await mutation.after?.()
      toast.success(mutation.successLabel)
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setBusy(false)
    }
  }

  const numInput = (val: number, set: (v: number) => void, min = 1, max = 9999999) => (
    <input
      type="number"
      min={min}
      max={max}
      value={val}
      onChange={event => set(Math.max(min, Math.min(max, parseInt(event.target.value) || min)))}
      className="rounded px-2 py-1.5 text-sm border w-28"
      style={inputStyle()}
    />
  )

  const actionRow = (label: string, inputs: ReactNode, btnLabel: string, mutation: ConfirmedMutation) => (
    <div className="flex items-end gap-3 py-3" style={{ borderBottom: '1px solid #1a1610' }}>
      <div className="w-36 shrink-0 text-sm" style={{ color: 'var(--color-text-dim)' }}>{label}</div>
      <div className="flex items-end gap-2 flex-1 flex-wrap">{inputs}</div>
      <Button size="sm" variant="ghost" onPress={() => runConfirmed(mutation)} isDisabled={busy}>{btnLabel}</Button>
    </div>
  )

  const onlineWarning = player.online_status === 'Online' ? (
    <div className="text-xs px-3 py-2 rounded mb-3" style={{ background: '#1a1200', border: '1px solid #c9820a', color: '#f0a830' }}>
      Player is online. Some changes may not take effect until reconnect, and unsafe mutations can cause client/server state drift.
    </div>
  ) : null

  const specRows = useMemo(() => XP_TRACKS.map(track => {
    const found = playerSpecs.find(row => row.track_type === track)
    return { track, currentXP: found ? found.xp : 0, inputVal: specXPInputs[track] ?? (found ? found.xp : 0) }
  }), [playerSpecs, specXPInputs])

  return (
    <>
      <Modal>
        <Modal.Backdrop isOpen={open} onOpenChange={value => !value && onClose()}>
          <Modal.Container size="full">
            <Modal.Dialog style={{ height: '92vh', display: 'flex', flexDirection: 'column' }}>
              <Modal.CloseTrigger />
              <Modal.Header>
                <Modal.Heading>
                  {player.name}
                  <span className="ml-2 text-xs font-mono font-normal" style={{ color: 'var(--color-text-dim)' }}>{target}</span>
                </Modal.Heading>
              </Modal.Header>
              <Modal.Body style={{ display: 'flex', gap: 0, overflow: 'hidden', padding: 0, flex: 1 }}>
                <div className="shrink-0 flex flex-col gap-1 p-3" style={{ borderRight: '1px solid #2a2418', background: '#0d0b07', minWidth: 120 }}>
                  {ACTION_SECTIONS.map(item => (
                    <button key={item.key} onClick={() => setSection(item.key)} className="text-left px-3 py-2 rounded text-sm transition-colors" style={{ background: section === item.key ? 'var(--color-primary)' : 'transparent', color: section === item.key ? '#fff' : 'var(--color-text)' }}>
                      {item.label}
                    </button>
                  ))}
                </div>

                <div className="flex-1 overflow-hidden flex flex-col p-4">
                  {section === 'resources' && (
                    <div className="overflow-y-auto flex-1 flex flex-col">
                      <div className="text-xs font-semibold uppercase tracking-widest px-1 py-2" style={{ color: 'var(--color-primary)', borderBottom: '1px solid #2a2418' }}>Currency &amp; Resources</div>
                      {actionRow('Give Currency', numInput(currency, setCurrency), 'Give', { path: '/api/v1/players/give-currency', title: 'Give currency', summary: `Give ${currency} Solari to ${player.name}.`, details: [`Amount: ${currency}`], confirmLabel: 'Give currency', successLabel: `Gave ${currency} Solari to ${player.name}`, run: reason => api.players.giveCurrency(player.controller_id, currency, reason) })}
                      {actionRow('Give Scrip', numInput(scrip, setScrip), 'Give', { path: '/api/v1/players/give-scrip', title: 'Give faction scrip', summary: `Give ${scrip} faction scrip to ${player.name}.`, details: [`Amount: ${scrip}`], confirmLabel: 'Give scrip', successLabel: `Gave ${scrip} scrip to ${player.name}`, run: reason => api.players.giveScrip(player.controller_id, scrip, reason) })}
                      {actionRow('Award Intel', numInput(intel, setIntel), 'Award', { path: '/api/v1/players/award-intel', title: 'Award intel', summary: `Award ${intel} intel to ${player.name}.`, details: [`Amount: ${intel}`], confirmLabel: 'Award intel', successLabel: `Awarded ${intel} intel to ${player.name}`, run: reason => api.players.awardIntel(player.controller_id, intel, reason) })}

                      <div className="text-xs font-semibold uppercase tracking-widest px-1 py-2 mt-4" style={{ color: 'var(--color-primary)', borderBottom: '1px solid #2a2418' }}>Character XP</div>
                      {onlineWarning}
                      {charXPCurrent && <div className="px-1 py-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>Current: <span style={{ color: 'var(--color-text)' }}>{charXPCurrent.xp.toLocaleString()} XP</span> — Level <span style={{ color: 'var(--color-text)' }}>{charXPCurrent.level}</span><span style={{ color: '#555' }}> / 200</span></div>}
                      {actionRow('Award Char XP', <div className="flex flex-col gap-0.5">{numInput(charXP, setCharXP, 0, 344440)}<span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Max 344,440</span></div>, 'Award', { path: '/api/v1/players/award-char-xp', title: 'Award character XP', summary: `Award ${charXP} character XP to ${player.name}.`, details: [`Amount: ${charXP}`, 'Online players may need to reconnect before XP state refreshes.'], confirmLabel: 'Award XP', successLabel: `Awarded ${charXP} character XP to ${player.name}`, run: reason => api.players.awardCharXP(player.id, charXP, reason), after: () => api.players.charXPCurrent(player.id).then(setCharXPCurrent).catch(() => {}) })}

                      <div className="text-xs font-semibold uppercase tracking-widest px-1 py-2 mt-4" style={{ color: 'var(--color-primary)', borderBottom: '1px solid #2a2418' }}>Faction Reputation</div>
                      <div className="flex items-center gap-2 py-3" style={{ borderBottom: '1px solid #1a1610' }}>
                        <div className="w-36 shrink-0 text-sm" style={{ color: 'var(--color-text-dim)' }}>Faction</div>
                        <select value={factionId} onChange={event => setFactionId(Number(event.target.value))} className="rounded px-2 py-1.5 text-sm border w-40" style={inputStyle()}>
                          {FACTIONS.map(faction => <option key={faction.id} value={faction.id}>{faction.name}</option>)}
                        </select>
                      </div>
                      {actionRow('Reputation', <div className="flex flex-col gap-0.5">{numInput(repDelta, setRepDelta, 0, 12474)}<span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Adds to current, max 12,474</span></div>, 'Give', { path: '/api/v1/players/give-faction-rep', title: 'Give faction reputation', summary: `Give ${repDelta} reputation to ${player.name}.`, details: [`Faction ID: ${factionId}`, `Amount: ${repDelta}`], confirmLabel: 'Give reputation', successLabel: `Gave ${repDelta} reputation to ${player.name}`, run: reason => api.players.giveFactionRep(player.controller_id, factionId, repDelta, reason) })}
                    </div>
                  )}

                  {section === 'specs' && (
                    <div className="overflow-y-auto flex-1 flex flex-col gap-3">
                      {onlineWarning}
                      <div className="flex items-center gap-2 shrink-0"><span className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Specializations</span><Button size="sm" variant="ghost" isDisabled={specsLoading} onPress={() => setSpecsLoaded(false)}>{specsLoading ? <Spinner size="sm" color="current" /> : 'Refresh'}</Button></div>
                      {specsLoading ? <div className="flex justify-center py-8"><Spinner size="lg" /></div> : (
                        <div className="overflow-auto rounded-lg" style={{ border: '1px solid #2a2418', flex: 1, minHeight: 0 }}>
                          <table className="w-full text-xs">
                            <thead><tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>{['Track', 'Current XP', 'Set XP', ''].map((h, idx) => <th key={idx} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>)}</tr></thead>
                            <tbody>{specRows.map((row, i) => <tr key={row.track} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}><td className="px-3 py-2 font-semibold" style={{ color: 'var(--color-text)' }}>{row.track}</td><td className="px-3 py-2 font-mono" style={{ color: 'var(--color-text-dim)' }}>{row.currentXP.toLocaleString()}</td><td className="px-3 py-2"><input type="number" min={0} max={44182} value={row.inputVal} onChange={event => setSpecXPInputs(prev => ({ ...prev, [row.track]: Math.max(0, Math.min(44182, parseInt(event.target.value) || 0)) }))} className="rounded px-2 py-1 text-sm border w-28" style={inputStyle()} /></td><td className="px-3 py-2"><Button size="sm" variant="ghost" isDisabled={busy} onPress={() => runConfirmed({ path: '/api/v1/players/set-spec-xp', title: 'Set specialization XP', summary: `Set ${row.track} XP to ${row.inputVal} for ${player.name}.`, details: [`Track: ${row.track}`, `XP: ${row.inputVal}`], confirmLabel: 'Set XP', successLabel: `Set ${row.track} XP to ${row.inputVal}`, run: reason => api.players.setSpecXP(player.controller_id, row.track, row.inputVal, reason), after: () => setPlayerSpecs(prev => prev.find(spec => spec.track_type === row.track) ? prev.map(spec => spec.track_type === row.track ? { ...spec, xp: row.inputVal } : spec) : [...prev, { player_id: player.controller_id, track_type: row.track, xp: row.inputVal, level: 0 }]) })}>Set</Button></td></tr>)}</tbody>
                          </table>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </Modal.Body>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
      {confirmationDialog}
    </>
  )
}
