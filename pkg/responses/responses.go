package responses

import "github.com/gin-gonic/gin"

type IResponse interface {
	Success(code int, data any)
	Error(code int, traceId, message string)
}

type response struct {
	StatusCode int
	Data       any
	ErrorRes   *responseError
	Context    *gin.Context
}

type responseError struct {
	TraceId string `json:"trace_id"`
	Message string `json:"message"`
}

func NewResponse(c *gin.Context) IResponse {
	return &response{Context: c}
}

func (r *response) Success(code int, data any) {
	r.StatusCode = code
	r.Data = data
	r.Context.JSON(r.StatusCode, gin.H{"data": r.Data})
}

func (r *response) Error(code int, traceId, message string) {
	r.StatusCode = code
	r.ErrorRes = &responseError{
		TraceId: traceId,
		Message: message,
	}
	r.Context.JSON(r.StatusCode, gin.H{"error": r.ErrorRes})
}
