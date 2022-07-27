package request

import (
	"fmt"
	"net/http"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
}

func (self HTTPRequest) Run() {
	fmt.Println(self.Method)
	fmt.Println(self.URL)
	fmt.Println(self.Headers)
	fmt.Println(self.Body)
	fmt.Println("----------")
}

func ParseHTTP(filename string) ([]HTTPRequest, error) {
	req1 := HTTPRequest{
		Method: http.MethodGet,
		URL:    "https://jsonplaceholder.typicode.com/comments",
	}

	req2 := HTTPRequest{
		Method: http.MethodPost,
		URL:    "https://jsonplaceholder.typicode.com/comments",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]any{},
	}

	return []HTTPRequest{req1, req2}, nil
}
