package response

import r "grape/src/common/dto/response"

type UserResponseDto struct {
	r.UuidResponseDto
	Organization OrganizationResponseDto `json:"organization" xml:"organization"`
}
