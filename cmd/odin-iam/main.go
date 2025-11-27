package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/fzalvarez/odin-iam/internal/api"
	"github.com/fzalvarez/odin-iam/internal/apikeys"
	"github.com/fzalvarez/odin-iam/internal/auth"
	dbconn "github.com/fzalvarez/odin-iam/internal/db"

	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/fzalvarez/odin-iam/internal/tenants"
	"github.com/fzalvarez/odin-iam/internal/users"

	_ "github.com/fzalvarez/odin-iam/docs"
)

// @title           Odin IAM API
// @version         1.0
// @description     Identity and Access Management Service with Multi-tenancy and RBAC.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@odin-iam.local

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// 1. Conectar a la base de datos
	conn, err := dbconn.Connect()
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}
	defer conn.Close()

	// 3. Inicializar repositorios
	// Nota: db/gen debe haber sido regenerado con sqlc antes de compilar
	credRepo := auth.NewCredentialsRepository(conn)
	sessionRepo := sessions.NewRepository(conn)
	tenantRepo := tenants.NewRepository(conn)
	userRepo := users.NewRepository(conn)
	roleRepo := roles.NewRepository(conn)
	apikeyRepo := apikeys.NewRepository(conn)

	// 4. Crear servicios
	// userRepo implementa AuthEmailsRepository (AddEmail, GetByEmail)
	authService := auth.NewService(
		userRepo,
		userRepo, // Usamos userRepo para emails
		credRepo,
		sessionRepo,
	)
	userService := users.NewService(userRepo)
	tenantService := tenants.NewService(tenantRepo)
	roleService := roles.NewRoleService(roleRepo)
	apikeyService := apikeys.NewService(apikeyRepo) // Nuevo servicio

	// 5. Crear router con dependencias
	r := api.NewRouter(api.RouterParams{
		AuthService:   authService,
		UserService:   userService,
		TenantService: tenantService,
		RoleService:   roleService,
		APIKeyService: apikeyService, // Inyección
	})

	// 6. Iniciar servidor
	log.Println("🚀 IAM service running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("❌ server error: %v", err)
	}
}
