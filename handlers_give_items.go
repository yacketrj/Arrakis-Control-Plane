package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

const (
	maxGiveItemRows      = 100
	maxGiveItemQty       = 9999
	maxGiveItemStackSize = 9999
)

type giveItemEntry struct {
	Template  string `json:"template"`
	Qty       int64  `json:"qty"`
	Quality   int64  `json:"quality"`
	StackSize int64  `json:"stack_size"`
}

type giveItemsRequest struct {
	PlayerID int64           `json:"player_id"`
	Template string          `json:"template"`
	Qty      int64           `json:"qty"`
	Quality  int64           `json:"quality"`
	Items    []giveItemEntry `json:"items"`
}

func normalizeGiveItemsRequest(req giveItemsRequest) ([]giveItemEntry, error) {
	items := req.Items
	if len(items) == 0 && strings.TrimSpace(req.Template) != "" {
		items = []giveItemEntry{{Template: req.Template, Qty: req.Qty, Quality: req.Quality, StackSize: 1}}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}
	if len(items) > maxGiveItemRows {
		return nil, fmt.Errorf("maximum %d item rows per request", maxGiveItemRows)
	}
	for i := range items {
		items[i].Template = strings.TrimSpace(items[i].Template)
		if items[i].Template == "" {
			return nil, fmt.Errorf("item %d template required", i+1)
		}
		if items[i].Qty <= 0 {
			return nil, fmt.Errorf("item %d quantity must be > 0", i+1)
		}
		if items[i].Qty > maxGiveItemQty {
			return nil, fmt.Errorf("item %d quantity must be <= %d", i+1, maxGiveItemQty)
		}
		if items[i].StackSize <= 0 {
			items[i].StackSize = 1
		}
		if items[i].StackSize > maxGiveItemStackSize {
			return nil, fmt.Errorf("item %d stack size must be <= %d", i+1, maxGiveItemStackSize)
		}
		if items[i].Quality < 0 || items[i].Quality > 5 {
			return nil, fmt.Errorf("item %d quality must be 0-5", i+1)
		}
		if items[i].Qty > math.MaxInt64/items[i].StackSize {
			return nil, fmt.Errorf("item %d total quantity is too large", i+1)
		}
	}
	return items, nil
}

func handleGiveItems(w http.ResponseWriter, r *http.Request) {
	limitBody(w, r, maxJSONBodyBytes)
	var req giveItemsRequest
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	if req.PlayerID == 0 {
		jsonErr(w, fmt.Errorf("player ID required"), http.StatusBadRequest)
		return
	}
	items, err := normalizeGiveItemsRequest(req)
	if err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	results := make([]string, 0, len(items))
	legacySingleItem := len(req.Items) == 0 && strings.TrimSpace(req.Template) != ""
	for i, item := range items {
		var msg msgMutate
		var ok bool
		if legacySingleItem {
			msg, ok = cmdGiveItem(req.PlayerID, item.Template, item.Qty, item.Quality)().(msgMutate)
		} else {
			msg, ok = cmdGiveItemStacks(req.PlayerID, item.Template, item.Qty, item.StackSize, item.Quality)().(msgMutate)
		}
		if !ok {
			jsonErr(w, fmt.Errorf("internal error"), http.StatusInternalServerError)
			return
		}
		if msg.err != nil {
			jsonErr(w, fmt.Errorf("item %d %s: %w", i+1, item.Template, msg.err), http.StatusInternalServerError)
			return
		}
		results = append(results, fmt.Sprintf("%d stack(s) x %d of %s grade %d", item.Qty, item.StackSize, item.Template, item.Quality))
	}
	jsonOK(w, map[string]string{"ok": fmt.Sprintf("Added %d item row(s): %s", len(results), strings.Join(results, "; "))})
}
