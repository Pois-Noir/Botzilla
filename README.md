# Botzilla

`Botzilla` is a lightweight library designed for seamless communication between nodes, optimized for minimal overhead and simplicity.

**Note:** To operate correctly, `Botzilla` requires the `Botzilla-Server` to be running on the network.

# Network Architecture
Botzilla revolves around the concept of **components**. Each component registers itself with the Botzilla server, using a unique name. Components can be grouped together, and multiple components can be instantiated within the same codebase.

The Botzilla server primarily serves as a discovery system; all messages are sent directly from one component to another without any intermediary, ensuring efficient and direct communication.

# Comminucation Types
## Message
A typical point-to-point (P2P) message exchange involves a request and a response. For usage, refer to the example provided in `register.go`.

## Broadcast
The broadcast function allows a message to be sent to all registered components on the network. If any components are not set up to listen for broadcasts, they will simply ignore the message.

# Security
To secure incoming messages, **HMAC** (Hash-based Message Authentication Code) with a shared secret key is used. This ensures the integrity and authenticity of the messages.

**Note:** The messages themselves are not encrypted, meaning that while HMAC ensures authenticity, the content of the messages could still be exposed if intercepted. Encryption for message content is planned for future releases.

# Functions

```
NewComponent(ServerAddr, secretKey, name, port, MessageListener) (Component, error)
```

**Parameters:**

- `ServerAddr:` A string representing the address where the Botzilla server is running (including the port), e.g., "localhost:8080".

- `secretKey:` The shared secret key known to all components and the server. **Do not send this via the protocol.**

- `name:` A unique name for the component.

- `port:` The TCP port to listen on.

- `MessageListener:` A function to handle incoming messages (see examples for how to define it).


<br />

```
SendMessage(name string, body map[string]string) (map[string]string, error)
```
- `Description`: Sends a message to a specified component and waits for a response.

- `Parameters`:
	- `name`: The unique identifier of the component to which the message is being sent.

	- `body`: The content of the message to send.
- `Returns`:  A response from the recipient component and any potential errors.


<br />

```
BroadcastMessage(message map[string]string) error
```
- `Description`: Sends a broadcast to all components registered in server.

- `Parameters`:
	- `body`: The body of the message.
- `Returns`: An error if the message could not be delivered to any of the followers.

<br />

```
GetComponents()
```
- `Returns`: An array of registered components by their names


# Road Map

- Streams / Realtime Comminucation
- C/C++/Java/C# bindings
- encrypted messages