package goGO

type User struct {
	Id        int    `json:"-"`
	FirstName string `json:"first_name" bindind:"required"`
	VkID      int64  `json:"vk_id"`
	LastName  string `json:"last_name" bindind:"required"`
	Username  string `json:"username" bindind:"required"`
	Email     string `json:"email" bindind:"required"`
}
