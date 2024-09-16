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
	commentRepository "app/internal/repositories/comment"
	postRepository "app/internal/repositories/post"
	userRepository "app/internal/repositories/user"
	usecases "app/internal/usecases"

	"github.com/gorilla/mux"

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
	userUsecase := usecases.NewUserUsecase(userRepo, configJWT, timeoutContext)
	userHandler := handler.NewUserHandler(userUsecase)

	postRepo := postRepository.NewPostRepository(db, timeoutContext)
	postUsecase := usecases.NewPostUsecase(postRepo, userRepo, timeoutContext)
	postHandler := handler.NewPostHandler(postUsecase)

	commentRepo := commentRepository.NewCommentRepository(db, timeoutContext)
	commentUsecase := usecases.NewCommentUsecase(commentRepo, postRepo, userRepo, timeoutContext)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	r := mux.NewRouter()

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/posts", configJWT.JWTMiddleware(postHandler.CreatePost)).Methods("POST")
	r.HandleFunc("/posts", postHandler.GetAllPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", postHandler.GetPostByID).Methods("GET")
	r.HandleFunc("/posts/{id}", configJWT.JWTMiddleware(postHandler.UpdatePost)).Methods("PUT")
	r.HandleFunc("/posts/{id}", configJWT.JWTMiddleware(postHandler.DeletePost)).Methods("DELETE")

	r.HandleFunc("/posts/{id}/comments", configJWT.JWTMiddleware(commentHandler.CreateComment)).Methods("POST")
	r.HandleFunc("/posts/{id}/comments", commentHandler.GetCommentsByPostID).Methods("GET")

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:    ":" + viper.GetString("SERVER_PORT"),
		Handler: r,
	}

	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
