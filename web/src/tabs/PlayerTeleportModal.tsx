import { useEffect, useState } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { Player, TeleportLocation } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

const inputStyle = {
  background: '#0d0b07',
  color: 'var(--color-text)',
  borderColor: '#3a3020',
  outline: 'none',
}

export default function PlayerTeleportModal({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [partitions, setPartitions] = useState<TeleportLocation[]>([])
  const [selectedPartition, setSelectedPartition] = useState('')
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  useEffect(() => {
    if (!open) {
      setSelectedPartition('')
      return
    }

    setLoading(true)
    api.players.partitions()
      .then(rows => {
        setPartitions(rows)
        setSelectedPartition(rows[0]?.name ?? '')
      })
      .catch((e: unknown) => toast.danger(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false))
  }, [open])

  const handleTeleport = async () => {
    if (!selectedPartition) {
      toast.warning('Select a destination first')
      return
    }

    if (player.online_status === 'Online') {
      toast.warning('Teleport is disabled while the player is online')
      return
    }

    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: '/api/v1/players/teleport',
        title: 'Move player to selected destination',
        summary: `Move ${player.name} to ${selectedPartition}.`,
        target: `actor:${player.id} · ctrl:${player.controller_id} · acct:${player.account_id} · fls:${player.fls_id}`,
        details: [
          `Destination: ${selectedPartition}`,
          `Current map: ${player.map || 'unknown'}`,
          `Online state: ${player.online_status || 'Offline'}`,
          'This action should be used only when the player is offline to avoid client/server location drift.',
        ],
        confirmLabel: 'Move player',
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
      await api.players.teleport(player.fls_id, selectedPartition, reason)
      toast.success(`Moved ${player.name} to ${selectedPartition}`)
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
        <Modal.Backdrop isOpen={open} onOpenChange={value => !value && onClose()}>
          <Modal.Container>
            <Modal.Dialog>
              <Modal.CloseTrigger />
              <Modal.Header><Modal.Heading>Move Player — {player.name}</Modal.Heading></Modal.Header>
              <Modal.Body>
                <div className="flex flex-col gap-3">
                  <div className="rounded p-3 text-xs" style={{ background: '#0f0d09', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
                    <div><strong style={{ color: 'var(--color-primary)' }}>Target:</strong> actor:{player.id} · ctrl:{player.controller_id} · acct:{player.account_id}</div>
                    <div className="mt-1"><strong style={{ color: 'var(--color-primary)' }}>Current map:</strong> {player.map || 'unknown'}</div>
                    <div className="mt-1"><strong style={{ color: 'var(--color-primary)' }}>Online state:</strong> {player.online_status || 'Offline'}</div>
                  </div>

                  {player.online_status === 'Online' && (
                    <div className="rounded p-3 text-xs" style={{ background: '#1a1200', border: '1px solid #c9820a', color: '#f0a830' }}>
                      Move is disabled while the player is online. Have the player disconnect first to reduce location/state drift risk.
                    </div>
                  )}

                  {loading ? (
                    <div className="flex justify-center py-6"><Spinner size="lg" /></div>
                  ) : (
                    <label className="flex flex-col gap-1 text-sm" style={{ color: 'var(--color-text-dim)' }}>
                      Destination
                      <select
                        value={selectedPartition}
                        onChange={event => setSelectedPartition(event.target.value)}
                        className="rounded px-2 py-2 text-sm border"
                        style={inputStyle}
                      >
                        {partitions.length === 0 && <option value="">No destinations available</option>}
                        {partitions.map(partition => <option key={partition.name} value={partition.name}>{partition.name}</option>)}
                      </select>
                    </label>
                  )}
                </div>
              </Modal.Body>
              <Modal.Footer>
                <Button variant="tertiary" onPress={onClose}>Cancel</Button>
                <Button onPress={handleTeleport} isDisabled={loading || submitting || !selectedPartition || player.online_status === 'Online'}>
                  {submitting ? <Spinner size="sm" color="current" /> : null}
                  Move Player
                </Button>
              </Modal.Footer>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
      {confirmationDialog}
    </>
  )
}
