package errs

import "net/http"

type Code struct {
	Code        uint16 `json:"code"`
	Description string `json:"description"`
	Status      int
}

var (
	InternalServerError = Code{Code: 100, Description: "Internal server error", Status: http.StatusInternalServerError}
	InvalidID           = Code{Code: 101, Description: "Invalid ID", Status: http.StatusBadRequest}
)
