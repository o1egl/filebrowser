package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WrapHH converts middleware to the gin middleware handler.
func WrapHH(fn func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var skip = true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			skip = false
		}
		fn(handler).ServeHTTP(c.Writer, c.Request)
		switch {
		case skip:
			c.Abort()
		default:
			c.Next()
		}
	}
}
