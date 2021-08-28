package tcp

type IPktProcessor interface {
	UnPack(session *Session) ([]byte, error)
	Pack(data []byte) []byte
}

type IPktProcessorMaker interface {
	CreatePktProcessor() IPktProcessor
}
