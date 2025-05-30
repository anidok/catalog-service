package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenSearchConfigShouldReturnDefaultConfig(t *testing.T) {
	os.Unsetenv("OPENSEARCH_HOST_SERVERS")
	c := Load()
	config := NewOpenSearchConfig(c)
	assert.EqualValues(t, []string{"http://localhost:9200"}, config.Host())
}

func TestOpenSearchConfigShouldReturnValidConfig(t *testing.T) {
	os.Setenv("OPENSEARCH_HOST_SERVERS", "localhost:9200,localhost:9201")
	c := Load()
	config := NewOpenSearchConfig(c)

	assert.EqualValues(t, []string{"localhost:9200", "localhost:9201"}, config.Host())
}
