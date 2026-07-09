package exporter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Scraper struct {
	command string
	timeout time.Duration
}

func NewScraper(command string, timeout time.Duration) Scraper {
	return Scraper{
		command: strings.TrimSpace(command),
		timeout: timeout,
	}
}

func (s Scraper) Scrape(ctx context.Context) ScrapeResult {
	when := time.Now()
	output, err := s.run(ctx)
	if err != nil {
		return ScrapeResult{When: when, Err: err}
	}

	return ScrapeResult{
		Status: ParsePwrstat(output),
		When:   when,
	}
}

func (s Scraper) run(ctx context.Context) (string, error) {
	parts := strings.Fields(s.command)
	if len(parts) == 0 {
		return "", errors.New("pwrstat command is empty")
	}

	timeout := s.timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(commandCtx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if commandCtx.Err() != nil {
		return "", fmt.Errorf("pwrstat command timed out after %s", timeout)
	}
	if err != nil {
		return "", fmt.Errorf("pwrstat command failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func MetricsHandler(scraper StatusScraper) http.Handler {
	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(scraper))
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
