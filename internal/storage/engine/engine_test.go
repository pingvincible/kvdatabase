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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kvDatabase := engine.New()
			kvDatabase.Set(tc.key, tc.value)
			value := kvDatabase.Get(tc.key)
			assert.Equal(t, tc.value, value)
			kvDatabase.Delete(tc.key)
			value = kvDatabase.Get(tc.key)
			assert.Empty(t, value)
		})
	}
}
