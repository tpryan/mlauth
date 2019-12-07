package main

import (
	"fmt"
	"os"

	"github.com/tpryan/mlauth/language"
)

const key = "Location"
const secret = "ambivalence"

func main() {

	b1 := "\033[1m"
	b2 := "\033[0m"
	fmt.Printf("****** %sLanguage auth%s ****** \n", b1, b2)
	if len(os.Args) < 2 {
		fmt.Printf("You didn't input a sentence to try. \n")
		return
	}

	if len(os.Args) < 3 {
		fmt.Printf("You didn't indicate a sentiment. \n")
		return
	}

	content := os.Args[1]
	fmt.Printf("Content accepted: %s.\n", content)

	positive := false
	sentiment := "Negative"
	if os.Args[2] == "positive" {
		sentiment = "Positive"
		positive = true
	}
	fmt.Printf("Sentiment accepted: %s.\n", sentiment)

	fmt.Printf("Testing content...")
	result, err := language.Auth(key, content, positive)

	if err != nil {
		fmt.Printf(" %sfailed%s. \nThere was an error testing the content: %s.\n", b1, b2, err)
		return
	}

	fmt.Printf(" %sdone%s.\n", b1, b2)

	if result {
		fmt.Printf("The secret is %s`%s'%s.\n", b1, secret, b2)
		return
	}
	fmt.Printf("The file did not unlock the secret.\n")

}
