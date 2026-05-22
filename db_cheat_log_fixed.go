package main

import (
	"context"
	"fmt"
)

func cmdFetchCheatLogFixed() Cmd {
	return func() Msg {
		if globalDB == nil {
			return msgCheatLog{err: fmt.Errorf("not connected")}
		}
		rows, err := globalDB.Query(context.Background(), `
			SELECT ct.fls_id,
			       ct.cheat_type::text,
			       to_char(ct.event_time AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:MI:SS'),
			       COALESCE(ps.character_name, ct.fls_id)
			FROM dune.cheater_tracking ct
			LEFT JOIN dune.encrypted_accounts e
			       ON convert_from(e.encrypted_funcom_id, 'UTF8') = ct.fls_id
			LEFT JOIN dune.player_state ps
			       ON ps.account_id = e.id
			WHERE ct.event_time > NOW() - INTERVAL '7 days'
			ORDER BY ct.event_time DESC
			LIMIT 500`)
		if err != nil {
			return msgCheatLog{err: err}
		}
		defer rows.Close()

		var out []cheatEntry
		for rows.Next() {
			var r cheatEntry
			if err := rows.Scan(&r.FLSID, &r.CheatType, &r.EventTime, &r.CharacterName); err != nil {
				continue
			}
			out = append(out, r)
		}
		if err := rows.Err(); err != nil {
			return msgCheatLog{err: err}
		}
		return msgCheatLog{rows: out}
	}
}
