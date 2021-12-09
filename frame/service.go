package frame

type Service struct {
	loop *Loop
}

func NewService(loopSize int) *Service {
	return &Service{
		loop: NewLoop(loopSize),
	}
}

func (service *Service) GetLoop() *Loop {
	return service.loop
}

func (service *Service) Run() {
	service.loop.Run()
}
