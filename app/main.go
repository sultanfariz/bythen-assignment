package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	commons "app/internal/commons"
	handler "app/internal/handlers"
	"app/internal/repositories"
	userRepository "app/internal/repositories/user"
	userUsecase "app/internal/usecases/user"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if viper.GetBool("debug") {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	lis, err := net.Listen("tcp", ":"+viper.GetString("SERVER_PORT"))
	timeoutContext := time.Duration(viper.GetInt("CONTEXT_TIMEOUT")) * time.Second

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	fmt.Println("Server running on port " + viper.GetString("SERVER_PORT"))

	configJWT := commons.ConfigJWT{
		SecretJWT:       viper.GetString("JWT_SECRET_KEY"),
		ExpiresDuration: viper.GetInt("JWT_EXPIRES_DURATION"),
	}

	configDB := repositories.DBConfig{
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Name:     viper.GetString("DB_NAME"),
	}

	db := repositories.InitDB(configDB)
	userRepo := userRepository.NewUserRepository(db, timeoutContext)
	userUsecase := userUsecase.NewUserUsecase(userRepo, configJWT, timeoutContext)
	userHandler := handler.NewUserHandler(*userUsecase)

	// Create a new HTTP mux
	mux := http.NewServeMux()

	// Register the handler function
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:    ":" + viper.GetString("SERVER_PORT"),
		Handler: mux,
	}

	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
