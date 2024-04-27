package manager

import (
	"sync"

	"github.com/syahyudi09/ChatPintar/App/repository"
)

type RepositoryManager interface {
	GetAuthRepo() repository.AuthRepository
	GetPrivateRepo() repository.PrivateChatReposiotry
}

type repositoryManager struct {
	infra    InfraManager
	authRepo repository.AuthRepository
	privateChatRepo repository.PrivateChatReposiotry
}

func NewRepoManager(infra InfraManager) RepositoryManager {
	return &repositoryManager{
		infra: infra,
	}
}

var onceLoadAuthRepo sync.Once
var onceLoadPrivateChatrepo sync.Once

func (rm *repositoryManager) GetAuthRepo() repository.AuthRepository{
	onceLoadAuthRepo.Do(func() {
		rm.authRepo = repository.NewAuthRepository(rm.infra.GetDB())
	})
	return rm.authRepo
}

func (rm *repositoryManager) GetPrivateRepo() repository.PrivateChatReposiotry{
	onceLoadPrivateChatrepo.Do(func() {
		rm.privateChatRepo = repository.NewPrivateChatReposiotry(rm.infra.GetDB())
	})
	return rm.privateChatRepo
}

