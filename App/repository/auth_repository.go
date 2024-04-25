package repository

import (
	"database/sql"
	"fmt"

	"github.com/syahyudi09/ChatPintar/App/model"
)

type AuthRepository interface {
	Register(newUser model.UserModel) error
	CheckPhoneNumber(PhoneNumber string) (bool, error)
	FindByPhoneNumber(PhoneNumber string) (model.UserModel , error) 
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

func (ar *authRepository) CheckPhoneNumber(PhoneNumber string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1)"

	var exists bool
	err := ar.db.QueryRow(query, PhoneNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("err on authRepository.CheckPhoneNumber: %w", err)
	}

	return exists, nil
}

func (ar *authRepository) FindByPhoneNumber(PhoneNumber string) (model.UserModel , error) {
	getQuery := "SELECT user_id, name, phone_number, password FROM users WHERE phone_number = $1"

	row := ar.db.QueryRow(getQuery, PhoneNumber)
	var user model.UserModel 
	err := row.Scan(&user.UserId, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, fmt.Errorf("error on authRepository.FindByPhoneNumber : %w", err)
	}
	return user, nil
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
