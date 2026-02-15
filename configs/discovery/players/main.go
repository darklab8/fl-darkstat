package players

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

type Player struct {
	Time   string `json:"time"`
	Name   string `json:"name"`
	System string `json:"system"`
	Region string `json:"region"`
}

type Config struct {
	file         *file.File
	TimestampReq time.Time `json:"-"`

	Error     interface{} `json:"error"`
	Players   []Player    `json:"players"`
	Timestamp string      `json:"timestamp"`
}

func Read(ctx context.Context, file *file.File) (*Config, error) {
	byteValue, err := file.ReadBytes()

	if logus.Log.CheckError(err, "failed to read file") {
		return nil, errors.New("failed to read file")
	}

	var conf *Config
	err = json.Unmarshal(byteValue, &conf)
	conf.TimestampReq = time.Now()

	if logus.Log.CheckError(err, "failed to read players") {
		return nil, errors.New("failed to read file")
	}

	conf.file = file
	return conf, nil
}
