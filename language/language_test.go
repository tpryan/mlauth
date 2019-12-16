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
package language

import "testing"

func TestAuth(t *testing.T) {
	cases := []struct {
		entityType string
		text       string
		positive   bool
		want       bool
	}{
		{"location", "I love staying at Marriott hotels.", true, true},
		{"location", "I hate staying at Marriott.", false, true},
		{"location", "I love staying at Marriott hotels.", false, false},
		{"location", "I hate staying at Marriott.", true, false},
		{"organization", "Google is a great place to work. ", true, true},
		{"organization", "Google is a great place to work. ", false, false},
		{"CONSUMER_GOOD", "Diet Vanilla Coke is absolutely the best beverage ever created.", true, true},
	}

	for _, c := range cases {
		got, _ := Auth(c.entityType, c.text, c.positive)
		if got != c.want {
			t.Errorf("Auth('%s', '%s', '%t') got %t, want %t", c.entityType, c.text, c.positive, got, c.want)
		}
	}
}
