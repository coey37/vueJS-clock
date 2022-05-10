
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
