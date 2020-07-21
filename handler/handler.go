package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/handler/recovery"
	"github.com/VolticFroogo/Froogo/handler/trade"
	"github.com/VolticFroogo/Froogo/handler/users"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/link"
	"github.com/VolticFroogo/Froogo/middleware"
	"github.com/VolticFroogo/Froogo/middleware/myJWT"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/VolticFroogo/Froogo/steam"
	"github.com/go-recaptcha/recaptcha"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	captchaSecret = os.Getenv("CAPTCHA_SECRET")
	captcha       = recaptcha.New(captchaSecret)
)

type loginData struct {
	Email, Password, Captcha string
}

// Start the server by handling the web server.
func Start() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Handle("/", http.HandlerFunc(index))

	r.Handle("/login", http.HandlerFunc(login)).Methods(http.MethodPost)

	r.Handle("/logout", negroni.New(
		negroni.HandlerFunc(middleware.Form),
		negroni.Wrap(http.H