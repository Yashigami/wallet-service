package db

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"gorm.io/gorm"
	"log"
)

// Migrate Выполнение миграций
func Migrate(db *gorm.DB) {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("failed to enable uuid-ossp: %v", err)
	}

	err := db.AutoMigrate(&model.Wallet{}, &model.WalletOperation{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
