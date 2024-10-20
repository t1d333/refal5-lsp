package server


type StartMode int 

const (
	TCP StartMode = iota
	WEB_SOCKET
	STDIO
)

type StartSettings struct {
	Mode StartMode
	Addres string
}
