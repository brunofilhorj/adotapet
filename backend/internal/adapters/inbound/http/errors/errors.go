package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	TraceID string         `json:"traceId,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func NotImplemented(w http.ResponseWriter, feature string) {
	WriteJSON(w, http.StatusNotImplemented, ErrorResponse{
		Code:    "NOT_IMPLEMENTED",
		Message: feature + " ainda nao foi implementado.",
	})
}
