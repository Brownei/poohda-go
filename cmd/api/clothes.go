package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
)

func (a *application) AllClothingRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/", a.CreateNewClothing)
		r.Get("/", a.GetAllClothings)
		r.Get("/{id}", a.GetOneClothing)
		r.Get("/search/{search}", a.GetClothesThroughName)
	})
}

func (a *application) CreateNewClothing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.ClothesDTO
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// Check if there are even categories at all Because without categories there cannot be any product
	categories, err := a.store.Categories.GetAllCategories()
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("You will need to create a category to create a product"))
			return
		}

		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	if len(categories) == 0 {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("You will need to create a category to create a product"))
		return
	}

	_, err = a.store.Categories.GetOneCategory(ctx, payload.CategoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("No category like this"))
			return
		}

		fmt.Print("Errorroor")
		utils.WriteError(w, http.StatusConflict, err)
		return
	}
	a.logger.Info(payload)
	newClothings, err := a.store.Clothes.CreateNewClothes(ctx, payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, newClothings)
}

func (a *application) GetAllClothings(w http.ResponseWriter, r *http.Request) {
	clothings, err := a.store.Clothes.GetAllClothes()
	if err != nil {
		if err == sql.ErrNoRows {
			a.logger.Infof("GetAllClpthng: %v", err)
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("No clothings at all"))
			return
		}

		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, clothings)
}

func (a *application) GetOneClothing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	clotheId, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}
	clothings, err := a.store.Clothes.GetOneClothes(ctx, clotheId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("No clothing like this exists!"))
			return
		}

		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, clothings)
}

func (a *application) GetClothesThroughName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	searchString := chi.URLParam(r, "search")

	clothes, err := a.store.Clothes.GetClothesThroughName(ctx, searchString)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, clothes)
}
