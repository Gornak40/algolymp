package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gornak40/algolymp/ejudge"
	"github.com/Gornak40/algolymp/polygon"
	"github.com/sirupsen/logrus"
)

type System struct {
	Editor string `json:"editor"`
}

type Config struct {
	Ejudge  ejudge.Config  `json:"ejudge"`
	Polygon polygon.Config `json:"polygon"`
	System  System         `json:"system"`
}

func NewConfig() *Config {
	confDir, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s/.config/algolymp/config.json", confDir)
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
