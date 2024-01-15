package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
	"github.com/GOLANG-NINJA/crud-app/internal/repository/psql"

	"github.com/gorilla/mux"
)

type User interface {
	SignUp(inp domain.SignUpInput) error
	SignIn(inp domain.SignInInput) (string, error)
}

type Handler struct {
	employee     *psql.Employees
	usersService User
	session      *psql.Sessions
}

func NewHandler(empls *psql.Employees, users User, session *psql.Sessions) *Handler {
	return &Handler{
		employee:     empls,
		usersService: users,
		session:      session,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/logout", h.logout).Methods(http.MethodGet)
	}

	books := r.PathPrefix("/employee").Subrouter()
	{
		books.Use(h.authMiddleware)

		books.HandleFunc("", h.createEmpl).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllEmpls).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getEmplByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteEmpl).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateEmpl).Methods(http.MethodPut)
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
