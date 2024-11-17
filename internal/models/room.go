package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Room struct {
	Id        int          `db:"id" json:"id"`
	Name      string       `db:"name" json:"name"`
	Settings  RoomSettings `json:"settings" db:"settings"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
}

type RoomResponse struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Settings RoomSettings `json:"settings"`
}

func (room *Room) ToResponse() *RoomResponse {
	return &RoomResponse{
		Id:       room.Id,
		Name:     room.Name,
		Settings: room.Settings,
	}
}

type RoomSettings struct {
	Capacity int `json:"capacity"`
}

// Value реализует интерфейс `driver.Valuer` для преобразования RoomSettings в JSON при добавлении в базу данных.
func (rs *RoomSettings) Value() (driver.Value, error) {
	return json.Marshal(rs)
}

// Scan реализует интерфейс `sql.Scanner` для преобразования JSON из базы данных в структуру RoomSettings.
func (rs *RoomSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to []byte")
	}

	return json.Unmarshal(bytes, rs)
}

type RoomUser struct {
	RoomId    int       `db:"room_id" json:"room_id"`
	UserId    int       `db:"user_id" json:"user_id"`
	RoleId    int       `db:"role_id" json:"role_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
