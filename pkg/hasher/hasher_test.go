package hasher

import (
	"testing"
)

func TestHashString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"Go is awesome!", "d557c06d48fd26fa66dfc2c327288fe815f537addfde447da9e70ae69ceae437"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := HashString(test.input)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
