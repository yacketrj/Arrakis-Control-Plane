import { useEffect, useState } from 'react'
import { Button, Spinner, toast } from '@heroui/react'
import { discordSelfServiceApi, type DiscordPlayerLink, type SelfPlayerCardSummary } from '../api/discordSelfService'

function fmt(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  if (typeof value === 'number') return value.toLocaleString()
  return String(value)
}

function pct(numerator: number, denominator: number): number {
  if (!denominator || denominator <= 0) return 0
  return Math.max(0, Math.min(100, Math.round((numerator / denominator) * 100)))
}

function ProgressBar({ value }: { value: number }) {
  return (
    <div className="h-2 rounded-full overflow-hidden" style={{ background: '#1a1610', border: '1px solid #2a2418' }}>
      <div className="h-full" style={{ width: `${value}%`, background: 'var(--color-primary)' }} />
    </div>
  )
}

export default function SelfPlayerCardTab() {
  const [link, setLink] = useState<DiscordPlayerLink | null>(null)
  const [card, setCard] = useState<SelfPlayerCardSummary | null>(null)
  const [loading, setLoading] = useState(false)

  const load = async () => {
    setLoading(true)
    try {
      const [nextLink, nextCard] = await Promise.all([
        discordSelfServiceApi.selfLink(),
        discordSelfServiceApi.selfPlayerCard(),
      ])
      setLink(nextLink)
      setCard(nextCard)
    } catch (e: unknown) {
      setLink(null)
      setCard(null)
      toast.danger(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [])

  const journeyComplete = card ? pct(card.journey_summary.complete_nodes, card.journey_summary.total_nodes) : 0
  const journeyRevealed = card ? pct(card.journey_summary.revealed_nodes, card.journey_summary.total_nodes) : 0
  const topSpecs = card?.specializations.slice().sort((a, b) => b.xp - a.xp).slice(0, 5) ?? []

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
        <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 className="text-lg font-semibold" style={{ color: 'var(--color-primary)' }}>My Player Card</h2>
            <p className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
              Read-only Discord self-service view for the player linked to your Discord identity. No admin actions or inventory mutations are available here.
            </p>
          </div>
          <Button variant="outline" size="sm" onPress={load} isDisabled={loading}>
            {loading ? <Spinner size="sm" color="current" /> : null}
            Refresh
          </Button>
        </div>
      </div>

      {loading && !card ? <div className="flex justify-center py-12"><Spinner size="lg" /></div> : null}

      {!loading && !card ? (
        <div className="rounded-lg p-6 text-sm" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418', color: 'var(--color-text-dim)' }}>
          No linked self-service player card is available for the current Discord session. Log in with Discord and ask an admin to link your Discord ID to your player actor ID.
        </div>
      ) : null}

      {card ? (
        <>
          <section className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
            <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
              <div>
                <h3 className="text-xl font-semibold" style={{ color: 'var(--color-primary)' }}>{fmt(card.player_name) || `Player ${card.player_id}`}</h3>
                <div className="text-xs mt-1" style={{ color: 'var(--color-text-dim)' }}>
                  Discord {card.discord_id} · Player {card.player_id} · Linked player {link?.player_id ?? card.player_id}
                </div>
              </div>
              <span className="rounded px-2 py-1 text-xs uppercase" style={{ border: '1px solid #2a2418', color: card.online_status === 'online' ? 'var(--color-primary)' : 'var(--color-text-dim)' }}>
                {card.online_status || 'unknown'}
              </span>
            </div>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mt-4 text-sm">
              <Stat label="Class" value={card.class} />
              <Stat label="Map" value={card.map || card.location?.map} />
              <Stat label="Level" value={card.character_xp?.level} />
              <Stat label="Vehicles" value={card.vehicle_count} />
            </div>
          </section>

          <section className="grid grid-cols-1 xl:grid-cols-3 gap-4">
            <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Inventory summary</h3>
              <div className="grid grid-cols-3 gap-2 text-xs">
                <Stat label="Rows" value={card.inventory_summary.total_items} />
                <Stat label="Stack total" value={card.inventory_summary.total_stack_size} />
                <Stat label="Templates" value={card.inventory_summary.unique_templates} />
              </div>
              <div className="mt-3 flex flex-col gap-1 text-xs">
                {card.inventory_summary.preview_items.slice(0, 6).map(item => (
                  <div key={item.id} className="flex justify-between gap-3" style={{ color: 'var(--color-text-dim)', borderBottom: '1px solid #1a1610' }}>
                    <span style={{ color: 'var(--color-text)' }}>{item.name}</span>
                    <span className="font-mono">x{fmt(item.stack_size)}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Journey</h3>
              <div className="text-xs flex flex-col gap-1">
                <div className="flex justify-between"><span style={{ color: 'var(--color-text-dim)' }}>Complete</span><span className="font-mono">{journeyComplete}%</span></div>
                <ProgressBar value={journeyComplete} />
                <div className="flex justify-between mt-2"><span style={{ color: 'var(--color-text-dim)' }}>Revealed</span><span className="font-mono">{journeyRevealed}%</span></div>
                <ProgressBar value={journeyRevealed} />
                <div className="mt-2" style={{ color: 'var(--color-text-dim)' }}>
                  {fmt(card.journey_summary.complete_nodes)} complete · {fmt(card.journey_summary.revealed_nodes)} revealed · {fmt(card.journey_summary.pending_rewards)} pending rewards
                </div>
              </div>
            </div>

            <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Currencies</h3>
              <div className="flex flex-col gap-1 text-xs">
                {card.currencies.length === 0 ? <span style={{ color: 'var(--color-text-dim)' }}>No currency rows.</span> : null}
                {card.currencies.slice(0, 8).map(row => (
                  <div key={`${row.player_id}-${row.currency_id}`} className="flex justify-between gap-3" style={{ color: 'var(--color-text-dim)', borderBottom: '1px solid #1a1610' }}>
                    <span>Currency {row.currency_id}</span>
                    <span className="font-mono" style={{ color: 'var(--color-text)' }}>{fmt(row.balance)}</span>
                  </div>
                ))}
              </div>
            </div>
          </section>

          <section className="grid grid-cols-1 xl:grid-cols-2 gap-4">
            <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Specializations</h3>
              <div className="flex flex-col gap-1 text-xs">
                {topSpecs.length === 0 ? <span style={{ color: 'var(--color-text-dim)' }}>No specialization rows.</span> : null}
                {topSpecs.map(spec => (
                  <div key={`${spec.track_type}-${spec.level}`} className="flex justify-between gap-3" style={{ color: 'var(--color-text-dim)', borderBottom: '1px solid #1a1610' }}>
                    <span style={{ color: 'var(--color-text)' }}>{spec.track_type}</span>
                    <span className="font-mono">Lv {fmt(spec.level)} · {fmt(spec.xp)} XP</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Factions</h3>
              <div className="flex flex-col gap-1 text-xs">
                {card.factions.length === 0 ? <span style={{ color: 'var(--color-text-dim)' }}>No faction rows.</span> : null}
                {card.factions.slice(0, 8).map(faction => (
                  <div key={`${faction.actor_id}-${faction.faction_id}`} className="flex justify-between gap-3" style={{ color: 'var(--color-text-dim)', borderBottom: '1px solid #1a1610' }}>
                    <span style={{ color: 'var(--color-text)' }}>{faction.faction_name}</span>
                    <span className="font-mono">Rep {fmt(faction.reputation)} · Scrip {fmt(faction.scrips)}</span>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {card.section_errors.length > 0 ? (
            <section className="rounded-lg p-4" style={{ background: 'var(--color-surface)', border: '1px solid #2a2418' }}>
              <h3 className="font-semibold mb-3" style={{ color: 'var(--color-primary)' }}>Partial data warnings</h3>
              <div className="flex flex-col gap-1 text-xs">
                {card.section_errors.map(error => (
                  <div key={`${error.section}-${error.error}`} style={{ color: 'var(--color-text-dim)' }}>{error.section}: {error.error}</div>
                ))}
              </div>
            </section>
          ) : null}
        </>
      ) : null}
    </div>
  )
}

function Stat({ label, value }: { label: string; value: unknown }) {
  return (
    <div className="rounded p-2" style={{ background: '#0d0b07', border: '1px solid #1a1610' }}>
      <div className="text-xs" style={{ color: 'var(--color-text-dim)' }}>{label}</div>
      <div className="font-semibold" style={{ color: 'var(--color-text)' }}>{fmt(value)}</div>
    </div>
  )
}
