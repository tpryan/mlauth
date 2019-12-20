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

// Package language wraps the Cloud Natural Language API and provides an auth
// method that allows you to check an input text for an input term.
package language

import (
	"context"
	"fmt"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Auth takes string of text and a term and compares them to each other to see
// if the an item matching the sentiment is contained in the audio file
func Auth(entityType, source string, positive bool) (AuthResult, error) {

	res, err := findSentiment(source)

	if err != nil {
		return res, err
	}

	if err := res.AuthTerm(entityType, positive); err != nil {
		return res, err
	}

	return res, nil

}

func findSentiment(source string) (AuthResult, error) {
	var err error
	ctx := context.Background()
	res := AuthResult{}

	client, err := language.NewClient(ctx)
	if err != nil {
		return res, err
	}
	defer client.Close()

	req := &languagepb.AnalyzeEntitySentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: source,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}

	resp, err := client.AnalyzeEntitySentiment(ctx, req)

	if err != nil {
		return res, fmt.Errorf("error analyzing text: %s", err)
	}

	res.Raw = resp

	return res, nil
}

// AuthResult is the return from auth operations. It allows us to show
// tbe pure result and the work.
type AuthResult struct {
	Result bool                                       `json:"result"`
	Raw    *languagepb.AnalyzeEntitySentimentResponse `json:"raw"`
}

// AuthTerm does the check to see if the language query worked
func (l *AuthResult) AuthTerm(entityType string, positive bool) error {

	res := make(map[string]float32)

	for _, v := range l.Raw.Entities {
		res[strings.ToUpper(v.Type.String())] = v.Sentiment.Score
	}

	sent, ok := res[strings.ToUpper(entityType)]
	if !ok {
		return nil
	}

	if positive {
		if sent > 0 {
			l.Result = true
		}
		return nil
	}

	if sent < 0 {
		l.Result = true
	}

	return nil
}
