package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
)

var (
	ErrEmplNotFound        = errors.New("book not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

// @Summary Get Employee By ID
// @Security ApiKeyAuth
// @Tags Employee
// @Description get employee by id
// @ID get-empl-by-id
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} domain.Employee
// @Failure 400,404 No content
// @Router /employee/{id} [get]
func (h *Handler) getEmplByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		logError("getEmplByID", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.employee.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, ErrEmplNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logError("getEmplByID", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		logError("getEmplByID", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create Employee
// @Security ApiKeyAuth
// @Tags Employee
// @Description create employee
// @ID create-empl
// @Accept json
// @Produce json
// @Param input body domain.Employee true "Employee info"
// @Success 200 No content
// @Failure 400,404 No content
// @Router /employee [post]
func (h *Handler) createEmpl(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError("createEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Employee
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		logError("createEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employee.Create(r.Context(), book)
	if err != nil {
		logError("createEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Delete Employee By ID
// @Security ApiKeyAuth
// @Tags Employee
// @Description delete employee by id
// @ID delete-empl
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 No content
// @Failure 400,404 No content
// @Router /employee/{id} [post]
func (h *Handler) deleteEmpl(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		logError("deleteEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employee.Delete(r.Context(), id)
	if err != nil {
		logError("deleteEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get All Employees
// @Security ApiKeyAuth
// @Tags Employee
// @Description get all employees
// @ID get-all-empls
// @Accept json
// @Produce json
// @Param id query int false "Employee ID"
// @Success 200 {array} domain.Employee
// @Failure 400,404 No content
// @Router /employee [get]
func (h *Handler) getAllEmpls(w http.ResponseWriter, r *http.Request) {
	empls, err := h.employee.GetAll(r.Context())
	if err != nil {
		logError("getAllEmpls", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(empls)
	if err != nil {
		logError("getAllEmpls", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Update Employee By ID
// @Security ApiKeyAuth
// @Tags Employee
// @Description update employee
// @ID update-empl
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 No content
// @Failure 400,404 No content
// @Router /employee/{id} [put]
func (h *Handler) updateEmpl(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		logError("updateEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError("updateEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateEmployee
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logError("updateEmpl", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.employee.Update(r.Context(), id, inp)
	if err != nil {
		logError("updateEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
