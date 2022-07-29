package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    map[string]any
}

func (self HTTPRequest) Run(data any) error {
  // TODO: Handle all http methods with `NewRequestWithContext`
	switch self.Method {
	case http.MethodGet:
		response, err := http.Get(self.URL)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		err = unmarshalBody(response.Body, data)

	case http.MethodPost:
		body, err := json.Marshal(self.Body)
		if err != nil {
			return err
		}

		contentType, ok := self.Headers["Content-Type"]
		if !ok {
			return errors.New("missing `Content-Type` header")
		}

		response, err := http.Post(self.URL, contentType, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer response.Body.Close()

		err = unmarshalBody(response.Body, data)
	}

	return nil
}

func ParseHTTP(filename string) ([]HTTPRequest, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	pieces := breakIntoPieces(lines)
	pieces = removeComments(pieces)
	pieces = removeEmptyLines(pieces)

	requests := []HTTPRequest{}

	for _, lines := range pieces {
		request := HTTPRequest{
			Headers: map[string]string{},
			Body:    map[string]any{},
		}

		haveHeadersEnded := false
		for i, line := range lines {
			if i == 0 {
				parts := strings.Split(line, " ")
				request.Method = parts[0]
				request.URL = parts[1]
				continue
			}
			if !haveHeadersEnded {
				if line == "" {
					haveHeadersEnded = true
					continue
				}

				parts := strings.Split(line, ": ")
				request.Headers[parts[0]] = parts[1]
				continue
			}
			if line == "" {
				continue
			}

			body := strings.Join(lines[i:], "\n")
			err := json.Unmarshal([]byte(body), &request.Body)
			if err != nil {
				return nil, err
			}
			break
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func unmarshalBody(body io.ReadCloser, data any) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	return nil
}

func breakIntoPieces(lines []string) [][]string {
	pieces := [][]string{}

	start := 0
	for i, line := range lines {
		if line != "###" {
			continue
		}

		pieces = append(pieces, lines[start:i])
		start = i + 1
	}
	pieces = append(pieces, lines[start:])

	return pieces
}

func removeComments(pieces [][]string) [][]string {
	for _, lines := range pieces {
		for i := len(lines) - 1; i >= 0; i-- {
			line := lines[i]
			if line == "" {
				continue
			}
			if string(line[0]) != "#" || line == "###" {
				continue
			}

			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	return pieces
}

func removeEmptyLines(pieces [][]string) [][]string {
	for i, lines := range pieces {
		firstMeaningful := 0
		for j, line := range lines {
			if line == "" {
				continue
			}

			firstMeaningful = j
			break
		}

		lastMeaningful := len(lines) - 1
		for j := len(lines) - 1; j >= 0; j-- {
			line := lines[j]
			if line == "" {
				continue
			}

			lastMeaningful = j
			break
		}

		pieces[i] = lines[firstMeaningful : lastMeaningful+1]
	}

	return pieces
}
