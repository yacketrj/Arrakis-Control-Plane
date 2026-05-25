import { useEffect, useState } from 'react'
import { Button, Modal, Spinner, toast } from '@heroui/react'
import { api } from '../api/client'
import type { InventoryItem, Player, VehicleRow } from '../api/client'
import { mutationConfirmationCancelledMessage, useMutationConfirmation } from '../hooks/useMutationConfirmation'

export default function InventoryModal({ player, open, onClose }: { player: Player; open: boolean; onClose: () => void }) {
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(false)
  const [vehicles, setVehicles] = useState<VehicleRow[]>([])
  const [vehiclesLoading, setVehiclesLoading] = useState(false)
  const { confirmMutation, confirmationDialog } = useMutationConfirmation()

  useEffect(() => {
    if (!open) {
      setVehicles([])
      return
    }
    setLoading(true)
    setVehiclesLoading(true)
    api.players.inventory(player.id)
      .then(setItems)
      .catch((e: unknown) => toast.danger(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false))
    api.players.vehicles(player.account_id)
      .then(setVehicles)
      .catch(() => {})
      .finally(() => setVehiclesLoading(false))
  }, [open, player.id, player.account_id])

  const itemLabel = (item: InventoryItem): string => item.name || item.template_id || `item:${item.id}`

  const targetLabel = (item: InventoryItem): string => [
    `actor:${player.id}`,
    `ctrl:${player.controller_id}`,
    `acct:${player.account_id}`,
    `item:${item.id}`,
  ].join(' · ')

  const handleDelete = async (item: InventoryItem) => {
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'DELETE',
        path: `/api/v1/players/item/${item.id}`,
        title: 'Delete inventory item',
        summary: `Delete ${itemLabel(item)} from ${player.name}'s inventory.`,
        target: targetLabel(item),
        details: [
          `Template: ${item.template_id}`,
          `Stack size: ${item.stack_size}`,
          `Quality: ${item.quality}`,
          player.online_status === 'Online'
            ? 'Player is online. Deleting items while online can cause inventory desyncs.'
            : 'Player is offline.',
        ],
        confirmLabel: 'Delete item',
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

    try {
      await api.players.deleteItem(item.id, reason)
      setItems(prev => prev.filter(i => i.id !== item.id))
      toast.success('Item deleted')
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  const handleRepair = async (item: InventoryItem) => {
    let reason: string | undefined
    try {
      reason = await confirmMutation({
        method: 'POST',
        path: '/api/v1/players/repair-item',
        title: 'Repair inventory item',
        summary: `Repair ${itemLabel(item)} in ${player.name}'s inventory.`,
        target: targetLabel(item),
        details: [
          `Template: ${item.template_id}`,
          `Current durability: ${item.durability}`,
          `Max durability: ${item.max_durability}`,
        ],
        confirmLabel: 'Repair item',
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

    try {
      await api.players.repairItem(item.id, reason)
      setItems(prev => prev.map(i => i.id === item.id ? { ...i, durability: i.max_durability } : i))
      toast.success(`Repaired ${itemLabel(item)}`)
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  return (
    <>
      <Modal>
        <Modal.Backdrop isOpen={open} onOpenChange={v => !v && onClose()}>
          <Modal.Container size="full">
            <Modal.Dialog>
              <Modal.CloseTrigger />
              <Modal.Header><Modal.Heading>{player.name} — Inventory</Modal.Heading></Modal.Header>
              <Modal.Body>
                {loading ? (
                  <div className="flex justify-center py-8"><Spinner size="lg" /></div>
                ) : items.length === 0 ? (
                  <p style={{ color: 'var(--color-text-dim)' }}>No items found.</p>
                ) : (
                  <>
                    <div className="overflow-auto rounded-lg" style={{ border: '1px solid #2a2418', maxHeight: '55vh', flex: 1, minHeight: 0 }}>
                      <table className="w-full text-xs">
                        <thead>
                          <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
                            {['Template', 'Stack', 'Quality', 'Durability', ''].map(h => (
                              <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {items.map((item, i) => (
                            <tr key={item.id} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                              <td className="px-3 py-1.5">
                                <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{item.name || item.template_id}</div>
                                {item.name && (
                                  <div className="text-xs font-mono" style={{ color: 'var(--color-text-dim)' }}>{item.template_id}</div>
                                )}
                              </td>
                              <td className="px-3 py-1.5" style={{ color: 'var(--color-text-dim)' }}>{item.stack_size}</td>
                              <td className="px-3 py-1.5" style={{ color: 'var(--color-text-dim)' }}>{item.quality}</td>
                              <td className="px-3 py-1.5" style={{ color: 'var(--color-text-dim)' }}>
                                {item.durability} / {item.max_durability}
                              </td>
                              <td className="px-3 py-1.5">
                                <div className="flex gap-1">
                                  {item.max_durability !== 'N/A' && (
                                    <Button size="sm" variant="ghost" onPress={() => handleRepair(item)}>Repair</Button>
                                  )}
                                  <Button size="sm" variant="danger-soft" onPress={() => handleDelete(item)}>X</Button>
                                </div>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>

                    <div className="mt-4">
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-sm font-semibold" style={{ color: 'var(--color-primary)' }}>Vehicles</span>
                        {vehiclesLoading && <Spinner size="sm" color="current" />}
                      </div>
                      <div className="overflow-auto rounded-lg" style={{ border: '1px solid #2a2418', maxHeight: '25vh' }}>
                        <table className="w-full text-xs">
                          <thead>
                            <tr style={{ background: '#1a1610', borderBottom: '1px solid #2a2418' }}>
                              {['Class', 'Location', 'Chassis', 'Name', 'Type'].map(h => (
                                <th key={h} className="text-left px-3 py-2 font-semibold uppercase tracking-wide" style={{ color: 'var(--color-primary)' }}>{h}</th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {vehicles.map((v, i) => (
                              <tr key={v.id} style={{ borderBottom: '1px solid #1a1610', background: i % 2 === 0 ? '#0d0b07' : '#0f0d09' }}>
                                <td className="px-3 py-1.5 font-semibold" style={{ color: 'var(--color-text)' }}>{v.class}</td>
                                <td className="px-3 py-1.5" style={{ color: 'var(--color-text-dim)' }}>{v.map || '—'}</td>
                                <td className="px-3 py-1.5" style={{ color: v.chassis_durability < 0.3 ? '#e88' : 'var(--color-text-dim)' }}>{Math.round(v.chassis_durability * 100)}%</td>
                                <td className="px-3 py-1.5" style={{ color: 'var(--color-text-dim)' }}>{v.vehicle_name || '—'}</td>
                                <td className="px-3 py-1.5">
                                  {v.is_backup && <span className="px-1.5 py-0.5 rounded text-xs" style={{ background: '#1a1a2a', color: '#8888ff', border: '1px solid #3a3a6a' }}>Backup</span>}
                                  {v.is_recovered && <span className="px-1.5 py-0.5 rounded text-xs" style={{ background: '#2a1a0a', color: '#f0a830', border: '1px solid #5a3a10' }}>Recovered</span>}
                                </td>
                              </tr>
                            ))}
                            {!vehiclesLoading && vehicles.length === 0 && (
                              <tr><td colSpan={5} className="px-3 py-6 text-center" style={{ color: 'var(--color-text-dim)' }}>No vehicles found</td></tr>
                            )}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </>
                )}
              </Modal.Body>
              <Modal.Footer>
                <Button onPress={onClose} variant="tertiary">Close</Button>
              </Modal.Footer>
            </Modal.Dialog>
          </Modal.Container>
        </Modal.Backdrop>
      </Modal>
      {confirmationDialog}
    </>
  )
}
