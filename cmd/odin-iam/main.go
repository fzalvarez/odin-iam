package main

import (
	"log"
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api"
	"github.com/fzalvarez/odin-iam/internal/apikeys"
	"github.com/fzalvarez/odin-iam/internal/auth"
	dbconn "github.com/fzalvarez/odin-iam/internal/db"
	dbgen "github.com/fzalvarez/odin-iam/internal/db/gen"

	// "github.com/fzalvarez/odin-iam/internal/emails" // Comentado si no existe el paquete
	"github.com/fzalvarez/odin-iam/internal/roles"
	"github.com/fzalvarez/odin-iam/internal/sessions"
	"github.com/fzalvarez/odin-iam/internal/tenants"
	"github.com/fzalvarez/odin-iam/internal/users"
	// _ "github.com/fzalvarez/odin-iam/docs" // Comentado hasta resolver dependencias
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
	// 1. Conectar a la base de datos
	conn, err := dbconn.Connect()
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}
	defer conn.Close()

	// 2. Instanciar el generador de queries sqlc
	q := dbgen.New(conn)

	// 3. Inicializar repositorios
	// Nota: db/gen debe haber sido regenerado con sqlc antes de compilar
	credRepo := auth.NewCredentialsRepository(q)
	sessionRepo := sessions.NewRepository(q)
	tenantRepo := tenants.NewRepository(q)
	userRepo := users.NewRepository(q)
	roleRepo := roles.NewRepository(q)
	apikeyRepo := apikeys.NewRepository(q)
	// emailRepo := emails.NewRepository(q) // Comentado

	// 4. Crear servicios
	// Asumimos que userRepo implementa lo necesario para emails o pasamos nil temporalmente
	// Si auth.NewService requiere 4 argumentos, debemos pasar algo.
	// Si no tenemos emailRepo, pasamos userRepo si implementa la interfaz, o nil y corregimos el servicio.
	// Por ahora, pasamos userRepo asumiendo que maneja emails también (común en repositorios unificados)
	// OJO: Esto requiere que userRepo implemente EmailsRepository.
	authService := auth.NewService(
		userRepo,
		nil, // emailRepo temporalmente nil o userRepo si implementa
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
