import { Suspense, lazy, useEffect, useState, type ReactNode } from 'react'
import { useStatus } from './hooks/useStatus'
import { api, getAdminToken, setAdminToken, type Status } from './api/client'

const AuditTab = lazy(() => import('./tabs/AuditTab'))
const BattlegroupTab = lazy(() => import('./tabs/BattlegroupTab'))
const PlayersTab = lazy(() => import('./tabs/PlayersTabWith360Launcher'))
const Player360Tab = lazy(() => import('./tabs/Player360Tab'))
const InventoryStudioTab = lazy(() => import('./tabs/InventoryStudioTab'))
const DatabaseTab = lazy(() => import('./tabs/DatabaseTab'))
const DbRoutinesTab = lazy(() => import('./tabs/DbRoutinesTab'))
const LogsTab = lazy(() => import('./tabs/LogsTab'))
const BlueprintsTab = lazy(() => import('./tabs/BlueprintsTab'))
const StorageTab = lazy(() => import('./tabs/StorageTab'))

type TabId = 'battlegroup' | 'players' | 'player-360' | 'inventory-studio' | 'database' | 'db-routines' | 'audit' | 'logs' | 'blueprints' | 'storage'

const tabs: Array<{ id: TabId; label: string }> = [
  { id: 'battlegroup', label: 'Battlegroup' },
  { id: 'players', label: 'Players' },
  { id: 'player-360', label: 'Player 360' },
  { id: 'inventory-studio', label: 'Inventory Studio' },
  { id: 'database', label: 'Database' },
  { id: 'db-routines', label: 'DB Routines' },
  { id: 'audit', label: 'Audit' },
  { id: 'logs', label: 'Logs' },
  { id: 'blueprints', label: 'Blueprints' },
  { id: 'storage', label: 'Storage' },
]

const dbBackedTabs = new Set<TabId>(['players', 'player-360', 'inventory-studio', 'database', 'db-routines', 'blueprints', 'storage'])

export default function App() {
  const status = useStatus()
  const [showBackendConfig, setShowBackendConfig] = useState(false)
  const [backendUrl, setBackendUrl] = useState(() => localStorage.getItem('dune_admin_backend') || '')
  const [tokenInput, setTokenInput] = useState(() => getAdminToken())
  const [activeTab, setActiveTab] = useState<TabId>('battlegroup')

  useEffect(() => {
    const openPlayer360 = () => setActiveTab('player-360')
    window.addEventListener('dune-admin-open-player-360', openPlayer360)
    return () => window.removeEventListener('dune-admin-open-player-360', openPlayer360)
  }, [])

  const saveBackendSettings = () => {
    const trimmedUrl = backendUrl.trim()
    if (trimmedUrl) localStorage.setItem('dune_admin_backend', trimmedUrl)
    else localStorage.removeItem('dune_admin_backend')
    setAdminToken(tokenInput)
    window.location.reload()
  }

  const resetBackendSettings = () => {
    localStorage.removeItem('dune_admin_backend')
    setAdminToken('')
    window.location.reload()
  }

  const dbUnavailableStatus = getAdminToken() && status && !status.db_connected ? status : null
  const dbUnavailable = dbUnavailableStatus !== null
  const activeTabBlocked = dbUnavailable && dbBackedTabs.has(activeTab)

  return (
    <div className="h-screen flex flex-col overflow-hidden" style={{ background: 'var(--color-background)' }}>
      <div
        className="flex items-center justify-between px-6 py-3 border-b shrink-0"
        style={{ borderColor: '#2a2418', background: 'var(--color-surface)' }}
      >
        <div className="flex items-center gap-3">
          <span
            className="text-xl font-bold tracking-widest uppercase"
            style={{ color: 'var(--color-primary)', letterSpacing: '0.2em' }}
          >
            DUNE ADMIN
          </span>
          {status?.ssh_host && (
            <span className="text-xs" style={{ color: 'var(--color-text-dim)' }}>
              {status.ssh_host}
            </span>
          )}
        </div>
        <div className="flex items-center gap-4 text-xs">
          <ConnectionBadge label="SSH" connected={status?.ssh_connected ?? false} />
          <ConnectionBadge label="DB" connected={status?.db_connected ?? false} />
          {status?.pod_ns && (
            <span style={{ color: 'var(--color-text-dim)' }}>ns: {status.pod_ns}</span>
          )}
          <button
            onClick={() => setShowBackendConfig(v => !v)}
            title="Configure backend URL and admin token"
            style={{
              background: 'transparent',
              border: '1px solid #2a2418',
              borderRadius: '4px',
              color: showBackendConfig ? 'var(--color-primary)' : 'var(--color-text-dim)',
              cursor: 'pointer',
              fontSize: '14px',
              padding: '2px 6px',
              lineHeight: 1,
            }}
          >
            ⚙
          </button>
        </div>
      </div>

      {dbUnavailableStatus && <DbUnavailableBanner status={dbUnavailableStatus} />}

      {showBackendConfig && (
        <div
          style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, zIndex: 100, background: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'flex-start', justifyContent: 'flex-end', paddingTop: '52px', paddingRight: '16px' }}
          onClick={e => { if (e.target === e.currentTarget) setShowBackendConfig(false) }}
        >
          <div style={{ background: '#0d0b07', border: '1px solid #2a2418', borderRadius: '8px', padding: '16px', width: '380px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ color: 'var(--color-primary)', fontWeight: 600, fontSize: '13px' }}>
                Backend Settings
              </span>
              <button onClick={() => setShowBackendConfig(false)} style={{ background: 'transparent', border: 'none', color: 'var(--color-text-dim)', cursor: 'pointer', fontSize: '16px', lineHeight: 1, padding: '0 2px' }}>
                ✕
              </button>
            </div>

            <div style={{ fontSize: '11px', color: 'var(--color-text-dim)' }}>
              Current URL:{' '}
              <span style={{ color: 'var(--color-text)', fontFamily: 'monospace' }}>
                {localStorage.getItem('dune_admin_backend') || 'http://localhost:8080'}
              </span>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <label style={{ fontSize: '11px', color: 'var(--color-text-dim)' }}>
                Backend URL
              </label>
              <input
                value={backendUrl}
                onChange={e => setBackendUrl(e.target.value)}
                placeholder="http://host:port"
                style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', borderRadius: '4px', color: 'var(--color-text)', fontFamily: 'monospace', fontSize: '12px', outline: 'none', padding: '6px 10px' }}
                onKeyDown={e => { if (e.key === 'Enter') saveBackendSettings() }}
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
              <label style={{ fontSize: '11px', color: 'var(--color-text-dim)' }}>
                Admin Token
              </label>
              <input
                type="password"
                value={tokenInput}
                onChange={e => setTokenInput(e.target.value)}
                placeholder="ADMIN_TOKEN from backend .env"
                style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', borderRadius: '4px', color: 'var(--color-text)', fontFamily: 'monospace', fontSize: '12px', outline: 'none', padding: '6px 10px' }}
                onKeyDown={e => { if (e.key === 'Enter') saveBackendSettings() }}
              />
              <span style={{ fontSize: '11px', color: 'var(--color-text-dim)' }}>
                Stored locally in this browser and sent to the backend as X-Admin-Token.
              </span>
            </div>

            <div style={{ display: 'flex', gap: '8px' }}>
              <button onClick={saveBackendSettings} style={{ background: 'var(--color-primary)', border: 'none', borderRadius: '4px', color: '#fff', cursor: 'pointer', fontSize: '12px', fontWeight: 600, padding: '6px 12px' }}>
                Save &amp; Reload
              </button>
              <button onClick={resetBackendSettings} style={{ background: 'transparent', border: '1px solid #2a2418', borderRadius: '4px', color: 'var(--color-text-dim)', cursor: 'pointer', fontSize: '12px', padding: '6px 12px' }}>
                Reset
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="flex-1 flex flex-col overflow-hidden min-h-0">
        <div className="px-4 pt-2 shrink-0" style={{ background: 'var(--color-background)' }}>
          <div role="tablist" aria-label="Admin sections" className="flex gap-1 overflow-x-auto">
            {tabs.map(tab => {
              const selected = activeTab === tab.id
              const blocked = dbUnavailable && dbBackedTabs.has(tab.id)
              return (
                <button
                  key={tab.id}
                  role="tab"
                  aria-selected={selected}
                  onClick={() => setActiveTab(tab.id)}
                  title={blocked ? 'DB connection required' : undefined}
                  style={{
                    background: selected ? 'var(--color-surface)' : 'transparent',
                    border: '1px solid #2a2418',
                    borderBottomColor: selected ? 'var(--color-primary)' : '#2a2418',
                    borderRadius: '6px 6px 0 0',
                    color: blocked ? '#6b5a48' : selected ? 'var(--color-primary)' : 'var(--color-text-dim)',
                    cursor: 'pointer',
                    fontSize: '12px',
                    fontWeight: selected ? 600 : 400,
                    padding: '8px 12px',
                    whiteSpace: 'nowrap',
                  }}
                >
                  {tab.label}{blocked ? ' ⚠' : ''}
                </button>
              )
            })}
          </div>
        </div>
        <div className={panelClass(activeTab)}>
          {activeTabBlocked && dbUnavailableStatus ? <DbBlockedPanel status={dbUnavailableStatus} /> : <LazyTab>{renderTab(activeTab)}</LazyTab>}
        </div>
      </div>
    </div>
  )
}

function renderTab(tab: TabId) {
  switch (tab) {
    case 'battlegroup': return <BattlegroupTab />
    case 'players': return <PlayersTab />
    case 'player-360': return <Player360Tab />
    case 'inventory-studio': return <InventoryStudioTab />
    case 'database': return <DatabaseTab />
    case 'db-routines': return <DbRoutinesTab />
    case 'audit': return <AuditTab />
    case 'logs': return <LogsTab />
    case 'blueprints': return <BlueprintsTab />
    case 'storage': return <StorageTab />
  }
}

function panelClass(tab: TabId) {
  switch (tab) {
    case 'battlegroup':
    case 'db-routines':
    case 'audit':
    case 'logs':
    case 'blueprints':
    case 'storage':
    case 'inventory-studio':
      return 'flex-1 overflow-hidden flex flex-col p-4'
    case 'players':
    case 'player-360':
    case 'database':
      return 'flex-1 overflow-auto p-4'
  }
}

function LazyTab({ children }: { children: ReactNode }) {
  return (
    <Suspense fallback={<TabFallback />}>
      {children}
    </Suspense>
  )
}

function TabFallback() {
  return (
    <div className="flex-1 flex items-center justify-center p-6" style={{ color: 'var(--color-text-dim)' }}>
      Loading section...
    </div>
  )
}

function ConnectionBadge({ label, connected }: { label: string; connected: boolean }) {
  return (
    <div className="flex items-center gap-1.5">
      <div className="w-2 h-2 rounded-full" style={{ background: connected ? 'var(--color-success)' : '#555' }} />
      <span style={{ color: connected ? 'var(--color-text)' : 'var(--color-text-dim)' }}>{label}</span>
    </div>
  )
}

function DbUnavailableBanner({ status }: { status: Status }) {
  return (
    <div className="px-6 py-2 text-xs" style={{ background: '#2a1a0a', borderBottom: '1px solid #5a3a10', color: '#f0a830' }}>
      <strong>DB unavailable.</strong> DB-backed tabs are gated until SSH/tunnel/database connectivity is restored.
      {status.startup_connect_error ? <span className="ml-2 font-mono">{status.startup_connect_error}</span> : null}
    </div>
  )
}

function DbBlockedPanel({ status }: { status: Status }) {
  const [busy, setBusy] = useState(false)
  const [message, setMessage] = useState('')
  const reconnect = async () => {
    setBusy(true)
    setMessage('')
    try {
      const next = await api.reconnect()
      setMessage(next.db_connected ? 'Reconnect succeeded. Reloading...' : 'Reconnect returned but DB is still unavailable.')
      if (next.db_connected) window.location.reload()
    } catch (e: unknown) {
      setMessage(e instanceof Error ? e.message : String(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="rounded-lg p-6 text-sm" style={{ background: '#0d0b07', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
      <h2 className="text-lg font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Database connection required</h2>
      <p>This section depends on live database access. The backend is running, but SSH/tunnel/database connectivity is not healthy.</p>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-2 mt-4 text-xs">
        <StatusTile label="SSH" value={status.ssh_connected ? 'connected' : 'not connected'} />
        <StatusTile label="DB" value={status.db_connected ? 'connected' : 'not connected'} />
        <StatusTile label="Tunnel Mode" value={status.tunnel_mode || 'unknown'} />
      </div>
      {status.startup_connect_error && <pre className="mt-4 rounded p-3 overflow-auto" style={{ background: '#080604', border: '1px solid #2a2418', color: '#f0a830' }}>{status.startup_connect_error}</pre>}
      <div className="mt-4 flex gap-2 items-center">
        <button onClick={reconnect} disabled={busy} style={{ background: 'var(--color-primary)', border: 'none', borderRadius: 4, color: '#fff', cursor: busy ? 'wait' : 'pointer', fontSize: 12, fontWeight: 600, padding: '6px 12px' }}>{busy ? 'Reconnecting...' : 'Retry reconnect'}</button>
        {message && <span style={{ color: 'var(--color-text-dim)' }}>{message}</span>}
      </div>
    </div>
  )
}

function StatusTile({ label, value }: { label: string; value: string }) {
  return <div className="rounded p-2" style={{ background: '#080604', border: '1px solid #2a2418' }}><div className="uppercase tracking-wide text-[10px]" style={{ color: 'var(--color-text-dim)' }}>{label}</div><div className="font-mono mt-1" style={{ color: 'var(--color-text)' }}>{value}</div></div>
}
