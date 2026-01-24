package main

import (
	"os"

	"github.com/btafoya/gomailserver/internal/commands"
)

func main() {
	if err := commands.Execute("dev"); err != nil {
		os.Exit(1)
	}
}
