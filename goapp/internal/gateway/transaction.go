package gateway

import "github.com/dionerweiss/fc-ms-wallet/goapp/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
