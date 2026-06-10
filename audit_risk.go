package main

import (
	"net/http"
	"strings"
)

var highRiskMutationPathMarkers = []string{
	"/reconnect",
	"/battlegroup/exec",
	"/database/sql",
	"/logs/stream-ticket",
	"/notify",
	"/players/item/",
	"/give-item",
	"/give-currency",
	"/give-faction-rep",
	"/give-scrip",
	"/grant-live",
	"/award-xp",
	"/award-char-xp",
	"/award-intel",
	"/kick",
	"/repair-item",
	"/teleport",
	"/journey/complete",
	"/set-faction",
	"/set-spec-xp",
	"/storage/",
}

func auditActionName(method, path string) string {
	name := strings.TrimPrefix(path, "/api/v1/")
	name = strings.Trim(name, "/")
	name = strings.ReplaceAll(name, "/", ".")
	if name == "" {
		name = "root"
	}
	return strings.ToLower(method) + ":" + name
}

func mutationRiskForRequest(method, path string) string {
	lower := strings.ToLower(path)
	if method == http.MethodDelete || strings.Contains(lower, "/wipe") || strings.Contains(lower, "/delete") || strings.Contains(lower, "/reset") || strings.Contains(lower, "/blueprints/import") {
		return "destructive"
	}
	for _, marker := range highRiskMutationPathMarkers {
		if strings.Contains(lower, marker) {
			return "high"
		}
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		return "medium"
	}
	return "low"
}
