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

func deref(value *float64) float64 {
	if value == nil {
		return -1
	}
	return *value
}
