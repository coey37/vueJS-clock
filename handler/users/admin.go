package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/middleware"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/gorilla/context"
)

type edit struct {
	ID, Privileges                            int
	CsrfSecret, Email, Password, Fname, Lname string
}

// Update is the handler for the update user request.
func Update(w http.ResponseWriter, r *http.Requ