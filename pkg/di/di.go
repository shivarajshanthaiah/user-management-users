package di

import (
	"log"

	"github.com/shivaraj-shanthaiah/user-management/config"
	"github.com/shivaraj-shanthaiah/user-management/pkg/db"
	"github.com/shivaraj-shanthaiah/user-management/pkg/handler"
	"github.com/shivaraj-shanthaiah/user-management/pkg/repo"
	"github.com/shivaraj-shanthaiah/user-management/pkg/server"
	"github.com/shivaraj-shanthaiah/user-management/pkg/service"
)

func Init() {
	cnfg := config.LoadConfig()

	redis, err := config.SetupRedis(cnfg)
	if err != nil {
		log.Fatalf("Failed to connect to redis")
	}
	log.Printf("Successfully connected to redis host: %s", redis.Client)

	db := db.ConnectDB(cnfg)
	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo, redis)
	userHandler := handler.NewUserHandler(userService)
	err = server.NewGrpcUserServer(cnfg.GrpcUserPort, userHandler)
	if err != nil {
		log.Fatalf("failed to start gRPC server %v", err.Error())
	}
}
