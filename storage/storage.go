package storage

const (
	OpAdd = iota
	OpDel = iota
	OpMod = iota
)

type Param struct {
	Name         string `json:"name"`
	Source       string `json:"source"`
	DefaultValue string `json:"default_value"`
	Required     bool   `json:"required"`
}

type APIChildren struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	Method  string   `json:"method"`
	Params  []*Param `json:"params"`
	Timeout int      `json:"timeout"`
}

type API struct {
	Url      string         `json:"url"`
	Method   string         `json:"method"`
	ExecType string         `json:"execType"`
	Children []*APIChildren `json:"children"`
}

type Project struct {
	Name string `json:"name"`
	APIs []*API `json:"apis"`
}

type Event struct {
	Op      int
	Project *Project
}

type IStorage interface {
	Init() error
	Get(project string) (p *Project, err error)
	GetAll() (projects map[string]*Project, err error)
	WatchEvent() (ch chan Event)
	Watch() (err error)
}

type Storage struct {
	IStorage
}

func NewStorage(o IStorage) *Storage {
	return &Storage{
		o,
	}
}
