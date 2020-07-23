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
	a.l.Println("Handle GET accounts")
	errorMessage := "Error reading accounts"
	data, err := a.ar.List()
	w.Header().Set("Content-type", "application/json")
	if err != nil && err != models.ErrAccountNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
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
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}

func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handle POST accounts")
	errorMessage := "Error creating account"
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
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	secret := []byte(input.Secret)
	hashSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
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
			http.Error(w, "CPF already registred", http.StatusConflict)
			return
		}
		http.Error(w, errorMessage, http.StatusInternalServerError)
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
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func (a *Accounts) GetBalance(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error reading balance"

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	a.l.Printf("Handle GET accounts/%d/balance", id)

	acc, err := a.ar.FindByID(uint(id))
	if err != nil && err != models.ErrAccountNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if acc == nil {
		http.Error(w, errorMessage, http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(w).Encode(acc.Balance); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}
