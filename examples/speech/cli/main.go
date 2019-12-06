package main

import (
	"fmt"
	"os"

	"github.com/tpryan/mlauth/speech"
)

const key = "Brooklyn"
const secret = "ambivalence"

func main() {
	b1 := "\033[1m"
	b2 := "\033[0m"
	fmt.Printf("****** %sSpeech auth%s ****** \n", b1, b2)
	if len(os.Args) < 2 {
		fmt.Printf("You didn't indicate a file to try. \n")
		return
	}

	arg := os.Args[1]
	fmt.Printf("File accepted: %s.\n", arg)
	fmt.Printf("Testing file...")
	result, err := speech.Auth(key, arg)

	if err != nil {
		fmt.Printf(" %sfailed%s. \nThere was an error testing the file: %s.\n", b1, b2, err)
		return
	}
	fmt.Printf(" %sdone%s.\n", b1, b2)

	if result {
		fmt.Printf("The secret is %s`%s'%s.\n", b1, secret, b2)
		return
	}
	fmt.Printf("The file did not unlock the secret.\n")
	return
}
