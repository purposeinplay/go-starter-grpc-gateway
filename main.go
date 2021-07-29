package main

import (
	"fmt"

	"github.com/purposeinplay/go-starter-grpc-gateway/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
