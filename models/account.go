package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"io"
	"log"
	"time"
)

var ErrAccountNotFound = errors.New("models: account not found")

type AccountService struct {
	db *sql.DB
}

func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Secret    string    `json:"secret"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Accounts []*Account

func (as *AccountService) Create(a *Account) error {
	log.Println(a)
	stmt, err := as.db.Prepare(`
		insert into account (name, cpf, secret, balance, created_at)
		values($1,$2,$3,$4,$5)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		a.Name,
		a.CPF,
		a.Secret,
		a.Balance,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (as *AccountService) List() ([]Account, error) {
	stmt, err := as.db.Prepare(`select id, name, cpf, created_at from account`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var accounts []Account
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a Account
		err = rows.Scan(&a.ID, &a.Name, &a.CPF, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (as *AccountService) FindByID(id int) (*Account, error) {
	stmt, err := as.db.Prepare(`select id, name, cpf, balance, created_at from account
	where id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var a Account
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Balance, &a.CreatedAt)
		if err != nil {
			return nil, ErrAccountNotFound
		}
	}
	return &a, nil
}

func (as *AccountService) FindByCPF(cpf string) (*Account, error) {
	log.Println("CHEGOU NO FIND CPF:", cpf)
	stmt, err := as.db.Prepare(`select id, name, cpf, secret, balance, created_at from account
	where cpf = $1`)
	if err != nil {
		log.Println("CHEGOU NO FIND CPF:", cpf)
		return nil, err
	}
	defer stmt.Close()
	var a Account
	rows, err := stmt.Query(cpf)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func (as *AccountService) UpdateBalance(a *Account) error {
	_, err := as.db.Exec(`
		UPDATE account SET balance = $1 WHERE id = $2`, &a.Balance, &a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func (a *Accounts) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}
