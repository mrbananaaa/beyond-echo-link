package response

import (
	"encoding/json"
	"net/http"
)

func JSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
}
