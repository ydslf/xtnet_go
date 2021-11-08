package frame

type LoopFun func()

type Loop struct {
	functions chan LoopFun
}

func NewLoop() *Loop {
	return &Loop{
		//TODO chan长度参数化，或设置为一个常量
		functions: make(chan LoopFun, 10000),
	}
}

func (loop *Loop) Post(f LoopFun) {
	//TODO 判断是否满了
	loop.functions <- f
}

func (loop *Loop) Run() {
	for {
		select {
		case f := <-loop.functions:
			f()
		}
	}
}

func (loop *Loop) RunOnce() {
	f := <-loop.functions
	//TODO 错误检测
	f()
}
