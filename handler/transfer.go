package handler

import (
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"log"
	"net/http"
)

type Transfers struct {
	l *log.Logger
}

func NewTransfers(l *log.Logger) *Transfers {
	return &Transfers{l}
}

func (t *Transfers) ListTransfers(w http.ResponseWriter, r *http.Request) {
	t.l.Println("Handle GET transfers")

	tr := models.GetTransfers()
	if err := tr.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
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

	origin_id := tr.AccountOriginID
	origin_acc, err := models.FindById(origin_id)
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if (origin_acc.Balance - tr.Amount) < 0 {
		http.Error(w, "Origin account do not have enough balance", http.StatusBadRequest)
		return
	}

	dest_id := tr.AccountDestinationID
	dest_acc, err := models.FindById(dest_id)
	if err == models.ErrAccountNotFound {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	origin_acc.Balance -= tr.Amount
	dest_acc.Balance += tr.Amount

	models.AddTransfer(tr)
	t.l.Printf("Acc: %#v", tr)
}
