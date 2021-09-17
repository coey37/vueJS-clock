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

func generateRandomBytes(size int) ([]by