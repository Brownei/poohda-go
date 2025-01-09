package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
)

func (a *application) AllAuthRoutes(r chi.Router) {
	r.Post("/login", a.Login)
}

func (a *application) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.LoginDto

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	if err := utils.ValidateJson(payload); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	if payload.Username != "treasurepooh" || payload.Password != "treasurepooh" {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Invalid password"))
		return
	}

	token := utils.JwtToken(payload.Username, ctx)
	utils.WriteJSON(w, http.StatusAccepted, token)
}
