package main

import (
	"log"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api"
	"github.com/fzalvarez/odin-iam/internal/auth"
	dbconn "github.com/fzalvarez/odin-iam/internal/db"
	dbgen "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/fzalvarez/odin-iam/internal/users"
)

func main() {
	// 1. Conectar a la base de datos
	conn, err := dbconn.Connect()
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}
	defer conn.Close()

	// 2. Instanciar el generador de queries sqlc
	q := dbgen.New(conn)

	// 3. Crear repositorios
	userRepo := users.NewRepository(q)
	emailRepo := users.NewEmailsRepository(q)
	credRepo := auth.NewCredentialsRepository(q)
	sessionRepo := sessions.NewRepository(q)

	// 4. Crear servicio de autenticación
	authService := auth.NewService(
		userRepo,
		emailRepo,
		credRepo,
		sessionRepo,
	)

	// 5. Crear router con dependencias
	r := api.NewRouter(api.RouterParams{
		AuthService: authService,
	})

	// 6. Iniciar servidor
	log.Println("🚀 IAM service running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ server error: %v", err)
	}
}
