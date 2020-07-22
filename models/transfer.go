package models

import (
	"errors"
	"time"
)

var ErrTransferNotFound = errors.New("Models: trasnfer not found")

type Transfer struct {
	ID                   int
	AccountOriginID      int
	AccountDestinationID int
	Amount               float64
	CreatedAt            time.Time
}

func GetTransfers() []*Transfer {
	return transfersList
}

func AddTransfer(t *Transfer) (int, error) {
	t.ID = getNextTransferID()
	transfersList = append(transfersList, t)
	return t.ID, nil
}

func FindTransfersByUserId(id int) ([]*Transfer, error) {
	var trs []*Transfer
	for i := range transfersList {
		t := transfersList[i]
		if t.AccountOriginID == id || t.AccountDestinationID == id {
			trs = append(trs, t)
		}
	}

	if len(trs) == 0 {
		return nil, ErrTransferNotFound
	}

	return trs, nil
}

func getNextTransferID() int {
	t := transfersList[len(transfersList)-1]
	return t.ID + 1
}

var transfersList = []*Transfer{
	{
		ID:                   1,
		AccountOriginID:      1,
		AccountDestinationID: 2,
		Amount:               200.50,
		CreatedAt:            time.Now(),
	},
}
