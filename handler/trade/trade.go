
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

		if len(item.Actions) != 0 {
			inspectable = true

			d := strings.TrimPrefix(item.Actions[0].Link, "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20S%owner_steamid%A%assetid%D")

			floatAPI, err = steam.FloatAPI(strconv.FormatInt(user.SteamID, 10), assetid, d)
			if err != nil {
				helpers.ThrowErr(w, r, "Error using float API", err)
				return
			}
		}
	}()

	var estimatedPoints int
	go func() {
		defer wg.Done()

		estimatedPoints, err = steam.GetItemValue(item.MarketHashName)
		if err != nil {
			helpers.ThrowErr(w, r, "Error getting item value", err)
			return
		}
	}()

	wg.Wait()

	variables := models.TemplateVariables{
		User:            user,
		CsrfSecret:      csrfSecret.Value,
		Inventory:       inventory,
		Item:            item,
		Assetid:         assetid,
		FloatAPI:        floatAPI,
		Inspectable:     inspectable,
		EstimatedPoints: estimatedPoints,
	}
	err = t.Execute(w, variables) // Execute temmplate with variables
	if err != nil {
		helpers.ThrowErr(w, r, "Template execution error", err)
	}
}

func Offer(w http.ResponseWriter, r *http.Request) {
	var data models.Offer                        // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		return
	}

	if !middleware.AJAX(w, r, models.AJAXData{CsrfSecret: data.CsrfSecret}) {
		// Failed middleware (invalid credentials)
		helpers.SuccessResponse(false, w, r)
		return
	}

	uuidString := context.Get(r, "uuid").(string)
	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error converting string to int", err)
		return
	}

	user, err := db.GetUserFromID(uuid)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting user from ID", err)
		return
	}

	data.Timestamp = time.Now().Unix()
	data.UserUUID = user.UUID

	receiver, err := db.GetUserFromUsername(data.User)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting user from username", err)
		return
	}

	if receiver.UUID == 0 {
		// If the user doesn't exist
		helpers.JSONResponse(models.ResponseWithIDInt{
			Success: false,
			ID:      UserNotFound,
		}, w)

		return
	}

	data.ReceiverUUID = receiver.UUID

	inventory, err := steam.GetInventory(user.SteamID)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error Getting Steam inventory", err)
		return
	}

	item, assetid := GetItem(data.Classid, inventory)
	data.Item = item

	d := strings.TrimPrefix(data.Item.Actions[0].Link, "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20S%owner_steamid%A%assetid%D")

	data.FloatAPI, err = steam.FloatAPI(strconv.FormatInt(user.SteamID, 10), assetid, d)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error using float API", err)
		return
	}

	err = db.AddOffer(data)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Adding offer error", err)
		return
	}

	helpers.SuccessResponse(true, w, r)
}

func Offers(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("handler/templates/panel/trade-offers.html", "handler/templates/nested.html") // Parse the HTML pages
	if err != nil {
		helpers.ThrowErr(w, r, "Error template parsing", err)
		return
	}

	csrfSecret, err := r.Cookie("csrfSecret")
	if err != nil {
		helpers.ThrowErr(w, r, "Error reading cookie", err)
		return
	}

	offers, err := db.GetOffers(user.UUID)
	if err != nil {
		helpers.ThrowErr(w, r, "Error getting trade offers", err)
		return
	}

	variables := models.TemplateVariables{
		User:       user,
		CsrfSecret: csrfSecret.Value,
		Offers:     offers,
	}
	err = t.Execute(w, variables) // Execute temmplate with variables
	if err != nil {
		helpers.ThrowErr(w, r, "Template execution error", err)
	}
}

func GetItem(classid string, inventory models.SteamInventory) (item models.Item, assetid string) {
	for _, steamItem := range inventory.Descriptions {
		if steamItem.Classid == classid {
			item = steamItem
			break
		}
	}

	for _, asset := range inventory.Assets {
		if asset.Classid == classid {
			assetid = asset.Assetid
			break
		}
	}

	return
}

func Cancel(w http.ResponseWriter, r *http.Request) {
	var data response                            // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		return
	}

	if !middleware.AJAX(w, r, models.AJAXData{CsrfSecret: data.CsrfSecret}) {
		// Failed middleware (invalid credentials)
		helpers.SuccessResponse(false, w, r)
		return
	}

	uuidString := context.Get(r, "uuid").(string)
	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error converting string to int", err)
		return
	}

	user, err := db.GetUserFromID(uuid)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting user from ID", err)
		return
	}

	offer, err := db.GetOffer(data.ID, user.UUID)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting offer", err)
		return
	}

	if offer.UserUUID != user.UUID && offer.ReceiverUUID != user.UUID {
		// This is exterior tampering
		helpers.SuccessResponse(false, w, r)
		return
	}

	err = db.DeleteOffer(data.ID)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error deleting offer", err)
		return
	}

	helpers.SuccessResponse(true, w, r)
}

func Accept(w http.ResponseWriter, r *http.Request) {
	var data response                            // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		return
	}

	if !middleware.AJAX(w, r, models.AJAXData{CsrfSecret: data.CsrfSecret}) {
		// Failed middleware (invalid credentials)
		helpers.SuccessResponse(false, w, r)
		return
	}

	uuidString := context.Get(r, "uuid").(string)
	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error converting string to int", err)
		return
	}

	user, err := db.GetUserFromID(uuid)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting user from ID", err)
		return
	}

	offer, err := db.GetOffer(data.ID, user.UUID)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting offer", err)
		return
	}

	if offer.ReceiverUUID != user.UUID {
		// This is exterior tampering
		helpers.SuccessResponse(false, w, r)
		return
	}

	err = db.SetOfferStatus(models.OfferStatusAccepted, data.ID)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error setting offer status", err)
		return
	}

	helpers.SuccessResponse(true, w, r)
}