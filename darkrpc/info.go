package darkrpc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/front/tab"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
	"gopkg.in/yaml.v3"
)

type GetInfoArgs struct {
	Query string
}

type GetInfoReply struct {
	Content []string
	Found   []InfoFound
}

func (t *ServerRpc) IsInfoFound(args GetInfoArgs, name string, nickname string) (bool, bool) {
	lowered_query := strings.ToLower(args.Query)
	if strings.Contains(strings.ToLower(name), lowered_query) {
		return true, false
	}

	if strings.Contains(strings.ToLower(nickname), lowered_query) {
		if _, ok := t.app_data.Configs.Hashes[nickname]; ok && nickname == lowered_query {
			return true, true // perfect nickname match
		}

		return true, false
	}

	first_line_in_infocard := tab.GetFirstLine(t.app_data.Configs.Infocards, configs_export.InfocardKey(nickname))
	if strings.Contains(strings.ToLower(first_line_in_infocard), lowered_query) {
		return true, false
	}

	return false, false
}

func EntityToYamlStrings(entity any) []string {
	var Content []string
	data, err := json.Marshal(entity)
	if err != nil {
		Content = append(Content, err.Error())
	}
	var hashmap map[string]interface{}
	err = json.Unmarshal(data, &hashmap)
	delete(hashmap, "BGCS_base_run_by")
	delete(hashmap, "file")
	delete(hashmap, "li01_01_base")
	delete(hashmap, "InfocardKey")
	delete(hashmap, "rephacks")
	delete(hashmap, "equipment_slots")
	delete(hashmap, "biggest_hardpoint")
	delete(hashmap, "ship_packages")
	delete(hashmap, "reputations")
	delete(hashmap, "bribe")
	delete(hashmap, "archetypes")
	delete(hashmap, "damage_bonuses")
	for key, _ := range hashmap {
		if strings.Contains(key, "_hash") {
			delete(hashmap, key)
		}
	}
	if err != nil {
		Content = append(Content, err.Error())
	}
	yaml_bytes, err := yaml.Marshal(hashmap)
	if err != nil {
		Content = append(Content, err.Error())
	}
	yaml_strs := strings.Split(string(yaml_bytes), "\n")
	Content = append(Content, "```yml")
	Content = append(Content, yaml_strs...)
	Content = append(Content, "```")
	return Content
}

type InfoFound struct {
	Nickname   string
	Name       string
	Entity     string
	FirstLine  string
	Obtainable bool
}

func (t *ServerRpc) NewInfoFound(Nickname string, Name string, Entity string, Obtainable bool) InfoFound {
	return InfoFound{
		Name:       Name,
		Nickname:   string(Nickname),
		FirstLine:  tab.GetFirstLine(t.app_data.Configs.Infocards, configs_export.InfocardKey(string(Nickname))),
		Obtainable: Obtainable,
		Entity:     Entity,
	}
}

func (t *ServerRpc) GetInfo(args GetInfoArgs, reply *GetInfoReply) error {

	if strings.ReplaceAll(args.Query, " ", "") == "" {
		reply.Content = []string{}
		reply.Content = append(reply.Content, "Input some name (or nickname) parts of a Freelancer item, base or pob")
		reply.Content = append(reply.Content, "for example: . info iw04_01_base")

		return nil
	}

	set_infocard := func(nickname string) {
		infocard := t.app_data.Configs.Infocards[configs_export.InfocardKey(nickname)]
		for _, line := range infocard {
			reply.Content = append(reply.Content, line.ToStr())
		}
	}

	for _, item := range t.app_data.Configs.Bases {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Base**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Base", false))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Ammos {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Ammo**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Ammo", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.MiningOperations {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Mining Operation**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Mining Operation", false))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Factions {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Faction**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Faction", false))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Commodities {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Commodity**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Commodity", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Guns {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Gun**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Gun", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Missiles {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Missile**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Missile", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Mines {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Mine**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Mine", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Shields {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Shield**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Shield", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Ships {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Ship**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Ship", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Thrusters {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Thruster**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Thruster", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Tractors {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Tractor**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Tractor", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Engines {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Engine**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Engine", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Cloaks {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **CloakingDevice**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Cloak", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.CMs {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Counter Measure**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "CM", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.Scanners {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Scanner**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "Scanner", t.app_data.Configs.Buyable(item.Bases)))
			reply.Content = append(reply.Content, EntityToYamlStrings(item)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}
	for _, item := range t.app_data.Configs.PoBs {
		if ok, is_perfect_nickname_match := t.IsInfoFound(args, item.Name, string(item.Nickname)); ok {
			reply.Content = []string{"entity: **Player Owned Base**"}
			reply.Found = append(reply.Found, t.NewInfoFound(string(item.Nickname), item.Name, "PoB", false))
			reply.Content = append(reply.Content, EntityToYamlStrings(item.PoBWithoutGoods)...)
			set_infocard(string(item.Nickname))
			if is_perfect_nickname_match {
				return nil
			}
		}
	}

	if len(reply.Found) == 0 {
		reply.Content = []string{}
		reply.Content = append(reply.Content, "no matching names or nicknames of entities were found.")
	}
	if len(reply.Found) > 1 {
		reply.Content = []string{}
		var sb strings.Builder

		sort.Slice(reply.Found, func(i, j int) bool {
			if reply.Found[i].Obtainable != reply.Found[j].Obtainable {
				return reply.Found[i].Obtainable
			}
			return reply.Found[i].Name < reply.Found[j].Name
		})

		sb.WriteString("Multiple entities were found possessing same name and nickname. ")
		sb.WriteString(fmt.Sprintf("Repeat request with more precise **name** or **nickname**. Printing no more than 10 matched entities (total matched: %d):", len(reply.Found)))
		reply.Content = append(reply.Content, sb.String())
		for i := 0; i < 10 && i < len(reply.Found); i++ {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("- **Name**: %s, **Nickname**: %s, **Type**: %s", reply.Found[i].Name, reply.Found[i].Nickname, reply.Found[i].Entity))
			if reply.Found[i].Obtainable {
				sb.WriteString(", **Obtainable**: true")
			}
			sb.WriteString(fmt.Sprintf(", **InfoName**: %s", reply.Found[i].FirstLine))
			reply.Content = append(reply.Content, sb.String())
		}
	}

	return nil
}

func (r *ClientRpc) GetInfo(args GetInfoArgs, reply *GetInfoReply) error {
	// Synchronous call
	// return r.client.Call("ServerRpc.GetBases", args, &reply)

	// // Asynchronous call
	client, err := r.getClient()
	if logus.Log.CheckWarn(err, "dialing:") {
		return err
	}

	divCall := client.Go("ServerRpc.GetInfo", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}
