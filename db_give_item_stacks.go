package main

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func cmdGiveItemStacks(playerID int64, template string, stacks, stackSize, quality int64, augments []giveItemAugmentEntry) Cmd {
	return func() Msg {
		if globalDB == nil {
			return msgMutate{err: fmt.Errorf("not connected")}
		}
		if playerID == 0 {
			return msgMutate{err: fmt.Errorf("player ID required")}
		}
		template = strings.TrimSpace(template)
		if template == "" {
			return msgMutate{err: fmt.Errorf("item template required")}
		}
		if stacks <= 0 {
			return msgMutate{err: fmt.Errorf("stacks must be > 0")}
		}
		if stackSize <= 0 {
			return msgMutate{err: fmt.Errorf("stack size must be > 0")}
		}
		if stacks > math.MaxInt64/stackSize {
			return msgMutate{err: fmt.Errorf("total quantity is too large")}
		}

		stats, err := buildAugmentedItemStatsJSON(augments)
		if err != nil {
			return msgMutate{err: err}
		}

		ctx := context.Background()
		var invID int64
		var maxSlots int
		var maxVolume float64
		err = globalDB.QueryRow(ctx, `
			SELECT id, COALESCE(max_item_count, -1), COALESCE(max_item_volume, -1)
			FROM dune.inventories
			WHERE actor_id = $1::bigint AND inventory_type = 0
			LIMIT 1`, playerID).Scan(&invID, &maxSlots, &maxVolume)
		if err != nil {
			err = globalDB.QueryRow(ctx,
				`SELECT id, COALESCE(max_item_count, -1), COALESCE(max_item_volume, -1)
				 FROM dune.inventories WHERE actor_id = $1::bigint LIMIT 1`, playerID).Scan(&invID, &maxSlots, &maxVolume)
			if err != nil {
				return msgMutate{err: fmt.Errorf("find inventory: %w", err)}
			}
		}

		hasSlotCap := maxSlots > 0
		hasVolumeCap := maxVolume > 0
		usedSlots := 0
		usedVolume := 0.0
		maxPos := int64(-1)

		rows, err := globalDB.Query(ctx, `
			SELECT template_id, stack_size, volume_override, position_index
			FROM dune.items
			WHERE inventory_id = $1::bigint`, invID)
		if err != nil {
			return msgMutate{err: err}
		}
		defer rows.Close()
		for rows.Next() {
			var tmpl string
			var existingStackSize int64
			var vol pgtype.Float8
			var pos int64
			if err := rows.Scan(&tmpl, &existingStackSize, &vol, &pos); err != nil {
				continue
			}
			usedSlots++
			if pos > maxPos {
				maxPos = pos
			}
			if hasVolumeCap {
				itemVol := 0.0
				if vol.Valid && vol.Float64 > 0 {
					itemVol = vol.Float64
				} else if itemData.Items != nil {
					if rule, ok := itemData.Items[strings.ToLower(tmpl)]; ok {
						itemVol = rule.Volume
					} else if itemData.DefaultVolume > 0 {
						itemVol = itemData.DefaultVolume
					}
				} else if itemData.DefaultVolume > 0 {
					itemVol = itemData.DefaultVolume
				}
				usedVolume += itemVol * float64(existingStackSize)
			}
		}
		if rows.Err() != nil {
			return msgMutate{err: rows.Err()}
		}

		if hasSlotCap {
			freeSlots := maxSlots - usedSlots
			if freeSlots < int(stacks) {
				return msgMutate{err: fmt.Errorf("inventory full: need %d free slots, have %d", stacks, freeSlots)}
			}
		}

		totalItems := stacks * stackSize
		if hasVolumeCap {
			perItemVol, err := resolveItemVolume(ctx, template)
			if err != nil {
				return msgMutate{err: err}
			}
			if perItemVol > 0 {
				availableVol := maxVolume - usedVolume
				if availableVol < 0 {
					availableVol = 0
				}
				maxByVolume := int64(math.Floor(availableVol / perItemVol))
				if maxByVolume < totalItems {
					return msgMutate{err: fmt.Errorf(
						"over weight limit: room for %d more %s (%.2f/%.2f volume used)",
						maxByVolume, template, usedVolume, maxVolume)}
				}
			}
		}

		tx, err := globalDB.Begin(ctx)
		if err != nil {
			return msgMutate{err: err}
		}
		defer tx.Rollback(ctx)

		nextPos := maxPos + 1
		for i := int64(0); i < stacks; i++ {
			_, err = tx.Exec(ctx, `
				INSERT INTO dune.items (inventory_id, stack_size, position_index, template_id, quality_level, stats)
				VALUES ($1::bigint, $2::bigint, $3::bigint, $4::text, $5::bigint, $6::jsonb)`,
				invID, stackSize, nextPos, template, quality, stats)
			if err != nil {
				return msgMutate{err: err}
			}
			nextPos++
		}

		if err := tx.Commit(ctx); err != nil {
			return msgMutate{err: err}
		}
		augSummary := ""
		if len(augments) > 0 {
			augSummary = fmt.Sprintf(" with %d augment(s)", len(augments))
		}
		return msgMutate{ok: fmt.Sprintf("Added %d stack(s) x %d of %s to player %d%s", stacks, stackSize, template, playerID, augSummary)}
	}
}
