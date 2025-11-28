package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/fzalvarez/odin-iam/internal/auth"
	"github.com/fzalvarez/odin-iam/internal/users"
	"github.com/google/uuid"
)

const (
	SuperAdminRoleID = "20000000-0000-0000-0000-000000000001"
	SystemTenantID   = "00000000-0000-0000-0000-000000000000"
)

// InitializeSystem crea el usuario admin inicial si no existe ningún usuario en el sistema
func InitializeSystem(ctx context.Context, db *sql.DB) error {
	// Verificar si ya existen usuarios
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// Si ya hay usuarios, no hacer nada
	if count > 0 {
		log.Println("✅ Sistema ya inicializado (usuarios existentes)")
		return nil
	}

	// Obtener credenciales del admin inicial desde variables de entorno
	adminEmail := os.Getenv("INITIAL_ADMIN_EMAIL")
	adminPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")
	adminName := os.Getenv("INITIAL_ADMIN_NAME")

	if adminEmail == "" || adminPassword == "" {
		log.Println("⚠️  No se encontraron INITIAL_ADMIN_EMAIL/PASSWORD, omitiendo creación de admin inicial")
		return nil
	}

	if adminName == "" {
		adminName = "System Administrator"
	}

	log.Println("🔧 Inicializando sistema: creando usuario admin inicial...")

	// Crear repositorios
	userRepo := users.NewRepository(db)
	credRepo := auth.NewCredentialsRepository(db)

	// Generar ID para el nuevo usuario
	adminID := uuid.New()
	tenantID, _ := uuid.Parse(SystemTenantID)

	// Crear usuario
	user, err := userRepo.CreateUser(ctx, tenantID, adminName, adminEmail)
	if err != nil {
		return err
	}

	// Hashear contraseña
	hashedPassword, err := auth.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	// Crear credencial
	err = credRepo.CreateCredential(ctx, user.ID.String(), hashedPassword)
	if err != nil {
		return err
	}

	// Asignar rol Super Admin
	roleID, _ := uuid.Parse(SuperAdminRoleID)
	_, err = db.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role_id, assigned_at) VALUES ($1, $2, NOW())",
		user.ID, roleID,
	)
	if err != nil {
		return err
	}

	log.Printf("✅ Usuario admin inicial creado: %s (ID: %s)", adminEmail, user.ID)
	log.Println("⚠️  IMPORTANTE: Cambia la contraseña después del primer login")

	return nil
}
