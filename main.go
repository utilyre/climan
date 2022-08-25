package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	chalk "github.com/fatih/color"
	"github.com/pborman/getopt/v2"
	"github.com/utilyre/climan/httpparser"
)

var (
	showHelp  *bool   = getopt.BoolLong("help", 'h', "show help")
	amVerbose *bool   = getopt.BoolLong("verbose", 'v', "output verbosely")
	color     *string = getopt.StringLong("color", 0, "auto", "determine when to use escape sequences", "WHEN")
	index     *int    = getopt.IntLong("index", 'i', 1, "determine which request to make", "NUM")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("climan: ")

	getopt.SetParameters("FILE")
	getopt.Parse()
	if *showHelp {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	switch *color {
	case "never":
		chalk.NoColor = true
	case "always":
		chalk.NoColor = false
	}

	filename := getopt.Arg(0)
	if filename == "" {
		log.Fatalln("missing file operand")
	} else if len(getopt.Args()) > 1 {
		log.Fatalln("too many files")
	}

	requests, err := httpparser.Parse(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if *index == 0 {
		log.Fatalln("index can not be zero")
	} else if *index < -len(requests) || *index > len(requests) {
		log.Fatalf("index must be greater than %d and less than %d\n", -len(requests)-1, len(requests)+1)
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

	if *amVerbose {
		statusColor := chalk.New(chalk.Bold, chalk.Underline)
		if response.StatusCode < 200 {
			statusColor.Add(chalk.FgMagenta)
		} else if response.StatusCode < 300 {
			statusColor.Add(chalk.FgGreen)
		} else if response.StatusCode < 400 {
			statusColor.Add(chalk.FgBlue)
		} else if response.StatusCode < 500 {
			statusColor.Add(chalk.FgRed)
		} else if response.StatusCode < 600 {
			statusColor.Add(chalk.FgYellow)
		}
		fmt.Println(statusColor.Sprint(response.Status))

		fmt.Println()

		keyColor := chalk.New(chalk.FgRed).SprintFunc()
		valueColor := chalk.New(chalk.FgYellow).SprintFunc()
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
