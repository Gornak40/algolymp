package config

import (
	"encoding/json"
	"os"

	"github.com/Gornak40/algolymp/ejudge"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Ejudge ejudge.Config `json:"ejudge"`
}

func NewConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read config")
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		logrus.WithError(err).Fatal("failed to unmarshal config")
	}
	return &cfg
}
