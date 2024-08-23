package main

import (
	"fmt"
	redisC "github.com/GOLANG-NINJA/crud-app/internal/cache/redis"
	"net/http"
	"os"

	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"
	"github.com/GOLANG-NINJA/crud-app/internal/service"
	"github.com/GOLANG-NINJA/crud-app/internal/transport/gRPC"
	"github.com/GOLANG-NINJA/crud-app/internal/transport/rest"
	"github.com/GOLANG-NINJA/crud-app/pkg/database"
	"github.com/GOLANG-NINJA/crud-app/pkg/hash"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	db, err := database.NewDB()
	if err != nil {

		log.Fatal(err)
	}
	// init deps
	hasher := hash.NewSHA1Hasher("salt")
	auditClient, err := gRPC.NewAuditClient(9000)
	if err != nil {
		log.Fatal(err)
	}
	employees := psql.NewEmpls(db, auditClient)
	cacheRedis := redisC.NewRedisCache("localhost:6379", "")
	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)

	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, auditClient, []byte("secret"))

	handler := rest.NewHandler(employees, usersService, cacheRedis)

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
