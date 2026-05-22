package main

import (
	"encoding/json"
	"testing"
)

func TestNormalizeGiveItemsRequestWithAugments(t *testing.T) {
	req := giveItemsRequest{
		PlayerID: 42,
		Items: []giveItemEntry{{
			Template:  "  ItemTemplateHeavyDart  ",
			Qty:       2,
			Quality:   5,
			StackSize: 1000,
			Augments: []giveItemAugmentEntry{{
				Name:      "  T6_Augment_Damage1  ",
				Grade:     5,
				RollCount: 2,
			}},
		}},
	}

	items, err := normalizeGiveItemsRequest(req)
	if err != nil {
		t.Fatalf("normalizeGiveItemsRequest returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item, got %d", len(items))
	}
	item := items[0]
	if item.Template != "ItemTemplateHeavyDart" {
		t.Fatalf("template was not trimmed: %q", item.Template)
	}
	if item.Qty != 2 || item.StackSize != 1000 || item.Quality != 5 {
		t.Fatalf("unexpected normalized item: %+v", item)
	}
	if len(item.Augments) != 1 {
		t.Fatalf("expected one augment, got %d", len(item.Augments))
	}
	if item.Augments[0].Name != "T6_Augment_Damage1" || item.Augments[0].Grade != 5 {
		t.Fatalf("unexpected augment: %+v", item.Augments[0])
	}
}

func TestNormalizeGiveItemsRequestLegacyAugments(t *testing.T) {
	req := giveItemsRequest{
		PlayerID:  42,
		Template:  "ItemTemplateWeapon",
		Qty:       1,
		Quality:   4,
		StackSize: 1,
		Augments: []giveItemAugmentEntry{{
			Name:  "T6_Augment_ReloadSpeed1",
			Roll:  0.75,
			Grade: 3,
		}},
	}
	items, err := normalizeGiveItemsRequest(req)
	if err != nil {
		t.Fatalf("normalizeGiveItemsRequest returned error: %v", err)
	}
	if len(items) != 1 || len(items[0].Augments) != 1 {
		t.Fatalf("expected legacy single item with one augment, got %+v", items)
	}
	if items[0].StackSize != 1 {
		t.Fatalf("expected stack size to be preserved/defaulted to 1, got %d", items[0].StackSize)
	}
}

func TestNormalizeGiveItemsRequestRejectsInvalidAugment(t *testing.T) {
	cases := []struct {
		name string
		req  giveItemsRequest
	}{
		{
			name: "blank augment name",
			req: giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: " "}}}}},
		},
		{
			name: "bad augment grade",
			req: giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 6}}}}},
		},
		{
			name: "bad roll",
			req: giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 1, Roll: 1.25}}}}},
		},
		{
			name: "too many rolls",
			req: giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 1, Rolls: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}}}}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := normalizeGiveItemsRequest(tc.req); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestBuildAugmentedItemStatsJSON(t *testing.T) {
	statsText, err := buildAugmentedItemStatsJSON([]giveItemAugmentEntry{
		{Name: "T6_Augment_Damage1", Grade: 5, Roll: 0.5, RollCount: 2},
		{Name: "T6_Augment_Magazinecapacity1", Grade: 4, Rolls: []float64{1, 0.75, 0}, EffectIndices: []int64{2}},
	})
	if err != nil {
		t.Fatalf("buildAugmentedItemStatsJSON returned error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(statsText), &decoded); err != nil {
		t.Fatalf("stats JSON is invalid: %v", err)
	}
	wrapper, ok := decoded["FAugmentedItemStats"].([]any)
	if !ok || len(wrapper) != 2 {
		t.Fatalf("expected FAugmentedItemStats wrapper with two entries, got %#v", decoded["FAugmentedItemStats"])
	}
	payload, ok := wrapper[1].(map[string]any)
	if !ok {
		t.Fatalf("expected payload object at wrapper index 1, got %#v", wrapper[1])
	}
	augments := payload["AppliedAugments"].([]any)
	rolls := payload["AppliedAugmentRollData"].([]any)
	qualities := payload["AppliedAugmentQualities"].([]any)
	if len(augments) != 2 || len(rolls) != 2 || len(qualities) != 2 {
		t.Fatalf("augment arrays are not aligned: %d %d %d", len(augments), len(rolls), len(qualities))
	}
	if qualities[0].(float64) != 5 || qualities[1].(float64) != 4 {
		t.Fatalf("unexpected qualities: %#v", qualities)
	}
}

func TestBuildAugmentedItemStatsJSONEmpty(t *testing.T) {
	statsText, err := buildAugmentedItemStatsJSON(nil)
	if err != nil {
		t.Fatalf("buildAugmentedItemStatsJSON returned error: %v", err)
	}
	if statsText != `{}` {
		t.Fatalf("expected empty stats object, got %s", statsText)
	}
}
