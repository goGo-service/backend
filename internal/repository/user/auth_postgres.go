package user

import (
	"fmt"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal"
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
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, user.VkID)
	if err := row.Scan(&id); err != nil {
		fmt.Println("Error executing query:", err)
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByVkId(vkId int64) (*goGO.User, error) {
	var user goGO.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE vk_id = $1", usersTable)
	row := r.db.QueryRow(query, vkId)
	if err := row.Scan(&user.Id, &user.VkID, &user.FirstName, &user.LastName, &user.Username, &user.Email); err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	return &user, nil
}
