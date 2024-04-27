package usecase

import (
	"fmt"
	"time"

	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/repository"
	"github.com/syahyudi09/ChatPintar/App/utils"
)

type PrivateChatUsecase interface {
	CreateMessage(input model.InputMessageModel) error
	UpdateMessageStatusBySender(senderID string, newStatus string) error
}

type privateChatUsecase struct {
	chat repository.PrivateChatReposiotry
	auth repository.AuthRepository
}

func (pcu *privateChatUsecase) CreateMessage(input model.InputMessageModel) error {
	// Periksa apakah chat sudah ada
	existingChat, err := pcu.chat.FindChatByUsers(input.SenderID, input.ReceiverID)
	if err != nil {
		return fmt.Errorf("error checking existing chat: %w", err)
	}

	var chatID string
	if existingChat == nil {
		chatID = utils.UuidGenerate()
		newChat := model.PrivateChatModel{
			ChatId:    chatID,
			UserID1:   input.SenderID,
			UserID2:   input.ReceiverID,
			CreatedAt: time.Now(),
		}

		chatID, err = pcu.chat.CreateChat(newChat)
		if err != nil {
			return fmt.Errorf("failed to create chat: %w", err)
		}
	} else {
		// Gunakan chat yang sudah ada
		chatID = existingChat.ChatId
	}

	// Buat pesan baru
	message := model.MessageModel{
		MessageID:      utils.UuidGenerate(),
		ChatID:         chatID,
		CreatedAt:      time.Now(),
		MessageContent: input.Message,
		Status:         "pending",
		SenderID:       input.SenderID,
	}

	err = pcu.chat.CreateMessage(message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (pcu *privateChatUsecase) UpdateMessageStatusBySender(senderID string, newStatus string) error {
	err := pcu.chat.UpdateMessageStatusBySender(senderID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}

func NewPrivateChatUsecase(chat repository.PrivateChatReposiotry, auth repository.AuthRepository) PrivateChatUsecase {
	return &privateChatUsecase{
		chat: chat,
		auth: auth,
	}
}
