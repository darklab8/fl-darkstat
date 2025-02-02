package pob_goods

import (
	"encoding/json"
	"html"
	"regexp"
	"strings"

	"github.com/darklab8/fl-darkstat/configs/configs_mapped/freelancer_mapped/data_mapped/initialworld/flhash"
	"github.com/darklab8/fl-darkstat/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-darkstat/configs/configs_settings/logus"
)

type ShopItem struct {
	Id        int `json:"id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
	SellPrice int `json:"sell_price"`
	MinStock  int `json:"min_stock"`
	MaxStock  int `json:"max_stock"`
}

func (good ShopItem) BaseSells() bool {
	return good.Quantity > good.MinStock
}
func (good ShopItem) BaseBuys() bool {
	return good.Quantity < good.MaxStock
}

type Base struct {
	Name      string
	Nickname  string
	ShopItems []ShopItem `json:"shop_items"`

	ForumThreadUrl *string `json:"thread"`
	CargoSpaceLeft *int    `json:"cargospace"`

	SystemHash         *flhash.HashCode `json:"system"`      //: 2745655887,
	Pos                *string          `json:"pos"`         //: "299016, 33, -178",
	AffiliationHash    *flhash.HashCode `json:"affiliation"` //: 2620,
	Level              *int             `json:"level"`       //: 1,
	Money              *int             `json:"money"`       //: 0,
	Health             *float64         `json:"health"`      //: 50,
	DefenseMode        *int             `json:"defensemode"` //: 1,
	InfocardParagraphs []string         `json:"infocard_paragraphs"`

	HostileFactionHashList []*flhash.HashCode `json:"hostile_list"`
	HostileTagList         []string           `json:"hostile_tag_list"`
	HostileNameList        []string           `json:"hostile_name_list"`
	AllyFactionHashList    []*flhash.HashCode `json:"ally_list"`
	AllyTagList            []string           `json:"ally_tag_list"`
	AllyNameList           []string           `json:"ally_name_list"`
	SrpFactionHashList     []*flhash.HashCode `json:"srp_list"`
	SrpTagList             []string           `json:"srp_tag_list"`
	SrpNameList            []string           `json:"srp_name_list"`
}

type Config struct {
	file        *file.File
	BasesByName map[string]*Base `json:"bases"`
	Timestamp   string           `json:"timestamp"`
	Bases       []*Base
}

func (c *Config) Refresh() {
	reread := Read(c.file)
	c.file = reread.file
	c.BasesByName = reread.BasesByName
	c.Timestamp = reread.Timestamp
	c.Bases = reread.Bases
}

func NameToNickname(name string) string {
	name = strings.ToLower(name)
	name = html.UnescapeString(name)
	name = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(name, "")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func Read(file *file.File) *Config {
	byteValue, err := file.ReadBytes()
	logus.Log.CheckFatal(err, "failed to read file")

	var conf *Config
	json.Unmarshal(byteValue, &conf)

	for base_name, base := range conf.BasesByName {
		base.Name = base_name

		hash := NameToNickname(base.Name)
		base.Nickname = hash
		conf.Bases = append(conf.Bases, base)
	}

	conf.file = file
	return conf
}

func (frelconfig *Config) Write() *file.File {
	return &file.File{}
}
