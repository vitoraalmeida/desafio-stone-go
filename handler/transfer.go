package handler

import (
	"encoding/json"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"github.com/vitoraalmeida/desafio-stone-go/pkg/auth"
	"log"
	"net/http"
)

type Transfers struct {
	l  *log.Logger
	as *models.AccountService
	ts *models.TransferService
}

func NewTransfers(
	l *log.Logger,
	as *models.AccountService,
	ts *models.TransferService,
) *Transfers {
	return &Transfers{
		l,
		as,
		ts,
	}
}

func (t *Transfers) ListTransfers(w http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle GET transfers")

	user_id, err := auth.ExtractTokenID(r)
	if err != nil {
		http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		return
	}

	trs, err := t.ts.FindByOriginID(int(user_id))
	if err != nil {
		http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		return
	}

	var toj []models.Transfer

	for _, t := range trs {
		toj = append(toj, models.Transfer{
			ID:                   t.ID,
			AccountOriginID:      t.AccountOriginID,
			AccountDestinationID: t.AccountDestinationID,
			Amount:               t.Amount,
			CreatedAt:            t.CreatedAt,
		})
	}

	err = json.NewEncoder(w).Encode(toj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *Transfers) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle POST transfers")

	tr := &models.Transfer{}

	if err := tr.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	if tr.Amount <= 0 {
		http.Error(w, "Amount 0 not allowed", http.StatusBadRequest)
		return
	}

	origin_id, err := auth.ExtractTokenID(r)
	if err != nil {
		http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		return
	}

	origin_acc, err := t.as.FindByID(int(origin_id))
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	t.l.Printf("ACHOU ORIGEM %#v", origin_acc)

	if (origin_acc.Balance - tr.Amount) < 0 {
		http.Error(w, "Origin account do not have enough balance", http.StatusBadRequest)
		return
	}

	dest_id := tr.AccountDestinationID
	dest_acc, err := t.as.FindByID(dest_id)
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	origin_acc.Balance -= tr.Amount
	dest_acc.Balance += tr.Amount

	if err := t.ts.Create(tr); err != nil {
		panic(err)
	}

	if err := t.as.UpdateBalance(origin_acc); err != nil {
		panic(err)
	}
	if err := t.as.UpdateBalance(dest_acc); err != nil {
		panic(err)
	}
}
