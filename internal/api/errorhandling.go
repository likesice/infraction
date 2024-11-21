package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	ierrors "infraction.mageis.net/internal/errors"
)

func handleIError(err error, c *gin.Context) {
	var errTyped ierrors.InfractionErrorI
	result := gin.H{"error": ierrors.ErrUnspecified}

	c.Error(err)

	if errors.As(err, &errTyped) {
		switch errTyped.GetKind() {
		case ierrors.USER_ERR:
			{
				result = gin.H{"error": errTyped}
			}
		case ierrors.SYS_ERR:
			{
				errTyped = ierrors.ErrUnspecified.Wrap(errTyped)
				result = gin.H{"error": errTyped}
			}
		}

	}
	c.AbortWithStatusJSON(errTyped.GetHttpCode(), result)
}
