package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64ToShortUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		seed     int64
		expected func(assert.TestingT, string)
	}{
		{
			name: "smallest lexicographically",
			seed: 0,
			expected: func(t assert.TestingT, translatedUrl string) {
				assert.Equal(t, strings.Repeat(string(characterSet[0]), tokenLength), translatedUrl)
			},
		},
		{
			name: "greatest lexicographically",
			seed: binPow(base, tokenLength) - 1,
			expected: func(t assert.TestingT, translatedUrl string) {
				assert.Equal(t, strings.Repeat(string(characterSet[base-1]), tokenLength), translatedUrl)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.expected(t, int64ToShortUrl(tc.seed))
		})
	}
}
