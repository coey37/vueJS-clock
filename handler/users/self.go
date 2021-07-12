package users

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/VolticFroogo/Froogo/middleware"
	"github.com/VolticFroogo/Froogo/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/zemirco/uid"
)

// Settings is the handler for a user editting their own settings.
func Settings(w http.ResponseWriter, r *http.Request) {
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

	if data.Password == "" {
		err = db.EditSelfNoPassword(user.UUID, data.Fname, data.Lname)
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

		err = db.EditSelf(user.UUID, password, data.Fname, data.Lname)
		if err != nil {
			helpers.SuccessResponse(false, w, r)
			helpers.ThrowErr(w, r, "Editing user error", err)
			return
		}
	}

	if data.Email != user.Email {
		user.Fname = data.Fname
		user.Lname = data.Lname

		if err := SendEmailVerification(user, data.Email); err != nil {
			helpers.SuccessResponse(false, w, r)
			helpers.ThrowErr(w, r, "Sending verification email error", err