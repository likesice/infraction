package api

import (
	"github.com/gin-gonic/gin"
	"infraction.mageis.net/internal/data"
	"infraction.mageis.net/internal/data/validator"
	ierrors "infraction.mageis.net/internal/errors"
	"strconv"
)

func (api *InfractionApi) createTransactionHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var transaction data.Transaction
	err := c.ShouldBindBodyWithJSON(&transaction)
	if err != nil {
		handleIError(err, c)
		return
	}

	v := validator.New()
	if valid := transaction.Validate(v); !valid {
		ierr := ierrors.ErrValidationFailed
		ierr.Args = v.Errors
		handleIError(ierr, c)
		return
	}
	err = api.repos.Transactions.Insert(&transaction, int64(id))
	if err != nil {
		handleIError(err, c)
		return
	}
	c.Header("location", "/infraction/"+strconv.FormatInt(int64(id), 10))
	c.JSON(201, transaction)
}
