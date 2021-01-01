package controller

import (
	"errors"
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// GenerateAdminKey handles a new admin_key generation.
func (h *handler) GenerateAdminKey(c *gin.Context) {
	// generate a new key
	key, err := h.ds.GenerateAdminKey(c)
	if err != nil {
		h.ds.LogError(c, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": data.ErrUnexpectedError,
		})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"admin_key": key,
	})
}

// RevokeAdminKey handles admin_key's cancellation.
func (h *handler) RevokeAdminKey(c *gin.Context) {
	// load prefix
	prefix := c.Query("prefix")
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing prefix query parameter",
		})
		return
	}

	// revoke
	if err := h.ds.RevokeAdminKey(c, prefix); err != nil {
		if errors.Is(err, data.ErrPrefixNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.ds.LogError(c, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": data.ErrUnexpectedError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"revoked_prefix": prefix,
	})
}
