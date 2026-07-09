# cyberpower-ups-exporter

[English README](README.en.md)

Prometheus exporter для ИБП CyberPower, которые управляются через
[PowerPanel for Linux](https://www.cyberpower.com/ru/ru/product/sku/powerpanel_for_linux).

Экспортёр запускает локальную команду `pwrstat`, разбирает статус ИБП и отдаёт
метрики через официальную Go-библиотеку Prometheus.

## Статус

Релиз 1.0. Текущая реализация рассчитана на распространённый формат вывода
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
- команда PowerPanel: `pwrstat -status`
- timeout сбора: `5s`

Пример запуска с явными параметрами:

```sh
./cyberpower-ups-exporter \
  --listen-address=:9833 \
  --metrics-path=/metrics \
  --pwrstat-command="pwrstat -status" \
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
PowerPanel for Linux, а команда `pwrstat -status` должна работать от имени
пользователя, под которым запущен экспортёр.
