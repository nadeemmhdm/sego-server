package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/config"
	"github.com/nadeemmhdm/sego-server/internal/db"
	"github.com/nadeemmhdm/sego-server/internal/httpapi"
	"github.com/nadeemmhdm/sego-server/internal/storage"
)

func main() {
	cfgPath := flag.String("config", "sego.json", "path to SEGO config")
	flag.Parse()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := config.Ensure(*cfgPath)
	if err != nil {
		log.Error("config", "error", err)
		os.Exit(1)
	}
	d, err := db.Open(cfg.DatabasePath)
	if err != nil {
		log.Error("database", "error", err)
		os.Exit(1)
	}
	defer d.Close()
	roots := make([]storage.Root, 0, len(cfg.StorageRoots))
	for _, r := range cfg.StorageRoots {
		roots = append(roots, storage.Root{ID: r.ID, Name: r.Name, Path: r.Path})
	}
	st, err := storage.New(roots)
	if err != nil {
		log.Error("storage", "error", err)
		os.Exit(1)
	}
	srv := httpapi.New(cfg, d, st, log)
	h := &http.Server{Addr: cfg.ListenAddr, Handler: srv.Router(), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, WriteTimeout: 0}
	log.Info("SEGO Server started", "listen", cfg.ListenAddr, "low_resource_mode", cfg.LowResourceMode)
	if err := h.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
