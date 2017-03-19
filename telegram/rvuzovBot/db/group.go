package db

import (
	"gopkg.in/mgo.v2/bson"
	"unicode/utf8"
)

type (
	GroupUsers struct {
		Allow []string `json:"allow" bson:"allow"`
		Deny  []string `json:"deny" bson:"deny"`
	}

	Group struct {
		Id           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
		Name         string        `json:"name"`
		Lessons      string        `json:"sсhedule" bson:"lessons,omitempty"`
		University   string        `json:"university"`
		Faculty      string        `json:"faculty"`
		Updateat     int           `json:"updateat"`
		Official     bool          `json:"official"`
		Parity       bool          `json:"parity"`
		Elder        string        `json:"elder"`
		ElderBans    []string      `json:"elder_bans" bson:"elder_bans"`
		Access       string        `json:"access" bson:"access"`
		LessonsEdit  string        `json:"lessons_edit" bson:"lessons_edit"`
		JournalsEdit string        `json:"journals_edit" bson:"journals_edit"`
		Users        GroupUsers    `json:"users" bson:"users"`
	}

	groups struct{}
)

func (g Group) GetID() string {
	return g.Id.Hex()
}

func (g Group) GetName() string {
	return g.Name
}

func GroupSearch(id, query string) (groups []Group, err error) {
	groups = make([]Group, 0)
	log.Debug("Search group:" + query + " in faculty:" + id)
	var find interface{}
	if query == "" {
		find = bson.M{"faculty": id}
	} else if utf8.RuneCountInString(query) < 3 {
		find = bson.M{"faculty": id, "name": &bson.RegEx{Pattern: "^" + query, Options: "si"}}
	} else {
		find = bson.M{"faculty": id, "name": &bson.RegEx{Pattern: query, Options: "si"}}
	}
	err = DBGroups.Find(find).Limit(50).Sort("name").All(&groups)
	refresh("group", err)
	return
}


