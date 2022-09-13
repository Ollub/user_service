package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
	"user_service/config"
	"user_service/internal/middleware"
	"user_service/internal/session"
	"user_service/internal/users/delivery"
	"user_service/internal/users/repo"
	"user_service/internal/users/usecase"
	db "user_service/pkg/db"
	"user_service/pkg/log"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func NewServer(cfg config.Config) http.Server {
	conn, err := db.GetPostgres(cfg.DbConf)
	if err != nil {
		panic(err)
	}
	user_repo := repo.NewPgRepository(conn)
	user_manager := usecase.NewManager(user_repo)
	session_manager := session.NewSessionsJWTVer(cfg.JwtKey, cfg.TokenTTLDays, user_manager)

	u := delivery.NewHandler(session_manager, user_manager)

	//mux := http.NewServeMux()
	//
	//mux.HandleFunc("/register", u.Register)
	//mux.HandleFunc("/users", u.List)
	////mux.HandleFunc("/user/logout", u.Logout)
	////mux.HandleFunc("/user/reg", u.Reg)
	////mux.HandleFunc("/user/change_pass", u.ChangePassword)
	//
	////http.Handle("/", middleware.AuthMiddleware(session_manager, mux))
	//
	//siteMux := middleware.SetupReqID(http.Handler(mux))
	//siteMux = middleware.InjectLogger(http.Handler(mux))
	//siteMux = middleware.SetupAccessLog(http.Handler(mux))
	//siteMux = middleware.AuthMiddleware(session_manager, siteMux)

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
