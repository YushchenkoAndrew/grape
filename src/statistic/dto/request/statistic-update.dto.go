package request

type StatisticUpdateDto struct {
	Kind string `json:"kind" xml:"kind" binding:"required,oneof=view click media"`
}
