package controller

import "github.com/gin-gonic/gin"

func (h *handler) AddRecord(c *gin.Context) {
	//TODO
}

func (h *handler) UpdateRecord(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")
	//TODO
}

func (h *handler) DeleteRecord(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")
	//TODO
}

func (h *handler) GetRecordById(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")
	//TODO
}

func (h *handler) GetRecordByShort(c *gin.Context) {

	// get record's Short
	short := c.Param("record_short")
	//TODO
}

func (h *handler) GetRecordByFull(c *gin.Context) {

	// get record's Full
	full := c.Param("record_full")
	//TODO
}

func (h *handler) GetRecordsLen(c *gin.Context) {
	// TODO
}

func (h *handler) GetAllRecords(c *gin.Context) {

	// get pagination data
	p := c.DefaultQuery("page", "1")
	pagin := c.DefaultQuery("pagin", "30")
	//TODO
}
