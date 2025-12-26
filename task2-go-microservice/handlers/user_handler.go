package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"task2-go-microservice/models"
	"task2-go-microservice/services"
	"task2-go-microservice/utils"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *services.UserService
	logger  *utils.Logger
}

func NewUserHandler(logger *utils.Logger) *UserHandler {
	return &UserHandler{service: services.NewUserService(), logger: logger}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.List()
	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, ok := h.service.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	saved := h.service.Create(user)
	h.logger.AsyncLog("CREATE user id=" + strconv.Itoa(saved.ID))
	respondJSON(w, http.StatusCreated, saved)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated, err := h.service.Update(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.AsyncLog("UPDATE user id=" + strconv.Itoa(updated.ID))
	respondJSON(w, http.StatusOK, updated)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ok := h.service.Delete(id); !ok {
		http.NotFound(w, r)
		return
	}
	h.logger.AsyncLog("DELETE user id=" + strconv.Itoa(id))
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
