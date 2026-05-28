package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const logStreamTicketRawBytes = 32
const logStreamTicketTTL = 60 * time.Second

type logStreamTicket struct {
	Namespace string
	Pod       string
	ExpiresAt time.Time
}

var logStreamTickets = struct {
	sync.Mutex
	values map[string]logStreamTicket
}{values: map[string]logStreamTicket{}}

type logStreamTicketRequest struct {
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
}

type logStreamTicketResponse struct {
	Ticket    string    `json:"ticket"`
	ExpiresAt time.Time `json:"expires_at"`
}

func generateLogStreamTicketValue() (string, error) {
	b := make([]byte, logStreamTicketRawBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate log stream ticket: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func pruneExpiredLogStreamTicketsLocked(now time.Time) {
	for ticket, info := range logStreamTickets.values {
		if !now.Before(info.ExpiresAt) {
			delete(logStreamTickets.values, ticket)
		}
	}
}

func issueLogStreamTicket(namespace, pod string) (logStreamTicketResponse, error) {
	if namespace == "" || pod == "" {
		return logStreamTicketResponse{}, fmt.Errorf("namespace and pod are required")
	}
	if !isValidRuntimeLogTarget(namespace, pod) {
		return logStreamTicketResponse{}, fmt.Errorf("invalid log stream target")
	}

	ticket, err := generateLogStreamTicketValue()
	if err != nil {
		return logStreamTicketResponse{}, err
	}
	info := logStreamTicket{Namespace: namespace, Pod: pod, ExpiresAt: time.Now().UTC().Add(logStreamTicketTTL)}

	logStreamTickets.Lock()
	defer logStreamTickets.Unlock()
	pruneExpiredLogStreamTicketsLocked(time.Now().UTC())
	logStreamTickets.values[ticket] = info

	return logStreamTicketResponse{Ticket: ticket, ExpiresAt: info.ExpiresAt}, nil
}

func consumeLogStreamTicket(ticket, namespace, pod string) bool {
	if !adminTokenPattern.MatchString(ticket) || namespace == "" || pod == "" {
		return false
	}

	now := time.Now().UTC()
	logStreamTickets.Lock()
	defer logStreamTickets.Unlock()
	pruneExpiredLogStreamTicketsLocked(now)

	info, ok := logStreamTickets.values[ticket]
	if !ok {
		return false
	}
	delete(logStreamTickets.values, ticket)

	if !now.Before(info.ExpiresAt) {
		return false
	}
	return info.Namespace == namespace && info.Pod == pod
}

func validateAndConsumeLogStreamTicket(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}
	return consumeLogStreamTicket(r.URL.Query().Get("ticket"), r.URL.Query().Get("ns"), r.URL.Query().Get("pod"))
}

func handleIssueLogStreamTicket(w http.ResponseWriter, r *http.Request) {
	limitBody(w, r, maxJSONBodyBytes)
	var req logStreamTicketRequest
	if err := decode(r, &req); err != nil {
		jsonErr(w, fmt.Errorf("invalid ticket request: %w", err), http.StatusBadRequest)
		return
	}
	resp, err := issueLogStreamTicket(req.Namespace, req.Pod)
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	jsonOK(w, resp)
}
