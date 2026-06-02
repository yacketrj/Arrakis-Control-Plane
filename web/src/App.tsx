import { Suspense, lazy, useEffect, useState, type ReactNode } from 'react'
import { useStatus } from './hooks/useStatus'
import { getAdminToken, setAdminToken, type Status } from './api/client'

const AuditTab = lazy(() => import('./tabs/AuditTab'))
const BattlegroupTab = lazy(() => import('./tabs/BattlegroupTab'))
const PlayersTab = lazy(() => import('./tabs/PlayersTabWith360Launcher'))
const PlayerCardsTab = lazy(() => import('./tabs/PlayerCardsTab'))
const Player360Tab = lazy(() => import('./tabs/Player360Tab'))
const InventoryStudioTab = lazy(() => import('./tabs/InventoryStudioTab'))
const FarmingRequestsTab = lazy(() => import('./tabs/FarmingRequestsTab'))
const DiscordPlayerLinksTab = lazy(() => import('./tabs/DiscordPlayerLinksTab'))
const SelfPlayerCardTab = lazy(() => import('./tabs/SelfPlayerCardTab'))
const DatabaseTab = lazy(() => import('./tabs/DatabaseTab'))
const DbRoutinesTab = lazy(() => import('./tabs/DbRoutinesTab'))
const LogsTab = lazy(() => import('./tabs/LogsTab'))
const BlueprintsTab = lazy(() => import('./tabs/BlueprintsTab'))
const StorageTab = lazy(() => import('./tabs/StorageTab'))

type TabId = 'battlegroup' | 'players' | 'player-cards' | 'player-360' | 'inventory-studio' | 'farming-requests' | 'discord-player-links' | 'self-player-card' | 'database' | 'db-routines' | 'audit' | 'logs' | 'blueprints' | 'storage'

const tabs: Array<{ id: TabId; label: string }> = [
  { id: 'audit', label: 'Audit' },
  { id: 'battlegroup', label: 'Battlegroup' },
  { id: 'players', label: 'Players' },
  { id: 'player-cards', label: 'Player Cards' },
  { id: 'player-360', label: 'Player 360' },
  { id: 'inventory-studio', label: 'Inventory Studio' },
  { id: 'farming-requests', label: 'Farming Requests' },
  { id: 'discord-player-links', label: 'Discord Links' },
  { id: 'self-player-card', label: 'My Player Card' },
  { id: 'database', label: 'Database' },
  { id: 'db-routines', label: 'DB Routines' },
  { id: 'logs', label: 'Logs' },
  { id: 'blueprints', label: 'Blueprints' },
  { id: 'storage', label: 'Storage' },
]

const dbBackedTabs = new Set<TabId>(['players', 'player-cards', 'player-360', 'inventory-studio', 'database', 'db-routines', 'blueprints', 'storage'])

export default function App() {
  const status = useStatus()
  const [showBackendConfig, setShowBackendConfig] = useState(false)
  const [backendUrl, setBackendUrl] = useState(() => localStorage.getItem('dune_admin_backend') || '')
  const [tokenInput, setTokenInput] = useState(() => getAdminToken())
  const [activeTab, setActiveTab] = useState<TabId>('audit')

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

  const browserAccessConfigured = Boolean(getAdminToken())
  const dbUnavailableStatus = browserAccessConfigured && status && !status.db_connected ? status : null
  const dbUnavailable = dbUnavailableStatus !== null
  const activeTabBlockedByAccess = !browserAccessConfigured && activeTab !== 'self-player-card'
  const activeTabBlockedByDb = dbUnavailable && dbBackedTabs.has(activeTab)

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
            title="Configure backend URL and browser access"
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

      {!browserAccessConfigured && <AccessBanner onConfigure={() => setShowBackendConfig(true)} />}
      {browserAccessConfigured && dbUnavailableStatus && <DbUnavailableBanner status={dbUnavailableStatus} />}

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
                Browser Access Key
              </label>
              <input
                type="password"
                value={tokenInput}
                onChange={e => setTokenInput(e.target.value)}
                placeholder="Paste value from backend configuration"
                style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', borderRadius: '4px', color: 'var(--color-text)', fontFamily: 'monospace', fontSize: '12px', outline: 'none', padding: '6px 10px' }}
                onKeyDown={e => { if (e.key === 'Enter') saveBackendSettings() }}
              />
              <span style={{ fontSize: '11px', color: 'var(--color-text-dim)' }}>
                Stored locally in this browser and sent to the backend with protected requests.
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
              const selfServiceTab = tab.id === 'self-player-card'
              const blocked = (!browserAccessConfigured && !selfServiceTab) || (dbUnavailable && dbBackedTabs.has(tab.id))
              return (
                <button
                  key={tab.id}
                  role="tab"
                  aria-selected={selected}
                  onClick={() => setActiveTab(tab.id)}
                  title={blocked ? 'Configuration required' : undefined}
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
          {activeTabBlockedByAccess ? <AccessSetupPanel onConfigure={() => setShowBackendConfig(true)} /> : activeTabBlockedByDb && dbUnavailableStatus ? <DbBlockedPanel status={dbUnavailableStatus} /> : <LazyTab>{renderTab(activeTab)}</LazyTab>}
        </div>
      </div>
    </div>
  )
}

function renderTab(tab: TabId) {
  switch (tab) {
    case 'battlegroup': return <BattlegroupTab />
    case 'players': return <PlayersTab />
    case 'player-cards': return <PlayerCardsTab />
    case 'player-360': return <Player360Tab />
    case 'inventory-studio': return <InventoryStudioTab />
    case 'farming-requests': return <FarmingRequestsTab />
    case 'discord-player-links': return <DiscordPlayerLinksTab />
    case 'self-player-card': return <SelfPlayerCardTab />
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
    case 'player-cards':
    case 'player-360':
    case 'farming-requests':
    case 'discord-player-links':
    case 'self-player-card':
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
    <span style={{ color: connected ? 'var(--color-primary)' : 'var(--color-text-dim)' }}>
      {label}: {connected ? 'connected' : 'offline'}
    </span>
  )
}

function AccessBanner({ onConfigure }: { onConfigure: () => void }) {
  return (
    <div className="px-4 py-2 text-xs" style={{ background: '#2a1a0a', color: '#f0a830', borderBottom: '1px solid #5a3a10' }}>
      Browser access is not configured. Protected sections are locked until a backend URL and access key are configured.{' '}
      <button onClick={onConfigure} style={{ color: '#ffd28a', textDecoration: 'underline', background: 'transparent', border: 0, cursor: 'pointer' }}>
        Configure access
      </button>
    </div>
  )
}

function AccessSetupPanel({ onConfigure }: { onConfigure: () => void }) {
  return (
    <div className="rounded-lg p-6 text-sm" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
      <h2 className="text-lg font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Access required</h2>
      <p>Set the backend URL and Browser Access Key before using protected admin sections.</p>
      <button onClick={onConfigure} className="mt-4 rounded px-3 py-2 text-sm" style={{ background: 'var(--color-primary)', color: '#fff' }}>
        Configure backend access
      </button>
    </div>
  )
}

function DbUnavailableBanner({ status }: { status: Status }) {
  return (
    <div className="px-4 py-2 text-xs" style={{ background: '#2a1a0a', color: '#f0a830', borderBottom: '1px solid #5a3a10' }}>
      Database-backed features are unavailable. SSH: {status.ssh_connected ? 'connected' : 'offline'} · DB: {status.db_connected ? 'connected' : 'offline'}
      {status.startup_connect_error ? ` · ${status.startup_connect_error}` : ''}
    </div>
  )
}

function DbBlockedPanel({ status }: { status: Status }) {
  return (
    <div className="rounded-lg p-6 text-sm" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
      <h2 className="text-lg font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Database unavailable</h2>
      <p>This section requires a working database connection.</p>
      <p className="mt-2">SSH: {status.ssh_connected ? 'connected' : 'offline'} · DB: {status.db_connected ? 'connected' : 'offline'}</p>
      {status.startup_connect_error && <p className="mt-2 font-mono text-xs">{status.startup_connect_error}</p>}
    </div>
  )
}
