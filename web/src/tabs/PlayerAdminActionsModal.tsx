import { Button, Modal, Spinner, toast } from '@heroui/react'
import { useState } from 'react'
import { api } from '../api/client'
import type { Player } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

type AdminMutation = {
  path: string
  title: string
  summary: string
  details: string[]
  confirmLabel: string
  successLabel: string
  run: (reason: string) => Promise<unknown>
}

export default function PlayerAdminActionsModal({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [busyAction, setBusyAction] = useState<string | null>(null)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  const target = `actor:${player.id} · ctrl:${player.controller_id} · acct:${player.account_id}`

  const runConfirmed = async (mutation: AdminMutation) => {
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: mutation.path,
        title: mutation.title,
        summary: mutation.summary,
        target,
        details: [
          `Player: ${player.name}`,
          `Online state: ${player.online_status || 'Offline'}`,
          ...mutation.details,
        ],
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

    setBusyAction(mutation.path)
    try {
      await mutation.run(reason)
      toast.success(mutation.successLabel)
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setBusyAction(null)
    }
  }

  const actions: AdminMutation[] = [
    {
      path: '/api/v1/players/journey/wipe',
      title: 'Clear all journey progress',
      summary: `Clear all journey progress for ${player.name}.`,
      details: ['This affects every journey node for the selected account.', 'Use only for support remediation after verifying the player target.'],
      confirmLabel: 'Clear journey',
      successLabel: `Cleared journey progress for ${player.name}`,
      run: reason => api.players.journeyWipe(player.account_id, reason),
    },
    {
      path: '/api/v1/players/delete-tutorials',
      title: 'Remove tutorial records',
      summary: `Remove tutorial completion records for ${player.name}.`,
      details: ['This removes tutorial completion state for the selected character.', 'Use only when tutorial state is known to be corrupt or blocking progression.'],
      confirmLabel: 'Remove tutorials',
      successLabel: `Removed tutorial records for ${player.name}`,
      run: reason => api.players.deleteTutorials(player.id, reason),
    },
    {
      path: '/api/v1/players/wipe-codex',
      title: 'Clear codex discoveries',
      summary: `Clear codex discoveries for ${player.name}.`,
      details: ['This affects codex discovery records for the selected account.', 'Use only after confirming the account target and support reason.'],
      confirmLabel: 'Clear codex',
      successLabel: `Cleared codex discoveries for ${player.name}`,
      run: reason => api.players.wipeCodex(player.account_id, reason),
    },
    {
      path: '/api/v1/players/kick',
      title: 'Disconnect player session',
      summary: `Disconnect ${player.name} from the server.`,
      details: ['This forcibly disconnects the selected player session.', 'Use when an operator needs the player offline before a state-changing support action.'],
      confirmLabel: 'Disconnect player',
      successLabel: `Disconnected ${player.name}`,
      run: reason => api.players.kick(player.id, reason),
    },
  ]

  return (
    <>
      <Modal>
        <Modal.Backdrop isOpen={open} onOpenChange={value => !value && onClose()}>
          <Modal.Container>
            <Modal.Dialog>
              <Modal.CloseTrigger />
              <Modal.Header><Modal.Heading>Admin Actions — {player.name}</Modal.Heading></Modal.Header>
              <Modal.Body>
                <div className="flex flex-col gap-3">
                  <div className="rounded p-3 text-xs" style={{ background: '#0f0d09', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                    <div><strong style={{ color: 'var(--color-primary)' }}>Target:</strong> {target}</div>
                    <div className="mt-1"><strong style={{ color: 'var(--color-primary)' }}>Current map:</strong> {player.map || 'unknown'}</div>
                    <div className="mt-1"><strong style={{ color: 'var(--color-primary)' }}>Online state:</strong> {player.online_status || 'Offline'}</div>
                  </div>

                  <div className="rounded p-3 text-xs" style={{ background: '#1a1200', border: '1px solid #c9820a', color: '#f0a830' }}>
                    These actions modify or disrupt player state. Verify the selected player, account, and support ticket before continuing. A reason is required and will be recorded in the audit log.
                  </div>

                  <div className="flex flex-col" style={{ border: '1px solid #2a2418', borderRadius: 8, overflow: 'hidden' }}>
                    {actions.map((action, index) => (
                      <div key={action.path} className="flex items-center gap-3 p-3" style={{ background: index % 2 === 0 ? '#0d0b07' : '#0f0d09', borderBottom: index === actions.length - 1 ? 'none' : '1px solid #1a1610' }}>
                        <div className="flex-1">
                          <div className="text-sm font-semibold" style={{ color: 'var(--color-text)' }}>{action.title}</div>
                          <div className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>{action.summary}</div>
                        </div>
                        <Button size="sm" variant="danger-soft" isDisabled={Boolean(busyAction)} onPress={() => runConfirmed(action)}>
                          {busyAction === action.path ? <Spinner size="sm" color="current" /> : null}
                          {action.confirmLabel}
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
              </Modal.Body>
              <Modal.Footer>
                <Button variant="tertiary" onPress={onClose}>Close</Button>
              </Modal.Footer>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
      {confirmationDialog}
    </>
  )
}
