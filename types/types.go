package types

type SubscribePayload struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Number string `json:"number" validate:"required"`
}
