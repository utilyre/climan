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

	reqs, err := httpparser.Parse(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if nth <= 0 {
		log.Fatalln("n must be greater than 0")
	}
	if nth > len(reqs) {
		log.Fatalln(fmt.Sprintf("n must be less than %d", len(reqs)))
	}

	res, err := http.DefaultClient.Do(&reqs[nth-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	bufBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var body any
	err = json.Unmarshal(bufBody, &body)
	if err != nil {
		log.Fatalln(err)
	}

	strBody, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range res.Header {
		fmt.Printf("%s: %s\n", k, v[0])
	}

	fmt.Println()
	fmt.Println(string(strBody))
}
