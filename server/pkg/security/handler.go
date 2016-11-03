package security

import "net/http"

// Handler checks whether an authorization JWT token was sent, verifies it against
// the authorization server and decodes it so that it can be used by other handlers
// to verify authorized scopes.
func Handler(h http.Handler) http.Handler {
	// TODO(c4milo): implement this
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
