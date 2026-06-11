package main

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"sync"
)

const defaultAdminAuditPath = "admin-audit.jsonl"

var adminAuditMu sync.Mutex

func adminAuditPath() string {
	if path := strings.TrimSpace(os.Getenv("ADMIN_AUDIT_LOG")); path != "" {
		return path
	}
	return defaultAdminAuditPath
}

func appendAdminAuditEvent(event adminAuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	adminAuditMu.Lock()
	defer adminAuditMu.Unlock()

	file, err := os.OpenFile(adminAuditPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	return err
}

func readAdminAuditEvents(limit int) ([]adminAuditEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	file, err := os.Open(adminAuditPath())
	if os.IsNotExist(err) {
		return []adminAuditEvent{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []adminAuditEvent
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event adminAuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			events = append(events, event)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].Timestamp > events[j].Timestamp })
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}
