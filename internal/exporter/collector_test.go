package exporter

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCollector(t *testing.T) {
	battery := 97.0
	runtime := 1800.0
	load := 15.0
	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(staticScraper{
		result: ScrapeResult{
			Status: UPSStatus{
				ModelName:               "CP1500",
				FirmwareNumber:          "CR021",
				State:                   "Normal",
				PowerSupply:             "Utility Power",
				BatteryCapacityPercent:  &battery,
				RemainingRuntimeSeconds: &runtime,
				LoadPercent:             &load,
			},
			When: time.Unix(1000, 0),
		},
	}))

	expected := `
# HELP cyberpower_ups_battery_capacity_percent UPS battery charge percentage.
# TYPE cyberpower_ups_battery_capacity_percent gauge
cyberpower_ups_battery_capacity_percent 97
# HELP cyberpower_ups_info Static UPS information reported by pwrstat.
# TYPE cyberpower_ups_info gauge
cyberpower_ups_info{firmware="CR021",model="CP1500",power_supply="Utility Power",state="Normal"} 1
# HELP cyberpower_ups_last_scrape_timestamp_seconds Unix timestamp of the last scrape.
# TYPE cyberpower_ups_last_scrape_timestamp_seconds gauge
cyberpower_ups_last_scrape_timestamp_seconds 1000
# HELP cyberpower_ups_load_percent UPS load percentage.
# TYPE cyberpower_ups_load_percent gauge
cyberpower_ups_load_percent 15
# HELP cyberpower_ups_on_battery Whether the UPS is currently supplying power from battery.
# TYPE cyberpower_ups_on_battery gauge
cyberpower_ups_on_battery 0
# HELP cyberpower_ups_remaining_runtime_seconds Estimated remaining battery runtime in seconds.
# TYPE cyberpower_ups_remaining_runtime_seconds gauge
cyberpower_ups_remaining_runtime_seconds 1800
# HELP cyberpower_ups_scrape_success Whether the last pwrstat scrape succeeded.
# TYPE cyberpower_ups_scrape_success gauge
cyberpower_ups_scrape_success 1
`

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expected),
		"cyberpower_ups_battery_capacity_percent",
		"cyberpower_ups_info",
		"cyberpower_ups_last_scrape_timestamp_seconds",
		"cyberpower_ups_load_percent",
		"cyberpower_ups_on_battery",
		"cyberpower_ups_remaining_runtime_seconds",
		"cyberpower_ups_scrape_success",
	); err != nil {
		t.Fatal(err)
	}
}

func TestCollectorWithScrapeError(t *testing.T) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(NewCollector(staticScraper{
		result: ScrapeResult{
			When: time.Unix(1000, 0),
			Err:  errors.New("boom"),
		},
	}))

	expected := `
# HELP cyberpower_ups_scrape_success Whether the last pwrstat scrape succeeded.
# TYPE cyberpower_ups_scrape_success gauge
cyberpower_ups_scrape_success 0
`

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expected), "cyberpower_ups_scrape_success"); err != nil {
		t.Fatal(err)
	}
}

type staticScraper struct {
	result ScrapeResult
}

func (s staticScraper) Scrape(_ context.Context) ScrapeResult {
	return s.result
}
