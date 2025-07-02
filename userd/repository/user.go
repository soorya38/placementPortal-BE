package repository

import (
	"backend/userd/entity"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(username, password, email, role string) (*entity.User, error) {
	query := `
		INSERT INTO users (id, username, password, email, role, created_at) 
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		RETURNING id, username, email, role, created_at`

	now := time.Now()
	row := r.db.QueryRow(query, username, password, email, role, now)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*entity.User, error) {
	query := `
		SELECT id, username, email, role, password, created_at 
		FROM users 
		WHERE username = $1`

	row := r.db.QueryRow(query, username)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ListUser() ([]*entity.User, error) {
	query := `
		SELECT id, username, email, role, created_at 
		FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *Repository) DeleteUser(id string) error {
	query := `
		DELETE FROM users 
		WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
