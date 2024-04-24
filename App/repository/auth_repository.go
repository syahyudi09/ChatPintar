package repository

import (
	"database/sql"
	"fmt"

	"github.com/syahyudi09/ChatPintar/App/model"
)

type AuthRepository interface {
	Register(newUser model.UserModel) error
	CheckPhoneNumber(PhoneNumber string) (bool, error)
}

type authRepository struct {
	db *sql.DB
}

func (ar *authRepository) Register(newUser model.UserModel) error {
	insertQuery := "INSERT INTO users(user_id, name, phone_number ,created_at, password) VALUES($1,$2,$3,$4,$5)"
	_, err := ar.db.Exec(insertQuery, newUser.UserId, newUser.Name, newUser.PhoneNumber, newUser.CreatedAt, newUser.Password)
	if err != nil {
		return fmt.Errorf("err on authRepository.Register: %v", err)
	}
	return nil
}

func (u *authRepository) CheckPhoneNumber(PhoneNumber string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1)"

	var exists bool
	err := u.db.QueryRow(query, PhoneNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("err on authRepository.CheckPhoneNumber: %w", err)
	}

	return exists, nil
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
