package app

import (
	"log"
	"net/http"
)

func App() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/recommend/", recommendHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
