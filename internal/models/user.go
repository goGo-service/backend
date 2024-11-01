package models

import "time"

type User struct {
	Id        int       `db:"id" json:"id"`
	VkID      int64     `db:"vk_id" json:"vk_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type VKIDUserInfo struct {
	Id         string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Sex        int    `json:"sex"`
	IsVerified bool   `json:"is_verified"`
	Birthday   string `json:"birthday"`
}
