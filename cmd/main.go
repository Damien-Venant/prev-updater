package main

import (
	"fmt"

	"github.com/prev-updater/internal/infra"
)

func init() {
	config := infra.LoadConfig()
	fmt.Println(config)
}

func main() {
	fmt.Println("Hello World")
}
