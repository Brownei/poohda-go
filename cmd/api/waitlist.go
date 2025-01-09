package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/utils"
)

func (a *application) AllWaitlistRoutes(r chi.Router) {
	r.Get("/", a.GetAllWaitlistParticipants)
}

func (a *application) GetAllWaitlistParticipants(w http.ResponseWriter, r *http.Request) {
	waitlist, err := a.store.Waitlist.GetAllWaitlistParticipants()
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, waitlist)
}
