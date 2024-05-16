package request

import (
	"grape/src/common/dto/request"
)

type TaskUpdateOrderDto struct {
	request.OrderUpdateDto

	StageID string `json:"stage_id" xml:"stage_id" binding:"omitempty,uuid4"`
}
