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
	var data edit