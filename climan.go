package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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
	if num < 0 {
		return "", errors.New("natural numbers can only be ordinalized")
	}

	switch num {
	case 1:
		return "1st", nil
	case 2:
		return "2nd", nil
	case 3:
		return "3rd", nil
	}

	return fmt.Sprintf("%dth", num), nil
}
