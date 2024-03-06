package src

import (
	"grape/src/common/client"
)

type indexService struct {
}

func NewIndexService(client *client.Clients) *indexService {
	return &indexService{}
}
