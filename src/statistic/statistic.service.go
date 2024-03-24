package statistic

import "grape/src/common/service"

type StatisticService struct {
}

func NewStatisticService(s *service.CommonService) *StatisticService {
	return &StatisticService{}
}
