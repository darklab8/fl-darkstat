package bases

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
	"github.com/darklab8/fl-darkstat/configs/discovery/pob_goods"
)

type Base struct {
	Affiliation string  `json:"affiliation"`
	Health      float64 `json:"health"`
	Tid         *int    `json:"tid"`
	Name        string  `json:"-"`
	Nickname    string
}

type Config struct {
	file *file.File
	// html.UnescapeString(pob.Base.Name) should be equal to this name
	BasesByHtmlUnescapeName map[string]*Base `json:"bases"`
	BasesByNick             map[string]*Base `json:"-"`
	Timestamp               time.Time        `json:"-"`
	Bases                   []*Base
}

func Read(ctx context.Context, file *file.File) (*Config, error) {
	byteValue, err := file.ReadBytes()
	if file_data := ctx.Value(cfg.CtxKey("pob_goods_data_override")); file_data != nil {
		byteValue = file_data.([]byte)
		err = nil
	}

	if logus.Log.CheckError(err, "failed to read file") {
		return nil, errors.New("failed to read file")
	}

	var conf *Config = &Config{}
	err = json.Unmarshal(byteValue, &conf.BasesByHtmlUnescapeName)
	conf.Timestamp = time.Now()

	if conf.BasesByNick == nil {
		conf.BasesByNick = make(map[string]*Base)
	}
	if logus.Log.CheckError(err, "failed to read pob goods") {
		return nil, errors.New("failed to read file")
	}

	for base_name, base := range conf.BasesByHtmlUnescapeName {
		base.Name = base_name

		if base.Health <= 0.001 {
			continue
		}

		hash := pob_goods.NameToNickname(base.Name)
		base.Nickname = hash
		conf.Bases = append(conf.Bases, base)
		conf.BasesByNick[base.Nickname] = base
	}
	conf.file = file
	return conf, nil
}

func (c *Config) Refresh() error {
	reread, err := Read(context.Background(), c.file)
	if logus.Log.CheckError(err, "failed to refresh") {
		return err
	}
	c.file = reread.file
	c.BasesByHtmlUnescapeName = reread.BasesByHtmlUnescapeName
	c.BasesByNick = reread.BasesByNick
	c.Timestamp = reread.Timestamp
	c.Bases = reread.Bases

	return nil
}
