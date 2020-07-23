package handler

import (
	"encoding/json"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"github.com/vitoraalmeida/desafio-stone-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Credentials struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

type Login struct {
	l  *log.Logger
	ar *models.AccountRepository
}

func NewLogin(l *log.Logger, ar *models.AccountRepository) *Login {
	return &Login{
		l,
		ar,
	}
}

func (ln *Login) SignIn(w http.ResponseWriter, r *http.Request) {
	ln.l.Println("Handle POST login")
	errorMessage := "Error signing in"
	credentials := &Credentials{}

	err := json.NewDecoder(r.Body).Decode(credentials)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	user, err := ln.ar.FindByCPF(credentials.CPF)
	if err != nil && err != models.ErrAccountNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Secret),
		[]byte(credentials.Secret),
	)
	if err != nil {
		http.Error(w, "Wrong credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}
