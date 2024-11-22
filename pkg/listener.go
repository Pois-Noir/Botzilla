package botzillaclient

type Listener interface {
	Command(body string) (string, error)
	Message(body string) error
	Stream() // Not sure yet how :(
}

type Config struct {
	Name        string
	CommandPort int
	MessagePort int
	StreamPort  int
}
