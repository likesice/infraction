package main

import (
	"log"
	"os"

	cli "infraction.mageis.net/internal/cli/infraction"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
