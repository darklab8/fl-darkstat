package stararch_mapped

import (
	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/semantic"
)

type Star struct {
	semantic.Model
	Nickname   *semantic.String
	StarGlow   *semantic.String
	StarCenter *semantic.String
	Radius     *semantic.Float
}

type Glow struct {
	semantic.Model
	Nickname   *semantic.String
	InnerColor *semantic.Vect
	OuterColor *semantic.Vect
	Scale      *semantic.Float
}

type Config struct {
	*iniload.IniLoader
	Stars       []*Star
	StarsByNick map[string]*Star
	GlowsByNick map[string]*Glow
}

func Read(input_files []*iniload.IniLoader) *Config {
	frelconfig := &Config{
		StarsByNick: make(map[string]*Star),
		GlowsByNick: make(map[string]*Glow),
	}
	for _, input_file := range input_files {
		for _, section := range input_file.SectionMap["[star]"] {
			star := &Star{
				Nickname:   semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				StarGlow:   semantic.NewString(section, cfg.Key("star_glow"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				StarCenter: semantic.NewString(section, cfg.Key("star_center"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				Radius:     semantic.NewFloat(section, cfg.Key("radius"), semantic.Precision(2)),
			}
			star.Map(section)
			frelconfig.Stars = append(frelconfig.Stars, star)
			frelconfig.StarsByNick[star.Nickname.Get()] = star

		}
		for _, section := range input_file.SectionMap["[star_glow]"] {
			glow := &Glow{
				Nickname:   semantic.NewString(section, cfg.Key("nickname"), semantic.WithLowercaseS(), semantic.WithoutSpacesS()),
				InnerColor: semantic.NewVector(section, cfg.Key("inner_color"), semantic.Precision(2)),
				OuterColor: semantic.NewVector(section, cfg.Key("outer_color"), semantic.Precision(2)),
				Scale:      semantic.NewFloat(section, cfg.Key("scale"), semantic.Precision(2)),
			}
			glow.Map(section)
			frelconfig.GlowsByNick[glow.Nickname.Get()] = glow

		}
	}
	return frelconfig

}

func (frelconfig *Config) Write() *file.File {
	inifile := frelconfig.Render()
	inifile.Write(inifile.File)
	return inifile.File
}
