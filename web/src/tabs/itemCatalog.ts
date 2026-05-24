export type TemplateOption = { id: string; name: string }

export type ItemTemplateCategory =
  | 'weapon'
  | 'armor-gear'
  | 'augment'
  | 'augmentation-station'
  | 'vehicle'
  | 'resource'
  | 'building-placeable'
  | 'other'

export type AugmentTarget = 'Weapon' | 'Armor / Gear' | 'Shield' | 'Vehicle' | 'Unknown'

const weaponWords = [
  'weapon', 'rifle', 'pistol', 'shotgun', 'sniper', 'smg', 'launcher', 'sword', 'knife', 'blade', 'maula', 'kindjal', 'lasgun', 'dart', 'grenade', 'damage', 'reload', 'headshot', 'magazine', 'firearm',
]
const armorWords = [
  'armor', 'armour', 'gear', 'helmet', 'head', 'chest', 'torso', 'glove', 'gauntlet', 'boot', 'shoe', 'leg', 'pants', 'stillsuit', 'cloak', 'robe', 'shieldbelt', 'shield belt',
]
const shieldWords = ['shield', 'shielddamage', 'shield damage']
const vehicleWords = ['vehicle', 'buggy', 'ornithopter', 'sandbike', 'crawler', 'carrier', 'chassis']
const resourceWords = ['resource', 'ore', 'ingot', 'fiber', 'spice', 'water', 'fuel', 'flour', 'plast', 'copper', 'iron', 'aluminum', 'aluminium', 'steel']
const stationWords = ['station', 'bench', 'workbench', 'fabricator', 'processor', 'refinery', 'assembler', 'crafting']
const buildingWords = ['placeable', 'building', 'wall', 'floor', 'door', 'gate', 'foundation', 'ramp', 'decor', 'deco', 'furniture', 'lamp']

function searchable(template: TemplateOption): string {
  return `${template.id} ${template.name}`.toLowerCase()
}

function hasAny(value: string, words: string[]): boolean {
  return words.some(word => value.includes(word))
}

export function isAugmentationStationTemplate(template: TemplateOption): boolean {
  const value = searchable(template)
  return value.includes('augment') && hasAny(value, stationWords)
}

export function isAugmentItemTemplate(template: TemplateOption): boolean {
  const value = searchable(template)
  return value.includes('augment') && !isAugmentationStationTemplate(template) && !hasAny(value, buildingWords)
}

export function guessAugmentTarget(template: TemplateOption): AugmentTarget {
  const value = searchable(template)
  if (hasAny(value, shieldWords)) return 'Shield'
  if (hasAny(value, armorWords)) return 'Armor / Gear'
  if (hasAny(value, vehicleWords)) return 'Vehicle'
  if (hasAny(value, weaponWords)) return 'Weapon'
  return 'Unknown'
}

export function categorizeTemplate(template: TemplateOption): ItemTemplateCategory {
  const value = searchable(template)
  if (isAugmentationStationTemplate(template)) return 'augmentation-station'
  if (isAugmentItemTemplate(template)) return 'augment'
  if (hasAny(value, weaponWords)) return 'weapon'
  if (hasAny(value, armorWords) || hasAny(value, shieldWords)) return 'armor-gear'
  if (hasAny(value, vehicleWords)) return 'vehicle'
  if (hasAny(value, resourceWords)) return 'resource'
  if (hasAny(value, buildingWords) || hasAny(value, stationWords)) return 'building-placeable'
  return 'other'
}

export function categoryLabel(category: ItemTemplateCategory): string {
  switch (category) {
    case 'weapon': return 'Weapons'
    case 'armor-gear': return 'Armor / Gear'
    case 'augment': return 'Item Augments'
    case 'augmentation-station': return 'Augmentation Stations'
    case 'vehicle': return 'Vehicles'
    case 'resource': return 'Resources'
    case 'building-placeable': return 'Buildings / Placeables'
    case 'other': return 'Other'
  }
}

export function categoryHelp(category: ItemTemplateCategory): string {
  switch (category) {
    case 'augment': return 'Attach these to a selected weapon, armor, shield, or gear item row.'
    case 'augmentation-station': return 'Crafting/placeable stations. These are not item augments and should not appear in the augment picker.'
    case 'weapon': return 'Weapon item templates. Select one as the item row, then attach compatible augments.'
    case 'armor-gear': return 'Armor, shields, stillsuits, and wearable gear. Select one as the item row, then attach compatible armor/gear augments.'
    default: return 'General item templates discovered from the database.'
  }
}
