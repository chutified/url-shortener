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
		v1.GET("/url/i/:record_short", h.GetRecordByShortPeek)

		authorized := v1.Group("/admin", h.AdminAuth())
		{
			authorized.GET("/url/short/:record_short", h.GetRecordByShort)
			authorized.GET("/url/id/:record_id", h.GetRecordByID)

			authorized.GET("/urls/l", h.GetRecordsLen)
			authorized.GET("/urls", h.GetAllRecords)

			authorized.POST("/url", h.AddRecord)
			authorized.PUT("/url/:record_id", h.UpdateRecord)
			authorized.DELETE("/url/:record_id", h.DeleteRecord)
			authorized.POST("/url/recovery/:record_id", h.RecordRecovery)
		}
	}

	return r
}
