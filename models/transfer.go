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
	return TransfersList
}

func AddTransfer(t *Transfer) {
	t.ID = getNextTransferID()
	TransfersList = append(TransfersList, t)
}

func (t *Transfers) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func getNextTransferID() int {
	t := TransfersList[len(TransfersList)-1]
	return t.ID + 1
}

var TransfersList = Transfers{
	&Transfer{
		ID:                   1,
		AccountOriginID:      1,
		AccountDestinationID: 2,
		Amount:               200.50,
		CreatedAt:            time.Now(),
	},
}
