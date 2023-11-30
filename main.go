package main

import (
	"github.com/patrickcping/pingone-sweep/cmd"
	"github.com/patrickcping/pingone-sweep/internal/logger"
)

func main() {
	l := logger.Get()

	l.Debug().Msg("Starting pingone-sweep")

	cmd.Execute()
}
