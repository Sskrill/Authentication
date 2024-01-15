package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(user domain.User) error
	GetByCredentials(email, password string) (domain.User, error)
}

type SessionsRepository interface {
	Create(session domain.Session) error
	GetId(id string) error
}

type Users struct {
	repo         UsersRepository
	sessionsRepo SessionsRepository
	hasher       PasswordHasher

	hmacSecret []byte
}

func NewUsers(repo UsersRepository, sessionsRepo SessionsRepository, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:         repo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		hmacSecret:   secret,
	}
}

func (s *Users) SignUp(inp domain.SignUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(user)
}

func (s *Users) SignIn(inp domain.SignInInput) (string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(inp.Email, password)
	if err != nil {
		return "", err
	}
	id, err := newSessionId()
	if err != nil {
		return "", err
	}
	session := domain.Session{Id: id, UserId: int(user.ID)}
	err = s.sessionsRepo.Create(session)
	return id, err
}

func newSessionId() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
