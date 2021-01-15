package frame

type LoopFun func()

type Loop struct {
	loopFuns chan LoopFun
}

func NewLoop() *Loop {
	return &Loop{
		loopFuns: make(chan LoopFun),
	}
}

func (loop *Loop) Post(f LoopFun) {
	loop.loopFuns <- f
}

func (loop *Loop) Run() {
	for f := range loop.loopFuns {
		f()
	}
}
