package service

import (
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/infra"
	"bira.io/template/repository"
)

type ServiceCollection struct {
	AuthService serviceContract.AuthService
	UserService serviceContract.UserService
}

func NewServiceCollection(repositoryCollection *repository.RepositoryCollection, infraCollection infra.InfraCollection) *ServiceCollection {
	serviceCollection := ServiceCollection{}

	serviceCollection.AuthService = NewAuthService(repositoryCollection.UserRepository)
	serviceCollection.UserService = NewUserService(repositoryCollection.UserRepository)

	return &serviceCollection
}
