package db

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type (
	Lesson struct {
		Id         bson.ObjectId               `json:"id,omitempty" bson:"_id,omitempty"`
		Subject    string                      `json:"subject"`
		Type       string                      `json:"type"`
		Time       Time                        `json:"time"`
		Date       string                      `json:"date"`
		Period     DateObj                     `json:"period,omitempty" bson:"period,omitempty"`
		Audiences  []Audience                  `json:"audiences"`
		Teachers   []Teacher                   `json:"teachers"`
		Subgroups  string                      `json:"subgroups"`
		Note       string                      `json:"note"`
		Group      StandartShortGroupCard      `json:"group"`
		University StandartShortUniversityCard `json:"university"`
		Faculty    StandartShortFacultyCard    `json:"faculty"`
		Skip       string                      `json:"skip"`
		Official   bool                        `json:"official"`
	}

	Time struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	DateStr string

	DateObj struct {
		Start   string `json:"start,omitempty"`
		End     string `json:"end,omitempty"`
		Weekday int    `json:"weekday,omitempty"`
		Week    int    `json:"week,omitempty"`
	}

	Audience struct {
		Name   string `json:"name"`
		Addr   string `json:"addr,omitempty"`
		LonLat string `json:"lonlat,omitempty"`
	}

	Teacher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func groupLessons(group string) (s string, err error) {
	result := make([]Lesson, 0)
	if err = Lessons.Find(bson.M{"group.id": group}).All(&result); err == nil {
		b, err1 := json.Marshal(result)
		s, err = string(b), err1
	}
	refresh("lesson", err)
	return
}

func GetLessonsByDate(groupID string, date string) (lessons []Lesson, err error) {
	if err = Lessons.Find(bson.M{"group.id": groupID}).All(&lessons); err != nil {
		return
	}

	for i := 0; i < len(lessons); i++ {
		if !strings.Contains(lessons[i].Date, date) {
			lessons = append(lessons[:i], lessons[i + 1:]...)
			i--
		}
	}
	return
}
