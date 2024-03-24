package request

type LoginDto struct {
	Username string `json:"username" xml:"username" binding:"required" example:"test"`
	Password string `json:"password" xml:"password" binding:"required" example:"test"`
}
