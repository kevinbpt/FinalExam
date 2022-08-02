package model

import "time"

type Photo struct {
	Id        int       `json:"Id,omitempty"`
	Title     string    `json:"Title,omitempty" validate:"required"`
	Caption   string    `json:"Caption,omitempty"`
	PhotoUrl  string    `json:"PhotoUrl,omitempty" validate:"required"`
	UserId    int       `json:"UserId,omitempty"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty"`
}

type PhotoResponse struct {
	Id        int       `json:"Id,omitempty"`
	Title     string    `json:"Title,omitempty" validate:"required"`
	Caption   string    `json:"Caption,omitempty"`
	PhotoUrl  string    `json:"PhotoUrl,omitempty" validate:"required"`
	UserId    int       `json:"UserId,omitempty"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty"`
	User      ResponseUser
}

type Comment struct {
	Id        int
	UserId    int
	PhotoId   int
	Message   string `json:"Message" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentResponse struct {
	Id        int
	UserId    int
	PhotoId   int
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      ResponseUser
	Photo     Photo
}
