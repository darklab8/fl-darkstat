package infocard

type Infocard struct {
	Content string
}

func NewInfocard(Content string) *Infocard {
	return &Infocard{Content: Content}
}

type Infoname string

type RecordKind string

const (
	TYPE_NAME    RecordKind = "NAME"
	TYPE_INFOCAD RecordKind = "INFOCARD"
)

type Config struct {
	Infonames map[int]Infoname
	Infocards map[int]*Infocard
}

func NewConfig() *Config {
	return &Config{
		Infocards: make(map[int]*Infocard),
		Infonames: make(map[int]Infoname),
	}
}
