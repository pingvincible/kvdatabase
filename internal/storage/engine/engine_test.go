package engine_test

import (
	"testing"

	"github.com/pingvincible/kvdatabase/internal/storage/engine"
	"github.com/stretchr/testify/assert"
)

func TestEngineMethods(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Test with key and value",
			key:   "key",
			value: "value",
		},
	}

	t.Parallel()

	for _, tc := range cases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			kvDatabase := engine.New()
			kvDatabase.Set(testCase.key, testCase.value)
			value := kvDatabase.Get(testCase.key)
			assert.Equal(t, testCase.value, value)
			kvDatabase.Delete(testCase.key)
			value = kvDatabase.Get(testCase.key)
			assert.Empty(t, value)
		})
	}
}
