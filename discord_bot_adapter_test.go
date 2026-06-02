package main

import (
	"testing"
	"time"
)

func TestAdaptDiscordBotPersonalRequest(t *testing.T) {
	action, err := adaptDiscordBotCommand(discordBotCommandInput{
		Name:            "/request-item",
		DiscordUserID:   "discord-1",
		DiscordUserName: "Ron",
		PlayerID:        "player-1",
		ItemName:        "Copper Ore",
		Quantity:        50,
		Notes:           "Need starter tools.",
	}, inventoryRequestStore{}, time.Date(2026, 6, 1, 1, 2, 3, 0, time.UTC))
	if err != nil {
		t.Fatalf("adapt personal request: %v", err)
	}
	if action.Action != discordBotActionCreateRequest || action.Request == nil {
		t.Fatalf("expected create request action: %#v", action)
	}
	if action.Request.Scope != inventoryRequestScopePersonal || action.Request.RequesterDiscordID != "discord-1" || action.Request.ItemName != "Copper Ore" || action.Request.Quantity != 50 {
		t.Fatalf("unexpected request: %#v", action.Request)
	}
}

func TestAdaptDiscordBotGuildRequest(t *testing.T) {
	action, err := adaptDiscordBotCommand(discordBotCommandInput{
		Name:            "guild request item",
		DiscordUserID:   "discord-1",
		DiscordUserName: "Ron",
		GuildID:         "guild-1",
		ItemName:        "Iron Ore",
		Quantity:        250,
	}, inventoryRequestStore{}, time.Now())
	if err != nil {
		t.Fatalf("adapt guild request: %v", err)
	}
	if action.Action != discordBotActionCreateRequest || action.Request == nil {
		t.Fatalf("expected create request action: %#v", action)
	}
	if action.Request.Scope != inventoryRequestScopeGuild || action.Request.GuildID != "guild-1" {
		t.Fatalf("unexpected guild request: %#v", action.Request)
	}
}

func TestAdaptDiscordBotFarmOrder(t *testing.T) {
	store := inventoryRequestStore{Requests: []inventoryRequest{
		{ID: "req-1", Scope: inventoryRequestScopeGuild, GuildID: "guild-1", Status: inventoryRequestStatusOpen},
	}}
	action, err := adaptDiscordBotCommand(discordBotCommandInput{
		Name:              "farm_order",
		DiscordUserID:     "discord-1",
		GuildID:           "guild-1",
		AssigneeDiscordID: "farmer-1",
		AssigneeName:      "Harvester Crew",
		RequestIDs:        []string{"req-1"},
		Notes:             "Run after storm.",
	}, store, time.Now())
	if err != nil {
		t.Fatalf("adapt farm order: %v", err)
	}
	if action.Action != discordBotActionCreateOrder || action.Order == nil {
		t.Fatalf("expected create order action: %#v", action)
	}
	if action.Order.Scope != inventoryRequestScopeGuild || action.Order.GuildID != "guild-1" || len(action.Order.RequestIDs) != 1 || action.Order.RequestIDs[0] != "req-1" {
		t.Fatalf("unexpected order: %#v", action.Order)
	}
}

func TestAdaptDiscordBotOrderUpdates(t *testing.T) {
	fill, err := adaptDiscordBotCommand(discordBotCommandInput{Name: "fill_order", AssigneeDiscordID: "farmer-1"}, inventoryRequestStore{}, time.Now())
	if err != nil {
		t.Fatalf("adapt fill order: %v", err)
	}
	if fill.Action != discordBotActionUpdateOrder || fill.OrderUpdatePayload == nil || fill.OrderUpdatePayload.Status != inventoryOrderStatusFilled {
		t.Fatalf("unexpected fill action: %#v", fill)
	}

	cancel, err := adaptDiscordBotCommand(discordBotCommandInput{Name: "cancel_order"}, inventoryRequestStore{}, time.Now())
	if err != nil {
		t.Fatalf("adapt cancel order: %v", err)
	}
	if cancel.Action != discordBotActionUpdateOrder || cancel.OrderUpdatePayload == nil || cancel.OrderUpdatePayload.Status != inventoryOrderStatusCancelled {
		t.Fatalf("unexpected cancel action: %#v", cancel)
	}
}

func TestAdaptDiscordBotRejectsUnsupportedCommand(t *testing.T) {
	if _, err := adaptDiscordBotCommand(discordBotCommandInput{Name: "dance"}, inventoryRequestStore{}, time.Now()); err == nil {
		t.Fatal("expected unsupported command to fail")
	}
}
