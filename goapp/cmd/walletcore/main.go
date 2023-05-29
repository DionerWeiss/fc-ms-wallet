package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/database"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/event"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/event/handler"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/usecase/create_account"
	create_client "github.com/dionerweiss/fc-ms-wallet/goapp/internal/usecase/create_client"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/usecase/create_transaction"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/web"
	"github.com/dionerweiss/fc-ms-wallet/goapp/internal/web/webserver"
	"github.com/dionerweiss/fc-ms-wallet/goapp/pkg/events"
	"github.com/dionerweiss/fc-ms-wallet/goapp/pkg/kafka"
	"github.com/dionerweiss/fc-ms-wallet/goapp/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}