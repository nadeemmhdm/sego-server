package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/updater"
	"github.com/nadeemmhdm/sego-server/internal/version"
)

func main() {
	cfg := flag.String("config", defaultConfigPath(), "SEGO config path")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	changed, latest, err := updater.ApplyLatest(ctx, version.Version, *cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "SEGO updater:", err)
		os.Exit(1)
	}
	if changed {
		fmt.Printf("SEGO updated to %s\n", latest)
	} else if latest == "" {
		fmt.Println("No published SEGO release is available yet.")
	} else {
		fmt.Printf("SEGO %s is already current.\n", version.Version)
	}
}

func defaultConfigPath() string {
	if v := os.Getenv("SEGO_CONFIG"); v != "" { return v }
	if runtime.GOOS == "windows" {
		if pd := os.Getenv("ProgramData"); pd != "" { return filepath.Join(pd,"SEGO","sego.json") }
	}
	return "/etc/sego/sego.json"
}
