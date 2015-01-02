package sys

import "github.com/lunny/tango"

var Tango *tango.Tango

func newTango() *tango.Tango {
	return tango.Classic()
}

func InitTango() {
	Tango = newTango()
}
