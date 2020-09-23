
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
		return // There is no captcha response.
	}
	captchaSuccess, err := captcha.Verify(data.Captcha, r.Header.Get("CF-Connecting-IP")) // Check the captcha.
	if err != nil {
		helpers.JSONResponse(response{Code: Recaptcha}, w)
		helpers.ThrowErr(w, r, "Recaptcha error", err)
		return
	}
	if !captchaSuccess {
		helpers.JSONResponse(response{Code: Recaptcha}, w)
		return // Unsuccessful captcha.
	}

	user, err := db.GetUserFromEmail(data.Email)
	if err != nil {
		helpers.JSONResponse(response{Code: Internal}, w)
		helpers.ThrowErr(w, r, "Getting user error", err)
		return
	}
	if user.UUID == 0 {
		helpers.JSONResponse(response{Code: InvalidEmail}, w)
		return
	}

	id := uid.New(64)

	err = db.AddRecovery(id, user.UUID, data.Email)
	if err != nil {
		helpers.JSONResponse(response{Code: Internal}, w)
		helpers.ThrowErr(w, r, "Adding recovery error", err)
		return
	}

	err = SendEmail(id, data.Email)
	if err != nil {
		helpers.JSONResponse(response{Code: SendingEmail}, w)
		helpers.ThrowErr(w, r, "Send email error", err)
		return
	}

	helpers.JSONResponse(response{Code: Success}, w)
}

// SendEmail sends the recovery email.
func SendEmail(id, email string) (err error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		return
	}

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("To recover your password please click this link: <a href=\"https://froogo.co.uk/password-recovery?code=" + id + "\">recover password</a>."),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("To recover your password please click this link: https://froogo.co.uk/password-recovery?code=" + id),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),