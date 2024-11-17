package models

import "time"

type Role struct {
	Id        int       `db:"id" json:"id"`
	RoleName  string    `db:"role_name" json:"role_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

const RoomOwner = 1
const RoomMember = 2
