package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/utilyre/climan/httpparser"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("climan: ")

	var nth int
	flag.IntVar(&nth, "n", 1, "run the nth request of <filename>")
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		log.Fatalln("missing filename")
	}

	requests, err := httpparser.Parse(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if nth <= 0 {
		log.Fatalln("n must be greater than 0")
	}
	if nth > len(requests) {
		log.Fatalln(fmt.Sprintf("n must be less than %d", len(requests)))
	}

	buf, err := doRequest(requests[nth-1])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(buf))
}

func doRequest(req http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data any

	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(data)
}
