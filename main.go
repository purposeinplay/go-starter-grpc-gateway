package main

import (
	"log"

	"github.com/purposeinplay/go-starter-grpc-gateway/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
