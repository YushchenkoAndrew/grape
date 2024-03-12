package filter

import (
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/filter/dto/request"
	"grape/src/filter/dto/response"
	repo "grape/src/filter/repositories"
)

type FilterService struct {
	Repository *repo.LocationRepositoryT
}

func NewFilterService(client *service.CommonService) *FilterService {
	return &FilterService{Repository: repo.NewLocationRepository(client.DB)}
}

func (c *FilterService) TraceIP(dto *request.LocationDto) (*response.LocationResponseDto, error) {
	res, err := c.Repository.GetOne(dto)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[response.LocationResponseDto](res), nil
}
