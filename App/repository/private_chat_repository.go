package repository

import (
	"database/sql"
	"fmt"

	"github.com/syahyudi09/ChatPintar/App/model"
)

type PrivateChatReposiotry interface {
	CreateMessage(massage model.MessageModel) error
	CreateChat(privateChat model.PrivateChatModel) (string, error)
	FindChatByUsers(userID1,userID2 string) (*model.PrivateChatModel, error)
	UpdateMessageStatusBySender(senderID string, newStatus string) error 
}

type privateChatRepository struct {
	db *sql.DB
}

func NewPrivateChatReposiotry(db *sql.DB) PrivateChatReposiotry {
	return &privateChatRepository{
		db: db,
	}
}

func (pr *privateChatRepository) CreateMessage(message model.MessageModel) error {
	_, err := pr.db.Exec(`
		INSERT INTO message (message_id, chat_id, message_content, created_at, status, sender_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		message.MessageID,
		message.ChatID, // Menggunakan chat_id yang benar
		message.MessageContent,
		message.CreatedAt,
		message.Status,
		message.SenderID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert into message: %w", err)
	}

	return nil
}


func (pr *privateChatRepository) CreateChat(privateChat model.PrivateChatModel) (string, error) {
	tx, err := pr.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}

	var chatID string
	err = tx.QueryRow(`
		INSERT INTO private_chat (chat_id, user_id1, user_id2, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING chat_id`,
		privateChat.ChatId,
		privateChat.UserID1,
		privateChat.UserID2,
		privateChat.CreatedAt,
		privateChat.UpdatedAt,
	).Scan(&chatID)

	if err != nil {
		tx.Rollback() // Rollback jika gagal
		return "", fmt.Errorf("failed to insert into private_chat: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return chatID, nil
}


func (pr *privateChatRepository) FindChatByUsers(userID1, userID2 string) (*model.PrivateChatModel, error) {
    var privateChat model.PrivateChatModel
    err := pr.db.QueryRow(`
        SELECT chat_id, user_id1, user_id2, created_at, updated_at
        FROM private_chat
        WHERE (user_id1 = $1 AND user_id2 = $2) OR (user_id1 = $2 AND user_id2 = $1)`,
        userID1, userID2,
    ).Scan(&privateChat.ChatId, &privateChat.UserID1, &privateChat.UserID2, &privateChat.CreatedAt, &privateChat.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, nil // Jika tidak ada hasil, kembalikan nil
    }

    if err != nil {
        return nil, err // Kembalikan kesalahan jika terjadi kesalahan
    }

    return &privateChat, nil
}


func (pcr *privateChatRepository) UpdateMessageStatusBySender(senderID string, newStatus string) error {
    query := "UPDATE message SET status = $1 WHERE sender_id = $2"

    _, err := pcr.db.Exec(query, newStatus, senderID)
    if err != nil {
        return fmt.Errorf("failed to update message status: %w", err)
    }

    return nil
}




