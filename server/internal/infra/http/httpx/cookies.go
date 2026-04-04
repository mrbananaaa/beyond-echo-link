package httpx

import (
	"net/http"
)

// TODO: refactor this to use env vars
var (
	refreshTokenCookieName = "auth_refresh_token"
)

func SetRefreshTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func GetRefreshTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func ClearRefreshTokenCokie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
