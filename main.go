package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/droundy/goopt"
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

	res, err := http.DefaultClient.Do(reqs[*nth-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if *verbose {
		fmt.Println(res.Status)

		for key := range res.Header {
			fmt.Printf("%s: %s\n", key, res.Header.Get(key))
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
