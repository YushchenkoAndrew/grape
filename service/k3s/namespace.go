package k3s

import (
	m "api/models"
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	k3s *kubernetes.Clientset
}

func NewNamespaceService(k3s *kubernetes.Clientset) *NamespaceService {
	return &NamespaceService{k3s}
}

func (c *NamespaceService) Create(dto *v1.Namespace) (*v1.Namespace, error) {
	return c.k3s.CoreV1().Namespaces().Create(context.Background(), dto, metaV1.CreateOptions{})
}

func (c *NamespaceService) Read(query *m.K3sListQueryDto) (*v1.NamespaceList, error) {
	return c.k3s.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{
		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		LabelSelector:       query.LabelSelector,
		FieldSelector:       query.FieldSelector,
		Watch:               query.Watch,
		AllowWatchBookmarks: query.AllowWatchBookmarks,
		ResourceVersion:     query.ResourceVersion,
	})
}

func (c *NamespaceService) Update(query *m.K3sUpdateQueryDto, dto *v1.Namespace) (*v1.Namespace, error) {
	return c.k3s.CoreV1().Namespaces().Update(context.Background(), dto, metaV1.UpdateOptions{
		TypeMeta:        metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		FieldManager:    query.FieldManager,
		FieldValidation: query.FieldValidation,
	})
}

func (c *NamespaceService) Delete(name string) error {
	return c.k3s.CoreV1().Namespaces().Delete(context.Background(), name, metaV1.DeleteOptions{})
}
