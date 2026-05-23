package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestAuditMiddlewareCapturesMutatingRequests(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	event := events[0]
	if event.Method != http.MethodPost {
		t.Fatalf("expected POST method, got %q", event.Method)
	}
	if event.Path != "/api/v1/players/give-item" {
		t.Fatalf("unexpected path: %q", event.Path)
	}
	if event.Action != "post:players.give-item" {
		t.Fatalf("unexpected action: %q", event.Action)
	}
	if event.Status != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", event.Status)
	}
	if event.Result != "success" {
		t.Fatalf("expected success result, got %q", event.Result)
	}
}

func TestAuditMiddlewareCapturesFailureStatus(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/players/item/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	if events[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", events[0].Result)
	}
	if events[0].Status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", events[0].Status)
	}
}

func TestAuditMiddlewareSkipsReadOnlyAndPublicRequests(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	requests := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/api/v1/players", nil),
		httptest.NewRequest(http.MethodPost, "/api/v1/public/status", nil),
	}
	for _, req := range requests {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no audit events, got %#v", events)
	}
}

func TestReadAdminAuditEventsLimitAndSort(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	oldEvent := adminAuditEvent{Timestamp: "2026-05-23T10:00:00Z", Method: http.MethodPost, Path: "/api/v1/old", Action: "post:old", Status: 200, Result: "success"}
	newEvent := adminAuditEvent{Timestamp: "2026-05-23T11:00:00Z", Method: http.MethodPost, Path: "/api/v1/new", Action: "post:new", Status: 200, Result: "success"}
	if err := appendAdminAuditEvent(oldEvent); err != nil {
		t.Fatalf("append old event: %v", err)
	}
	if err := appendAdminAuditEvent(newEvent); err != nil {
		t.Fatalf("append new event: %v", err)
	}

	events, err := readAdminAuditEvents(1)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one event after limit, got %d", len(events))
	}
	if events[0].Action != "post:new" {
		t.Fatalf("expected newest event first, got %#v", events[0])
	}
}
