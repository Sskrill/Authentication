package domain

import "errors"

var (
	ErrEmplNotFound        = errors.New("employee not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
