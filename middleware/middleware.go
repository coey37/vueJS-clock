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

	if authTokenString.Value != "" 