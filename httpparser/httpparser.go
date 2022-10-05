package httpparser

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const RequestSeparator string = "###"

func Parse(src string, index int) (*http.Request, error) {
	if index == 0 {
		return nil, errors.New("httpparser: index must be nonzero")
	}

	lines := strings.Split(src, "\n")
	lines = trimSpaces(lines)
	lines = removeComments(lines)
	pieces := breakIntoPieces(lines)

	index, err := normalizeIndex(len(pieces), index)
	if err != nil {
		return nil, err
	}

	pieces = removeEmptyLines(pieces)
	piece := applyEnv(pieces[index])

	request, err := extractRequest(piece)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func ParseFile(name string, index int) (*http.Request, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return Parse(string(content), index)
}

func trimSpaces(lines []string) []string {
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

func normalizeIndex(length, index int) (int, error) {
	if index < -length || index > length {
		return 0, fmt.Errorf("httpparser: index must be within %d and %d", -length, length)
	}
	if index < 0 {
		return length + index, nil
	}

	return index - 1, nil
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

func applyEnv(lines []string) []string {
	text := strings.Join(lines, "\n")
	rgx := regexp.MustCompile("[^\\\\]\\$\\{.+?\\}")

	text = rgx.ReplaceAllStringFunc(text, func(match string) string {
		dollarIndex := strings.Index(match, "$")
		prefix := match[:dollarIndex]

		key := match[dollarIndex+2 : len(match)-1]
		env, ok := os.LookupEnv(key)
		if !ok {
			return match
		}

		return prefix + env
	})

	return strings.Split(text, "\n")
}

func extractRequest(lines []string) (*http.Request, error) {
	method := ""
	url := ""
	header := map[string]string{}
	body := ""

	hasHeaderEnded := false
	for i, line := range lines {
		if i == 0 {
			parts := strings.Split(line, " ")
			if len(parts) != 2 {
				return nil, errors.New("httpparser: expected a http method followed by a url separated by one space")
			}
			if !isValidMethod(parts[0]) {
				return nil, fmt.Errorf("httpparser: expected a http method instead of '%s'", parts[0])
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
				return nil, errors.New("httpparser: expected a key follows by a value in header")
			}

			header[parts[0]] = parts[1]
			continue
		}
		if line == "" {
			continue
		}

		body = strings.Join(lines[i:], "\n")
		break
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	return request, nil
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
