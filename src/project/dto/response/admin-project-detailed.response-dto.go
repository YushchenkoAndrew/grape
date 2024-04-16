package response

import (
	att "grape/src/attachment/dto/response"
	ln "grape/src/link/dto/response"
	st "grape/src/statistic/dto/response"
)

type AdminProjectDetailedResponseDto struct {
	AdminProjectBasicResponseDto

	Footer      string                              `json:"footer" xml:"footer" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
	Statistic   st.StatisticBasicResponseDto        `json:"statistic" xml:"statistic"`
	Attachments []att.AttachmentAdvancedResponseDto `json:"attachments" xml:"attachments"`
	Links       []ln.LinkBasicResponseDto           `json:"links" xml:"links"`
}
