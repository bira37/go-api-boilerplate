package repository

import (
	repositoryContract "bira.io/template/contract/repository"
	"bira.io/template/infra"
)

type RepositoryCollection struct {
	UserRepository repositoryContract.UserRepository
}

func NewRepositoryCollection(infraCollection infra.InfraCollection) *RepositoryCollection {
	return &RepositoryCollection{
		UserRepository: NewUserRepository(infraCollection),
	}
}
