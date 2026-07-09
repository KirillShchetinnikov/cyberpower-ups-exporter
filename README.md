# cyberpower-ups-exporter

[English README](README.en.md)

Prometheus exporter для ИБП CyberPower, которые управляются через
[PowerPanel for Linux](https://www.cyberpower.com/ru/ru/product/sku/powerpanel_for_linux).

Экспортёр запускает локальные команды `pwrstat -status`, `pwrstat -config` и
`pwrstat -version`, разбирает статус ИБП и настройки PowerPanel daemon, затем
отдаёт метрики через официальную Go-библиотеку Prometheus.

## Статус

Релиз 1.1. Текущая реализация рассчитана на распространённый формат вывода
`pwrstat -status` и использует `github.com/prometheus/client_golang` для сбора и
HTTP-публикации метрик.

## Метрики

| Метрика | Описание |
| --- | --- |
| `cyberpower_ups_info` | Информационные labels по ИБП. |
| `cyberpower_ups_scrape_success` | `1`, если последний запуск `pwrstat` завершился успешно. |
| `cyberpower_ups_last_scrape_timestamp_seconds` | Unix timestamp последнего сбора. |
| `cyberpower_ups_on_battery` | `1`, если ИБП работает от батареи. |
| `cyberpower_ups_battery_capacity_percent` | Заряд батареи в процентах. |
| `cyberpower_ups_remaining_runtime_seconds` | Оценка оставшегося времени работы в секундах. |
| `cyberpower_ups_utility_voltage_volts` | Входное напряжение сети. |
| `cyberpower_ups_output_voltage_volts` | Выходное напряжение ИБП. |
| `cyberpower_ups_load_watts` | Текущая нагрузка в ваттах. |
| `cyberpower_ups_load_percent` | Текущая нагрузка в процентах. |
| `cyberpower_ups_pwrstat_info` | Версия установленной утилиты `pwrstat`. |
| `cyberpower_ups_config_scrape_success` | `1`, если последний запуск `pwrstat -config` завершился успешно. |
| `cyberpower_ups_version_scrape_success` | `1`, если последний запуск `pwrstat -version` завершился успешно. |
| `cyberpower_ups_alarm_enabled` | Включена ли звуковая сигнализация ИБП в настройках PowerPanel. |
| `cyberpower_ups_hibernate_enabled` | Включён ли hibernate вместо shutdown. |
| `cyberpower_ups_cloud_enabled` | Включена ли CyberPower cloud-интеграция. |
| `cyberpower_ups_power_failure_delay_seconds` | Задержка перед действием при пропадании питания. |
| `cyberpower_ups_power_failure_script_enabled` | Включён ли запуск скрипта при пропадании питания. |
| `cyberpower_ups_power_failure_action_info` | Labels с путём скрипта для события пропадания питания. |
| `cyberpower_ups_power_failure_command_duration_seconds` | Длительность выполнения скрипта при пропадании питания. |
| `cyberpower_ups_power_failure_shutdown_enabled` | Включён ли shutdown системы при пропадании питания. |
| `cyberpower_ups_low_battery_runtime_threshold_seconds` | Порог remaining runtime для события low battery. |
| `cyberpower_ups_low_battery_capacity_threshold_percent` | Порог заряда батареи для события low battery. |
| `cyberpower_ups_low_battery_script_enabled` | Включён ли запуск скрипта при low battery. |
| `cyberpower_ups_low_battery_action_info` | Labels с путём скрипта для события low battery. |
| `cyberpower_ups_low_battery_command_duration_seconds` | Длительность выполнения скрипта при low battery. |
| `cyberpower_ups_low_battery_shutdown_enabled` | Включён ли shutdown системы при low battery. |

## Сборка

```sh
go build ./cmd/cyberpower-ups-exporter
```

Требуется Go 1.23 или новее.

## Запуск

```sh
./cyberpower-ups-exporter
```

Значения по умолчанию:

- адрес прослушивания: `:9833`
- путь метрик: `/metrics`
- команда статуса PowerPanel: `pwrstat -status`
- команда настроек PowerPanel: `pwrstat -config`
- команда версии PowerPanel: `pwrstat -version`
- timeout сбора: `5s`

Пример запуска с явными параметрами:

```sh
./cyberpower-ups-exporter \
  --listen-address=:9833 \
  --metrics-path=/metrics \
  --pwrstat-command="pwrstat -status" \
  --pwrstat-config-command="pwrstat -config" \
  --pwrstat-version-command="pwrstat -version" \
  --pwrstat-timeout=5s
```

Проверка метрик:

```sh
curl http://localhost:9833/metrics
```

## Grafana

Готовый dashboard для Grafana находится в
[`dashboards/grafana/cyberpower-ups-exporter.json`](dashboards/grafana/cyberpower-ups-exporter.json).

Импортируйте JSON в Grafana и выберите Prometheus datasource, который собирает
метрики этого экспортёра.

## Разработка

```sh
go test ./...
```

## Примечания

Экспортёр не обращается к ИБП напрямую. На хосте должен быть установлен
PowerPanel for Linux, а команды `pwrstat -status`, `pwrstat -config` и
`pwrstat -version` должны работать от имени пользователя, под которым запущен
экспортёр.
