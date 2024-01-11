package transport

import (
	"errors"
	"net/http"

	"github.com/Sskrill/Authentication.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input domain.SignUpInput

	err := c.BindJSON(&input)
	if err != nil {
		loggerError("sign up", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = input.Validate()
	if err != nil {
		loggerError("sign up", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.userService.SignUp(input)
	if err != nil {
		loggerError("sign up", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput

	err := c.BindJSON(&input)
	if err != nil {
		loggerError("sign in", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		loggerError("sign in", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.userService.SignIn(input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(c, err)
			return
		}
		loggerError("sign in", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"token": token})
}
func handleNotFoundError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
}
