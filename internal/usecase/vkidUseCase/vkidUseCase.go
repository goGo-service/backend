package authUseCase

import (
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

func (u VKIDUseCase) getUserIdAndAT(code string, deviceId string, state string) (int64, string, error) {
	resp, err := u.services.VKID.ExchangeCode(code, deviceId, state)
	if err != nil {
		return 0, "", err
	}
}
