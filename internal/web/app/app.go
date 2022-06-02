package app

import (
	"log"
	"net/http"
)

func App() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/evaluate/", evaluateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
