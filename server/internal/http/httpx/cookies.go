package httpx

import "net/http"

// TODO: refactor this to use env vars
var (
	refreshTokenAge        = 60 * 60 * 24 * 30
	refreshTokenCookieName = "tokennamegng"
)

func GetRefreshTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
