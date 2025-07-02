package user

import (
	"backend/userd/entity"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Usecase {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(username, password, email, role string) (*entity.User, error) {
	user, err := s.repo.CreateUser(username, password, email, role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByUsername(username, password string) (*entity.User, error) {
	// Mock user data
	// users := &entity.User{
	// 	ID:        "1",
	// 	Username:  "test",
	// 	Email:     "admin@company.com",
	// 	Role:      "Admin",
	// 	CreatedAt: "2024-01-01T00:00:00Z",
	// }

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	users := &entity.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	if user.Password != password {
		return nil, errors.New("Invalid username or password")
	}

	return users, nil
}

func (s *Service) ListUser() ([]*entity.User, error) {
	users, err := s.repo.ListUser()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
