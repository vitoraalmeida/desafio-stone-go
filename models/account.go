package models

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	ErrAccountNotFound  = errors.New("models: account not found")
	ErrAlreadyRegistred = errors.New("models: cpf already registred")
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

type Account struct {
	ID        uint
	Name      string
	CPF       string
	Secret    string
	Balance   float64
	CreatedAt time.Time
}

type Accounts []*Account

func (as *AccountRepository) Create(a *Account) (uint, error) {
	var id uint
	err := as.db.QueryRow(`
		insert into account (name, cpf, secret, balance, created_at)
		values($1,$2,$3,$4,$5) RETURNING id`,
		a.Name,
		a.CPF,
		a.Secret,
		a.Balance,
		time.Now().Format("2006-01-02")).Scan(&id)

	if err != nil {
		if id == 0 {
			return 0, ErrAlreadyRegistred
		}
		return 0, err
	}

	return id, nil
}

func (as *AccountRepository) List() ([]*Account, error) {
	log.Println("Chegou na database")
	stmt, err := as.db.Prepare(`select id, name, cpf, created_at from account`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var accounts []*Account
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
		log.Println(accounts)
		accounts = append(accounts, &a)

	}

	if len(accounts) == 0 {
		return nil, ErrAccountNotFound
	}

	log.Println(accounts)
	return accounts, nil
}

func (as *AccountRepository) FindByID(id uint) (*Account, error) {
	stmt, err := as.db.Prepare(`select id, name, cpf, balance, created_at from account
	where id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var a Account
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Balance, &a.CreatedAt)
		if err != nil {
			return nil, ErrAccountNotFound
		}
	}
	return &a, nil
}

func (ar *AccountRepository) FindByCPF(cpf string) (*Account, error) {
	stmt, err := ar.db.Prepare(`
		select id, name, cpf, secret, balance, created_at from account
		where cpf = $1`)
	if err != nil {
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
			return nil, ErrAccountNotFound
		}
	}
	return &a, nil
}
