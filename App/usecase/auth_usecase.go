package usecase

import (
	"fmt"
	"time"

	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/repository"
	"github.com/syahyudi09/ChatPintar/App/utils"
	"github.com/syahyudi09/ChatPintar/App/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(input model.UserInput) error
	PhoneNumberExits(PhoneNumber string) (bool, error)
	Login(input model.UserInputLogin) (model.UserFormatter, error)
	FindByUserId(UserId string) (model.UserModel, error)
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
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("err %w", err)
	}
	user.Password = string(hashedPassword)
	return au.repository.Register(user)
}

func (au *authUsecase) PhoneNumberExits(PhoneNumber string) (bool, error) {
	return au.repository.CheckPhoneNumber(PhoneNumber)
}

func (au *authUsecase) FindByUserId(UserId string) (model.UserModel, error) {
	user, err := au.repository.FindByPhoneNumber(UserId)
	if err != nil {
		return model.UserModel{}, fmt.Errorf("phone number not found")
	}
	return user, nil
}

func (au *authUsecase) Login(input model.UserInputLogin) (model.UserFormatter, error) {
	phoneNumber := input.PhoneNumber
	password := input.Password

	user, err := au.repository.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return model.UserFormatter{}, fmt.Errorf("phone number not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return model.UserFormatter{}, fmt.Errorf("invalid password")
	}

	accessToken, err := token.GenerateToken(user.UserId, user.PhoneNumber)
	if err != nil {
		return model.UserFormatter{}, fmt.Errorf("failed to generate token: %w", err)
	}

	formatter := model.UserFormatter{
		ID:           user.UserId,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		AccessToken:  accessToken,
	}

	return formatter, nil
}

func NewAuthUsecase(repository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		repository: repository,
	}
}
