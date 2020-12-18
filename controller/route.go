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
		v1.GET("/url/short/:record_short", h.GetRecordByShort)
		v1.GET("/url/id/:record_id", h.GetRecordByID)
		v1.GET("/url/i/:record_short", h.GetRecordByShortPeek)

		v1.GET("/urls/l", h.GetRecordsLen)
		v1.GET("/urls", h.GetAllRecords)

		v1.POST("/url", h.AddRecord)
		v1.PUT("/url/:record_id", h.UpdateRecord)
		v1.DELETE("/url/:record_id", h.DeleteRecord)
		v1.POST("/url/recovery/:record_id", h.RecordRecovery)
	}

	return r
}
