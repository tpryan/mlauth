package vision

import (
	"os"
	"strings"
	"testing"
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
		term string
		file string
	}{
		{"Golden Retriever", "testdata/dog.jpg"},
		{"dishware", "testdata/plate.jpg"},
		{"floral", "testdata/decoration.jpg"},
		{"fork", "testdata/fork.jpg"},
		{"soda", "testdata/soda.jpg"},
		{"can", "testdata/soda.jpg"},
		{"Coca-cola", "testdata/soda.jpg"},
		{"Golden Retriever", gcs + "/dog.jpg"},
		{"dishware", gcs + "/plate.jpg"},
		{"floral", gcs + "/decoration.jpg"},
		{"fork", gcs + "/fork.jpg"},
		{"soda", gcs + "/soda.jpg"},
		{"can", gcs + "/soda.jpg"},
		{"Coca-cola", gcs + "/soda.jpg"},
		{"Golden Retriever", web + "/dog.jpg"},
		{"dishware", web + "/plate.jpg"},
		{"floral", web + "/decoration.jpg"},
		{"fork", web + "/fork.jpg"},
		{"soda", web + "/soda.jpg"},
		{"can", web + "/soda.jpg"},
		{"Coca-cola", web + "/soda.jpg"},
	}

	for _, c := range cases {
		got, err := findLabels(c.file)
		if err != nil && c.file != "" {
			t.Errorf("findLabels(%s) threw error: %s", c.file, err)
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
