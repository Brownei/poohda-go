package api

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/go-playground/validator/v10"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
	"gopkg.in/gomail.v2"
)

var (
	content       embed.FS
	ZOHO_EMAIL    = os.Getenv("ZOHO_EMAIL")
	ZOHO_PASSWORD = os.Getenv("ZOHO_PASSWORD")
)

type ImageUrlForMail struct {
	LogoUrl string
}

func (a *application) SendMail(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	var body bytes.Buffer
	var payload types.SubscribePayload
	cld, err := utils.InitializeCloudinary()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	result, _ := cld.Admin.AssetByAssetID(r.Context(), admin.AssetByAssetIDParams{
		AssetID: "b3fcb62e3a906ed8af10449f240fdf9c",
	})

	a.logger.Info(result)
	mailImages := ImageUrlForMail{
		LogoUrl: result.URL,
	}

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

	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/public/index.html", cwd))
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}

	err = tmpl.Execute(&body, mailImages)
	if err != nil {
		log.Fatal("Failed to execute template:", err)
	}

	a.logger.Info("Sending email.......")
	// a.logger.Info(body.String())
	chanErr := make(chan error, 1)

	err = a.store.Waitlist.AddToWaitlist(r.Context(), payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", err))
		return
	}

	go sendMailGoRoutine(chanErr, payload)

	if err := <-chanErr; err != nil {
		fmt.Print(err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	} else {
		utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("%s with email: %s have joined the waitlist", payload.Name, payload.Email))
		a.logger.Info("Sent email.......")
	}
}

func sendMailGoRoutine(chanErr chan error, payload types.SubscribePayload) {
	m := gomail.NewMessage()
	// cwd, _ := os.Getwd()
	// poohdaLogo := fmt.Sprintf("%s/public/PoohDa White green.png", cwd)
	// firstName := strings.Split(payload.Name, " ")
	m.SetHeader("From", "noreply@poohda.com")
	m.SetHeader("To", payload.Email)
	m.SetAddressHeader("Cc", payload.Email, payload.Name)
	m.SetHeader("Subject", "You’re on the Waitlist to Be Da Difference")
	// m.SetBody("text/html", body.String())
	m.SetBody("text/html", `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Email Template</title>
  <style type="text/css">
    @font-face {
      font-family: 'CustomFont';
      src: url('https://res.cloudinary.com/brownson/raw/upload/v1734014059/rw9mzwor2rzzhzgylrtc.ttf') format('truetype');
      font-weight: normal;
      font-style: normal;
    }
  </style>
  <!--<link href="https://res.cloudinary.com/brownson/raw/upload/v1734014059/rw9mzwor2rzzhzgylrtc.ttf" rel="stylesheet">-->
</head>

<body style="color: #ffffff;  background-color: #f4f4f4; padding: 10px;">
  <div style="max-width: 600px; margin: 0 auto; background-color: #000000; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);">

    <div style="text-align: center;">
      <img alt="PoohDa" src="https://res.cloudinary.com/brownson/image/upload/v1734001320/pmbizybnu0aeentwkcza.png"
        style="width: 300px;  padding: 0px; margin: -50px;" />
      <h1 style="font-family: Helvetica, Arial, sans-serif; font-size: 30px; margin-top: -50px; color: #008000;">
        Welcome To Poohda</h1>
    </div>

    <div style="line-height: 1.6;">
      <p style="margin-bottom: 16px;">Hey Dauntless!</p>

      <p style="margin-bottom: 16px;">You made it to the waitlist to be Da Difference—the leader of the pack! That means
        you’ll be the first to know when our exclusive pieces go live on the PooHDa website.</p>

      <p style="margin-bottom: 16px;">So get ready for the launch.</p>

      <p style="margin-bottom: 16px;">Every piece is crafted to break the rules and elevate your style—rare, limited,
        and unimagined. And you’re at the front of the line.</p>

      <p style="margin-bottom: 16px;">Stay close. Your access to Da Difference is just around the corner—gear up to
        elevate your wardrobe with PooHDa—fashion that’s all about being 100%% you, no compromises.</p>


      <p style="margin-top: 40px;">
        <span style="display: block;">Catch you soon,</span>
        <span style="display: block; font-weight: bold;">POOH</span>
        <span style="display: block;">Creative Director, PooHDa</span>
      </p>
    </div>

    <div style="text-align: center; margin-top: 20px; font-size: 12px; color: #aaa;">
      <p>&copy; 2024 PooHDa. All rights reserved.</p>
    </div>
  </div>
</body>

</html>
    `)

	d := gomail.NewDialer("smtppro.zoho.com", 465, ZOHO_EMAIL, ZOHO_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
		chanErr <- err
	}

	chanErr <- nil
}
