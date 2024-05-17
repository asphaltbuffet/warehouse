package configurator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/wherehouse/pkg/configurator"
)

func TestConfigKey_String(t *testing.T) {
	tests := []struct {
		name string
		key  configurator.ConfigKey
		want string
	}{
		{"one word", configurator.ConfigDirKey, "config-dir"},
		{"nested", configurator.DatabaseFilenameKey, "database.file"},
		{"unknown", "fake", "fake"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.key.String())
		})
	}
}
