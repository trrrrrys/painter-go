package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/trrrrrys/painter-go/internal/painter"
)

func main() {
	if err := painter.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
