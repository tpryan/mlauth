package speech

import (
	"os"
	"strings"
	"testing"
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
		if got != c.want {
			t.Errorf("Auth('%s', '%s') got %t, want %t", c.term, c.file, got, c.want)
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
		got, err := findContent(c.file)
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
	}

	for _, c := range cases {
		got := isValidURL(c.in)
		if got != c.want {
			t.Errorf("isValidURL('%s') got %t, want %t", c.in, got, c.want)
		}
	}
}
