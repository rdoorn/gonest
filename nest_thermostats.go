package gonest

import (
	"fmt"
	"time"
)

type Thermostat struct {
	Humidity                  float64   `json:"humidity"`
	Locale                    string    `json:"locale"`
	TemperatureScale          string    `json:"temperature_scale"`
	IsUsingEmergencyHeat      bool      `json:"is_using_emergency_heat"`
	HasFan                    bool      `json:"has_fan"`
	SoftwareVersion           string    `json:"software_version"`
	HasLeaf                   bool      `json:"has_leaf"`
	WhereID                   string    `json:"where_id"`
	DeviceID                  string    `json:"device_id"`
	Name                      string    `json:"name"`
	CanHeat                   bool      `json:"can_heat"`
	CanCool                   bool      `json:"can_cool"`
	TargetTemperatureC        float64   `json:"target_temperature_c"`
	TargetTemperatureF        float64   `json:"target_temperature_f"`
	TargetTemperatureHighC    float64   `json:"target_temperature_high_c"`
	TargetTemperatureHighF    float64   `json:"target_temperature_high_f"`
	TargetTemperatureLowC     float64   `json:"target_temperature_low_c"`
	TargetTemperatureLowF     float64   `json:"target_temperature_low_f"`
	AmbientTemperatureC       float64   `json:"ambient_temperature_c"`
	AmbientTemperatureF       float64   `json:"ambient_temperature_f"`
	AwayTemperatureHighC      float64   `json:"away_temperature_high_c"`
	AwayTemperatureHighF      float64   `json:"away_temperature_high_f"`
	AwayTemperatureLowC       float64   `json:"away_temperature_low_c"`
	AwayTemperatureLowF       float64   `json:"away_temperature_low_f"`
	EcoTemperatureHighC       float64   `json:"eco_temperature_high_c"`
	EcoTemperatureHighF       float64   `json:"eco_temperature_high_f"`
	EcoTemperatureLowC        float64   `json:"eco_temperature_low_c"`
	EcoTemperatureLowF        float64   `json:"eco_temperature_low_f"`
	IsLocked                  bool      `json:"is_locked"`
	LockedTempMinC            float64   `json:"locked_temp_min_c"`
	LockedTempMinF            float64   `json:"locked_temp_min_f"`
	LockedTempMaxC            float64   `json:"locked_temp_max_c"`
	LockedTempMaxF            float64   `json:"locked_temp_max_f"`
	SunlightCorrectionActive  bool      `json:"sunlight_correction_active"`
	SunlightCorrectionEnabled bool      `json:"sunlight_correction_enabled"`
	StructureID               string    `json:"structure_id"`
	FanTimerActive            bool      `json:"fan_timer_active"`
	FanTimerTimeout           time.Time `json:"fan_timer_timeout"`
	FanTimerDuration          int64     `json:"fan_timer_duration"`
	PreviousHVACMode          string    `json:"previous_hvac_mode"`
	HVACMode                  string    `json:"hvac_mode"`
	TimeToTarget              string    `json:"time_to_target"`
	TimeToTargetTraining      string    `json:"time_to_target_training"`
	WhereName                 string    `json:"where_name"`
	Label                     string    `json:"label"`
	NameLong                  string    `json:"name_long"`
	IsOnline                  bool      `json:"is_online"`
	LastConnection            string    `json:"last_connection"`
	HVACState                 string    `json:"hvac_state"`
}

func (h *Handler) ReadThermostats() (map[string]Thermostat, error) {
	n, err := h.Get()
	if err != nil {
		return nil, err
	}

	thermostats := make(map[string]Thermostat)
	for tid, t := range n.Devices.Thermostats {
		thermostats[tid] = t
	}
	return thermostats, nil
}

func (h *Handler) SetTemperature(temperature float64) error {
	themostats, err := h.ReadThermostats()
	if err != nil {
		return err
	}
	for tid := range themostats {
		if err := h.Set(fmt.Sprintf("devices/thermostats/%s", tid), fmt.Sprintf(`{"target_temperature_c": "%f"}`, temperature)); err != nil {
			return err
		}
	}
	return nil
}
