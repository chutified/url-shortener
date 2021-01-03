package controller

import (
	"errors"
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AddRecord adds a new record.
func (h *handler) AddRecord(c *gin.Context) {
	// bind record
	var newRecord data.Record
	err := c.ShouldBindJSON(&newRecord)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// add record
	r, err := h.ds.AddRecord(c, &newRecord)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidRecord):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrUnavailableShort):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// record successfully added
	c.JSON(http.StatusOK, r)
}

// UpdateRecord replace the record with given the certain ID.
func (h *handler) UpdateRecord(c *gin.Context) { // get record's ID
	id := c.Param("record_id")

	// bind record
	var newRecord data.ShortRecord
	err := c.ShouldBindJSON(&newRecord)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// update record
	r, err := h.ds.UpdateRecord(c, id, &newRecord)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUnavailableShort):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrIDNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		// server error
		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// record successfully updated
	c.JSON(http.StatusOK, r)
}

// DeleteRecord removes the record with the certain ID.
func (h *handler) DeleteRecord(c *gin.Context) { // get record's ID
	id := c.Param("record_id")

	// delete record
	did, err := h.ds.DeleteRecord(c, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrIDNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// record successfully deleted
	c.JSON(http.StatusOK, gin.H{
		"delete_record_id": did,
	})
}

// GetRecordByID serves a record with the certain ID.
func (h *handler) GetRecordByID(c *gin.Context) { // get record's ID
	id := c.Param("record_id")

	// get record
	r, err := h.ds.GetRecordByID(c, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrIDNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// record successfully retrieved
	c.JSON(http.StatusOK, r)
}

// GetRecordByShort serves a record with the certain Short value.
func (h *handler) GetRecordByShort(c *gin.Context) { // get record's Short
	short := c.Param("record_short")

	// get record
	r, err := h.ds.GetRecordByShort(c, short)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrShortNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// record successfully retrieved
	c.JSON(http.StatusOK, r)
}

// GetRecordByShortPeek serves a full url of the shortcut.
func (h *handler) GetRecordByShortPeek(c *gin.Context) { // get short
	short := c.Param("record_short")

	// get full url
	full, err := h.ds.GetRecordByShortPeek(c, short)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrShortNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			h.ds.LogError(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": data.ErrUnexpectedError,
			})
		}

		return
	}

	// found
	c.JSON(http.StatusOK, gin.H{
		"url": full,
	})
}

// GetRecordsLen returns a total number of records.
func (h *handler) GetRecordsLen(c *gin.Context) { // get length
	l, err := h.ds.GetRecordsLen(c)
	if err != nil {
		h.ds.LogError(c, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": data.ErrUnexpectedError,
		})

		return
	}

	// records length successfully retrieved
	c.JSON(http.StatusOK, gin.H{
		"len": l,
	})
}

// GetAllRecords returns xth page with a certain number of records.
func (h *handler) GetAllRecords(c *gin.Context) {
	// get records
	rs, err := h.ds.GetAllRecords(c)
	if err != nil {
		h.ds.LogError(c, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": data.ErrUnexpectedError,
		})

		return
	}

	// records successfully retrieved
	c.JSON(http.StatusOK, rs)
}

// RecordRecovery tries to recover a softly deleted record.
func (h *handler) RecordRecovery(c *gin.Context) { // load id
	id := c.Param("record_id")

	// recover a record
	rid, err := h.ds.RecordRecovery(c, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotDeleted):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, data.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		// server error
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	// record successfully recovered
	c.JSON(http.StatusOK, gin.H{
		"recovered_id": rid,
	})
}
