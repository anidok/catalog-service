package config

import "strings"

type OpenSearchConfig struct {
	host []string
}

func NewOpenSearchConfig(cfg *AppConfig) *OpenSearchConfig {
	return &OpenSearchConfig{
		host: strings.Split(cfg.GetValue("OPENSEARCH_HOST_SERVERS"), ","),
	}
}

func (c *OpenSearchConfig) Host() []string {
	return c.host
}
