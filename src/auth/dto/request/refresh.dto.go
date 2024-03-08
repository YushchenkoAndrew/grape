package request

type RefreshDto struct {
	RefreshToken string `json:"refresh_token" xml:"refresh_token" binding:"required"`
}
