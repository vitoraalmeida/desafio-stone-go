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
		http.Error(w, "Login required", http.StatusUnauthorized)
		return
	}

	trs, err := t.tr.FindByOriginID(uint(user_id))
	w.Header().Set("Content-Type", "application/json")
	if err != nil && err != models.ErrTransferNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return

	}

	if trs == nil {
		http.Error(w, "This account don't have transfers", http.StatusNotFound)
		return
	}

	var toJ []*TransferPresenter
	for _, t := range trs {
		toJ = append(toJ, &TransferPresenter{
			ID:                   t.ID,
			AccountOriginID:      t.AccountOriginID,
			AccountDestinationID: t.AccountDestinationID,
			Amount:               t.Amount,
			CreatedAt:            t.CreatedAt,
		})
	}
	if err = json.NewEncoder(w).Encode(toJ); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func (t *Transfers) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	errorMessage := "Error creating transfer"
	t.l.Println("Handle POST transfers")

	originID, err := auth.ExtractTokenID(r)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	var input struct {
		AccountDestinationID uint    `json:"account_destination_id"`
		Amount               float64 `json:"amount"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if input.Amount <= 0 {
		http.Error(w, "Amount 0 not allowed", http.StatusBadRequest)
		return
	}

	originAcc, err := t.ar.FindByID(uint(originID))
	if err != nil && err != models.ErrAccountNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if originAcc == nil {
		http.Error(w, errorMessage, http.StatusNotFound)
		return
	}

	if (originAcc.Balance - input.Amount) < 0 {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	destID := input.AccountDestinationID
	destAcc, err := t.ar.FindByID(uint(destID))
	if err != nil && err != models.ErrAccountNotFound {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if originAcc == nil {
		http.Error(w, errorMessage, http.StatusNotFound)
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
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	originAcc.Balance -= tr.Amount
	destAcc.Balance += tr.Amount
	if err := t.ar.UpdateBalance(originAcc.ID, originAcc.Balance); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if err := t.ar.UpdateBalance(destAcc.ID, destAcc.Balance); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

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
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}
