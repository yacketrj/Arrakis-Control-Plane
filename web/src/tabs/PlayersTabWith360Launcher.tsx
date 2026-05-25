import { useEffect, useRef } from 'react'
import PlayersTab from './PlayersTab'

const player360EventName = 'dune-admin-open-player-360'
const player360StorageKey = 'dune_admin_player_360_id'
const launcherClassName = 'dune-admin-player-360-launcher'

function openPlayer360(playerId: string) {
  localStorage.setItem(player360StorageKey, playerId)
  window.dispatchEvent(new CustomEvent(player360EventName, { detail: { playerId } }))
}

function installLaunchButtons(root: HTMLElement) {
  const rows = Array.from(root.querySelectorAll('tbody tr')) as HTMLTableRowElement[]
  for (const row of rows) {
    if (row.querySelector(`.${launcherClassName}`)) continue
    const cells = row.querySelectorAll('td')
    if (cells.length < 6) continue
    const playerId = cells[0]?.textContent?.trim()
    const actionsCell = cells[5]
    if (!playerId || !/^\d+$/.test(playerId) || !actionsCell) continue

    const button = document.createElement('button')
    button.type = 'button'
    button.textContent = '360'
    button.className = launcherClassName
    button.setAttribute('aria-label', `Open Player 360 profile for player ${playerId}`)
    button.style.border = '1px solid #2a2418'
    button.style.borderRadius = '4px'
    button.style.background = 'transparent'
    button.style.color = 'var(--color-primary)'
    button.style.cursor = 'pointer'
    button.style.fontSize = '12px'
    button.style.padding = '2px 8px'
    button.addEventListener('click', event => {
      event.preventDefault()
      event.stopPropagation()
      openPlayer360(playerId)
    })

    const actionWrapper = actionsCell.querySelector('.flex') ?? actionsCell
    actionWrapper.appendChild(button)
  }
}

export default function PlayersTabWith360Launcher() {
  const rootRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    const root = rootRef.current
    if (!root) return

    installLaunchButtons(root)
    const observer = new MutationObserver(() => installLaunchButtons(root))
    observer.observe(root, { childList: true, subtree: true })
    return () => observer.disconnect()
  }, [])

  return (
    <div ref={rootRef} className="h-full">
      <PlayersTab />
    </div>
  )
}
