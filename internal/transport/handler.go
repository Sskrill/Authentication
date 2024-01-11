package transport

import (
	"errors"
	"strconv"

	repos "github.com/Sskrill/Authentication.git/internal/repository"
	"github.com/Sskrill/Authentication.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService service.UserService
	employee    *repos.Employees
}

func NewHandler(users service.UserService, employees *repos.Employees) *Handler {
	return &Handler{userService: users, employee: employees}
}

func (h *Handler) InitRout() {
	rout := gin.Default()
	rout.Use(loggingMidleware())

	employee := rout.Group("/employee")
	employee.Use(h.authMidleWare())
	employee.GET("/:id", h.get)
	employee.POST("/create", h.create)
	employee.PUT("/:id", h.update)
	employee.DELETE("/:id", h.delete)

	auth := rout.Group("/auth")

	auth.GET("/sign-in", h.signIn)
	auth.POST("/sign-up", h.signUp)
	rout.Run(":8080")
}

func getIdFromReq(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("id cant be zero")
	}
	return id, nil
}
