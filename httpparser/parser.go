package httpparser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Parse(filename string) ([]http.Request, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	lines = trimSpace(lines)
	lines = removeComments(lines)
	pieces := breakIntoPieces(lines)
	pieces = removeEmptyLines(pieces)

	reqs := []http.Request{}

	for i, lines := range pieces {
		var method string
		var url string
		var header = header{}
		var body body

		hasHeaderEnded := false
		for j, line := range lines {
			if j == 0 {
				parts := strings.Split(line, " ")
				if len(parts) != 2 {
					return nil, errors.New(fmt.Sprintf("expected a http method followed by a url separated by one space in the #%d request", i+1))
				}
				if !isValidMethod(parts[0]) {
					return nil, errors.New(fmt.Sprintf("expected a http method instead of `%s` in the #%d request", parts[0], i+1))
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
					return nil, errors.New(fmt.Sprintf("expected a key follows by a value in header of #%d request", i+i))
				}

				header[parts[0]] = parts[1]
				continue
			}
			if line == "" {
				continue
			}

			err := json.Unmarshal([]byte(strings.Join(lines[j:], "\n")), &body)
			if err != nil {
				return nil, err
			}
			break
		}

		buf, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(method, url, bytes.NewBuffer(buf))
		if err != nil {
			return nil, err
		}

		for k, v := range header {
			req.Header.Set(k, v)
		}

		reqs = append(reqs, *req)
	}

	return reqs, nil
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
		if string(line[0]) != "#" || line == "###" {
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
		if line != "###" {
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
