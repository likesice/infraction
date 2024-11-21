package api

import (
	"github.com/gin-gonic/gin"
	"infraction.mageis.net/internal/data"
	"infraction.mageis.net/internal/data/validator"
	ierrors "infraction.mageis.net/internal/errors"
	"net/http"
	"strconv"
)

func (api *InfractionApi) createInfractionHandler(c *gin.Context) {
	var infraction data.Infraction
	err := c.ShouldBindBodyWithJSON(&infraction)
	if err != nil {
		handleIError(err, c)
		return
	}

	v := validator.New()
	if valid := infraction.Validate(v); !valid {
		ierr := ierrors.ErrValidationFailed
		ierr.Args = v.Errors
		handleIError(ierr, c)
		return
	}
	err = api.repos.Infractions.Insert(&infraction)
	if err != nil {
		handleIError(err, c)
		return
	}
	c.Header("location", "/infraction/"+strconv.FormatInt(infraction.Id, 10))
	c.JSONP(http.StatusCreated, infraction)
}

func (api *InfractionApi) getInfractionHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	infraction, err := api.repos.Infractions.Select(int64(id))
	if err != nil {
		handleIError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": infraction})
}

func (api *InfractionApi) getAllInfractionHandler(c *gin.Context) {
	infractions, err := api.repos.Infractions.SelectAll()
	if err != nil {
		handleIError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": infractions})
}

func (api *InfractionApi) deleteInfractionHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := api.repos.Infractions.Delete(int64(id))
	if err != nil {
		handleIError(err, c)
		return
	}
	c.Status(http.StatusOK)
}
