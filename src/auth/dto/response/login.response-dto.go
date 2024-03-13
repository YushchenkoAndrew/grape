package response

import r "grape/src/user/dto/response"

type LoginResponseDto struct {
	AccessToken  string                    `json:"access_token" xml:"access_token" example:"temp"`
	RefreshToken string                    `json:"refresh_token" xml:"refresh_token" example:"temp"`
	User         r.UserAdvancedResponseDto `json:"user" xml:"user"`
}
