package models

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"time"
)

var ErrTransferNotFound = errors.New("Models: transfer not found")

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

type Transfer struct {
	ID                   uint
	AccountOriginID      uint
	AccountDestinationID uint
	Amount               float64
	CreatedAt            time.Time
}

func (ts *TransferRepository) Create(t *Transfer) (uint, error) {
	var id uint
	err := ts.db.QueryRow(`
		insert into transfer (acc_origin_id, acc_destination_id, amount, created_at)
		values($1,$2,$3,$4) RETURNING id_transfer`,
		t.AccountOriginID,
		t.AccountDestinationID,
		t.Amount,
		time.Now().Format("2006-01-02")).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ts *TransferRepository) FindByOriginID(uid uint) ([]*Transfer, error) {
	stmt, err := ts.db.Prepare(`
		select id_transfer, acc_origin_id, acc_destination_id, amount,  created_at from transfer
		where acc_origin_id = $1 OR acc_destination_id = $1`)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var transfers []*Transfer
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
			return nil, err
		}
		transfers = append(transfers, &t)
	}
	if len(transfers) == 0 {
		return nil, ErrTransferNotFound
	}
	return transfers, nil
}
