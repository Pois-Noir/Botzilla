package botzillaclient

type UserHandler interface {
	Message(body string, sender string) (string, error)
	Broadcast(body string, sender string) error
}
