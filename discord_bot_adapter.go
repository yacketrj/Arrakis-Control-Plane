package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	discordBotActionCreateRequest = "create_request"
	discordBotActionCreateOrder   = "create_order"
	discordBotActionUpdateOrder   = "update_order"
)

type discordBotCommandInput struct {
	Name               string
	DiscordUserID      string
	DiscordUserName    string
	GuildID            string
	PlayerID           string
	ItemName           string
	ItemTemplateID     string
	Quantity           int
	Notes              string
	RequestIDs         []string
	AssigneeDiscordID  string
	AssigneeName       string
	OrderStatus        string
}

type discordBotCommandAction struct {
	Action             string
	Request            *inventoryRequest
	Order              *inventoryOrder
	OrderUpdatePayload *inventoryOrderUpdatePayload
}

func normalizeDiscordBotCommandName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.TrimPrefix(name, "/")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func adaptDiscordBotCommand(input discordBotCommandInput, store inventoryRequestStore, now time.Time) (discordBotCommandAction, error) {
	switch normalizeDiscordBotCommandName(input.Name) {
	case "request_item", "personal_request", "request":
		request, err := normalizeInventoryRequestPayload(inventoryRequestCreatePayload{
			Scope:              inventoryRequestScopePersonal,
			RequesterDiscordID: input.DiscordUserID,
			RequesterName:      input.DiscordUserName,
			PlayerID:           input.PlayerID,
			ItemName:           input.ItemName,
			ItemTemplateID:     input.ItemTemplateID,
			Quantity:           input.Quantity,
			Notes:              input.Notes,
		}, now)
		if err != nil {
			return discordBotCommandAction{}, err
		}
		return discordBotCommandAction{Action: discordBotActionCreateRequest, Request: &request}, nil

	case "guild_request", "guild_item_request", "guild_request_item":
		request, err := normalizeInventoryRequestPayload(inventoryRequestCreatePayload{
			Scope:              inventoryRequestScopeGuild,
			RequesterDiscordID: input.DiscordUserID,
			RequesterName:      input.DiscordUserName,
			GuildID:            input.GuildID,
			PlayerID:           input.PlayerID,
			ItemName:           input.ItemName,
			ItemTemplateID:     input.ItemTemplateID,
			Quantity:           input.Quantity,
			Notes:              input.Notes,
		}, now)
		if err != nil {
			return discordBotCommandAction{}, err
		}
		return discordBotCommandAction{Action: discordBotActionCreateRequest, Request: &request}, nil

	case "farm_order", "create_order", "assign_order":
		order, err := normalizeInventoryOrderPayload(inventoryOrderCreatePayload{
			Scope:             coalesceInventoryScope(input.GuildID),
			GuildID:           input.GuildID,
			RequesterDiscordID: input.DiscordUserID,
			AssigneeDiscordID: input.AssigneeDiscordID,
			AssigneeName:      input.AssigneeName,
			RequestIDs:        input.RequestIDs,
			Notes:             input.Notes,
		}, store, now)
		if err != nil {
			return discordBotCommandAction{}, err
		}
		return discordBotCommandAction{Action: discordBotActionCreateOrder, Order: &order}, nil

	case "fill_order", "complete_order", "cancel_order":
		status := strings.TrimSpace(input.OrderStatus)
		if status == "" && normalizeDiscordBotCommandName(input.Name) == "cancel_order" {
			status = inventoryOrderStatusCancelled
		}
		if status == "" {
			status = inventoryOrderStatusFilled
		}
		if !validInventoryOrderStatus(status) {
			return discordBotCommandAction{}, fmt.Errorf("invalid order status")
		}
		payload := inventoryOrderUpdatePayload{
			Status:            status,
			AssigneeDiscordID: input.AssigneeDiscordID,
			AssigneeName:      input.AssigneeName,
			Notes:             input.Notes,
		}
		return discordBotCommandAction{Action: discordBotActionUpdateOrder, OrderUpdatePayload: &payload}, nil

	default:
		return discordBotCommandAction{}, fmt.Errorf("unsupported discord bot command")
	}
}

func coalesceInventoryScope(guildID string) string {
	if strings.TrimSpace(guildID) != "" {
		return inventoryRequestScopeGuild
	}
	return inventoryRequestScopePersonal
}
