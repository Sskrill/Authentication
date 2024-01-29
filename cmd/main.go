package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"
	"github.com/GOLANG-NINJA/crud-app/internal/service"
	"github.com/GOLANG-NINJA/crud-app/internal/transport/rest"
	cfgMQ "github.com/GOLANG-NINJA/crud-app/pkg/NewMQ"
	"github.com/GOLANG-NINJA/crud-app/pkg/database"
	"github.com/GOLANG-NINJA/crud-app/pkg/hash"
	"github.com/streadway/amqp"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfgMQ, err := cfgMQ.NewCfgMQ()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := amqp.Dial(cfgMQ.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	que, err := ch.QueueDeclare("log", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	// init deps
	hasher := hash.NewSHA1Hasher("salt")

	employees := psql.NewEmpls(db)

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte("secret"), ch, que.Name)

	handler := rest.NewHandler(employees, usersService)

	// init & run server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: handler.InitRouter(),
	}

	log.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
