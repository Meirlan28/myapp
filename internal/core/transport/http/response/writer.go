package response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type StatusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *StatusRecorder) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *StatusRecorder) StatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.statusCode
}
