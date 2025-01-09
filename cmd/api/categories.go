package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
)

func (a *application) AllCategoryRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/", a.GetAllCategories)
		r.Post("/", a.CreateNewCategory)
		r.Get("/{category}/clothes", a.GetAllClothingReferenceToCategory)
	})
}

func (a *application) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.CategoryDTO

	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Printf("Parsing error: %v", err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	if err := utils.ValidateJson(payload); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	category, err := a.store.Categories.CreateNewCategory(ctx, payload)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"category_name_key\"" {
			utils.WriteError(w, http.StatusConflict, fmt.Errorf("There's already a category like this!"))
			return
		}
		log.Printf("Creating category error: %v", err)
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, category)
}

func (a *application) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := a.store.Categories.GetAllCategories()
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, "Nothing here")
		}
		utils.WriteError(w, http.StatusBadGateway, err)
	}

	utils.WriteJSON(w, http.StatusOK, categories)
}

func (a *application) GetAllClothingReferenceToCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryName := chi.URLParam(r, "category")
	clothings, err := a.store.Categories.GetAllClothesReferenceToACategory(ctx, categoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("No clothings in this category"))
			return
		}

		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, clothings)
}
