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

	asyncparser.SendRequest(config)

	// readFile("")
	// process("", []string{"", ""}, 1)
}

func readFile(filePath string) {}

func process(url string, fields []string, sizeParam int) {
	return
}

func makeRequest(url string, fields []string, sizeParam int) {
	// response, err := http.Get(fmt.Sprintf("%s/%s", url, accountID))
	// response.StatusCode
}
