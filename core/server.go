package core

import "github.com/lunny/tango"

type CoreServer struct {
	*tango.Tango
}

func NewCoreServer() *CoreServer {
	server := tango.Classic()
	return &CoreServer{server}
}
