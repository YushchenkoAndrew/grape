package deployment

// import (
// 	m "grape/models"
// 	"context"

// 	v1 "k8s.io/api/apps/v1"
// 	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// type DeploymentService struct {
// 	k3s *kubernetes.Clientset
// }

// func NewDeploymentService(k3s *kubernetes.Clientset) *DeploymentService {
// 	return &DeploymentService{k3s}
// }

// func (c *DeploymentService) Create(namespace string, dto *v1.Deployment) (*v1.Deployment, error) {
// 	return c.k3s.AppsV1().Deployments(namespace).Create(context.Background(), dto, metaV1.CreateOptions{})
// }

// func (c *DeploymentService) Read(namespace string, query *m.K3sListQueryDto) (*v1.DeploymentList, error) {
// 	return c.k3s.AppsV1().Deployments(namespace).List(context.Background(), metaV1.ListOptions{
// 		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
// 		LabelSelector:       query.LabelSelector,
// 		FieldSelector:       query.FieldSelector,
// 		Watch:               query.Watch,
// 		AllowWatchBookmarks: query.AllowWatchBookmarks,
// 		ResourceVersion:     query.ResourceVersion,
// 	})
// }

// func (c *DeploymentService) Update(namespace string, query *m.K3sUpdateQueryDto, dto *v1.Deployment) (*v1.Deployment, error) {
// 	return c.k3s.AppsV1().Deployments(namespace).Update(context.Background(), dto, metaV1.UpdateOptions{
// 		TypeMeta:        metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
// 		FieldManager:    query.FieldManager,
// 		FieldValidation: query.FieldValidation,
// 	})
// }

// func (c *DeploymentService) Delete(namespace, name string) error {
// 	return c.k3s.AppsV1().Deployments(namespace).Delete(context.Background(), name, metaV1.DeleteOptions{})
// }
