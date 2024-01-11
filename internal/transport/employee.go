package transport

import (
	"errors"
	"net/http"

	"github.com/Sskrill/Authentication.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) create(c *gin.Context) {
	var empl domain.Employee
	err := c.BindJSON(&empl)
	if err != nil {
		loggerError("create", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.employee.Create(empl)
	if err != nil {
		loggerError("create", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Writer.WriteHeader(http.StatusCreated)
}

func (h *Handler) get(c *gin.Context) {
	id, err := getIdFromReq(c)
	if err != nil {
		loggerError("get", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	empl, err := h.employee.Get(id)
	if err != nil {
		if errors.Is(domain.ErrEmplNotFound, err) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		loggerError("get", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, empl)
}
func (h *Handler) update(c *gin.Context) {
	id, err := getIdFromReq(c)
	if err != nil {
		loggerError("update", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var empl domain.UpdateEmployee
	err = c.BindJSON(&empl)
	if err != nil {
		loggerError("update", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.employee.Update(id, empl)
	if err != nil {
		loggerError("update", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
func (h *Handler) delete(c *gin.Context) {
	id, err := getIdFromReq(c)
	if err != nil {
		loggerError("delete", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.employee.Delete(id)
	if err != nil {
		loggerError("delete", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
