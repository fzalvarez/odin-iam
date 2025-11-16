package main

import (
	"log"
	"net/http"

	//"github.com/go-chi/chi/v5"
	"github.com/fzalvarez/odin-iam/internal/api"
)

func main() {
	//r := chi.NewRouter()
	r := api.NewRouter()

	/* r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	}) */

	log.Println("IAM service running on :8080")
	http.ListenAndServe(":8080", r)
}
