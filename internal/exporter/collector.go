package exporter

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ScrapeResult struct {
	Status         UPSStatus
	Config         PwrstatConfig
	PwrstatVersion string
	When           time.Time
	Err            error
	ConfigErr      error
	VersionErr     error
}

type StatusScraper interface {
	Scrape(ctx context.Context) ScrapeResult
}

type Collector struct {
	scraper StatusScraper

	infoDesc             *prometheus.Desc
	pwrstatInfoDesc      *prometheus.Desc
	scrapeSuccessDesc    *prometheus.Desc
	configSuccessDesc    *prometheus.Desc
	versionSuccessDesc   *prometheus.Desc
	lastScrapeDesc       *prometheus.Desc
	onBatteryDesc        *prometheus.Desc
	batteryCapacityDesc  *prometheus.Desc
	remainingRuntimeDesc *prometheus.Desc
	utilityVoltageDesc   *prometheus.Desc
	outputVoltageDesc    *prometheus.Desc
	loadWattsDesc        *prometheus.Desc
	loadPercentDesc      *prometheus.Desc
	alarmEnabledDesc     *prometheus.Desc
	hibernateEnabledDesc *prometheus.Desc
	cloudEnabledDesc     *prometheus.Desc

	powerFailureDelayDesc           *prometheus.Desc
	powerFailureScriptEnabledDesc   *prometheus.Desc
	powerFailureActionInfoDesc      *prometheus.Desc
	powerFailureCommandDurationDesc *prometheus.Desc
	powerFailureShutdownEnabledDesc *prometheus.Desc

	lowBatteryRuntimeThresholdDesc  *prometheus.Desc
	lowBatteryCapacityThresholdDesc *prometheus.Desc
	lowBatteryScriptEnabledDesc     *prometheus.Desc
	lowBatteryActionInfoDesc        *prometheus.Desc
	lowBatteryCommandDurationDesc   *prometheus.Desc
	lowBatteryShutdownEnabledDesc   *prometheus.Desc
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
		pwrstatInfoDesc: prometheus.NewDesc(
			"cyberpower_ups_pwrstat_info",
			"Static pwrstat program information.",
			[]string{"version"},
			nil,
		),
		scrapeSuccessDesc: prometheus.NewDesc(
			"cyberpower_ups_scrape_success",
			"Whether the last pwrstat status scrape succeeded.",
			nil,
			nil,
		),
		configSuccessDesc: prometheus.NewDesc(
			"cyberpower_ups_config_scrape_success",
			"Whether the last pwrstat config scrape succeeded.",
			nil,
			nil,
		),
		versionSuccessDesc: prometheus.NewDesc(
			"cyberpower_ups_version_scrape_success",
			"Whether the last pwrstat version scrape succeeded.",
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
		alarmEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_alarm_enabled",
			"Whether the UPS alarm is enabled in PowerPanel daemon configuration.",
			nil,
			nil,
		),
		hibernateEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_hibernate_enabled",
			"Whether hibernate is enabled instead of shutdown in PowerPanel daemon configuration.",
			nil,
			nil,
		),
		cloudEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_cloud_enabled",
			"Whether CyberPower cloud integration is enabled in PowerPanel daemon configuration.",
			nil,
			nil,
		),
		powerFailureDelayDesc: prometheus.NewDesc(
			"cyberpower_ups_power_failure_delay_seconds",
			"Delay before power failure action execution.",
			nil,
			nil,
		),
		powerFailureScriptEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_power_failure_script_enabled",
			"Whether script execution is enabled for power failure events.",
			nil,
			nil,
		),
		powerFailureActionInfoDesc: prometheus.NewDesc(
			"cyberpower_ups_power_failure_action_info",
			"Power failure action configuration labels.",
			[]string{"script_path"},
			nil,
		),
		powerFailureCommandDurationDesc: prometheus.NewDesc(
			"cyberpower_ups_power_failure_command_duration_seconds",
			"Configured power failure script execution duration.",
			nil,
			nil,
		),
		powerFailureShutdownEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_power_failure_shutdown_enabled",
			"Whether system shutdown is enabled for power failure events.",
			nil,
			nil,
		),
		lowBatteryRuntimeThresholdDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_runtime_threshold_seconds",
			"Remaining runtime threshold used to identify low battery events.",
			nil,
			nil,
		),
		lowBatteryCapacityThresholdDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_capacity_threshold_percent",
			"Battery capacity threshold used to identify low battery events.",
			nil,
			nil,
		),
		lowBatteryScriptEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_script_enabled",
			"Whether script execution is enabled for low battery events.",
			nil,
			nil,
		),
		lowBatteryActionInfoDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_action_info",
			"Low battery action configuration labels.",
			[]string{"script_path"},
			nil,
		),
		lowBatteryCommandDurationDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_command_duration_seconds",
			"Configured low battery script execution duration.",
			nil,
			nil,
		),
		lowBatteryShutdownEnabledDesc: prometheus.NewDesc(
			"cyberpower_ups_low_battery_shutdown_enabled",
			"Whether system shutdown is enabled for low battery events.",
			nil,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.pwrstatInfoDesc
	ch <- c.scrapeSuccessDesc
	ch <- c.configSuccessDesc
	ch <- c.versionSuccessDesc
	ch <- c.lastScrapeDesc
	ch <- c.onBatteryDesc
	ch <- c.batteryCapacityDesc
	ch <- c.remainingRuntimeDesc
	ch <- c.utilityVoltageDesc
	ch <- c.outputVoltageDesc
	ch <- c.loadWattsDesc
	ch <- c.loadPercentDesc
	ch <- c.alarmEnabledDesc
	ch <- c.hibernateEnabledDesc
	ch <- c.cloudEnabledDesc
	ch <- c.powerFailureDelayDesc
	ch <- c.powerFailureScriptEnabledDesc
	ch <- c.powerFailureActionInfoDesc
	ch <- c.powerFailureCommandDurationDesc
	ch <- c.powerFailureShutdownEnabledDesc
	ch <- c.lowBatteryRuntimeThresholdDesc
	ch <- c.lowBatteryCapacityThresholdDesc
	ch <- c.lowBatteryScriptEnabledDesc
	ch <- c.lowBatteryActionInfoDesc
	ch <- c.lowBatteryCommandDurationDesc
	ch <- c.lowBatteryShutdownEnabledDesc
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
	ch <- prometheus.MustNewConstMetric(c.pwrstatInfoDesc, prometheus.GaugeValue, 1, result.PwrstatVersion)
	ch <- prometheus.MustNewConstMetric(c.scrapeSuccessDesc, prometheus.GaugeValue, boolFloat(result.Err == nil))
	ch <- prometheus.MustNewConstMetric(c.configSuccessDesc, prometheus.GaugeValue, boolFloat(result.ConfigErr == nil))
	ch <- prometheus.MustNewConstMetric(c.versionSuccessDesc, prometheus.GaugeValue, boolFloat(result.VersionErr == nil))
	ch <- prometheus.MustNewConstMetric(c.lastScrapeDesc, prometheus.GaugeValue, float64(when.Unix()))
	ch <- prometheus.MustNewConstMetric(c.onBatteryDesc, prometheus.GaugeValue, boolFloat(isOnBattery(result.Status)))

	collectOptional(ch, c.batteryCapacityDesc, result.Status.BatteryCapacityPercent)
	collectOptional(ch, c.remainingRuntimeDesc, result.Status.RemainingRuntimeSeconds)
	collectOptional(ch, c.utilityVoltageDesc, result.Status.UtilityVoltageVolts)
	collectOptional(ch, c.outputVoltageDesc, result.Status.OutputVoltageVolts)
	collectOptional(ch, c.loadWattsDesc, result.Status.LoadWatts)
	collectOptional(ch, c.loadPercentDesc, result.Status.LoadPercent)
	collectOptionalBool(ch, c.alarmEnabledDesc, result.Config.AlarmEnabled)
	collectOptionalBool(ch, c.hibernateEnabledDesc, result.Config.HibernateEnabled)
	collectOptionalBool(ch, c.cloudEnabledDesc, result.Config.CloudEnabled)

	collectOptional(ch, c.powerFailureDelayDesc, result.Config.PowerFailure.DelaySeconds)
	collectOptionalBool(ch, c.powerFailureScriptEnabledDesc, result.Config.PowerFailure.ScriptEnabled)
	collectInfo(ch, c.powerFailureActionInfoDesc, result.Config.PowerFailure.ScriptPath)
	collectOptional(ch, c.powerFailureCommandDurationDesc, result.Config.PowerFailure.CommandDurationSecond)
	collectOptionalBool(ch, c.powerFailureShutdownEnabledDesc, result.Config.PowerFailure.ShutdownEnabled)

	collectOptional(ch, c.lowBatteryRuntimeThresholdDesc, result.Config.LowBattery.RuntimeThresholdSeconds)
	collectOptional(ch, c.lowBatteryCapacityThresholdDesc, result.Config.LowBattery.CapacityThresholdPercent)
	collectOptionalBool(ch, c.lowBatteryScriptEnabledDesc, result.Config.LowBattery.ScriptEnabled)
	collectInfo(ch, c.lowBatteryActionInfoDesc, result.Config.LowBattery.ScriptPath)
	collectOptional(ch, c.lowBatteryCommandDurationDesc, result.Config.LowBattery.CommandDurationSecond)
	collectOptionalBool(ch, c.lowBatteryShutdownEnabledDesc, result.Config.LowBattery.ShutdownEnabled)
}

func collectOptional(ch chan<- prometheus.Metric, desc *prometheus.Desc, value *float64) {
	if value == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, *value)
}

func collectOptionalBool(ch chan<- prometheus.Metric, desc *prometheus.Desc, value *bool) {
	if value == nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, boolFloat(*value))
}

func collectInfo(ch chan<- prometheus.Metric, desc *prometheus.Desc, labelValue string) {
	if labelValue == "" {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1, labelValue)
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
