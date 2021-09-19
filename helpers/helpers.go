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

/