package model

import "time"

type PrivateChatModel struct {
	ChatId    string
	UserID1   string
	UserID2   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageModel struct {
	MessageID      string `json:"id"`      // ID unik pesan
	ChatID         string `json:"chat_id"` // ID percakapa
	MessageContent string `json:"content"`
	Status         StatusEnum
	CreatedAt      time.Time `json:"created_at"` // Waktu pembuatan pesan
	SenderID string
}

type InputMessageModel struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Status StatusEnum
	Message    string `json:"message"`
}

type InputMassageStatus struct{
	Status StatusEnum
}
