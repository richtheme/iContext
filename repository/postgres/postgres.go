package postgres

import (
	"database/sql"
	"fmt"
	"iContext/internal/models"
	"log"

	_ "github.com/lib/pq"
)

const (
	usersTable = "users"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURI)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db
	log.Println("Database connection created successfully!")
	return nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (storage *Storage) CreateUser(u *models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, age) values ($1, $2) RETURNING id", usersTable)

	row := storage.db.QueryRow(query, u.Name, u.Age)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// for tests only (purpose: to empty database)
func (storage *Storage) Exec(a string) error {
	_, err := storage.db.Exec(a)
	if err != nil {
		return err
	}
	return nil
}
