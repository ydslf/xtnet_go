package frame

//service 可以创建 service， 与别的service通讯

type Service struct {
	loop *Loop
	Rpc  string
}

func (service *Service) Run() {
	service.loop.Run()
}

func (service *Service) RegService() {

}
