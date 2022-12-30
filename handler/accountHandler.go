package handler

import (
	"account-api/model"
	"account-api/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

type accountsResource struct {
	AccountService   service.AccountService
	accountValidator service.BoValidator
	kafkaSvc         service.EventSender
}

func NewAccountsResource(accountService service.AccountService, accountValidator service.BoValidator, kafkaSvc service.EventSender) *accountsResource {
	return &accountsResource{AccountService: accountService, accountValidator: accountValidator, kafkaSvc: kafkaSvc}
}

func (ar *accountsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ar.List)    // GET /accounts - read a list of accounts
	r.Post("/", ar.Create) // POST /accounts - create a new account and persist it

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a users map, and lets actually load/manipulate
		r.Get("/", ar.Get)       // GET /accounts/{id} - read a single account by :id
		r.Put("/", ar.Update)    // PUT /accounts/{id} - update a single account by :id
		r.Delete("/", ar.Delete) // DELETE /accounts/{id} - delete a single account by :id
	})

	return r
}

func (ar *accountsResource) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := ar.AccountService.GetAccounts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, accounts)
}

func (ar *accountsResource) Create(w http.ResponseWriter, r *http.Request) {
	account := model.NewAccount()

	if err := render.DecodeJSON(r.Body, &account); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	if err := ar.accountValidator.Validate(account); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	user, _, _ := r.BasicAuth()
	account.UpdatedBy = user

	msgKey := model.NewFourEyesMessageKey(account.ID, "cnc-account", model.CREATE)

	if err := ar.kafkaSvc.Send(*msgKey, account); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	log.Printf("Account %d sent to 4yes service", account.ID)

	render.Status(r, http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", account.ID))
}

func (ar *accountsResource) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	account, err := ar.AccountService.GetAccount(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, account)
}

func (ar *accountsResource) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	account := model.Account{}

	if err := render.DecodeJSON(r.Body, &account); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	if err := ar.accountValidator.Validate(account); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	getAccount, err2 := ar.AccountService.GetAccount(uint(id))

	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Account with id %s not found", id)))
		return
	}

	getAccount.Update(account)

	user, _, _ := r.BasicAuth()
	getAccount.UpdatedBy = user

	account, err = ar.AccountService.UpdateAccount(getAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, account)
}

func (ar *accountsResource) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	err = ar.AccountService.DeleteAccount(uint(id))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
