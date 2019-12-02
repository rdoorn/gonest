package gonest

type Device struct {
	Thermostats map[string]Thermostat `json:"thermostats"`
}
