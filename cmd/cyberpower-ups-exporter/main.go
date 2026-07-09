package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KirillShchetinnikov/cyberpower-ups-exporter/internal/exporter"
)

const version = "0.1.0"

func main() {
	var (
		listenAddress = flag.String("listen-address", ":9833", "HTTP listen address")
		metricsPath   = flag.String("metrics-path", "/metrics", "Prometheus metrics path")
		pwrstatCmd    = flag.String("pwrstat-command", "pwrstat -status", "PowerPanel command to execute without shell")
		pwrstatTO     = flag.Duration("pwrstat-timeout", 5*time.Second, "PowerPanel command timeout")
		showVersion   = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	scraper := exporter.NewScraper(*pwrstatCmd, *pwrstatTO)
	mux := http.NewServeMux()
	mux.Handle(*metricsPath, exporter.MetricsHandler(scraper))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	server := &http.Server{
		Addr:              *listenAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown failed: %v", err)
		}
	}()

	log.Printf("starting cyberpower-ups-exporter version=%s listen=%s metrics_path=%s command=%q", version, *listenAddress, *metricsPath, *pwrstatCmd)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
