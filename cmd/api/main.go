package main

import (
	"context"
	"database/sql"

	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PhantomXD-nepal/goauthtemplate/db/generated/sqlc"
	"github.com/PhantomXD-nepal/goauthtemplate/internal/config"
	"github.com/PhantomXD-nepal/goauthtemplate/internal/server"
	logger "github.com/PhantomXD-nepal/goauthtemplate/package"
)

func main() {
	cfg := config.Envs
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddress,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to db with err: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping db with err: " + err.Error())
	}
	logger.Mascot()
	logger.Info("Connected to database successfully")

	queries := sqlc.New(db)
	_ = queries
	apiServer := server.NewAPIServer(fmt.Sprintf(":%s", cfg.Port), db)

	logger.Info("Starting API server on port " + cfg.Port)

	if err := apiServer.Start(); err != nil {
		logger.Error("Failed to start server with err: " + err.Error())
	}

}
