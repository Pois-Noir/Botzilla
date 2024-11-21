package botzillaclient

func Start(name string, port string) error {
	return nil
}

func SendCommand(follower string, body string) (string, error) {
	return "", nil
}

func BroadcastMessage(followers []string, body string) error {
	return nil
}

func StartCommandListener(listener CommandListener) error {
	return nil
}

func StartMessageListener(listener MessageListener) error {
	return nil
}
