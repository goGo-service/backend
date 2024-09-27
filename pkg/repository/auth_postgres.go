package repository

import (
	"fmt"

	goGO "github.com/goGo-service/back"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user goGO.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, username, email, vk_id) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	// TODO:  доставать vk_id из API
	row := r.db.QueryRow(query, user.Name, user.Surname, user.Username, user.Email, 1)
	if err := row.Scan(&id); err != nil {
		fmt.Println("Error executing query:", err)
		return 0, err
	}

	return id, nil
}
