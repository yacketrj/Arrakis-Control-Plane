import type { GiveItemAugment, GiveItemRow } from '../api/client'
import { findAugmentPreset } from './augmentPresets'

export type GiveItemDraftForPayload = {
  template: string
  qty: number
  quality: number
  stack_size: number
  augments: GiveItemAugment[]
}

export function clampInt(value: string, min: number, max: number, fallback: number): number {
  const parsed = parseInt(value, 10)
  if (!Number.isFinite(parsed)) return fallback
  return Math.max(min, Math.min(max, parsed))
}

export function clampFloat(value: string, min: number, max: number, fallback: number): number {
  const parsed = parseFloat(value)
  if (!Number.isFinite(parsed)) return fallback
  return Math.max(min, Math.min(max, parsed))
}

export function parseRollsCsv(value: string): number[] | undefined {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  const rolls = trimmed
    .split(',')
    .map(part => part.trim())
    .filter(Boolean)
    .map(part => clampFloat(part, 0, 1, 1))
  return rolls.length > 0 ? rolls : undefined
}

export function presetAugment(name: string): GiveItemAugment {
  const preset = findAugmentPreset(name)
  return {
    name,
    grade: preset?.defaultGrade ?? 5,
    roll: preset?.defaultRoll ?? 1,
    rolls: preset?.defaultRolls,
    roll_count: preset?.defaultRollCount ?? 1,
    effect_indices: [],
  }
}

export function toGiveItemPayload(row: GiveItemDraftForPayload): GiveItemRow {
  return {
    template: row.template.trim(),
    qty: row.qty,
    quality: row.quality,
    stack_size: row.stack_size,
    augments: row.augments
      .filter(aug => aug.name.trim())
      .map(aug => ({
        name: aug.name.trim(),
        grade: aug.grade,
        roll: aug.roll,
        rolls: aug.rolls,
        roll_count: aug.roll_count,
        effect_indices: aug.effect_indices ?? [],
      })),
  }
}
