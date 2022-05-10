package middleware

import (
	"net/http"
	"time"

	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/middleware/myJWT"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/gorilla/context"
)

// Panel handles authentication for authenticated pages.
func Panel(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authTokenString, err := r.Cookie("authToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	refreshTokenString, err := r.Cookie("refreshToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	if authTokenString.Value != "" {
		authTokenValid, uuid, err := myJWT.CheckToken(authTokenString.Value, "", false, false)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if authTokenValid {
			context.Set(r, "uuid", uuid)
			next(w, r)
			return
		}
	}

	if refreshTokenString.Value != "" {
		refreshTokenValid, uuid, err := myJWT.CheckToken(refreshTokenString.Value, "", true, false)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if refreshTokenValid {
			newAuthTokenString, newRefreshTokenString, newCsrfSecret, err := myJWT.RefreshTokens(refreshTokenString.Value)
			if err != nil {
				helpers.ThrowErr(w, r, "Creating new tokens error", err)
				return
			}

			WriteNewAuth(w, r, newAuthTokenString, newRefreshTokenString, newCsrfSecret)

			context.Set(r, "uuid", uuid)
			next(w, r)
			return
		}
	}

	RedirectToLogin(w, r)
}

// Form is the function used to protect forms.
func Form(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authTokenString, err := r.Cookie("authToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	refreshTokenString, err := r.Cookie("refreshToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	csrfSecret := r.FormValue("csrfSecret")

	if authTokenString.Value != "" {
		authTokenValid, uuid, err := myJWT.CheckToken(authTokenString.Value, csrfSecret, false, true)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if authTokenValid {
			context.Set(r, "uuid", uuid)
			next(w, r)
			return
		}
	}

	if refreshTokenString.Value != "" {
		refreshTokenValid, uuid, err := myJWT.CheckToken(refreshTokenString.Value, csrfSecret, true, true)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if refreshTokenValid {
			newAuthTokenString, newRefreshTokenString, newCsrfSecret, err := myJWT.RefreshTokens(refreshTokenString.Value)
			if err != nil {
				helpers.ThrowErr(w, r, "Creating new tokens error", err)
				return
			}

			WriteNewAuth(w, r, newAuthTokenString, newRefreshTokenString, newCsrfSecret)

			context.Set(r, "uuid", uuid)
			next(w, r)
			return
		}
	}

	RedirectToLogin(w, r)
}

// AJAX is the function used to protect AJAX requests.
func AJAX(w http.ResponseWriter, r *http.Request, data models.AJAXData) (valid bool) {
	valid = false

	authTokenString, err := r.Cookie("authToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	refreshTokenString, err := r.Cookie("refreshToken")
	if err != nil {
		helpers.ThrowErr(w, r, "Reading cookie error", err)
		return
	}

	if authTokenString.Value != "" {
		authTokenValid, uuid, err := myJWT.CheckToken(authTokenString.Value, data.CsrfSecret, false, true)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if authTokenValid {
			context.Set(r, "uuid", uuid)
			return true
		}
	}

	if refreshTokenString.Value != "" {
		refreshTokenValid, uuid, err := myJWT.CheckToken(refreshTokenString.Value, data.CsrfSecret, true, true)
		if err != nil {
			helpers.ThrowErr(w, r, "Checking token error", err)
			return
		}

		if refreshTokenValid {
			newAuthTokenString, newRefreshTokenString, newCsrfSecret, err := myJWT.RefreshTokens(refreshTokenString.Value)
			if err != nil {
				helpers.ThrowErr(w, r, "Creating new tokens error", err)
				return
			}

			WriteNewAuth(w, r, newAuthTokenString, newRefreshTokenString, newCsrfSecret)

			context.Set(r, "uuid", uuid)
			return true
		}
	}

	return
}

// WriteNewAuth writes authentication to a user's browser.
func WriteNewAuth(w http.ResponseWriter, r *http.Request, authTokenString, refreshTokenString, csrfSecret string) {
	expiration := time.Now().Add(time.Hour * 24 * 365)

	cookie := http.Cookie{Name: "authToken", Value: authTokenString, Expires: expiration, Path: "/", HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{Name: "refreshToken", Value: refreshTokenString, Expires: expiration, Path: "/", HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{Name: "csrfSecret", Value: csrfSecret, Expires: expiration, Path: "/", HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)

	return
}

// RedirectToLogin redirects the client to the login.
func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
