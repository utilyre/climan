package main

import (
	"errors"
	"flag"
	"fmt"
)

func main() {
	nth := flag.Int("nth", -1, "only run the nth request in the file")
	flag.Parse()

	if ordinalized, err := ordinalize(*nth); err == nil {
		fmt.Printf("Running the %s request...\n", ordinalized)
		fmt.Println("{\"message\":\"Hello world!\"}")
		return
	}

	fmt.Println("Running all requests...")
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
