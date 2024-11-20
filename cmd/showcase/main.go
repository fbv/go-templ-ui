package main

import (
	"log"

	"github.com/fbv/go-templ-ui/server"
)

func main() {
	s := server.New()
	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
