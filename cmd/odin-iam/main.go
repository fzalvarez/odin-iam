package main

import (
	"log"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api"
	"github.com/fzalvarez/odin-iam/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	r := api.NewRouter()

	addr := ":" + cfg.Port
	log.Printf("odin-iam running on %s\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
