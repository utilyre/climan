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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	var nth int
	flag.IntVar(&nth, "n", 0, "run the nth request of <filename>")

	flag.Parse()
	if nth <= 0 {
		panic("n must be greater than 0")
	}

	filename := flag.Arg(0)
	if filename == "" {
		panic("missing filename")
	}

	requests, err := request.ParseHTTP(filename)
	if err != nil {
		panic(err)
	}
	if nth > len(requests) {
		panic(fmt.Sprintf("n must be less than %d", len(requests)))
	}

	var data any
	err = requests[nth-1].Run(&data)
	if err != nil {
		panic(err)
	}

	raw, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(raw))
}
