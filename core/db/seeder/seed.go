package seeder

import (
	"log"

	models "github.com/antoniusDoni/monorepo/shared/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Auto migrate all tables
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Seed Permissions
	permNames := []string{"view_users", "edit_users", "delete_users"}
	var permissions []models.Permission
	for _, name := range permNames {
		var p models.Permission
		if err := db.FirstOrCreate(&p, models.Permission{Name: name}).Error; err != nil {
			log.Fatalf("Failed to seed permission %s: %v", name, err)
		}
		permissions = append(permissions, p)
	}

	// Seed Roles
	adminRole := models.Role{Name: "admin"}
	if err := db.FirstOrCreate(&adminRole, models.Role{Name: adminRole.Name}).Error; err != nil {
		log.Fatalf("Failed to seed role admin: %v", err)
	}

	// Assign Permissions to Role (RolePermission)
	if err := db.Model(&adminRole).Association("Permissions").Replace(permissions); err != nil {
		log.Printf("Failed to assign permissions to admin role: %v", err)
	}

	// Seed User
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Password hash generation failed: %v", err)
	}

	adminUser := models.User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: string(hash),
	}

	if err := db.FirstOrCreate(&adminUser, models.User{Email: adminUser.Email}).Error; err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// Assign Role to User (UserRole)
	if err := db.Model(&adminUser).Association("Roles").Replace([]models.Role{adminRole}); err != nil {
		log.Printf("Failed to assign role to admin user: %v", err)
	}

	log.Println("Seeding completed.")
}
