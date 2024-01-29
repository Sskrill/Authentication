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
	err = h.usersService.Publish("Method: GET ,Enity: Employee")
	if err != nil {
		logError("getEmplByID", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

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
	err = h.usersService.Publish("Method: POST ,Enity: Employee")
	if err != nil {
		logError("createEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

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
	err = h.usersService.Publish("Method: DELETE ,Enity: Employee")
	if err != nil {
		logError("deleteEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

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
	err = h.usersService.Publish("Method: PUT ,Enity: Employee")
	if err != nil {
		logError("updateEmpl", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
