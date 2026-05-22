package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestNormalizeGiveItemsRequestSingleItemCompatibility(t *testing.T) {
	items, err := normalizeGiveItemsRequest(giveItemsRequest{
		PlayerID: 42,
		Template: "  tpl-sword  ",
		Qty:      3,
		Quality:  4,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Template != "tpl-sword" || items[0].Qty != 3 || items[0].Quality != 4 || items[0].StackSize != 1 {
		t.Fatalf("unexpected normalized item: %+v", items[0])
	}
}

func TestNormalizeGiveItemsRequestBatchPayload(t *testing.T) {
	items, err := normalizeGiveItemsRequest(giveItemsRequest{
		PlayerID: 42,
		Items: []giveItemEntry{
			{Template: " tpl-a ", Qty: 2, Quality: 0, StackSize: 10},
			{Template: "tpl-b", Qty: 1, Quality: 5, StackSize: 50},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Template != "tpl-a" || items[0].StackSize != 10 {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
	if items[1].Template != "tpl-b" || items[1].Quality != 5 || items[1].StackSize != 50 {
		t.Fatalf("unexpected second item: %+v", items[1])
	}
}

func TestNormalizeGiveItemsRequestAllowsMaximumBounds(t *testing.T) {
	items, err := normalizeGiveItemsRequest(giveItemsRequest{
		PlayerID: 42,
		Items: []giveItemEntry{{Template: "tpl", Qty: maxGiveItemQty, Quality: 5, StackSize: maxGiveItemStackSize}},
	})
	if err != nil {
		t.Fatalf("expected no error at configured max bounds, got %v", err)
	}
	if items[0].Qty != maxGiveItemQty || items[0].StackSize != maxGiveItemStackSize {
		t.Fatalf("unexpected item at bounds: %+v", items[0])
	}
}

func TestNormalizeGiveItemsRequestDefaultsStackSize(t *testing.T) {
	items, err := normalizeGiveItemsRequest(giveItemsRequest{
		PlayerID: 42,
		Items: []giveItemEntry{{Template: "tpl", Qty: 1, Quality: 1, StackSize: 0}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if items[0].StackSize != 1 {
		t.Fatalf("expected default stack size 1, got %d", items[0].StackSize)
	}
}

func TestNormalizeGiveItemsRequestValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		req  giveItemsRequest
		want string
	}{
		{name: "empty payload", req: giveItemsRequest{PlayerID: 42}, want: "at least one item"},
		{name: "blank template", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: " ", Qty: 1, Quality: 1, StackSize: 1}}}, want: "template required"},
		{name: "zero qty", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: 0, Quality: 1, StackSize: 1}}}, want: "quantity must be > 0"},
		{name: "negative qty", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: -1, Quality: 1, StackSize: 1}}}, want: "quantity must be > 0"},
		{name: "qty too high", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: maxGiveItemQty + 1, Quality: 1, StackSize: 1}}}, want: "quantity must be <="},
		{name: "stack too high", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: 1, Quality: 1, StackSize: maxGiveItemStackSize + 1}}}, want: "stack size must be <="},
		{name: "negative quality", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: 1, Quality: -1, StackSize: 1}}}, want: "quality must be 0-5"},
		{name: "high quality", req: giveItemsRequest{PlayerID: 42, Items: []giveItemEntry{{Template: "tpl", Qty: 1, Quality: 6, StackSize: 1}}}, want: "quality must be 0-5"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeGiveItemsRequest(tc.req)
			if err == nil {
				t.Fatalf("expected error containing %q", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q to contain %q", err.Error(), tc.want)
			}
		})
	}
}

func TestNormalizeGiveItemsRequestRejectsMoreThanOneHundredRows(t *testing.T) {
	rows := make([]giveItemEntry, maxGiveItemRows+1)
	for i := range rows {
		rows[i] = giveItemEntry{Template: fmt.Sprintf("tpl-%d", i), Qty: 1, Quality: 1, StackSize: 1}
	}
	_, err := normalizeGiveItemsRequest(giveItemsRequest{PlayerID: 42, Items: rows})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "maximum 100 item rows") {
		t.Fatalf("unexpected error: %v", err)
	}
}
