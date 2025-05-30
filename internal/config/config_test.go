package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndGet(t *testing.T) {
	os.Setenv("SOME_ENV_KEY", "some-env-value")
	cfg := Load()
	assert.Equal(t, "catalog-service", AppName())
	assert.Equal(t, "some-value", cfg.GetOptionalValue("SOME_STR_KEY", "some-value"))
	assert.Equal(t, "some-env-value", cfg.GetOptionalValue("SOME_ENV_KEY", "default-value"))
	assert.Equal(t, 42, cfg.GetOptionalIntValue("SOME_INT_KEY", 0))
	assert.Equal(t, 0, cfg.GetOptionalIntValue("SOME_OTHER_INT_KEY", 0))
	os.Unsetenv("SOME_ENV_KEY")
}

func TestShouldGetOpenSearchConfig(t *testing.T) {
	c := Load()
	assert.Equal(t, NewOpenSearchConfig(c), OpenSearch())
}
