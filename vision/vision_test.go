package vision

import (
	"os"
	"strings"
	"testing"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var gcs = "gs://" + os.Getenv("mlauth_bucket") + "/vision"
var web = "https://storage.googleapis.com/" + os.Getenv("mlauth_bucket") + "/vision"

func TestAuth(t *testing.T) {
	cases := []struct {
		term string
		file string
		want bool
	}{
		{"", "", false},
		{"Golden Retriever", "testdata/dog.jpg", true},
		{"cat", "testdata/dog.jpg", false},
		{"soda", "testdata/soda.jpg", true},
		{"soda", "testdata/soda.jpg", true},
		{"can", "testdata/soda.jpg", true},
		{"coca-cola", "testdata/soda.jpg", true},
		{"pepsi", "testdata/soda.jpg", false},
		{"dishware", "testdata/plate.jpg", true},
		{"floral", "testdata/decoration.jpg", true},
		{"fork", "testdata/fork.jpg", true},

		{"Golden Retriever", gcs + "/dog.jpg", true},
		{"cat", gcs + "/dog.jpg", false},
		{"soda", gcs + "/soda.jpg", true},
		{"soda", gcs + "/soda.jpg", true},
		{"can", gcs + "/soda.jpg", true},
		{"coca-cola", gcs + "/soda.jpg", true},
		{"pepsi", gcs + "/soda.jpg", false},
		{"dishware", gcs + "/plate.jpg", true},
		{"floral", gcs + "/decoration.jpg", true},
		{"fork", gcs + "/fork.jpg", true},

		{"Golden Retriever", web + "/dog.jpg", true},
		{"cat", web + "/dog.jpg", false},
		{"soda", web + "/soda.jpg", true},
		{"soda", web + "/soda.jpg", true},
		{"can", web + "/soda.jpg", true},
		{"coca-cola", web + "/soda.jpg", true},
		{"pepsi", web + "/soda.jpg", false},
		{"dishware", web + "/plate.jpg", true},
		{"floral", web + "/decoration.jpg", true},
		{"fork", web + "/fork.jpg", true},
	}

	for _, c := range cases {
		got, _ := Auth(c.term, c.file)
		if got != c.want {
			t.Errorf("Auth('%s', '%s') got %t, want %t", c.term, c.file, got, c.want)
		}
	}
}

func TestFindLabels(t *testing.T) {
	cases := []struct {
		term      string
		file      string
		shouldErr bool
	}{
		{"", "", true},
		{"", "testdata/blank.txt", true},
		{"Golden Retriever", "testdata/dog.jpg", false},
		{"dishware", "testdata/plate.jpg", false},
		{"floral", "testdata/decoration.jpg", false},
		{"fork", "testdata/fork.jpg", false},
		{"soda", "testdata/soda.jpg", false},
		{"can", "testdata/soda.jpg", false},
		{"Coca-cola", "testdata/soda.jpg", false},
		{"Golden Retriever", gcs + "/dog.jpg", false},
		{"dishware", gcs + "/plate.jpg", false},
		{"floral", gcs + "/decoration.jpg", false},
		{"fork", gcs + "/fork.jpg", false},
		{"soda", gcs + "/soda.jpg", false},
		{"can", gcs + "/soda.jpg", false},
		{"Coca-cola", gcs + "/soda.jpg", false},
		{"Golden Retriever", web + "/dog.jpg", false},
		{"dishware", web + "/plate.jpg", false},
		{"floral", web + "/decoration.jpg", false},
		{"fork", web + "/fork.jpg", false},
		{"soda", web + "/soda.jpg", false},
		{"can", web + "/soda.jpg", false},
		{"Coca-cola", web + "/soda.jpg", false},
	}

	for _, c := range cases {

		var image *pb.Image
		var err error

		if isValidURL(c.file) {
			image = vision.NewImageFromURI(c.file)
		} else {
			image, err = imageFromFile(c.file)
			if err != nil {
				if !c.shouldErr {
					t.Errorf("findLabels(%s) threw error: %s", c.file, err)
				}
			}
		}

		got, err := findLabels(image)
		if err != nil {
			if !c.shouldErr {
				t.Errorf("findLabels(%s) threw error: %s", c.file, err)
			}
			continue
		}

		found := false
		for _, r := range got {
			if strings.Contains(strings.ToUpper(r), strings.ToUpper(c.term)) {
				found = true
				break
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
		{"/Users/username/file", false},
		{"gs://bucket/file", true},
	}

	for _, c := range cases {
		got := isValidURL(c.in)
		if got != c.want {
			t.Errorf("isValidURL('%s') got %t, want %t", c.in, got, c.want)
		}
	}
}
