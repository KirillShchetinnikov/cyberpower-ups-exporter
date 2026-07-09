package exporter

import (
	"regexp"
	"strconv"
	"strings"
)

var keyValueLine = regexp.MustCompile(`^(.+?)(?:\.{2,}|:)\s*(.+)$`)

type UPSStatus struct {
	ModelName               string
	FirmwareNumber          string
	State                   string
	PowerSupply             string
	UtilityVoltageVolts     *float64
	OutputVoltageVolts      *float64
	BatteryCapacityPercent  *float64
	RemainingRuntimeSeconds *float64
	LoadWatts               *float64
	LoadPercent             *float64
}

func ParsePwrstat(output string) UPSStatus {
	var status UPSStatus

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := keyValueLine.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		key := normalizeKey(matches[1])
		value := strings.TrimSpace(matches[2])

		switch key {
		case "modelname":
			status.ModelName = value
		case "firmwarenumber":
			status.FirmwareNumber = value
		case "state":
			status.State = value
		case "powersupplyby", "powersupply":
			status.PowerSupply = value
		case "utilityvoltage":
			status.UtilityVoltageVolts = firstNumber(value)
		case "outputvoltage":
			status.OutputVoltageVolts = firstNumber(value)
		case "batterycapacity":
			status.BatteryCapacityPercent = firstNumber(value)
		case "remainingruntime":
			status.RemainingRuntimeSeconds = runtimeSeconds(value)
		case "load":
			status.LoadWatts = loadWatts(value)
			status.LoadPercent = loadPercent(value)
		}
	}

	return status
}

func normalizeKey(value string) string {
	value = strings.ToLower(value)
	var builder strings.Builder
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func firstNumber(value string) *float64 {
	fields := regexp.MustCompile(`[-+]?\d+(?:\.\d+)?`).FindStringSubmatch(value)
	if len(fields) == 0 {
		return nil
	}

	parsed, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func runtimeSeconds(value string) *float64 {
	number := firstNumber(value)
	if number == nil {
		return nil
	}

	lower := strings.ToLower(value)
	seconds := *number
	switch {
	case strings.Contains(lower, "hour"), strings.Contains(lower, " hr"), strings.Contains(lower, "hrs"):
		seconds *= 3600
	case strings.Contains(lower, "sec"):
		seconds *= 1
	default:
		seconds *= 60
	}
	return &seconds
}

func loadWatts(value string) *float64 {
	lower := strings.ToLower(value)
	wattIndex := strings.Index(lower, "watt")
	if wattIndex == -1 {
		return nil
	}
	return firstNumber(value[:wattIndex])
}

func loadPercent(value string) *float64 {
	percentIndex := strings.Index(value, "%")
	if percentIndex == -1 {
		return nil
	}
	prefix := value[:percentIndex]
	lastOpenParen := strings.LastIndex(prefix, "(")
	if lastOpenParen != -1 {
		prefix = prefix[lastOpenParen+1:]
	}
	return firstNumber(prefix)
}
