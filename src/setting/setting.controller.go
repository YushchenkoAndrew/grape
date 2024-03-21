package setting

type SettingController struct {
	service *SettingService
}

func NewSettingController(s *SettingService) *SettingController {
	return &SettingController{service: s}
}
