package main

import (
	"fmt"

	"github.com/asciifaceman/godock/pkg/server"
)

func main() {

	s, err := server.NewServer()
	if err != nil {
		fmt.Printf("ERRO: Failed to start server - %v", err)
	}

	err = s.RegisterHandlers()
	if err != nil {
		fmt.Printf("ERRO: Failed to register handlers - %v", err)
	}

	err = s.Run()
	if err != nil {
		fmt.Printf("ERRO: Failed to start server - %v", err)
	}

	return
}
