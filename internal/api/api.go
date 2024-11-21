package api

import (
	"infraction.mageis.net/internal/config"
	"infraction.mageis.net/internal/data"
	"log/slog"
)

type InfractionApi struct {
	logger *slog.Logger
	repos  data.Repositories
	cfg    *config.Config
}

func New(logger *slog.Logger, cfg *config.Config, repos data.Repositories) *InfractionApi {
	return &InfractionApi{
		logger: logger,
		repos:  repos,
		cfg:    cfg,
	}
}
