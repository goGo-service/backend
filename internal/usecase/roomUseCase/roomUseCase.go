package roomUseCase

import (
	"fmt"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
	"github.com/sirupsen/logrus"
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
	roomId, err := u.services.Room.CreateRoom(userId, room)
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

func (u *RoomUseCase) GetUserRooms(userId int) ([]*models.Room, error) {
	return u.services.Room.GetUserRooms(userId)
}

func (u *RoomUseCase) GetRoomPresence(userId int, roomId int) ([]*models.RoomUserPresence, error) {
	roomUserIds, err := u.services.Room.GetRoomPresence(roomId)
	if err != nil {
		return nil, err
	}
	var users []*models.RoomUserPresence
	for _, roomUserId := range roomUserIds {
		if userId == roomUserId { //Не стоит отдавать на фронт присутствие самого клиента
			//continue
		}
		user, err := u.services.User.GetUser(roomUserId)
		if err != nil {
			return nil, err
		}
		users = append(users, &models.RoomUserPresence{
			UserId:   user.Id,
			Username: user.Username,
		})
	}
	users = append(users, &models.RoomUserPresence{
		UserId:   1,
		Username: "Pavel Durov",
	})
	return users, nil
}

func (u *RoomUseCase) PublishMessage(userId int, roomId int, message string) error {
	haveRight, err := u.services.Room.HaveAccess(userId, roomId)
	if err != nil {
		logrus.Errorf("Error services.Room.HaveAccess %s", err.Error())
		return err
	}
	if !haveRight {
		return fmt.Errorf("dont have right")
	}

	err = u.services.Room.PublishMessage(userId, roomId, message)
	if err != nil {
		logrus.Errorf("Error  Room.PublishMessage %s", err.Error())
		return fmt.Errorf("can`t publish message")
	}

	return nil
}
