package request

type StageCreateDto struct {
	Name string `json:"name" xml:"name" binding:"required"`
}
