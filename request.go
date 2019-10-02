package asyncparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"gonum.org/v1/gonum/stat"
)

type RequestClient struct {
	URL         string
	Query       string
	ParamString string
	SizeParam   string
	Verbose     bool
}

type ConfigData struct {
	URL         string   `json:"url"`
	Query       string   `json:"query"`
	ParamString string   `json:"paramString"`
	Fields      []string `json:"fields"`
	SizeParam   string   `json:"sizeParam"`
	SizeValue   int      `json:"sizeValue"`
}

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

func (r *RequestClient) sendRequest(fields []string, sizeValue int) {
	type Wrapper struct {
		OperationName *string   `json:"operationName"`
		Variables     *struct{} `json:"variables"`
		Query         *string   `json:"query"`
	}

	params := fmt.Sprintf("%s,%s:%d", r.ParamString, r.SizeParam, sizeValue)

	body := fmt.Sprintf("{%s(%s){%s}}", r.Query, params, strings.Join(fields, ","))
	var variables struct{}
	payload := Wrapper{nil, &variables, &body}
	b, err := json.Marshal(payload)

	if r.Verbose == true {
		fmt.Println(string(b))
	}

	response, err := http.Post(r.URL, "application/json", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println(fmt.Errorf("%T %s", err, err))
	} else if r.Verbose == true {
		fmt.Println(response.StatusCode)
		if b, err := ioutil.ReadAll(response.Body); err == nil {
			fmt.Println(string(b))
		}
	}
}

func Benchmark(config *ConfigData) {
	nRequests := 10
	concurrency := 3
	// - all fields
	// - fields one by one
	// - default count
	// - various count values
	// ==>> measure time
	// ==>> async (how many simulatenous)

	client := RequestClient{
		config.URL,
		config.Query,
		config.ParamString,
		config.SizeParam,
		false,
	}

	ch := make(chan time.Duration, concurrency)

	varyingSizeRequest := func(size int) []float64 {
		return processRun(nRequests, concurrency, ch, func() {
			sendBenchmarkedRequest(&client, config.Fields, size, ch)
		})
	}

	min := 1
	max := 10
	for i := min; i <= max; i++ {
		results := varyingSizeRequest(i)

		fmt.Printf("For size = %d\n", i)
		fmt.Println(results)
		sort.Float64s(results)
		fmt.Printf("Mean: %fms\n", stat.Mean(results, nil))
		fmt.Printf("Median: %fms\n", stat.Quantile(0.5, stat.Empirical, results, nil))
	}

	// results := processRun(nRequests, concurrency, ch, func() {
	// sendBenchmarkedRequest(&client, config.Fields, config.SizeValue, ch)
	// })

	// fmt.Println(results)
	// sort.Float64s(results)
	// fmt.Printf("Mean: %fms\n", stat.Mean(results, nil))
	// fmt.Printf("Median: %fms\n", stat.Quantile(0.5, stat.Empirical, results, nil))
}

// func varyRequestSize(ch chan time.Duration, min int, max int, fun func(size int) []float64) []float64 {
// 	for i := min; i <= max; i++ {
// 		results := fun(1)
// 	}
// }

func sendBenchmarkedRequest(client *RequestClient, fields []string, size int, c chan time.Duration) {
	start := time.Now()
	client.sendRequest(fields, size)
	c <- time.Since(start)
}

func processRun(nRequests int, concurrency int, ch chan time.Duration, fun func()) []float64 {
	results := []float64{}

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
