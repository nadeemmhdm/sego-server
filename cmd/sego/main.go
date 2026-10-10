package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/auth"
	"github.com/nadeemmhdm/sego-server/internal/config"
	segodb "github.com/nadeemmhdm/sego-server/internal/db"
	"github.com/nadeemmhdm/sego-server/internal/storage"
	"golang.org/x/term"
)

func main() {
	cfgPath := flag.String("config", "sego.json", "config path")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		return
	}

	cfg, err := config.Ensure(*cfgPath)
	fatalIf(err)

	// Configuration-only commands should not require opening the database.
	if len(args) == 2 && args[0] == "mode" {
		switch args[1] {
		case "show":
			fmt.Printf("Mode: %s\nUpload concurrency: %d\n", config.Mode(cfg), cfg.UploadConcurrency)
			return
		case config.ModeLowResource, config.ModeStandard:
			fatalIf(config.ApplyMode(&cfg, args[1]))
			fatalIf(config.Save(*cfgPath, cfg))
			fmt.Printf("[OK] mode set to %s\n", args[1])
			return
		}
	}
	if len(args) == 3 && args[0] == "mode" && args[1] == "set" {
		fatalIf(config.ApplyMode(&cfg, args[2]))
		fatalIf(config.Save(*cfgPath, cfg))
		fmt.Printf("[OK] mode set to %s\n", args[2])
		return
	}

	d, err := segodb.Open(cfg.DatabasePath)
	fatalIf(err)
	defer d.Close()

	switch strings.Join(args, " ") {
	case "status":
		fmt.Printf("SEGO Server config OK\nListen: %s\nDatabase: %s\nRoots: %d\nMode: %s\n", cfg.ListenAddr, cfg.DatabasePath, len(cfg.StorageRoots), config.Mode(cfg))
	case "doctor":
		doctor(cfg)
	case "storage list":
		for _, r := range cfg.StorageRoots {
			fmt.Printf("%s\t%s\t%s\n", r.ID, r.Name, r.Path)
		}
	case "admin create":
		createAdmin(d)
	default:
		usage()
	}
}

func createAdmin(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Admin username: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	if len(user) < 3 {
		fatalIf(fmt.Errorf("username must be at least 3 characters"))
	}
	p1 := readSecret("Password (12+ chars): ")
	p2 := readSecret("Confirm password: ")
	if p1 != p2 {
		fatalIf(fmt.Errorf("passwords do not match"))
	}
	hash, err := auth.HashPassword(p1)
	fatalIf(err)
	_, err = d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES(?,?,?,?)`, user, hash, "super_admin", time.Now().UTC().Format(time.RFC3339))
	fatalIf(err)
	fmt.Println("[OK] Super Admin created")
}

func usage() {
	fmt.Println(`SEGO Server CLI
Usage:
  sego --config sego.json status
  sego --config sego.json doctor
  sego --config sego.json storage list
  sego --config sego.json admin create
  sego --config sego.json mode show
  sego --config sego.json mode set low-resource
  sego --config sego.json mode set standard`)
}

func doctor(cfg config.Config) {
	roots := make([]storage.Root, 0, len(cfg.StorageRoots))
	for _, r := range cfg.StorageRoots {
		roots = append(roots, storage.Root{ID: r.ID, Name: r.Name, Path: r.Path})
	}
	_, err := storage.New(roots)
	fatalIf(err)
	fmt.Println("[OK] configuration")
	fmt.Println("[OK] database path configured")
	fmt.Println("[OK] storage roots accessible")
	fmt.Printf("[INFO] mode: %s\n", config.Mode(cfg))
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "SEGO:", err)
		os.Exit(1)
	}
}

func readSecret(prompt string) string {
	fmt.Print(prompt)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err == nil {
			return string(b)
		}
	}
	s, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(s)
}
