package httpparser

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const RequestSeparator string = "###"

func Parse(filename string) ([]*http.Request, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	lines = trimSpace(lines)
	lines = removeComments(lines)
	pieces := breakIntoPieces(lines)
	pieces = removeEmptyLines(pieces)

	return extractRequests(pieces)
}

func trimSpace(lines []string) []string {
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return lines
}

func removeComments(lines []string) []string {
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if line == "" {
			continue
		}
		if line == RequestSeparator {
			continue
		}
		if string(line[0]) != "#" {
			continue
		}

		lines = append(lines[:i], lines[i+1:]...)
	}

	return lines
}

func breakIntoPieces(lines []string) [][]string {
	pieces := [][]string{}

	start := 0
	for i, line := range lines {
		if line != RequestSeparator {
			continue
		}

		pieces = append(pieces, lines[start:i])
		start = i + 1
	}
	pieces = append(pieces, lines[start:])

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

func extractRequests(pieces [][]string) ([]*http.Request, error) {
	requests := []*http.Request{}

	for i, lines := range pieces {
		method := ""
		url := ""
		header := map[string]string{}
		body := ""

		hasHeaderEnded := false
		for j, line := range lines {
			if j == 0 {
				parts := strings.Split(line, " ")
				if len(parts) != 2 {
					return nil, errors.New(fmt.Sprintf("httpparser: expected a http method followed by a url separated by one space in the #%d request", i+1))
				}
				if !isValidMethod(parts[0]) {
					return nil, errors.New(fmt.Sprintf("httpparser: expected a http method instead of `%s` in the #%d request", parts[0], i+1))
				}

				method = parts[0]
				url = parts[1]
				continue
			}
			if !hasHeaderEnded {
				if line == "" {
					hasHeaderEnded = true
					continue
				}

				parts := strings.SplitN(line, ": ", 2)
				if len(parts) < 2 {
					return nil, errors.New(fmt.Sprintf("httpparser: expected a key follows by a value in header of #%d request", i+i))
				}

				header[parts[0]] = parts[1]
				continue
			}
			if line == "" {
				continue
			}

			body = strings.Join(lines[j:], "\n")
			break
		}

		request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
		if err != nil {
			return nil, err
		}

		for key, value := range header {
			request.Header.Set(key, value)
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func isValidMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodDelete ||
		method == http.MethodConnect ||
		method == http.MethodOptions ||
		method == http.MethodTrace ||
		method == http.MethodPatch
}
