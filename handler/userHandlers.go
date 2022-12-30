package handler

import (
	"account-api/model"
	"account-api/service"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type usersResource struct {
	UserService service.UserService
}

func NewUsersResource(userService service.UserService) *usersResource {
	return &usersResource{UserService: userService}
}

func (ur *usersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ur.List)    // GET /users - read a list of users
	r.Post("/", ur.Create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ur.Get)       // GET /users/{id} - read a single users by :id
		r.Put("/", ur.Update)    // PUT /users/{id} - update a single users by :id
		r.Delete("/", ur.Delete) // DELETE /users/{id} - delete a single users by :id
	})
	return r
}

func (ur *usersResource) List(w http.ResponseWriter, r *http.Request) {
	users, err := ur.UserService.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, users)
}

func (ur *usersResource) Create(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	user, err := ur.UserService.CreateUser(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, user)
}

func (ur *usersResource) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	err = ur.UserService.DeleteUser(uint(id))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (ur *usersResource) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	user, err := ur.UserService.GetUser(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, user)
}

func (ur *usersResource) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}
	user := model.User{}

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	getUser, err2 := ur.UserService.GetUser(uint(id))

	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("User with id %s not found", id)))
		return
	}

	getUser.Update(user)

	user, err = ur.UserService.UpdateUser(getUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, user)
}
