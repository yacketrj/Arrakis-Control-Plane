import { useEffect, useRef, useState } from 'react'
import { api } from '../api/client'
import type { Player } from '../api/client'
import InventoryModal from './InventoryModal'
import PlayerActionsModalConfirmed from './PlayerActionsModalConfirmed'
import PlayerTeleportModal from './PlayerTeleportModal'
import PlayersTab from './PlayersTab'

const player360EventName = 'dune-admin-open-player-360'
const player360StorageKey = 'dune_admin_player_360_id'
const launcherClassName = 'dune-admin-player-360-launcher'
const teleportLauncherClassName = 'dune-admin-player-move-launcher'

type InterceptedAction = 'Inventory' | 'Actions' | 'Move'

function openPlayer360(playerId: string) {
  localStorage.setItem(player360StorageKey, playerId)
  window.dispatchEvent(new CustomEvent(player360EventName, { detail: { playerId } }))
}

function makeLauncherButton(text: string, className: string, ariaLabel: string) {
  const button = document.createElement('button')
  button.type = 'button'
  button.textContent = text
  button.className = className
  button.setAttribute('aria-label', ariaLabel)
  button.style.border = '1px solid #2a2418'
  button.style.borderRadius = '4px'
  button.style.background = 'transparent'
  button.style.color = 'var(--color-primary)'
  button.style.cursor = 'pointer'
  button.style.fontSize = '12px'
  button.style.padding = '2px 8px'
  return button
}

function installLaunchButtons(root: HTMLElement) {
  const rows = Array.from(root.querySelectorAll('tbody tr')) as HTMLTableRowElement[]
  for (const row of rows) {
    const cells = row.querySelectorAll('td')
    if (cells.length < 6) continue
    const playerId = cells[0]?.textContent?.trim()
    const actionsCell = cells[5]
    if (!playerId || !/^\d+$/.test(playerId) || !actionsCell) continue

    const actionWrapper = actionsCell.querySelector('.flex') ?? actionsCell

    if (!row.querySelector(`.${launcherClassName}`)) {
      const button = makeLauncherButton('360', launcherClassName, `Open Player 360 profile for player ${playerId}`)
      button.addEventListener('click', event => {
        event.preventDefault()
        event.stopPropagation()
        openPlayer360(playerId)
      })
      actionWrapper.appendChild(button)
    }

    if (!row.querySelector(`.${teleportLauncherClassName}`)) {
      const button = makeLauncherButton('Move', teleportLauncherClassName, `Move player ${playerId}`)
      actionWrapper.appendChild(button)
    }
  }
}

function playerActionFromClick(target: EventTarget | null): { playerId: number; action: InterceptedAction } | null {
  if (!(target instanceof HTMLElement)) return null
  const button = target.closest('button')
  if (!button) return null
  const action = button.textContent?.trim()
  if (action !== 'Inventory' && action !== 'Actions' && action !== 'Move') return null

  const row = button.closest('tr')
  const playerId = row?.querySelector('td')?.textContent?.trim()
  if (!playerId || !/^\d+$/.test(playerId)) return null
  return { playerId: Number(playerId), action }
}

export default function PlayersTabWith360Launcher() {
  const rootRef = useRef<HTMLDivElement | null>(null)
  const [players, setPlayers] = useState<Player[]>([])
  const [inventoryPlayer, setInventoryPlayer] = useState<Player | null>(null)
  const [showInventory, setShowInventory] = useState(false)
  const [actionsPlayer, setActionsPlayer] = useState<Player | null>(null)
  const [showActions, setShowActions] = useState(false)
  const [teleportPlayer, setTeleportPlayer] = useState<Player | null>(null)
  const [showTeleport, setShowTeleport] = useState(false)

  useEffect(() => {
    let active = true
    api.players.list()
      .then(rows => { if (active) setPlayers(rows) })
      .catch(() => {})
    return () => { active = false }
  }, [])

  useEffect(() => {
    const root = rootRef.current
    if (!root) return

    const handlePlayerActionClick = async (event: MouseEvent) => {
      const clicked = playerActionFromClick(event.target)
      if (!clicked) return

      event.preventDefault()
      event.stopPropagation()
      event.stopImmediatePropagation()

      let player = players.find(row => row.id === clicked.playerId) ?? null
      if (!player) {
        try {
          const refreshed = await api.players.list()
          setPlayers(refreshed)
          player = refreshed.find(row => row.id === clicked.playerId) ?? null
        } catch {
          player = null
        }
      }

      if (!player) return
      if (clicked.action === 'Inventory') {
        setInventoryPlayer(player)
        setShowInventory(true)
      } else if (clicked.action === 'Actions') {
        setActionsPlayer(player)
        setShowActions(true)
      } else {
        setTeleportPlayer(player)
        setShowTeleport(true)
      }
    }

    installLaunchButtons(root)
    root.addEventListener('click', handlePlayerActionClick, true)
    const observer = new MutationObserver(() => installLaunchButtons(root))
    observer.observe(root, { childList: true, subtree: true })
    return () => {
      observer.disconnect()
      root.removeEventListener('click', handlePlayerActionClick, true)
    }
  }, [players])

  return (
    <div className="h-full">
      <div ref={rootRef} className="h-full">
        <PlayersTab />
      </div>
      {inventoryPlayer && (
        <InventoryModal player={inventoryPlayer} open={showInventory} onClose={() => setShowInventory(false)} />
      )}
      {actionsPlayer && (
        <PlayerActionsModalConfirmed player={actionsPlayer} open={showActions} onClose={() => setShowActions(false)} />
      )}
      {teleportPlayer && (
        <PlayerTeleportModal player={teleportPlayer} open={showTeleport} onClose={() => setShowTeleport(false)} />
      )}
    </div>
  )
}
