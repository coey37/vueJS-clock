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
func Update(w http.ResponseWriter, r *http.Request) {
	var data edit                                // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		return
	}

	if !middleware.AJAX(w, r, models.AJAXData{CsrfSecret: data.CsrfSecret}) {
		// Failed middleware (invalid credentials)
		helpers.SuccessResponse(false, w, r)
		return
	}

	uuidString := context.Get(r, "uuid").(string)
	uuid, err := strconv.Atoi(uuidString)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error converting string to int", err)
	}

	user, err := db.GetUserFromID(uuid)
	if err != nil {
		helpers.SuccessResponse(false, w, r)
		helpers.ThrowErr(w, r, "Error getting user from ID", err)
	}

	if user.Priv != models.PrivSuperAdmin {
		// User isn't a super admin.
		helpers.SuccessResponse(false, w, r)
		return
	}

	if data.Password == "" {
		err = db.EditUserNoPassword(data.ID, data.Email, data.Fname, data.Lname, data.Privileges)
		if err != nil {
			helpers.SuccessResponse(false, w, r)
			helpers.ThrowErr(w, r, "Editing user (no password) error", err)
			return
		}
	} else {
		password, err := helpers.HashPassword(data.Password)
		if err != nil {
			helpers.SuccessResponse(false, w, r)
			helpers.ThrowErr(w, r, "Hashing password error", err)
			return
		}

		err = db.EditUser(data.ID, data.Email, password, data.Fname, data.Lname, data.Privileges)
		if err != nil {
			helpers.SuccessResponse(false, w, r)
			helpers.ThrowErr(w, r, "Editing user error", err)
			return
		}
	}

	helpers.SuccessResponse(true, w, r)
	if err != nil {
		helpers.ThrowErr(w, r, "JSON encoding error", err)
	}
}

// New is the handler for the new user request.
func New(w http.ResponseWriter, r *http.Request) {
	var data edit                                // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		helpers.SuccessResponse(