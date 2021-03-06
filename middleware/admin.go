package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chutommy/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AdminLogin validates login data.
func AdminLogin(s data.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load login
		username, ok := c.GetPostForm("username")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing username form field",
			})
			c.Abort()

			return
		}

		password, ok := c.GetPostForm("password")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing password form field",
			})
			c.Abort()

			return
		}

		// authentication
		if err := s.AuthenticateAdmin(username, password); errors.Is(err, data.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Errorf("authentication error: %w", err),
			})
			c.Abort()

			return
		} else if err != nil {
			s.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
			c.Abort()

			return
		}

		c.Next()
	}
}

// ValidateAdminKey middleware checks if a request is authorized.
func ValidateAdminKey(s data.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// load admin key
		key := c.Query("admin_key")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing admin_key query parameter",
			})
			c.Abort()

			return
		}

		// validate admin key
		err := s.ValidateAdminKey(c, key)
		if errors.Is(err, data.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid admin_key query parameter",
			})
			c.Abort()

			return
		} else if err != nil {
			s.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
			c.Abort()

			return
		}

		c.Next()
	}
}
