package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handleMutationSafetyClassify(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("method")))
	path := strings.TrimSpace(r.URL.Query().Get("path"))
	if method == "" {
		method = http.MethodPost
	}
	if path == "" {
		jsonErr(w, fmt.Errorf("path is required"), http.StatusBadRequest)
		return
	}
	jsonOK(w, mutationSafetyForPath(method, path))
}
