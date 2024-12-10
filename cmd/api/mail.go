package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
	"gopkg.in/gomail.v2"
)

var (
	ZOHO_EMAIL    = os.Getenv("ZOHO_EMAIL")
	ZOHO_PASSWORD = os.Getenv("ZOHO_PASSWORD")
)

func (a *application) SendMail(w http.ResponseWriter, r *http.Request) {
	var payload types.SubscribePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot be able to parse json"))
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		fmt.Printf("Error: %s", errors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	a.logger.Info("Sending email.......")
	a.logger.Info(payload)
	chanErr := make(chan error, 1)

	go sendMailGoRoutine(chanErr, payload)

	if err := <-chanErr; err != nil {
		fmt.Print(err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	} else {
		utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Mail sent successfully"))
		a.logger.Info("Sent email.......")
	}
}

func sendMailGoRoutine(chanErr chan error, payload types.SubscribePayload) {
	encoded, err := utils.ChangeFontToBase64("public/fonts/Hellion.ttf")
	if err != nil {
		chanErr <- err
	}

	m := gomail.NewMessage()
	// firstName := strings.Split(payload.Name, " ")
	m.SetHeader("From", "noreply@poohda.com")
	m.SetHeader("To", payload.Email)
	m.SetAddressHeader("Cc", payload.Email, payload.Name)
	m.SetHeader("Subject", "You’re on the Waitlist to Be Da Difference")
	m.SetBody("text/html", fmt.Sprintf(`
    <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
      @font-face {
        font-family: 'Hellion';
        src: url('http://localhost:8000/public/fonts/Hellion.ttf') format('truetype');
      };

      body {
        font-family: 'Hellion', Arial, sans-serif;
      };
    </style>
</head>
<body style="line-height: 1.6; color: #333; padding: 20px; background-color: #f9f9f9;">
    <table style="width: 100%%; max-width: 600px; margin: 0 auto; background: #000; color: #fff; border-radius: 8px; overflow: hidden;">
        <!-- Header Section -->
        <tr style="text-align: center;">
            <td style="padding: 20px; background-image: url('/public/poohda.png'); background-size: cover; background-position: center;">
                <h1 style="margin: 0; position: relative; z-index: 1;">The Leader of the Pack</h1>
            </td>
        </tr>

        <!-- Body Section -->
        <tr>
            <td style="padding: 20px;">
                <p>
                    You made it to the waitlist to be <strong>Da Difference</strong>—the leader of the pack! That means you’ll be the first to know when our exclusive pieces go live on the PooHDa website.
                </p>
                <p><strong>So get ready for the launch!</strong></p>
                <p>
                    Every piece is crafted to break the rules and elevate your style—rare, limited, and unimagined. And you’re at the front of the line.
                </p>
                <p>
                    Stay close. Your access to Da Difference is just around the corner—gear up to elevate your wardrobe with PooHDa—fashion that’s all about being 100%% you, no compromises.
                </p>
            </td>
        </tr>

        <!-- CTA Button -->
        <tr>
            <td style="text-align: center; padding: 20px;">
                <a href="#" style="
                    background: #86EFAC; 
                    color: #000; 
                    text-decoration: none; 
                    padding: 15px 30px; 
                    font-size: 16px; 
                    font-weight: bold; 
                    border-radius: 5px;">
                    Be Da Difference
                </a>
            </td>
        </tr>

        <!-- Footer Section -->
        <tr>
            <td style="padding: 20px; text-align: center; font-size: 14px; color: #fff;">
                <p>Catch you soon,</p>
                <p><strong>POOH</strong><br>Creative Director, PooHDa</p>
            </td>
        </tr>
    </table>
</body>
</html>
`))

	d := gomail.NewDialer("smtppro.zoho.com", 465, ZOHO_EMAIL, ZOHO_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
		chanErr <- err
	}

	chanErr <- nil
}
