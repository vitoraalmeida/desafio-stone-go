package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AccountPresenter struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	CreatedAt time.Time `json:"created_at"`
}

type Accounts struct {
	l  *log.Logger
	ar *models.AccountRepository
}

func NewAccounts(l *log.Logger, ar *models.AccountRepository) *Accounts {
	return &Accounts{
		l,
		ar,
	}
}

func (a *Accounts) ListAccounts(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error reading accounts"
	a.l.Println("Handle GET accounts")
	data, err := a.ar.List()
	w.Header().Set("Content-type", "application/json")
	if err != nil && err != models.ErrAccountNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	var toJ []*AccountPresenter
	for _, a := range data {
		toJ = append(toJ, &AccountPresenter{
			ID:        a.ID,
			Name:      a.Name,
			CPF:       a.CPF,
			CreatedAt: a.CreatedAt,
		})
	}
	if err := json.NewEncoder(w).Encode(toJ); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
	}
}

func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error creating account"
	a.l.Println("Handle POST accounts")
	acc := &models.Account{}

	var input struct {
		Name    string  `json:"name"`
		CPF     string  `json:"cpf"`
		Secret  string  `json:"secret"`
		Balance float64 `json:"balance"`
	}

	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	secret := []byte(input.Secret)
	hashSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	input.Secret = string(hashSecret)
	acc = &models.Account{
		Name:      input.Name,
		CPF:       input.CPF,
		Secret:    input.Secret,
		Balance:   input.Balance,
		CreatedAt: time.Now(),
	}

	acc.ID, err = a.ar.Create(acc)
	if err != nil {
		if err == models.ErrAlreadyRegistred {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("CPF already registred"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	toJ := &AccountPresenter{
		ID:        acc.ID,
		Name:      acc.Name,
		CPF:       acc.CPF,
		CreatedAt: acc.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toJ); err != nil {
		a.l.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

}

func (a *Accounts) GetBalance(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error reading balance"
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	a.l.Printf("Handle GET accounts/%d/balance", id)

	acc, err := a.ar.FindByID(uint(id))

	w.Header().Set("Content-Type", "application/json")
	if err != nil && err != models.ErrAccountNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	if acc == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}

	if err = json.NewEncoder(w).Encode(acc.Balance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
	}

}
