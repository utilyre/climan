package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/droundy/goopt"
	"github.com/fatih/color"
	"github.com/utilyre/climan/httpparser"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("climan: ")

	goopt.Version = "0.0.0"
	goopt.Summary = "climan [OPTIONS] -- <file>"
	goopt.Description = func() string { return "A file based HTTP client" }
	verbose := goopt.Flag([]string{"-v", "--verbose"}, []string{"-q", "--quiet"}, "output verbosely", "be quiet")
	nth := goopt.Int([]string{"-r", "--request"}, 1, "determines which request to do")

	goopt.Parse(nil)
	if len(goopt.Args) != 1 {
		log.Fatalln("missing file")
	}

	reqs, err := httpparser.Parse(goopt.Args[0])
	if err != nil {
		log.Fatalln(err)
	}

	if *nth <= 0 {
		log.Fatalln("n must be greater than 0")
	}
	if *nth > len(reqs) {
		log.Fatalln(fmt.Sprintf("n must be less than %d", len(reqs)+1))
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(reqs[*nth-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if *verbose {
		statusColor := color.New(color.Bold)
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
		statusColor.Println(res.Status)

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
