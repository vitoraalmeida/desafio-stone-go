package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

type Accounts struct {
	l  *log.Logger
	as *models.AccountService
}

func NewAccounts(l *log.Logger, as *models.AccountService) *Accounts {
	return &Accounts{
		l,
		as,
	}
}

func (a *Accounts) ListAccounts(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle GET accounts")
	acc, err := a.as.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var toj []models.Account

	for _, a := range acc {
		toj = append(toj, models.Account{
			ID:        a.ID,
			Name:      a.Name,
			CPF:       a.CPF,
			CreatedAt: a.CreatedAt,
		})
	}
	a.l.Println(toj)
	err = json.NewEncoder(w).Encode(toj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle POST accounts")

	acc := &models.Account{}

	if err := acc.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	a.l.Println(acc)

	secret := []byte(acc.Secret)
	hashSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Unable to hash secret", http.StatusInternalServerError)
		return
	}

	acc.Secret = string(hashSecret)

	err = a.as.Create(acc)
	if err != nil {
		a.l.Println(err)
	}
}

func (a *Accounts) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "unable to convert id", http.StatusBadRequest)
	}
	a.l.Printf("Handle GET accounts/%d/balance", id)

	acc, err := a.as.FindByID(id)
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Account not found", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(acc.Balance); err != nil {
		http.Error(w, "Account not found", http.StatusInternalServerError)
	}

}
