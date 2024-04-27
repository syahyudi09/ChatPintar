package repository

import (
	"database/sql"
	"fmt"

	"github.com/syahyudi09/ChatPintar/App/model"
)

type ChatGroupRepository interface {
	CreateGroup(newGroup model.GroupModel) error
	CreateGroupUser(newGorupUsers model.GroupUsers) error
	CheckUserInGroup(userID, groupID string) (bool, error)
	RemoveGroupUser(userID, groupID string) error
}

type chatGroupRepository struct {
	db *sql.DB
}

func (cgr *chatGroupRepository) CreateGroup(newGroup model.GroupModel) error {
	insertQuery := "INSERT INTO groups(group_id, group_name, group_description, created_at) VALUES ($1, $2, $3, $4)"

	_, err := cgr.db.Exec(insertQuery, newGroup.GroupID, newGroup.GroupName, newGroup.Description, newGroup.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	return nil
}

func (cgr *chatGroupRepository) CreateGroupUser(newGorupUsers model.GroupUsers) error {
	insertQuery := "INSERT INTO group_users(phone_number, group_id, role, joined_at, added_by) VALUES ($1, $2, $3, $4, $5)"

	_, err := cgr.db.Exec(insertQuery, newGorupUsers.PhoneNumber, newGorupUsers.GroupId, newGorupUsers.Role, newGorupUsers.CreatedAt, newGorupUsers.AddByPhoneNumber)
	if err != nil {
		return fmt.Errorf("user is already in the group")
	}
	return nil
}

func (cgr *chatGroupRepository) RemoveGroupUser(phoneNumber, groupID string) error {
	deleteQuery := "DELETE FROM group_users WHERE phone_number = $1 AND group_id = $2"

	// Menjalankan perintah DELETE
	_, err := cgr.db.Exec(deleteQuery, phoneNumber, groupID)
	if err != nil {
		return fmt.Errorf("error removing user from group: %v", err)
	}
	return nil
}

func (cgr *chatGroupRepository) CheckUserInGroup(userID, groupID string) (bool, error) {
	var count int

	// SQL untuk memeriksa apakah pengguna ada dalam grup
	query := "SELECT COUNT(*) FROM group_users WHERE phone_number = $1 AND group_id = $2"

	err := cgr.db.QueryRow(query, userID, groupID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is in group: %v", err)
	}

	// Jika count > 0, berarti pengguna ada dalam grup
	return count > 0, nil
}

func NewChatGroupRepository(db *sql.DB) ChatGroupRepository {
	return &chatGroupRepository{
		db: db,
	}
}
