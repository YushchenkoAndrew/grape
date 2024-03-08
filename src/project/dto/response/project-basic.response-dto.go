package response

import r "grape/src/common/dto/response"

// Name        string `gorm:"not null" xml:"name" example:"Code Rain"`
// Description string `json:"desc" xml:"desc" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
// Type        t.ProjectTypeEnum `json:"" xml:"flag" example:"js"`
// Footer      string `json:"note" xml:"note" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`

type ProjectBasicResponseDto struct {
	*r.UuidResponseDto

	Description string `json:"description" xml:"description" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Type        string `copier:"GetType" json:"type" xml:"type" example:"html"`
	Footer      string `json:"footer" xml:"footer" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
}

func NewProjectBasicResponseDto() *ProjectBasicResponseDto {
	return &ProjectBasicResponseDto{
		UuidResponseDto: &r.UuidResponseDto{},
	}
}
