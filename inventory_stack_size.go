package main

import (
	"context"
	"fmt"
	"net/http"
)

type setItemStackSizeRequest struct {
	ID        int64 `json:"id"`
	StackSize int64 `json:"stack_size"`
}

func handleSetItemStackSize(w http.ResponseWriter, r *http.Request) {
	limitBody(w, r, maxJSONBodyBytes)
	var req setItemStackSizeRequest
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, http.StatusBadRequest)
		return
	}
	msg, ok := cmdSetItemStackSize(req.ID, req.StackSize)().(msgMutate)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), http.StatusInternalServerError)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"ok": msg.ok})
}

func cmdSetItemStackSize(itemID int64, stackSize int64) Cmd {
	return func() Msg {
		if globalDB == nil {
			return msgMutate{err: fmt.Errorf("not connected")}
		}
		if itemID <= 0 {
			return msgMutate{err: fmt.Errorf("item ID required")}
		}
		if stackSize <= 0 {
			return msgMutate{err: fmt.Errorf("stack size must be greater than zero")}
		}
		if stackSize > maxGiveItemStackSize {
			return msgMutate{err: fmt.Errorf("stack size must be <= %d", maxGiveItemStackSize)}
		}
		ctx := context.Background()
		var oldStack int64
		var template string
		err := globalDB.QueryRow(ctx, `
			SELECT stack_size, template_id
			FROM dune.items
			WHERE id = $1::bigint`, itemID).Scan(&oldStack, &template)
		if err != nil {
			return msgMutate{err: fmt.Errorf("find item %d: %w", itemID, err)}
		}
		res, err := globalDB.Exec(ctx, `
			UPDATE dune.items
			SET stack_size = $1::bigint
			WHERE id = $2::bigint`, stackSize, itemID)
		if err != nil {
			return msgMutate{err: fmt.Errorf("set item stack size: %w", err)}
		}
		if res.RowsAffected() == 0 {
			return msgMutate{err: fmt.Errorf("item %d not found", itemID)}
		}
		return msgMutate{ok: fmt.Sprintf("Set item %d stack size from %d to %d (%s)", itemID, oldStack, stackSize, template)}
	}
}
