package steamweb

import (
	"errors"
	"fmt"
	"time"
)

var ErrConfigUndefinedParam = errors.New("config param is not defined")

const (
	DefaultSteamURL            = "https://api.steampowered.com"
	DefaultTimeout             = 10 * time.Second
	DefaultTLSHandshakeTimeout = 5 * time.Second
	DefaultDialerTimeout       = 5 * time.Second
	DefaultLimit               = 50000
)

type (
	Config struct {
		Disabled bool `json:"disabled" yaml:"disabled"`

		// Key is access api key for Steam requests.
		Key string `json:"key" yaml:"key"`

		// URL is a Steam Web API url string.
		URL string `json:"url" yaml:"url"`

		// Timeout specifies a time limit for requests made by this
		// Client. The timeout includes connection time, any
		// redirects, and reading the response body. The timer remains
		// running after Get, Head, Post, or Do return and will
		// interrupt reading of the Response.Body.
		//
		// The default is 10 seconds.
		Timeout time.Duration `json:"timeout" yaml:"timeout"`

		// Transport is a configuration settings for the implementation
		// of RoundTripper that supports Openapi, HTTPS, and Openapi proxies.
		Transport struct {
			// A Dialer contains options for connecting to an address.
			Dialer Dialer `json:"dialer" yaml:"dialer"`

			// TLSHandshakeTimeout specifies the maximum amount of time waiting to
			// wait for a TLS handshake. Zero means no timeout.
			TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`
		} `json:"transport" yaml:"transport"`

		Limit int `json:"limit" yaml:"limit"`

		DefaultServerNames []string `json:"default_server_names" yaml:"default_server_names"`
	}

	Dialer struct {
		// Timeout is the maximum amount of time a dial will wait for
		// a connect to complete. If Deadline is also set, it may fail
		// earlier.
		//
		// The default is 5 seconds.
		//
		// When using TCP and dialing a host name with multiple IP
		// addresses, the timeout may be divided between them.
		//
		// With or without a timeout, the operating system may impose
		// its own earlier timeout. For instance, TCP timeouts are
		// often around 3 minutes.
		Timeout time.Duration `json:"timeout" yaml:"timeout"`

		// Deadline is the absolute point in time after which dials
		// will fail. If Timeout is set, it may fail earlier.
		// Zero means no deadline, or dependent on the operating system
		// as with the Timeout option.
		Deadline time.Time `json:"deadline" yaml:"deadline"`

		// FallbackDelay specifies the length of time to wait before
		// spawning a fallback connection, when DualStack is enabled.
		// If zero, a default delay of 300ms is used.
		FallbackDelay time.Duration `json:"fallback_delay" yaml:"fallback_delay"`

		// KeepAlive specifies the keep-alive period for an active
		// network connection.
		// If zero, keep-alives are not enabled. Network protocols
		// that do not support keep-alives ignore this field.
		KeepAlive time.Duration `json:"keep_alive" yaml:"keep_alive"`
	}
)

func (cfg *Config) Validate() error {
	// Do not validate config for disabled client.
	if cfg.Disabled {
		return nil
	}

	if cfg.Key == "" {
		return fmt.Errorf("%w: %s", ErrConfigUndefinedParam, "key")
	}

	if cfg.URL == "" {
		return fmt.Errorf("%w: %s", ErrConfigUndefinedParam, "url")
	}

	return nil
}

func (cfg *Config) SetDefaults() {
	if cfg.URL == "" {
		cfg.URL = DefaultSteamURL
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	if cfg.Transport.TLSHandshakeTimeout == 0 {
		cfg.Transport.TLSHandshakeTimeout = DefaultTLSHandshakeTimeout
	}

	if cfg.Transport.Dialer.Timeout == 0 {
		cfg.Transport.Dialer.Timeout = DefaultDialerTimeout
	}

	if cfg.Limit == 0 {
		cfg.Limit = DefaultLimit
	}
}
