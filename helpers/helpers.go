package helpers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Success bool `json:"success"`
}

func generateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	// Note that err == nil only if we read len(bytes) bytes.
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GenerateRandomString returns a random string with the size specified.
func GenerateRandomString(size int) (string, error) {
	b, err := generateRandomBytes(size)
	return base64.URLEncoding.EncodeToString(b), err
}

// HashPassword hashes a password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword checks a password against a hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ThrowErr throws an HTTP error and logs it to the server.
func ThrowErr(w http.ResponseWriter, r *http.Request, errName string, err error) {
	log.Printf("%v: %v\n", errName, err)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// JSO