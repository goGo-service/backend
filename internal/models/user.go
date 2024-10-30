package models

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name" bindind:"required"`
	VkID      int64  `json:"vk_id"`
	LastName  string `json:"last_name" bindind:"required"`
	Username  string `json:"username" bindind:"required"`
	Email     string `json:"email" bindind:"required"`
}

type VKIDUserInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
