package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ProjectDto struct {
	request.PageDto

	Query string `form:"query,omitempty" example:"test"`

	ProjectIds []string

	// 	CreatedTo   time.Time `form:"created_to,omitempty" time_format:"2006-01-02" example:"2021-08-06"`
	// 	CreatedFrom time.Time `form:"created_from,omitempty" time_format:"2006-01-02" example:"2021-08-06"`

	// 	Page  int `form:"page,omitempty,default=-1" example:"1"`
	// 	Limit int `form:"limit,omitempty" example:"10"`

	// 	Link struct {
	// 		ID   uint32 `form:"link[id],omitempty"`
	// 		Name string `form:"link[name],omitempty" example:"main"`

	// 		Page  int `form:"link[page],omitempty,default=-1" example:"1"`
	// 		Limit int `form:"link[limit],omitempty" example:"10"`
	// 	}

	// 	File struct {
	// 		ID   uint32 `form:"file[id],omitempty"`
	// 		Name string `form:"file[name],omitempty" example:"main"`
	// 		Role string `form:"file[role],omitempty" example:"src"`
	// 		Path string `form:"file[path],omitempty" example:"/test"`

	// 		Page  int `form:"file[page],omitempty,default=-1" example:"1"`
	// 		Limit int `form:"file[limit],omitempty" example:"10"`
	// 	}

	// 	Subscription struct {
	// 		ID     uint32 `form:"subscription[id],omitempty"`
	// 		Name   string `form:"subscription[name],omitempty" example:"main"`
	// 		CronID string `form:"subscription[cron_id],omitempty" example:"main"`

	// 		Page  int `form:"subscription[page],omitempty,default=-1" example:"1"`
	// 		Limit int `form:"subscription[limit],omitempty" example:"10"`
	// 	}

	// 	Metrics struct {
	// 		ID            uint32 `form:"metrics[id],omitempty"`
	// 		Name          string `form:"metrics[name],omitempty" example:"main"`
	// 		Namespace     string `form:"metrics[namespace],omitempty" example:"void-deployment-8985bd57d-k9n5g"`
	// 		ContainerName string `form:"metrics[container_name],omitempty" example:"void"`

	// 		CreatedTo   time.Time `form:"metrics[created_to],omitempty" time_format:"2006-01-02" example:"2021-08-06"`
	// 		CreatedFrom time.Time `form:"metrics[created_from],omitempty" time_format:"2006-01-02" example:"2021-08-06"`

	//		Page  int `form:"metrics[page],omitempty,default=-1" example:"1"`
	//		Limit int `form:"metrics[limit],omitempty" example:"10"`
	//	}
}

func NewProjectDto(user *entities.UserEntity) *ProjectDto {
	return &ProjectDto{PageDto: request.NewPageDto(user)}
}
