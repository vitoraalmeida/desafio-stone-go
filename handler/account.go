package handler

import (
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"log"
	"net/http"
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
	}
	models.AddAccount(acc)
	a.l.Printf("Prod: %#v", acc)
}
