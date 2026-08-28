package rpc

import (
	"errors"
	"time"
	"xtnet/net"
	"xtnet/net/packet"
)

/*
	|                      pktBody                        |
	|  msgID |             msgBody                        |
	| rpcType | contextID  | msgType | OP |   msgBody     |
*/

const (
	rtDirect   int8 = iota + 1 //直接，没有RPC
	rtRequest                  //RPC请求
	rtResponse                 //异步RPC回应
)

const maxContextID int32 = 0x7FFFFFFF

var ErrRequestTimeout = errors.New("rpc request timeout")

type OnRpcDirect func(session net.ISession, rpk *packet.ReadPacket)
type OnRpcRequest func(session net.ISession, contextID int32, rpk *packet.ReadPacket)
type RequestCallback func(rpk *packet.ReadPacket, err error)

type IRpc interface {
	SetOnRpcDirect(onRpcDirect OnRpcDirect)
	SetOnRpcRequest(nRpcRequest OnRpcRequest)
	HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket)
	WriteDirectHead(wpk *packet.WritePacket)
	SendDirect(session net.ISession, wpk *packet.WritePacket)
	SendDirectRaw(session net.ISession, wpk *packet.WritePacket)
	RequestAsync(session net.ISession, wpk *packet.WritePacket, expireMS time.Duration, cb RequestCallback)
	RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS time.Duration) (rpk *packet.ReadPacket, err error)
	Respond(session net.ISession, contextID int32, wpk *packet.WritePacket)
}
