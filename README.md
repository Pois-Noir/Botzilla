# Botzilla

`Botzilla` is a lightweight library designed for seamless communication between nodes, optimized for minimal overhead and simplicity.

# Network Architecture
Botzilla revolves around the concept of **components**. Each component registers itself on the local network, using a unique name.

# Comminucation Types
## Message
A typical point-to-point (P2P) message exchange involves a request and a response. For usage, refer to the example provided in `register.go`.

# Security
To secure incoming messages, **HMAC** (Hash-based Message Authentication Code) with a shared secret key is used. This ensures the integrity and authenticity of the messages.

**Note:** The messages themselves are not encrypted, meaning that while HMAC ensures authenticity, the content of the messages could still be exposed if intercepted. Encryption for message content is planned for future releases.

# Functions

```
NewComponent(name, secretKey) (Component, error)
```

**Parameters:**
- `name:` A unique name for the component.

- `secretKey:` The shared secret key known to all components and the server. **Do not send this via the protocol.**

<br />

```
Component.SendMessage(name string, body map[string]string) (map[string]string, error)
```
- `Description`: Sends a message to a specified component and waits for a response.

- `Parameters`:
	- `name`: The unique identifier of the component to which the message is being sent.

	- `body`: The content of the message to send.
- `Returns`:  A response from the recipient component and any potential errors.

<br />

```
GetComponents()
```
- `Returns`: A map with component name as key and component ip address as value


# Road Map
- Native Android support
- Realtime data