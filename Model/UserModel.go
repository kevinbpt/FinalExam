package model

import "time"

type CreateUser struct {
	Id       int    `json:"Id,omitempty"`
	Username string `json:"Username,omitempty" validate:"required"`
	Password string `json:"Password,omitempty" validate:"min=6"`
	Email    string `json:"Email,omitempty" validate:"required,email"`
	Age      int    `json:"Age" validate:"numeric,gt=8"`
}

type UpdateUser struct {
	Id        int
	Username  string `json:"Username" validate:"required"`
	Password  string `json:"Password" validate:"min=6"`
	Email     string `json:"Email" validate:"required,email"`
	Age       int    `json:"Age" validate:"numeric,gt=8"`
	UpdatedAt time.Time
}

type ResponseUser struct {
	Id       int    `json:"Id,omitempty"`
	Username string `json:"Username,omitempty"`
	Email    string `json:"Email,omitempty"`
}

type SocialMedia struct {
	Id             int
	Name           string `json:"Name" validate:"required"`
	SocialMediaUrl string `json:"SocialMediaUrl" validate:"required"`
	UserId         int
}
