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

ty