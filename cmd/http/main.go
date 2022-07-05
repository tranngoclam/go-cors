package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
