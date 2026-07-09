package exporter

import "testing"

func TestParsePwrstatStatus(t *testing.T) {
	output := `
The UPS information shows as following:

	Properties:
		Model Name................... CP1500EPFCLCD
		Firmware Number.............. CR021

	Current UPS status:
		State........................ Normal
		Power Supply by.............. Utility Power
		Utility Voltage.............. 230 V
		Output Voltage............... 230 V
		Battery Capacity............. 100 %
		Remaining Runtime............ 42 min.
		Load......................... 120 Watt(13 %)
`

	status := ParsePwrstat(output)

	if status.ModelName != "CP1500EPFCLCD" {
		t.Fatalf("unexpected model name: %q", status.ModelName)
	}
	if got := deref(status.UtilityVoltageVolts); got != 230 {
		t.Fatalf("unexpected utility voltage: %v", got)
	}
	if got := deref(status.BatteryCapacityPercent); got != 100 {
		t.Fatalf("unexpected battery capacity: %v", got)
	}
	if got := deref(status.RemainingRuntimeSeconds); got != 2520 {
		t.Fatalf("unexpected runtime seconds: %v", got)
	}
	if got := deref(status.LoadWatts); got != 120 {
		t.Fatalf("unexpected load watts: %v", got)
	}
	if got := deref(status.LoadPercent); got != 13 {
		t.Fatalf("unexpected load percent: %v", got)
	}
}

func TestParseColonSeparatedOutput(t *testing.T) {
	output := `
State: Battery Mode
Power Supply: Battery Power
Remaining Runtime: 90 sec.
`

	status := ParsePwrstat(output)

	if status.State != "Battery Mode" {
		t.Fatalf("unexpected state: %q", status.State)
	}
	if got := deref(status.RemainingRuntimeSeconds); got != 90 {
		t.Fatalf("unexpected runtime seconds: %v", got)
	}
}

func TestParsePwrstatConfig(t *testing.T) {
	output := `
Daemon Configuration:

Alarm .............................................. Off
Hibernate .......................................... Off
Cloud .............................................. Off

Action for Power Failure:

        Delay time since Power failure ............. 30 sec.
        Run script command ......................... On
        Path of script command ..................... /etc/pwrstatd-powerfail.sh
        Duration of command running ................ 0 sec.
        Enable shutdown system ..................... On

Action for Battery Low:

        Remaining runtime threshold ................ 300 sec.
        Battery capacity threshold ................. 35 %.
        Run script command ......................... On
        Path of command ............................ /etc/pwrstatd-lowbatt.sh
        Duration of command running ................ 0 sec.
        Enable shutdown system ..................... On
`

	config := ParsePwrstatConfig(output)

	if got := derefBool(config.AlarmEnabled); got {
		t.Fatalf("unexpected alarm enabled: %v", got)
	}
	if got := derefBool(config.HibernateEnabled); got {
		t.Fatalf("unexpected hibernate enabled: %v", got)
	}
	if got := derefBool(config.CloudEnabled); got {
		t.Fatalf("unexpected cloud enabled: %v", got)
	}
	if got := deref(config.PowerFailure.DelaySeconds); got != 30 {
		t.Fatalf("unexpected power failure delay: %v", got)
	}
	if got := derefBool(config.PowerFailure.ScriptEnabled); !got {
		t.Fatalf("unexpected power failure script enabled: %v", got)
	}
	if config.PowerFailure.ScriptPath != "/etc/pwrstatd-powerfail.sh" {
		t.Fatalf("unexpected power failure script path: %q", config.PowerFailure.ScriptPath)
	}
	if got := deref(config.PowerFailure.CommandDurationSecond); got != 0 {
		t.Fatalf("unexpected power failure duration: %v", got)
	}
	if got := derefBool(config.PowerFailure.ShutdownEnabled); !got {
		t.Fatalf("unexpected power failure shutdown enabled: %v", got)
	}
	if got := deref(config.LowBattery.RuntimeThresholdSeconds); got != 300 {
		t.Fatalf("unexpected low battery runtime threshold: %v", got)
	}
	if got := deref(config.LowBattery.CapacityThresholdPercent); got != 35 {
		t.Fatalf("unexpected low battery capacity threshold: %v", got)
	}
	if got := derefBool(config.LowBattery.ScriptEnabled); !got {
		t.Fatalf("unexpected low battery script enabled: %v", got)
	}
	if config.LowBattery.ScriptPath != "/etc/pwrstatd-lowbatt.sh" {
		t.Fatalf("unexpected low battery script path: %q", config.LowBattery.ScriptPath)
	}
	if got := deref(config.LowBattery.CommandDurationSecond); got != 0 {
		t.Fatalf("unexpected low battery duration: %v", got)
	}
	if got := derefBool(config.LowBattery.ShutdownEnabled); !got {
		t.Fatalf("unexpected low battery shutdown enabled: %v", got)
	}
}

func TestParsePwrstatVersion(t *testing.T) {
	output := `
version:
pwrstat version 1.4.2
`

	if got := ParsePwrstatVersion(output); got != "1.4.2" {
		t.Fatalf("unexpected pwrstat version: %q", got)
	}
}

func deref(value *float64) float64 {
	if value == nil {
		return -1
	}
	return *value
}

func derefBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}
