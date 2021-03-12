
package trade

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/middleware"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/VolticFroogo/Froogo/steam"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const (
	UserNotFound = iota
)

type response struct {
	ID         int
	CsrfSecret string
}

func Begin(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("handler/templates/panel/trade.html", "handler/templates/nested.html") // Parse the HTML pages
	if err != nil {
		helpers.ThrowErr(w, r, "Error template parsing", err)
		return
	}

	csrfSecret, err := r.Cookie("csrfSecret")
	if err != nil {
		helpers.ThrowErr(w, r, "Error reading cookie", err)
		return
	}

	inventory, err := steam.GetInventory(user.SteamID)
	if err != nil {
		helpers.ThrowErr(w, r, "Error Getting Steam inventory", err)
		return
	}

	vars := mux.Vars(r)
	classid := vars["Classid"]

	item, assetid := GetItem(classid, inventory)

	var wg sync.WaitGroup
	wg.Add(2)

	inspectable := false
	var floatAPI models.FloatAPI
	go func() {
		defer wg.Done()
