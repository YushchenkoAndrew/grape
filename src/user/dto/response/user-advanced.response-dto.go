package response

type UserAdvancedResponseDto struct {
	UserBasicResponseDto
	Organization OrganizationResponseDto `json:"organization" xml:"organization"`
}
