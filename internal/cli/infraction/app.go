package infraction

import (
	"github.com/urfave/cli/v2"
	v "infraction.mageis.net/internal/version"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:        "infraction",
		Usage:       "the infraction REST-API",
		Description: "Infraction handles splitting costs for group of friends",
		Version:     v.GetVersionString(),

		HideHelpCommand: true,
		Commands: []*cli.Command{
			newServeCommand(),
		},
	}
}
