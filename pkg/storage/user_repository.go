package storage

import (
	"EM_test_task/internal/entities"
	"fmt"
	"log/slog"
)

//type Authorization interface {
//	CreateUser(user entities.User) (int, error)
//	GetUser(username string) (entities.User, error)
//}

func (usr *personRepository) CreateUser(user entities.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (username,password) VALUES ($1, $2) RETURNING id")
	row := usr.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (usr *personRepository) GetUser(username string) (entities.User, error) {
	query := fmt.Sprintf("SELECT id, password FROM users WHERE username=$1")
	row := usr.db.QueryRow(query, username)

	var user entities.User

	if err := row.Scan(
		&user.Username,
		&user.Password,
	); err != nil {
		slog.Error("GetUser scan process Error")
		return user, fmt.Errorf("error scanning person row: %v", err)
	}
	return user, nil
}
