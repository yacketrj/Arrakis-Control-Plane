package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type templateOut struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func handleGetPlayers(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchPlayers().(msgPlayers)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []playerInfo{}
	}
	jsonOK(w, rows)
}

func handleGetOnlineState(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchOnlineState().(msgOnlineState)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	// Serialize as JSON-friendly structs
	type onlineRow struct {
		PlayerID int64  `json:"player_id"`
		Name     string `json:"name"`
		Map      string `json:"map"`
		Status   string `json:"status"`
		LastSeen string `json:"last_seen"`
	}
	rows := make([]onlineRow, 0, len(msg.rows))
	for _, r := range msg.rows {
		rows = append(rows, onlineRow{
			PlayerID: r.PlayerID,
			Name:     r.Name,
			Map:      r.Map,
			Status:   r.Status,
			LastSeen: r.LastSeen,
		})
	}
	jsonOK(w, rows)
}

func handleGetCurrency(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchCurrency().(msgCurrency)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []currencyRow{}
	}
	jsonOK(w, rows)
}

func handleGetFactions(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchFactions().(msgFactions)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []factionRep{}
	}
	jsonOK(w, rows)
}

func handleGetSpecs(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchSpecs().(msgSpecs)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []specTrack{}
	}
	jsonOK(w, rows)
}

func templateRows() []templateOut {
	rows := make([]templateOut, len(dbItemTemplates))
	for i, t := range dbItemTemplates {
		rows[i] = templateOut{ID: t, Name: itemData.Names[strings.ToLower(t)]}
	}
	return rows
}

func handleGetTemplates(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, templateRows())
}

func handleRefreshTemplates(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchItemTemplates().(msgItemTemplates)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	mergeItemTemplates(msg.templates)
	jsonOK(w, templateRows())
}

func handleGetInventory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonErr(w, fmt.Errorf("invalid id"), 400)
		return
	}
	msg, ok := cmdFetchInventory(id)().(msgInventory)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []itemInfo{}
	}
	jsonOK(w, rows)
}

func handleGetJourney(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.PathValue("id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		jsonErr(w, fmt.Errorf("invalid accountId"), 400)
		return
	}
	msg, ok := cmdFetchJourneyNodes(accountID)().(msgJourney)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	rows := msg.rows
	if rows == nil {
		rows = []journeyNode{}
	}
	// Serialize with JSON tags (journeyNode fields are unexported-named but have tags in new model.go)
	type jNode struct {
		NodeID           string `json:"node_id"`
		IsComplete       bool   `json:"is_complete"`
		IsRevealed       bool   `json:"is_revealed"`
		HasPendingReward bool   `json:"has_pending_reward"`
	}
	out := make([]jNode, 0, len(rows))
	for _, n := range rows {
		out = append(out, jNode{
			NodeID:           n.NodeID,
			IsComplete:       n.IsComplete,
			IsRevealed:       n.IsRevealed,
			HasPendingReward: n.HasPendingReward,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func handleGiveItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID int64  `json:"player_id"`
		Template string `json:"template"`
		Qty      int64  `json:"qty"`
		Quality  int64  `json:"quality"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdGiveItem(req.PlayerID, req.Template, req.Qty, req.Quality)().(msgMutate)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]string{"ok": msg.ok})
}

func handleGrantLive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ControllerID int64  `json:"controller_id"`
		Template     string `json:"template"`
		Amount       int64  `json:"amount"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	if req.Template == "" {
		jsonErr(w, fmt.Errorf("template required"), 400)
		return
	}
	if req.Amount <= 0 {
		req.Amount = 1
	}
	msg, ok := cmdGrantLive(req.ControllerID, req.Template, req.Amount)().(msgMutate)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]string{"ok": msg.ok})
}

func handleGiveCurrency(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID int64 `json:"player_id"`
		Amount   int64 `json:"amount"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdGiveCurrency(req.PlayerID, req.Amount)().(msgMutate)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]string{"ok": msg.ok})
}

func handleGiveFactionRep(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ActorID   int64 `json:"actor_id"`
		FactionID int16 `json:"faction_id"`
		Delta     int32 `json:"delta"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdGiveFactionRep(req.ActorID, req.FactionID, req.Delta)().(msgMutate)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]string{"ok": msg.ok})
}
