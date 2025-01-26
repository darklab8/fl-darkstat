package infocard

type NotParsedInfocard struct {
	content string
}

type Infocard struct {
	Lines []string
}

func NewInfocard(Lines []string) *Infocard {
	return &Infocard{Lines: Lines}
}

func NewNotParsedInfocard(Content string) *NotParsedInfocard {
	return &NotParsedInfocard{content: Content}
}

type Infoname string

type RecordKind string

const (
	TYPE_NAME    RecordKind = "NAME"
	TYPE_INFOCAD RecordKind = "INFOCARD"
)

type Config struct {
	Infonames         map[int]Infoname
	NotParsedInfocard map[int]*NotParsedInfocard
	Infocards         map[int]*Infocard
}

func NewConfig() *Config {
	return &Config{
		Infocards:         make(map[int]*Infocard),
		Infonames:         make(map[int]Infoname),
		NotParsedInfocard: make(map[int]*NotParsedInfocard),
	}
}
