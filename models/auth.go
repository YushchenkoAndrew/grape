package models

type Auth struct {
	AccessToken   string
	RefreshToken  string
	AccessUUID    string
	RefreshUUID   string
	AccessExpire  int64
	RefreshExpire int64
}

func NewAuth() *Auth {
	return &Auth{}
}

func (c *Auth) Fill(auth map[string]interface{}) *Auth {
	if val, ok := auth["access_token"]; ok && val != "" {
		c.AccessToken = val.(string)
	}

	if val, ok := auth["refresh_token"]; ok && val != "" {
		c.RefreshToken = val.(string)
	}

	if val, ok := auth["access_uuid"]; ok && val != "" {
		c.AccessUUID = val.(string)
	}

	if val, ok := auth["refresh_uuid"]; ok && val != "" {
		c.RefreshUUID = val.(string)
	}

	if val, ok := auth["access_expire"]; ok && val != 0 {
		c.AccessExpire = val.(int64)
	}

	if val, ok := auth["refresh_expire"]; ok && val != "" {
		c.RefreshExpire = val.(int64)
	}

	return c
}

type LoginDto struct {
	User string `json:"user" xml:"user" binding:"required"`
	Pass string `json:"pass" xml:"pass" binding:"required"`
}

func (c *LoginDto) IsOK() bool {
	return c.User != "" || c.Pass != ""
}

type TokenEntity struct {
	Status       string `json:"status" xml:"status" example:"OK"`
	AccessToken  string `json:"access_token" xml:"access_token"`
	RefreshToken string `json:"refresh_token" xml:"refresh_token"`
}

type TokenDto struct {
	AccessToken  string `json:"access_token" xml:"access_token"`
	RefreshToken string `json:"refresh_token" xml:"refresh_token" binding:"required"`
}
