package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("listening on localhost:3000...")
	http.ListenAndServe(":3000", nil)
}
