package manager

import (
	"sync"

	"github.com/syahyudi09/ChatPintar/App/repository"
)

type RepositoryManager interface {
	GetAuthRepo() repository.AuthRepository
}

type repositoryManager struct {
	infra    InfraManager
	authRepo repository.AuthRepository
}

func NewRepoManager(infra InfraManager) RepositoryManager {
	return &repositoryManager{
		infra: infra,
	}
}

var onceLoadAuthRepo sync.Once

func (rm repositoryManager) GetAuthRepo() repository.AuthRepository{
	onceLoadAuthRepo.Do(func() {
		rm.authRepo = repository.NewAuthRepository(rm.infra.GetDB())
	})
	return rm.authRepo
}