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
	statusCommand  string
	configCommand  string
	versionCommand string
	timeout        time.Duration
}

func NewScraper(statusCommand string, configCommand string, versionCommand string, timeout time.Duration) Scraper {
	return Scraper{
		statusCommand:  strings.TrimSpace(statusCommand),
		configCommand:  strings.TrimSpace(configCommand),
		versionCommand: strings.TrimSpace(versionCommand),
		timeout:        timeout,
	}
}

func (s Scraper) Scrape(ctx context.Context) ScrapeResult {
	when := time.Now()
	result := ScrapeResult{When: when}

	statusOutput, err := s.run(ctx, s.statusCommand)
	if err != nil {
		result.Err = err
	} else {
		result.Status = ParsePwrstat(statusOutput)
	}

	configOutput, err := s.run(ctx, s.configCommand)
	if err != nil {
		result.ConfigErr = err
	} else {
		result.Config = ParsePwrstatConfig(configOutput)
	}

	versionOutput, err := s.run(ctx, s.versionCommand)
	if err != nil {
		result.VersionErr = err
	} else {
		result.PwrstatVersion = ParsePwrstatVersion(versionOutput)
	}

	return result
}

func (s Scraper) run(ctx context.Context, command string) (string, error) {
	parts := strings.Fields(command)
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
