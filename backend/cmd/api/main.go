package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"backend/internal/config"
	"backend/internal/http/router"
	"backend/internal/repository/sqlite"

)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	// SQLite için klasör yoksa oluştur
	if err := os.MkdirAll("./data", 0o755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	// schema init
if err := sqlite.InitSchema(db); err != nil {
	log.Fatalf("schema init failed: %v", err)
}


	addr := ":" + cfg.AppPort
	h := router.New(db, cfg)



	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	fmt.Printf("API listening on http://localhost%s (env=%s)\n", addr, cfg.AppEnv)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
