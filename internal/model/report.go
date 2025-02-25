package model

type URL struct {
	URL    string
	Method string
	Params  []Param
}

type Report struct {
	Id   string
	URLs []URL
}

type Param struct {
	Name     string
	Values []string
	Patterns []string
}
