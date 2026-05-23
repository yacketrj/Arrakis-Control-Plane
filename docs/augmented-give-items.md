# Augmented Give Items

Dune Admin supports augmented item grants through `POST /api/v1/players/give-item`.

The endpoint remains backward compatible with legacy single-item payloads and the existing multi-item stack workflow.

## Request shape

```json
{
  "player_id": 123,
  "items": [
    {
      "template": "ItemTemplateWeaponExample",
      "qty": 1,
      "quality": 5,
      "stack_size": 1,
      "augments": [
        {
          "name": "T6_Augment_Damage1",
          "grade": 5,
          "roll": 1.0,
          "roll_count": 1,
          "effect_indices": []
        },
        {
          "name": "T6_Augment_Magazinecapacity1",
          "grade": 5,
          "rolls": [1.0, 1.0, 1.0],
          "effect_indices": []
        }
      ]
    }
  ]
}
```

## Field behavior

### Item fields

- `template`: item template identifier.
- `qty`: number of stacks to create.
- `quality`: item grade, accepted range `0-5`.
- `stack_size`: items per stack.
- `augments`: optional list of augment definitions.

### Augment fields

- `name`: augment template name, for example `T6_Augment_Damage1`.
- `grade`: augment grade, accepted range `1-5`.
- `quality`: backward-compatible alias for `grade`.
- `roll`: single normalized roll value in the range `0.0-1.0`.
- `rolls`: explicit roll array. Use this for augments such as magazine capacity that appear to store more than one roll value.
- `roll_count`: number of times to repeat `roll` when `rolls` is not supplied.
- `effect_indices`: optional selected effect indices. Send an empty array when unused.

## Stored JSON shape

Augments are written into `dune.items.stats` under the observed Unreal-style wrapper shape:

```json
{
  "FAugmentedItemStats": [
    [],
    {
      "AppliedAugments": [
        { "Name": "T6_Augment_Damage1" }
      ],
      "AppliedAugmentRollData": [
        {
          "StatRolls": [1.0],
          "AppliedEffectIndices": []
        }
      ],
      "AppliedAugmentQualities": [5]
    }
  ]
}
```

The arrays are index-aligned:

```text
AppliedAugments[n]
AppliedAugmentRollData[n]
AppliedAugmentQualities[n]
```

## Validation

The backend validates:

- Maximum item rows per request.
- Item template presence.
- Stack count and stack size limits.
- Item grade range.
- Maximum augments per item.
- Augment name presence.
- Augment grade range.
- Roll value range.
- Maximum roll count and explicit roll array length.

## Reliability notes

The current backend implementation writes the augmented item JSON directly when creating new stack rows. Existing stack top-up is intentionally bypassed for augmented batch grants so augmented items do not merge into plain stacks or differently augmented stacks.

## UI status

The active Players tab now opens `web/src/tabs/GiveItemModalAugmented.tsx` from the Give Item button. The prior embedded modal remains exported as `LegacyGiveItemModal` only for short-term rollback/reference while the larger player tab is split into smaller components.
