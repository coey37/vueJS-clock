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
		negroni.Wrap(http.HandlerFunc(logout)),
	))

	r.Handle("/panel", negroni.New(
		negroni.HandlerFunc(middleware.Panel),
		negroni.Wrap(http.HandlerFunc(panel)),
	))

	r.Handle("/panel/settings/update", http.HandlerFunc(users.Settings))

	r.Handle("/panel/user/new", http.HandlerFunc(users.New))
	r.Handle("/panel/user/update", http.HandlerFunc(users.Update))
	r.Handle("/panel/user/delete", http.HandlerFunc(users.Delete))

	r.Handle("/link/steam", http.HandlerFunc(link.Steam))
	r.Handle("/link/steam/callback", negroni.New(
		negroni.HandlerFunc(middleware.Panel),
		negroni.Wrap(http.HandlerFunc(link.SteamCallback)),
	))

	r.Handle("/panel/trade-offers", negroni.New(
		negroni.HandlerFunc(middleware.Panel),
		negroni.Wrap(http.HandlerFunc(trade.Offers)),
	))
	r.Handle("/panel/trade/{Classid}", negroni.New(
		negroni.HandlerFunc(middleware.Panel),
		negroni.Wrap(http.HandlerFunc(trade.Begin)),
	)).Methods(http.MethodGet)
	r.Handle("/panel/trade", http.HandlerFunc(trade.Offer)).Methods(http.MethodPost)
	r.Handle("/panel/trade/accept", http.HandlerFunc(trade.Accept))
	r.Handle("/panel/trade/cancel", http.HandlerFunc(trade.Cancel))

	r.Handle("/verify-email/{code}", http.HandlerFunc(users.VerifyEmail))
	r.Handle("/forgot-password", http.HandlerFunc(recovery.Begin)).Methods(http.MethodPost)
	r.Handle("/password-recovery", http.HandlerFunc(recovery.End)).Methods(http.MethodPost)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Printf("Server started...")
	http.ListenAndServe(":82", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("handler/templates/index.html", "handler/templates/nested.html") // Parse the HTML pages
	if err != nil {
		helpers.ThrowErr(w, r, "Template parsing error", err)
		return
	}

	variables := models.TemplateVariables{}
	err = t.Execute(w, variables) // Execute temmplate with variables
	if err != nil {
		helpers.ThrowErr(w, r, "Template execution error", err)
	}
}

func panel(w http.ResponseWriter, r *http.Request) {
	uuidString := context.Get(r, "uuid").(string)

	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.ThrowErr(w, r, "Error converting string to int", err)
		return
	}

	user, err := db.GetUserFromID(uuid)
	if err != nil {
		helpers.ThrowErr(w, r, "Error getting user from ID", err)
		return
	}

	execPanel(w, r, user, "panel")
}

func execPanel(w http.ResponseWriter, r *http.Request, user models.User, templateName string) {
	t, err := template.ParseFiles("handler/templates/panel/"+templateName+".html", "handler/templates/nested.html") // Parse the HTML pages
	if err != nil {
		helpers.ThrowErr(w, r, "Error template parsing", err)
		return
	}

	csrfSecret, err := r.Cookie("csrfSecret")
	if err != nil {
		helpers.ThrowErr(w, r, "Error reading cookie", err)
		return
	}

	var inventory models.SteamInventory
	if user.SteamID != 0 {
		// Get Steam inventory.
		inventory, err 