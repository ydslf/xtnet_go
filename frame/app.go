package frame

import (
	xtNet "xtnet/net"
)

type App struct {
	netServers map[int]xtNet.IServer
}
