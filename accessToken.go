package goGO

type AccessToken struct {
	Id          int    `json:"-"`
	AccessToken string `json:"access_token" bindind:"required"`
}
