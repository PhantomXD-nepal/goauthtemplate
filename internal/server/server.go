package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/PhantomXD-nepal/goauthtemplate/db/generated/sqlc"
	"github.com/PhantomXD-nepal/goauthtemplate/internal/services/user"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr    string
	db      *sql.DB
	queries *sqlc.Queries
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	//Initialaize the routers
	authRouter := router.PathPrefix("/auth").Subrouter()

	userService := user.NewService(sqlc.New(s.db))
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(authRouter)

	return http.ListenAndServe(s.addr, router)
}
