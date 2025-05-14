package main

import (
	"github.com/Yashigami/wallet-service/internal/config"
	"github.com/Yashigami/wallet-service/internal/db"
	"github.com/Yashigami/wallet-service/internal/handler"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/Yashigami/wallet-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	cfg           *config.Config
	dbConnection  *gorm.DB
	repoWallet    *repository.WalletRepo
	walletService *service.WalletService
	serverEngine  *gin.Engine
)

func main() {
	err := serverEngine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cfg = config.Load()

	dbConnection = dbConn()

	serverEngine = gin.Default()

	repoWallet = repository.NewWalletRepo(dbConnection)
	walletService = service.NewWalletService(repoWallet)

	handler.NewWalletHandler(walletService, serverEngine)
}

func dbConn() *gorm.DB {
	dbConn, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	var DB, _ = dbConn.DB()
	for i := 0; i < 10; i++ {
		if err := DB.Ping(); err == nil {
			db.Migrate(dbConn)
			return dbConn
		}
		log.Printf("Waiting for DB... (%d)", i+1)
		time.Sleep(2 * time.Second)
	}
	panic(err)
	/*if err != nil {
		panic(err)
	}
	panic(err)*/
}
