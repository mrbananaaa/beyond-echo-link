package response

import "net/http"

type ErrorResponse struct {
	Error struct {
		Message string   `json:"message"`
		Errors  []string `json:"erors"`
	} `json:"error"`
}

func Error(
	w http.ResponseWriter,
	status int,
	msg string,
	errs []string,
) {
	resp := ErrorResponse{}
	resp.Error.Message = msg
	resp.Error.Errors = errs

	JSON(w, status, resp)
}

func BadRequest(
	w http.ResponseWriter,
	msg string,
	errs []string,
) {
	Error(
		w,
		http.StatusBadRequest,
		msg,
		errs,
	)
}

func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error", nil)
}
