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
	alarm := false
	hibernate := false
	cloud := false
	powerFailureDelay := 30.0
	powerFailureScript := true
	powerFailureDuration := 0.0
	powerFailureShutdown := true
	lowBatteryRuntime := 300.0
	lowBatteryCapacity := 35.0
	lowBatteryScript := true
	lowBatteryDuration := 0.0
	lowBatteryShutdown := true
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
			Config: PwrstatConfig{
				AlarmEnabled:     &alarm,
				HibernateEnabled: &hibernate,
				CloudEnabled:     &cloud,
				PowerFailure: PowerFailureConfig{
					DelaySeconds:          &powerFailureDelay,
					ScriptEnabled:         &powerFailureScript,
					ScriptPath:            "/etc/pwrstatd-powerfail.sh",
					CommandDurationSecond: &powerFailureDuration,
					ShutdownEnabled:       &powerFailureShutdown,
				},
				LowBattery: LowBatteryConfig{
					RuntimeThresholdSeconds:  &lowBatteryRuntime,
					CapacityThresholdPercent: &lowBatteryCapacity,
					ScriptEnabled:            &lowBatteryScript,
					ScriptPath:               "/etc/pwrstatd-lowbatt.sh",
					CommandDurationSecond:    &lowBatteryDuration,
					ShutdownEnabled:          &lowBatteryShutdown,
				},
			},
			PwrstatVersion: "1.4.2",
			When:           time.Unix(1000, 0),
		},
	}))

	expected := `
# HELP cyberpower_ups_alarm_enabled Whether the UPS alarm is enabled in PowerPanel daemon configuration.
# TYPE cyberpower_ups_alarm_enabled gauge
cyberpower_ups_alarm_enabled 0
# HELP cyberpower_ups_battery_capacity_percent UPS battery charge percentage.
# TYPE cyberpower_ups_battery_capacity_percent gauge
cyberpower_ups_battery_capacity_percent 97
# HELP cyberpower_ups_cloud_enabled Whether CyberPower cloud integration is enabled in PowerPanel daemon configuration.
# TYPE cyberpower_ups_cloud_enabled gauge
cyberpower_ups_cloud_enabled 0
# HELP cyberpower_ups_config_scrape_success Whether the last pwrstat config scrape succeeded.
# TYPE cyberpower_ups_config_scrape_success gauge
cyberpower_ups_config_scrape_success 1
# HELP cyberpower_ups_hibernate_enabled Whether hibernate is enabled instead of shutdown in PowerPanel daemon configuration.
# TYPE cyberpower_ups_hibernate_enabled gauge
cyberpower_ups_hibernate_enabled 0
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
# HELP cyberpower_ups_power_failure_action_info Power failure action configuration labels.
# TYPE cyberpower_ups_power_failure_action_info gauge
cyberpower_ups_power_failure_action_info{script_path="/etc/pwrstatd-powerfail.sh"} 1
# HELP cyberpower_ups_power_failure_command_duration_seconds Configured power failure script execution duration.
# TYPE cyberpower_ups_power_failure_command_duration_seconds gauge
cyberpower_ups_power_failure_command_duration_seconds 0
# HELP cyberpower_ups_power_failure_delay_seconds Delay before power failure action execution.
# TYPE cyberpower_ups_power_failure_delay_seconds gauge
cyberpower_ups_power_failure_delay_seconds 30
# HELP cyberpower_ups_power_failure_script_enabled Whether script execution is enabled for power failure events.
# TYPE cyberpower_ups_power_failure_script_enabled gauge
cyberpower_ups_power_failure_script_enabled 1
# HELP cyberpower_ups_power_failure_shutdown_enabled Whether system shutdown is enabled for power failure events.
# TYPE cyberpower_ups_power_failure_shutdown_enabled gauge
cyberpower_ups_power_failure_shutdown_enabled 1
# HELP cyberpower_ups_pwrstat_info Static pwrstat program information.
# TYPE cyberpower_ups_pwrstat_info gauge
cyberpower_ups_pwrstat_info{version="1.4.2"} 1
# HELP cyberpower_ups_remaining_runtime_seconds Estimated remaining battery runtime in seconds.
# TYPE cyberpower_ups_remaining_runtime_seconds gauge
cyberpower_ups_remaining_runtime_seconds 1800
# HELP cyberpower_ups_scrape_success Whether the last pwrstat status scrape succeeded.
# TYPE cyberpower_ups_scrape_success gauge
cyberpower_ups_scrape_success 1
# HELP cyberpower_ups_low_battery_action_info Low battery action configuration labels.
# TYPE cyberpower_ups_low_battery_action_info gauge
cyberpower_ups_low_battery_action_info{script_path="/etc/pwrstatd-lowbatt.sh"} 1
# HELP cyberpower_ups_low_battery_capacity_threshold_percent Battery capacity threshold used to identify low battery events.
# TYPE cyberpower_ups_low_battery_capacity_threshold_percent gauge
cyberpower_ups_low_battery_capacity_threshold_percent 35
# HELP cyberpower_ups_low_battery_command_duration_seconds Configured low battery script execution duration.
# TYPE cyberpower_ups_low_battery_command_duration_seconds gauge
cyberpower_ups_low_battery_command_duration_seconds 0
# HELP cyberpower_ups_low_battery_runtime_threshold_seconds Remaining runtime threshold used to identify low battery events.
# TYPE cyberpower_ups_low_battery_runtime_threshold_seconds gauge
cyberpower_ups_low_battery_runtime_threshold_seconds 300
# HELP cyberpower_ups_low_battery_script_enabled Whether script execution is enabled for low battery events.
# TYPE cyberpower_ups_low_battery_script_enabled gauge
cyberpower_ups_low_battery_script_enabled 1
# HELP cyberpower_ups_low_battery_shutdown_enabled Whether system shutdown is enabled for low battery events.
# TYPE cyberpower_ups_low_battery_shutdown_enabled gauge
cyberpower_ups_low_battery_shutdown_enabled 1
# HELP cyberpower_ups_version_scrape_success Whether the last pwrstat version scrape succeeded.
# TYPE cyberpower_ups_version_scrape_success gauge
cyberpower_ups_version_scrape_success 1
`

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expected),
		"cyberpower_ups_alarm_enabled",
		"cyberpower_ups_battery_capacity_percent",
		"cyberpower_ups_cloud_enabled",
		"cyberpower_ups_config_scrape_success",
		"cyberpower_ups_hibernate_enabled",
		"cyberpower_ups_info",
		"cyberpower_ups_last_scrape_timestamp_seconds",
		"cyberpower_ups_load_percent",
		"cyberpower_ups_on_battery",
		"cyberpower_ups_power_failure_action_info",
		"cyberpower_ups_power_failure_command_duration_seconds",
		"cyberpower_ups_power_failure_delay_seconds",
		"cyberpower_ups_power_failure_script_enabled",
		"cyberpower_ups_power_failure_shutdown_enabled",
		"cyberpower_ups_pwrstat_info",
		"cyberpower_ups_remaining_runtime_seconds",
		"cyberpower_ups_scrape_success",
		"cyberpower_ups_low_battery_action_info",
		"cyberpower_ups_low_battery_capacity_threshold_percent",
		"cyberpower_ups_low_battery_command_duration_seconds",
		"cyberpower_ups_low_battery_runtime_threshold_seconds",
		"cyberpower_ups_low_battery_script_enabled",
		"cyberpower_ups_low_battery_shutdown_enabled",
		"cyberpower_ups_version_scrape_success",
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
# HELP cyberpower_ups_scrape_success Whether the last pwrstat status scrape succeeded.
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
