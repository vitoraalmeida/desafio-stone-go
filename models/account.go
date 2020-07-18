package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

var ErrAccountNotFound = errors.New("models: account not found")

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Secret    string    `json:"secret"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Account) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

type Accounts []*Account

func (a *Accounts) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func GetAccounts() Accounts {
	return accounts
}

func AddAccount(a *Account) {
	a.ID = getNextID()
	accounts = append(accounts, a)
}

func FindById(id int) (*Account, error) {
	for i := range accounts {
		a := accounts[i]
		if a.ID == id {
			return a, nil
		}
	}
	return nil, ErrAccountNotFound
}

func getNextID() int {
	a := accounts[len(accounts)-1]
	return a.ID + 1
}

var accounts = Accounts{
	&Account{
		ID:        1,
		Name:      "José da Silva",
		CPF:       "01111111111",
		Secret:    "hash",
		Balance:   5000.00,
		CreatedAt: time.Now(),
	},
}
