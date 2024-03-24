package setting

import "grape/src/common/service"

type SettingService struct{}

func NewSettingService(s *service.CommonService) *SettingService {
	return &SettingService{}
}
