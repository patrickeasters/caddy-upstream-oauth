package upstream_oauth

import (
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// ServeHTTP implements caddyhttp.MiddlewareHandler
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Fetch token
	token, err := m.tokenSource.Token()
	if err != nil {
		m.logger.Sugar().Errorw("failed to get oauth token", "error", err)
	}

	// Insert token for use upstream
	token.SetAuthHeader(r)
	return next.ServeHTTP(w, r)
}
