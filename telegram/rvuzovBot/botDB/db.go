package botDB

import (
	"github.com/AlexeySpiridonov/goapp-config"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2"
	"strings"
)

var (
	log = logging.MustGetLogger("botDB")

	ctx context
)

type context struct {
	Session *mgo.Session
	DB      *mgo.Database

	Chats *mgo.Collection
}

func (c *context) Setup(session *mgo.Session) {
	c.Session = session
	c.DB = session.DB(config.Get("dbNameBot"))
	c.Chats = c.DB.C("chats")
}

func (c *context) Refresh(source string, err error) {
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
		c.Session.Refresh()
	}
}

func Init() (*mgo.Session, error) {
	log.Info("Connect to DB: " + config.Get("dbHostBot") + " " + config.Get("dbNameBot"))
	session, err := mgo.Dial(config.Get("dbHostBot"))
	if err != nil {
		log.Panic("Cant't connect to mongoDB. Server is stopped")
	}
	log.Info("DB ok")
	ctx.Setup(session)
	return session, err
}
