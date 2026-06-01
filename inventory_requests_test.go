package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func useTempInventoryRequestStore(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "inventory-requests.json")
	t.Setenv("INVENTORY_REQUEST_STORE", path)
	return path
}

func TestNormalizeInventoryRequestPayload(t *testing.T) {
	request, err := normalizeInventoryRequestPayload(inventoryRequestCreatePayload{
		Scope:              inventoryRequestScopePersonal,
		RequesterDiscordID: "discord-123",
		RequesterName:      "Sihaya",
		PlayerID:           "player-1",
		ItemName:           "Solari",
		Quantity:           25,
		Notes:              "Need for field kit.",
	}, time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("normalize request: %v", err)
	}
	if request.ID == "" || request.Status != inventoryRequestStatusOpen {
		t.Fatalf("unexpected request identity/status: %#v", request)
	}
	if request.Scope != inventoryRequestScopePersonal || request.RequesterDiscordID != "discord-123" || request.Quantity != 25 {
		t.Fatalf("unexpected request payload: %#v", request)
	}
	if request.CreatedAt == "" || request.UpdatedAt == "" {
		t.Fatalf("expected timestamps: %#v", request)
	}
}

func TestNormalizeInventoryRequestPayloadRejectsInvalidValues(t *testing.T) {
	tests := []inventoryRequestCreatePayload{
		{Scope: "", RequesterDiscordID: "discord-123", ItemName: "Water", Quantity: 1},
		{Scope: "bad", RequesterDiscordID: "discord-123", ItemName: "Water", Quantity: 1},
		{Scope: inventoryRequestScopePersonal, ItemName: "Water", Quantity: 1},
		{Scope: inventoryRequestScopeGuild, ItemName: "Water", Quantity: 1},
		{Scope: inventoryRequestScopePersonal, RequesterDiscordID: "discord-123", ItemName: "", Quantity: 1},
		{Scope: inventoryRequestScopePersonal, RequesterDiscordID: "discord-123", ItemName: "Water", Quantity: 0},
	}
	for _, payload := range tests {
		if _, err := normalizeInventoryRequestPayload(payload, time.Now()); err == nil {
			t.Fatalf("expected invalid payload to fail: %#v", payload)
		}
	}
}

func TestInventoryRequestHandlersCreateListAndOrderLifecycle(t *testing.T) {
	useTempInventoryRequestStore(t)
	mux := http.NewServeMux()
	registerRoutes(mux)

	createRequestBody := inventoryRequestCreatePayload{
		Scope:              inventoryRequestScopeGuild,
		GuildID:            "guild-1",
		RequesterDiscordID: "discord-123",
		RequesterName:      "Ron",
		ItemName:           "Copper Ore",
		Quantity:           100,
		Notes:              "Guild smelting run.",
	}
	requestRecorder := httptest.NewRecorder()
	mux.ServeHTTP(requestRecorder, jsonRequest(http.MethodPost, "/api/v1/inventory/requests", createRequestBody))
	if requestRecorder.Code != http.StatusOK {
		t.Fatalf("create request got %d: %s", requestRecorder.Code, requestRecorder.Body.String())
	}
	createdRequest := decodeJSONResponse[inventoryRequest](t, requestRecorder)
	if createdRequest.ID == "" || createdRequest.Status != inventoryRequestStatusOpen {
		t.Fatalf("unexpected created request: %#v", createdRequest)
	}

	listRecorder := httptest.NewRecorder()
	mux.ServeHTTP(listRecorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/requests?scope=guild&guild_id=guild-1", nil))
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list request got %d: %s", listRecorder.Code, listRecorder.Body.String())
	}
	listedRequests := decodeJSONResponse[[]inventoryRequest](t, listRecorder)
	if len(listedRequests) != 1 || listedRequests[0].ID != createdRequest.ID {
		t.Fatalf("unexpected listed requests: %#v", listedRequests)
	}

	createOrderBody := inventoryOrderCreatePayload{
		Scope:             inventoryRequestScopeGuild,
		GuildID:           "guild-1",
		AssigneeDiscordID: "farmer-1",
		AssigneeName:      "Harvester Team",
		RequestIDs:        []string{createdRequest.ID},
		Notes:             "Assigned to west route.",
	}
	orderRecorder := httptest.NewRecorder()
	mux.ServeHTTP(orderRecorder, jsonRequest(http.MethodPost, "/api/v1/inventory/orders", createOrderBody))
	if orderRecorder.Code != http.StatusOK {
		t.Fatalf("create order got %d: %s", orderRecorder.Code, orderRecorder.Body.String())
	}
	createdOrder := decodeJSONResponse[inventoryOrder](t, orderRecorder)
	if createdOrder.ID == "" || createdOrder.Status != inventoryOrderStatusOpen || len(createdOrder.RequestIDs) != 1 {
		t.Fatalf("unexpected created order: %#v", createdOrder)
	}

	storeAfterOrder, err := loadInventoryRequestStore()
	if err != nil {
		t.Fatalf("load store after order: %v", err)
	}
	if len(storeAfterOrder.Requests) != 1 || storeAfterOrder.Requests[0].Status != inventoryRequestStatusOrdered || storeAfterOrder.Requests[0].OrderID != createdOrder.ID {
		t.Fatalf("request was not marked ordered: %#v", storeAfterOrder.Requests)
	}

	patchRecorder := httptest.NewRecorder()
	mux.ServeHTTP(patchRecorder, jsonRequest(http.MethodPatch, "/api/v1/inventory/orders/"+createdOrder.ID, inventoryOrderUpdatePayload{Status: inventoryOrderStatusFilled}))
	if patchRecorder.Code != http.StatusOK {
		t.Fatalf("patch order got %d: %s", patchRecorder.Code, patchRecorder.Body.String())
	}
	updatedOrder := decodeJSONResponse[inventoryOrder](t, patchRecorder)
	if updatedOrder.Status != inventoryOrderStatusFilled || updatedOrder.CompletedAt == "" {
		t.Fatalf("unexpected filled order: %#v", updatedOrder)
	}

	storeAfterFill, err := loadInventoryRequestStore()
	if err != nil {
		t.Fatalf("load store after fill: %v", err)
	}
	if len(storeAfterFill.Requests) != 1 || storeAfterFill.Requests[0].Status != inventoryRequestStatusFulfilled {
		t.Fatalf("request was not marked fulfilled: %#v", storeAfterFill.Requests)
	}
}

func TestCreateInventoryOrderRejectsMissingRequest(t *testing.T) {
	useTempInventoryRequestStore(t)
	mux := http.NewServeMux()
	registerRoutes(mux)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, jsonRequest(http.MethodPost, "/api/v1/inventory/orders", inventoryOrderCreatePayload{
		Scope:      inventoryRequestScopePersonal,
		RequestIDs: []string{"missing-request"},
	}))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for missing request id, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func jsonRequest(method string, path string, value any) *http.Request {
	body, _ := json.Marshal(value)
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeJSONResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(recorder.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	return value
}
