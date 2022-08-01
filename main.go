package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/utilyre/climan/request"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("climan: ")

	var nth int
	flag.IntVar(&nth, "n", 0, "run the nth request of <filename>")

	flag.Parse()
	if nth <= 0 {
		log.Fatalln("n must be greater than 0")
	}

	filename := flag.Arg(0)
	if filename == "" {
		log.Fatalln("missing filename")
	}

	requests, err := request.ParseHTTP(filename)
	if err != nil {
		log.Fatalln(err)
	}
	if nth > len(requests) {
		log.Fatalln(fmt.Sprintf("n must be less than %d", len(requests)))
	}

	var data any
	err = requests[nth-1].Run(&data)
	if err != nil {
		log.Fatalln(err)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(raw))
}
