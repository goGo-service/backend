package vkidUseCase

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
)

type VKIDUseCase struct {
	services *service.Service
}

func NewVKIDUseCase(service *service.Service) *VKIDUseCase {
	return &VKIDUseCase{
		services: service,
	}
}

func (u *VKIDUseCase) GetUserIdAndAT(code string, deviceId string, state string) (int64, string, error) {
	resp, err := u.services.VKID.ExchangeCode(code, deviceId, state)
	if err != nil {
		return 0, "", err
	}

	return resp.UserId, resp.AccessToken, nil
}

func (u *VKIDUseCase) GetUserInfo(accessToken string, code string) (*models.VKIDUserInfo, error) {
	info, err := u.services.VKID.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	userInfo := &models.VKIDUserInfo{
		ID:        info.Response[0].ID,
		FirstName: info.Response[0].FirstName,
		LastName:  info.Response[0].LastName,
		Email:     info.Response[0].Email,
	}
	err = u.services.VKID.CacheVKID(code, userInfo.ID)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (u *VKIDUseCase) GetVKID(code string) (int64, error) {
	vkid, err := u.services.VKID.GetCachedVKID(code)
	if err != nil {
		return 0, err
	}
	return vkid, nil
}

func (u *VKIDUseCase) DeleteVKID(code string) error {
	return u.services.VKID.DeleteCachedVKID(code)
}
