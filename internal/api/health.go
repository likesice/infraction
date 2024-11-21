package api

import (
	"github.com/gin-gonic/gin"
	"infraction.mageis.net/internal/version"
)

func (api *InfractionApi) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"health": map[string]interface{}{
		"status":      "available",
		"environment": api.cfg.Env,
		"version":     version.GetVersion(),
	}})
}
