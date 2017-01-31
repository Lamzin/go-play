package db

type StandartShortGroupCard struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type StandartShortFacultyCard struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type StandartShortUserCard struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Access string `json:"access"`
	Url    string `json:"url"`
}

type StandartShortUniversityCard struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}
