package main

import (
	"fmt"
	"io"
	"os"

	"github.com/prev-updater/internal/infra"
)

var (
	TOKEN string = ""
)

func init() {
	config := infra.LoadConfig()
	fmt.Println(config)

	if file, err := os.OpenFile("./api-key", os.O_RDONLY, os.ModeAppend); err != nil {
		panic(err)
	} else {
		defer file.Close()
		result, _ := io.ReadAll(file)
		TOKEN = string(result)
	}
}

func main() {
	fmt.Print(TOKEN)
}
