package filesave

// Provided by marko_oktabyr (discord) at https://gist.github.com/dwmunster/fec794a4a938967d25ebdd4c29091350
// Originally was hooked to https://github.com/dwmunster/ini lib
// with import "github.com/dwmunster/ini"
// But looks like it was merged to main repo https://github.com/go-ini/ini/issues/284

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"gopkg.in/ini.v1"
)

func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var cs []string
	chunk := make([]rune, chunkSize)
	l := 0
	for _, r := range s {
		chunk[l] = r
		l++
		if l == chunkSize {
			cs = append(cs, string(chunk))
			l = 0
		}
	}
	if l > 0 {
		cs = append(cs, string(chunk[:l]))
	}
	return cs
}

type CodepointString string

func (cs CodepointString) String() string {
	var points []uint16
	for _, chunk := range chunks(string(cs), 4) {
		u, err := strconv.ParseUint(chunk, 16, 16)
		if err != nil {
			log.Printf("unable to decode codepoint: %s\n", chunk)
			u = 0
		}
		points = append(points, uint16(u))
	}
	return string(utf16.Decode(points))
}

func CPSFromString(s string) CodepointString {
	var b strings.Builder

	// Preallocate memory
	b.Grow(len(s) * 4)

	for _, c := range s {
		_, _ = fmt.Fprintf(&b, "%04X", c)
	}

	return CodepointString(b.String())

}

type TStamp []uint

const (
	winTicksEpochDifference = 116444736000000000
	nsPerWinTick            = 100
)

func (t TStamp) Time() time.Time {
	// 100-nanosecond intervals since January 1, 1601
	nsec := int64(t[0])<<32 + int64(t[1])
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= winTicksEpochDifference
	// convert into nanoseconds
	nsec *= nsPerWinTick
	return time.Unix(0, nsec)
}

func (t TStamp) String() string {
	return t.Time().String()
}

func TStampFromTime(t time.Time) TStamp {
	nsec := t.UnixNano()
	winTicks := (nsec / nsPerWinTick) + winTicksEpochDifference
	u := uint(winTicks)
	return TStamp{(u >> 32) & 0xFFFFFFFF, u & 0xFFFFFFFF}
}

type HashCode uint32

type Equipment struct {
	ID     HashCode
	Mount  string
	Health float64
}

func (e Equipment) INILine() string {
	return fmt.Sprintf("%d, %s, 1", e.ID, e.Mount)
}

type Cargo struct {
	ID       HashCode
	Quantity int
	Health   float64
}

type ShipKills struct {
	ID       HashCode
	Quantity int
}

type Model struct {
	Body      string `json:"body"`
	Head      string `json:"head"`
	LeftHand  string `json:"left_hand"`
	RightHand string `json:"right_hand"`
}

type Player struct {
	Name                    string             `json:"name"`
	Description             string             `json:"description"`
	Timestamp               time.Time          `json:"timestamp"`
	Rank                    int                `json:"rank"`
	Reputations             map[string]float64 `json:"reputations"`
	ReputationGroup         string             `json:"reputation_group"`
	Money                   int                `json:"money"`
	ComModel                Model              `json:"com_model"`
	Model                   Model              `json:"model"`
	System                  string             `json:"system"`
	Base                    string             `json:"base,omitempty"`
	Position                []float64          `json:"position,omitempty"`
	Rotation                []float64          `json:"rotation,omitempty"`
	Ship                    HashCode           `json:"ship"`
	Equipment               []Equipment        `json:"equipment"`
	Cargo                   []Cargo            `json:"cargo"`
	LastBase                string             `json:"last_base"`
	BaseHullStatus          float64            `json:"base_hull_status"`
	BaseCollisionGroups     []string           `json:"base_collision_groups"`
	BaseEquipment           []Equipment        `json:"base_equipment"`
	BaseCargo               []Cargo            `json:"base_cargo"`
	Visited                 []string           `json:"visited"`
	TimePlayed              time.Duration      `json:"time_played"`
	SystemsVisited          []string           `json:"systems_visited"`
	BasesVisited            []string           `json:"bases_visited"`
	HolesVisited            []string           `json:"holes_visited"`
	ShipKills               []ShipKills        `json:"ship_kills"`
	VisitedNPC              []string           `json:"vnpc,omitempty"`
	RandomMissionsCompleted []string           `json:"random_missions_completed,omitempty"`
	RandomMissionsAborted   []string           `json:"random_missions_aborted,omitempty"`
	RandomMissionsFailed    []string           `json:"random_missions_failed,omitempty"`
	Rumors                  []string           `json:"rumors,omitempty"`
}

func UpdateINI(file *ini.File, p Player) error {
	pSec := file.Section("Player")

	pSec.Key("name").SetValue(string(CPSFromString(p.Name)))
	pSec.Key("description").SetValue(string(CPSFromString(p.Description)))

	t := TStampFromTime(p.Timestamp)
	pSec.Key("tstamp").SetValue(fmt.Sprintf("%d,%d", t[0], t[1]))
	pSec.Key("rank").SetValue(fmt.Sprint(p.Rank))

	pSec.DeleteKey("house")
	hk := pSec.Key("house")
	anySet := false
	for k, v := range p.Reputations {
		if !anySet {
			anySet = true
			hk.SetValue(fmt.Sprintf("%0f, %s", v, k))
			continue
		}
		err := hk.AddShadow(fmt.Sprintf("%01g, %s", v, k))
		if err != nil {
			return fmt.Errorf("unable to add shadow key to 'house': %w", err)
		}
	}
	pSec.DeleteKey("rep_group")
	if p.ReputationGroup != "" {
		pSec.Key("rep_group").SetValue(p.ReputationGroup)
	}

	pSec.Key("money").SetValue(fmt.Sprint(p.Money))

	pSec.Key("com_body").SetValue(p.ComModel.Body)
	pSec.Key("com_head").SetValue(p.ComModel.Head)
	pSec.Key("com_lefthand").SetValue(p.ComModel.LeftHand)
	pSec.Key("com_righthand").SetValue(p.ComModel.RightHand)

	pSec.Key("body").SetValue(p.Model.Body)
	pSec.Key("head").SetValue(p.Model.Head)
	pSec.Key("lefthand").SetValue(p.Model.LeftHand)
	pSec.Key("righthand").SetValue(p.Model.RightHand)

	pSec.Key("system").SetValue(p.System)
	// Prefer Base to Position
	if p.Base != "" {
		pSec.DeleteKey("pos")
		pSec.DeleteKey("rotate")
		pSec.Key("base").SetValue(p.Base)
	} else {
		pSec.DeleteKey("base")
		pSec.Key("pos").SetValue(fmt.Sprintf("%01g,%01g,%01g", p.Position[0], p.Position[1], p.Position[2]))
		pSec.Key("rotate").SetValue(fmt.Sprintf("%01g,%01g,%01g", p.Rotation[0], p.Rotation[1], p.Rotation[2]))

	}

	pSec.Key("ship_archetype").SetValue(fmt.Sprintf("%d", p.Ship))

	pSec.DeleteKey("equip")
	if len(p.Equipment) > 0 {
		ek := pSec.Key("equip")
		for i, eq := range p.Equipment {
			if i == 0 {
				ek.SetValue(fmt.Sprintf("%d, %s, 1", eq.ID, eq.Mount))
				continue
			}
			err := ek.AddShadow(fmt.Sprintf("%d, %s, 1", eq.ID, eq.Mount))
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'equip': %w", err)
			}
		}
	}

	pSec.DeleteKey("cargo")
	if len(p.Cargo) > 0 {
		ck := pSec.Key("cargo")
		for i, c := range p.Cargo {
			if i == 0 {
				ck.SetValue(fmt.Sprintf("%d, %d, , %01f, 0", c.ID, c.Quantity, c.Health))
				continue
			}
			err := ck.AddShadow(fmt.Sprintf("%d, %d, , %01f, 0", c.ID, c.Quantity, c.Health))
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'cargo': %w", err)
			}
		}
	}

	pSec.Key("last_base").SetValue(p.LastBase)
	pSec.Key("base_hull_status").SetValue(fmt.Sprintf("%01f", p.BaseHullStatus))
	pSec.DeleteKey("base_collision_group")
	if len(p.BaseCollisionGroups) > 0 {
		ck := pSec.Key("base_collision_group")
		for i, c := range p.BaseCollisionGroups {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'base_collision_group': %w", err)
			}
		}
	}

	pSec.DeleteKey("base_equip")
	if len(p.BaseEquipment) > 0 {
		bek := pSec.Key("base_equip")
		for i, eq := range p.BaseEquipment {
			if i == 0 {
				bek.SetValue(fmt.Sprintf("%d, %s, %01f", eq.ID, eq.Mount, eq.Health))
				continue
			}
			err := bek.AddShadow(fmt.Sprintf("%d, %s, %01f", eq.ID, eq.Mount, eq.Health))
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'base_equip': %w", err)
			}
		}
	}

	pSec.DeleteKey("base_cargo")
	if len(p.BaseCargo) > 0 {
		bck := pSec.Key("base_cargo")
		for i, c := range p.BaseCargo {
			if i == 0 {
				bck.SetValue(fmt.Sprintf("%d, %d, , %01f, 0", c.ID, c.Quantity, c.Health))
				continue
			}
			err := bck.AddShadow(fmt.Sprintf("%d, %d, , %01f, 0", c.ID, c.Quantity, c.Health))
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'base_cargo': %w", err)
			}
		}
	}

	pSec.DeleteKey("visit")
	if len(p.Visited) > 0 {
		ck := pSec.Key("visit")
		sort.Slice(p.Visited, func(i, j int) bool {
			pi := strings.Split(p.Visited[i], ",")
			pj := strings.Split(p.Visited[j], ",")
			ii, _ := strconv.Atoi(pi[0])
			ij, _ := strconv.Atoi(pj[0])
			return ii < ij
		})
		for i, c := range p.Visited {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'visit': %w", err)
			}
		}
	}

	mSec := file.Section("mPlayer")

	mSec.Key("total_time_played").SetValue(fmt.Sprintf("%01f", p.TimePlayed.Seconds()))

	mSec.DeleteKey("sys_visited")
	if len(p.SystemsVisited) > 0 {
		ck := mSec.Key("sys_visited")
		for i, c := range p.SystemsVisited {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'sys_visited': %w", err)
			}
		}
	}

	mSec.DeleteKey("base_visited")
	if len(p.BasesVisited) > 0 {
		ck := mSec.Key("base_visited")
		for i, c := range p.BasesVisited {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'base_visited': %w", err)
			}
		}
	}

	mSec.DeleteKey("holes_visited")
	if len(p.HolesVisited) > 0 {
		ck := mSec.Key("holes_visited")
		for i, c := range p.HolesVisited {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'holes_visited': %w", err)
			}
		}
	}

	mSec.DeleteKey("ship_type_killed")
	if len(p.ShipKills) > 0 {
		sk := mSec.Key("ship_type_killed")
		for i, k := range p.ShipKills {
			if i == 0 {
				sk.SetValue(fmt.Sprintf("%d, %d", k.ID, k.Quantity))
				continue
			}
			err := sk.AddShadow(fmt.Sprintf("%d, %d", k.ID, k.Quantity))
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'ship_type_killed': %w", err)
			}
		}
	}

	mSec.DeleteKey("vnpc")
	if len(p.VisitedNPC) > 0 {
		ck := mSec.Key("vnpc")
		for i, c := range p.VisitedNPC {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'vnpc': %w", err)
			}
		}
	}

	mSec.DeleteKey("rm_completed")
	if len(p.RandomMissionsCompleted) > 0 {
		ck := mSec.Key("rm_completed")
		for i, c := range p.RandomMissionsCompleted {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'rm_completed': %w", err)
			}
		}
	}

	mSec.DeleteKey("rm_aborted")
	if len(p.RandomMissionsAborted) > 0 {
		ck := mSec.Key("rm_aborted")
		for i, c := range p.RandomMissionsAborted {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'rm_aborted': %w", err)
			}
		}
	}

	mSec.DeleteKey("rm_failed")
	if len(p.RandomMissionsFailed) > 0 {
		ck := mSec.Key("rm_failed")
		for i, c := range p.RandomMissionsFailed {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'rm_failed': %w", err)
			}
		}
	}

	mSec.DeleteKey("rumor")
	if len(p.Rumors) > 0 {
		ck := mSec.Key("rumor")
		for i, c := range p.Rumors {
			if i == 0 {
				ck.SetValue(c)
				continue
			}
			err := ck.AddShadow(c)
			if err != nil {
				return fmt.Errorf("unable to add shadow key to 'rumor': %w", err)
			}
		}
	}

	return nil
}

type rawPlayer struct {
	Name               CodepointString `ini:"name"`
	Description        CodepointString `ini:"description"`
	Timestamp          TStamp          `ini:"tstamp"`
	Rank               int             `ini:"rank"`
	Reputations        []string        `ini:"house,allowshadow"`
	ReputationGroup    string          `ini:"rep_group,omitempty"`
	Money              int             `ini:"money"`
	NumKills           int             `ini:"num_kills"`
	NumMissionSuccess  int             `ini:"num_misn_successes"`
	NumMissionFailure  int             `ini:"num_misn_failures"`
	Voice              string          `ini:"voice"`
	ComBody            string          `ini:"com_body"`
	ComHead            string          `ini:"com_head"`
	ComLeftHand        string          `ini:"com_lefthand"`
	ComRightHand       string          `ini:"com_righthand"`
	Body               string          `ini:"body"`
	Head               string          `ini:"head"`
	LeftHand           string          `ini:"lefthand"`
	RightHand          string          `ini:"righthand"`
	System             string          `ini:"system"`
	Base               string          `ini:"base,omitempty"`
	Position           []float64       `ini:"pos,omitempty"`
	Rotation           []float64       `ini:"rotate,omitempty"`
	Ship               HashCode        `ini:"ship_archetype"`
	Equipment          []string        `ini:"equip,allowshadow"`
	Cargo              []string        `ini:"cargo,allowshadow"`
	LastBase           string          `ini:"last_base"`
	BaseHullStatus     float64         `ini:"base_hull_status"`
	BaseCollisionGroup []string        `ini:"base_collision_group,allowshadow,omitempty" delim:"|"`
	BaseEquipment      []string        `ini:"base_equip,allowshadow"`
	BaseCargo          []string        `ini:"base_cargo,allowshadow"`
	Visited            []string        `ini:"visit,allowshadow" delim:"|"`
}

type multiPlayer struct {
	TimePlayed              float64  `ini:"total_time_played"`
	SystemsVisited          []string `ini:"sys_visited,allowshadow,omitempty"`
	BasesVisited            []string `ini:"base_visited,allowshadow,omitempty"`
	HolesVisited            []string `ini:"holes_visited,allowshadow,omitempty"`
	ShipKills               []int    `ini:"ship_type_killed,allowshadow,omitempty"`
	VNPC                    []string `ini:"vnpc,allowshadow,omitempty" delim:"|"`
	RandomMissionsCompleted []string `ini:"rm_completed,allowshadow,omitempty" delim:"|"`
	RandomMissionsAborted   []string `ini:"rm_aborted,allowshadow,omitempty" delim:"|"`
	RandomMissionsFailed    []string `ini:"rm_failed,allowshadow,omitempty" delim:"|"`
	Rumor                   []string `ini:"rumor,allowshadow,omitempty" delim:"|"`
}

func repParser(repEntries []string) (map[string]float64, error) {
	reps := make(map[string]float64)

	if len(repEntries)%2 != 0 {
		return reps, fmt.Errorf("repEntries must have an even number of repEntries")
	}

	for i := 0; i < len(repEntries); i += 2 {
		val, err := strconv.ParseFloat(repEntries[i], 64)
		if err != nil {
			return reps, fmt.Errorf("error decoding '%s' into float: %w", repEntries[i], err)
		}
		reps[repEntries[i+1]] = val
	}

	return reps, nil
}

func equipParser(equipEntries []string) ([]Equipment, error) {
	if len(equipEntries)%3 != 0 {
		return []Equipment{}, fmt.Errorf("len(equipEntries) must be divisible by three")
	}
	eq := make([]Equipment, 0, len(equipEntries)/3)

	for i := 0; i < len(equipEntries); i += 3 {
		code, err := strconv.ParseUint(equipEntries[i], 10, 32)
		if err != nil {
			return eq, fmt.Errorf("error decoding '%s' into uint: %w", equipEntries[i], err)
		}
		health, err := strconv.ParseFloat(equipEntries[i+2], 64)
		if err != nil {
			return eq, fmt.Errorf("error decoding '%s' into float64: %w", equipEntries[i+2], err)
		}
		eq = append(eq, Equipment{
			ID:     HashCode(code),
			Mount:  equipEntries[i+1],
			Health: health,
		})
	}
	return eq, nil
}

func killParser(killEntries []int) ([]ShipKills, error) {
	if len(killEntries)%2 != 0 {
		return []ShipKills{}, fmt.Errorf("len(killEntries) must be divisible by two")
	}
	kills := make([]ShipKills, 0, len(killEntries)/2)

	for i := 0; i < len(killEntries); i += 2 {
		kills = append(kills, ShipKills{
			ID:       HashCode(killEntries[i]),
			Quantity: killEntries[i+1],
		})
	}
	return kills, nil
}

func cargoParser(cargoEntries []string) ([]Cargo, error) {
	if len(cargoEntries)%5 != 0 {
		return []Cargo{}, fmt.Errorf("len(cargoEntries) must be divisible by five")
	}
	cargo := make([]Cargo, 0, len(cargoEntries)/5)

	for i := 0; i < len(cargoEntries); i += 5 {
		code, err := strconv.ParseUint(cargoEntries[i], 10, 32)
		if err != nil {
			return cargo, fmt.Errorf("error decoding '%s' into uint: %w", cargoEntries[i], err)
		}
		q, err := strconv.ParseUint(cargoEntries[i+1], 10, 32)
		if err != nil {
			return cargo, fmt.Errorf("error decoding '%s' into uint: %w", cargoEntries[i], err)
		}
		health := 1.
		if cargoEntries[i+3] != "" {
			health, err = strconv.ParseFloat(cargoEntries[i+3], 64)
			if err != nil {
				return cargo, fmt.Errorf("error decoding '%s' into float64: %w", cargoEntries[i+3], err)
			}
		}

		cargo = append(cargo, Cargo{
			ID:       HashCode(code),
			Quantity: int(q),
			Health:   health,
		})
	}
	return cargo, nil
}

type saveFile struct {
	RawPlayer   rawPlayer   `ini:"Player"`
	Multiplayer multiPlayer `ini:"mPlayer"`
}

func (s saveFile) Player() (Player, error) {
	reps, err := repParser(s.RawPlayer.Reputations)
	if err != nil {
		return Player{}, err
	}

	equip, err := equipParser(s.RawPlayer.Equipment)
	if err != nil {
		return Player{}, err
	}

	bEquip, err := equipParser(s.RawPlayer.BaseEquipment)
	if err != nil {
		return Player{}, err
	}

	kills, err := killParser(s.Multiplayer.ShipKills)
	if err != nil {
		return Player{}, err
	}

	cargo, err := cargoParser(s.RawPlayer.Cargo)
	if err != nil {
		return Player{}, err
	}

	bCargo, err := cargoParser(s.RawPlayer.BaseCargo)
	if err != nil {
		return Player{}, err
	}

	return Player{
		Name:            s.RawPlayer.Name.String(),
		Description:     s.RawPlayer.Description.String(),
		Timestamp:       s.RawPlayer.Timestamp.Time(),
		Rank:            s.RawPlayer.Rank,
		Reputations:     reps,
		ReputationGroup: s.RawPlayer.ReputationGroup,
		Money:           s.RawPlayer.Money,
		ComModel: Model{
			Body:      s.RawPlayer.ComBody,
			Head:      s.RawPlayer.ComHead,
			LeftHand:  s.RawPlayer.ComLeftHand,
			RightHand: s.RawPlayer.ComRightHand,
		},
		Model: Model{
			Body:      s.RawPlayer.Body,
			Head:      s.RawPlayer.Head,
			LeftHand:  s.RawPlayer.LeftHand,
			RightHand: s.RawPlayer.RightHand,
		},
		System:                  s.RawPlayer.System,
		Base:                    s.RawPlayer.Base,
		Position:                s.RawPlayer.Position,
		Rotation:                s.RawPlayer.Rotation,
		Ship:                    s.RawPlayer.Ship,
		Equipment:               equip,
		Cargo:                   cargo,
		LastBase:                s.RawPlayer.LastBase,
		BaseHullStatus:          s.RawPlayer.BaseHullStatus,
		BaseCollisionGroups:     s.RawPlayer.BaseCollisionGroup,
		BaseEquipment:           bEquip,
		BaseCargo:               bCargo,
		Visited:                 s.RawPlayer.Visited,
		TimePlayed:              time.Duration(s.Multiplayer.TimePlayed) * time.Second,
		SystemsVisited:          s.Multiplayer.SystemsVisited,
		BasesVisited:            s.Multiplayer.BasesVisited,
		HolesVisited:            s.Multiplayer.HolesVisited,
		ShipKills:               kills,
		VisitedNPC:              s.Multiplayer.VNPC,
		RandomMissionsCompleted: s.Multiplayer.RandomMissionsCompleted,
		RandomMissionsAborted:   s.Multiplayer.RandomMissionsAborted,
		RandomMissionsFailed:    s.Multiplayer.RandomMissionsFailed,
		Rumors:                  s.Multiplayer.Rumor,
	}, nil
}

func LoadPlayer(file *ini.File) (Player, error) {
	var s saveFile
	err := file.MapTo(&s)
	if err != nil {
		return Player{}, fmt.Errorf("unable to map ini to save: %w", err)
	}

	p, err := s.Player()
	if err != nil {
		return Player{}, fmt.Errorf("unable to map save to Player: %w", err)
	}
	return p, nil
}
