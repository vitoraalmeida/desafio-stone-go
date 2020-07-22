package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vitoraalmeida/desafio-stone-go/handler"
	"github.com/vitoraalmeida/desafio-stone-go/middlewares"
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
		log.Println("Getting the env values...")
	}

	dsn := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"), os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	l := log.New(os.Stdout, "desafio-stone-go", log.LstdFlags)

	ar := models.NewAccountRepository(db)
	tr := models.NewTransferRepository(db)

	ah := handler.NewAccounts(l, ar)
	th := handler.NewTransfers(l, ar, tr)
	lh := handler.NewLogin(l, ar)
	r := mux.NewRouter()

	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/accounts", ah.ListAccounts)
	getRouter.HandleFunc("/accounts/{id:[0-9]+}/balance", ah.GetBalance)
	getRouter.HandleFunc("/transfers", middlewares.Authentication(th.ListTransfers))

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/accounts", ah.CreateAccount)
	postRouter.HandleFunc("/transfers", middlewares.Authentication(th.CreateTransfer))
	postRouter.HandleFunc("/login", lh.SignIn)

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":3000",
		Handler:      r,
	}

	log.Println("listening on localhost:3000...")
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
