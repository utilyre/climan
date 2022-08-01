package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    map[string]any
}

var methods []string = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
	http.MethodPatch,
}

func (self HTTPRequest) Run(data any) error {
	body, err := json.Marshal(self.Body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(self.Method, self.URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for k, v := range self.Headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
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
	// TODO: removeTrailingWhitespace

	requests := []HTTPRequest{}

	for i, lines := range pieces {
		request := HTTPRequest{
			Headers: map[string]string{},
			Body:    map[string]any{},
		}

		haveHeadersEnded := false
		for j, line := range lines {
			if j == 0 {
				parts := strings.Split(line, " ")
				if len(parts) != 2 {
					return nil, errors.New(fmt.Sprintf("expected a http method followed by a url separated by one space in the #%d request", i+1))
				}
				if !slices.Contains(methods, parts[0]) {
					return nil, errors.New(fmt.Sprintf("expected a http method instead of `%s` in the #%d request", parts[0], i+1))
				}

				request.Method = parts[0]
				request.URL = parts[1]
				continue
			}
			if !haveHeadersEnded {
				if line == "" {
					haveHeadersEnded = true
					continue
				}

				parts := strings.SplitN(line, ": ", 2)
				if len(parts) < 2 {
					return nil, errors.New(fmt.Sprintf("expected a key follows by a value in headers of #%d request", i+i))
				}

				request.Headers[parts[0]] = parts[1]
				continue
			}
			if line == "" {
				continue
			}

			body := strings.Join(lines[j:], "\n")
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
