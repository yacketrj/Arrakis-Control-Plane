export type AugmentPreset = {
  name: string
  label: string
  defaultGrade: number
  defaultRoll: number
  defaultRollCount: number
  defaultRolls?: number[]
}

export const AUGMENT_PRESETS: AugmentPreset[] = [
  { name: 'T6_Augment_Damage1', label: 'Damage I', defaultGrade: 5, defaultRoll: 1, defaultRollCount: 1 },
  { name: 'T6_Augment_ReloadSpeed1', label: 'Reload Speed I', defaultGrade: 5, defaultRoll: 1, defaultRollCount: 1 },
  { name: 'T6_Augment_Shielddamage1', label: 'Shield Damage I', defaultGrade: 5, defaultRoll: 1, defaultRollCount: 1 },
  { name: 'T6_Augment_Headshotdamage1', label: 'Headshot Damage I', defaultGrade: 5, defaultRoll: 1, defaultRollCount: 1 },
  { name: 'T6_Augment_Magazinecapacity1', label: 'Magazine Capacity I', defaultGrade: 5, defaultRoll: 1, defaultRollCount: 3, defaultRolls: [1, 1, 1] },
]

export function findAugmentPreset(name: string): AugmentPreset | undefined {
  const normalized = name.trim().toLowerCase()
  return AUGMENT_PRESETS.find(p => p.name.toLowerCase() === normalized)
}
