package service

import (
	repositoryContract "bira.io/template/contract/repository"
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/dto"
	"bira.io/template/infra"
)

type userService struct {
	userRepository repositoryContract.UserRepository
}

func NewUserService(ur repositoryContract.UserRepository) serviceContract.UserService {
	return &userService{
		userRepository: ur,
	}
}

func (s *userService) GetLoggedUser(getLoggedUserRequest dto.GetLoggedUserRequest) (dto.GetLoggedUserResponse, error) {
	response, err := s.userRepository.FindUserByUsername(getLoggedUserRequest.Username)

	if err != nil {
		dberr, _ := err.(*infra.SqlDbError)
		if dberr.Code == infra.ErrDbNotFound {
			return dto.GetLoggedUserResponse{}, NewHttpErrNotFound("User not found.")
		}
		return dto.GetLoggedUserResponse{}, NewHttpErrInternalServer("Internal error.")
	}

	return dto.GetLoggedUserResponse{
		Username: response.Username,
		Name:     response.Name,
		Email:    response.Email,
		Id:       response.Id,
	}, nil
}
