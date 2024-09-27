package goGO

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name" bindind:"required"`
	Surname  string `json:"surname" bindind:"required"`
	Username string `json:"username" bindind:"required"`
	Email    string `json:"email" bindind:"required"`
}
