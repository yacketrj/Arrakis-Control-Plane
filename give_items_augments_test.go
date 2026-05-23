package main

import (
	"encoding/json"
	"net/http/httptest"
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

func TestNormalizeGiveItemsRequestAugmentQualityAlias(t *testing.T) {
	req := giveItemsRequest{
		PlayerID: 1,
		Items: []giveItemEntry{{
			Template:  "Item",
			Qty:       1,
			StackSize: 1,
			Augments:  []giveItemAugmentEntry{{Name: "Aug", Quality: 4}},
		}},
	}
	items, err := normalizeGiveItemsRequest(req)
	if err != nil {
		t.Fatalf("normalizeGiveItemsRequest returned error: %v", err)
	}
	if items[0].Augments[0].Grade != 4 {
		t.Fatalf("expected quality alias to populate grade, got %+v", items[0].Augments[0])
	}
}

func TestNormalizeGiveItemsRequestRejectsInvalidAugment(t *testing.T) {
	cases := []struct {
		name string
		req  giveItemsRequest
	}{
		{
			name: "blank augment name",
			req:  giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: " "}}}}},
		},
		{
			name: "bad augment grade",
			req:  giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 6}}}}},
		},
		{
			name: "bad roll",
			req:  giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 1, Roll: 1.25}}}}},
		},
		{
			name: "too many rolls",
			req:  giveItemsRequest{PlayerID: 1, Items: []giveItemEntry{{Template: "Item", Qty: 1, StackSize: 1, Augments: []giveItemAugmentEntry{{Name: "Aug", Grade: 1, Rolls: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}}}}}},
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

func TestNormalizeAugmentRollsDefaultsAndExplicitRolls(t *testing.T) {
	defaultRolls := normalizeAugmentRolls(giveItemAugmentEntry{Name: "Aug", Grade: 5})
	if len(defaultRolls) != 1 || defaultRolls[0] != 1.0 {
		t.Fatalf("expected default max roll, got %#v", defaultRolls)
	}

	repeated := normalizeAugmentRolls(giveItemAugmentEntry{Name: "Aug", Grade: 5, Roll: 0.5, RollCount: 3})
	if len(repeated) != 3 || repeated[0] != 0.5 || repeated[1] != 0.5 || repeated[2] != 0.5 {
		t.Fatalf("expected repeated roll values, got %#v", repeated)
	}

	explicit := normalizeAugmentRolls(giveItemAugmentEntry{Name: "Aug", Grade: 5, Roll: 0.25, RollCount: 3, Rolls: []float64{1, 0.75, 0}})
	if len(explicit) != 3 || explicit[0] != 1 || explicit[1] != 0.75 || explicit[2] != 0 {
		t.Fatalf("expected explicit rolls to win, got %#v", explicit)
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

func TestMergeItemTemplatesAndHandleGetTemplates(t *testing.T) {
	oldTemplates := dbItemTemplates
	oldItemData := itemData
	defer func() {
		dbItemTemplates = oldTemplates
		itemData = oldItemData
	}()

	itemData = itemDataFile{
		Names: map[string]string{
			"json_template": "Curated JSON Template",
			"db_template":   "DB Template Friendly Name",
		},
		Items: map[string]itemRule{
			"item_rule_template": {Name: "Item Rule Template", StackMax: 10},
		},
	}

	mergeItemTemplates([]string{"DB_TEMPLATE", "live_template"})
	if len(dbItemTemplates) != 3 {
		t.Fatalf("expected three merged templates, got %#v", dbItemTemplates)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/players/templates", nil)
	handleGetTemplates(rec, req)
	if rec.Code != 200 {
		t.Fatalf("expected 200 response, got %d", rec.Code)
	}
	var rows []templateOut
	if err := json.Unmarshal(rec.Body.Bytes(), &rows); err != nil {
		t.Fatalf("invalid template response JSON: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected three template rows, got %#v", rows)
	}
	foundFriendlyName := false
	for _, row := range rows {
		if row.ID == "DB_TEMPLATE" && row.Name == "DB Template Friendly Name" {
			foundFriendlyName = true
		}
	}
	if !foundFriendlyName {
		t.Fatalf("expected DB_TEMPLATE to use lower-case friendly name lookup, got %#v", rows)
	}
}
