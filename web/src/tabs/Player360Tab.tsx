import { useMemo, useState, type ReactNode } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { getPlayerProfile, type PlayerProfile } from '../api/playerProfile'

function fmt(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  if (typeof value === 'number') return value.toLocaleString()
  return String(value)
}

function Section({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="rounded-lg p-3" style={{ background: '#0d0b07', border: '1px solid #2a2418' }}>
      <h3 className="text-sm font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>{title}</h3>
      {children}
    </section>
  )
}

function KV({ label, value }: { label: string; value: unknown }) {
  return (
    <div className="flex justify-between gap-3 text-xs py-1" style={{ borderBottom: '1px solid #1a1610' }}>
      <span style={{ color: 'var(--color-text-dim)' }}>{label}</span>
      <span className="text-right font-mono" style={{ color: 'var(--color-text)' }}>{fmt(value)}</span>
    </div>
  )
}

function PartialData({ profile }: { profile: PlayerProfile }) {
  if (!profile.section_errors.length) return null
  return (
    <div className="rounded-lg p-3 text-xs" style={{ background: '#2a1a0a', border: '1px solid #5a3a10', color: '#f0a830' }}>
      <div className="font-semibold mb-1">Partial data</div>
      {profile.section_errors.map((e, idx) => (
        <div key={`${e.section}-${idx}`}>{e.section}: {e.error}</div>
      ))}
    </div>
  )
}

export default function Player360Tab() {
  const [playerIdInput, setPlayerIdInput] = useState('')
  const [profile, setProfile] = useState<PlayerProfile | null>(null)
  const [loading, setLoading] = useState(false)

  const loadProfile = async () => {
    const id = Number.parseInt(playerIdInput, 10)
    if (!Number.isFinite(id) || id <= 0) {
      toast.danger('Enter a valid player actor ID')
      return
    }
    setLoading(true)
    try {
      setProfile(await getPlayerProfile(id))
    } catch (e: unknown) {
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  const previewItems = useMemo(() => profile?.inventory_summary.preview_items ?? [], [profile])
  const previewJourney = useMemo(() => profile?.journey_summary.preview_nodes ?? [], [profile])

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <div className="flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 className="text-lg font-semibold" style={{ color: 'var(--color-primary)' }}>Player 360 Profile</h2>
            <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
              Protected read-only support view. Enter a PlayerCharacter actor ID to aggregate identity, online state, currency, inventory summary, vehicles, faction status, specs, journey, events, and dungeon history.
            </p>
          </div>
          <div className="flex gap-2">
            <input
              className="rounded px-3 py-1.5 text-sm border w-48"
              style={{ background: '#0d0b07', color: 'var(--color-text)', borderColor: '#2a2418', outline: 'none' }}
              placeholder="Player actor ID"
              value={playerIdInput}
              onChange={e => setPlayerIdInput(e.target.value)}
              onKeyDown={e => { if (e.key === 'Enter') loadProfile() }}
            />
            <Button variant="outline" size="sm" onPress={loadProfile} isDisabled={loading}>
              {loading ? <Spinner size="sm" color="current" /> : null}
              Load
            </Button>
          </div>
        </div>
      </div>

      {loading && <div className="flex justify-center py-12"><Spinner size="lg" /></div>}

      {!loading && !profile && (
        <div className="rounded-lg p-6 text-sm" style={{ background: '#0d0b07', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
          No profile loaded yet.
        </div>
      )}

      {!loading && profile && (
        <div className="flex flex-col gap-3">
          <PartialData profile={profile} />

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-3">
            <Section title="Player Info">
              <KV label="Name" value={profile.identity?.name} />
              <KV label="Actor ID" value={profile.identity?.id ?? profile.player_id} />
              <KV label="Account ID" value={profile.identity?.account_id} />
              <KV label="Controller ID" value={profile.identity?.controller_id} />
              <KV label="FLS ID" value={profile.identity?.fls_id} />
              <KV label="Class" value={profile.identity?.class} />
            </Section>

            <Section title="Online and Location">
              <KV label="Status" value={profile.online_state?.status ?? profile.identity?.online_status} />
              <KV label="Map" value={profile.location?.map ?? profile.online_state?.map ?? profile.identity?.map} />
              <KV label="Location Source" value={profile.location?.source} />
              <KV label="Last Seen" value={profile.online_state?.last_seen} />
            </Section>

            <Section title="Character Progression">
              <KV label="Character Level" value={profile.character_xp?.level} />
              <KV label="Character XP" value={profile.character_xp?.xp} />
              <KV label="Spec Tracks" value={profile.specializations.length} />
              <KV label="Journey Nodes" value={profile.journey_summary.total_nodes} />
            </Section>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-3">
            <Section title="Currencies">
              {profile.currencies.length === 0 ? <p className="text-xs" style={{ color: 'var(--color-text-dim)' }}>No currency rows.</p> : profile.currencies.map(c => (
                <KV key={`${c.player_id}-${c.currency_id}`} label={`Currency ${c.currency_id}`} value={c.balance} />
              ))}
            </Section>

            <Section title="Faction Status">
              {profile.factions.length === 0 ? <p className="text-xs" style={{ color: 'var(--color-text-dim)' }}>No faction rows.</p> : profile.factions.map(f => (
                <div key={`${f.actor_id}-${f.faction_id}`} className="text-xs py-1" style={{ borderBottom: '1px solid #1a1610' }}>
                  <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{f.faction_name}</div>
                  <div style={{ color: 'var(--color-text-dim)' }}>Rep {fmt(f.reputation)} · Scrips {fmt(f.scrips)}</div>
                </div>
              ))}
            </Section>

            <Section title="Inventory Summary">
              <KV label="Item Rows" value={profile.inventory_summary.total_items} />
              <KV label="Total Stack Size" value={profile.inventory_summary.total_stack_size} />
              <KV label="Unique Templates" value={profile.inventory_summary.unique_templates} />
            </Section>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <Section title="Inventory Preview">
              {previewItems.length === 0 ? <p className="text-xs" style={{ color: 'var(--color-text-dim)' }}>No preview items.</p> : (
                <div className="overflow-auto" style={{ maxHeight: 220 }}>
                  <table className="w-full text-xs">
                    <tbody>
                      {previewItems.map(item => (
                        <tr key={item.id} style={{ borderBottom: '1px solid #1a1610' }}>
                          <td className="py-1 pr-2" style={{ color: 'var(--color-text)' }}>{item.name || item.template_id}</td>
                          <td className="py-1 text-right font-mono" style={{ color: 'var(--color-text-dim)' }}>x{item.stack_size}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </Section>

            <Section title="Journey Summary">
              <KV label="Complete" value={profile.journey_summary.complete_nodes} />
              <KV label="Revealed" value={profile.journey_summary.revealed_nodes} />
              <KV label="Pending Rewards" value={profile.journey_summary.pending_rewards} />
              {previewJourney.length > 0 && (
                <div className="mt-2 text-xs" style={{ color: 'var(--color-text-dim)' }}>
                  {previewJourney.slice(0, 5).map(node => <div key={node.node_id}>{node.node_id}</div>)}
                </div>
              )}
            </Section>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-3">
            <Section title="Vehicles">
              <KV label="Vehicle Rows" value={profile.vehicles.length} />
              {profile.vehicles.slice(0, 5).map(v => (
                <div key={v.id} className="text-xs py-1" style={{ borderBottom: '1px solid #1a1610', color: 'var(--color-text-dim)' }}>{v.class} · {v.map || 'no map'}</div>
              ))}
            </Section>

            <Section title="Recent Events">
              <KV label="Event Rows" value={profile.recent_events.length} />
              {profile.recent_events.slice(0, 5).map((event, idx) => (
                <div key={`${event.universe_time}-${idx}`} className="text-xs py-1" style={{ borderBottom: '1px solid #1a1610', color: 'var(--color-text-dim)' }}>{event.universe_time} · type {event.event_type}</div>
              ))}
            </Section>

            <Section title="Dungeon History">
              <KV label="Dungeon Rows" value={profile.dungeon_history.length} />
              {profile.dungeon_history.slice(0, 5).map(d => (
                <div key={d.completion_id} className="text-xs py-1" style={{ borderBottom: '1px solid #1a1610', color: 'var(--color-text-dim)' }}>{d.dungeon_id} · {d.difficulty}</div>
              ))}
            </Section>
          </div>
        </div>
      )}
    </div>
  )
}
