package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	Validator          = validator.New()
	CLOUDINARY_SECRET  = os.Getenv("CLOUDINARY_SECRET")
	CLOUDINARY_API_KEY = os.Getenv("CLOUDINARY_API_KEY ")
)

func ChangeFontToBase64(laptopPath string) (string, error) {
	path, err := filepath.Abs(laptopPath)
	if err != nil {
		return "", err
	}

	// Read the font file
	fontData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Encode the font file data to base64
	encodedFont := base64.StdEncoding.EncodeToString(fontData)

	// Print the base64 string
	return encodedFont, nil
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("No body in this request")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func ValidateJson(payload any) error {
	// Validate the payload
	if err := Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	return nil
}

func VerifyPassword(encryptedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func JwtToken(email string, ctx context.Context) string {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	expiryTime := time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days in seconds
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,    // Subject (user identifier)
		"iss": "poohda", // Issuer
		"exp": expiryTime,
		"iat": time.Now().Unix(), // Issued at
	})

	token, _ := claims.SignedString(secretKey)
	fmt.Printf("Token claims added: %+v\n", token)
	return token
}

func VerifyToken(token string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Error: %v", ok)
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Error in the verified token: %s\n%v", err.Error(), verifiedToken)
	}

	// Check if the token is valid
	if !verifiedToken.Valid {
		return "", fmt.Errorf("Not Valid: %v", err)
	}

	email, _ := verifiedToken.Claims.GetSubject()
	// Return the verified token
	log.Printf("VerifiedToken: %v\n", email)

	return email, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func InitializeCloudinary() (*cloudinary.Cloudinary, error) {
	fmt.Print(CLOUDINARY_API_KEY)
	fmt.Print(CLOUDINARY_SECRET)
	cld, err := cloudinary.NewFromParams("brownson", "491131249681159", CLOUDINARY_SECRET)

	return cld, err
}
