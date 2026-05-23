package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	maxGiveItemAugments      = 5
	maxGiveItemAugmentRolls  = 8
	defaultGiveItemRollValue = 1.0
)

type giveItemAugmentEntry struct {
	Name          string    `json:"name"`
	Grade         int64     `json:"grade"`
	Quality       int64     `json:"quality"`
	Roll          float64   `json:"roll"`
	Rolls         []float64 `json:"rolls"`
	RollCount     int       `json:"roll_count"`
	EffectIndices []int64   `json:"effect_indices"`
}

type augmentedItemStats struct {
	FAugmentedItemStats []any `json:"FAugmentedItemStats"`
}

type appliedAugmentName struct {
	Name string `json:"Name"`
}

type appliedAugmentRollData struct {
	StatRolls            []float64 `json:"StatRolls"`
	AppliedEffectIndices []int64   `json:"AppliedEffectIndices"`
}

type augmentedItemStatsPayload struct {
	AppliedAugments         []appliedAugmentName     `json:"AppliedAugments"`
	AppliedAugmentRollData  []appliedAugmentRollData `json:"AppliedAugmentRollData"`
	AppliedAugmentQualities []int64                  `json:"AppliedAugmentQualities"`
}

func normalizeGiveItemAugments(rowIndex int, augments []giveItemAugmentEntry) error {
	if len(augments) > maxGiveItemAugments {
		return fmt.Errorf("item %d augments must be <= %d", rowIndex, maxGiveItemAugments)
	}
	for j := range augments {
		aug := &augments[j]
		aug.Name = strings.TrimSpace(aug.Name)
		if aug.Name == "" {
			return fmt.Errorf("item %d augment %d name required", rowIndex, j+1)
		}
		if aug.Grade == 0 && aug.Quality > 0 {
			aug.Grade = aug.Quality
		}
		if aug.Grade == 0 {
			aug.Grade = 5
		}
		if aug.Grade < 1 || aug.Grade > 5 {
			return fmt.Errorf("item %d augment %d grade must be 1-5", rowIndex, j+1)
		}
		if len(aug.Rolls) > maxGiveItemAugmentRolls {
			return fmt.Errorf("item %d augment %d rolls must be <= %d", rowIndex, j+1, maxGiveItemAugmentRolls)
		}
		if aug.RollCount < 0 || aug.RollCount > maxGiveItemAugmentRolls {
			return fmt.Errorf("item %d augment %d roll_count must be 0-%d", rowIndex, j+1, maxGiveItemAugmentRolls)
		}
		if aug.Roll < 0 || aug.Roll > 1 {
			return fmt.Errorf("item %d augment %d roll must be 0.0-1.0", rowIndex, j+1)
		}
		for k, roll := range aug.Rolls {
			if roll < 0 || roll > 1 {
				return fmt.Errorf("item %d augment %d roll %d must be 0.0-1.0", rowIndex, j+1, k+1)
			}
		}
	}
	return nil
}

func buildAugmentedItemStatsJSON(augments []giveItemAugmentEntry) (string, error) {
	if len(augments) == 0 {
		return `{}`, nil
	}

	payload := augmentedItemStatsPayload{
		AppliedAugments:         make([]appliedAugmentName, 0, len(augments)),
		AppliedAugmentRollData:  make([]appliedAugmentRollData, 0, len(augments)),
		AppliedAugmentQualities: make([]int64, 0, len(augments)),
	}

	for _, aug := range augments {
		rolls := normalizeAugmentRolls(aug)
		payload.AppliedAugments = append(payload.AppliedAugments, appliedAugmentName{Name: aug.Name})
		payload.AppliedAugmentRollData = append(payload.AppliedAugmentRollData, appliedAugmentRollData{
			StatRolls:            rolls,
			AppliedEffectIndices: append([]int64(nil), aug.EffectIndices...),
		})
		payload.AppliedAugmentQualities = append(payload.AppliedAugmentQualities, aug.Grade)
	}

	stats := augmentedItemStats{FAugmentedItemStats: []any{[]any{}, payload}}
	data, err := json.Marshal(stats)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func normalizeAugmentRolls(aug giveItemAugmentEntry) []float64 {
	rolls := append([]float64(nil), aug.Rolls...)
	if len(rolls) > 0 {
		return rolls
	}

	roll := aug.Roll
	if roll == 0 {
		roll = defaultGiveItemRollValue
	}
	count := aug.RollCount
	if count <= 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		rolls = append(rolls, roll)
	}
	return rolls
}
