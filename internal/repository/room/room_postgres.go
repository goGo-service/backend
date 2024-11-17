package room

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	db *sqlx.DB
}

func NewRoomPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (r *Postgres) SaveRoom(room models.Room) (int, error) {
	query := `INSERT INTO rooms (settings, name) VALUES ($1, $2) RETURNING id`

	// Преобразуем settings в driver.Value (JSON строку)
	settingsValue, err := room.Settings.Value()
	if err != nil {
		logrus.Error("failed to marshal settings: ", err)
		return 0, err
	}

	err = r.db.QueryRow(query, settingsValue, room.Name).Scan(&room.Id)
	if err != nil {
		logrus.Error("failed to execute query: ", err)
		return 0, err
	}

	return room.Id, nil
}

func (r *Postgres) FetchRoomById(id int) (*models.Room, error) {
	var fetchedRoom models.Room
	query := `SELECT id, settings, name, created_at, updated_at  FROM rooms WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&fetchedRoom.Id, &fetchedRoom.Settings, &fetchedRoom.Name, &fetchedRoom.CreatedAt, &fetchedRoom.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logrus.Error(err)
		return nil, err
	}
	return &fetchedRoom, nil
}

func (r *Postgres) SaveRoomUser(userId int, roomId int, roleId int) error {
	query := `
		INSERT INTO rooms_users (user_id, room_id, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, room_id) DO UPDATE 
		SET role_id = $3`

	_, err := r.db.Exec(query, userId, roomId, roleId)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}

	return nil
}

func (r *Postgres) FetchRoomUser(userId int, roomId int) (*models.RoomUser, error) {
	query := `
		SELECT user_id, room_id, role_id, created_at, updated_at
		FROM rooms_users
		WHERE user_id = $1 AND room_id = $2
	`

	var roomUser models.RoomUser
	err := r.db.QueryRow(query, userId, roomId).Scan(
		&roomUser.UserId,
		&roomUser.RoomId,
		&roomUser.RoleId,
		&roomUser.CreatedAt,
		&roomUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Если записи нет, возвращаем nil вместо ошибки
		}
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	return &roomUser, nil
}

func (r *Postgres) FetchRoomsByUserId(userId int) ([]*models.Room, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.name, r.settings, r.created_at, r.updated_at
		FROM %s r
		INNER JOIN %s ru ON ru.room_id = r.id
		WHERE ru.user_id = $1
	`, internal.RoomsTable, internal.RoomsUsersTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room

	for rows.Next() {
		var room models.Room
		var roomSettings models.RoomSettings
		err := rows.Scan(
			&room.Id,
			&room.Name,
			&roomSettings,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		room.Settings = roomSettings
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, nil
	}

	return rooms, nil
}
