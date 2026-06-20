package settings

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/darklab8/fl-darkstat/configs/cfg"
	"github.com/darklab8/fl-darkstat/darkcore/envers/darkflag"
	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
)

type DarkcoreEnvVars struct {
	utils_settings.UtilsEnvs
	Password            string
	Secret              string
	ExtraCookieHost     string
	IsDiscoOauthEnabled bool
	CacheControl        string
	BuildFolder         string
	EnableUnixSockets   bool
	Enver               *enverant.Enverant
	WebPort             int
	AppStart            time.Time
	DisableAppDevMode   bool
}

var Env DarkcoreEnvVars

func GetEnvs() DarkcoreEnvVars {
	envs := enverant.NewEnverant(enverant.WithPrefix("DARKCORE_"), enverant.WithDescription("DARKCORE set of envs for a web framework based on templ to implement static site generator with backend fallback"))

	Env = DarkcoreEnvVars{
		UtilsEnvs:           utils_settings.GetEnvs(),
		Password:            envs.GetStrOr("PASSWORD", *darkflag.ArgPassword, enverant.WithDesc("protect access to web interface of darkstat with ?password=query_param")),
		CacheControl:        envs.GetStrOr("CACHE_CONTROL", ""), // refactor to boolean and set as true
		IsDiscoOauthEnabled: envs.GetBool("DISCO_OAUTH", enverant.WithDesc("an option to turn auth of darkstat for Discovery freelancer a protected dev instance of darkstat")),
		Secret:              envs.GetStrOr("SECRET", "passphrasewhichneedstobe32bytes!", enverant.WithDesc("secret to persist authentifications with query param password or oauth, required if using auths")),
		EnableUnixSockets:   envs.GetBoolOr("ENABLE_UNIX_SOCKETS", *darkflag.ArgEnableUnixSockets, enverant.WithDesc("creating unix sockets, requires /tmp/darkstat or /tmp/darkstat-{environment} folder defined")),
		WebPort:             envs.GetIntOr("WEB_PORT", *darkflag.ArgWebPort, enverant.WithDesc("specify web port")),
		BuildFolder:         envs.GetStrOr("BUILD_FOLDER", "build", enverant.WithDesc("output to which build folder name")),
		DisableAppDevMode:   envs.GetBoolOr("DISABLE_DEV_MODE", *darkflag.DisableDevMode, enverant.WithDesc("Add ability to show extra information. Map has dev mode in infocard at least")),

		AppStart: time.Now(),
		Enver:    envs,
	}
	return Env
}

//go:embed extra_disco_pob_coords.yml
var ExtraPoBCoordsStr string

type HardcodedPob struct {
	Nick         string `yaml:"nick"`
	CoordsStr    string `yaml:"coords_str"`
	SystemNick   string `yaml:"sys_nick"`
	Infocard     string `yaml:"infocard"`
	SnapshotTime string `yaml:"snapshot_time"`
	Level        int    `yaml:"level"`
	Coords       cfg.Vector
}

type HardcodedPoBConf struct {
	Pobs       []HardcodedPob `yaml:"pobs"`
	PobsByNick map[string]HardcodedPob
}

var HardcodedPoBs HardcodedPoBConf

func init() {
	err := yaml.Unmarshal([]byte(ExtraPoBCoordsStr), &HardcodedPoBs)
	if err != nil {
		log.Fatal("failed to unmarshal extra pobs")
	}
	HardcodedPoBs.PobsByNick = map[string]HardcodedPob{}
	for _, pob := range HardcodedPoBs.Pobs {

		coords := strings.Split(pob.CoordsStr, " ")
		pob.Coords.X, err = strconv.ParseFloat(coords[0], 64)
		if err != nil {
			log.Fatal("failed to unmarshal extra pobs, X")
		}
		pob.Coords.Y, err = strconv.ParseFloat(coords[1], 64)
		if err != nil {
			log.Fatal("failed to unmarshal extra pobs, Y")
		}
		pob.Coords.Z, err = strconv.ParseFloat(coords[2], 64)
		if err != nil {
			log.Fatal("failed to unmarshal extra pobs, Z")
		}

		HardcodedPoBs.PobsByNick[pob.Nick] = pob
	}

	Env = GetEnvs()
}
