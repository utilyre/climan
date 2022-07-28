package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
}

func (self HTTPRequest) Run(data any) error {
	switch self.Method {
	case http.MethodGet:
		response, err := http.Get(self.URL)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseHTTP(filename string) ([]HTTPRequest, error) {
	req1 := HTTPRequest{
		Method: http.MethodGet,
		URL:    "https://jsonplaceholder.typicode.com/comments",
	}

	req2 := HTTPRequest{
		Method: http.MethodGet,
		URL:    "https://jsonplaceholder.typicode.com/comments/1",
	}

	req3 := HTTPRequest{
		Method: http.MethodPost,
		URL:    "https://jsonplaceholder.typicode.com/comments",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]any{
			"id":   1,
			"name": "Utilyre",
			"body": "This is Amirabbas.",
		},
	}

	return []HTTPRequest{req1, req2, req3}, nil
}
