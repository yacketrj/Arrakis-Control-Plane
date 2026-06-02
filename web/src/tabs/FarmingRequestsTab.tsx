import { useEffect, useState, type CSSProperties } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import {
  inventoryRequestsApi,
  type InventoryOrder,
  type InventoryOrderStatus,
  type InventoryRequest,
  type InventoryRequestScope,
  type InventoryRequestStatus,
} from '../api/inventoryRequests'

type RequestForm = {
  scope: InventoryRequestScope
  requester_discord_id: string
  requester_name: string
  guild_id: string
  player_id: string
  item_name: string
  item_template_id: string
  quantity: string
  notes: string
}

type OrderForm = {
  scope: InventoryRequestScope
  guild_id: string
  assignee_discord_id: string
  assignee_name: string
  notes: string
}

const requestStatusOptions: Array<InventoryRequestStatus | 'all'> = ['all', 'open', 'ordered', 'fulfilled', 'cancelled']
const orderStatusOptions: Array<InventoryOrderStatus | 'all'> = ['all', 'open', 'filled', 'cancelled']

function fmt(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  return String(value)
}

function fmtDate(value: string): string {
  if (!value) return '—'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

function statusStyle(status: string): CSSProperties {
  const active = status === 'open' || status === 'ordered'
  return {
    border: '1px solid #2a2418',
    color: active ? 'var(--color-primary)' : 'var(--color-text-dim)',
  }
}

export default function FarmingRequestsTab() {
  const [requests, setRequests] = useState<InventoryRequest[]>([])
  const [orders, setOrders] = useState<InventoryOrder[]>([])
  const [loading, setLoading] = useState(false)
  const [requestStatus, setRequestStatus] = useState<InventoryRequestStatus | 'all'>('all')
  const [orderStatus, setOrderStatus] = useState<InventoryOrderStatus | 'all'>('all')
  const [scopeFilter, setScopeFilter] = useState<InventoryRequestScope | 'all'>('all')
  const [selectedRequestIds, setSelectedRequestIds] = useState<string[]>([])
  const [requestForm, setRequestForm] = useState<RequestForm>({
    scope: 'personal',
    requester_discord_id: '',
    requester_name: '',
    guild_id: '',
    player_id: '',
    item_name: '',
    item_template_id: '',
    quantity: '1',
    notes: '',
  })
  const [orderForm, setOrderForm] = useState<OrderForm>({
    scope: 'guild',
    guild_id: '',
    assignee_discord_id: '',
    assignee_name: '',
    notes: '',
  })

  const load = async () => {
    setLoading(true)
    try {
      const [nextRequests, nextOrders] = await Promise.all([
        inventoryRequestsApi.listRequests({ scope: scopeFilter, status: requestStatus }),
        inventoryRequestsApi.listOrders({ scope: scopeFilter, status: orderStatus }),
      ])
      setRequests(nextRequests)
      setOrders(nextOrders)
      setSelectedRequestIds(ids => ids.filter(id => nextRequests.some(request => request.id === id && request.status === 'open')))
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [scopeFilter, requestStatus, orderStatus])

  const createRequest = async () => {
    try {
      const quantity = Number.parseInt(requestForm.quantity, 10)
      await inventoryRequestsApi.createRequest({
        scope: requestForm.scope,
        requester_discord_id: requestForm.requester_discord_id.trim() || undefined,
        requester_name: requestForm.requester_name.trim() || undefined,
        guild_id: requestForm.guild_id.trim() || undefined,
        player_id: requestForm.player_id.trim() || undefined,
        item_name: requestForm.item_name.trim(),
        item_template_id: requestForm.item_template_id.trim() || undefined,
        quantity,
        notes: requestForm.notes.trim() || undefined,
      })
      toast.success('Inventory request created')
      setRequestForm(form => ({ ...form, item_name: '', item_template_id: '', quantity: '1', notes: '' }))
      await load()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  const createOrder = async () => {
    try {
      await inventoryRequestsApi.createOrder({
        scope: orderForm.scope,
        guild_id: orderForm.guild_id.trim() || undefined,
        assignee_discord_id: orderForm.assignee_discord_id.trim() || undefined,
        assignee_name: orderForm.assignee_name.trim() || undefined,
        request_ids: selectedRequestIds,
        notes: orderForm.notes.trim() || undefined,
      })
      toast.success('Farming order created')
      setSelectedRequestIds([])
      setOrderForm(form => ({ ...form, notes: '' }))
      await load()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  const updateOrderStatus = async (order: InventoryOrder, status: InventoryOrderStatus) => {
    try {
      await inventoryRequestsApi.updateOrder(order.id, { status })
      toast.success(`Order marked ${status}`)
      await load()
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    }
  }

  const toggleRequest = (id: string) => {
    setSelectedRequestIds(ids => ids.includes(id) ? ids.filter(existing => existing !== id) : [...ids, id])
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 className="text-lg font-semibold" style={{ color: 'var(--color-primary)' }}>Farming Requests</h2>
            <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
              Coordination-only request and farming-order ledger. This tab does not write player inventory, guild storage, claim rewards, or Player 360 data.
            </p>
          </div>
          <div className="flex flex-wrap gap-2">
            <select value={scopeFilter} onChange={e => setScopeFilter(e.target.value as InventoryRequestScope | 'all')} className="rounded px-3 py-1.5 text-sm border" style={inputStyle()}>
              <option value="all">All scopes</option>
              <option value="personal">Personal</option>
              <option value="guild">Guild</option>
            </select>
            <select value={requestStatus} onChange={e => setRequestStatus(e.target.value as InventoryRequestStatus | 'all')} className="rounded px-3 py-1.5 text-sm border" style={inputStyle()}>
              {requestStatusOptions.map(status => <option key={status} value={status}>Requests: {status}</option>)}
            </select>
            <select value={orderStatus} onChange={e => setOrderStatus(e.target.value as InventoryOrderStatus | 'all')} className="rounded px-3 py-1.5 text-sm border" style={inputStyle()}>
              {orderStatusOptions.map(status => <option key={status} value={status}>Orders: {status}</option>)}
            </select>
            <Button variant="outline" size="sm" onPress={load} isDisabled={loading}>
              {loading ? <Spinner size="sm" color="current" /> : null}
              Refresh
            </Button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
        <section className="rounded-lg p-4 flex flex-col gap-3" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
          <h3 className="font-semibold" style={{ color: 'var(--color-primary)' }}>Create request</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
            <select value={requestForm.scope} onChange={e => setRequestForm({ ...requestForm, scope: e.target.value as InventoryRequestScope })} className="rounded px-3 py-1.5 text-sm border" style={inputStyle()}>
              <option value="personal">Personal</option>
              <option value="guild">Guild</option>
            </select>
            <input value={requestForm.quantity} onChange={e => setRequestForm({ ...requestForm, quantity: e.target.value })} placeholder="Quantity" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={requestForm.requester_discord_id} onChange={e => setRequestForm({ ...requestForm, requester_discord_id: e.target.value })} placeholder="Requester Discord ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={requestForm.requester_name} onChange={e => setRequestForm({ ...requestForm, requester_name: e.target.value })} placeholder="Requester name" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={requestForm.guild_id} onChange={e => setRequestForm({ ...requestForm, guild_id: e.target.value })} placeholder="Guild ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={requestForm.player_id} onChange={e => setRequestForm({ ...requestForm, player_id: e.target.value })} placeholder="Player ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={requestForm.item_name} onChange={e => setRequestForm({ ...requestForm, item_name: e.target.value })} placeholder="Item name" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
            <input value={requestForm.item_template_id} onChange={e => setRequestForm({ ...requestForm, item_template_id: e.target.value })} placeholder="Item template ID optional" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
            <textarea value={requestForm.notes} onChange={e => setRequestForm({ ...requestForm, notes: e.target.value })} placeholder="Notes" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
          </div>
          <Button variant="primary" size="sm" onPress={createRequest}>Create request</Button>
        </section>

        <section className="rounded-lg p-4 flex flex-col gap-3" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
          <h3 className="font-semibold" style={{ color: 'var(--color-primary)' }}>Create farming order</h3>
          <p className="text-xs" style={{ color: 'var(--color-text-dim)' }}>Select open requests below, then create an order for a farmer or crew.</p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
            <select value={orderForm.scope} onChange={e => setOrderForm({ ...orderForm, scope: e.target.value as InventoryRequestScope })} className="rounded px-3 py-1.5 text-sm border" style={inputStyle()}>
              <option value="personal">Personal</option>
              <option value="guild">Guild</option>
            </select>
            <input value={orderForm.guild_id} onChange={e => setOrderForm({ ...orderForm, guild_id: e.target.value })} placeholder="Guild ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={orderForm.assignee_discord_id} onChange={e => setOrderForm({ ...orderForm, assignee_discord_id: e.target.value })} placeholder="Assignee Discord ID" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <input value={orderForm.assignee_name} onChange={e => setOrderForm({ ...orderForm, assignee_name: e.target.value })} placeholder="Assignee name" className="rounded px-3 py-1.5 text-sm border" style={inputStyle()} />
            <textarea value={orderForm.notes} onChange={e => setOrderForm({ ...orderForm, notes: e.target.value })} placeholder="Order notes" className="rounded px-3 py-1.5 text-sm border md:col-span-2" style={inputStyle()} />
          </div>
          <Button variant="primary" size="sm" onPress={createOrder} isDisabled={selectedRequestIds.length === 0}>Create order from {selectedRequestIds.length} request(s)</Button>
        </section>
      </div>

      <section className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Requests</h3>
        {loading && requests.length === 0 ? <Spinner size="lg" /> : null}
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-3">
          {requests.map(request => (
            <article key={request.id} className="rounded p-3 text-sm" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
              <div className="flex justify-between gap-3">
                <div>
                  <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{request.item_name}</div>
                  <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{request.id} · {request.scope} · qty {request.quantity}</div>
                </div>
                <span className="rounded px-2 py-1 text-xs uppercase" style={statusStyle(request.status)}>{request.status}</span>
              </div>
              <div className="grid grid-cols-2 gap-2 mt-3 text-xs" style={{ color: 'var(--color-text-dim)' }}>
                <div>Requester: {fmt(request.requester_name || request.requester_discord_id)}</div>
                <div>Guild: {fmt(request.guild_id)}</div>
                <div>Player: {fmt(request.player_id)}</div>
                <div>Order: {fmt(request.order_id)}</div>
                <div className="col-span-2">Created: {fmtDate(request.created_at)}</div>
              </div>
              {request.notes ? <p className="mt-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{request.notes}</p> : null}
              {request.status === 'open' ? (
                <label className="mt-3 flex items-center gap-2 text-xs" style={{ color: 'var(--color-text)' }}>
                  <input type="checkbox" checked={selectedRequestIds.includes(request.id)} onChange={() => toggleRequest(request.id)} />
                  Add to next farming order
                </label>
              ) : null}
            </article>
          ))}
        </div>
      </section>

      <section className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Orders</h3>
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-3">
          {orders.map(order => (
            <article key={order.id} className="rounded p-3 text-sm" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
              <div className="flex justify-between gap-3">
                <div>
                  <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{order.id}</div>
                  <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{order.scope} · {order.request_ids.length} request(s)</div>
                </div>
                <span className="rounded px-2 py-1 text-xs uppercase" style={statusStyle(order.status)}>{order.status}</span>
              </div>
              <div className="grid grid-cols-2 gap-2 mt-3 text-xs" style={{ color: 'var(--color-text-dim)' }}>
                <div>Assignee: {fmt(order.assignee_name || order.assignee_discord_id)}</div>
                <div>Guild: {fmt(order.guild_id)}</div>
                <div>Created: {fmtDate(order.created_at)}</div>
                <div>Completed: {fmtDate(order.completed_at || '')}</div>
                <div className="col-span-2">Requests: {order.request_ids.join(', ')}</div>
              </div>
              {order.notes ? <p className="mt-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>{order.notes}</p> : null}
              {order.status === 'open' ? (
                <div className="flex gap-2 mt-3">
                  <Button size="sm" variant="secondary" onPress={() => updateOrderStatus(order, 'filled')}>Mark filled</Button>
                  <Button size="sm" variant="outline" onPress={() => updateOrderStatus(order, 'cancelled')}>Cancel</Button>
                </div>
              ) : null}
            </article>
          ))}
        </div>
      </section>
    </div>
  )
}

function inputStyle(): CSSProperties {
  return {
    background: '#0d0b07',
    color: 'var(--color-text)',
    borderColor: '#2a2418',
    outline: 'none',
  }
}
