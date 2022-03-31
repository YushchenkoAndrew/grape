package k3s

import (
	"api/config"
	m "api/models"
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentService struct {
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

func (c *DeploymentService) Create(namespace string, dto *v1.Deployment) (*v1.Deployment, error) {
	return config.K3s.AppsV1().Deployments(namespace).Create(context.Background(), dto, metaV1.CreateOptions{})
}

func (c *DeploymentService) Read(namespace string, query *m.K3sQueryDto) (*v1.DeploymentList, error) {
	return config.K3s.AppsV1().Deployments(namespace).List(context.Background(), metaV1.ListOptions{
		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		LabelSelector:       query.LabelSelector,
		FieldSelector:       query.FieldSelector,
		Watch:               query.Watch,
		AllowWatchBookmarks: query.AllowWatchBookmarks,
		ResourceVersion:     query.ResourceVersion,
	})
}

func (c *DeploymentService) Update(query *m.CronQueryDto, model *m.CronDto) ([]m.CronEntity, error) {
	return nil, fmt.Errorf("Not implimented")
}

func (c *DeploymentService) Delete(query *m.CronQueryDto) (int, error) {
	return -1, fmt.Errorf("Not implimented")
}
