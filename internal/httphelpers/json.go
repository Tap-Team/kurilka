package httphelpers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	return err
}
