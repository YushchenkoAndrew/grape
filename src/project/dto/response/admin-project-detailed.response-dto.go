package response

import (
	att "grape/src/attachment/dto/response"
	st "grape/src/statistic/dto/response"
)

type AdminProjectDetailedResponseDto struct {
	AdminProjectBasicResponseDto

	Footer      string                           `json:"footer" xml:"footer" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
	Attachments []att.AttachmentBasicResponseDto `json:"attachments" xml:"attachments"`
	Statistic   st.StatisticBasicResponseDto     `json:"statistic" xml:"statistic"`
}
