package models

import (
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token lifetimes
const (
	// AuthTokenValidTime is the lifetime of an auth token.
	AuthTokenValidTime = time.Minute * 15
	// RefreshTokenValidTime is the lifetime of a refresh token.
	RefreshTokenValidTime = time.Hour * 72
)

// Privileges
const (
	PrivNone = iota
	PrivUser
	PrivAdmin
	PrivSuperAdmin
)

const (
	OfferStatusPending = iota
	OfferStatusAccepted
	OfferStatusCompleted
)

// User is a user retrieved from a Database.
type User struct {
	UUID, Priv, Points                      int
	Creation, SteamID                       int64
	Email, Password, Fname, Lname, Username string
}

// Users is an array of User for the admin page.
type Users []User

// TokenClaims are the claims in a token.
type TokenClaims struct {
	jwt.StandardClaims
	CSRF string `json:"csrf"`
}

// TemplateVariables is the struct used when executing a template.
type TemplateVariables struct {
	CsrfSecret, Assetid string
	User                User
	SteamUser           SteamUser
	UnixTime            int64
	Inventory           SteamInventory
	Item                Item
	FloatAPI            FloatAPI
	Inspectable         bool
	EstimatedPoints     int
	Offers              []Offer
}

// AJAXData is the struct used with the AJAX middleware.
type AJAXData struct {
	CsrfSecret string
}

// JTI is the struct used for JTIs in the DB.
type JTI struct {
	ID     int
	Expiry int64
	JTI    string
}

// ResponseWithID is a simple struct for responding to an AJAX request.
type ResponseWithID struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// ResponseWithIDInt is a simple struct for responding to an AJAX request.
type ResponseWithIDInt struct {
	Success bool `json:"success"`
	ID      int  `json:"id"`
}

type SteamUser struct {
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
}

type SteamInventory struct {
	Assets []struct {
		Appid      int    `json:"appid"`
		Contextid  string `json:"contextid"`
		Assetid    string `json:"assetid"`
		Classid    string `json:"classid"`
		Instanceid string `json:"instanceid"`
		Amount     string `json:"amount"`
	} `json:"assets"`
	Descriptions        []Item `json:"descriptions"`
	TotalInventoryCount int    `json:"total_inventory_count"`
	Success             int    `json:"success"`
	Rwgrsn              int    `json:"rwgrsn"`
}

type Item struct {
	Appid           int    `json:"appid"`
	Classid         string `json:"classid"`
	Instanceid      string `json:"instanceid"`
	Currency        int    `json:"currency"`
	BackgroundColor string `json:"background_color"`
	IconURL         string `json:"icon_url"`
	Descriptions    []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Color string `json:"color,omitempty"`
	} `json:"descriptions"`
	Tradable                  int    `json:"tradable"`
	Name                      string `json:"name"`
	NameColor                 string `json:"name_color"`
	Type                      string `json:"type"`
	MarketName                string `json:"market_name"`
	MarketHashName            string `json:"market_hash_name"`
	Commodity                 int    `json:"commodity"`
	MarketTradableRestriction int    `json:"market_tradable_restriction"`
	Marketable                int    `json:"marketable"`
	Tags                      []struct {
		Category              string `json:"category"`
		InternalName      