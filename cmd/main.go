package main

import (
	"flag"
	"fmt"

	"github.com/tdidierjean/asyncparser"
)

func main() {
	var verbose = flag.Bool("verbose", false, "Show extra output")
	var filePath = flag.String("file", "", "Path to config file")
	flag.Parse()

	config, err := asyncparser.ParseConfigFile(*filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	asyncparser.Benchmark(config, *verbose)
}

func readFile(filePath string) {}

func process(url string, fields []string, sizeParam int) {
	return
}
