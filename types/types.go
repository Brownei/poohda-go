package types

type SubscribePayload struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Number string `json:"number" validate:"required"`
}

type Category struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsFeatured  bool     `json:"is_featured"`
	Pictures    []string `json:"pictures"`
}

type CategoryDTO struct {
	Name        string   `json:"name" validate:"required,min=3"`
	Description string   `json:"description" validate:"required"`
	IsFeatured  bool     `json:"is_featured"`
	Pictures    []string `json:"pictures" validate:"required"`
}

type Clothes struct {
	Id          int      `json:"id" `
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	Quantity    int      `json:"quantity"`
	CategoryId  int      `json:"category_id"`
	Pictures    []string `json:"pictures"`
	Sizes       []string `json:"sizes"`
}

type ClothesDTO struct {
	CategoryId  int      `json:"category_id" validate:"required"`
	Name        string   `json:"name" validate:"required,min=3"`
	Price       int      `json:"price" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Quantity    int      `json:"quantity" validate:"required"`
	Pictures    []string `json:"pictures" validate:"required"`
	Sizes       []string `json:"sizes" validate:"required"`
}

type Sizes struct {
	Id   int    `json:"id"`
	Size string `json:"size"`
}

type Pictures struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Waitlist struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Number string `json:"number"`
}

type Order struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	IsDelivered   bool     `json:"is_delivered"`
	Quantity      int      `json:"quantity"`
	Price         int      `json:"price"`
	ClothesBought []string `json:"clothes_bought"`
}

type OrderDTO struct {
	Name          string          `json:"name" validate:"required,min=3"`
	Address       string          `json:"address" validate:"required,min=3"`
	IsDelivered   bool            `json:"is_delivered"`
	Quantity      int             `json:"quantity" validate:"required"`
	Price         int             `json:"price" validate:"required"`
	ClothesBought []ClothesBought `json:"clothes_bought" validate:"required"`
}

type ClothesBought struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
