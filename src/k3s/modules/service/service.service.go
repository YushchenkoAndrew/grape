package service

// import (
// 	m "grape/models"
// 	"context"

// 	v1 "k8s.io/api/core/v1"
// 	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// type ServiceService struct {
// 	k3s *kubernetes.Clientset
// }

// func NewServiceService(k3s *kubernetes.Clientset) *ServiceService {
// 	return &ServiceService{k3s}
// }

// func (c *ServiceService) Create(namespace string, dto *v1.Service) (*v1.Service, error) {
// 	return c.k3s.CoreV1().Services(namespace).Create(context.Background(), dto, metaV1.CreateOptions{})
// }

// func (c *ServiceService) Read(namespace string, query *m.K3sListQueryDto) (*v1.ServiceList, error) {
// 	return c.k3s.CoreV1().Services(namespace).List(context.Background(), metaV1.ListOptions{
// 		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
// 		LabelSelector:       query.LabelSelector,
// 		FieldSelector:       query.FieldSelector,
// 		Watch:               query.Watch,
// 		AllowWatchBookmarks: query.AllowWatchBookmarks,
// 		ResourceVersion:     query.ResourceVersion,
// 	})
// }

// func (c *ServiceService) Update(namespace string, query *m.K3sUpdateQueryDto, dto *v1.Service) (*v1.Service, error) {
// 	return c.k3s.CoreV1().Services(namespace).Update(context.Background(), dto, metaV1.UpdateOptions{
// 		TypeMeta:        metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
// 		FieldManager:    query.FieldManager,
// 		FieldValidation: query.FieldValidation,
// 	})
// }

// func (c *ServiceService) Delete(namespace, name string) error {
// 	return c.k3s.CoreV1().Services(namespace).Delete(context.Background(), name, metaV1.DeleteOptions{})
// }
