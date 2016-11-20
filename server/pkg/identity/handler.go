package identity

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	idapi "github.com/hooklift/apis/go/identity"
)

var (
	errUnauthorized = errors.New("Required Authorization. Token not found")
	accessTkCookie  = "hatk"
)

// Handler verifies authorization tokens against Hooklift Identity service.
func Handler(h http.Handler, clientID string) http.Handler {
	accounts := idapi.NewAccountsClient(Connection(clientID))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip gRPC requests if this handler gets misplaced.
		if strings.Contains(r.Header.Get("content-type"), "application/grpc") {
			h.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("authorization")
		var tokenValue string
		if token == "" {
			sessionCookie, err := r.Cookie(accessTkCookie)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

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
			ClientId: clientID,
			Token:    tokenValue,
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
