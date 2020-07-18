package models

import (
	"encoding/json"
	"io"
	"time"
)

type Transfer struct {
	ID                   int       `json:"id"`
	AccountOriginID      int       `json:"account_origin_id"`
	AccountDestinationID int       `json:"account_destination_id"`
	Amount               float32   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func (t *Transfer) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

type Transfers []*Transfer

func GetTransfers() Transfers {
	return transfersList
}

func AddTransfer(t *Transfer) {
	t.ID = getNextTransferID()
	transfersList = append(transfersList, t)
}

func FindTransfersByUserId(id int) Transfers {
	var trs Transfers
	for i := range transfersList {
		t := transfersList[i]
		if t.AccountOriginID == id {
			trs = append(trs, t)
		}
	}

	return trs
}

func (t *Transfers) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func getNextTransferID() int {
	t := transfersList[len(transfersList)-1]
	return t.ID + 1
}

var transfersList = Transfers{
	&Transfer{
		ID:                   1,
		AccountOriginID:      1,
		AccountDestinationID: 2,
		Amount:               200.50,
		CreatedAt:            time.Now(),
	},
}
