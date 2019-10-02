package main

import (
	"fmt"
	"os"

	"github.com/tdidierjean/asyncparser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing argument: path to config file")
		return
	}

	filePath := os.Args[1]
	config, err := asyncparser.ParseConfigFile(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	asyncparser.Benchmark(config)
}

func readFile(filePath string) {}

func process(url string, fields []string, sizeParam int) {
	return
}
