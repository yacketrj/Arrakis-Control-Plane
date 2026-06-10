package main

import (
	"bytes"
	"io"
	"net/http"
)

type inspectedRequestBody struct {
	Data     []byte
	TooLarge bool
	Empty    bool
}

func inspectAndRestoreRequestBody(r *http.Request, maxBytes int64) (inspectedRequestBody, error) {
	if r == nil || r.Body == nil || r.ContentLength == 0 {
		return inspectedRequestBody{Empty: true}, nil
	}
	if r.ContentLength > maxBytes {
		return inspectedRequestBody{TooLarge: true}, nil
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBytes+1))
	r.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return inspectedRequestBody{Data: body}, err
	}
	if int64(len(body)) > maxBytes {
		return inspectedRequestBody{Data: body, TooLarge: true}, nil
	}
	return inspectedRequestBody{Data: body, Empty: len(bytes.TrimSpace(body)) == 0}, nil
}
