package main

import (
	"log"
	"time"

	repos "github.com/Sskrill/Authentication.git/internal/repository"
	"github.com/Sskrill/Authentication.git/internal/service"
	"github.com/Sskrill/Authentication.git/internal/transport"
	"github.com/Sskrill/Authentication.git/pkg/hash"
	pqsql "github.com/Sskrill/Authentication.git/pkg/pqSQL"
)

func main() {
	db, err := pqsql.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	hasher := hash.NewPSHasher("salt")
	employees := repos.NewEmployees(db)
	users := repos.NewUsers(db)
	userService := service.NewUsersService(users, hasher, []byte("secret"), time.Minute*10)

	handler := transport.NewHandler(userService, employees)

	handler.InitRout()
}
