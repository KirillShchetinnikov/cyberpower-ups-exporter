package exporter

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ScrapeResult struct {
	Status UPSStatus
	When   time.Time
	Err    error
}

type StatusScraper interface {
	Scrape(ctx context.Context) ScrapeResult
}

type Collector struct {
	scraper StatusScraper

	infoDesc             *prometheus.Desc
	scrapeSuccessDesc    *prometheus.Desc
	lastScrapeDesc       *prometheus.Desc
	onBatteryDesc        *prometheus.Desc
	batteryCapacityDesc  *prometheus.Desc
	remainingRuntimeDesc *prometheus.Desc
	utilityVoltageDesc   *prometheus.Desc
	outputVoltageDesc    *prometheus.Desc
	loadWattsDesc        *prometheus.Desc
	loadPercentDesc      *prometheus.Desc
}

func NewCollector(scraper StatusScraper) *Collector {
	return &Collector{
		scraper: scraper,
		infoDesc: prometheus.NewDesc(
			"cyberpower_ups_info",
			"Static UPS information reported by pwrstat.",
			[]string{"model", "firmware", "state", "power_supply"},
			nil,
		),
		scrapeSuccessDesc: prometheus.NewDesc(
			"cyberpower_ups_scrape_success",
			"Whether the last pwrstat scrape succeeded.",
			nil,
			nil,
		),
		lastScrapeDesc: prometheus.NewDesc(
			"cyberpower_ups_last_scrape_timestamp_seconds",
			"Unix timestamp of the last scrape.",
			nil,
			nil,
		),
		onBatteryDesc: prometheus.NewDesc(
			"cyberpower_ups_on_battery",
			"Whether the UPS is currently supplying power from battery.",
			nil,
			nil,
		),
		batteryCapacityDesc: prometheus.NewDesc(
			"cyberpower_ups_battery_capacity_percent",
			"UPS battery charge percentage.",
			nil,
			nil,
		),
		remainingRuntimeDesc: prometheus.NewDesc(
			"cyberpower_ups_remaining_runtime_seconds",
			"Estimated remaining battery runtime in seconds.",
			nil,
			nil,
		),
		utilityVoltageDesc: prometheus.NewDesc(
			"cyberpower_ups_utility_voltage_volts",
			"Utility input voltage in volts.",
			nil,
			nil,
		),
		outputVoltageDesc: prometheus.NewDesc(
			"cyberpower_ups_output_voltage_volts",
			"UPS output voltage in volts.",
			nil,
			nil,
		),
		loadWattsDesc: prometheus.NewDesc(
			"cyberpower_ups_load_watts",
			"UPS load in watts.",
			nil,
			nil,
		),
		loadPercentDesc: prometheus.NewDesc(
			"cyberpower_ups_load_percent",
			"UPS load percentage.",
			nil,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.scrapeSuccessDesc
	ch <- c.lastScrapeDesc
	ch <- c.onBatteryDesc
	ch <- c.batteryCapacityDesc
	ch <- c.remainingRuntimeDesc
	ch <- c.utilityVoltageDesc
	ch <- c.outputVoltageDesc
	ch <- c.loadWattsDesc
	ch <- c.loadPercentDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	result := c.scraper.Scrape(context.Background())
	when := result.When
	if when.IsZero() {
		when = time.Now()
	}

	ch <- prometheus.MustNewConstMetric(
		c.infoDesc,
		prometheus.GaugeValue,
		1,
		result.Status.ModelName,
		result.Status.FirmwareNumber,
		result.Status.State,
		result.Status.PowerSupply,
	)
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, boolFloat(result.Err == nil))
	ch <- prometheus.MustNewConstMetric(c.lastScrapeDesc, prometheus.GaugeValue, float64(when.Unix()))
	ch <- prometheus.MustNewConstMetric(c.onBatteryDesc, prometheus.GaugeValue, boolFloat(isOnBattery(result.Status)))

	collectOptional(ch, c.batteryCapacityDesc, result.Status.BatteryCapacityPercent)
	collectOptional(ch, c.remainingRuntimeDesc, result.Status.RemainingRuntimeSeconds)
	collectOptional(ch, c.utilityVoltageDesc, result.Status.UtilityVoltageVolts)
	collectOptional(ch, c.outputVoltageDesc, result.Status.OutputVoltageVolts)
	collectOptional(ch, c.loadWattsDesc, result.Status.LoadWatts)
	collectOptional(ch, c.loadPercentDesc, result.Status.LoadPercent)
}

func collectOptional(ch chan<- prometheus.Metric, desc *prometheus.Desc, value *float64) {
	if value == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, *value)
}

func boolFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func isOnBattery(status UPSStatus) bool {
	source := strings.ToLower(status.PowerSupply + " " + status.State)
	return strings.Contains(source, "battery")
}
