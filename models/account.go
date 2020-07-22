package models

import (
	"errors"
	"time"
)

var ErrAccountNotFound = errors.New("models: account not found")

type Account struct {
	ID        int
	Name      string
	CPF       string
	Secret    string
	Balance   float64
	CreatedAt time.Time
}

func GetAccounts() []*Account {
	return accounts
}

func AddAccount(a *Account) int {
	a.ID = getNextID()
	accounts = append(accounts, a)
	return a.ID
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

func FindByCPF(cpf string) (*Account, error) {
	for i := range accounts {
		a := accounts[i]
		if a.CPF == cpf {
			return a, nil
		}
	}
	return nil, ErrAccountNotFound
}

func getNextID() int {
	a := accounts[len(accounts)-1]
	return a.ID + 1
}

var accounts = []*Account{
	{
		ID:        1,
		Name:      "José da Silva",
		CPF:       "01111111111",
		Secret:    "hash",
		Balance:   9999.99,
		CreatedAt: time.Now(),
	},
}
