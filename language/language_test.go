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
	}

	for _, c := range cases {
		got, _ := Auth(c.entityType, c.text, c.positive)
		if got != c.want {
			t.Errorf("Auth('%s', '%s', '%t') got %t, want %t", c.entityType, c.text, c.positive, got, c.want)
		}
	}
}
