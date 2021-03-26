package frame

type Rpc struct {
	loop *Loop
	ret  chan func()
}
