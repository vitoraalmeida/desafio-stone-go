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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	user, err := ln.ar.FindByCPF(credentials.CPF)
	if err != nil && err != models.ErrAccountNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Account not found"))
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Secret),
		[]byte(credentials.Secret),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error: Wrong Credentials"))
		return
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	if err = json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
	}
}
