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
