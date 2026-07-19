package request

type Register struct {
	Name     string `json:"name" validate:"required,min=2,max=100" example:"Ada Lovelace"`
	Email    string `json:"email" validate:"required,email,max=254" example:"ada@example.com"`
	Password string `json:"password" validate:"required,min=12,max=72" example:"correct horse battery staple"`
}
