package main

import (
	"fmt"
	_ "github.com/GOLANG-NINJA/crud-app/docs"
	"net/http"
	"os"

	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"
	"github.com/GOLANG-NINJA/crud-app/internal/service"
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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(1)
	// init deps
	hasher := hash.NewSHA1Hasher("salt")

	employees := psql.NewEmpls(db)

	usersRepo := psql.NewUsers(db)
	tokensRepo := psql.NewTokens(db)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte("secret"))

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
