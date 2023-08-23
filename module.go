package upstream_oauth

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	caddy.RegisterModule(Middleware{})
}

// Middleware implements an HTTP handler that manages OAuth client credentials
// and injects an access token in the request upstream
type Middleware struct {
	TokenURL     string   `json:"token_url"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`

	tokenSource oauth2.TokenSource
	logger      *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.upstream_oauth",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	if len(m.TokenURL) == 0 {
		return fmt.Errorf("token url must be provided")
	}
	if len(m.ClientID) == 0 || len(m.ClientSecret) == 0 {
		return fmt.Errorf("client id and secret must be provided")
	}

	conf := clientcredentials.Config{
		TokenURL:     m.TokenURL,
		ClientID:     m.ClientID,
		ClientSecret: m.ClientSecret,
		Scopes:       m.Scopes,
	}

	m.tokenSource = conf.TokenSource(ctx.Context)
	m.logger = ctx.Logger(m)

	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	if m.tokenSource == nil {
		return fmt.Errorf("no token source")
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)
