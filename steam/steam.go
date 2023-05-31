
package steam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VolticFroogo/Froogo/models"
)

type UserResponse struct {
	Response struct {
		Players []struct {
			Steamid                  string `json:"steamid"`
			Communityvisibilitystate int    `json:"communityvisibilitystate"`
			Profilestate             int    `json:"profilestate"`
			Personaname              string `json:"personaname"`
			Lastlogoff               int    `json:"lastlogoff"`
			Commentpermission        int    `json:"commentpermission"`
			Profileurl               string `json:"profileurl"`
			Avatar                   string `json:"avatar"`
			Avatarmedium             string `json:"avatarmedium"`
			Avatarfull               string `json:"avatarfull"`
			Personastate             int    `json:"personastate"`
			Realname                 string `json:"realname"`
			Primaryclanid            string `json:"primaryclanid"`
			Timecreated              int    `json:"timecreated"`
			Personastateflags        int    `json:"personastateflags"`
			Loccountrycode           string `json:"loccountrycode"`
		} `json:"players"`
	} `json:"response"`
}

type ItemValue struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

var (
	steamAPIKey = os.Getenv("STEAM_API_KEY")
)

func GetUser(steamID int64) (user models.SteamUser, err error) {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%v&steamids=%v", steamAPIKey, steamID), nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "froogo")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var response UserResponse
	err = json.Unmarshal(body, &response)

	user = response.Response.Players[0]
	return
}

func GetInventory(steamID int64) (inventory models.SteamInventory, err error) {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs