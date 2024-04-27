package usecase

import (
	"fmt"
	"time"

	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/repository"
	"github.com/syahyudi09/ChatPintar/App/utils"
)

type ChatGroupUsecase interface {
	CreateGroup(input model.InputChatGroupModel) (string, error)
	AddUserToGroup(input model.InputUserToGroup) error
	DeleteUserGrou(phoneNumber, groupID string) error
	CreateMessageGroup(message model.MessageGroup) error
	IsUserMemberOfGroup(userID, groupID string) (bool, error)
}

type chatGroupUsecase struct {
	chatGroup repository.ChatGroupRepository
}

func (cgu *chatGroupUsecase) CreateGroup(input model.InputChatGroupModel) (string, error) {
	var group model.GroupModel

	group.GroupName = input.GroupName
	group.Description = input.Description // Memperbaiki nama kolom
	group.GroupID = utils.UuidGenerate()  // Menggunakan UUID untuk ID unik
	group.CreatedAt = time.Now()

	// Menyimpan grup ke database melalui repository
	err := cgu.chatGroup.CreateGroup(group)
	if err != nil {
		return "", fmt.Errorf("failed to create group: %v", err)
	}

	return group.GroupID, nil
}

func (cgu *chatGroupUsecase) AddUserToGroup(input model.InputUserToGroup) error {
	var users model.GroupUsers
	users.PhoneNumber = input.PhoneNumber // Sesuaikan dengan model Anda
	users.GroupId = input.GroupId         // Pastikan GroupID juga ditentukan
	users.Role = input.Role
	users.CreatedAt = time.Now()

	// Menambahkan pengguna ke grup melalui repository
	err := cgu.chatGroup.CreateGroupUser(users)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %v", err)
	}

	return nil
}

func (cgu *chatGroupUsecase) CreateMessageGroup(message model.MessageGroup) error {
	// Membuat ID pesan yang unik
	message.MessageID = utils.UuidGenerate()

	// Menyimpan pesan ke dalam database
	if err := cgu.chatGroup.CreateMessageGroup(message); err != nil {
		return fmt.Errorf("failed to create group message: %v", err)
	}

	return nil
}

func (cgu *chatGroupUsecase) DeleteUserGrou(phoneNumber, groupID string) error {
	if phoneNumber == "" || groupID == "" {
		return fmt.Errorf("phone number or group ID cannot be empty")
	}

	err := cgu.chatGroup.RemoveGroupUser(phoneNumber, groupID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %v", err)
	}

	return nil
}

func (cgu *chatGroupUsecase) IsUserMemberOfGroup(userID, groupID string) (bool, error) {
	if userID == "" || groupID == "" {
		return false, fmt.Errorf("userID and groupID cannot be empty")
	}

	return cgu.chatGroup.CheckUserInGroup(userID, groupID)
}

func NewChatGroupUsecase(chatGroup repository.ChatGroupRepository) ChatGroupUsecase {
	return &chatGroupUsecase{
		chatGroup: chatGroup,
	}
}
