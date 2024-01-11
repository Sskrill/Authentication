package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Sskrill/Authentication.git/internal/domain"
	rep "github.com/Sskrill/Authentication.git/internal/repository"
	"github.com/Sskrill/Authentication.git/pkg/hash"
	"github.com/golang-jwt/jwt"
)

type UserService struct {
	repos      *rep.Users
	hasher     hash.PSHasher
	hmacSecret []byte
	ttl        time.Duration
}

func NewUsersService(repos *rep.Users, hasher hash.PSHasher, secret []byte, time time.Duration) UserService {
	return UserService{repos: repos, hasher: hasher, hmacSecret: secret, ttl: time}
}

func (us *UserService) SignUp(input domain.SignUpInput) error {
	password, err := us.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	user := domain.User{
		Login:    input.Login,
		Email:    input.Email,
		Password: password,
	}
	return us.repos.Create(user)
}

func (us *UserService) SignIn(input domain.SignInInput) (string, error) {
	password, err := us.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}
	user, err := us.repos.Get(password, input.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(us.ttl).Unix(),
	})

	return token.SignedString(us.hmacSecret)
}
func (us *UserService) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return us.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}
