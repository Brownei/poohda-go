package types

type SubscribePayload struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Number string `json:"number" validate:"required"`
}

type Category struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type CategoryDTO struct {
	Name    string `json:"name" validate:"required,min=3"`
	Picture string `json:"picture"`
}

type Clothes struct {
	Id          int    `json:"id" `
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	CategoryId  int    `json:"category_id"`
}

type ClothesDTO struct {
	CategoryId  int    `json:"category_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=3"`
	Price       int    `json:"price" validate:"required"`
	Description string `json:"description" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required"`
}

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
