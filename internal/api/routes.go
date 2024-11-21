package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *InfractionApi) Routes() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(api.RequestLogger())

	r.GET("/health", api.healthHandler)

	v1 := r.Group("/v1")
	{
		v1.POST("/infraction", api.createInfractionHandler)
		v1.GET("/infraction/:id", api.getInfractionHandler)
		v1.DELETE("/infraction/:id", api.deleteInfractionHandler)
		v1.GET("/infraction", api.getAllInfractionHandler)
		v1.POST("/infraction/:id/transaction", api.createTransactionHandler)
	}
	return r
}
