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
	}
}
