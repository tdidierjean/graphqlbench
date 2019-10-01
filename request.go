package asyncparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func SendRequest(config *ConfigData) {

	type Wrapper struct {
		OperationName *string   `json:"operationName"`
		Variables     *struct{} `json:"variables"`
		Query         *string   `json:"query"`
	}

	params := fmt.Sprintf("%s,%s:%d", config.ParamString, config.SizeParam, config.SizeValue)

	body := fmt.Sprintf("{%s(%s){%s}}", config.Query, params, strings.Join(config.Fields, ","))
	var variables struct{}
	payload := Wrapper{nil, &variables, &body}
	b, err := json.Marshal(payload)

	fmt.Println(string(b))

	response, err := http.Post(config.URL, "application/json", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println(fmt.Errorf("%T %s", err, err))
	} else {
		fmt.Println(response.StatusCode)
		if b, err := ioutil.ReadAll(response.Body); err == nil {
			fmt.Println(string(b))
		}
	}
}
