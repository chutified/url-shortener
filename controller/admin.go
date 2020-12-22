package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenerateAdminKey handles a new admin_key generation.
func (h *handler) GenerateAdminKey(c *gin.Context) {

	// generate a new key
	key, err := h.ds.GenerateAdminKey(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unexpected internal server error",
		})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"admin_key": key,
	})
}
