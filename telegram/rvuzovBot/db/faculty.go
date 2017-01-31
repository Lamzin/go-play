package db

import (
	"gopkg.in/mgo.v2/bson"
	"unicode/utf8"
)

type Faculty struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name"`
	University string        `json:"university"`
	Updateat   int           `json:"updateat"`
	Official   bool          `json:"official"`
}

func (f Faculty) GetID() string {
	return f.Id.Hex()
}

func (f Faculty) GetName() string {
	return f.Name
}

func FacultySearch(id, query string) (faculties []Faculty, err error) {
	faculties = make([]Faculty, 0)
	log.Debug("Search faculty:" + query + " in unversity:" + id)
	if query == "" {
		find := bson.M{"university": id}
		err = Faculties.Find(find).Sort("name").Limit(100).All(&faculties)
	} else if utf8.RuneCountInString(query) < 3 {
		find := bson.M{"university": id, "name": &bson.RegEx{Pattern: "^" + query, Options: "si"}}
		err = Faculties.Find(find).Sort("name").Limit(20).All(&faculties)
	} else {
		find := bson.M{"university": id, "name": &bson.RegEx{Pattern: query, Options: "i"}}
		err = Faculties.Find(find).Sort("name").Limit(20).All(&faculties)
	}
	refresh("faculty", err)
	return
}
