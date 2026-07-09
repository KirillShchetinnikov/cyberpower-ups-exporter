# cyberpower-ups-exporter

Prometheus exporter for CyberPower UPS devices managed by
[PowerPanel for Linux](https://www.cyberpower.com/ru/ru/product/sku/powerpanel_for_linux).

The exporter runs the local `pwrstat` CLI, parses UPS status output, and exposes
the values through the official Prometheus Go client library.

## Status

MVP. The current implementation targets the common `pwrstat -status` output
format and uses `github.com/prometheus/client_golang` for collection and HTTP
exposition.

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
- PowerPanel command: `pwrstat -status`
- scrape timeout: `5s`

Custom command example:

```sh
./cyberpower-ups-exporter \
  --listen-address=:9833 \
  --metrics-path=/metrics \
  --pwrstat-command="pwrstat -status" \
  --pwrstat-timeout=5s
```

Scrape it:

```sh
curl http://localhost:9833/metrics
```

## Development

```sh
go test ./...
```

## Notes

The exporter does not talk to UPS hardware directly. PowerPanel for Linux must
be installed and `pwrstat -status` must work for the same user that runs the
exporter.
