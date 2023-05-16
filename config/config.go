package config

import "time"

type Status struct {
	StartTime time.Time     `json:"start_time"`
	UpTime    time.Duration `json:"up_time"`
	Version   string        `json:"version"`
}

type Policy struct {
	Config struct {
		Username      string                 `mapstructure:"username,omitempty" yaml:"username,omitempty"`
		Password      string                 `mapstructure:"password,omitempty" yaml:"password,omitempty"`
		Port          int                    `mapstructure:"port,omitempty" yaml:"port,omitempty"`
		Timeout       string                 `mapstructure:"timeout,omitempty" yaml:"timeout,omitempty"`
		SkipVerify    bool                   `mapstructure:"skip-verify,omitempty" yaml:"skip-verify,omitempty"`
		Encoding      string                 `mapstructure:"encoding,omitempty" yaml:"encoding,omitempty"`
		Targets       map[string]interface{} `mapstructure:"targets,omitempty" json:"targets,omitempty" yaml:"targets,omitempty"`
		Subscriptions map[string]interface{} `mapstructure:"subscriptions,omitempty" json:"subscriptions,omitempty" yaml:"subscriptions,omitempty"`
		Outputs       map[string]interface{} `mapstructure:"outputs,omitempty" json:"outputs,omitempty" yaml:"outputs,omitempty"`
		Inputs        map[string]interface{} `mapstructure:"inputs,omitempty" json:"inputs,omitempty" yaml:"inputs,omitempty"`
		Processors    map[string]interface{} `mapstructure:"processors,omitempty" json:"processors,omitempty" yaml:"processors,omitempty"`
	} `mapstructure:"config,omitempty" yaml:"config,omitempty"`
}

type Config struct {
	Debug         bool   `mapstructure:"orb_gnmic_debug"`
	SelfTelemetry bool   `mapstructure:"orb_gnmic_self_telemetry"`
	ServerHost    string `mapstructure:"orb_gnmic_server_host"`
	ServerPort    uint64 `mapstructure:"orb_gnmic_server_port"`
}
