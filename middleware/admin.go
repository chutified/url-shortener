package middleware

import (
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AdminLogin validates login data.
func AdminLogin(s data.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// load login
		username, ok := c.GetPostForm("username")
		if !ok {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing username field",
			})
			c.Abort()
			return
		}

		password, ok := c.GetPostForm("password")
		if !ok {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing password field",
			})
			c.Abort()
			return
		}

		// authentication
		if err := s.AuthenticateAdmin(c, username, password); err == data.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
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

// ValidateAdminKey middleware checks if a request is authorized.
func ValidateAdminKey(s data.Service) gin.HandlerFunc {
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
		err := s.ValidateAdminKey(c, key)
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
