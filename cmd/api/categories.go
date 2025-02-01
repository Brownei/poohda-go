package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
)

func (a *application) AllCategoryRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/", a.GetAllCategories)
		r.Post("/", a.CreateNewCategory)
		r.Get("/{category}", a.GetOneCategory)
		r.Get("/{category}/clothes", a.GetAllClothingReferenceToCategory)
		r.Put("/{category}", a.EditCategory)
		r.Delete("/{category}", a.DeleteCategory)
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

func (a *application) GetOneCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryName := chi.URLParam(r, "category")
	id, err := strconv.Atoi(categoryName)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot convert to int"))
		return
	}

	category, err := a.store.Categories.GetOneCategory(ctx, id)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot convert to int"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, category)
}

func (a *application) GetAllClothingReferenceToCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryName := chi.URLParam(r, "category")
	id, err := strconv.Atoi(categoryName)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot convert to int"))
		return
	}
	clothings, err := a.store.Categories.GetAllClothesReferenceToACategory(ctx, id)
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

func (a *application) EditCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := chi.URLParam(r, "category")
	var payload types.CategoryDTO
	ctx := r.Context()
	id, err := strconv.Atoi(categoryName)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot convert to int"))
		return
	}

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

	category, err := a.store.Categories.EditCategory(ctx, id, payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, category)
}

func (a *application) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := chi.URLParam(r, "category")
	ctx := r.Context()
	id, err := strconv.Atoi(categoryName)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("Cannot convert to int"))
		return
	}

	category, err := a.store.Categories.DeleteCategory(ctx, id)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, category)
}
