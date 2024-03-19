package customer

import (
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/customer/dto/request"
	"grape/src/customer/dto/response"
	repo "grape/src/customer/repositories"
)

type CustomerService struct {
	Repository *repo.LocationRepositoryT
}

func NewCustomerService(client *service.CommonService) *CustomerService {
	return &CustomerService{Repository: repo.NewLocationRepository(client.DB)}
}

func (c *CustomerService) TraceIP(dto *request.LocationDto) (*response.LocationResponseDto, error) {
	res, err := c.Repository.GetOne(dto)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[response.LocationResponseDto](res), nil
}
