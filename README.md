# Botzilla Client

BotzillaClient is a library designed for interacting with Botzilla. To learn more about <a href="https://github.com/rima1881/Botzilla">Botzilla</a> itself, refer to its documentation. Note that BotzillaClient requires Botzilla to be actively running on the network to function properly.

# Functions

```
Start( name string , port string ) error
```
- This function is used to connect a component to the Botzilla server.
- Each component must invoke this function to establish a connection.
- The name parameter must be unique for each component.
- Opens a TCP server on the specified port.

<br />

```
SendCommnad( follower string , body string ) string
```
- Sends a command to a specified follower and returns a response.

<br />

```
SendMessage( followers string[] , body string )
```
- Sends a message to multiple followers specified in the followers list.

<br />

### Receiving Commands and Messages
To handle incoming commands or messages, you need to implement the following two interfaces:

`CommandListener`
```
type CommandListener interface {
	OnReceive() string
}
```

<br />

`MessageListener`
```
type MessageListener interface {
	OnReceive() string
}
```

### Setting Up Listeners
Once the interfaces are implemented, you can initialize the listeners using the following functions:

`Start Command Listener`
```
func StartCommandListener( listener CommandListener )
```

<br />

`Start Message Listener`
```
func StartMessageListener( listener MessageListener)
```

by setting up these listeners, your components can effectively handle incoming commands and messages.