package asyncparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestClient struct {
	URL         string
	Query       string
	ParamString string
	SizeParam   string
	Verbose     bool
}

// SendRequest sends a post request to a GraphQL endpoint, returns status code
func (r *RequestClient) SendRequest(fields []string, sizeValue int) (int, error) {
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
		return 0, err
	} else if r.Verbose == true {
		fmt.Println(response.StatusCode)
		if b, err := ioutil.ReadAll(response.Body); err == nil {
			fmt.Println(string(b))
		}
	}

	return response.StatusCode, nil
}
