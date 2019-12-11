package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tpryan/mlauth/language"
	"github.com/tpryan/mlauth/speech"
	"github.com/tpryan/mlauth/vision"
)

const b1 = "\033[1m"
const b2 = "\033[0m"
const keyVision = "Golden Retriever"
const keySpeech = "Brooklyn"
const keyLanguage = "location"
const positiveLanguage = true
const secret = "Redeem this key (4D03873C-3729B) for free Google Cloud credits."

const header = `
********************************************************************************
Google Clould ML Scavenger Hunt
                                        
              -:o+oooo+o:-              
          .:ooooooooooooooooo.          
       :oooooooo+o:--:o+oooooooo:       
      oooo++++:          -+ooooooo      
  -:::::::::::::::-         ooooooooo+  
 -::::::-     .:-             .oooooooo 
 ::::::-                        ooooooo 
  ::::::oooooooooooo+++++++++oooooooo+  
     -+sssssssssssssooooooooooooo+:     
        .-::::::::::::::::::---. 

A small experiment that uses ML APIs to make a scavenger hunt. 
********************************************************************************
`

var errorAuthFalse = errors.New("the content did not pass authentication")

var step = 0
var runClear = true

func main() {

	clearScreen()

	reader := bufio.NewReader(os.Stdin)

	for {
		if runClear {
			clearScreen()
			runClear = false

			switch step {
			case 0:
				step00()
			case 1:
				step01()
			case 2:
				step02()
			case 3:
				step03()
			}
		}

		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		input = strings.TrimSuffix(input, "\n")

		if input == "exit" {
			clearScreen()
			fmt.Printf("Thanks for trying it. \n")
			os.Exit(0)
		}

		switch step {
		case 0:
			err = authVision(input)
		case 1:
			err = authSpeech(input)
		case 2:
			err = authLanguage(input)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func clearScreen() {
	fmt.Printf("\033[H\033[2J")
}

func step00() {
	fmt.Printf("%s\n", header)
	o := `The first thing we need is an image of a '%s%s%s'. 

You can take a picture with your camera, or find it in a search, but an image of 
some sort is required. 

Enter the path to a image file that contains a '%s%s%s'.
`
	o = fmt.Sprintf(o, b1, keyVision, b2, b1, keyVision, b2)

	fmt.Printf("%s\n", o)
}

func step01() {
	fmt.Printf("%s\n", header)
	o := `Next we need a file with audio that contains the word '%s%s%s'. 

You can record it with your phone, or find it in a search, but an audio file of 
some sort is required. 

Enter the path to a audio file that contains a '%s%s%s'.`
	o = fmt.Sprintf(o, b1, keySpeech, b2, b1, keySpeech, b2)

	fmt.Printf("%s\n", o)
}

func step02() {
	fmt.Printf("%s\n", header)
	o := `Complete this by entering a positive sentence about a '%s%s%s'. 
Something like 'I love visiting the park'.`
	o = fmt.Sprintf(o, b1, keyLanguage, b2)

	fmt.Printf("%s\n", o)
}

func step03() {
	fmt.Printf("%s\n", header)
	o := `
SUCCESS 
The secret is %s'%s'%s
	`
	o = fmt.Sprintf(o, b1, secret, b2)

	fmt.Printf("%s\n", o)
	os.Exit(0)
}

func authVision(filename string) error {

	fmt.Printf("**************************** %sPicture auth%s ********************************* \n", b1, b2)

	fmt.Printf("File accepted: %s\n", filename)
	fmt.Printf("Testing file...")
	result, err := vision.Auth(keyVision, filename)

	if err != nil {
		return errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		runClear = true
		step++
		fmt.Printf("passed.\n")
		fmt.Printf("Now move on to the next auth step.\n")
		time.Sleep(4 * time.Second)
		return nil
	}
	return errorAuthFalse
}

func authSpeech(filename string) error {
	fmt.Printf("**************************** %sSpeech auth%s ********************************** \n", b1, b2)

	fmt.Printf("File accepted: %s.\n", filename)
	fmt.Printf("Testing file...")
	result, err := speech.Auth(keySpeech, filename)

	if err != nil {
		return errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		runClear = true
		step++
		fmt.Printf("passed.\n")
		fmt.Printf("Now move on to the next auth step.\n")
		time.Sleep(4 * time.Second)
		return nil
	}
	return errorAuthFalse
}

func authLanguage(content string) error {

	fmt.Printf("**************************** %sLanguage auth%s ********************************* \n", b1, b2)

	fmt.Printf("Content accepted: %s.\n", content)

	fmt.Printf("Testing content...")
	result, err := language.Auth(keyLanguage, content, true)

	if err != nil {
		return errors.Wrap(err, "could not authenticate content: ")
	}

	if result {
		runClear = true
		step++
		fmt.Printf("passed.\n")
		time.Sleep(4 * time.Second)
		return nil
	}
	return errorAuthFalse

}
