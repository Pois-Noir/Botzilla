package botzillaclient

type CommandListener interface {
	OnReceive() string
}

type MessageListener interface {
	OnReceive() string
}
