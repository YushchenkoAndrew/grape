package k3s

import (
	m "api/models"
	"context"

	v1beta1 "k8s.io/api/networking/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IngressService struct {
	k3s *kubernetes.Clientset
}

func NewIngressService(k3s *kubernetes.Clientset) *IngressService {
	return &IngressService{k3s}
}

func (c *IngressService) Create(namespace string, dto *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return c.k3s.NetworkingV1beta1().Ingresses(namespace).Create(context.Background(), dto, metaV1.CreateOptions{})
}

func (c *IngressService) Read(namespace string, query *m.K3sListQueryDto) (*v1beta1.IngressList, error) {
	return c.k3s.NetworkingV1beta1().Ingresses(namespace).List(context.Background(), metaV1.ListOptions{
		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		LabelSelector:       query.LabelSelector,
		FieldSelector:       query.FieldSelector,
		Watch:               query.Watch,
		AllowWatchBookmarks: query.AllowWatchBookmarks,
		ResourceVersion:     query.ResourceVersion,
	})
}

func (c *IngressService) Update(namespace string, query *m.K3sUpdateQueryDto, dto *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return c.k3s.NetworkingV1beta1().Ingresses(namespace).Update(context.Background(), dto, metaV1.UpdateOptions{
		TypeMeta:        metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		FieldManager:    query.FieldManager,
		FieldValidation: query.FieldValidation,
	})
}

func (c *IngressService) Delete(namespace, name string) error {
	return c.k3s.NetworkingV1beta1().Ingresses(namespace).Delete(context.Background(), name, metaV1.DeleteOptions{})
}
