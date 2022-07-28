package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/utilyre/climan/request"
)

func main() {
	var nth int
	flag.IntVar(&nth, "n", 0, "only run the nth request in the file")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	filename := flag.Arg(0)
	requests, _ := request.ParseHTTP(filename)
	if nth > len(requests) {
		fmt.Printf("Please enter a number less than %d.\n", len(requests))
		os.Exit(1)
	}

	if nth > 0 {
		runRequest(requests[nth-1])
		return
	}

	for _, request := range requests {
		runRequest(request)
	}
}

func runRequest(request request.HTTPRequest) {
	var data any

	err := request.Run(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	prettified, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(prettified))
}
