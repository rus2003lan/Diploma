package report


type url struct {
	URL    string   `json:"url"`
	Method string   `json:"method"`
	Params []param  `json:"params"`
}

type report struct {
	Id   string `json:"id"`
	URLs []url  `json:"urls"`
}

type param struct {
	Name     string   `json:"name"`
	Values   []string `json:"values"`
	Patterns []string `json:"patterns"`
}