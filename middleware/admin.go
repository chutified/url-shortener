package middleware

import (
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AdminAuth middleware checks if a request is authorized.
func AdminAuth(s data.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// load admin key
		key := c.Query("admin_key")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing admin_key",
			})
			c.Abort()
			return
		}

		// validate admin key
		err := s.AdminAuth(c, key)
		if err == data.ErrUnauthorized {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid admin_key",
			})
			c.Abort()
			return

		} else if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unexpected server error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
