
package myJWT

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/dgrijalva/jwt-go"
)

// Variables/
var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

// InitKeys defines the signing and verification RSA keys for JWT.
func InitKeys() error {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

// DeleteJTI deletes a JTI when given a refresh token.
func DeleteJTI(tokenString string) (err error) {
	token, _ := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	tokenClaims, _ := token.Claims.(*models.TokenClaims)
	err = db.DeleteJTI(tokenClaims.StandardClaims.Id)
	return
}

/*
	Refreshing tokens and all related functions.
*/

// RefreshTokens returns new fresh tokens with a CSRF Secret.
func RefreshTokens(oldRefreshTokenString string) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {
	token, err := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})
	if err != nil {
		return
	}

	oldTokenClaims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	return CreateNewTokens(oldTokenClaims.StandardClaims.Subject)
}

/*
	Validating tokens and all related functions.
*/

// CheckToken checks the validity of a token.
func CheckToken(tokenString, csrfSecret string, refresh, checkCsrf bool) (valid bool, uuid string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	tokenClaims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	if csrfSecret != tokenClaims.CSRF && checkCsrf {
		return false, "", fmt.Errorf("csrf token doesn't match jwt")
	}

	if refresh {
		jti, err := db.GetJTI(tokenClaims.StandardClaims.Id)
		if err != nil {
			return false, "", fmt.Errorf("getting jti error")
		}

		jtiValid, err := db.CheckJTI(jti)
		if err != nil {
			return false, "", fmt.Errorf("checking jti error")
		}

		if jtiValid {
			err = db.DeleteJTI(tokenClaims.StandardClaims.Id) // There will be a new JTI created in it's place by the middleware.
			if err != nil {
				return true, tokenClaims.StandardClaims.Subject, err
			}

			return true, tokenClaims.StandardClaims.Subject, nil
		}
	}

	return token.Valid, tokenClaims.StandardClaims.Subject, nil
}

/*
	Creating tokens and all related functions.
*/

// CreateNewTokens creates an auth and refresh token.
func CreateNewTokens(uuid string) (authTokenString, refreshTokenString, csrfSecret string, err error) {
	// Generate the CSRF Secret
	csrfSecret, err = generateCSRFSecret()
	if err != nil {
		return
	}

	// Generate the refresh token
	refreshTokenString, err = createRefreshTokenString(uuid, csrfSecret)
	if err != nil {
		return
	}

	// Generate the auth token
	authTokenString, err = createAuthTokenString(uuid, csrfSecret)

	return
}

func createRefreshTokenString(uuid, csrfSecret string) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
	refreshJti, err := db.StoreRefreshToken()
	if err != nil {
		return
	}

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        refreshJti.JTI,  // Token Id
			Subject:   uuid,            // Universally Unique Identifier
			ExpiresAt: refreshTokenExp, // Expiry time in UNIX
		},
		csrfSecret, // CSRF Secret to prevent CSRF
	}

	// Make a new unsigned token
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	// Sign token
	refreshTokenString, err = unsignedToken.SignedString(signKey)

	return
}

func createAuthTokenString(uuid, csrfSecret string) (authTokenString string, err error) {
	authTokenExp := time.Now().Add(models.AuthTokenValidTime).Unix()

	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: authTokenExp,
		},
		csrfSecret,
	}

	// Make a new unsigned token
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodRS256, authClaims)
	// Sign token
	authTokenString, err = unsignedToken.SignedString(signKey)

	return
}

func generateCSRFSecret() (csrfSecret string, err error) {
	return helpers.GenerateRandomString(32)
}