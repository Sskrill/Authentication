package rest

import (
	"context"
	"errors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"strconv"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"

	"github.com/gorilla/mux"
)

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	employee     *psql.Employees
	usersService User
}

func NewHandler(empls *psql.Employees, users User) *Handler {
	return &Handler{
		employee:     empls,
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
	}

	employee := r.PathPrefix("/employee").Subrouter()
	{
		employee.Use(h.authMiddleware)

		employee.HandleFunc("", h.createEmpl).Methods(http.MethodPost)
		employee.HandleFunc("", h.getAllEmpls).Methods(http.MethodGet)
		employee.HandleFunc("/{id:[0-9]+}", h.getEmplByID).Methods(http.MethodGet)
		employee.HandleFunc("/{id:[0-9]+}", h.deleteEmpl).Methods(http.MethodDelete)
		employee.HandleFunc("/{id:[0-9]+}", h.updateEmpl).Methods(http.MethodPut)
	}

	return r
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
