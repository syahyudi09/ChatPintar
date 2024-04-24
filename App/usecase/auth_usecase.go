package usecase

import (
	"fmt"
	"time"

	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/repository"
	"github.com/syahyudi09/ChatPintar/App/utils"
)

type AuthUsecase interface {
	Register(input model.UserInput) error 
	PhoneNumberExits(PhoneNumber string) (bool, error)
}

type authUsecase struct {
	repository repository.AuthRepository
}

func (au *authUsecase) Register(input model.UserInput) error {
	user := model.UserModel{}
	user.UserId = utils.UuidGenerate()
	user.Name = input.Name
	user.PhoneNumber = input.PhoneNumber
	user.CreatedAt = time.Now()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("err %w", err)
	}
	user.Password = string(hashedPassword)
	return au.repository.Register(user)
}

func (au *authUsecase) PhoneNumberExits(PhoneNumber string) (bool, error) {
	return au.repository.CheckPhoneNumber(PhoneNumber)
}

func NewAuthUsecase(repository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		repository: repository,
	}
}
