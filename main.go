package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/utilyre/climan/httpparser"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("climan: ")

	flag.Usage = func() {
		fmt.Println("climan - A file based HTTP client")

		fmt.Println()

		fmt.Println("Usage:")
		fmt.Println("  climan [OPTIONS]... -- FILE")

		fmt.Println()

		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	isVerbose := flag.Bool("verbose", false, "output verbosely")
	index := flag.Int("index", 0, "determines which request to make")
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatalln("missing file operand")
	}

	requests, err := httpparser.Parse(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	if *index < -len(requests) || *index > len(requests) {
		log.Fatalf("index must be greater than %d and less than %d\n", -len(requests)-1, len(requests)+1)
	}
	if *index == 0 {
		log.Fatalln("index can not be zero")
	}

	if *index < 0 {
		*index = len(requests) + 1 + *index
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(requests[*index-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	if *isVerbose {
		statusColor := color.New(color.Bold, color.Underline)
		if response.StatusCode < 200 {
			statusColor.Add(color.FgMagenta)
		} else if response.StatusCode < 300 {
			statusColor.Add(color.FgGreen)
		} else if response.StatusCode < 400 {
			statusColor.Add(color.FgBlue)
		} else if response.StatusCode < 500 {
			statusColor.Add(color.FgRed)
		} else if response.StatusCode < 600 {
			statusColor.Add(color.FgYellow)
		}
		fmt.Println(statusColor.Sprint(response.Status))

		fmt.Println()

		keyColor := color.New(color.FgRed).SprintFunc()
		valueColor := color.New(color.FgYellow).SprintFunc()
		for key := range response.Header {
			fmt.Printf("%s: %s\n", keyColor(key), valueColor(response.Header.Get(key)))
		}

		fmt.Println()
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body := prettifyBody(raw, response.Header.Get("Content-Type"))
	fmt.Println(body)
}

func prettifyBody(raw []byte, contentType string) string {
	switch contentType {
	case "application/json":
		prettified, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return string(raw)
		}

		return string(prettified)
	}

	return string(raw)
}
