package http

const (
	INTERVAL_SERVER_ERROR = "A problem has occurred"
	TIMEOUT_ERROR         = "Server is busy"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
