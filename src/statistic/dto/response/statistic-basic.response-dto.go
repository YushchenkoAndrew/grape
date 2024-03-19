package response

type StatisticBasicResponseDto struct {
	Views  int `json:"views" xml:"views" example:"0"`
	Clicks int `json:"clicks" xml:"clicks" example:"0"`
	Media  int `json:"media" xml:"media" example:"0"`
}
