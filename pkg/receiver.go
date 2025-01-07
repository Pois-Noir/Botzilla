package botzilla

type UserHandler interface {
	Message(body map[string]string, sender string) (map[string]string, error)
}
