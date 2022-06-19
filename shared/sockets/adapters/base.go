package adapters

type MessageHandler = func(message []byte) bool

type SocketAdapter interface {
	Close() error

	// Boolean indicating whether the error was unexpected (true) or not (false)
	Serve() (error, bool)

	WriteControlMessage(message []byte) error
	SetControlMessageHandler(handler MessageHandler)

	WriteDataMessage(message []byte) error
	SetDataMessageHandler(handler MessageHandler)

	WritePingMessage() error
	SetPongHandler(handler func())
}

type socketBase struct {
	controlMessageHandler MessageHandler
	dataMessageHandler    MessageHandler
	pongHandler           func()
}

func (s *socketBase) SetControlMessageHandler(handler MessageHandler) {
	s.controlMessageHandler = handler
}

func (s *socketBase) SetDataMessageHandler(handler MessageHandler) {
	s.dataMessageHandler = handler
}

func (s *socketBase) SetPongHandler(handler func()) {
	s.pongHandler = handler
}