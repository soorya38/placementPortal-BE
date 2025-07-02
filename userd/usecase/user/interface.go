package user

import "backend/userd/entity"

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	GetUserByUsername(username string) (*entity.User, error)
	ListUser() ([]*entity.User, error)
}

type Writer interface {
	CreateUser(username, password, email, role string) (*entity.User, error)
	DeleteUser(id string) error
}

type Usecase interface {
	GetUserByUsername(username, password string) (*entity.User, error)
	CreateUser(username, password, email, role string) (*entity.User, error)
	ListUser() ([]*entity.User, error)
	DeleteUser(id string) error
}
