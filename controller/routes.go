package controller

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *handler) getHTTPHandler() http.Handler {

	// set router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// V1
	v1 := r.Group("/v1")
	{
		v1.GET("/book/short/:record_short", h.GetRecordByShort)
		v1.GET("/book/full/:record_full", h.GetRecordByFull)
		v1.GET("/book/id/:record_id", h.GetRecordByID)

		v1.GET("/books/len", h.GetRecordsLen)
		v1.GET("/books", h.GetAllRecords)

		v1.POST("/book", h.AddRecord)
		v1.PUT("/book/:record_id", h.UpdateRecord)
		v1.DELETE("/book/:record_id", h.DeleteRecord)
	}

	return r
}
