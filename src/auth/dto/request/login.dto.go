package request

type LoginDto struct {
	Name string `json:"name" xml:"name" binding:"required"`
	Pass string `json:"pass" xml:"pass" binding:"required"`
}
