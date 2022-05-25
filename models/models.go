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

// Users is a