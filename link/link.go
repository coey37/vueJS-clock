package link

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/gorilla/context"
	openid "github.com/yohcop/openid-go"
)

var (
	nonceStore     = openid.NewSimpleNonceStore()
	discoveryCache = openid.NewSimpleDiscoveryCache()
)

func Steam(w http.ResponseWriter, r *http.Request) {
	if url, err := openid.RedirectURL("https://steamcommunity.com/openid", "https://froogo.co.uk/link/steam/callback", "https://froogo.co.uk/"); err == nil {
		http.Redirect(w, r, url, 303)
	} else {
		log.Print(err)
	}
}

func SteamCallback(w http.ResponseWriter, r *http.Request) {
	fullUrl := "https://froogo.co.uk" + r.URL.String()
	id, err := openid.Verify(fullUrl, discoveryCache, nonceStore)
	if err != nil {
		helpers.ThrowErr(w, r, "Error Verifying OpenID", err)
		return
	}

	steamIDString := strings.TrimPrefix(id, "https://steamcommunity.com/openid/id/")
	steamID, err := strconv.ParseInt(steamIDString, 10, 64)
	if err != nil {
		helpers.ThrowErr(w, r, "Error Converting SteamID to int64", err)
		return
	}

	uuidString := context.Get(r, "uuid").(string)

	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.ThrowErr(w, r, "Error converting string to int", err)
		return
	}

	err = db.LinkSteam(uuid, steamID)
	if err != nil {
		helpers.ThrowErr(w, r, "Error link