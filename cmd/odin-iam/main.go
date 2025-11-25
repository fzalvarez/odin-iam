package main

import (
	"log"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api"
	"github.com/fzalvarez/odin-iam/internal/auth"
	dbconn "github.com/fzalvarez/odin-iam/internal/db"
	dbgen "github.com/fzalvarez/odin-iam/internal/db/gen"
	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/fzalvarez/odin-iam/internal/tenants"
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
	tenantRepo := tenants.NewRepository(q)
	roleRepo := roles.NewMockRepository() // TODO: Reemplazar con implementación SQL real

	// 4. Crear servicios
	authService := auth.NewService(
		userRepo,
		emailRepo,
		credRepo,
		sessionRepo,
	)
	userService := users.NewService(userRepo)
	tenantService := tenants.NewService(tenantRepo)
	roleService := roles.NewRoleService(roleRepo)

	// 5. Crear router con dependencias
	r := api.NewRouter(api.RouterParams{
		AuthService:   authService,
		UserService:   userService,
		TenantService: tenantService,
		RoleService:   roleService,
	})

	// 6. Iniciar servidor
	log.Println("🚀 IAM service running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ server error: %v", err)
	}
}
