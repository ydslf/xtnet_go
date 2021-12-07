package frame

type Service struct {
	loop *Loop
}

func NewService() *Service {
	return &Service{
		loop: NewLoop(),
	}
}

func (service *Service) GetLoop() *Loop {
	return service.loop
}

func (service *Service) Run() {
	service.loop.Run()
}
