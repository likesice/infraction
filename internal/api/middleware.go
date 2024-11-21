package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	ierrors "infraction.mageis.net/internal/errors"
	"log/slog"
	"net/http"
	"time"
)

func (api *InfractionApi) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		host := c.Request.Host
		route := c.FullPath()
		end := time.Now()
		latency := end.Sub(start)
		ip := c.ClientIP()
		referer := c.Request.Referer()

		baseAttributes := []slog.Attr{}

		requestAttributes := []slog.Attr{
			slog.Time("time", start),
			slog.String("method", method),
			slog.String("host", host),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.String("route", route),
			slog.String("ip", ip),
			slog.String("referer", referer),
		}

		responseAttributes := []slog.Attr{
			slog.Time("time", end),
			slog.Duration("latency", latency),
			slog.Int("status", status),
		}

		level := slog.LevelInfo
		msg := ""
		var errorAttributes = []slog.Attr{}
		if len(c.Errors) != 0 {
			for _, err := range c.Errors {
				var errTyped ierrors.InfractionErrorI
				var cause slog.Attr
				if errors.As(err, &errTyped) {
					if errTyped.Unwrap() != nil {
						cause = slog.String("ierror_cause", errTyped.Unwrap().Error())
					}
					errorAttributes = append(errorAttributes, []slog.Attr{
						slog.String("ierror_code", errTyped.GetCode()),
						slog.String("ierror_message", errTyped.GetMessage()),
						slog.String("ierror_description", errTyped.GetDescription()),
						slog.Any("ierror_args", errTyped.GetArgs()),
						cause,
					}...)
				} else {
					errTyped = ierrors.ErrUnspecified.Wrap(err)
					if errTyped.Unwrap() != nil {
						cause = slog.String("ierror_cause", errTyped.Unwrap().Error())
					}
					errorAttributes = append(errorAttributes, []slog.Attr{
						slog.String("ierror_code", errTyped.GetCode()),
						slog.String("ierror_message", errTyped.GetMessage()),
						slog.String("ierror_description", errTyped.GetDescription()),
					}...)
				}
			}
		}

		attributes := append(
			[]slog.Attr{
				{
					Key:   "request",
					Value: slog.GroupValue(requestAttributes...),
				},
				{
					Key:   "response",
					Value: slog.GroupValue(responseAttributes...),
				},
				{
					Key:   "ierror",
					Value: slog.GroupValue(errorAttributes...),
				},
			},
			baseAttributes...,
		)
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			level = slog.LevelInfo
			msg = c.Errors.String()
		} else if status >= http.StatusInternalServerError {
			level = slog.LevelError
			msg = c.Errors.String()
		}

		api.logger.LogAttrs(c.Request.Context(), level, msg, attributes...)
	}
}
