package gateway

import "github.com/dionerweiss/fc-ms-wallet/goapp/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
