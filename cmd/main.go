package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ollub/user_service/config"
	"github.com/Ollub/user_service/internal/middleware"
	"github.com/Ollub/user_service/internal/session"
	"github.com/Ollub/user_service/internal/users/delivery"
	"github.com/Ollub/user_service/internal/users/repo"
	"github.com/Ollub/user_service/internal/users/usecase"
	"github.com/Ollub/user_service/pkg/db"
	"github.com/Ollub/user_service/pkg/log"
	"github.com/gorilla/mux"
)

func NewServer(cfg config.Config) http.Server {
	conn, err := db.GetPostgres(cfg.DbConf)
	if err != nil {
		panic(err)
	}
	user_repo := repo.NewPgRepository(conn)
	user_manager := usecase.NewManager(user_repo)
	session_manager := session.NewSessionsJWTVer(cfg.JwtKey, cfg.TokenTTLDays, user_manager)

	u := delivery.NewHandler(session_manager, user_manager)

	apiHandler := mux.NewRouter()

	apiHandler.HandleFunc("/register", u.Register).Methods("POST")
	apiHandler.HandleFunc("/login", u.Login).Methods("POST")
	apiHandler.HandleFunc("/users", u.List).Methods("GET")
	apiHandler.HandleFunc("/users/{id}", u.Update).Methods("PUT")

	apiHandler.Use(
		middleware.Authentication(session_manager),
		middleware.SetupReqID,
		middleware.InjectLogger,
		middleware.SetupAccessLog,
	)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", apiHandler)

	return http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      siteMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func Init() {
	config.Load()

	logLevel := log.InfoLevel
	if config.Cfg.Debug {
		logLevel = log.DebugLevel
	}
	log.SetupLogger(logLevel)
}

func main() {
	Init()
	server := NewServer(config.Cfg)
	log.Info("Start server")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
