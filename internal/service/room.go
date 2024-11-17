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
	return s.repo.FetchRoomById(id)
}

func (s *RoomService) AddOwnerToRoom(userId int, roomId int) error {
	return s.repo.SaveRoomUser(userId, roomId, models.RoomOwner)
}

func (s *RoomService) HaveAccess(userId int, roomId int) (bool, error) {
	roomUser, err := s.repo.FetchRoomUser(userId, roomId)
	if err != nil {
		return false, err
	}
	if roomUser == nil {
		return false, nil
	}

	return true, nil
}

func (s *RoomService) GetUserRooms(userID int) ([]*models.Room, error) {
	return s.repo.FetchRoomsByUserId(userID)
}
