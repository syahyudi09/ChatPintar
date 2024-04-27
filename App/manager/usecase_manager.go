package manager

import (
	"sync"

	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type UseceaseManager interface {
	GetAuthUsecase() usecase.AuthUsecase
	GetPrivateUsecase() usecase.PrivateChatUsecase
	GetChatGroupUsecase() usecase.ChatGroupUsecase
}

type usecaseManager struct {
	rm              RepositoryManager
	authUsecase     usecase.AuthUsecase
	privateChatUsecase usecase.PrivateChatUsecase
	chatGroupUsecase usecase.ChatGroupUsecase
}

var onceLoadAuthUsecase sync.Once
var onceLoadPrivateChatUsecase sync.Once
var onceLoadChatGroupUsecase sync.Once

func (um *usecaseManager) GetAuthUsecase() usecase.AuthUsecase {
	onceLoadAuthUsecase.Do(func() {
		um.authUsecase = usecase.NewAuthUsecase(
			um.rm.GetAuthRepo(),
		)
	})
	return um.authUsecase
}

func (um *usecaseManager) GetPrivateUsecase() usecase.PrivateChatUsecase{
	onceLoadPrivateChatUsecase.Do(func() {
		um.privateChatUsecase = usecase.NewPrivateChatUsecase(
			um.rm.GetPrivateRepo(),
			um.rm.GetAuthRepo(),
		)
	})
	return um.privateChatUsecase
}

func (um *usecaseManager) GetChatGroupUsecase() usecase.ChatGroupUsecase{
	onceLoadChatGroupUsecase.Do(func() {
		um.chatGroupUsecase = usecase.NewChatGroupUsecase(
			um.rm.GetChatGroupRepo(),
		)
	})
	return um.chatGroupUsecase
}

func NewUsecaseManager(rm RepositoryManager) UseceaseManager {
	return &usecaseManager{
		rm: rm,
	}
}