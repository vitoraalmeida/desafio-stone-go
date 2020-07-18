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

type Transfers []*Transfer

func (t *Transfers) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

var transfers = Transfers{
	&Transfer{
		ID:                   1,
		AccountOriginID:      1,
		AccountDestinationID: 2,
		Amount:               200.00,
		CreatedAt:            time.Now(),
	},
}
