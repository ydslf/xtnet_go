package frame

//service 可以创建 service， 与别的service通讯

type Service struct {
	loop       Loop
	rpc        string //service之间的rpc
	ServiceMgr string
}
