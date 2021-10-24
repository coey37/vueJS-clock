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

// JSONResponse sends a client a JSON response.
func JSONResponse(data interface{}, w http.ResponseWriter) (err error) {
	dataJSON, err := json.Marshal(data) // Encode response into JSON.
	if err != nil {
		return
	}
	w.Write(dataJSON) // Write JSON data to response writer.
	return
}

// SuccessResponse is a JSON response with a success boolean.
func SuccessResponse(valid bool, w http.ResponseWriter, r *http.Request) {
	res := response{
		Success: valid,
	}
	resEnc, err := json.Marshal(res) // Encode response into JSON.
	if err != nil {
		ThrowErr(w, r, "Sending success response error: %v", err)
	}
	w.Write(resEnc) // Write JSON data to response writer.
	return
}

// CheckEmail checks if an email is valid.
func CheckEmail(email string) (err error) {
	err = checkmail.ValidateFormat(email)
	if err != nil {
		return
	}

	err = checkmail.ValidateHost(email)
	return
}
