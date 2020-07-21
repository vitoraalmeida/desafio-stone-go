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

var ErrTransferNotFound = errors.New("Transfer not found")

type TransferService struct {
	db *sql.DB
}

func NewTransferService(db *sql.DB) *TransferService {
	return &TransferService{
		db: db,
	}
}

type Transfer struct {
	ID                   int       `json:"id"`
	AccountOriginID      int       `json:"account_origin_id"`
	AccountDestinationID int       `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

type Transfers []*Transfer

func (ts *TransferService) Create(t *Transfer) error {
	log.Println(t)

	stmt, err := ts.db.Prepare(`
		insert into transfer (acc_origin_id, acc_destination_id, amount, created_at)
		values($1,$2,$3,$4)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		t.AccountOriginID,
		t.AccountDestinationID,
		t.Amount,
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

func (ts *TransferService) FindByOriginID(uid int) ([]Transfer, error) {
	stmt, err := ts.db.Prepare(`
		select id, acc_origin_id, acc_destination_id,  created_at from account
		where acc_origin_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var transfers []Transfer
	rows, err := stmt.Query(uid)
	if err != nil {
		return nil, ErrTransferNotFound
	}
	for rows.Next() {
		var t Transfer
		err = rows.Scan(
			&t.ID,
			&t.AccountOriginID,
			&t.AccountDestinationID,
			&t.Amount,
			&t.CreatedAt)

		if err != nil {
			return nil, ErrAccountNotFound
		}

		transfers = append(transfers, t)
	}
	return transfers, nil
}

func (t *Transfer) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (t *Transfers) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}
