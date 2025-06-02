package config

import "strings"

type OpenSearchConfig struct {
	host    []string
	timeout int
}

func NewOpenSearchConfig(cfg *AppConfig) *OpenSearchConfig {
	return &OpenSearchConfig{
		host:    strings.Split(cfg.GetValue("OPENSEARCH_HOST_SERVERS"), ","),
		timeout: cfg.GetOptionalIntValue("OPENSEARCH_TIMEOUT_MS", 15000),
	}
}

func (c *OpenSearchConfig) Host() []string {
	return c.host
}

func (c *OpenSearchConfig) Timeout() int {
	return c.timeout
}
