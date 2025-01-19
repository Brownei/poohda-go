package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/poohda-go/types"
	"github.com/poohda-go/utils"
)

func (a *application) AllOrdersRoutes(r chi.Router) {
	r.Get("/", a.GetAllOrders)
	r.Post("/", a.CreateANewOrder)
	r.Get("/{order}", a.GetASingleOrder)
}

func (a *application) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := a.store.Orders.GetAllOrders()
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

func (a *application) CreateANewOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.OrderDTO
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	err = utils.ValidateJson(payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	areThereClothes, err := a.store.Clothes.GetAllClothes()
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	if len(areThereClothes) == 0 {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("You cannot order when there are no clothes"))
		return
	}

	for _, clotheId := range payload.ClothesBought {
		_, err := a.store.Clothes.GetOneClothes(ctx, clotheId.Id)
		if err != nil {
			utils.WriteError(w, http.StatusNotAcceptable, err)
			return
		}
	}

	newOrder, err := a.store.Orders.CreateANewOrder(ctx, payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, newOrder)
}

func (a *application) GetASingleOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	param := chi.URLParam(r, "order")
	orderId, err := strconv.Atoi(param)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	order, err := a.store.Orders.GetASingleOrder(ctx, orderId)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, order)
}
