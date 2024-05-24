package config_consts

type Relationship string

const (
	RepFriend  Relationship = "friend"
	RepNeutral Relationship = "neutral"
	RepEnemy   Relationship = "enemy"
)

func (rep Relationship) ToStr() string {
	return string(rep)
}

func GetRelationshipStatus(rep_value float64) Relationship {
	switch rep := rep_value; {
	case rep > 0.6:
		return RepFriend
	case rep < -0.6:
		return RepEnemy
	default:
		return RepNeutral
	}
}
