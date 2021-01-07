package frame

type LoopFun func()

type Loop struct {
	loopFuns chan LoopFun
}

func NewLoop() *Loop {
	return &Loop{}
}

func (loop *Loop) Post(f LoopFun) {
	loop.loopFuns <- f
}

func (loop *Loop) Start() {
	for f := range loop.loopFuns {
		f()
	}
}
