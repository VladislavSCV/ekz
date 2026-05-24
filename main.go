package main

import (
	"fmt"
	"os"

	"github.com/VladislavSCV/ekz/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}
