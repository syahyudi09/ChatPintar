package model

import "time"

type GroupModel struct {
	GroupID     string
	GroupName   string
	Description string
	CreatedAt   time.Time
}

type GroupUsers struct {
	AddByPhoneNumber string
	PhoneNumber string
	GroupId     string
	Role        RoleEnum
	CreatedAt   time.Time
	LeftAt      time.Time
}

type MessageGroup struct {
	MessageID      string
	MessageContent string
	SenderID       string
	GroupId        string
	Status         StatusEnum
	CreatedAt      time.Time
}

type InputChatGroupModel struct {
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}
type InputUserToGroup struct {
	AddPhoneNumber string
	PhoneNumber    string   `json:"phone_number"`
	GroupId        string   `json:"group_id"`
	Role           RoleEnum `json:"role"`
}

type InputDeleteUsersGroup struct {
	PhoneNumber string `json:"phone_number"`
	GroupID     string `json:"group_id"`
}