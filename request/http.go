package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

		err = unmarshalBody(response.Body, data)

	case http.MethodPost:
		body, err := json.Marshal(self.Body)
		if err != nil {
			return err
		}

		response, err := http.Post(self.URL, self.Headers["Content-Type"], bytes.NewBuffer(body))
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
	pieces = formatPieces(pieces)

	return []HTTPRequest{}, nil
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

func formatPieces(pieces [][]string) [][]string {
	for _, lines := range pieces {
		for i, line := range lines {
			if line == "" {
				continue
			}
			if string(line[0]) != "#" || line == "###" {
				continue
			}

			lines = append(lines[:i], lines[i+1:]...)
		}
	}

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
