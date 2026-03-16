package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string            `json:"error,omitempty"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func ReadJson(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dst)
}

func WriteError(w http.ResponseWriter, status int, errR ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errR)
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
