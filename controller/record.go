package controller

import (
	"net/http"

	"github.com/chutified/url-shortener/data"
	"github.com/gin-gonic/gin"
)

// AddRecord adds a new record.
func (h *handler) AddRecord(c *gin.Context) {

	// bind record
	var newr data.Record
	err := c.ShouldBindJSON(&newr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add record
	r, err := h.ds.AddRecord(c, &newr)
	if err != nil {
		switch err {

		// invalid record (e.g. missing keys)
		case data.ErrInvalidRecord:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		// short already in use
		case data.ErrUnavailableShort:
			c.JSON(http.StatusConflict, gin.H{
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

	// record successfully added
	c.JSON(http.StatusOK, r)
}

// UpdateRecord replace the record with given the certain ID.
func (h *handler) UpdateRecord(c *gin.Context) {

	// get record's ID
	// id := c.Param("record_id")
	//TODO
}

// DeleteRecord removes the record with the certain ID.
func (h *handler) DeleteRecord(c *gin.Context) {

	// get record's ID
	// id := c.Param("record_id")
	//TODO
}

// GetRecordById serves a record with the certain ID.
func (h *handler) GetRecordById(c *gin.Context) {

	// get record's ID
	// id := c.Param("record_id")
	//TODO
}

// GetRecordByShort serves a record with the certain Short value.
func (h *handler) GetRecordByShort(c *gin.Context) {

	// get record's Short
	// short := c.Param("record_short")
	//TODO
}

// GetRecordByFull serves a record with the sertain Full value.
func (h *handler) GetRecordByFull(c *gin.Context) {

	// get record's Full
	// full := c.Param("record_full")
	//TODO
}

// GetRecordsLen returns a total number of records.
func (h *handler) GetRecordsLen(c *gin.Context) {
	// TODO
}

// func (h *handler) GetAllRecords(c *gin.Context) {

//     // get pagination data
//     p := c.DefaultQuery("page", "1")
//     pagin := c.DefaultQuery("pagin", "30")
//     //TODO
// }
