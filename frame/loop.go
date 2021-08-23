package frame

type LoopFun func()

type Loop struct {
	loopFuns chan LoopFun
}

func NewLoop() *Loop {
	return &Loop{
		//TODO chan长度参数化，或设置为一个常量
		loopFuns: make(chan LoopFun, 10000),
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

func (loop *Loop) RunOnce() {
	f := <-loop.loopFuns
	//TODO 错误检测
	f()
}
