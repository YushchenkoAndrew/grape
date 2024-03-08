package request

type LoginDto struct {
	User string `json:"user" xml:"user" binding:"required,not_empty"`
	Pass string `json:"pass" xml:"pass" binding:"required,not_empty"`
}
