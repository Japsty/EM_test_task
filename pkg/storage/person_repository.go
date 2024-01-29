package storage

import (
	"EM_test_task/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/pingcap/log"
	"log/slog"
	"time"
)

// Repository - интерфейс, предоставляющий все доступные методы для работы с базой данных
type Repository interface {
	CreatePerson(ctx context.Context, arg PersonParams) (entities.Person, error)
	GetPeople(page, perPage int) ([]entities.Person, error)
	GetPerson(id int) (entities.Person, error)
	GetPeopleFiltered(sort string, from int, to int, page int, perPage int) ([]entities.Person, error)
	UpdatePerson(id int, person entities.Person) error
	DeletePerson(id int) error
}

type personRepository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	s := &personRepository{
		db: db,
	}

	return s
}

type PersonParams struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Filter struct {
	SortBy string `json:"sort_by"`
	From   int    `json:"from,omitempty"`
	To     int    `json:"to,omitempty"`
}

// CreatePerson - метод для создания человека по параметрам name,surname и по желанию patronymic
func (s *personRepository) CreatePerson(ctx context.Context, arg PersonParams) (entities.Person, error) {
	var person entities.Person

	query := `
		INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING name, surname, patronymic, age, gender, nationality
	`

	err := s.db.QueryRowContext(ctx, query,
		arg.Name,
		arg.Surname,
		arg.Patronymic,
		arg.Age,
		arg.Gender,
		arg.Nationality,
	).Scan(
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)

	if err != nil {
		log.Error("CreatePerson Error")
		return entities.Person{}, fmt.Errorf("error creating person: %v", err)
	}

	slog.Debug("CreatePerson created person with id:", person.ID)
	return person, nil
}

// GetPeople - метод для получения списка людей с пагинацией
func (s *personRepository) GetPeople(page, perPage int) ([]entities.Person, error) {
	offset := (page - 1) * perPage

	query := `
		SELECT * FROM people
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.QueryContext(context.Background(), query, perPage, offset)
	if err != nil {
		slog.Error("GetPeople QueryContext Error")
		return nil, fmt.Errorf("error getting people: %v", err)
	}
	defer rows.Close()

	var people []entities.Person

	for rows.Next() {
		var person entities.Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
			&person.CreatedAt,
			&person.UpdatedAt,
		); err != nil {
			slog.Error("GetPeople scan process Error")
			return nil, fmt.Errorf("error scanning person row: %v", err)
		}
		people = append(people, person)
	}

	slog.Debug("GetPeople found people")
	return people, nil
}

// GetPerson - метод для получения одного человека по id
func (s *personRepository) GetPerson(id int) (entities.Person, error) {
	query := `
		SELECT * FROM people
		WHERE id = $1
	`

	var person entities.Person

	err := s.db.QueryRowContext(context.Background(), query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		slog.Error("GetPerson Error")
		return entities.Person{}, fmt.Errorf("error getting person: %v", err)
	}

	slog.Debug("GetPerson found person with id:", id)
	return person, nil
}

// GetPeopleFiltered - метод для получения списка людей с применением пагинации и фильтрации
func (s *personRepository) GetPeopleFiltered(sort string, from int, to int, page int, perPage int) ([]entities.Person, error) {
	slog.Info("", sort, from, to)

	var query string
	var args []interface{}

	switch sort {
	case "age":
		if from < 0 || to == 0 {
			slog.Error("From or To field is empty")
			return nil, nil
		}
		query = `
		SELECT * FROM people
		WHERE age > $1 AND age < $2
		ORDER BY age
		LIMIT $3 OFFSET $4
	`
		args = append(args, from, to, perPage, (page-1)*perPage)
	case "id":
		if from == 0 || to == 0 {
			slog.Error("From or To field is empty")
			return nil, nil
		}
		query = `
		SELECT * FROM people
		WHERE id > $1 AND id < $2
		ORDER BY id
		LIMIT $3 OFFSET $4
	`
		args = append(args, from, to, perPage, (page-1)*perPage)
	case "nation":
		query = `
		SELECT * FROM people
		ORDER BY nationality
		LIMIT $1 OFFSET $2
	`
		args = append(args, perPage, (page-1)*perPage)
	case "gender":
		query = `
		SELECT * FROM people
		ORDER BY gender
		LIMIT $1 OFFSET $2
	`
		args = append(args, perPage, (page-1)*perPage)
	default:
		query = `
		SELECT * FROM people
		LIMIT $1 OFFSET $2
	`
		args = append(args, perPage, (page-1)*perPage)
	}

	slog.Info(query, args)
	rows, err := s.db.QueryContext(context.Background(), query, args...)
	slog.Info("", rows)
	if err != nil {
		slog.Error("GetPeopleFiltered QueryContext Error ")
		return nil, fmt.Errorf("error getting filtered people: %v", err)
	}
	defer rows.Close()

	var people []entities.Person

	for rows.Next() {
		var person entities.Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
			&person.CreatedAt,
			&person.UpdatedAt,
		); err != nil {
			slog.Error("GetPeopleFiltered scan process Error")
			return nil, fmt.Errorf("error scanning person row: %v", err)
		}
		people = append(people, person)
	}

	slog.Debug("GetPeopleFiltered found people")
	return people, nil
}

// UpdatePerson - метод для обновление данных о человеке по его id
func (s *personRepository) UpdatePerson(id int, person entities.Person) error {
	query := `
        UPDATE people
        SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, nationality = $7, updated_at = $8
        WHERE id = $1
    `
	_, err := s.db.Exec(query, id, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.UpdatedAt)
	if err != nil {
		slog.Error("UpdatePerson ExecQuery Error")
		return fmt.Errorf("error updating person: %v", err)
	}

	slog.Debug("UpdatePerson updated person with id:", id)
	return err
}

// DeletePerson - метод для удаления человека из базы по его id
func (s *personRepository) DeletePerson(id int) error {
	query := `
        DELETE FROM people
        WHERE id = $1
    `
	_, err := s.db.Exec(query, id)
	if err != nil {
		slog.Error("DeletePerson ExecQuery Error")
		return fmt.Errorf("error deleting person: %v", err)
	}

	slog.Debug("DeletePerson deleted person with id:", id)
	return nil
}
