package repositories

import (
	"gopkg.in/mgo.v2"
)

type connectionChecker struct {
	Database *mgo.Database
}

func (cc *connectionChecker) Execute(logic func() error) error {
	if cc.Database.Session.Ping() != nil {
		cc.Database.Session.Refresh()
		if err := cc.Database.Session.Ping(); err != nil {
			return err
		}
	}
	return logic()
}
