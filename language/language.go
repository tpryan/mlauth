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

import (
	"context"
	"fmt"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func Auth(entityType, source string, positive bool) (bool, error) {

	ann, err := findSentiment(source)

	if err != nil {
		return false, err
	}

	termSent, ok := ann[strings.ToUpper(entityType)]
	if !ok {
		return false, nil
	}

	if positive {
		if termSent > 0 {
			return true, nil
		}
		return false, nil
	}

	if termSent < 0 {
		return true, nil
	}

	return false, nil
}

func findSentiment(source string) (map[string]float32, error) {
	var err error
	ctx := context.Background()

	client, err := language.NewClient(ctx)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("error analyzing text: %s", err)
	}

	res := make(map[string]float32)

	for _, v := range resp.Entities {
		res[strings.ToUpper(v.Type.String())] = v.Sentiment.Score
	}

	return res, nil
}
