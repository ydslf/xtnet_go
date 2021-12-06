package rpc

import (
	"xtnet/net"
	"xtnet/net/packet"
)

/*
	|                      pktBody                     |
	|  msgID |             msgBody                     |
	| rpcType | contextID  |       msgBody             |
*/

const (
	rtDirect   int8 = iota //直接，没有RPC
	rtRequest              //RPC请求
	rtResponse             //异步RPC回应
)

const maxContextID int32 = 0x8FFFFFFF

type OnRpcDirect func(session net.ISession, rpk *packet.ReadPacket)
type OnRpcRequest func(session net.ISession, contextID int32, rpk *packet.ReadPacket)
type RequestCallback func(rpk *packet.ReadPacket)

type IRpc interface {
	HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket)
	SendDirect(session net.ISession, wpk *packet.WritePacket)
	RequestAsync(session net.ISession, wpk *packet.WritePacket, cb RequestCallback)
	RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS int) (rpk *packet.ReadPacket, err error)
	Respond(session net.ISession, contextID int32, wpk *packet.WritePacket)
}
