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
	index := flag.Int("request", 1, "determines which request to make")
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatalln("missing file operand")
	}

	reqs, err := httpparser.Parse(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	if *index <= 0 {
		log.Fatalln("n must be greater than 0")
	}
	if *index > len(reqs) {
		log.Fatalln(fmt.Sprintf("n must be less than %d", len(reqs)+1))
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(reqs[*index-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if *isVerbose {
		statusColor := color.New(color.Bold, color.Underline)
		if res.StatusCode < 200 {
			statusColor.Add(color.FgMagenta)
		} else if res.StatusCode < 300 {
			statusColor.Add(color.FgGreen)
		} else if res.StatusCode < 400 {
			statusColor.Add(color.FgBlue)
		} else if res.StatusCode < 500 {
			statusColor.Add(color.FgRed)
		} else if res.StatusCode < 600 {
			statusColor.Add(color.FgYellow)
		}
		fmt.Println(statusColor.Sprint(res.Status))

		fmt.Println()

		keyColor := color.New(color.FgRed).SprintFunc()
		valueColor := color.New(color.FgYellow).SprintFunc()
		for key := range res.Header {
			fmt.Printf("%s: %s\n", keyColor(key), valueColor(res.Header.Get(key)))
		}

		fmt.Println()
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body := prettifyBody(raw, res.Header.Get("Content-Type"))
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
