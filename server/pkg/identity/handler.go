package identity

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	idapi "github.com/hooklift/apis/go/identity"
)

var (
	errUnauthorized = errors.New("Authorization required. Access token not found")
	accessTkCookie  = "hatk"
)

// Handler verifies authorization tokens against Hooklift Identity service.
func Handler(h http.Handler, clientURI string) http.Handler {
	accounts := idapi.NewAccountsClient(Connection(clientURI))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip gRPC requests if this handler gets misplaced.
		if strings.Contains(r.Header.Get("content-type"), "application/grpc") {
			h.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("authorization")
		var tokenValue string
		if token == "" {
			// If we didn't find an authorization header, we check the http-only cookie "hatk" for a valid access token.
			sessionCookie, err := r.Cookie(accessTkCookie)
			// Lets specific endpoint handler to reject if authorization is required.
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			// Since we have a token within a cookie, it means the user signed in from the browser. So
			// we need to take care of refreshing the access token when it is expired.
			// TODO(c4milo): Verify if token expired and refresh it
			// TODO(c4milo): Set new tokens in an http-only cookie.
			tokenValue = sessionCookie.Value
		} else {
			tokenParts := strings.Fields(token)
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, errUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			tokenValue = tokenParts[1]
		}

		if tokenValue == "" {
			h.ServeHTTP(w, r)
			return
		}

		res, err := accounts.VerifyToken(r.Context(), &idapi.VerifyTokenRequest{
			ClientUri: clientURI,
			Token:     tokenValue,
		})
		if err != nil {
			glog.Errorf("%+v", errors.Wrapf(err, "failed verifying token"))
			http.Error(w, errUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		ctx := NewContext(r.Context(), res.Token)
		req := r.WithContext(ctx)

		h.ServeHTTP(w, req)
	})
}
