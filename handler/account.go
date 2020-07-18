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
	l *log.Logger
}

func NewAccounts(l *log.Logger) *Accounts {
	return &Accounts{l}
}

func (a *Accounts) ListAccounts(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle GET accounts")
	acc := models.GetAccounts()
	if err := acc.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle POST accounts")

	acc := &models.Account{}

	if err := acc.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	secret := []byte(acc.Secret)
	hashSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Unable to hash secret", http.StatusInternalServerError)
		return
	}

	acc.Secret = string(hashSecret)

	models.AddAccount(acc)
	a.l.Printf("Acc: %#v", acc)
}

func (a *Accounts) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "unable to convert id", http.StatusBadRequest)
	}
	a.l.Printf("Handle GET accounts/%d/balance", id)

	acc, err := models.FindById(id)

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
