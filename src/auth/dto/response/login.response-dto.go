package response

type LoginResponseDto struct {
	AccessToken  string `json:"access_token" xml:"access_token" example:"temp"`
	RefreshToken string `json:"refresh_token" xml:"refresh_token" example:"temp"`
	User         string `json:"user" xml:"user" example:"Not Found"`
}
