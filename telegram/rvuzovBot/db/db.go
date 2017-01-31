package db

import (
	"github.com/AlexeySpiridonov/goapp-config"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2"
	"strings"
)

const (
	universityDbName  = "university"
	lessonDbName      = "lesson"
	groupDbName       = "group"
	facultyDbName     = "faculty"
)

var (
	log = logging.MustGetLogger("db")

	context Context

	Universities   *mgo.Collection
	Lessons        *mgo.Collection
	DBGroups       *mgo.Collection
	Faculties      *mgo.Collection
)

type Context struct {
	Session *mgo.Session
	Db      *mgo.Database
}

func Get() Context {
	return context
}

func Init() (*mgo.Session, error) {
	log.Info("Connect to DB: " + config.Get("dbHost") + " " + config.Get("dbName"))
	mongo, err := mgo.Dial(config.Get("dbHost"))
	if err != nil {
		log.Panic("Cant't connect to mongoDB. Server is stopped")
	}
	log.Info("DB ok")
	set(mongo, mongo.DB(config.Get("dbName")))

	return mongo, err
}

func set(session *mgo.Session, db *mgo.Database) {
	context = Context{session, db}

	Universities = context.Db.C(universityDbName)
	Lessons = context.Db.C(lessonDbName)
	DBGroups = context.Db.C(groupDbName)
	Faculties = context.Db.C(facultyDbName)
}

func refresh(source string, err error) {
	if err == nil {
		return
	}

	if err.Error() == "not found" {
		log.Notice(source + " " + err.Error())
	} else {
		log.Error(source + " " + err.Error())
	}

	if err.Error() == "EOF" || strings.Contains(err.Error(), "connection reset by peer") {
		log.Warning("DB connect autoRefresh")
		context.Session.Refresh()
	}
}
