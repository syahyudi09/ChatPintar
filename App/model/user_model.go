package model

import "time"

type UserModel struct {
	UserId      string
	Name        string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
}

type UserInput struct {
	Name        string `json:"name" validate:"required,max=30"`
	PhoneNumber string `json:"phone_number" validate:"required,max=11"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserInputLogin struct {
	PhoneNumber string `json:"phone_number" validate:"required,max=11"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UserFormatter struct {
	ID           string
	Name         string
	PhoneNumber string
	AccessToken  string
	RefreshToken string
}
