package main

import (
	"os"

	"github.com/urawa72/pqlt/tui"
)

func start() int {
	if err := tui.New().Run(); err != nil {
		return 1
	}

	return 0
}

func main() {
	os.Exit(start())
}
