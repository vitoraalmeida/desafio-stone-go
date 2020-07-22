package handler

import (
	"encoding/json"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"github.com/vitoraalmeida/desafio-stone-go/pkg/auth"
	"log"
	"net/http"
	"time"
)

type TransferPresenter struct {
	ID                   uint      `json:"transfer_id"`
	AccountOriginID      uint      `json:"account_origin_id"`
	AccountDestinationID uint      `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

type Transfers struct {
	l  *log.Logger
	ar *models.AccountRepository
	tr *models.TransferRepository
}

func NewTransfers(
	l *log.Logger,
	ar *models.AccountRepository,
	tr *models.TransferRepository,
) *Transfers {
	return &Transfers{
		l,
		ar,
		tr,
	}
}

func (t *Transfers) ListTransfers(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error reading transfers"
	t.l.Println("Handle GET transfers")

	user_id, err := auth.ExtractTokenID(r)

	w.Header().Set("Content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized operation: Login required"))
		return
	}

	t.l.Println("chegou no find by origin")
	trs, err := t.tr.FindByOriginID(uint(user_id))
	w.Header().Set("Content-Type", "application/json")
	if err != nil && err != models.ErrTransferNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	if trs == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This account don't have transfers"))
		return
	}

	if err = json.NewEncoder(w).Encode(trs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
	}

}

func (t *Transfers) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error creating transfer"
	t.l.Println("Handle POST transfers")

	originID, err := auth.ExtractTokenID(r)

	w.Header().Set("Content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized operation: Login required"))
		return
	}

	var input struct {
		AccountDestinationID uint    `json:"account_destination_id"`
		Amount               float64 `json:"amount"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	if input.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Amount 0 not allowed"))
		return
	}

	originAcc, err := t.ar.FindByID(uint(originID))
	if err != nil && err != models.ErrAccountNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	if originAcc == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}
	if (originAcc.Balance - input.Amount) < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Insufficient balance"))
		return
	}

	destID := input.AccountDestinationID
	destAcc, err := t.ar.FindByID(uint(destID))
	if err != nil && err != models.ErrAccountNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	if originAcc == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}

	tr := &models.Transfer{
		AccountOriginID:      uint(originID),
		AccountDestinationID: input.AccountDestinationID,
		Amount:               input.Amount,
		CreatedAt:            time.Now(),
	}

	tr.ID, err = t.tr.Create(tr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}

	originAcc.Balance -= tr.Amount
	destAcc.Balance += tr.Amount

	toJ := &TransferPresenter{
		ID:                   tr.ID,
		AccountOriginID:      tr.AccountOriginID,
		AccountDestinationID: tr.AccountDestinationID,
		Amount:               tr.Amount,
		CreatedAt:            tr.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toJ); err != nil {
		t.l.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
}
