package src

import (
	"grape/src/common/service"
)

type indexService struct {
}

func NewIndexService(s *service.CommonService) *indexService {
	return &indexService{}
}
