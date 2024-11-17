package service

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
)

type RoomService struct {
	repo repository.Room
}

func NewRoomService(repo repository.Room) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(room models.Room) (int, error) {
	return s.repo.SaveRoom(room)
}

func (s *RoomService) GetRoom(id int) (*models.Room, error) {
	//TODO: implement me
	panic("implement me")
}

func (s *RoomService) AddOwnerToRoom(userId int, roomId int) error {
	return s.repo.SaveRoomUser(userId, roomId, models.RoomOwner)
}
