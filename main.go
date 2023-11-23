package main

import (
	"github.com/patrickcping/pingone-clean-config/cmd"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
)

func main() {
	l := logger.Get()

	l.Debug().Msg("Starting pingone-cleanconfig")

	cmd.Execute()
}

func init() {}
