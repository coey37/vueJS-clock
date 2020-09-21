
package recovery

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-recaptcha/recaptcha"
	"github.com/zemirco/uid"
)

var (
	captchaSecret = os.Getenv("CAPTCHA_SECRET")
	captcha       = recaptcha.New(captchaSecret)
)

// Response codes.
const (
	Success = iota
	InvalidEmail
	Recaptcha
	Internal
	SendingEmail
	InvalidCode
)

type message struct {
	Code, Password, Email, Captcha string
}

type response struct {
	Code int
}

// Begin is the function called after an AJAX request is sent from the forgot password page.
func Begin(w http.ResponseWriter, r *http.Request) {
	var data message                             // Create struct to store data.
	err := json.NewDecoder(r.Body).Decode(&data) // Decode response to struct.
	if err != nil {
		helpers.JSONResponse(response{Code: Internal}, w)
		helpers.ThrowErr(w, r, "JSON decoding error", err)
		return
	}

	if data.Captcha == "" {
		helpers.JSONResponse(response{Code: Recaptcha}, w)