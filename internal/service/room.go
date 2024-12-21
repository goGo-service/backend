package service

import (
	"fmt"
	"github.com/goGo-service/back/internal/adapter"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"strconv"
)

type RoomService struct {
	repo             repository.Room
	centrifugoClient *adapter.Client
}

func NewRoomService(repo repository.Room, centrifugoClient *adapter.Client) *RoomService {
	return &RoomService{repo: repo, centrifugoClient: centrifugoClient}
}

func (s *RoomService) CreateRoom(userId int, room models.Room) (int, error) {
	rooms, err := s.repo.FetchRoomsByUserId(userId)
	if err != nil {
		return 0, err
	}

	if len(rooms) > 4 {
		return 0, fmt.Errorf("to much rooms, %d", len(rooms))
	}

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

	return roomUser != nil, nil
}

func (s *RoomService) GetUserRooms(userID int) ([]*models.Room, error) {
	return s.repo.FetchRoomsByUserId(userID)
}

func (s *RoomService) GetRoomPresence(roomId int) ([]int, error) {
	return s.centrifugoClient.GetPresence(getChannelName(roomId))
}

func getChannelName(roomId int) string {
	return "room:" + strconv.Itoa(roomId)
}

type Message struct {
	Type string `json:"type"`
	User int    `json:"user"`
}

func (s *RoomService) PublishMessage(userId int, roomId int, message string) error {
	msg := &Message{Type: message, User: userId}
	return s.centrifugoClient.PublishMessage(getChannelName(roomId), msg)
}
