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

func (u *RoomUseCase) GetRoom(roomId int, userId int) (*models.Room, error) {
	haveRight, err := u.services.Room.HaveAccess(userId, roomId)
	if err != nil {
		return nil, err
	}
	if !haveRight {
		return nil, nil //TODO: мб стоит ошибку добавить другую
	}

	room, err := u.services.Room.GetRoom(roomId)
	if err != nil {
		return nil, err
	}

	return room, nil
}
