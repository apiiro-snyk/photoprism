package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/get"
	"github.com/photoprism/photoprism/pkg/authn"
	"github.com/photoprism/photoprism/pkg/header"
	"github.com/photoprism/photoprism/pkg/i18n"
)

// OIDCLogin redirects a browser to the login page of the configured OpenID Connect provider, if any.
//
// GET /api/v1/oidc/login
func OIDCLogin(router *gin.RouterGroup) {
	router.GET("/oidc/login", func(c *gin.Context) {
		// Prevent CDNs from caching this endpoint.
		if header.IsCdn(c.Request) {
			AbortNotFound(c)
			return
		}

		// Disable caching of responses.
		c.Header(header.CacheControl, header.CacheControlNoStore)

		// Get client IP address for logs and rate limiting checks.
		clientIp := ClientIP(c)
		actor := "unknown client"
		action := "login"

		// Get global config.
		conf := get.Config()

		// Abort in public mode and if OIDC is disabled.
		if get.Config().Public() {
			event.AuditErr([]string{clientIp, "oidc", actor, action, authn.ErrDisabledInPublicMode.Error()})
			Abort(c, http.StatusForbidden, i18n.ErrForbidden)
			return
		} else if !conf.OIDCEnabled() {
			event.AuditErr([]string{clientIp, "oidc", actor, action, authn.ErrAuthenticationDisabled.Error()})
			Abort(c, http.StatusMethodNotAllowed, i18n.ErrUnsupported)
			return
		}

		// Get OIDC provider.
		provider := get.OIDC()

		if provider == nil {
			event.AuditErr([]string{clientIp, "oidc", actor, action, authn.ErrAuthenticationDisabled.Error()})
			Abort(c, http.StatusInternalServerError, i18n.ErrConnectionFailed)
			return
		}

		// Handle OIDC login request.
		provider.AuthCodeUrlHandler(c)
	})
}
