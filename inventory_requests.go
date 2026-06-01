package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	inventoryRequestScopePersonal = "personal"
	inventoryRequestScopeGuild    = "guild"

	inventoryRequestStatusOpen      = "open"
	inventoryRequestStatusOrdered   = "ordered"
	inventoryRequestStatusFulfilled = "fulfilled"
	inventoryRequestStatusCancelled = "cancelled"

	inventoryOrderStatusOpen      = "open"
	inventoryOrderStatusFilled    = "filled"
	inventoryOrderStatusCancelled = "cancelled"
)

type inventoryRequest struct {
	ID                 string `json:"id"`
	Scope              string `json:"scope"`
	RequesterDiscordID string `json:"requester_discord_id,omitempty"`
	RequesterName      string `json:"requester_name,omitempty"`
	GuildID            string `json:"guild_id,omitempty"`
	PlayerID           string `json:"player_id,omitempty"`
	ItemName           string `json:"item_name"`
	ItemTemplateID     string `json:"item_template_id,omitempty"`
	Quantity           int    `json:"quantity"`
	Notes              string `json:"notes,omitempty"`
	Status             string `json:"status"`
	OrderID            string `json:"order_id,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type inventoryOrder struct {
	ID                 string   `json:"id"`
	Scope              string   `json:"scope"`
	GuildID            string   `json:"guild_id,omitempty"`
	RequesterDiscordID string   `json:"requester_discord_id,omitempty"`
	AssigneeDiscordID  string   `json:"assignee_discord_id,omitempty"`
	AssigneeName       string   `json:"assignee_name,omitempty"`
	RequestIDs         []string `json:"request_ids"`
	Status             string   `json:"status"`
	Notes              string   `json:"notes,omitempty"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
	CompletedAt        string   `json:"completed_at,omitempty"`
}

type inventoryRequestStore struct {
	Requests []inventoryRequest `json:"requests"`
	Orders   []inventoryOrder   `json:"orders"`
}

type inventoryRequestCreatePayload struct {
	Scope              string `json:"scope"`
	RequesterDiscordID string `json:"requester_discord_id"`
	RequesterName      string `json:"requester_name"`
	GuildID            string `json:"guild_id"`
	PlayerID           string `json:"player_id"`
	ItemName           string `json:"item_name"`
	ItemTemplateID     string `json:"item_template_id"`
	Quantity           int    `json:"quantity"`
	Notes              string `json:"notes"`
}

type inventoryRequestUpdatePayload struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
	Notes   string `json:"notes"`
}

type inventoryOrderCreatePayload struct {
	Scope              string   `json:"scope"`
	GuildID            string   `json:"guild_id"`
	RequesterDiscordID string   `json:"requester_discord_id"`
	AssigneeDiscordID  string   `json:"assignee_discord_id"`
	AssigneeName       string   `json:"assignee_name"`
	RequestIDs         []string `json:"request_ids"`
	Notes              string   `json:"notes"`
}

type inventoryOrderUpdatePayload struct {
	Status             string `json:"status"`
	AssigneeDiscordID  string `json:"assignee_discord_id"`
	AssigneeName       string `json:"assignee_name"`
	Notes              string `json:"notes"`
}

func inventoryRequestStorePath() string {
	if path := strings.TrimSpace(os.Getenv("INVENTORY_REQUEST_STORE")); path != "" {
		return path
	}
	return "inventory-requests.json"
}

func loadInventoryRequestStore() (inventoryRequestStore, error) {
	path := inventoryRequestStorePath()
	store := inventoryRequestStore{}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return store, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return store, nil
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return store, fmt.Errorf("decode inventory request store: %w", err)
	}
	return store, nil
}

func saveInventoryRequestStore(store inventoryRequestStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(inventoryRequestStorePath(), data, 0o600)
}

func newInventoryTrackingID(prefix string) string {
	token, err := randomURLToken(10)
	if err != nil || token == "" {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())
	}
	return prefix + "_" + token
}

func cleanInventoryText(value string, maxLen int, field string, required bool) (string, error) {
	trimmed := strings.TrimSpace(value)
	if required && trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if trimmed != "" && containsUnsafeControl(trimmed) {
		return "", fmt.Errorf("%s contains unsupported control characters", field)
	}
	if len(trimmed) > maxLen {
		return "", fmt.Errorf("%s exceeds %d characters", field, maxLen)
	}
	return trimmed, nil
}

func validInventoryRequestScope(scope string) bool {
	return scope == inventoryRequestScopePersonal || scope == inventoryRequestScopeGuild
}

func validInventoryRequestStatus(status string) bool {
	switch status {
	case inventoryRequestStatusOpen, inventoryRequestStatusOrdered, inventoryRequestStatusFulfilled, inventoryRequestStatusCancelled:
		return true
	default:
		return false
	}
}

func validInventoryOrderStatus(status string) bool {
	switch status {
	case inventoryOrderStatusOpen, inventoryOrderStatusFilled, inventoryOrderStatusCancelled:
		return true
	default:
		return false
	}
}

func normalizeInventoryRequestPayload(payload inventoryRequestCreatePayload, now time.Time) (inventoryRequest, error) {
	scope, err := cleanInventoryText(payload.Scope, 32, "scope", true)
	if err != nil {
		return inventoryRequest{}, err
	}
	if !validInventoryRequestScope(scope) {
		return inventoryRequest{}, fmt.Errorf("scope must be personal or guild")
	}
	if payload.Quantity <= 0 || payload.Quantity > 999999 {
		return inventoryRequest{}, fmt.Errorf("quantity must be between 1 and 999999")
	}
	requesterID, err := cleanInventoryText(payload.RequesterDiscordID, 128, "requester_discord_id", scope == inventoryRequestScopePersonal)
	if err != nil {
		return inventoryRequest{}, err
	}
	guildID, err := cleanInventoryText(payload.GuildID, 128, "guild_id", scope == inventoryRequestScopeGuild)
	if err != nil {
		return inventoryRequest{}, err
	}
	requesterName, err := cleanInventoryText(payload.RequesterName, 128, "requester_name", false)
	if err != nil {
		return inventoryRequest{}, err
	}
	playerID, err := cleanInventoryText(payload.PlayerID, 128, "player_id", false)
	if err != nil {
		return inventoryRequest{}, err
	}
	itemName, err := cleanInventoryText(payload.ItemName, 160, "item_name", true)
	if err != nil {
		return inventoryRequest{}, err
	}
	itemTemplateID, err := cleanInventoryText(payload.ItemTemplateID, 256, "item_template_id", false)
	if err != nil {
		return inventoryRequest{}, err
	}
	notes, err := cleanInventoryText(payload.Notes, 1000, "notes", false)
	if err != nil {
		return inventoryRequest{}, err
	}
	stamp := now.UTC().Format(time.RFC3339Nano)
	return inventoryRequest{
		ID:                 newInventoryTrackingID("req"),
		Scope:              scope,
		RequesterDiscordID: requesterID,
		RequesterName:      requesterName,
		GuildID:            guildID,
		PlayerID:           playerID,
		ItemName:           itemName,
		ItemTemplateID:     itemTemplateID,
		Quantity:           payload.Quantity,
		Notes:              notes,
		Status:             inventoryRequestStatusOpen,
		CreatedAt:          stamp,
		UpdatedAt:          stamp,
	}, nil
}

func normalizeInventoryOrderPayload(payload inventoryOrderCreatePayload, store inventoryRequestStore, now time.Time) (inventoryOrder, error) {
	scope, err := cleanInventoryText(payload.Scope, 32, "scope", true)
	if err != nil {
		return inventoryOrder{}, err
	}
	if !validInventoryRequestScope(scope) {
		return inventoryOrder{}, fmt.Errorf("scope must be personal or guild")
	}
	guildID, err := cleanInventoryText(payload.GuildID, 128, "guild_id", scope == inventoryRequestScopeGuild)
	if err != nil {
		return inventoryOrder{}, err
	}
	requesterID, err := cleanInventoryText(payload.RequesterDiscordID, 128, "requester_discord_id", false)
	if err != nil {
		return inventoryOrder{}, err
	}
	assigneeID, err := cleanInventoryText(payload.AssigneeDiscordID, 128, "assignee_discord_id", false)
	if err != nil {
		return inventoryOrder{}, err
	}
	assigneeName, err := cleanInventoryText(payload.AssigneeName, 128, "assignee_name", false)
	if err != nil {
		return inventoryOrder{}, err
	}
	notes, err := cleanInventoryText(payload.Notes, 1000, "notes", false)
	if err != nil {
		return inventoryOrder{}, err
	}
	requestIDs := normalizeRequestIDs(payload.RequestIDs)
	if len(requestIDs) == 0 {
		return inventoryOrder{}, fmt.Errorf("request_ids must include at least one request id")
	}
	if err := validateOrderRequestIDs(store, requestIDs, scope, guildID); err != nil {
		return inventoryOrder{}, err
	}
	stamp := now.UTC().Format(time.RFC3339Nano)
	return inventoryOrder{
		ID:                 newInventoryTrackingID("ord"),
		Scope:              scope,
		GuildID:            guildID,
		RequesterDiscordID: requesterID,
		AssigneeDiscordID:  assigneeID,
		AssigneeName:       assigneeName,
		RequestIDs:         requestIDs,
		Status:             inventoryOrderStatusOpen,
		Notes:              notes,
		CreatedAt:          stamp,
		UpdatedAt:          stamp,
	}, nil
}

func normalizeRequestIDs(ids []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || containsUnsafeControl(id) || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func validateOrderRequestIDs(store inventoryRequestStore, requestIDs []string, scope string, guildID string) error {
	requestSet := map[string]bool{}
	for _, id := range requestIDs {
		requestSet[id] = true
	}
	found := map[string]bool{}
	for _, request := range store.Requests {
		if !requestSet[request.ID] {
			continue
		}
		if request.Scope != scope {
			return fmt.Errorf("request %s scope does not match order scope", request.ID)
		}
		if scope == inventoryRequestScopeGuild && request.GuildID != guildID {
			return fmt.Errorf("request %s guild_id does not match order guild_id", request.ID)
		}
		if request.Status == inventoryRequestStatusFulfilled || request.Status == inventoryRequestStatusCancelled {
			return fmt.Errorf("request %s is not orderable", request.ID)
		}
		found[request.ID] = true
	}
	for id := range requestSet {
		if !found[id] {
			return fmt.Errorf("request %s was not found", id)
		}
	}
	return nil
}

func handleListInventoryRequests(w http.ResponseWriter, r *http.Request) {
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	guildID := strings.TrimSpace(r.URL.Query().Get("guild_id"))
	requesterID := strings.TrimSpace(r.URL.Query().Get("requester_discord_id"))
	var out []inventoryRequest
	for _, request := range store.Requests {
		if scope != "" && request.Scope != scope {
			continue
		}
		if status != "" && request.Status != status {
			continue
		}
		if guildID != "" && request.GuildID != guildID {
			continue
		}
		if requesterID != "" && request.RequesterDiscordID != requesterID {
			continue
		}
		out = append(out, request)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	jsonOK(w, out)
}

func handleCreateInventoryRequest(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	defer r.Body.Close()
	var payload inventoryRequestCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	request, err := normalizeInventoryRequestPayload(payload, time.Now())
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	store.Requests = append(store.Requests, request)
	if err := saveInventoryRequestStore(store); err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, request)
}

func handleUpdateInventoryRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" || containsUnsafeControl(id) {
		jsonErr(w, fmt.Errorf("valid request id is required"), http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	defer r.Body.Close()
	var payload inventoryRequestUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	for i := range store.Requests {
		if store.Requests[i].ID != id {
			continue
		}
		if status := strings.TrimSpace(payload.Status); status != "" {
			if !validInventoryRequestStatus(status) {
				jsonErr(w, fmt.Errorf("invalid request status"), http.StatusBadRequest)
				return
			}
			store.Requests[i].Status = status
		}
		orderID, err := cleanInventoryText(payload.OrderID, 128, "order_id", false)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if orderID != "" {
			store.Requests[i].OrderID = orderID
		}
		notes, err := cleanInventoryText(payload.Notes, 1000, "notes", false)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if notes != "" {
			store.Requests[i].Notes = notes
		}
		store.Requests[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		if err := saveInventoryRequestStore(store); err != nil {
			jsonErr(w, err, http.StatusInternalServerError)
			return
		}
		jsonOK(w, store.Requests[i])
		return
	}
	jsonErr(w, fmt.Errorf("inventory request not found"), http.StatusNotFound)
}

func handleListInventoryOrders(w http.ResponseWriter, r *http.Request) {
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	guildID := strings.TrimSpace(r.URL.Query().Get("guild_id"))
	var out []inventoryOrder
	for _, order := range store.Orders {
		if scope != "" && order.Scope != scope {
			continue
		}
		if status != "" && order.Status != status {
			continue
		}
		if guildID != "" && order.GuildID != guildID {
			continue
		}
		out = append(out, order)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	jsonOK(w, out)
}

func handleCreateInventoryOrder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	defer r.Body.Close()
	var payload inventoryOrderCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	order, err := normalizeInventoryOrderPayload(payload, store, time.Now())
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	store.Orders = append(store.Orders, order)
	markRequestsForOrder(store.Requests, order.ID, order.RequestIDs, inventoryRequestStatusOrdered)
	if err := saveInventoryRequestStore(store); err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, order)
}

func handleUpdateInventoryOrder(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" || containsUnsafeControl(id) {
		jsonErr(w, fmt.Errorf("valid order id is required"), http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	defer r.Body.Close()
	var payload inventoryOrderUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	store, err := loadInventoryRequestStore()
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	for i := range store.Orders {
		if store.Orders[i].ID != id {
			continue
		}
		status := strings.TrimSpace(payload.Status)
		if status != "" {
			if !validInventoryOrderStatus(status) {
				jsonErr(w, fmt.Errorf("invalid order status"), http.StatusBadRequest)
				return
			}
			store.Orders[i].Status = status
			if status == inventoryOrderStatusFilled {
				stamp := time.Now().UTC().Format(time.RFC3339Nano)
				store.Orders[i].CompletedAt = stamp
				markRequestsForOrder(store.Requests, store.Orders[i].ID, store.Orders[i].RequestIDs, inventoryRequestStatusFulfilled)
			}
			if status == inventoryOrderStatusCancelled {
				markRequestsForOrder(store.Requests, store.Orders[i].ID, store.Orders[i].RequestIDs, inventoryRequestStatusOpen)
			}
		}
		assigneeID, err := cleanInventoryText(payload.AssigneeDiscordID, 128, "assignee_discord_id", false)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if assigneeID != "" {
			store.Orders[i].AssigneeDiscordID = assigneeID
		}
		assigneeName, err := cleanInventoryText(payload.AssigneeName, 128, "assignee_name", false)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if assigneeName != "" {
			store.Orders[i].AssigneeName = assigneeName
		}
		notes, err := cleanInventoryText(payload.Notes, 1000, "notes", false)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if notes != "" {
			store.Orders[i].Notes = notes
		}
		store.Orders[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
		if err := saveInventoryRequestStore(store); err != nil {
			jsonErr(w, err, http.StatusInternalServerError)
			return
		}
		jsonOK(w, store.Orders[i])
		return
	}
	jsonErr(w, fmt.Errorf("inventory order not found"), http.StatusNotFound)
}

func markRequestsForOrder(requests []inventoryRequest, orderID string, requestIDs []string, status string) {
	ids := map[string]bool{}
	for _, id := range requestIDs {
		ids[id] = true
	}
	stamp := time.Now().UTC().Format(time.RFC3339Nano)
	for i := range requests {
		if !ids[requests[i].ID] {
			continue
		}
		requests[i].OrderID = orderID
		requests[i].Status = status
		requests[i].UpdatedAt = stamp
		if status == inventoryRequestStatusOpen {
			requests[i].OrderID = ""
		}
	}
}
