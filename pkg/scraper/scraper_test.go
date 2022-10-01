package scraper

import "testing"

func TestNormalizePrice(t *testing.T) {
	tests := []struct {
		price string
		want  string
	}{
		{"$1,000", "1000"},
		{"$1,000,000", "1000000"},
	}

	for _, test := range tests {
		if got := normalizePrice(test.price); got != test.want {
			t.Errorf("normalizePrice(%q) = %q; want %q", test.price, got, test.want)
		}
	}
}
