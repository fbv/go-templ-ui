package main

import (
	"log"

	"github.com/fbv/go-templ-ui/server"
	"github.com/fbv/go-templ-ui/view"
)

func main() {
	s := server.New(view.Content)
	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
