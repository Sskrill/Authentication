package repos

import (
	"database/sql"

	"github.com/Sskrill/Authentication.git/internal/domain"
)

type Users struct {
	DB *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{DB: db}
}
func (u *Users) Get(password string, login string) (domain.User, error) {
	var user domain.User
	err := u.DB.QueryRow("SELECT id,login,email,password FROM users WHERE password=$1 AND login=$2", password, login).Scan(&user.Id, &user.Login, &user.Email, &user.Password)

	return user, err
}
func (u *Users) Create(user domain.User) error {
	_, err := u.DB.Exec("INSERT INTO users (login,email,password) VALUES($1,$2,$3)", user.Login, user.Email, user.Password)
	return err
}
