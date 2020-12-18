package controller

import (
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AdminAuth middleware checks if a request is authorized.
func (h *handler) AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// load admin key
		key, err := c.Cookie("admin_key")
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing admin_key",
			})
			return
		}

		// validate admin key
		err = h.ds.AdminAuth(key)
		if err == data.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid admin_key",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unexpected server error",
			})
			return
		}

		c.Next()
	}
}
