package models

type RedirectUrl struct {
	AppId         int    `json:"app_id"`
	RedirectUrl   string `json:"redirect_url"`
	State         string `json:"state"`
	CodeChallenge string `json:"code_challenge"`
}
