package main

import (
	"bytes"
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
	amRaw     *bool   = getopt.BoolLong("raw", 'r', "do not try to parse response body")
	color     *string = getopt.StringLong("color", 0, "auto", "determine when to use escape sequences", "WHEN")
	index     *int    = getopt.IntLong("index", 'i', 1, "determine which request to make", "NUM")

	filename string = ""
)

func main() {
	setupCLI()

	request := getRequest()
	response := sendRequest(request)
	defer response.Body.Close()

	if *amVerbose {
		printStatus(response)
		fmt.Println()

		printHeader(response)
		fmt.Println()
	}

	printBody(response)
}

func setupCLI() {
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

	filename = getopt.Arg(0)
	if filename == "" {
		log.Fatalln("missing file operand")
	} else if len(getopt.Args()) > 1 {
		log.Fatalln("too many files")
	}
}

func getRequest() *http.Request {
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

	return requests[*index-1]
}

func sendRequest(request *http.Request) *http.Response {
	client := &http.Client{
		CheckRedirect: func(request *http.Request, requests []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	return response
}

func printStatus(response *http.Response) {
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
}

func printHeader(response *http.Response) {
	keyColor := chalk.New(chalk.FgRed)
	valueColor := chalk.New(chalk.FgYellow)

	for key := range response.Header {
		fmt.Printf("%s: %s\n", keyColor.Sprint(key), valueColor.Sprint(response.Header.Get(key)))
	}
}

func printBody(response *http.Response) {
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if *amRaw {
		goto resort
	}

	switch response.Header.Get("Content-Type") {
	case "application/json":
		var prettified bytes.Buffer
		if err = json.Indent(&prettified, raw, "", "\t"); err != nil {
			goto resort
		}

		fmt.Println(string(prettified.Bytes()))
		return
	}

resort:
	fmt.Println(string(raw))
}
