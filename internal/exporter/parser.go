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

type PwrstatConfig struct {
	AlarmEnabled     *bool
	HibernateEnabled *bool
	CloudEnabled     *bool
	PowerFailure     PowerFailureConfig
	LowBattery       LowBatteryConfig
}

type PowerFailureConfig struct {
	DelaySeconds          *float64
	ScriptEnabled         *bool
	ScriptPath            string
	CommandDurationSecond *float64
	ShutdownEnabled       *bool
}

type LowBatteryConfig struct {
	RuntimeThresholdSeconds  *float64
	CapacityThresholdPercent *float64
	ScriptEnabled            *bool
	ScriptPath               string
	CommandDurationSecond    *float64
	ShutdownEnabled          *bool
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

func ParsePwrstatConfig(output string) PwrstatConfig {
	var config PwrstatConfig
	section := ""

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		normalizedLine := normalizeKey(line)
		switch normalizedLine {
		case "actionforpowerfailure":
			section = "power_failure"
			continue
		case "actionforbatterylow":
			section = "battery_low"
			continue
		}

		matches := keyValueLine.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		key := normalizeKey(matches[1])
		value := strings.TrimSpace(matches[2])

		switch section {
		case "power_failure":
			parsePowerFailureConfig(&config.PowerFailure, key, value)
		case "battery_low":
			parseLowBatteryConfig(&config.LowBattery, key, value)
		default:
			switch key {
			case "alarm":
				config.AlarmEnabled = onOff(value)
			case "hibernate":
				config.HibernateEnabled = onOff(value)
			case "cloud":
				config.CloudEnabled = onOff(value)
			}
		}
	}

	return config
}

func ParsePwrstatVersion(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "pwrstat version") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return fields[len(fields)-1]
			}
		}
	}
	return ""
}

func parsePowerFailureConfig(config *PowerFailureConfig, key string, value string) {
	switch key {
	case "delaytimesincepowerfailure":
		config.DelaySeconds = firstNumber(value)
	case "runscriptcommand":
		config.ScriptEnabled = onOff(value)
	case "pathofscriptcommand":
		config.ScriptPath = value
	case "durationofcommandrunning":
		config.CommandDurationSecond = firstNumber(value)
	case "enableshutdownsystem":
		config.ShutdownEnabled = onOff(value)
	}
}

func parseLowBatteryConfig(config *LowBatteryConfig, key string, value string) {
	switch key {
	case "remainingruntimethreshold":
		config.RuntimeThresholdSeconds = firstNumber(value)
	case "batterycapacitythreshold":
		config.CapacityThresholdPercent = firstNumber(value)
	case "runscriptcommand":
		config.ScriptEnabled = onOff(value)
	case "pathofcommand":
		config.ScriptPath = value
	case "durationofcommandrunning":
		config.CommandDurationSecond = firstNumber(value)
	case "enableshutdownsystem":
		config.ShutdownEnabled = onOff(value)
	}
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

func onOff(value string) *bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "on", "enable", "enabled", "true", "yes", "1":
		enabled := true
		return &enabled
	case "off", "disable", "disabled", "false", "no", "0":
		enabled := false
		return &enabled
	default:
		return nil
	}
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
