package model

type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}

type LoginUserRequest struct {
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}
