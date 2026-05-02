package fbd

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type sessionCheckTransport struct {
	base     http.RoundTripper
	jar      http.CookieJar
	baseURL  *url.URL
	enforcer func() bool
}

func (t *sessionCheckTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.enforcer != nil && t.enforcer() {
		needsSession := req != nil &&
			req.URL != nil &&
			t.baseURL != nil &&
			req.URL.Host == t.baseURL.Host &&
			!strings.HasPrefix(req.URL.Path, "/signin")

		hasSession := false
		if needsSession && t.jar != nil {
			for _, cookie := range t.jar.Cookies(t.baseURL) {
				if cookie.Name == sessionCookieName && cookie.Value != "" {
					hasSession = true
					break
				}
			}
		}

		if needsSession && !hasSession {
			return nil, errors.New("missing ASP.NET_SessionId cookie for authenticated request")
		}
	}
	return t.base.RoundTrip(req)
}
