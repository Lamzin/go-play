package db

import (
	"gopkg.in/mgo.v2/bson"
	"unicode/utf8"
	"fmt"
)

type University struct {
	Id              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string        `json:"name"`
	Abbr            string        `json:"abbr"`
	Token           string        `json:"token" bson:"token,omitempty"`
	Email           string        `json:"email"`
	Updateat        int           `json:"updateat"`
	Fcount          int           `json:"fcount" bson:"fcount,omitempty"`
	Gcount          int           `json:"gcount" bson:"gcount,omitempty"`
	Ucount          int           `json:"ucount" bson:"ucount,omitempty"`
	DataUrl         string        `json:"dataurl"`
	IntegrationType string        `json:"integrationtype"`
	Official        bool          `json:"official"`
	Whitelabel      bool          `json:"whitelabel"`
}

func (u University) GetID() string {
	return u.Id.Hex()
}

func (u University) GetName() string {
	return fmt.Sprintf("%s, %s", u.Name, u.Abbr)
}

func UniversityList() (universities []University, err error) {
	universities = make([]University, 0)
	err = Universities.Find(nil).Limit(50).Sort("ucount").All(&universities)
	refresh("university", err)
	return
}

func UniversitySearch(query string) (universities []University, err error) {
	universities = make([]University, 0)
	if utf8.RuneCountInString(query) < 3 {
		find := bson.M{
			"gcount": bson.M{"$gt": 0},
			"$or": []interface{}{
				bson.M{"name": &bson.RegEx{Pattern: "^" + query, Options: "si"}},
				bson.M{"abbr": &bson.RegEx{Pattern: "^" + query, Options: "si"}}},
			"whitelabel": bson.M{"$ne": true}}
		err = Universities.Find(find).Sort("name").Limit(20).All(&universities)
	} else {
		find := bson.M{
			"gcount": bson.M{"$gt": 0},
			"$or": []interface{}{
				bson.M{"name": &bson.RegEx{Pattern: query, Options: "i"}},
				bson.M{"abbr": &bson.RegEx{Pattern: query, Options: "i"}}},
			"whitelabel": bson.M{"$ne": true}}
		err = Universities.Find(find).Sort("name").Limit(20).All(&universities)
	}
	refresh("university", err)
	return
}
