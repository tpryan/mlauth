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

// Package vision wraps the Cloud Vision API and provides an auth method that
// allows you to check an input image for an input term.
package vision

import (
	"context"
	"io"
	"net/url"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Auth takes a picture and a term and compares them to each other to see if
// the an item matching the input term is contained in the image
func Auth(term, file string) (bool, error) {

	var image *pb.Image
	var err error

	if isValidURL(file) {
		image = vision.NewImageFromURI(file)
	} else {
		image, err = imageFromFile(file)
		if err != nil {
			return false, err
		}
	}

	return compareAuth(image, term)
}

// AuthFromReader takes a picture and a term and compares them to each other
// to see if the an item matching the input term is contained in the image
func AuthFromReader(term string, file io.Reader) (bool, error) {

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return false, err
	}

	return compareAuth(image, term)
}

func compareAuth(image *pb.Image, term string) (bool, error) {
	ann, err := findLabels(image)

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

func findLabels(image *pb.Image) ([]string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}
	var labels []string
	for _, annotation := range annotations {
		labels = append(labels, annotation.Description)
	}
	return labels, nil
}

func imageFromFile(path string) (*pb.Image, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func isValidURL(toTest string) bool {

	r, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	if r.Scheme == "http" || r.Scheme == "https" || r.Scheme == "gs" {
		return true
	}
	return false

}
