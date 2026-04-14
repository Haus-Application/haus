// Command validate-devices runs Buster + GOB Q&A validation against every KB
// device and writes JSON reports to disk. Ships as a CLI so you can run a full
// pass on demand, CI, or in-process from the HTTP server.
//
// Usage:
//
//	validate-devices [-kb docs/devices] [-out validation] [-concurrency 8] [-only SLUG]
//
// The API key is read from ANTHROPIC_API_KEY, or baked in at release time via
// -ldflags "-X main.defaultAPIKey=...".
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coalson/haus/internal/kb"
	"github.com/coalson/haus/internal/validation"
)

// defaultAPIKey is set at build time via -ldflags for release builds.
var defaultAPIKey string

func main() {
	var (
		kbDir       = flag.String("kb", "docs/devices", "Path to the device knowledge base directory")
		outDir      = flag.String("out", "validation", "Output directory for JSON reports")
		concurrency = flag.Int("concurrency", 8, "Parallel device validations")
		only        = flag.String("only", "", "Run just one device by slug (skip summary write)")
		model       = flag.String("model", "claude-haiku-4-5-20251001", "Claude model")
	)
	flag.Parse()

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		apiKey = defaultAPIKey
	}
	if apiKey == "" {
		log.Fatal("[validate] No API key. Set ANTHROPIC_API_KEY or build with -ldflags main.defaultAPIKey=...")
	}

	log.Printf("[validate] Loading KB from %s", *kbDir)
	catalog, err := kb.Load(*kbDir)
	if err != nil {
		log.Fatalf("[validate] Failed to load KB: %v", err)
	}

	runner := validation.NewRunner(catalog, apiKey, *model, *outDir, *concurrency)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Println("[validate] Interrupt received, cancelling...")
		cancel()
	}()

	if *only != "" {
		if err := os.MkdirAll(fmt.Sprintf("%s/devices", *outDir), 0o755); err != nil {
			log.Fatalf("[validate] mkdir: %v", err)
		}
		rep, err := runner.RunOne(ctx, *only)
		if err != nil {
			log.Fatalf("[validate] %s: %v", *only, err)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(rep)
		log.Printf("[validate] %s: %d%% (B %d/%d G %d/%d)", rep.Slug, rep.Score.TotalPct, rep.Score.Buster, rep.Score.BusterMax, rep.Score.GOB, rep.Score.GOBMax)
		return
	}

	log.Printf("[validate] Running %d devices with concurrency %d (model=%s)", len(catalog.All), *concurrency, *model)
	sum, err := runner.RunAll(ctx, nil)
	if err != nil {
		log.Fatalf("[validate] run failed: %v", err)
	}
	log.Printf("[validate] Done. %d devices, avg %d%%, %d passing, %d warning, %d failing.", sum.TotalDevices, sum.AvgScore, sum.Passing, sum.Warning, len(sum.Failing))
}
