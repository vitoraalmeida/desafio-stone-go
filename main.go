package main

import (
	"github.com/gorilla/mux"
	"github.com/vitoraalmeida/desafio-stone-go/handler"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	l := log.New(os.Stdout, "desafio-stone-go", log.LstdFlags)

	ah := handler.NewAccounts(l)
	r := mux.NewRouter()

	r.HandleFunc("/accounts", ah.ListAccounts)

	log.Println("listening on localhost:3000...")

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":3000",
		Handler:      r,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
