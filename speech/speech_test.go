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
	"os"
	"strings"
	"testing"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

var gcs = "gs://" + os.Getenv("mlauth_bucket") + "/speech"
var web = "https://storage.googleapis.com/" + os.Getenv("mlauth_bucket") + "/speech"

func TestAuth(t *testing.T) {
	cases := []struct {
		term string
		file string
		want bool
	}{
		{"Brooklyn", "testdata/audio.raw", true},
		{"Walker", "testdata/audio.raw", false},
		{"quit", "testdata/quit.raw", true},
		{"Walker", "testdata/quit.raw", false},
		{"Brooklyn", gcs + "/audio.raw", true},
		{"Walker", gcs + "/audio.raw", false},
		{"quit", gcs + "/quit.raw", true},
		{"Walker", gcs + "/quit.raw", false},
		{"Brooklyn", web + "/audio.raw", false},
		{"Walker", web + "/audio.raw", false},
		{"quit", web + "/quit.raw", false},
		{"Walker", web + "/quit.raw", false},
		{"", "", false},
	}

	for _, c := range cases {
		got, _ := Auth(c.term, c.file)
		if got.Result != c.want {
			t.Errorf("Auth('%s', '%s') got %t, want %t", c.term, c.file, got.Result, c.want)
		}
	}
}

func TestFindContent(t *testing.T) {
	cases := []struct {
		term      string
		file      string
		shouldErr bool
	}{
		{"Brooklyn", "testdata/audio.raw", false},
		{"quit", "testdata/quit.raw", false},
		{"conference", "testdata/voicememo.m4a", true},

		{"Brooklyn", gcs + "/audio.raw", false},
		{"quit", gcs + "/quit.raw", false},
		{"conference", gcs + "/voicememo.m4a", true},

		{"Brooklyn", web + "/audio.raw", true},
		{"quit", web + "/quit.raw", true},
		{"conference", web + "/voicememo.m4a", true},

		{"", "", true},
	}

	for _, c := range cases {

		var req *speechpb.RecognizeRequest
		var err error

		if isValidURL(c.file) {

			if c.file[0:4] == "gs://" {
				if !c.shouldErr {
					t.Errorf("findLabels(%s) threw error: %s", c.file, err)
				}
			}

			req = getReqFromURI(c.file)
		} else {
			req, err = getReqFromFile(c.file)
			if err != nil {
				if !c.shouldErr {
					t.Errorf("findLabels(%s) threw error: %s", c.file, err)
				}
			}
		}

		got, err := findContent(req)
		if err != nil {
			if !c.shouldErr {
				t.Errorf("findLabels(%s) threw error: %s", c.file, err)
			}
			continue
		}

		found := false

		for _, result := range got.Raw.Results {
			for _, alt := range result.Alternatives {
				if strings.Contains(strings.ToUpper(alt.Transcript), strings.ToUpper(c.term)) {
					found = true
					break
				}
			}
		}

		if !found {
			t.Errorf("findLabels(%s) should have found: %s in %v", c.file, c.term, got)
		}

	}

}

func TestIsValidURL(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"fjcvj48fhr74hr8f", false},
		{"http://dwdwf.com", true},
		{"https://dwdwf.com", true},
		{"http://dwdwf", true},
		{"https://dwdwf", true},
	}

	for _, c := range cases {
		got := isValidURL(c.in)
		if got != c.want {
			t.Errorf("isValidURL('%s') got %t, want %t", c.in, got, c.want)
		}
	}
}
