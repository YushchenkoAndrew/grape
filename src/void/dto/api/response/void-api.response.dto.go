package response

type VoidApiResponseDto struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}
