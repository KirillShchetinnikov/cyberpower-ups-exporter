# cyberpower-ups-exporter

Prometheus exporter for CyberPower UPS devices managed by
[PowerPanel for Linux](https://www.cyberpower.com/ru/ru/product/sku/powerpanel_for_linux).

The exporter runs the local `pwrstat -status`, `pwrstat -config`, and
`pwrstat -version` commands, parses UPS status and PowerPanel daemon settings,
and exposes the values through the official Prometheus Go client library.

## Status

Release 1.1. The current implementation targets the common `pwrstat -status`
output format and uses `github.com/prometheus/client_golang` for collection and
HTTP exposition.

## Metrics

| Metric | Description |
| --- | --- |
| `cyberpower_ups_info` | Static UPS information labels. |
| `cyberpower_ups_scrape_success` | `1` when the last `pwrstat` scrape succeeded. |
| `cyberpower_ups_last_scrape_timestamp_seconds` | Unix timestamp of the scrape. |
| `cyberpower_ups_on_battery` | `1` when UPS power source looks like battery power. |
| `cyberpower_ups_battery_capacity_percent` | Battery charge percentage. |
| `cyberpower_ups_remaining_runtime_seconds` | Estimated remaining runtime. |
| `cyberpower_ups_utility_voltage_volts` | Utility input voltage. |
| `cyberpower_ups_output_voltage_volts` | UPS output voltage. |
| `cyberpower_ups_load_watts` | Current load in watts. |
| `cyberpower_ups_load_percent` | Current load percentage. |
| `cyberpower_ups_pwrstat_info` | Installed `pwrstat` utility version. |
| `cyberpower_ups_config_scrape_success` | `1` when the last `pwrstat -config` scrape succeeded. |
| `cyberpower_ups_version_scrape_success` | `1` when the last `pwrstat -version` scrape succeeded. |
| `cyberpower_ups_alarm_enabled` | Whether UPS alarm is enabled in PowerPanel settings. |
| `cyberpower_ups_hibernate_enabled` | Whether hibernate is enabled instead of shutdown. |
| `cyberpower_ups_cloud_enabled` | Whether CyberPower cloud integration is enabled. |
| `cyberpower_ups_power_failure_delay_seconds` | Delay before power failure action execution. |
| `cyberpower_ups_power_failure_script_enabled` | Whether script execution is enabled for power failure events. |
| `cyberpower_ups_power_failure_action_info` | Labels with the power failure script path. |
| `cyberpower_ups_power_failure_command_duration_seconds` | Configured power failure script execution duration. |
| `cyberpower_ups_power_failure_shutdown_enabled` | Whether system shutdown is enabled for power failure events. |
| `cyberpower_ups_low_battery_runtime_threshold_seconds` | Remaining runtime threshold for low battery events. |
| `cyberpower_ups_low_battery_capacity_threshold_percent` | Battery capacity threshold for low battery events. |
| `cyberpower_ups_low_battery_script_enabled` | Whether script execution is enabled for low battery events. |
| `cyberpower_ups_low_battery_action_info` | Labels with the low battery script path. |
| `cyberpower_ups_low_battery_command_duration_seconds` | Configured low battery script execution duration. |
| `cyberpower_ups_low_battery_shutdown_enabled` | Whether system shutdown is enabled for low battery events. |

## Build

```sh
go build ./cmd/cyberpower-ups-exporter
```

Requires Go 1.23 or newer.

## Run

```sh
./cyberpower-ups-exporter
```

Defaults:

- listen address: `:9833`
- metrics path: `/metrics`
- PowerPanel status command: `pwrstat -status`
- PowerPanel config command: `pwrstat -config`
- PowerPanel version command: `pwrstat -version`
- scrape timeout: `5s`

Custom command example:

```sh
./cyberpower-ups-exporter \
  --listen-address=:9833 \
  --metrics-path=/metrics \
  --pwrstat-command="pwrstat -status" \
  --pwrstat-config-command="pwrstat -config" \
  --pwrstat-version-command="pwrstat -version" \
  --pwrstat-timeout=5s
```

Scrape it:

```sh
curl http://localhost:9833/metrics
```

## Grafana

A ready-to-import dashboard is available at
[`dashboards/grafana/cyberpower-ups-exporter.json`](dashboards/grafana/cyberpower-ups-exporter.json).

Import it in Grafana and select the Prometheus datasource used to scrape this
exporter.

## Development

```sh
go test ./...
```

## Notes

The exporter does not talk to UPS hardware directly. PowerPanel for Linux must
be installed and `pwrstat -status`, `pwrstat -config`, and `pwrstat -version`
must work for the same user that runs the exporter.
