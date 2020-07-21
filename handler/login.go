package handler

import (
	"encoding/json"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"github.com/vitoraalmeida/desafio-stone-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Login struct {
	l  *log.Logger
	as *models.AccountService
}

type Credentials struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func NewLogin(l *log.Logger, as *models.AccountService) *Login {
	return &Login{
		l,
		as,
	}
}

func (ln *Login) SignIn(w http.ResponseWriter, r *http.Request) {
	ln.l.Println("Handle POST login")
	credentials := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(credentials)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	ln.l.Println("VAI PROCURAR PELO CPF", credentials.CPF)
	user, err := ln.as.FindByCPF(credentials.CPF)
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Secret),
		[]byte(credentials.Secret),
	)

	if err != nil {
		http.Error(w, "Wrong Credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed create token", http.StatusInternalServerError)
	}
	if err = json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, "Account not found", http.StatusInternalServerError)
	}
}
