package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/vitoraalmeida/desafio-stone-go/handler"
	"github.com/vitoraalmeida/desafio-stone-go/models"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env: %v", err)
	} else {
		log.Println("We are getting the env values")
	}

	l := log.New(os.Stdout, "desafio-stone-go", log.LstdFlags)

	ah := handler.NewAccounts(l)
	th := handler.NewTransfers(l)
	lh := handler.NewLogin(l)

	r := mux.NewRouter()

	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/accounts", ah.ListAccounts)
	getRouter.HandleFunc("/accounts/{id:[0-9]+}/balance", ah.GetBalance)
	getRouter.HandleFunc("/transfers", th.ListTransfers)

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/accounts", ah.CreateAccount)
	postRouter.HandleFunc("/transfers", th.CreateTransfer)
	postRouter.HandleFunc("/login", lh.SignIn)

	log.Println("listening on localhost:3000...")

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":3000",
		Handler:      r,
	}

	log.Println(models.TransfersList[0])

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
