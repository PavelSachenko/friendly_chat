package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":3001", http.FileServer(http.Dir("./public"))))
}
