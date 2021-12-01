package frame

import (
	xtnetNet "xtnet/net"
)

type App struct {
	netServers map[int]xtnetNet.IServer
}
