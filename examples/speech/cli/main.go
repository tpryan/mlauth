// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/json"
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

	verbose := ""
	input := os.Args[1]
	if len(os.Args) > 2 {
		verbose = os.Args[2]
	}

	fmt.Printf("File accepted: %s.\n", input)
	fmt.Printf("Testing file...")
	result, err := speech.Auth(key, input)

	if err != nil {
		fmt.Printf(" %sfailed%s. \nThere was an error testing the file: %s.\n", b1, b2, err)
		if verbose == "-v" {
			txt, _ := json.MarshalIndent(result, " ", "   ")
			fmt.Printf("%s\n", txt)
		}
		return
	}
	fmt.Printf(" %sdone%s.\n", b1, b2)

	if result.Result {
		fmt.Printf("The secret is %s`%s'%s.\n", b1, secret, b2)
		if verbose == "-v" {
			txt, _ := json.MarshalIndent(result, " ", "   ")
			fmt.Printf("%s\n", txt)
		}
		return
	}
	fmt.Printf("The file did not unlock the secret.\n")
	if verbose == "-v" {
		txt, _ := json.MarshalIndent(result, " ", "   ")
		fmt.Printf("%s\n", txt)
	}
	return
}
