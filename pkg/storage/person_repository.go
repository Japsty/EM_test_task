package storage

import (
	"EM_test_task/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/pingcap/log"
)

type Repository interface {
	CreatePerson(ctx context.Context, arg PersonParams) (entities.Person, error)
	GetPeople() ([]entities.Person, error)
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
}

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

	return person, nil
}

func (s *personRepository) GetPeople() ([]entities.Person, error) {
	query := `
		SELECT * FROM people
	`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		log.Error("GetPeople QueryContext Error")
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
		); err != nil {
			log.Error("GetPeople scan process Error")
			return nil, fmt.Errorf("error scanning person row: %v", err)
		}
		people = append(people, person)
	}

	return people, nil
}

func (s *personRepository) UpdatePerson(id int, person entities.Person) error {
	query := `
        UPDATE people
        SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, nationality = $7
        WHERE id = $1
    `
	_, err := s.db.Exec(query, id, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		log.Error("UpdatePerson ExecQuery Error")
		return fmt.Errorf("error updating person: %v", err)
	}
	return err
}

func (s *personRepository) DeletePerson(id int) error {
	query := `
        DELETE FROM people
        WHERE id = $1
    `
	_, err := s.db.Exec(query, id)
	if err != nil {
		log.Error("DeletePerson ExecQuery Error")
		return fmt.Errorf("error deleting person: %v", err)
	}
	return nil
}
