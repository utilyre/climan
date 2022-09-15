package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	chalk "github.com/fatih/color"
	"github.com/pborman/getopt/v2"
	"github.com/utilyre/climan/httpparser"
)

const version string = "0.3.0"

var (
	showHelp    *bool   = getopt.BoolLong("help", 'h', "show help")
	showVersion *bool   = getopt.BoolLong("version", 'V', "show version")
	amVerbose   *bool   = getopt.BoolLong("verbose", 'v', "output verbosely")
	amRaw       *bool   = getopt.BoolLong("raw", 'r', "do not format response body")
	color       *string = getopt.StringLong("color", 0, "auto", "whether to use escape sequences (auto|never|always)", "WHEN")
	index       *int    = getopt.IntLong("index", 'i', 1, "determine which request to make", "NUM")
	timeout     *int    = getopt.IntLong("timeout", 't', 0, "set maximum time for a request to take", "MS")

	filename string = ""
)

func main() {
	setupCLI()

	request := getRequest()
	response := sendRequest(request)
	defer response.Body.Close()

	if *amVerbose {
		printStatus(response)
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

	switch {
	case *showHelp:
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	case *showVersion:
		fmt.Println(version)
		os.Exit(0)
	}

	switch *color {
	case "auto":
	case "never":
		chalk.NoColor = true
	case "always":
		chalk.NoColor = false
	default:
		log.Fatalf("invalid color '%s'\n", *color)
	}

	filename = getopt.Arg(0)
	if filename == "" {
		log.Fatalln("missing file operand")
	} else if len(getopt.Args()) > 1 {
		log.Fatalln("too many files")
	}
}

func getRequest() *http.Request {
	request, err := httpparser.ParseFile(filename, *index)
	if err != nil {
		log.Fatalln(err)
	}

	return request
}

func sendRequest(request *http.Request) *http.Response {
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Millisecond,
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
	protoColor := chalk.New(chalk.Bold, chalk.FgBlue)

	statusColor := chalk.New(chalk.Bold)
	switch {
	case response.StatusCode < 200:
		statusColor.Add(chalk.FgMagenta)
	case response.StatusCode < 300:
		statusColor.Add(chalk.FgGreen)
	case response.StatusCode < 400:
		statusColor.Add(chalk.FgCyan)
	case response.StatusCode < 500:
		statusColor.Add(chalk.FgRed)
	case response.StatusCode < 600:
		statusColor.Add(chalk.FgYellow)
	}

	fmt.Printf("%s %s\n", protoColor.Sprint(response.Proto), statusColor.Sprint(response.Status))
}

func printHeader(response *http.Response) {
	keyColor := chalk.New(chalk.FgRed)
	valueColor := chalk.New(chalk.FgYellow)

	for key := range response.Header {
		fmt.Printf("%s: %s\n", keyColor.Sprint(key), valueColor.Sprint(response.Header.Get(key)))
	}
}

func printBody(response *http.Response) {
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if *amRaw {
		fmt.Println(string(raw))
		return
	}

	contentType := response.Header.Get("Content-Type")
	contentType = strings.Split(contentType, ";")[0]

	switch contentType {
	default:
		fmt.Println(string(raw))

	case "application/json":
		var data any
		if err := json.Unmarshal(raw, &data); err != nil {
			fmt.Println(string(raw))
			return
		}

		prettified, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			fmt.Println(string(raw))
			return
		}

		fmt.Println(string(prettified))
	}
}
