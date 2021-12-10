package frame

type Service struct {
	loop *Loop
}

func NewService(loopSize int, fullWarn bool) *Service {
	return &Service{
		loop: NewLoop(loopSize, fullWarn),
	}
}

func (service *Service) GetLoop() *Loop {
	return service.loop
}

func (service *Service) Run() {
	service.loop.Run()
}
