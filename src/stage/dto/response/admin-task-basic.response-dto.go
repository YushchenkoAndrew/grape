package response

type AdminTaskBasicResponseDto struct {
	TaskBasicResponseDto

	Order  int    `json:"order" xml:"order" example:"0"`
	Status string `copier:"GetStatus" json:"status" xml:"status" example:"true"`
}
