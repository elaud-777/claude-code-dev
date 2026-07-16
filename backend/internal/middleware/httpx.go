package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	apierrors "taskflow-backend/internal/errors"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Wrap adapts a HandlerFunc (which may return an *errors.ApiError) into a
// standard http.HandlerFunc, serializing errors as { error: { code, message, meta? } }.
func Wrap(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			WriteError(w, err)
		}
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierrors.ApiError); ok {
		WriteJSON(w, apiErr.Status, map[string]any{"error": apiErr})
		return
	}
	log.Printf("unexpected error: %v", err)
	WriteJSON(w, 500, map[string]any{"error": map[string]string{"code": "INTERNAL_ERROR", "message": "서버 오류가 발생했습니다"}})
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
