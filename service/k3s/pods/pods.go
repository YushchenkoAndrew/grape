package pods

import (
	m "api/models"
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type PodsService struct {
	metrics *metrics.Clientset
}

func NewPodsService(metrics *metrics.Clientset) *PodsService {
	return &PodsService{metrics}
}

func (c *PodsService) Read(namespace string, query *m.K3sListQueryDto) (*v1beta1.PodMetricsList, error) {
	return c.metrics.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metaV1.ListOptions{
		TypeMeta:            metaV1.TypeMeta{Kind: query.Kind, APIVersion: query.APIVersion},
		LabelSelector:       query.LabelSelector,
		FieldSelector:       query.FieldSelector,
		Watch:               query.Watch,
		AllowWatchBookmarks: query.AllowWatchBookmarks,
		ResourceVersion:     query.ResourceVersion,
	})
}
