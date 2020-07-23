package handler

import (
	"encoding/json"
	"fmt"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"github.com/vitoraalmeida/desafio-stone-go/pkg/auth"
	"log"
	"net/http"
	"time"
)

//TransferPresenter defines how data will be shown in JSON representation
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

	//Extracts from the token the account id that intends to perform the listing
	user_id, err := auth.ExtractTokenID(r)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
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

	//Prepare data for response in JSON
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

	//Extracts from the token the account id that intends to perform the listing
	originID, err := auth.ExtractTokenID(r)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	//defines the expected structure of the request body
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

	//checks that can be done before searching the database
	if uint(originID) == input.AccountDestinationID {
		http.Error(w, "Destination must be another account", http.StatusBadRequest)
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

	//check that can be done before searching the database again
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

	//Save old balances for possible reversals
	prevOriginBal := originAcc.Balance
	prevDestBal := destAcc.Balance

	originAcc.Balance -= tr.Amount
	destAcc.Balance += tr.Amount

	if err := t.ar.UpdateBalance(originAcc.ID, originAcc.Balance); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if err := t.ar.UpdateBalance(destAcc.ID, destAcc.Balance); err != nil {
		//If the update of the balance at the destination goes wrong,
		//the withdrawal at the source reverses
		t.l.Printf(`Error on updating destination account balance, update on origin acc:%v`, tr)
		if err = t.ar.UpdateBalance(originAcc.ID, prevOriginBal); err != nil {
			t.l.Printf("Error reverting Origin account balance update:\n %v", tr)
		}
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	tr.ID, err = t.tr.Create(tr)
	if err != nil {
		errorMessage := fmt.Sprintf(
			`Error on tranfer from account id: %d, to account id: %d amount: %.2f`,
			tr.AccountOriginID, tr.AccountDestinationID, tr.Amount)
		t.l.Println(errorMessage)
		t.l.Println("Reverting balance updates on origin and destination accounts")

		//If there is a problem creating the transaction in database,the changes
		//in the balances are reversed
		err = t.revertAllBalanceUpdate(originAcc, destAcc, prevOriginBal, prevDestBal)
		if err != nil {
			t.l.Println("Failed to reverse balances updates")
		}
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//Reverts the source and destination account balances to the state prior to the
//failed transfer. If the reversal is unsuccessful, report the error;
func (t *Transfers) revertAllBalanceUpdate(
	originAcc, destAcc *models.Account,
	prevOriginBal, prevDestBal float64,
) error {
	if err := t.ar.UpdateBalance(originAcc.ID, prevOriginBal); err != nil {
		t.l.Printf(`Error reverting balance update on bad transfer:
			Account: %d; Previous balance: %.2f; New balance %.2f`,
			originAcc.ID, prevOriginBal, originAcc.Balance)
		return err
	}
	if err := t.ar.UpdateBalance(destAcc.ID, prevDestBal); err != nil {
		t.l.Printf(`Error reverting balance update on bad transfer:
			Account: %d; Previous balance: %.2f; New balance %.2f`,
			destAcc.ID, prevDestBal, destAcc.Balance)
		return err
	}
	return nil
}
