package manager

import (
	"sync"

	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type UseceaseManager interface {
	GetAuthUsecase() usecase.AuthUsecase
}

type usecaseManager struct {
	rm              RepositoryManager
	authUsecase     usecase.AuthUsecase
}

var onceLoadAuthUsecase sync.Once

func (um *usecaseManager) GetAuthUsecase() usecase.AuthUsecase {
	onceLoadAuthUsecase.Do(func() {
		um.authUsecase = usecase.NewAuthUsecase(
			um.rm.GetAuthRepo(),
		)
	})
	return um.authUsecase
}

func NewUsecaseManager(rm RepositoryManager) UseceaseManager {
	return &usecaseManager{
		rm: rm,
	}
}