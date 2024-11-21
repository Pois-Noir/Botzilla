package botzillaclient

type CommandListener interface {
	OnReceive(body string) (string, error)
}

type MessageListener interface {
	OnReceive(body string)
}
