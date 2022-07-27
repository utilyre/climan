package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	nth := flag.Int("nth", -1, "only run the nth request in the file")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	filename := flag.Arg(0)

	if ordinalized, err := ordinalize(*nth); err == nil {
		fmt.Printf("Running the %s request of %s...\n\n", ordinalized, filename)

		fmt.Println("{\"message\":\"Hello world!\"}")
		return
	}

	fmt.Printf("Running all requests of %s...\n", filename)
}

func ordinalize(num int) (string, error) {
	if num <= 0 {
		return "", errors.New("natural numbers can only be ordinalized")
	}

	strNum := strconv.Itoa(num)
	switch strNum[len(strNum)-1:] {
	case "1":
		return fmt.Sprintf("%dst", num), nil
	case "2":
		return fmt.Sprintf("%dnd", num), nil
	case "3":
		return fmt.Sprintf("%drd", num), nil
	}

	return fmt.Sprintf("%dth", num), nil
}
