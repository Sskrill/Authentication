package domain

import "errors"

var (
	ErrEmplNotFound        = errors.New("book not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
