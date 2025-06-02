package config

import (
	"strings"
	"time"
)

type OpenSearchConfig struct {
	hosts               []string
	maxIdleConns        int
	maxIdleConnsPerHost int
	idleConnTimeout     time.Duration
	dialTimeout         time.Duration
	keepAlive           time.Duration
	tlsHandshakeTimeout time.Duration
}

func NewOpenSearchConfig(cfg *AppConfig) *OpenSearchConfig {
	return &OpenSearchConfig{
		hosts:               strings.Split(cfg.GetValue("OPENSEARCH_HOST_SERVERS"), ","),
		maxIdleConns:        cfg.GetOptionalIntValue("OPENSEARCH_MAX_IDLE_CONNS", 100),
		maxIdleConnsPerHost: cfg.GetOptionalIntValue("OPENSEARCH_MAX_IDLE_CONNS_PER_HOST", 100),
		idleConnTimeout:     time.Duration(cfg.GetOptionalIntValue("OPENSEARCH_IDLE_CONN_TIMEOUT_MS", 90000)) * time.Millisecond,
		dialTimeout:         time.Duration(cfg.GetOptionalIntValue("OPENSEARCH_DIAL_TIMEOUT_MS", 30000)) * time.Millisecond,
		keepAlive:           time.Duration(cfg.GetOptionalIntValue("OPENSEARCH_KEEP_ALIVE_MS", 30000)) * time.Millisecond,
		tlsHandshakeTimeout: time.Duration(cfg.GetOptionalIntValue("OPENSEARCH_TLS_HANDSHAKE_TIMEOUT_MS", 10000)) * time.Millisecond,
	}
}

func (c *OpenSearchConfig) Host() []string {
	return c.hosts
}
func (c *OpenSearchConfig) MaxIdleConns() int {
	return c.maxIdleConns
}
func (c *OpenSearchConfig) MaxIdleConnsPerHost() int {
	return c.maxIdleConnsPerHost
}
func (c *OpenSearchConfig) IdleConnTimeout() time.Duration {
	return c.idleConnTimeout
}
func (c *OpenSearchConfig) DialTimeout() time.Duration {
	return c.dialTimeout
}
func (c *OpenSearchConfig) KeepAlive() time.Duration {
	return c.keepAlive
}
func (c *OpenSearchConfig) TLSHandshakeTimeout() time.Duration {
	return c.tlsHandshakeTimeout
}
