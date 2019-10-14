package asyncparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ConfigData struct {
	URL         string   `json:"url"`
	Query       string   `json:"query"`
	ParamString string   `json:"paramString"`
	Fields      []string `json:"fields"`
	SizeParam   string   `json:"sizeParam"`
	SizeValue   int      `json:"sizeValue"`
}

// ParseConfigFile populate a ConfigData instance based on a json config file
func ParseConfigFile(path string) (*ConfigData, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully Opened file")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var configData ConfigData

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	if json.Unmarshal(byteValue, &configData) != nil {
		log.Fatal(err)
	}

	return &configData, nil
}

// Benchmark runs a benchmark based on config data
func Benchmark(config *ConfigData, verbose bool) {
	nRequests := 5
	concurrency := 1

	client := RequestClient{
		config.URL,
		config.Query,
		config.ParamString,
		config.SizeParam,
		verbose,
	}

	formatter := NewFormatter()

	ch := make(chan time.Duration, concurrency)

	varyingSizeRequest := func(size int) []float64 {
		return processRun(nRequests, concurrency, ch, func() {
			sendBenchmarkedRequest(&client, config.Fields, size, ch)
		})
	}

	min := 1
	max := 15
	for i := min; i <= max; i++ {
		results := varyingSizeRequest(i)
		formatter.AddSizeResults(i, results)

		fmt.Printf("For size = %d, results in ms => %v\n", i, results)
	}

	formatter.FormatSizes()
}

// sendBenchmarkedRequest send request and feed elapsed time to the channel
func sendBenchmarkedRequest(client *RequestClient, fields []string, size int, c chan time.Duration) {
	start := time.Now()
	client.SendRequest(fields, size)
	c <- time.Since(start)
}

// processRun batch the sending of request based on concurrency param and total number of request
func processRun(nRequests int, concurrency int, ch chan time.Duration, fun func()) []float64 {
	results := make([]float64, 0, nRequests)

	n := nRequests
	for n > 0 {
		for i := 0; i < concurrency; i++ {
			if n > 0 {
				go fun()
				n--
			}
		}

		for i := 0; i < concurrency; i++ {
			if len(results) < nRequests {
				results = append(results, float64(<-ch)/float64(time.Millisecond))
			}
		}
	}

	return results
}
