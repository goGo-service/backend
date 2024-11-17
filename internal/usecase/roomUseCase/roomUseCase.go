package roomUseCase

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
)

type RoomUseCase struct {
	services *service.Service
}

func NewRoomUseCase(service *service.Service) *RoomUseCase {
	return &RoomUseCase{
		services: service,
	}
}

func (u *RoomUseCase) CreateNewRoom(room models.Room, userId int) (int, error) {
	roomId, err := u.services.Room.CreateRoom(room)
	if err != nil {
		return 0, err
	}
	err = u.services.Room.AddOwnerToRoom(userId, roomId)
	if err != nil {
		return 0, err
	}
	return roomId, nil
}
