package seeder

import (
	"log"
	"time"

	warehouseModels "github.com/antoniusDoni/monorepo/modules/warehouse/model"
	models "github.com/antoniusDoni/monorepo/shared/model"
	"github.com/google/uuid"
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
		&warehouseModels.Office{},    // add Office
		&warehouseModels.Branch{},    // add Branch
		&warehouseModels.Warehouse{}, // updated Warehouse with OfficeID and BranchID
		&warehouseModels.CategoryProduct{},
		&warehouseModels.Product{},
		&warehouseModels.UnitProduct{},
		&warehouseModels.StockEntry{},
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

	units := []warehouseModels.UnitProduct{
		// Metric
		{ID: uuid.New(), Code: "pcs", Name: "Piece"},
		{ID: uuid.New(), Code: "kg", Name: "Kilogram"},
		{ID: uuid.New(), Code: "g", Name: "Gram"},
		{ID: uuid.New(), Code: "mg", Name: "Milligram"},
		{ID: uuid.New(), Code: "l", Name: "Liter"},
		{ID: uuid.New(), Code: "ml", Name: "Milliliter"},
		{ID: uuid.New(), Code: "m", Name: "Meter"},
		{ID: uuid.New(), Code: "cm", Name: "Centimeter"},
		{ID: uuid.New(), Code: "mm", Name: "Millimeter"},

		// Imperial
		{ID: uuid.New(), Code: "lb", Name: "Pound"},
		{ID: uuid.New(), Code: "oz", Name: "Ounce"},
		{ID: uuid.New(), Code: "gal", Name: "Gallon"},
		{ID: uuid.New(), Code: "qt", Name: "Quart"},
		{ID: uuid.New(), Code: "pt", Name: "Pint"},
		{ID: uuid.New(), Code: "ft", Name: "Foot"},
		{ID: uuid.New(), Code: "in", Name: "Inch"},

		// Packaged/Other
		{ID: uuid.New(), Code: "box", Name: "Box"},
		{ID: uuid.New(), Code: "bag", Name: "Bag"},
		{ID: uuid.New(), Code: "btl", Name: "Bottle"},
		{ID: uuid.New(), Code: "can", Name: "Can"},
		{ID: uuid.New(), Code: "roll", Name: "Roll"},
		{ID: uuid.New(), Code: "pack", Name: "Pack"},
		{ID: uuid.New(), Code: "carton", Name: "Carton"},
		{ID: uuid.New(), Code: "set", Name: "Set"},
	}

	for _, unit := range units {
		var existing warehouseModels.UnitProduct
		err := db.Where("code = ?", unit.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			unit.CreatedAt = time.Now()
			unit.UpdatedAt = time.Now()
			if err := db.Create(&unit).Error; err != nil {
				log.Printf("❌ Failed to seed unit %s: %v", unit.Code, err)
			} else {
				log.Printf("✅ Seeded unit: %s", unit.Code)
			}
		}
	}

	log.Println("Seeding completed.")
}
