package main

import (
	"os"

	"github.com/fergloragain/discobeat/cmd"

	_ "github.com/fergloragain/discobeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
