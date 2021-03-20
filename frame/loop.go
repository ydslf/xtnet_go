package frame

type LoopFun func()

type Loop struct {
	loopFuns chan LoopFun
}

func NewLoop() *Loop {
	return &Loop{
		loopFuns: make(chan LoopFun, 100),
	}
}

func (loop *Loop) Post(f LoopFun) {
	loop.loopFuns <- f
}

func (loop *Loop) Run() {
	for {
		select {
		case f := <-loop.loopFuns:
			f()
		}
	}
}
