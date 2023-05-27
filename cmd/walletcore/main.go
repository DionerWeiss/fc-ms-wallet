package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dionerweiss/fc-ms-wallet/internal/database"
	"github.com/dionerweiss/fc-ms-wallet/internal/event"
	"github.com/dionerweiss/fc-ms-wallet/internal/event/handler"
	"github.com/dionerweiss/fc-ms-wallet/internal/usecase/create_account"
	create_client "github.com/dionerweiss/fc-ms-wallet/internal/usecase/create_client"
	"github.com/dionerweiss/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/dionerweiss/fc-ms-wallet/internal/web"
	"github.com/dionerweiss/fc-ms-wallet/internal/web/webserver"
	"github.com/dionerweiss/fc-ms-wallet/pkg/events"
	"github.com/dionerweiss/fc-ms-wallet/pkg/kafka"
	"github.com/dionerweiss/fc-ms-wallet/pkg/uow"
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
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	// transactionDb := database.NewTransactionDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	transactionCreatedEvent := event.NewTransactionCreated()

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUSeCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandler(*createTransactionUSeCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()
}
