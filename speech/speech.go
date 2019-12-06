package speech

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

// Auth takes an audio file and a term and compares them to each other to see if
// the an item matching the input term is contained in the audio file
func Auth(term, file string) (bool, error) {

	ann, err := findContent(file)

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

func findContent(file string) ([]string, error) {
	var err error
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &speechpb.RecognizeRequest{}

	if isValidURL(file) {
		req = getReqFromURI(file)
	} else {
		req, err = getReqFromFile(file)
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.Recognize(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("error recogniziong content file: %s", err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("unable to extract text from audio, usually wrong format: %s", err)
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
