package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const maxDBQueryParamLength = 128

func handleDBTables(w http.ResponseWriter, r *http.Request) {
	msg, ok := cmdFetchTables().(msgTables)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	type tableOut struct {
		Name     string `json:"name"`
		RowCount int64  `json:"row_count"`
	}
	rows := make([]tableOut, 0, len(msg.rows))
	for _, r := range msg.rows {
		rows = append(rows, tableOut{Name: r.Name, RowCount: r.RowCount})
	}
	jsonOK(w, rows)
}

func handleDBDescribe(w http.ResponseWriter, r *http.Request) {
	table, err := requiredDBQueryParam(r, "table")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdDescribeTable(table)().(msgDescribe)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	type colOut struct {
		Name     string `json:"name"`
		DataType string `json:"data_type"`
		Nullable string `json:"nullable"`
	}
	cols := make([]colOut, 0, len(msg.cols))
	for _, c := range msg.cols {
		cols = append(cols, colOut{Name: c.Name, DataType: c.DataType, Nullable: c.Nullable})
	}
	jsonOK(w, map[string]any{"table": msg.table, "columns": cols})
}

func handleDBSample(w http.ResponseWriter, r *http.Request) {
	table, err := requiredDBQueryParam(r, "table")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	msg, ok := cmdSampleTable(table, limit)().(msgSample)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]any{
		"table":   msg.table,
		"headers": msg.headers,
		"rows":    redactDBStringRows(msg.rows),
	})
}

func handleDBSearch(w http.ResponseWriter, r *http.Request) {
	term, err := requiredDBQueryParam(r, "term")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdSearchColumns(term)().(msgSearchCols)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]any{
		"headers": msg.headers,
		"rows":    redactDBStringRows(msg.rows),
	})
}

func handleDBFunctions(w http.ResponseWriter, r *http.Request) {
	term, err := optionalDBQueryParam(r, "term")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	category, err := optionalDBQueryParam(r, "category")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	msg, ok := cmdFetchDBFunctions(term, category)().(msgDBFunctions)
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
		rows = []dbFunctionRow{}
	}
	jsonOK(w, rows)
}

func handleDBFunctionInspect(w http.ResponseWriter, r *http.Request) {
	oid, err := requiredDBQueryParam(r, "oid")
	if err != nil {
		jsonErr(w, err, 400)
		return
	}
	if _, err := strconv.ParseInt(oid, 10, 64); err != nil {
		jsonErr(w, fmt.Errorf("oid must be numeric"), 400)
		return
	}
	msg, ok := cmdInspectDBFunction(oid)().(msgDBFunctionInspection)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, msg.row)
}

func handleDBSQL(w http.ResponseWriter, r *http.Request) {
	limitBody(w, r, maxJSONBodyBytes)
	var req struct {
		SQL string `json:"sql"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	req.SQL = strings.TrimSpace(req.SQL)
	if req.SQL == "" {
		jsonErr(w, fmt.Errorf("sql required"), 400)
		return
	}
	if !isReadOnlySQL(req.SQL) {
		jsonErr(w, fmt.Errorf("only single-statement read-only SQL is allowed"), 400)
		return
	}
	msg, ok := cmdRunSQL(req.SQL)().(msgSQL)
	if !ok {
		jsonErr(w, fmt.Errorf("internal error"), 500)
		return
	}
	if msg.err != nil {
		jsonErr(w, msg.err, 500)
		return
	}
	jsonOK(w, map[string]string{"result": RedactSensitiveText(msg.result)})
}

func requiredDBQueryParam(r *http.Request, name string) (string, error) {
	value, err := optionalDBQueryParam(r, name)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("%s required", name)
	}
	return value, nil
}

func optionalDBQueryParam(r *http.Request, name string) (string, error) {
	value := strings.TrimSpace(r.URL.Query().Get(name))
	if value == "" {
		return "", nil
	}
	if len(value) > maxDBQueryParamLength {
		return "", fmt.Errorf("%s is too long", name)
	}
	if containsUnsafeControl(value) {
		return "", fmt.Errorf("%s contains unsupported control characters", name)
	}
	return value, nil
}

func redactDBStringRows(rows [][]string) [][]string {
	if rows == nil {
		return nil
	}
	out := make([][]string, len(rows))
	for i, row := range rows {
		out[i] = make([]string, len(row))
		for j, value := range row {
			out[i][j] = RedactSensitiveText(value)
		}
	}
	return out
}
