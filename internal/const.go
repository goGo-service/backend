package internal

import (
	"errors"
)

const (
	vkIdUsersTable  = "vk_id_users"
	UsersTable      = "users"
	UserTokensTable = "user_tokens"
	roomsTable      = "rooms"
	rolesTable      = "roles"
	roomsUsersTable = "rooms_users"
)

var (
	AccessTokenRequiredError = errors.New("access token is required")
	InternalServiceError     = errors.New("internal service error")
	UserNotFoundError        = errors.New("user not found")
)
