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
func Auth(term, file string) (bool, error) {

	var err error
	req := &speechpb.RecognizeRequest{}

	if isValidURL(file) {

		if file[0:4] == "gs://" {
			return false, fmt.Errorf("speech api only works on Cloud Storage URI's")
		}

		req = getReqFromURI(file)
	} else {
		req, err = getReqFromFile(file)
		if err != nil {
			return false, err
		}
	}

	return compareAuth(req, term)
}

// AuthFromReader takes a reader containing an audio file and a term and
// compares them to each other to see if the an item matching the input term is
// contained in the audio file
func AuthFromReader(term string, file io.Reader) (bool, error) {

	var err error
	req := &speechpb.RecognizeRequest{}

	req, err = getReqFromReader(file)
	if err != nil {
		return false, err
	}

	return compareAuth(req, term)
}

func compareAuth(req *speechpb.RecognizeRequest, term string) (bool, error) {
	ann, err := findContent(req)

	if err != nil {
		return false, err
	}

	for _, v := range ann {
		if strings.Contains(strings.ToUpper(v), strings.ToUpper(term)) {
			return true, nil
		}
	}

	return false, nil
}

func findContent(req *speechpb.RecognizeRequest) ([]string, error) {
	var err error
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := client.Recognize(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("error recogniziong content file: %s", err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("unable to extract text from audio, usually wrong format")
	}

	res := []string{}
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			res = append(res, alt.Transcript)
		}
	}

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
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
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
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
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
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
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
