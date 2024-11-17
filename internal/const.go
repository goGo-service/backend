package internal

import (
	"errors"
)

const (
	vkIdUsersTable  = "vk_id_users"
	UsersTable      = "users"
	roomsTable      = "rooms"
	rolesTable      = "roles"
	RoomsUsersTable = "rooms_users"
)

var (
	AccessTokenRequiredError = errors.New("access token is required")
	InternalServiceError     = errors.New("internal service error")
	UserNotFoundError        = errors.New("user not found")
	InvalidUserIDError       = errors.New("invalid user id")
	UserIDNotFoundError      = errors.New("user not found")
)
