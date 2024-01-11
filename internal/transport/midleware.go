package transport

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CtxValue int

const (
	ctxUserId CtxValue = iota
)

func loggingMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"uri":    c.Request.RequestURI,
		}).Info()
		c.Next()
	}
}
func (h *Handler) authMidleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			loggerError("authMidleWare", err)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, err := h.userService.ParseToken(c.Request.Context(), token)
		if err != nil {
			loggerError("authMidleWare", err)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(c.Request.Context(), ctxUserId, userId)
		c.Request.WithContext(ctx)
		c.Next()
	}
}
func getTokenFromRequest(r *gin.Context) (string, error) {
	header := r.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty authoriztion")
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization")
	}
	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}
	return headerParts[1], nil
}
