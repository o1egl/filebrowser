package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func maintainDefaultRequestID() func() {
	original := RequestIDHeader

	return func() {
		RequestIDHeader = original
	}
}

func TestRequestID(t *testing.T) {
	tests := map[string]struct {
		requestIDHeader  string
		request          func() *http.Request
		expectedResponse string
	}{
		"Retrieves Request Id from default header": {
			"X-Request-Id",
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Request-Id", "req-123456")

				return req
			},
			"RequestID: req-123456",
		},
		"Retrieves Request Id from custom header": {
			"X-Trace-Id",
			func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Add("X-Trace-Id", "trace:abc123")

				return req
			},
			"RequestID: trace:abc123",
		},
	}

	defer maintainDefaultRequestID()()

	for _, test := range tests {
		w := httptest.NewRecorder()

		r := gin.New()

		RequestIDHeader = test.requestIDHeader

		r.Use(RequestID)

		r.GET("/", func(c *gin.Context) {
			requestID := GetReqID(c.Request.Context())
			response := fmt.Sprintf("RequestID: %s", requestID)

			_, _ = w.Write([]byte(response))
		})
		r.ServeHTTP(w, test.request())

		if w.Body.String() != test.expectedResponse {
			t.Fatalf("RequestID was not the expected value")
		}
	}
}
