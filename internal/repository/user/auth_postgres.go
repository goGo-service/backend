package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) CreateUser(user models.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, username, email, vk_id) values ($1, $2, $3, $4, $5) RETURNING id", internal.UsersTable)
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, user.VkID)
	if err := row.Scan(&id); err != nil {
		fmt.Println("Error executing query:", err)
		return 0, err
	}

	return id, nil
}

func (r *Postgres) GetUserByVkId(vkId int64) (*models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE vk_id = $1", internal.UsersTable)
	row := r.db.QueryRow(query, vkId)
	if err := row.Scan(&user.Id, &user.VkID, &user.FirstName, &user.LastName, &user.Username, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		fmt.Println("Error executing query:", err)

		return nil, err
	}

	return &user, nil
}

func (r *Postgres) GetUserById(userId int) (*models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", internal.UsersTable)
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&user.Id, &user.VkID, &user.FirstName, &user.LastName, &user.Username, &user.Email); err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	return &user, nil
}
