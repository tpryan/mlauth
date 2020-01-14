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

// Package speech wraps the Cloud Speech API and provides an auth
// method that allows you to check an input audio file for an input term.
package speech

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

// Auth takes an audio file and a term and compares them to each other to see if
// the an item matching the input term is contained in the audio file
func Auth(term string, rate, channels int32, file string) (AuthResult, error) {

	var err error
	req := &speechpb.RecognizeRequest{}

	if isValidURL(file) {

		if file[0:4] == "gs://" {
			return AuthResult{}, fmt.Errorf("speech api only works on Cloud Storage URI's")
		}

		req = getReqFromURI(file)
	} else {
		req, err = getReqFromFile(file)
		if err != nil {
			return AuthResult{}, err
		}
	}

	req.Config.SampleRateHertz = rate
	req.Config.AudioChannelCount = channels

	return compareAuth(req, term)
}

// AuthFromReader takes a reader containing an audio file and a term and
// compares them to each other to see if the an item matching the input term is
// contained in the audio file
func AuthFromReader(term string, rate, channels int32, file io.Reader) (AuthResult, error) {

	var err error
	req := &speechpb.RecognizeRequest{}

	req, err = getReqFromReader(file)
	if err != nil {
		return AuthResult{}, err
	}

	req.Config.SampleRateHertz = rate
	req.Config.AudioChannelCount = channels

	return compareAuth(req, term)
}

func compareAuth(req *speechpb.RecognizeRequest, term string) (AuthResult, error) {
	resp, err := findContent(req)

	if err != nil {
		return resp, err
	}

	resp.AuthTerm(term)

	return resp, nil
}

func findContent(req *speechpb.RecognizeRequest) (AuthResult, error) {
	var err error
	ctx := context.Background()

	res := AuthResult{}

	client, err := speech.NewClient(ctx)
	if err != nil {
		return res, err
	}
	defer client.Close()

	resp, err := client.Recognize(ctx, req)

	if err != nil {
		return res, fmt.Errorf("error recogniziong content file: %s", err)
	}

	if len(resp.Results) == 0 {
		res.Raw = resp
		return res, fmt.Errorf("unable to extract text from audio, no words or wrong format")
	}

	res.Raw = resp
	return res, nil
}

func getReqFromFile(file string) (*speechpb.RecognizeRequest, error) {
	req := &speechpb.RecognizeRequest{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return req, fmt.Errorf("error reading file: %s", err)
	}

	req = &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}
	return req, nil
}

func getReqFromReader(file io.Reader) (*speechpb.RecognizeRequest, error) {

	req := &speechpb.RecognizeRequest{}

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	data := buf.Bytes()

	req = &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}
	return req, nil
}

func getReqFromURI(path string) *speechpb.RecognizeRequest {

	return &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: path},
		},
	}
}

func isValidURL(toTest string) bool {

	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true
}

// AuthResult is the return from auth operations. It allows us to show
// tbe pure result and the work.
type AuthResult struct {
	Result bool                        `json:"result"`
	Raw    *speechpb.RecognizeResponse `json:"raw"`
}

// AuthTerm does the check to see if the language query worked
func (l *AuthResult) AuthTerm(term string) error {

	res := []string{}
	for _, result := range l.Raw.Results {
		for _, alt := range result.Alternatives {
			res = append(res, alt.Transcript)
		}
	}

	for _, v := range res {
		if strings.Contains(strings.ToUpper(v), strings.ToUpper(term)) {
			l.Result = true
		}
	}

	return nil
}
