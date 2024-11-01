package vkidUseCase

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
	"github.com/spf13/viper"
	"strconv"
)

type VKIDUseCase struct {
	services *service.Service
}

func NewVKIDUseCase(service *service.Service) *VKIDUseCase {
	return &VKIDUseCase{
		services: service,
	}
}

func (u *VKIDUseCase) GetRedirectUrl() (*models.RedirectUrl, error) {
	url := viper.GetString("VKID_REDIRECT_URL")
	appId := viper.GetInt("VKID_APP_ID")

	state, codeChallenge, err := u.services.VKID.GenerateStateAndCodeChallenge()
	if err != nil {
		return nil, err
	}

	redirectUrl := &models.RedirectUrl{
		AppId:         appId,
		RedirectUrl:   url,
		State:         state,
		CodeChallenge: codeChallenge,
	}

	return redirectUrl, nil
}

func (u *VKIDUseCase) GetUserIdAndAT(code string, deviceId string, state string) (int64, string, error) {
	resp, err := u.services.VKID.ExchangeCode(code, deviceId, state)
	if err != nil {
		return 0, "", err
	}

	return resp.UserId, resp.AccessToken, nil
}

func (u *VKIDUseCase) GetUserInfo(accessToken string, code string) (*models.VKIDUserInfo, error) {
	userInfo, err := u.services.VKID.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	vkid, err := strconv.ParseInt(userInfo.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	err = u.services.VKID.CacheVKID(code, vkid)
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
