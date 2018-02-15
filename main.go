package main

import (
	"log"

	"github.com/makpoc/hades-sheet/server"
)

func main() {
	err := server.Start()
	if err != nil {
		log.Fatalf("Server start returned an error: %v", err)
	}
}
