package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var dbFunctionOIDPattern = regexp.MustCompile(`^[0-9]+$`)

type dbFunctionRow struct {
	OID        string   `json:"oid"`
	Schema     string   `json:"schema"`
	Name       string   `json:"name"`
	Arguments  string   `json:"arguments"`
	ResultType string   `json:"result_type"`
	Language   string   `json:"language"`
	Volatility string   `json:"volatility"`
	Category   string   `json:"category"`
	References []string `json:"references"`
	Summary    string   `json:"summary"`
}

type dbFunctionInspection struct {
	dbFunctionRow
	Definition string   `json:"definition"`
	Risk       string   `json:"risk"`
	Notes      []string `json:"notes"`
}

type msgDBFunctions struct {
	rows []dbFunctionRow
	err  error
}

type msgDBFunctionInspection struct {
	row dbFunctionInspection
	err error
}

func cmdFetchDBFunctions(term, category string) Cmd {
	return func() Msg {
		if globalDB == nil {
			return msgDBFunctions{err: fmt.Errorf("not connected")}
		}
		term = strings.TrimSpace(term)
		category = strings.TrimSpace(category)

		rows, err := globalDB.Query(context.Background(), `
			SELECT p.oid::text,
			       n.nspname,
			       p.proname,
			       pg_get_function_identity_arguments(p.oid),
			       pg_get_function_result(p.oid),
			       l.lanname,
			       p.provolatile::text,
			       COALESCE(obj_description(p.oid, 'pg_proc'), ''),
			       pg_get_functiondef(p.oid)
			FROM pg_proc p
			JOIN pg_namespace n ON n.oid = p.pronamespace
			JOIN pg_language l ON l.oid = p.prolang
			WHERE n.nspname = $1::text
			  AND (
			    $2::text = ''
			    OR p.proname ILIKE '%' || $2::text || '%'
			    OR pg_get_function_identity_arguments(p.oid) ILIKE '%' || $2::text || '%'
			    OR pg_get_function_result(p.oid) ILIKE '%' || $2::text || '%'
			    OR pg_get_functiondef(p.oid) ILIKE '%' || $2::text || '%'
			  )
			ORDER BY p.proname
			LIMIT 300`, dbSchema, term)
		if err != nil {
			return msgDBFunctions{err: err}
		}
		defer rows.Close()

		var out []dbFunctionRow
		for rows.Next() {
			var oid, schema, name, args, resultType, language, volatility, description, definition string
			if err := rows.Scan(&oid, &schema, &name, &args, &resultType, &language, &volatility, &description, &definition); err != nil {
				return msgDBFunctions{err: err}
			}
			row := buildDBFunctionRow(oid, schema, name, args, resultType, language, volatility, description, definition)
			if category != "" && !strings.EqualFold(row.Category, category) {
				continue
			}
			out = append(out, row)
		}
		if err := rows.Err(); err != nil {
			return msgDBFunctions{err: err}
		}
		return msgDBFunctions{rows: out}
	}
}

func cmdInspectDBFunction(oid string) Cmd {
	return func() Msg {
		if globalDB == nil {
			return msgDBFunctionInspection{err: fmt.Errorf("not connected")}
		}
		oid = strings.TrimSpace(oid)
		if !dbFunctionOIDPattern.MatchString(oid) {
			return msgDBFunctionInspection{err: fmt.Errorf("invalid function oid")}
		}

		var schema, name, args, resultType, language, volatility, description, definition string
		err := globalDB.QueryRow(context.Background(), `
			SELECT n.nspname,
			       p.proname,
			       pg_get_function_identity_arguments(p.oid),
			       pg_get_function_result(p.oid),
			       l.lanname,
			       p.provolatile::text,
			       COALESCE(obj_description(p.oid, 'pg_proc'), ''),
			       pg_get_functiondef(p.oid)
			FROM pg_proc p
			JOIN pg_namespace n ON n.oid = p.pronamespace
			JOIN pg_language l ON l.oid = p.prolang
			WHERE p.oid = $1::oid
			  AND n.nspname = $2::text`, oid, dbSchema).Scan(&schema, &name, &args, &resultType, &language, &volatility, &description, &definition)
		if err != nil {
			return msgDBFunctionInspection{err: err}
		}

		base := buildDBFunctionRow(oid, schema, name, args, resultType, language, volatility, description, definition)
		inspection := dbFunctionInspection{
			dbFunctionRow: base,
			Definition:    definition,
			Risk:          classifyDBFunctionRisk(definition),
			Notes:         buildDBFunctionNotes(base, definition),
		}
		return msgDBFunctionInspection{row: inspection}
	}
}

func buildDBFunctionRow(oid, schema, name, args, resultType, language, volatility, description, definition string) dbFunctionRow {
	refs := detectDBFunctionReferences(definition)
	return dbFunctionRow{
		OID:        oid,
		Schema:     schema,
		Name:       name,
		Arguments:  args,
		ResultType: resultType,
		Language:   language,
		Volatility: volatilityLabel(volatility),
		Category:   classifyDBFunctionCategory(name, args, definition),
		References: refs,
		Summary:    summarizeDBFunction(description, refs, definition),
	}
}

func classifyDBFunctionCategory(name, args, definition string) string {
	text := strings.ToLower(name + " " + args + " " + definition)
	checks := []struct {
		category string
		terms    []string
	}{
		{"Item / Inventory", []string{"item", "inventory", "craft", "loot", "durability", "stack_size"}},
		{"Reward / Claim", []string{"reward", "claim", "landsraad"}},
		{"Notification / Event / Queue", []string{"notify", "event", "queue", "message", "outbox", "rabbit", "amqp"}},
		{"Player Movement", []string{"move", "teleport", "partition", "location", "position"}},
		{"Guild / Faction", []string{"guild", "faction", "reputation"}},
		{"Currency", []string{"currency", "solaris", "scrip"}},
		{"Journey / Progression", []string{"journey", "codex", "tutorial", "specialization", "xp", "level"}},
	}
	for _, check := range checks {
		for _, term := range check.terms {
			if strings.Contains(text, term) {
				return check.category
			}
		}
	}
	return "Other"
}

func detectDBFunctionReferences(definition string) []string {
	text := strings.ToLower(definition)
	refs := map[string]bool{}
	checks := map[string][]string{
		"dune.items":                 {"dune.items", " items"},
		"dune.inventories":           {"dune.inventories", " inventories"},
		"rewards/claims":             {"reward", "claim", "landsraad"},
		"notification/event/queue":   {"notify", "event", "queue", "message", "outbox", "rabbit", "amqp"},
		"player_state":               {"player_state"},
		"actors":                     {"actors"},
		"currency":                   {"currency", "solaris", "scrip"},
		"guild/faction":              {"guild", "faction", "reputation"},
		"journey/progression":        {"journey", "codex", "tutorial", "specialization", "xp"},
		"movement/partition/location": {"partition", "teleport", "location", "position"},
	}
	for label, terms := range checks {
		for _, term := range terms {
			if strings.Contains(text, term) {
				refs[label] = true
				break
			}
		}
	}
	out := make([]string, 0, len(refs))
	for ref := range refs {
		out = append(out, ref)
	}
	sort.Strings(out)
	return out
}

func summarizeDBFunction(description string, refs []string, definition string) string {
	if strings.TrimSpace(description) != "" {
		return strings.TrimSpace(description)
	}
	if len(refs) == 0 {
		return "No obvious table or event references detected. Inspect definition before use."
	}
	return "References: " + strings.Join(refs, ", ")
}

func classifyDBFunctionRisk(definition string) string {
	text := strings.ToLower(definition)
	mutatingTerms := []string{"insert ", "update ", "delete ", "truncate ", "drop ", "alter ", "perform ", "execute "}
	for _, term := range mutatingTerms {
		if strings.Contains(text, term) {
			if strings.Contains(text, "notify") || strings.Contains(text, "queue") || strings.Contains(text, "event") || strings.Contains(text, "outbox") {
				return "Mutating with possible live/server side effects"
			}
			return "Mutating"
		}
	}
	return "Read-only or unknown"
}

func buildDBFunctionNotes(row dbFunctionRow, definition string) []string {
	var notes []string
	text := strings.ToLower(definition)
	if strings.Contains(text, "dune.items") || strings.Contains(text, "inventory") {
		notes = append(notes, "Touches item or inventory state; verify whether online players need a relog.")
	}
	if strings.Contains(text, "reward") || strings.Contains(text, "claim") || strings.Contains(text, "landsraad") {
		notes = append(notes, "Touches reward/claim concepts; treat as claim queue semantics unless the function body proves direct inventory mutation.")
	}
	if strings.Contains(text, "notify") || strings.Contains(text, "queue") || strings.Contains(text, "event") || strings.Contains(text, "outbox") {
		notes = append(notes, "References notification/event/queue behavior; this is a candidate for live-player refresh research.")
	}
	if row.Category == "Item / Inventory" && !containsString(row.References, "notification/event/queue") {
		notes = append(notes, "No obvious notification/event reference detected; function may still require relog for online players.")
	}
	if len(notes) == 0 {
		notes = append(notes, "Review function definition and test on a non-production character before using operationally.")
	}
	return notes
}

func volatilityLabel(v string) string {
	switch v {
	case "i":
		return "immutable"
	case "s":
		return "stable"
	case "v":
		return "volatile"
	default:
		return v
	}
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
