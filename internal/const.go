package internal

import (
	"errors"
)

const (
	VkIdUsersTable  = "vk_id_users"
	UsersTable      = "users"
	RoomsTable      = "rooms"
	RolesTable      = "roles"
	RoomsUsersTable = "rooms_users"
)

var (
	AccessTokenRequiredError = errors.New("access token is required")
	InternalServiceError     = errors.New("internal service error")
	UserNotFoundError        = errors.New("user not found")
	InvalidUserIDError       = errors.New("invalid user id")
	UserIDNotFoundError      = errors.New("user not found")
)
